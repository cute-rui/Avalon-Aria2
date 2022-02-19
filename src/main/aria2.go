package main

import (
	"avalon-aria2/src/aria2"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Aria2AgentServiceWorker struct {
	aria2.UnimplementedAria2AgentServer
}

func (a Aria2AgentServiceWorker) AwaitDownload(param *aria2.Param, server aria2.Aria2Agent_AwaitDownloadServer) error {
	var NEWDOWNLOAD bool

	d, g := param.GetDownloadInfoList(), param.GetGIDList()
	if len(d) == 0 && len(g) == 0 {
		return status.Error(codes.InvalidArgument, `Invalid Filename`)
	}

	TokenMap := map[string]string{}
	if len(g) == 0 {
		NEWDOWNLOAD = true
		for i := range d {
			gid, err := NewDownloadTask(d[i].GetURL(), getDownloadOption(d[i])...).DoDownload(d[i].GetToken())
			if err != nil {
				return err
			}

			g = append(g, gid)
			TokenMap[d[i].GetToken()] = gid
		}
	}

	TaskList := NewAria2TaskGIDList(g...)
	var errMsg string

	FileInfo := map[string]*aria2.FileInfo{}
	for k, v := range TokenMap {
		if _, ok := TaskList.Status[v]; ok {
			FileInfo[k] = &aria2.FileInfo{GID: v, IsFinished: false}
			continue
		}
		errMsg = StringBuilder(errMsg, `TOKEN `, k, ` GID NOT FOUND ERROR`)
	}

	if errMsg != `` {
		if err := server.Send(&aria2.Result{
			Code:     400,
			Msg:      errMsg,
			FileInfo: FileInfo,
		}); err != nil {
			return err
		}
	}

	if err := server.Send(&aria2.Result{
		Code:     1,
		Msg:      `Received`,
		FileInfo: FileInfo,
	}); err != nil {
		return err
	}

	errMsg = TaskList.WaitUntilFinished()

	if NEWDOWNLOAD {
		for k, v := range TokenMap {
			if b, ok := TaskList.Status[v]; ok {
				FileInfo[k] = &aria2.FileInfo{GID: v, IsFinished: b}
				continue
			}
			errMsg += StringBuilder(errMsg, `TOKEN `, k, ` GID NOT FOUND ERROR`)
		}
	} else {
		for k, v := range TaskList.Status {
			FileInfo[k] = &aria2.FileInfo{GID: k, IsFinished: v}
		}
	}

	if errMsg == `` {
		if err := server.Send(&aria2.Result{
			Code:     0,
			Msg:      `OK`,
			FileInfo: FileInfo,
		}); err != nil {
			return err
		}
	} else {
		if err := server.Send(&aria2.Result{
			Code:     400,
			Msg:      errMsg,
			FileInfo: FileInfo,
		}); err != nil {
			return err
		}
	}

	return nil

}

func (a Aria2AgentServiceWorker) CheckDownload(ctx context.Context, param *aria2.Param) (*aria2.Result, error) {
	GIDList := param.GetGIDList()
	if len(GIDList) == 0 {
		return nil, status.Error(codes.InvalidArgument, `Invalid Filename`)
	}

	TaskList := NewAria2TaskGIDList(GIDList...)
	errMsg := TaskList.CheckStatus()

	FileInfo := map[string]*aria2.FileInfo{}
	for k, v := range TaskList.Status {
		FileInfo[k] = &aria2.FileInfo{IsFinished: v}
	}

	if errMsg == `` {
		return &aria2.Result{Code: 0, Msg: `OK`, FileInfo: FileInfo}, nil
	} else {
		return &aria2.Result{Code: 400, Msg: errMsg, FileInfo: FileInfo}, nil
	}
}

func getDownloadOption(info *aria2.DownloadInfo) []DownloadOption {
	var DownloadOption []DownloadOption
	if m := info.DownloadOption.GetWithHeader(); m != nil {
		for k, v := range m {
			DownloadOption = append(DownloadOption, WithHeader(k, v))
		}
	}

	if dst := info.GetDestination(); dst != `` {
		DownloadOption = append(DownloadOption, WithDestination(dst))
	}

	if name := info.GetFileName(); name != `` {
		DownloadOption = append(DownloadOption, WithOutName(name))
	}

	var t Method
	switch info.GetDownloadType() {
	case aria2.DownloadType_HTTP:
		t = HTTPDownload
	case aria2.DownloadType_MAGLINK:
		t = HTTPDownload
	case aria2.DownloadType_METALINK:
		t = MetalinkDownload
	case aria2.DownloadType_TORRENT:
		t = TorrentDownload
	}

	DownloadOption = append(DownloadOption, WithTaskType(t))

	return DownloadOption
}
