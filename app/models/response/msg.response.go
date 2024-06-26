package response

type HttpCode int64

const (
	Success                    HttpCode = 100200
	Failed                     HttpCode = 100500
	DataNotExist               HttpCode = 100100
	GroupNotExist              HttpCode = 100201
	UserPasswordError          HttpCode = 100202
	AuthNotExist               HttpCode = 100301
	AuthFail                   HttpCode = 100302
	SameGroupName              HttpCode = 100401
	RequestParamError          HttpCode = 100501
	UserNameNotExist           HttpCode = 100502
	UserNameExist              HttpCode = 100503
	PhoneExist                 HttpCode = 100504
	EmailExist                 HttpCode = 100505
	TokenBuildError            HttpCode = 100506
	TokenTimeOut               HttpCode = 100507
	AddDataError               HttpCode = 100508
	WsCreateError              HttpCode = 100509
	SqlExecuteError            HttpCode = 100510
	DeleteSuccess              HttpCode = 100511
	ExistSameData              HttpCode = 100512
	SoftwareCreateError        HttpCode = 100513
	DataNotNeedUpdate          HttpCode = 100514
	DataDeleteFail             HttpCode = 100515
	LogoutSuccess              HttpCode = 100516
	LogoutFail                 HttpCode = 100517
	DateUpdateError            HttpCode = 100518
	FileNotDir                 HttpCode = 100519
	ReadFileListError          HttpCode = 100520
	EmailSendError             HttpCode = 100521
	CaptchaExpire              HttpCode = 100522
	CaptchaError               HttpCode = 100523
	RoleCreateError            HttpCode = 100524
	ResourceCreateError        HttpCode = 100525
	UserUpdateError            HttpCode = 100526
	UserAuthRoleError          HttpCode = 100527
	UserPasswdResetError       HttpCode = 100528
	UserDeleteError            HttpCode = 100529
	ScriptGetError             HttpCode = 100530
	GetDataError               HttpCode = 100531
	ScriptGroupDeleteError     HttpCode = 100532
	ScriptGroupCreateError     HttpCode = 100533
	ScriptItemGetError         HttpCode = 100534
	DockerContainerListError   HttpCode = 100535
	DockerImageListError       HttpCode = 100536
	DockerConnectionError      HttpCode = 100537
	DockerConnectionStartError HttpCode = 100538
	DockerConnectionStopError  HttpCode = 100539
	DockerImageSearchError     HttpCode = 100540
	ExistRunningTaskError      HttpCode = 100541
	ImageTagsGetError          HttpCode = 100542
	ImageTagSetError           HttpCode = 100543
	ImageDeleteError           HttpCode = 100544
	RedisClientError           HttpCode = 100545
	RedisClientCloseError      HttpCode = 100546
)

var Menus = map[HttpCode]string{
	Success:                    "操作成功",
	Failed:                     "操作失败",
	DataNotExist:               "数据不存在",
	GroupNotExist:              "当前分组不存在",
	UserPasswordError:          "账号密码不正确",
	AuthNotExist:               "认证信息不正确",
	AuthFail:                   "校验认证信息失败",
	SameGroupName:              "存在相同分组名称",
	RequestParamError:          "请求参数错误",
	UserNameNotExist:           "用户名不存在",
	UserNameExist:              "账号已经存在",
	PhoneExist:                 "用户手机号已经存在",
	EmailExist:                 "用户邮箱已经注册",
	TokenBuildError:            "生成Token错误",
	TokenTimeOut:               "认证信息过期",
	AddDataError:               "增加数据失败",
	WsCreateError:              "websocket创建链接失败",
	SqlExecuteError:            "SQL执行错误",
	DeleteSuccess:              "删除成功",
	ExistSameData:              "存在相同的数据",
	SoftwareCreateError:        "软件创建失败",
	DataNotNeedUpdate:          "数据不需要修改",
	DataDeleteFail:             "数据删除失败",
	LogoutSuccess:              "退出登录成功",
	LogoutFail:                 "退出登录失败",
	DateUpdateError:            "数据更新失败",
	FileNotDir:                 "当前资源不是目录",
	ReadFileListError:          "读取文件列表失败",
	EmailSendError:             "验证码发送失败",
	CaptchaExpire:              "验证码过期",
	CaptchaError:               "验证码错误",
	RoleCreateError:            "创建角色失败",
	ResourceCreateError:        "创建资源失败",
	UserUpdateError:            "用户更新失败",
	UserAuthRoleError:          "用户授权角色失败",
	UserPasswdResetError:       "用户重置密码失败",
	UserDeleteError:            "用户删除失败",
	ScriptGetError:             "服务列表数据获取失败",
	GetDataError:               "获取数据失败",
	ScriptGroupDeleteError:     "分组删除失败",
	ScriptGroupCreateError:     "分组创建失败",
	ScriptItemGetError:         "服务项获取去失败",
	DockerContainerListError:   "Docker容器列表获取失败",
	DockerImageListError:       "Docker容器镜像列表获取失败",
	DockerConnectionError:      "Docker连接失败",
	DockerConnectionStartError: "Docker容器启动失败",
	DockerConnectionStopError:  "Docker容器停止失败",
	DockerImageSearchError:     "Docker镜像搜索错误",
	ExistRunningTaskError:      "存在正在运行的任务",
	ImageTagsGetError:          "获取镜像Tags失败",
	ImageTagSetError:           "设置镜像Tag失败",
	ImageDeleteError:           "镜像已经被容器使用或者正在被构建",
	RedisClientError:           "Redis连接失败",
	RedisClientCloseError:      "Redis连接关闭失败",
}

// Message 消息
type Message struct {
	Code HttpCode `json:"code"`
	Msg  string   `json:"message"`
	Data any      `json:"data"`
}

func Ok(data any) Message {
	return Message{
		Code: Success,
		Msg:  Menus[Success],
		Data: data,
	}
}

func Fail(data any) Message {
	return Message{
		Code: Failed,
		Msg:  Menus[Failed],
		Data: data,
	}
}

func ResultCustom(err *BusinessError) Message {
	return Message{
		Code: err.Code,
		Msg:  Menus[err.Code],
		Data: err.Error(),
	}
}

func Result(code HttpCode, data any) Message {
	return Message{
		Code: code,
		Msg:  Menus[code],
		Data: data,
	}
}

type PageData struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Data  any   `json:"records"`
}

func NewPageData(page, size int, total int64, data any) *PageData {
	return &PageData{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  data,
	}
}
