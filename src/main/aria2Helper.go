package main

import (
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/url"
)

type Method string

type Aria2Instance struct {
	InfoChan   *chan Aria2WebSocketEvent
	ActionChan *chan Aria2WebSocketEvent
}

var Aria2Controller Aria2Instance

const (
	DownloadStart      Method = "aria2.onDownloadStart"
	DownloadPause      Method = "aria2.onDownloadPause"
	DownloadStop       Method = "aria2.onDownloadStop"
	DownloadComplete   Method = "aria2.onDownloadComplete"
	DownloadError      Method = "aria2.onDownloadError"
	BtDownloadComplete Method = "aria2.onBtDownloadComplete"
)

type Aria2WebSocketEvent struct {
	Id         string `json:"id"`
	RPCVersion string `json:"jsonrpc"`
	Method     `json:"method,omitempty"`
	Params     []Aria2Param         `json:"params,omitempty"`
	Result     *jsoniter.RawMessage `json:"result,omitempty"`
	Error      Aria2Error           `json:"error,omitempty"`
}

type Aria2Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Aria2WebSocketGIDInfo struct {
	Gid string `json:"gid"`
}

func (I Aria2Instance) StartAria2Websocket(url *url.URL) {
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	go func() {
		for {
			select {
			case action := <-*I.ActionChan:
				err = c.WriteJSON(action)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}
		m, err := deserializeAria2Event(msg)
		if err != nil {
			log.Println(err)
			continue
		}
		*I.InfoChan <- m
	}
}

func deserializeAria2Event(msg []byte) (Aria2WebSocketEvent, error) {

	var e Aria2WebSocketEvent

	err := jsoniter.Unmarshal(msg, &e)
	if err != nil {
		return e, err
	}

	return e, nil
}

func (M Method) IsDownloadEvent() bool {
	return (M == BtDownloadComplete || M == DownloadStart || M == DownloadComplete || M == DownloadPause || M == DownloadError || M == DownloadStop)
}

type Aria2Param interface{}

type Aria2DownloadParam struct {
	Dir    string   `json:"dir,omitempty"`
	Header []string `json:"header,omitempty"`
	Out    string   `json:"out,omitempty"`
}

type Aria2DownloadTask struct {
	Source      string
	Destination string
	OutName     string
	Headers     []string
	Type        Method
}

func (T Aria2DownloadTask) Download(token string) Aria2WebSocketEvent {
	if T.Type == "aria2.addTorrent" {
		return Aria2WebSocketEvent{
			Id:         token,
			RPCVersion: "2.0",
			Method:     T.Type,
			Params: []Aria2Param{T.Source, []string{}, Aria2DownloadParam{
				Dir:    T.Destination,
				Header: T.Headers,
				Out:    T.OutName,
			}},
		}
	} else {
		return Aria2WebSocketEvent{
			Id:         token,
			RPCVersion: "2.0",
			Method:     T.Type,
			Params: []Aria2Param{[]string{T.Source}, Aria2DownloadParam{
				Dir:    T.Destination,
				Header: T.Headers,
				Out:    T.OutName,
			}},
		}
	}

}

func Aria2RemoveTask(gid, token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.remove",
		Params:     []Aria2Param{gid},
	}
}

func Aria2ForceRemoveTask(gid, token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.forceRemove",
		Params:     []Aria2Param{gid},
	}
}

func Aria2PauseTask(gid, token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.pause",
		Params:     []Aria2Param{gid},
	}
}

func Aria2PauseAllTask(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.pauseAll",
		Params:     []Aria2Param{},
	}
}

func Aria2ForcePauseTask(gid, token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.forcePause",
		Params:     []Aria2Param{gid},
	}
}

func Aria2ForcePauseAllTask(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.forcePauseAll",
		Params:     []Aria2Param{},
	}
}

func Aria2UnpauseTask(gid, token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.unpause",
		Params:     []Aria2Param{gid},
	}
}

func Aria2UnpauseAllTask(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.unpauseAll",
		Params:     []Aria2Param{},
	}
}

func Aria2TellStatusTasks(gid, token string, options ...string) Aria2WebSocketEvent {
	var P []Aria2Param
	P = append(P, gid)

	if len(options) > 0 {
		P = append(P, options)
	}
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.tellStatus",
		Params:     P,
	}
}

func Aria2GetUris(gid, token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.getUris",
		Params:     []Aria2Param{gid},
	}
}

func Aria2GetFiles(gid, token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.getFiles",
		Params:     []Aria2Param{gid},
	}
}

func Aria2GetPeers(gid, token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.getPeers",
		Params:     []Aria2Param{gid},
	}
}

func Aria2GetServers(gid, token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.getServers",
		Params:     []Aria2Param{gid},
	}
}

func Aria2TellActiveTasks(token string, options ...string) Aria2WebSocketEvent {
	var P []Aria2Param

	if len(options) > 0 {
		P = append(P, options)
	}
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.tellActive",
		Params:     P,
	}
}

func Aria2TellWaitingTasks(offset, number int, token string, options ...string) Aria2WebSocketEvent {
	var P []Aria2Param
	P = append(P, offset, number)

	if len(options) > 0 {
		P = append(P, options)
	}
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.tellWaiting",
		Params:     P,
	}
}

func Aria2TellStoppedTasks(offset, number int, token string, options ...string) Aria2WebSocketEvent {
	var P []Aria2Param
	P = append(P, offset, number)

	if len(options) > 0 {
		P = append(P, options)
	}
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.tellStopped",
		Params:     P,
	}
}

const (
	Aria2PositionSET     string = "POS_SET"
	Aria2PositionCURRENT string = "POS_CURRENT"
	Aria2PositionEND     string = "POS_END"
)

func Aria2ChangePosition(gid, action, token string, position int) Aria2WebSocketEvent {
	var P []Aria2Param
	P = append(P, gid, position, action)

	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.changePosition",
		Params:     P,
	}
}

func Aria2ChangeUri(gid, token string, fileIndex int, delUri []string, addUri []string) Aria2WebSocketEvent {
	var P []Aria2Param
	P = append(P, gid, fileIndex, delUri, addUri)

	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.changeUri",
		Params:     P,
	}
}

func Aria2GetOption(gid, token string) Aria2WebSocketEvent {
	var P []Aria2Param
	P = append(P, gid)

	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.getOption",
		Params:     P,
	}
}

const (
//    name =
)

func Aria2ChangeOption(gid, token string, options map[string]string) Aria2WebSocketEvent {
	var P []Aria2Param
	P = append(P, gid, options)

	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.changeOption",
		Params:     P,
	}
}

func Aria2GetGlobalOption(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.getGlobalOption",
		Params:     []Aria2Param{},
	}
}

func Aria2ChangeGlobalOption(token string, options map[string]string) Aria2WebSocketEvent {
	var P []Aria2Param
	P = append(P, options)

	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.changeGlobalOption",
		Params:     P,
	}
}

func Aria2GetGlobalStat(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.getGlobalStat",
		Params:     []Aria2Param{},
	}
}

func Aria2PurgeDownloadResult(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.purgeDownloadResult",
		Params:     []Aria2Param{},
	}
}

func Aria2RemoveDownloadResult(gid, token string) Aria2WebSocketEvent {
	var P []Aria2Param
	P = append(P, gid)

	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.removeDownloadResult",
		Params:     P,
	}
}

func Aria2GetVersion(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.getVersion",
		Params:     []Aria2Param{},
	}
}

func Aria2GetSessionInfo(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.getSessionInfo",
		Params:     []Aria2Param{},
	}
}

func Aria2Shutdown(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.shutdown",
		Params:     []Aria2Param{},
	}
}

func Aria2ForceShutdown(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.forceShutdown",
		Params:     []Aria2Param{},
	}
}

func Aria2SaveSession(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "aria2.saveSession",
		Params:     []Aria2Param{},
	}
}

type MultiCallFunction struct {
	MethodName string       `json:"methodName"`
	Params     []Aria2Param `json:"params"`
}

func Aria2SystemMultiCall(token string, f []MultiCallFunction) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "system.multicall",
		Params:     []Aria2Param{f},
	}
}

func Aria2SystemListMethods(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "system.listMethods",
	}
}

func Aria2SystemListNotifications(token string) Aria2WebSocketEvent {
	return Aria2WebSocketEvent{
		Id:         token,
		RPCVersion: "2.0",
		Method:     "system.listNotifications",
	}
}

type Aria2UniversalTaskCallback struct {
	Result string `json:"result"`
}

type Aria2AddMetalinkTaskCallback struct {
	Result []string `json:"result"`
}

type Aria2TaskStatusCallback struct {
	Result Aria2TaskStatus `json:"result"`
}

type Aria2TaskListStatusCallback struct {
	Result []Aria2TaskStatus `json:"result"`
}

type Aria2TaskStatus struct {
	Gid                    string           `json:"gid,omitempty"`
	Status                 string           `json:"status,omitempty"`
	TotalLength            string           `json:"totalLength,omitempty"`
	CompletedLength        string           `json:"completedLength,omitempty"`
	UploadLength           string           `json:"uploadLength,omitempty"`
	Bitfield               string           `json:"bitfield,omitempty"`
	DownloadSpeed          string           `json:"downloadSpeed,omitempty"`
	UploadSpeed            string           `json:"uploadSpeed,omitempty"`
	InfoHash               string           `json:"infoHash,omitempty"`
	NumSeeders             string           `json:"numSeeders,omitempty"`
	Seeder                 string           `json:"seeder,omitempty"`
	PieceLength            string           `json:"pieceLength,omitempty"`
	NumPieces              string           `json:"NumPieces,omitempty"`
	Connections            string           `json:"connections,omitempty"`
	ErrorCode              string           `json:"errorCode,omitempty"`
	ErrorMessage           string           `json:"errorMessage,omitempty"`
	FollowedBy             string           `json:"followedBy,omitempty"`
	Following              string           `json:"following,omitempty"`
	BelongsTo              string           `json:"belongsTo,omitempty"`
	Dir                    string           `json:"dir,omitempty"`
	VerifiedLength         string           `json:"verifiedLength,omitempty"`
	VerifyIntegrityPending string           `json:"verifyIntegrityPending,omitempty"`
	Files                  []Aria2FileInfo  `json:"files,omitempty"`
	Bittorrent             Aria2TorrentInfo `json:"bittorrent,omitempty"`
}

type Aria2FilesStatus struct {
	Result []Aria2FileInfo `json:"result"`
}

type Aria2FileInfo struct {
	Index           string         `json:"index"`
	Path            string         `json:"path"`
	Length          string         `json:"length"`
	CompletedLength string         `json:"completedLength"`
	Selected        string         `json:"selected"`
	URIs            []Aria2UriInfo `json:"uris"`
}

type Aria2TorrentInfo struct {
	AnnounceList [][]string           `json:"announceList,omitempty"`
	Comment      string               `json:"comment,omitempty"`
	CreationDate int                  `json:"creationDate,omitempty"`
	Mode         string               `json:"mode,omitempty"`
	Info         Aria2TorrentInfoName `json:"info,omitempty"`
}

type Aria2TorrentInfoName struct {
	Name string `json:"name"`
}

type Aria2UriStatus struct {
	Result []Aria2UriInfo `json:"result"`
}

type Aria2UriInfo struct {
	Status string `json:"status"`
	Uri    string `json:"uri"`
}

type Aria2GlobalStatCallback struct {
	Result Aria2GlobalStat `json:"result"`
}

type Aria2GlobalStat struct {
	DownloadSpeed   string `json:"downloadSpeed"`
	UploadSpeed     string `json:"uploadSpeed"`
	NumActive       string `json:"numActive"`
	NumWaiting      string `json:"numWaiting"`
	NumStopped      string `json:"numStopped"`
	NumStoppedTotal string `json:"numStoppedTotal"`
}

type Aria2GetVersionCallback struct {
	Result Aria2VersionInfo `json:"result"`
}

type Aria2VersionInfo struct {
	Version         string   `json:"version"`
	EnabledFeatures []string `json:"enabledFeatures"`
}

type Aria2GetSessionCallback struct {
	Result Aria2SessionInfo `json:"result"`
}

type Aria2SessionInfo struct {
	SessionId string `json:"sessionId"`
}

type Aria2GetOptionCallback struct {
	Result Aria2OptionInfo `json:"result"`
}

type Aria2OptionInfo struct {
	AllowOverwrite                string `json:"allow-overwrite,omitempty"`
	AllowPieceLengthChange        string `json:"allow-piece-length-change,omitempty"`
	AlwaysResume                  string `json:"always-resume,omitempty"`
	AsyncDns                      string `json:"async-dns,omitempty"`
	AutoFileRenaming              string `json:"auto-file-renaming,omitempty"`
	BtEnableHookAfterHashCheck    string `json:"bt-enable-hook-after-hash-check,omitempty"`
	BtEnableLpd                   string `json:"bt-enable-lpd,omitempty"`
	BtExternalIp                  string `json:"bt-external-ip,omitempty"`
	BtForceEncryption             string `json:"bt-force-encryption,omitempty"`
	BtHashCheckSeed               string `json:"bt-hash-check-seed,omitempty"`
	BtLoadSavedMetadata           string `json:"bt-load-saved-metadata,omitempty"`
	BtMaxPeers                    string `json:"bt-max-peers,omitempty"`
	BtMetadataOnly                string `json:"bt-metadata-only,omitempty"`
	BtMinCryptoLevel              string `json:"bt-min-crypto-level,omitempty"`
	BtRemoveUnselectedFile        string `json:"bt-remove-unselected-file,omitempty"`
	BtRequestPeerSpeedLimit       string `json:"bt-request-peer-speed-limit,omitempty"`
	BtRequireCrypto               string `json:"bt-require-crypto,omitempty"`
	BtSaveMetadata                string `json:"bt-save-metadata,omitempty"`
	BtSeedUnverified              string `json:"bt-seed-unverified,omitempty"`
	BtStopTimeout                 string `json:"bt-stop-timeout,omitempty"`
	BtTracker                     string `json:"bt-tracker,omitempty"`
	BtTrackerConnectTimeout       string `json:"bt-tracker-connect-timeout,omitempty"`
	BtTrackerInterval             string `json:"bt-tracker-interval,omitempty"`
	BtTrackerTimeout              string `json:"bt-tracker-timeout,omitempty"`
	CheckIntegrity                string `json:"check-integrity,omitempty"`
	ConditionalGet                string `json:"conditional-get,omitempty"`
	ConnectTimeout                string `json:"connect-timeout,omitempty"`
	ContentDispositionDefaultUtf8 string `json:"content-disposition-default-utf8,omitempty"`
	Continue                      string `json:"continue,omitempty"`
	Dir                           string `json:"dir,omitempty"`
	DryRun                        string `json:"dry-run,omitempty"`
	EnableHttpKeepAlive           string `json:"enable-http-keep-alive,omitempty"`
	EnableHttpPipelining          string `json:"enable-http-pipelining,omitempty"`
	EnableMmap                    string `json:"enable-mmap,omitempty"`
	EnablePeerExchange            string `json:"enable-peer-exchange,omitempty"`
	FileAllocation                string `json:"file-allocation,omitempty"`
	FollowMetalink                string `json:"follow-metalink,omitempty"`
	FollowTorrent                 string `json:"follow-torrent,omitempty"`
	ForceSave                     string `json:"force-save,omitempty"`
	FtpPasv                       string `json:"ftp-pasv,omitempty"`
	FtpReuseConnection            string `json:"ftp-reuse-connection,omitempty"`
	FtpType                       string `json:"ftp-type,omitempty"`
	HashCheckOnly                 string `json:"hash-check-only,omitempty"`
	HttpAcceptGzip                string `json:"http-accept-gzip,omitempty"`
	HttpAuthChallenge             string `json:"http-auth-challenge,omitempty"`
	HttpNoCache                   string `json:"http-no-cache,omitempty"`
	LowestSpeedLimit              string `json:"lowest-speed-limit,omitempty"`
	MaxConnectionPerServer        string `json:"max-connection-per-server,omitempty"`
	MaxDownloadLimit              string `json:"max-download-limit,omitempty"`
	MaxFileNotFound               string `json:"max-file-not-found,omitempty"`
	MaxMmapLimit                  string `json:"max-mmap-limit,omitempty"`
	MaxResumeFailureTries         string `json:"max-resume-failure-tries,omitempty"`
	MaxTries                      string `json:"max-tries,omitempty"`
	MaxUploadLimit                string `json:"max-upload-limit,omitempty"`
	MetalinkEnableUniqueProtocol  string `json:"metalink-enable-unique-protocol,omitempty"`
	MetalinkPreferredProtocol     string `json:"metalink-preferred-protocol,omitempty"`
	MinSplitSize                  string `json:"min-split-size,omitempty"`
	NoFileAllocationLimit         string `json:"no-file-allocation-limit,omitempty"`
	NoNetrc                       string `json:"no-netrc,omitempty"`
	ParameterizedUri              string `json:"parameterized-uri,omitempty"`
	PauseMetadata                 string `json:"pause-metadata,omitempty"`
	PieceLength                   string `json:"piece-length,omitempty"`
	ProxyMethod                   string `json:"proxy-method,omitempty"`
	RealtimeChunkChecksum         string `json:"realtime-chunk-checksum,omitempty"`
	RemoteTime                    string `json:"remote-time,omitempty"`
	RemoveControlFile             string `json:"remove-control-file,omitempty"`
	RetryWait                     string `json:"retry-wait,omitempty"`
	ReuseUri                      string `json:"reuse-uri,omitempty"`
	RpcSaveUploadMetadata         string `json:"rpc-save-upload-metadata,omitempty"`
	SaveNotFound                  string `json:"save-not-found,omitempty"`
	SeedRatio                     string `json:"seed-ratio,omitempty"`
	Split                         string `json:"split,omitempty"`
	StreamPieceSelector           string `json:"stream-piece-selector,omitempty"`
	Timeout                       string `json:"timeout,omitempty"`
	UriSelector                   string `json:"uri-selector,omitempty"`
	UseHead                       string `json:"use-head,omitempty"`
	UserAgent                     string `json:"user-agent,omitempty"`
}

type Aria2GetGlobalOptionCallback struct {
	Result Aria2GlobalOptionInfo `json:"result"`
}

type Aria2GlobalOptionInfo struct {
	AllowOverwrite                string `json:"allow-overwrite,omitempty"`
	AllowPieceLengthChange        string `json:"allow-piece-length-change,omitempty"`
	AlwaysResume                  string `json:"always-resume,omitempty"`
	AsyncDns                      string `json:"async-dns,omitempty"`
	AutoFileRenaming              string `json:"auto-file-renaming,omitempty"`
	AutoSaveInterval              string `json:"auto-save-interval,omitempty"`
	BtDetachSeedOnly              string `json:"bt-detach-seed-only,omitempty"`
	BtEnableHookAfterHashCheck    string `json:"bt-enable-hook-after-hash-check,omitempty"`
	BtEnableLpd                   string `json:"bt-enable-lpd,omitempty"`
	BtExternalIp                  string `json:"bt-external-ip,omitempty"`
	BtForceEncryption             string `json:"bt-force-encryption,omitempty"`
	BtHashCheckSeed               string `json:"bt-hash-check-seed,omitempty"`
	BtLoadSavedMetadata           string `json:"bt-load-saved-metadata,omitempty"`
	BtMaxOpenFiles                string `json:"bt-max-open-files,omitempty"`
	BtMaxPeers                    string `json:"bt-max-peers,omitempty"`
	BtMetadataOnly                string `json:"bt-metadata-only,omitempty"`
	BtMinCryptoLevel              string `json:"bt-min-crypto-level,omitempty"`
	BtRemoveUnselectedFile        string `json:"bt-remove-unselected-file,omitempty"`
	BtRequestPeerSpeedLimit       string `json:"bt-request-peer-speed-limit,omitempty"`
	BtRequireCrypto               string `json:"bt-require-crypto,omitempty"`
	BtSaveMetadata                string `json:"bt-save-metadata,omitempty"`
	BtSeedUnverified              string `json:"bt-seed-unverified,omitempty"`
	BtStopTimeout                 string `json:"bt-stop-timeout,omitempty"`
	BtTracker                     string `json:"bt-tracker,omitempty"`
	BtTrackerConnectTimeout       string `json:"bt-tracker-connect-timeout,omitempty"`
	BtTrackerInterval             string `json:"bt-tracker-interval,omitempty"`
	BtTrackerTimeout              string `json:"bt-tracker-timeout,omitempty"`
	CaCertificate                 string `json:"ca-certificate,omitempty"`
	CheckCertificate              string `json:"check-certificate,omitempty"`
	CheckIntegrity                string `json:"check-integrity,omitempty"`
	ConditionalGet                string `json:"conditional-get,omitempty"`
	ConfPath                      string `json:"conf-path,omitempty"`
	ConnectTimeout                string `json:"connect-timeout,omitempty"`
	ConsoleLogLevel               string `json:"console-log-level,omitempty"`
	ContentDispositionDefaultUtf8 string `json:"content-disposition-default-utf8,omitempty"`
	Continue                      string `json:"continue,omitempty"`
	Daemon                        string `json:"daemon,omitempty"`
	DeferredInput                 string `json:"deferred-input,omitempty"`
	DhtFilePath                   string `json:"dht-file-path,omitempty"`
	DhtFilePath6                  string `json:"dht-file-path6,omitempty"`
	DhtListenPort                 string `json:"dht-listen-port,omitempty"`
	DhtMessageTimeout             string `json:"dht-message-timeout,omitempty"`
	Dir                           string `json:"dir,omitempty"`
	DisableIpv6                   string `json:"disable-ipv6,omitempty"`
	DiskCache                     string `json:"disk-cache,omitempty"`
	DownloadResult                string `json:"download-result,omitempty"`
	DryRun                        string `json:"dry-run,omitempty"`
	Dscp                          string `json:"dscp,omitempty"`
	EnableColor                   string `json:"enable-color,omitempty"`
	EnableDht                     string `json:"enable-dht,omitempty"`
	EnableDht6                    string `json:"enable-dht6,omitempty"`
	EnableHttpKeepAlive           string `json:"enable-http-keep-alive,omitempty"`
	EnableHttpPipelining          string `json:"enable-http-pipelining,omitempty"`
	EnableMmap                    string `json:"enable-mmap,omitempty"`
	EnablePeerExchange            string `json:"enable-peer-exchange,omitempty"`
	EnableRpc                     string `json:"enable-rpc,omitempty"`
	EventPoll                     string `json:"event-poll,omitempty"`
	FileAllocation                string `json:"file-allocation,omitempty"`
	FollowMetalink                string `json:"follow-metalink,omitempty"`
	FollowTorrent                 string `json:"follow-torrent,omitempty"`
	ForceSave                     string `json:"force-save,omitempty"`
	FtpPasv                       string `json:"ftp-pasv,omitempty"`
	FtpReuseConnection            string `json:"ftp-reuse-connection,omitempty"`
	FtpType                       string `json:"ftp-type,omitempty"`
	HashCheckOnly                 string `json:"hash-check-only,omitempty"`
	Help                          string `json:"help,omitempty"`
	HttpAcceptGzip                string `json:"http-accept-gzip,omitempty"`
	HttpAuthChallenge             string `json:"http-auth-challenge,omitempty"`
	HttpNoCache                   string `json:"http-no-cache,omitempty"`
	HumanReadable                 string `json:"human-readable,omitempty"`
	KeepUnfinishedDownloadResult  string `json:"keep-unfinished-download-result,omitempty"`
	ListenPort                    string `json:"listen-port,omitempty"`
	LogLevel                      string `json:"log-level,omitempty"`
	LowestSpeedLimit              string `json:"lowest-speed-limit,omitempty"`
	MaxConcurrentDownloads        string `json:"max-concurrent-downloads,omitempty"`
	MaxConnectionPerServer        string `json:"max-connection-per-server,omitempty"`
	MaxDownloadLimit              string `json:"max-download-limit,omitempty"`
	MaxDownloadResult             string `json:"max-download-result,omitempty"`
	MaxFileNotFound               string `json:"max-file-not-found,omitempty"`
	MaxMmapLimit                  string `json:"max-mmap-limit,omitempty"`
	MaxOverallDownloadLimit       string `json:"max-overall-download-limit,omitempty"`
	MaxOverallUploadLimit         string `json:"max-overall-upload-limit,omitempty"`
	MaxResumeFailureTries         string `json:"max-resume-failure-tries,omitempty"`
	MaxTries                      string `json:"max-tries,omitempty"`
	MaxUploadLimit                string `json:"max-upload-limit,omitempty"`
	MetalinkEnableUniqueProtocol  string `json:"metalink-enable-unique-protocol,omitempty"`
	MetalinkPreferredProtocol     string `json:"metalink-preferred-protocol,omitempty"`
	MinSplitSize                  string `json:"min-split-size,omitempty"`
	MinTlsVersion                 string `json:"min-tls-version,omitempty"`
	NetrcPath                     string `json:"netrc-path,omitempty"`
	NoConf                        string `json:"no-conf,omitempty"`
	NoFileAllocationLimit         string `json:"no-file-allocation-limit,omitempty"`
	NoNetrc                       string `json:"no-netrc,omitempty"`
	OnBtDownloadComplete          string `json:"on-bt-download-complete,omitempty"`
	OnDownloadComplete            string `json:"on-download-complete,omitempty"`
	OnDownloadError               string `json:"on-download-error,omitempty"`
	OptimizeConcurrentDownloads   string `json:"optimize-concurrent-downloads,omitempty"`
	ParameterizedUri              string `json:"parameterized-uri,omitempty"`
	PauseMetadata                 string `json:"pause-metadata,omitempty"`
	PeerAgent                     string `json:"peer-agent,omitempty"`
	PeerIdPrefix                  string `json:"peer-id-prefix,omitempty"`
	PieceLength                   string `json:"piece-length,omitempty"`
	ProxyMethod                   string `json:"proxy-method,omitempty"`
	Quiet                         string `json:"quiet,omitempty"`
	RealtimeChunkChecksum         string `json:"realtime-chunk-checksum,omitempty"`
	RemoteTime                    string `json:"remote-time,omitempty"`
	RemoveControlFile             string `json:"remove-control-file,omitempty"`
	RetryWait                     string `json:"retry-wait,omitempty"`
	ReuseUri                      string `json:"reuse-uri,omitempty"`
	RlimitNofile                  string `json:"rlimit-nofile,omitempty"`
	RpcAllowOriginAll             string `json:"rpc-allow-origin-all,omitempty"`
	RpcListenAll                  string `json:"rpc-listen-all,omitempty"`
	RpcListenPort                 string `json:"rpc-listen-port,omitempty"`
	RpcMaxRequestSize             string `json:"rpc-max-request-size,omitempty"`
	RpcSaveUploadMetadata         string `json:"rpc-save-upload-metadata,omitempty"`
	RpcSecure                     string `json:"rpc-secure,omitempty"`
	SaveNotFound                  string `json:"save-not-found,omitempty"`
	SaveSession                   string `json:"save-session,omitempty"`
	SaveSessionInterval           string `json:"save-session-interval,omitempty"`
	SeedRatio                     string `json:"seed-ratio,omitempty"`
	ServerStatTimeout             string `json:"server-stat-timeout,omitempty"`
	ShowConsoleReadout            string `json:"show-console-readout,omitempty"`
	ShowFiles                     string `json:"show-files,omitempty"`
	SocketRecvBufferSize          string `json:"socket-recv-buffer-size,omitempty"`
	Split                         string `json:"split,omitempty"`
	Stderr                        string `json:"stderr,omitempty"`
	Stop                          string `json:"stop,omitempty"`
	StreamPieceSelector           string `json:"stream-piece-selector,omitempty"`
	SummaryInterval               string `json:"summary-interval,omitempty"`
	Timeout                       string `json:"timeout,omitempty"`
	TruncateConsoleReadout        string `json:"truncate-console-readout,omitempty"`
	UriSelector                   string `json:"uri-selector,omitempty"`
	UseHead                       string `json:"use-head,omitempty"`
	UserAgent                     string `json:"user-agent,omitempty"`
}
