package main

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"
)

/*const (
    Aria2GetFileInfo   string = "GetFileInfo"
    Aria2GetGlobalInfo string = "GetGlobalInfo"
    Aria2PauseTask     string = "PauseTask"
    Aria2StopTask      string = "StopTask"
    Aria2UnpauseTask   string = "UnpauseTask"

)*/

const (
	HTTPDownload     Method = "aria2.addUri" //Including Magnet Link
	TorrentDownload  Method = "aria2.addTorrent"
	MetalinkDownload Method = "aria2.addMetalink"
)

type Aria2Tasks struct {
	Gids   []string
	Status map[string]bool
}

var isAria2Running = false
var workQueue map[string]*chan Aria2WebSocketEvent
var workLock sync.RWMutex

var subscriptionQueue map[Method][]*chan Aria2WebSocketGIDInfo
var subLock sync.RWMutex

var TaskStucked []string
var TaskPaused []string
var totalActiveAllowed int

func init() {
	workQueue = make(map[string]*chan Aria2WebSocketEvent)
	subscriptionQueue = make(map[Method][]*chan Aria2WebSocketGIDInfo)
}

func SetAria2Listener() {

	infoChan := make(chan Aria2WebSocketEvent)
	actionChan := make(chan Aria2WebSocketEvent)

	Aria2Controller = Aria2Instance{
		InfoChan:   &infoChan,
		ActionChan: &actionChan,
	}

	u := url.URL{
		Scheme: "ws",
		Host:   Conf.GetString(`aria2.server.host`),
		Path:   Conf.GetString(`aria2.server.path`),
	}

	go Aria2Controller.StartAria2Websocket(&u)
	go dispatchAria2Event()

	isAria2Running = true
	updatePausedData()

	for {
		autoSorting()
		time.Sleep(Conf.GetDuration("aria2.server.autoPause") * time.Second)
	}
}

type DownloadOption func(*Aria2DownloadTask)

func WithHeader(n, v string) DownloadOption {
	return func(t *Aria2DownloadTask) {
		t.Headers = append(t.Headers, StringBuilder(n, ":", v))
	}
}

func WithOutName(o string) DownloadOption {
	return func(t *Aria2DownloadTask) {
		t.OutName = o
	}
}

func WithDestination(d string) DownloadOption {
	return func(t *Aria2DownloadTask) {
		t.Destination = d
	}
}

func WithTaskType(t Method) DownloadOption {
	return func(task *Aria2DownloadTask) {
		task.Type = t
	}
}

func NewDownloadTask(source string, options ...DownloadOption) *Aria2DownloadTask {
	Default := &Aria2DownloadTask{
		Destination: Conf.GetString("aria2.server.tempPath"),
		OutName:     "tmpfile",
		Type:        HTTPDownload,
	}

	Default.Source = source

	for i := range options {
		options[i](Default)
	}

	return Default
}

func (t *Aria2DownloadTask) DoDownload(token string) (string, error) {
	if token == `` {
		token = genToken()
	}
	pipe := make(chan Aria2WebSocketEvent)

	workLock.Lock()
	workQueue[token] = &pipe
	workLock.Unlock()

	e := t.Download(token)
	*Aria2Controller.ActionChan <- e

	select {
	case info := <-pipe:
		close(pipe)
		if info.Error.Code != 0 {
			return "", errors.New(info.Error.Message)
		}
		var r string
		err := jsoniter.Unmarshal(*info.Result, &r)
		return r, err
	}
}

func NewAria2TaskGIDList(gid ...string) *Aria2Tasks {
	status := map[string]bool{}
	for i := range gid {
		if _, ok := status[gid[i]]; !ok {
			status[gid[i]] = false
		}
	}
	return &Aria2Tasks{Gids: gid, Status: status}
}

func (Tasks Aria2Tasks) GetGlobalStat() (*Aria2GlobalStat, error) {
	if !isAria2Running {
		return nil, errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)

	checkStatus(&pipe)

	var err string
	select {
	case info := <-pipe:
		if info.Error.Code == 0 {
			var r Aria2GlobalStat
			err := jsoniter.Unmarshal(*info.Result, &r)
			return &r, err
		} else {
			err = info.Error.Message
		}

		close(pipe)
		break
	}

	return nil, errors.New(err)
}

func (Tasks Aria2Tasks) GetFilesInfo() ([][]Aria2FileInfo, error) {
	if !isAria2Running {
		return nil, errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)
	if len(Tasks.Gids) == 0 {
		return nil, errors.New("task gid required")
	}

	for i := range Tasks.Gids {
		go checkStatus(&pipe, Tasks.Gids[i])
	}

	var infos [][]Aria2FileInfo
	for {
		select {
		case info := <-pipe:
			var i []Aria2FileInfo
			err := jsoniter.Unmarshal(*info.Result, i)
			if err == nil {
				infos = append(infos, i)
			} else {
				log.Println(i, info.Result, err)
			}

			if len(infos) == len(Tasks.Gids) {
				close(pipe)
				return infos, nil
			}
		}
	}

}

func (Tasks Aria2Tasks) Pause(isForce bool) error {
	if !isAria2Running {
		return errors.New("aria2 is not running")
	}

	applyAll := len(Tasks.Gids) == 0
	pipe := make(chan Aria2WebSocketEvent)

	pauseTask(isForce, applyAll, &pipe, Tasks.Gids...)

	var (
		err string
		num int
	)

	for {
		select {
		case info := <-pipe:
			if info.Error.Code != 0 {
				err += info.Error.Message
			}

			num++

			if num >= len(Tasks.Gids) {
				close(pipe)
				if err != "" {
					return errors.New(err)
				}

				return nil
			}
		}
	}
}

func (Tasks Aria2Tasks) Remove(isForce bool) error {
	if !isAria2Running {
		return errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)

	removeTask(isForce, &pipe, Tasks.Gids...)

	var (
		err string
		num int
	)

	for {
		select {
		case info := <-pipe:
			if info.Error.Code != 0 {
				err += info.Error.Message
			}

			num++

			if num >= len(Tasks.Gids) {
				close(pipe)
				if err != "" {
					return errors.New(err)
				}

				return nil
			}
		}
	}

}

func (Tasks Aria2Tasks) RemoveDownloadResult() error {
	if !isAria2Running {
		return errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)

	removeDownloadResult(&pipe, Tasks.Gids...)

	var (
		err string
		num int
	)

	for {
		select {
		case info := <-pipe:
			if info.Error.Code != 0 {
				err += info.Error.Message
			}

			num++

			if num >= len(Tasks.Gids) {
				close(pipe)
				if err != "" {
					return errors.New(err)
				}

				return nil
			}
		}
	}

}

func (Tasks Aria2Tasks) Unpause() error {
	if !isAria2Running {
		return errors.New("aria2 is not running")
	}

	applyAll := len(Tasks.Gids) == 0
	pipe := make(chan Aria2WebSocketEvent)

	unpauseTask(applyAll, &pipe, Tasks.Gids...)

	var (
		err string
		num int
	)

	for {
		select {
		case info := <-pipe:
			if info.Error.Code != 0 {
				err += info.Error.Message
			}

			num++

			if num >= len(Tasks.Gids) {
				close(pipe)
				if err != "" {
					return errors.New(err)
				}

				return nil
			}
		}
	}

}

func (Tasks Aria2Tasks) TellStatus(options ...string) ([]Aria2TaskStatus, error) {
	if !isAria2Running {
		return nil, errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)
	var num int
	for i := range Tasks.Gids {
		go tellTaskStatus(&pipe, Tasks.Gids[i], options...)
	}

	var infos []Aria2TaskStatus
	var err string
	for {
		select {
		case info := <-pipe:
			if info.Error.Code == 0 {
				var i Aria2TaskStatus
				e := jsoniter.Unmarshal(*info.Result, &i)
				if e == nil {
					infos = append(infos, i)
				} else {
					err += e.Error()
				}
			} else {
				err += info.Error.Message
			}

			num++

			if num >= len(Tasks.Gids) {
				close(pipe)
				if err != "" {
					return infos, errors.New(err)
				}

				return infos, nil
			}
		}
	}
}

func (Tasks Aria2Tasks) TellActive(options ...string) ([]Aria2TaskStatus, error) {
	if !isAria2Running {
		return nil, errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)

	tellActiveTask(&pipe, options...)

	var err string
	select {
	case info := <-pipe:
		if info.Error.Code == 0 {
			var i []Aria2TaskStatus
			e := jsoniter.Unmarshal(*info.Result, &i)
			return i, e
		} else {
			err = info.Error.Message
		}
		close(pipe)
		return nil, errors.New(err)
	}

}

func (Tasks Aria2Tasks) TellWaiting(offset, amount int, options ...string) ([]Aria2TaskStatus, error) {
	if !isAria2Running {
		return nil, errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)

	tellWaitingTask(offset, amount, &pipe, options...)

	var err string
	select {
	case info := <-pipe:
		if info.Error.Code == 0 {
			var i []Aria2TaskStatus
			e := jsoniter.Unmarshal(*info.Result, &i)
			return i, e
		} else {
			err = info.Error.Message
		}
		close(pipe)
		return nil, errors.New(err)
	}
}

func (Tasks Aria2Tasks) TellStopped(offset, amount int, options ...string) ([]Aria2TaskStatus, error) {
	if !isAria2Running {
		return nil, errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)

	tellStoppedTask(offset, amount, &pipe, options...)

	var err string
	select {
	case info := <-pipe:
		if info.Error.Code == 0 {
			var i []Aria2TaskStatus
			e := jsoniter.Unmarshal(*info.Result, &i)
			return i, e
		} else {
			err = info.Error.Message
		}
		close(pipe)
		return nil, errors.New(err)
	}
}

func (Tasks Aria2Tasks) ChangePosition(action string, pos int) error {
	if !isAria2Running {
		return errors.New("aria2 is not running")
	}
	var isSetGroup = len(Tasks.Gids) > 1 && action == "POS_SET"

	pipe := make(chan Aria2WebSocketEvent)

	for i := range Tasks.Gids {
		go changePosition(Tasks.Gids[i], action, pos, &pipe)
		if isSetGroup {
			pos++
		}
	}

	var num int
	var err string
	for {
		select {
		case info := <-pipe:
			if info.Error.Code != 0 {
				err += info.Error.Message
			}

			num++

			if num >= len(Tasks.Gids) {
				close(pipe)
				if err != "" {
					return errors.New(err)
				}

				return nil
			}
		}
	}

}

func (Tasks Aria2Tasks) GetTaskOptions() ([]Aria2OptionInfo, error) {
	if !isAria2Running {
		return nil, errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)
	var num int
	for i := range Tasks.Gids {
		go getAria2TaskSettings(Tasks.Gids[i], &pipe)
	}

	var infos []Aria2OptionInfo
	var err string
	for {
		select {
		case info := <-pipe:
			if info.Error.Code == 0 {
				var i Aria2OptionInfo
				e := jsoniter.Unmarshal(*info.Result, &i)
				if e == nil {
					infos = append(infos, i)
				} else {
					err += e.Error()
				}
			} else {
				err += info.Error.Message
			}

			num++

			if num >= len(Tasks.Gids) {
				close(pipe)
				if err != "" {
					return infos, errors.New(err)
				}

				return infos, nil
			}
		}
	}

}

func (Tasks Aria2Tasks) ChangeTaskOptions(options map[string]string) error {
	if !isAria2Running {
		return errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)
	var num int
	for i := range Tasks.Gids {
		go changeAria2TaskSettings(Tasks.Gids[i], &pipe, options)
	}

	var err string
	for {
		select {
		case info := <-pipe:
			if info.Error.Code != 0 {
				err += info.Error.Message
			}

			num++

			if num >= len(Tasks.Gids) {
				close(pipe)
				if err != "" {
					return errors.New(err)
				}

				return nil
			}
		}
	}

}

func (Tasks Aria2Tasks) GetGlobalOptions() (*Aria2GlobalOptionInfo, error) {
	if !isAria2Running {
		return nil, errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)

	getAria2Settings(&pipe)

	var err string
	select {
	case info := <-pipe:
		if info.Error.Code == 0 {
			var i Aria2GlobalOptionInfo
			e := jsoniter.Unmarshal(*info.Result, &i)
			return &i, e
		} else {
			err = info.Error.Message
		}
		close(pipe)
		return nil, errors.New(err)
	}

}

func (Tasks Aria2Tasks) ChangeGlobalOptions(options map[string]string) error {
	if !isAria2Running {
		return errors.New("aria2 is not running")
	}

	pipe := make(chan Aria2WebSocketEvent)

	changeAria2Settings(&pipe, options)

	var err string
	select {
	case info := <-pipe:
		if info.Error.Code != 0 {
			err = info.Error.Message
		} else {
			return nil
		}

		close(pipe)
		return errors.New(err)
	}

}

func (Tasks Aria2Tasks) CheckStatus() string {
	if !isAria2Running {
		return "aria2 is not running"
	}

	var mutex sync.Mutex

	for i := range Tasks.Gids {
		Tasks.Status[Tasks.Gids[i]] = false
	}

	i := 0
	errStr := ""

	Tasks.taskStatusLookup(&i, &errStr, &mutex)

	return errStr
}

func (Tasks Aria2Tasks) WaitUntilFinished() string {
	if !isAria2Running {
		return "aria2 is not running"
	}
	pipe := make(chan Aria2WebSocketGIDInfo, len(Tasks.Gids)+1)

	var mutex sync.Mutex
	var WG sync.WaitGroup

	for i := range Tasks.Gids {
		Tasks.Status[Tasks.Gids[i]] = false
	}

	i := 0
	errStr := ""

	time.Sleep(1 * time.Second)

	Tasks.taskStatusLookup(&i, &errStr, &mutex)

	if i == len(Tasks.Gids) {
		goto Over
	}

	subLock.Lock()
	subscriptionQueue[DownloadComplete] = append(subscriptionQueue[DownloadComplete], &pipe)
	subscriptionQueue[BtDownloadComplete] = append(subscriptionQueue[BtDownloadComplete], &pipe)
	subLock.Unlock()

	WG.Add(1)
	go func() {
		for {
			select {
			case e := <-pipe:
				if _, ok := Tasks.Status[e.Gid]; ok {
					mutex.Lock()
					i++
					Tasks.Status[e.Gid] = true
					mutex.Unlock()
				}

				if i == len(Tasks.Gids) {
					WG.Done()
					return
				}
			case <-time.After(10 * time.Minute):
				Tasks.taskStatusLookup(&i, &errStr, &mutex)

				if errStr != "" {
					WG.Done()
					return
				}
			}
		}

	}()
	WG.Wait()
Over:
	go deletePipeInSubscriptionQueue(&pipe, BtDownloadComplete, DownloadComplete)

	if errStr != "" {
		return errStr
	}

	return ``

}

func (Tasks Aria2Tasks) taskStatusLookup(index *int, errStr *string, mutex *sync.Mutex) {
	status, err := Tasks.TellStatus()
	if err != nil {
		*errStr += err.Error()
		return
	}

	if status == nil {
		*errStr += `empty status`
		return
	}

	for j := range status {
		if _, ok := Tasks.Status[(status)[j].Gid]; ok {
			switch (status)[j].Status {
			case "active":
				if (status)[j].CompletedLength == (status)[j].TotalLength && (status)[j].TotalLength != `0` {
					mutex.Lock()
					*index++
					Tasks.Status[(status)[j].Gid] = true
					mutex.Unlock()
				}
			case "waiting":
				mutex.Lock()
				Tasks.Status[(status)[j].Gid] = false
				mutex.Unlock()
			case "paused":
				mutex.Lock()
				*index++
				Tasks.Status[(status)[j].Gid] = false
				*errStr += (status)[j].Gid + "Task Paused\n"
				mutex.Unlock()
			case "complete":
				mutex.Lock()
				*index++
				Tasks.Status[(status)[j].Gid] = true
				mutex.Unlock()
			case "removed":
				mutex.Lock()
				*index++
				Tasks.Status[(status)[j].Gid] = false
				*errStr += (status)[j].Gid + "Task Removed\n"
				mutex.Unlock()
			case "error":
				mutex.Lock()
				*index++
				Tasks.Status[(status)[j].Gid] = false
				*errStr += (status)[j].Gid + "Task Error\n"
				mutex.Unlock()
			}
		}
	}
}

func (Tasks Aria2Tasks) IsFinished() (map[string]bool, error) {
	if !isAria2Running {
		return nil, errors.New("aria2 is not running")
	}
	pipe := make(chan Aria2WebSocketGIDInfo)

	var mutex sync.Mutex

	for i := range Tasks.Gids {
		Tasks.Status[Tasks.Gids[i]] = false
	}

	i := 0
	errStr := ""

	Tasks.taskStatusLookup(&i, &errStr, &mutex)

	subLock.Lock()
	subscriptionQueue[DownloadComplete] = append(subscriptionQueue[DownloadComplete], &pipe)
	subscriptionQueue[BtDownloadComplete] = append(subscriptionQueue[BtDownloadComplete], &pipe)
	subLock.Unlock()

	go func() {
		for {
			select {
			case e := <-pipe:
				if _, ok := Tasks.Status[e.Gid]; ok {
					mutex.Lock()
					i++
					Tasks.Status[e.Gid] = true
					mutex.Unlock()
				}

				if i == len(Tasks.Gids) {
					return
				}
			case <-time.After(1 * time.Second):
				Tasks.taskStatusLookup(&i, &errStr, &mutex)
				return
			}
		}
	}()

	go deletePipeInSubscriptionQueue(&pipe, DownloadComplete, BtDownloadComplete)

	if errStr != "" {
		return Tasks.Status, errors.New(errStr)
	}

	return Tasks.Status, nil
}

func checkStatus(pipe *chan Aria2WebSocketEvent, gid ...string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	var e Aria2WebSocketEvent

	if len(gid) > 0 {
		for i := range gid {
			e = Aria2GetFiles(gid[i], t)
			*Aria2Controller.ActionChan <- e
		}
	} else {
		e = Aria2GetGlobalStat(t)
		*Aria2Controller.ActionChan <- e
	}
}

func pauseTask(isForce, applyAll bool, pipe *chan Aria2WebSocketEvent, gid ...string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	var e Aria2WebSocketEvent
	if applyAll {
		if isForce {
			e = Aria2ForcePauseAllTask(t)
		} else {
			e = Aria2PauseAllTask(t)
		}
		*Aria2Controller.ActionChan <- e
	} else {
		if isForce {
			for i := range gid {
				e = Aria2ForcePauseTask(gid[i], t)
				*Aria2Controller.ActionChan <- e
			}
		} else {
			for i := range gid {
				e = Aria2PauseTask(gid[i], t)
				*Aria2Controller.ActionChan <- e
			}
		}
	}
}

func removeTask(isForce bool, pipe *chan Aria2WebSocketEvent, gid ...string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	var e Aria2WebSocketEvent
	if isForce {
		for i := range gid {
			e = Aria2ForceRemoveTask(gid[i], t)
			*Aria2Controller.ActionChan <- e
		}
	} else {
		for i := range gid {
			e = Aria2RemoveTask(gid[i], t)
			*Aria2Controller.ActionChan <- e
		}
	}
}

func removeDownloadResult(pipe *chan Aria2WebSocketEvent, gid ...string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	for i := range gid {
		e := Aria2RemoveDownloadResult(gid[i], t)
		*Aria2Controller.ActionChan <- e
	}
}

func unpauseTask(applyAll bool, pipe *chan Aria2WebSocketEvent, gid ...string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	var e Aria2WebSocketEvent
	if applyAll {
		e = Aria2UnpauseAllTask(t)
		*Aria2Controller.ActionChan <- e
	} else {
		for i := range gid {
			e = Aria2UnpauseTask(gid[i], t)
			*Aria2Controller.ActionChan <- e
		}
	}
}

func tellTaskStatus(pipe *chan Aria2WebSocketEvent, gid string, options ...string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	e := Aria2TellStatusTasks(gid, t, options...)

	*Aria2Controller.ActionChan <- e
}

func tellActiveTask(pipe *chan Aria2WebSocketEvent, options ...string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	e := Aria2TellActiveTasks(t, options...)

	*Aria2Controller.ActionChan <- e
}

func tellWaitingTask(offset, amount int, pipe *chan Aria2WebSocketEvent, options ...string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	e := Aria2TellWaitingTasks(offset, amount, t, options...)

	*Aria2Controller.ActionChan <- e
}

func tellStoppedTask(offset, amount int, pipe *chan Aria2WebSocketEvent, options ...string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	e := Aria2TellStoppedTasks(offset, amount, t, options...)

	*Aria2Controller.ActionChan <- e
}

func changePosition(gid, action string, pos int, pipe *chan Aria2WebSocketEvent) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	e := Aria2ChangePosition(gid, action, t, pos)

	*Aria2Controller.ActionChan <- e
}

func getAria2TaskSettings(gid string, pipe *chan Aria2WebSocketEvent) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	e := Aria2GetOption(gid, t)

	*Aria2Controller.ActionChan <- e
}

func getAria2Settings(pipe *chan Aria2WebSocketEvent) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	e := Aria2GetGlobalOption(t)

	*Aria2Controller.ActionChan <- e
}

func changeAria2TaskSettings(gid string, pipe *chan Aria2WebSocketEvent, options map[string]string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	e := Aria2ChangeOption(gid, t, options)

	*Aria2Controller.ActionChan <- e
}

func changeAria2Settings(pipe *chan Aria2WebSocketEvent, options map[string]string) {
	t := genToken()

	workLock.Lock()
	workQueue[t] = pipe
	workLock.Unlock()

	e := Aria2ChangeGlobalOption(t, options)

	*Aria2Controller.ActionChan <- e
}

func updatePausedData() {
	tasks := NewAria2TaskGIDList()

	paused, err := tasks.TellWaiting(0, 1000)
	if err != nil {
		log.Println(err)
		return
	}

	if paused != nil {
		for i := range paused {
			TaskPaused = append(TaskPaused, (paused)[i].Gid)
		}
	}
}

func autoSorting() {
	cleanAllFinished()
	tasks := NewAria2TaskGIDList()

	if totalActiveAllowed == 0 {
		i, err := tasks.GetGlobalOptions()
		if err != nil {
			log.Println(err)
			return
		}

		totalActiveAllowed, err = strconv.Atoi(i.MaxConcurrentDownloads)

		if err != nil {
			log.Println(err)
			return
		}
	}

	infos, err := tasks.TellActive()
	if err != nil {
		log.Println(err)
		return
	}

	BTUpdating := 0

	taskNeedAction := NewAria2TaskGIDList()
	if totalActiveAllowed <= len(infos)+1 { //Lock download progress when other task waiting
		for i := range infos {
			if infos[i].TotalLength == infos[i].CompletedLength {
				BTUpdating++
			}

			if infos[i].DownloadSpeed == "0" && (infos[i].TotalLength != infos[i].CompletedLength || infos[i].CompletedLength == "0") {
				contained := false
				for j := range TaskStucked {
					if TaskStucked[j] == infos[i].Gid {
						contained = true
						break
					}
				}
				if contained {
					taskNeedAction.Gids = append(taskNeedAction.Gids, infos[i].Gid)
				} else {
					TaskStucked = append(TaskStucked, infos[i].Gid)
				}
			}
		}

		if taskNeedAction.Gids != nil {
			err = taskNeedAction.Pause(true)
			if err != nil {
				log.Println(err)
			} else {
				TaskPaused = append(taskNeedAction.Gids)
			}
		}
	}
	//Unlock download progress when waiting list is empty
	if k := len(TaskPaused); k != 0 {
		j := 0
		for i := len(infos) + 1; i <= totalActiveAllowed+BTUpdating; i++ {
			if j < k {
				taskNeedAction.Gids = append(taskNeedAction.Gids, TaskPaused[j])
				j++
			} else {
				break
			}
		}
		err = taskNeedAction.Unpause()

		if err != nil {
			log.Println(err)
		} else {
			TaskPaused = TaskPaused[j:]
		}
	}

}

func cleanAllFinished() {
	tasks := NewAria2TaskGIDList()

	infos, err := tasks.TellStopped(0, 100)
	if infos == nil || err != nil {
		log.Println(err)
		return
	}

	if len(infos) == 0 {
		return
	}

	for i := range infos {
		tasks.Gids = append(tasks.Gids, (infos)[i].Gid)
	}

	err = tasks.RemoveDownloadResult()
	if err != nil {
		log.Println(err)
	}
}

func dispatchAria2Event() {
	for {
		select {
		case e := <-*Aria2Controller.InfoChan:
			workLock.RLock()
			if c, ok := workQueue[e.Id]; ok {
				*c <- e
			} else if e.Method.IsDownloadEvent() {
				go subscribeProcess(e)
			}
			workLock.RUnlock()
		}
	}
}

func subscribeProcess(e Aria2WebSocketEvent) {
	if e.Params == nil {
		log.Println(`download event didn't contain a gid`)
		return
	}

	i, ok := e.Params[0].(map[string]interface{})
	if !ok {
		log.Println(`download event parse failed`)
		return
	}

	gid, ok := i[`gid`].(string)
	if !ok {
		log.Println(`download event parse failed`)
		return
	}

	subLock.RLock()
	slice := subscriptionQueue[e.Method]
	subLock.RUnlock()

	info := Aria2WebSocketGIDInfo{Gid: gid}
	for i := range slice {
		*slice[i] <- info
	}
}

func genToken() string {
	return RandString(32)
}

func deletePipeInSubscriptionQueue(pipe *chan Aria2WebSocketGIDInfo, t ...Method) {
	subLock.Lock()
	for j := range t {
		for i := 0; i < len(subscriptionQueue[t[j]]); i++ {
			if subscriptionQueue[t[j]][i] == pipe {
				subscriptionQueue[t[j]] = append(subscriptionQueue[t[j]][:i], subscriptionQueue[t[j]][i+1:]...)
				i-- // maintain the correct index
			}
		}
	}
	close(*pipe)
	subLock.Unlock()
}
