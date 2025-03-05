package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Maalog 定义日志信息结构
type Maalog struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Code      int    `json:"code"`
	Class    string `json:"class,omitempty"`
}

// AsstMsg 枚举定义（日志编码）
// Global Info
const (
	InternalError     = 0  // 内部错误
	InitFailed        = 1  // 初始化失败
	ConnectionInfo    = 2  // 连接相关信息
	AllTasksCompleted = 3  // 全部任务完成
	AsyncCallInfo     = 4  // 外部异步调用信息
	Destroyed         = 5  // 实例已销毁
)

// TaskChain Info
const (
	TaskChainError     = 10000 // 任务链执行/识别错误
	TaskChainStart     = 10001 // 任务链开始
	TaskChainCompleted = 10002 // 任务链完成
	TaskChainExtraInfo = 10003 // 任务链额外信息
	TaskChainStopped   = 10004 // 任务链手动停止
)

// SubTask Info
const (
	SubTaskError     = 20000 // 原子任务执行/识别错误
	SubTaskStart     = 20001 // 原子任务开始
	SubTaskCompleted = 20002 // 原子任务完成
	SubTaskExtraInfo = 20003 // 原子任务额外信息
	SubTaskStopped   = 20004 // 原子任务手动停止
)

// codeDescriptions 存放各个日志编码对应的描述信息
var codeDescriptions = map[int]string{
	// Global Info
	InternalError:     "内部错误",
	InitFailed:        "初始化失败",
	ConnectionInfo:    "连接相关信息",
	AllTasksCompleted: "全部任务完成",
	AsyncCallInfo:     "外部异步调用信息",
	Destroyed:         "实例已销毁",

	// TaskChain Info
	TaskChainError:     "任务链执行/识别错误",
	TaskChainStart:     "任务链开始",
	TaskChainCompleted: "任务链完成",
	TaskChainExtraInfo: "任务链额外信息",
	TaskChainStopped:   "任务链手动停止",

	// SubTask Info
	SubTaskError:     "子任务执行/识别错误",
	SubTaskStart:     "子任务开始",
	SubTaskCompleted: "子任务完成",
	SubTaskExtraInfo: "子任务额外信息",
	SubTaskStopped:   "子任务手动停止",
}

// taskchainTranslation 存放英文 taskchain 到中文的映射
var taskchainTranslation = map[string]string{
	"StartUp":              "开始唤醒",
	"CloseDown":            "关闭游戏",
	"Fight":                "刷理智",
	"Mall":                 "信用点及购物",
	"Recruit":              "自动公招",
	"Infrast":              "基建换班",
	"Award":                "领取日常奖励",
	"Roguelike":            "无限刷肉鸽",
	"Copilot":              "自动抄作业",
	"SSSCopilot":           "自动抄保全作业",
	"Depot":                "仓库识别",
	"OperBox":              "干员 box 识别",
	"ReclamationAlgorithm": "生息演算",
	"Custom":               "自定义任务",
	"SingleStep":           "单步任务",
	"VideoRecognition":     "视频识别任务",
	"Debug":                "调试",
}

type SSEClientConn struct {
    Req       *ghttp.Request
    Ch        chan Maalog
    CreatedAt time.Time
}


// NewMaalogFromDetail 将 detail 中的关键信息解析成中文描述，并设置 Level 字段
func NewMaalogFromDetail(detail string, code int) *Maalog {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(detail), &data)
	if err != nil {
		// 如果不能解析为 JSON，则直接将原始字符串作为描述信息
		data = map[string]interface{}{
			"raw": detail,
		}
	}

	// 根据 code 查找对应的描述（如果没有则默认为“未知类型”）
	description, ok := codeDescriptions[code]
	if !ok {
		description = "未知类型"
	}
	eclass := ""
	var msg string
	// 针对不同的编码做对应的 JSON 字段解析
	switch code {
	// Global Info 类型：内部错误、初始化失败
	case InternalError, InitFailed:
		// 预期 JSON 格式：
		// {
		//   "what": string, "why": string, "details": object
		// }
		what, _ := data["what"].(string)
		why, _ := data["why"].(string)
		msg = description
		if what != "" {
			msg += "，错误类型：" + what
		}
		if why != "" {
			msg += "，原因：" + why
		}
		if det, ok := data["details"]; ok {
			detBytes, _ := json.Marshal(det)
			msg += "，详细信息：" + string(detBytes)
		}

	// ConnectionInfo
	case ConnectionInfo:
		// JSON 格式：
		// {
		//   "what": string,
		//   "why": string,
		//   "uuid": string,
		//   "details": {
		//       "adb": string, "address": string, "config": string
		//   }
		// }
		what, _ := data["what"].(string)
		why, _ := data["why"].(string)
		uuid, _ := data["uuid"].(string)
		var detInfo string
		if d, exists := data["details"]; exists {
			if dMap, ok := d.(map[string]interface{}); ok {
				adb, _ := dMap["adb"].(string)
				address, _ := dMap["address"].(string)
				config, _ := dMap["config"].(string)
				detInfo = fmt.Sprintf("adb路径：%s，地址：%s，配置：%s", adb, address, config)
			}
		}
		msg = description
		if what != "" {
			msg += "，信息类型：" + what
		}
		if why != "" {
			msg += "，原因：" + why
		}
		if uuid != "" {
			msg += "，设备ID：" + uuid
		}
		if detInfo != "" {
			msg += "，详细信息：" + detInfo
		}

	// AsyncCallInfo
	case AsyncCallInfo:
		// JSON 格式：
		// {
		//   "uuid": string,
		//   "what": string,
		//   "async_call_id": int,
		//   "details": {
		//       "ret": bool, "cost": int64
		//   }
		// }
		what, _ := data["what"].(string)
		uuid, _ := data["uuid"].(string)
		asyncCallID := 0
		if aid, ok := data["async_call_id"].(float64); ok {
			asyncCallID = int(aid)
		}
		var detInfo string
		if d, exists := data["details"]; exists {
			if dMap, ok := d.(map[string]interface{}); ok {
				ret := false
				if r, ok := dMap["ret"].(bool); ok {
					ret = r
				}
				cost := 0
				if c, ok := dMap["cost"].(float64); ok {
					cost = int(c)
				}
				detInfo = fmt.Sprintf("调用返回值：%v，耗时：%d毫秒", ret, cost)
			}
		}
		msg = description
		if what != "" {
			msg += "，回调类型：" + what
		}
		if uuid != "" {
			msg += "，设备ID：" + uuid
		}
		msg += fmt.Sprintf("，异步调用ID：%d", asyncCallID)
		if detInfo != "" {
			msg += "，" + detInfo
		}

	// AllTasksCompleted
	case AllTasksCompleted:
		// 预期 JSON 格式：
		// {
		//   "taskchain": string,
		//   "uuid": string,
		//   "finished_tasks": [int, ...]
		// }
		taskchain, _ := data["taskchain"].(string)
		// 翻译 taskchain 字段
		if translated, ok := taskchainTranslation[taskchain]; ok {
			taskchain = translated
		}
		uuid, _ := data["uuid"].(string)
		finishedTasks := ""
		if tasks, ok := data["finished_tasks"].([]interface{}); ok {
			for i, t := range tasks {
				finishedTasks += fmt.Sprintf("%v", t)
				if i < len(tasks)-1 {
					finishedTasks += "、"
				}
			}
		}
		msg = description
		if taskchain != "" {
			msg += "，任务链：" + taskchain
		}
		if uuid != "" {
			msg += "，设备ID：" + uuid
		}
		if finishedTasks != "" {
			msg += "，已完成任务ID：" + finishedTasks
		}

	// TaskChain 类型
	case TaskChainError, TaskChainStart, TaskChainCompleted, TaskChainExtraInfo, TaskChainStopped:
		// 预期 JSON 格式：
		// {
		//   "taskchain": string, "taskid": int, "uuid": string
		// }
		taskchain, _ := data["taskchain"].(string)
		// 翻译 taskchain 字段
		if translated, ok := taskchainTranslation[taskchain]; ok {
			taskchain = translated
		}
		taskid := 0
		if t, ok := data["taskid"].(float64); ok {
			taskid = int(t)
		}
		uuid, _ := data["uuid"].(string)
		msg = description
		if taskchain != "" {
			msg += "：" + taskchain
		}
		msg += fmt.Sprintf("，任务ID：%d", taskid)
		if uuid != "" {
			msg += "，设备ID：" + uuid
		}

	// SubTask 类型
	case SubTaskError, SubTaskStart, SubTaskCompleted, SubTaskExtraInfo, SubTaskStopped:
		// 预期 JSON 格式：
		// {
		//   "subtask": string,
		//   "class": string,
		//   "taskchain": string,
		//   "taskid": int,
		//   "details": object,
		//   "uuid": string
		// }
		subtask, _ := data["subtask"].(string)
		class, _ := data["class"].(string)
		taskchain, _ := data["taskchain"].(string)
		
		// 翻译 taskchain 字段
		if translated, ok := taskchainTranslation[taskchain]; ok {
			taskchain = translated
		}
		taskid := 0
		if t, ok := data["taskid"].(float64); ok {
			taskid = int(t)
		}
		uuid, _ := data["uuid"].(string)
		var detInfo string
		if d, exists := data["details"]; exists {
			detBytes, _ := json.Marshal(d)
			detInfo = string(detBytes)
		}
		msg = description
		if subtask != "" {
			msg += "，子任务：" + subtask
		}
		if class != "" {
			msg += "，类：" + class
		}
		if taskchain != "" {
			msg += "，任务链：" + taskchain
		}
		msg += fmt.Sprintf("，任务ID：%d", taskid)
		if uuid != "" {
			msg += "，设备ID：" + uuid
		}
		if detInfo != "" {
			msg += "，详细信息：" + detInfo
		}
		eclass = class

	default:
		// 若为未知的编码，则直接拼接原始字符串
		msg = description + "，详情：" + detail
	}

	// 根据日志编码设置 Level
	var level string
	switch {
	case code == InternalError || code == TaskChainError || code == SubTaskError || code == InitFailed:
		level = "ERROR"
	case code == Destroyed || code == AllTasksCompleted:
		level = "WARN"
	default:
		level = "INFO"
	}
	if eclass == "asst::StageDropsTaskPlugin" {
		level = "ERROR"
	}

	return &Maalog{
		Level:     level,
		Message:   msg,
		Timestamp: time.Now().Unix(),
		Code:      code,
	}
}


func Filterd(code int) bool {
	switch code {
	case ConnectionInfo, AsyncCallInfo, Destroyed, SubTaskCompleted, 
		 SubTaskError, SubTaskStart, SubTaskStopped:
		return true
	}


	return false
}


// GetDescription 根据 Maalog 中的 Code 返回对应的日志描述信息
func (m *Maalog) GetDescription() string {
	if desc, ok := codeDescriptions[m.Code]; ok {
		return desc
	}
	
	return "未知日志类型"
}


// GetAsstMsgDescription 通过输入日志编码，返回描述信息
func GetAsstMsgDescription(code int) string {
	if desc, ok := codeDescriptions[code]; ok {
		return desc
	}
	return "未知日志类型"
}