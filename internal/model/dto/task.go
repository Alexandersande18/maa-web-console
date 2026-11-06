package dto

import "encoding/json"

type TaskType int
type TaskID int32

const (
	StartUp   TaskType = iota //开始唤醒
	Fight                     //刷理智
	Recruit                   //公开招募
	Infrast                   //基建换班
	Visit                     //访问好友
	Mall                      //领取信用及商店购物。
	Award                     //领取日常奖励
	Roguelike                 //无限刷肉鸽
	Copilot                   //自动抄作业
	Depot                     //仓库识别
	CloseDown                 //关闭游戏
)

func (t TaskType) String() string {
	switch t {
	case StartUp:
		return "StartUp"
	case Fight:
		return "Fight"
	case Recruit:
		return "Recruit"
	case Infrast:
		return "Infrast"
	case Mall:
		return "Mall"
	case Award:
		return "Award"
	case Roguelike:
		return "Roguelike"
	case Copilot:
		return "Copilot"
	case Depot:
		return "Depot"
	case CloseDown:
		return "CloseDown"
	default:
		return ""
	}
}

var TaskFixed = map[string]TaskType{
	"StartUp":   StartUp,
	"Fight":     Fight,
	"Recruit":   Recruit,
	"Infrast":   Infrast,
	"Mall":      Mall,
	"Award":     Award,
	"Roguelike": Roguelike,
	"Copilot":   Copilot,
	"Depot":     Depot,
	"CloseDown": CloseDown,
}

type TaskerString string

func (t TaskerString) Serialize() string {
	return string(t)
}

// 添加的任务需要实现该接口
type Tasker interface {
	Serialize() string
	GetType() TaskType
}

// 1. 开始唤醒（StartUp）
// 对应任务参数：
// {
//     "enable": bool,              // 是否启用本任务，可选，默认为 true
//     "client_type": string,       // 客户端版本，可选，默认为空
//                                  // 选项："Official" | "Bilibili" | "txwy" | "YoStarEN" | "YoStarJP" | "YoStarKR"
//     "start_game_enabled": bool,  // 是否自动启动客户端，可选，默认不启动
//     "account_name": string       // 切换账号，可选，默认不切换
//                                  // 仅支持切换至已登录的账号，使用登录名进行查找，保证输入内容在所有已登录账号唯一即可
//                                  // 官服：123****4567，可输入 123****4567、4567、123、3****4567
//                                  // B服：张三，可输入 张三、张、三
// }
type TaskStartUp struct {
	Enable           bool   `json:"enable"           v:"boolean"            d:"true"                                                                 dc:"是否启用本任务，可选，默认为 true"`
	ClientType       string `json:"client_type"      v:"in:Official,Bilibili,txwy,YoStarEN,YoStarJP,YoStarKR" d:""                                dc:"客户端版本，可选，选项：Official,Bilibili,txwy,YoStarEN,YoStarJP,YoStarKR"`
	StartGameEnabled bool   `json:"start_game_enabled" v:""                 d:"false"                                                               dc:"是否自动启动客户端，可选，默认不启动"`
	AccountName      string `json:"account_name"     v:""                   d:""                                                                   dc:"切换账号，可选，支持模糊匹配已登录账号"`
}

func (t *TaskStartUp) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *TaskStartUp) GetType() TaskType {
	return TaskType(StartUp)
}

// 2. 关闭游戏（CloseDown）
// 对应任务参数：
//
//	{
//	    "enable": bool,
//	    "client_type": string,       // "Official" | "Bilibili" | "txwy" | "YoStarEN" | "YoStarJP" | "YoStarKR"
//	}
type TaskCloseDown struct {
	Enable     bool   `json:"enable"      v:"boolean" d:"true"                                                                 dc:"是否启用本任务"`
	ClientType string `json:"client_type" v:"in:Official,Bilibili,txwy,YoStarEN,YoStarJP,YoStarKR" d:""                      dc:"客户端版本，选项：Official,Bilibili,txwy,YoStarEN,YoStarJP,YoStarKR"`
}

func (t *TaskCloseDown) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *TaskCloseDown) GetType() TaskType {
	return TaskType(CloseDown)
}

// 3. 刷理智（Fight）
// 对应任务参数：
// 对应的任务参数
// {
//     "enable": bool,             // 是否启用本任务，可选，默认为 true
//     "stage": string,            // 关卡名，可选，默认为空，识别当前/上次的关卡。不支持运行中设置
//                                 // 支持全部主线关卡，如 "1-7"、"S3-2"等
//                                 // 可在关卡结尾输入Normal/Hard表示需要切换标准与磨难难度
//                                 // 剿灭作战，必须输入 "Annihilation"
//                                 // 当期 SS 活动 后三关，必须输入完整关卡编号
//     "medicine": int,            // 最大使用理智药数量，可选，默认 0
//     "expiring_medicine": int,   // 最大使用 48 小时内过期理智药数量，可选，默认 0
//     "stone": int,               // 最大吃石头数量，可选，默认 0
//     "times": int,               // 指定次数，可选，默认无穷大
//     "series": int,              // 连战次数，可选，0~10，0为自动选择
//     "drops": {                  // 指定掉落数量，可选，默认不指定
//         "30011": int,           // key - item_id, value 数量. key 可参考 resource/item_index.json 文件
//         "30062": int            // 是或的关系
//     },
//     /* 以上全部是或的关系，即任一达到即停止任务 */

//     "report_to_penguin": bool,  // 是否汇报企鹅数据，可选，默认 false
//     "penguin_id": string,       // 企鹅数据汇报 id, 可选，默认为空。仅在 report_to_penguin 为 true 时有效
//     "server": string,           // 服务器，可选，默认 "CN", 会影响掉落识别及上传
//                                 // 选项："CN" | "US" | "JP" | "KR"
//     "client_type": string,      // 客户端版本，可选，默认为空。用于游戏崩溃时重启并连回去继续刷，若为空则不启用该功能
//                                 // 选项："Official" | "Bilibili" | "txwy" | "YoStarEN" | "YoStarJP" | "YoStarKR"
//     "DrGrandet": bool,          // 节省理智碎石模式，可选，默认 false，仅在可能产生碎石效果时生效。
//                                 // 在碎石确认界面等待，直到当前的 1 点理智恢复完成后再立刻碎石
// }
type TaskFight struct {
	Enable           bool           `json:"enable"            v:"boolean"         d:"true"                                                                  dc:"是否启用本任务"`
	Stage            string         `json:"stage"             v:""                d:""                                                                      dc:"关卡名称，支持主线/剿灭/活动关卡"`
	Medicine         int            `json:"medicine"          v:"min:0"           d:"0"                                                                     dc:"最大使用理智药数量"`
	ExpiringMedicine int            `json:"expiring_medicine" v:"min:0"           d:"0"                                                                     dc:"最大使用临期理智药数量"`
	Stone            int            `json:"stone"             v:"min:0"           d:"0"                                                                     dc:"最大吃石头数量"`
	Times            int            `json:"times"             v:"min:-1"          d:"-1"                                                                    dc:"执行次数，2147483647表示无限"`
	Series           int            `json:"series"            v:"min:0|max:10"     d:"1"                                                                     dc:"连战次数，范围0-10，0为自动选择"`
	Drops            map[string]int `json:"drops"             v:""                d:""                                                                      dc:"掉落数量限制，满足任一即停"`
	ReportToPenguin  bool           `json:"report_to_penguin" v:""                d:"false"                                                                 dc:"是否上报企鹅数据"`
	PenguinID        string         `json:"penguin_id"        v:"required-if:ReportToPenguin,true" d:""                                                    dc:"企鹅数据ID，上报时必填"`
	Server           string         `json:"server"            v:"in:CN,US,JP,KR"  d:"CN"                                                                    dc:"服务器地区"`
	ClientType       string         `json:"client_type"       v:"in:Official,Bilibili,txwy,YoStarEN,YoStarJP,YoStarKR" d:""                                 dc:"客户端类型"`
	DrGrandet        bool           `json:"DrGrandet"         v:""                d:"false"                                                                 dc:"智能碎石模式"`
}

func (t *TaskFight) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *TaskFight) GetType() TaskType {
	return TaskType(Fight)
}

// 4. 公开招募（Recruit）
// 对应任务参数：
// {
//     "enable": bool,             // 是否启用本任务，可选，默认为 true
//     "refresh": bool,            // 是否刷新三星 Tags, 可选，默认 false
//     "select": [                 // 会去点击标签的 Tag 等级，必选
//         int,
//         ...
//     ],
//     "confirm": [                // 会去点击确认的 Tag 等级，必选。若仅公招计算，可设置为空数组
//         int,
//         ...
//     ],
//     "first_tags": [             // 首选 Tags，仅在 Tag 等级为 3 时有效。可选，默认为空
//         string,                 // 当 Tag 等级为 3 时，会尽可能多地选择这里的 Tags（如果有）
//         ...                     // 而且是强制选择，也就是会忽略所有“让 3 星 Tag 不被选择”的设置
//     ],
//     "extra_tags_mode": int,     // 选择更多的 Tags, 可选, 默认为 0
//                                 // 0 - 默认行为
//                                 // 1 - 选 3 个 Tags, 即使可能冲突
//                                 // 2 - 如果可能, 同时选择更多的高星 Tag 组合, 即使可能冲突
//     "times": int,               // 招募多少次，可选，默认 0。若仅公招计算，可设置为 0
//     "set_time": bool,           // 是否设置招募时限。仅在 times 为 0 时生效，可选，默认 true
//     "expedite": bool,           // 是否使用加急许可，可选，默认 false
//     "expedite_times": int,      // 加急次数，仅在 expedite 为 true 时有效。
//                                 // 可选，默认无限使用（直到 times 达到上限）
//     "skip_robot": bool,         // 是否在识别到小车词条时跳过，可选，默认跳过
//     "recruitment_time": {       // Tag 等级（大于等于 3）和对应的希望招募时限，单位为分钟，默认值都为 540（即 09:00:00）
//         "3": int,
//         "4": int,
//         ...
//     },
//     "report_to_penguin": bool,  // 是否汇报企鹅数据，可选，默认 false
//     "penguin_id": string,       // 企鹅数据汇报 id, 可选，默认为空。仅在 report_to_penguin 为 true 时有效
//     "report_to_yituliu": bool,  // 是否汇报一图流数据，可选，默认 false
//     "yituliu_id": string,       // 一图流汇报 id, 可选，默认为空。仅在 report_to_yituliu 为 true 时有效
//     "server": string,           // 服务器，可选，默认 "CN", 会影响上传
//                                 // 选项："CN" | "US" | "JP" | "KR"
// }
type TaskRecruit struct {
	Enable          bool           `json:"enable"           v:"boolean"         d:"true"                                 dc:"是否启用本任务"`
	Refresh         bool           `json:"refresh"          v:""                d:"false"                                dc:"是否刷新三星标签"`
	Select          []int          `json:"select"           v:"required"                                                 dc:"点击标签等级（必选）"`
	Confirm         []int          `json:"confirm"          v:""                                                         dc:"确认标签等级"`
	FirstTags       []string       `json:"first_tags"       v:""                d:""                                     dc:"首选标签"`
	ExtraTagsMode   int            `json:"extra_tags_mode"  v:"in:0,1,2"        d:"0"                                    dc:"额外标签模式：0-默认,1-选3标签,2-组合高星"`
	Times           int            `json:"times"            v:"min:0"           d:"0"                                    dc:"招募次数"`
	SetTime         bool           `json:"set_time"         v:""                d:"true"                                 dc:"是否设置招募时间"`
	Expedite        bool           `json:"expedite"         v:""                d:"false"                                dc:"是否使用加急许可"`
	ExpediteTimes   int            `json:"expedite_times"   v:"min:0"           d:"0"                                    dc:"加急次数限制"`
	SkipRobot       bool           `json:"skip_robot"       v:""                d:"true"                                 dc:"是否跳过机器人标签"`
	RecruitmentTime map[string]int `json:"recruitment_time" v:""                d:"3:540,4:540,5:540,6:540"             dc:"各星级招募时间（分钟）"`
	ReportToPenguin bool           `json:"report_to_penguin" v:""               d:"false"                                dc:"上报企鹅数据"`
	PenguinID       string         `json:"penguin_id"       v:"required-if:ReportToPenguin,true" d:""                   dc:"企鹅数据ID"`
	ReportToYituliu bool           `json:"report_to_yituliu" v:""               d:"false"                                dc:"上报一图流数据"`
	YituliuID       string         `json:"yituliu_id"       v:"required-if:ReportToYituliu,true" d:""                   dc:"一图流数据ID"`
	Server          string         `json:"server"           v:"in:CN,US,JP,KR"  d:"CN"                                   dc:"服务器地区"`
}

func (t *TaskRecruit) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *TaskRecruit) GetType() TaskType {
	return TaskType(Recruit)
}

// 5. 基建换班（Infrast）
// 对应任务参数：
// {
//     "enable": bool,         // 是否启用本任务，可选，默认为 true
//     "mode": int,            // 换班工作模式，可选，默认 0
//                             // 0 - 默认换班模式，单设施最优解
//                             // 10000 - 自定义换班模式，读取用户配置，可参考 protocol/base-scheduling-schema.md

//     "facility": [           // 要换班的设施（有序），必选。不支持运行中设置
//         string,             // 设施名，"Mfg" | "Trade" | "Power" | "Control" | "Reception" | "Office" | "Dorm"
//         ...
//     ],
//     "drones": string,       // 无人机用途，可选项，默认 _NotUse
//                             // mode==10000 时该字段无效（会被忽略）
//                             // "_NotUse"、"Money"、"SyntheticJade"、"CombatRecord"、"PureGold"、"OriginStone"、"Chip"
//     "threshold": float,     // 工作心情阈值，可选，取值范围 [0, 1.0]，默认 0.3
//                             // mode==10000 时该字段仅针对 "autofill" 有效
//     "replenish": bool,      // 贸易站“源石碎片”是否自动补货，可选，默认 false

//     "dorm_notstationed_enabled": bool, // 是否启用宿舍“未进驻”选项，可选，默认 false
//     "dorm_trust_enabled": bool, // 是否将宿舍剩余位置填入信赖未满干员，可选，默认 false

//     /* 以下参数仅在 mode=10000 时生效，否则会被忽略 */
//     "filename": string,     // 自定义配置路径，必选。不支持运行中设置
//     "plan_index": int,      // 使用配置中的方案序号，必选。不支持运行中设置
// }
type TaskInfrast struct {
	Enable                  bool     `json:"enable"                   v:"boolean"         d:"true"                                                                  dc:"是否启用本任务"`
	Mode                    int      `json:"mode"                     v:"in:0,10000"      d:"0"                                                                     dc:"模式：0-最优解,10000-自定义"`
	Facility                []string `json:"facility"                 v:"required"                                                                                  dc:"换班设施列表（必填）"`
	Drones                  string   `json:"drones"                   v:"in:_NotUse,Money,SyntheticJade,CombatRecord,PureGold,OriginStone,Chip" d:"_NotUse"       dc:"无人机用途"`
	Threshold               float64  `json:"threshold"                v:"min:0|max:1.0"   d:"0.3"                                                                   dc:"心情阈值"`
	Replenish               bool     `json:"replenish"                v:""                d:"false"                                                                 dc:"贸易站自动补货"`
	DormNotstationedEnabled bool     `json:"dorm_notstationed_enabled" v:""               d:"false"                                                                 dc:"启用未进驻宿舍"`
	DormTrustEnabled        bool     `json:"dorm_trust_enabled"       v:""                d:"false"                                                                 dc:"填入信赖未满干员"`
	Filename                string   `json:"filename"                 v:"required-if:Mode,10000"                                                                   dc:"自定义配置文件路径"`
	PlanIndex               int      `json:"plan_index"               v:"required-if:Mode,10000"                                                                   dc:"方案序号"`
}

func (t *TaskInfrast) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *TaskInfrast) GetType() TaskType {
	return TaskType(Infrast)
}

// 6. 领取信用及商店购物（Mall）
// 对应任务参数：
//
// {
//     "enable": bool,         // 是否启用本任务，可选，默认为 true
//     "shopping": bool,       // 是否购物，可选，默认 false。不支持运行中设置
//     "buy_first": [          // 优先购买列表，可选。不支持运行中设置
//         string,             // 商品名，如 "招聘许可"、"龙门币" 等
//         ...
//     ],
//     "blacklist": [          // 黑名单列表，可选。不支持运行中设置
//         string,             // 商品名，如 "加急许可"、"家具零件" 等
//         ...
//     ],
//    "force_shopping_if_credit_full": bool // 是否在信用溢出时无视黑名单，默认为 true
//     "only_buy_discount": bool // 是否只购买折扣物品，只作用于第二轮购买，默认为 false
//     "reserve_max_credit": bool // 是否在信用点低于300时停止购买，只作用于第二轮购买，默认为 false
// }
type TaskMall struct {
	Enable                    bool     `json:"enable"                     v:"boolean" d:"true"                                                                  dc:"是否启用本任务"`
	Shopping                  bool     `json:"shopping"                   v:""        d:"false"                                                                 dc:"是否执行购物"`
	VisitFriends              bool     `json:"visit_friends"              v:""        d:"false"                                                                 dc:"是否访问好友"`
	BuyFirst                  []string `json:"buy_first"                  v:""        d:""                                                                      dc:"优先购买列表"`
	Blacklist                 []string `json:"blacklist"                  v:""        d:""                                                                      dc:"黑名单列表"`
	ForceShoppingIfCreditFull bool     `json:"force_shopping_if_credit_full" v:""     d:"true"                                                                  dc:"信用溢出时强制购物"`
	OnlyBuyDiscount           bool     `json:"only_buy_discount"          v:""        d:"false"                                                                 dc:"仅购折扣商品"`
	ReserveMaxCredit          bool     `json:"reserve_max_credit"         v:""        d:"false"                                                                 dc:"保留最低信用"`
}

func (t *TaskMall) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *TaskMall) GetType() TaskType {
	return TaskType(Mall)
}

// 7. 领取各种奖励（Award）
// 对应任务参数：
//
//	{
//	    "enable": bool,            // 是否启用本任务，可选，默认为 true
//	    "award": bool,             // 领取每日/每周任务奖励，默认为 true
//	    "mail": bool,              // 领取所有邮件奖励，默认为 false
//	    "recruit": bool,           // 领取限定池子赠送的每日免费单抽，默认为 false
//	    "orundum": bool,           // 领取幸运墙的合成玉奖励，默认为 false
//	    "mining": bool,            // 领取限时开采许可的合成玉奖励，默认为 false
//	    "specialaccess": bool      // 领取五周年赠送的月卡奖励，默认为 false
//	}
type TaskAward struct {
	Enable        bool `json:"enable"         v:"boolean" d:"true"                                                                  dc:"是否启用本任务"`
	Award         bool `json:"award"          v:""        d:"true"                                                                  dc:"领取任务奖励"`
	Mail          bool `json:"mail"           v:""        d:"false"                                                                 dc:"领取邮件奖励"`
	Recruit       bool `json:"recruit"        v:""        d:"false"                                                                 dc:"领取公招奖励"`
	Orundum       bool `json:"orundum"        v:""        d:"false"                                                                 dc:"领取合成玉"`
	Mining        bool `json:"mining"         v:""        d:"false"                                                                 dc:"领取开采许可"`
	Specialaccess bool `json:"specialaccess"  v:""        d:"false"                                                                 dc:"领取月卡奖励"`
}

func (t *TaskAward) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *TaskAward) GetType() TaskType {
	return TaskType(Award)
}

// 8. 无限刷肉鸽（Roguelike）
// 对应任务参数较多，示例如下：
//
// {
//     "enable": bool,  // 是否启用本任务，可选，默认值 true
//     "theme": string, // 主题，可选，默认值 "Phantom"
//                      //   Phantom   - 傀影与猩红血钻
//                      //   Mizuki    - 水月与深蓝之树
//                      //   Sami      - 探索者的银凇止境
//                      //   Sarkaz    - 萨卡兹的无终奇语
//                      //   JieGarden - 岁的界园志异
//                      //   BlackFlow - 黑流树海
//     "mode": int,     // 模式，可选，默认值 0
//                      //   0 - 刷分/奖励点数，尽可能稳定地打更多层数
//                      //   1 - 刷源石锭，第一层投资完就退出
//                      //   2 - 【已弃用】兼顾模式 0 与 1，投资过后再退出，没有投资就继续往后打
//                      //   3 - 开发中...
//                      //   4 - 凹开局，先在 0 难度下到达第三层后重开，再到指定难度下凹开局奖励，若不为热水壶或希望则回到 0 难度下重新来过；
//                      //       若在 Phantom 主题下则不切换难度，仅在当前难度下尝试到达第三层、重开、凹开局
//                      //   5 - 刷坍缩范式；仅适用于 Sami 主题；通过战斗漏怪等方式加快坍缩值积累，
//                      //       若遇到的第一个的坍缩范式在 expected_collapsal_paradigms 列表中则停止任务，否则重开
//                      //   6 - 刷月度小队蚊子腿，除了针对模式的适配以外和模式0相同
//                      //   7 - 刷深入调查蚊子腿，除了针对模式的适配以外和模式0相同
//     "squad": string,                // 开局分队名，可选，默认值 "指挥分队"；
//     "roles": string,                // 开局职业组，可选，默认值 "取长补短"；
//     "core_char": string,            // 开局干员名，可选；仅支持单个干员**中文名**，无论区服；若留空或设置为空字符串 "" 则根据练度自动选择
//     "use_support": bool,            // 开局干员是否为助战干员，可选，默认值 false
//     "use_nonfriend_support": bool,  // 是否可以是非好友助战干员，可选，默认值 false；仅在 use_support 为 true 时有效
//     "starts_count": int,                // 开始探索的次数，可选，默认3，最大3；达到后自动停止任务
//     "difficulty": int,                  // 指定难度等级，可选，默认值 0；仅适用于**除 Phantom 以外**的主题；
//                                         // 若未解锁难度，则会选择当前已解锁的最高难度
//     "stop_at_final_boss": bool,         // 是否在第 5 层险路恶敌节点前停止任务，可选，默认值 false；仅适用于**除 Phantom 以外**的主题
//     "stop_at_max_level": bool,          // 是否在肉鸽等级刷满后停止任务，可选，默认值 false
//     "investment_enabled": bool,         // 是否投资源石锭，可选，默认值 true
//     "investments_count": int,           // 投资源石锭的次数，可选，默认值 INT_MAX；达到后自动停止任务
//     "stop_when_investment_full": bool,  // 是否在投资到达上限后自动停止任务，可选，默认值 false
//     "investment_with_more_score": bool, // 是否在投资后尝试购物，可选，默认值 false；仅适用于模式 1
//     "start_with_elite_two": bool,       // 是否在凹开局的同时凹干员精二直升，可选，默认值 false；仅适用于模式 4
//     "only_start_with_elite_two": bool,  // 是否只凹开局干员精二直升而忽视其他开局条件，可选，默认值 false；
//                                         // 仅在模式为 4 且 start_with_elite_two 为 true 时有效
//     "refresh_trader_with_dice": bool,   // 是否用骰子刷新商店购买特殊商品，可选，默认值 false；仅适用于 Mizuki 主题，用于刷指路鳞
//     "first_floor_foldartal": string,    // 希望在第一层远见阶段得到的密文版，可选；仅适用于 Sami 主题，不限模式；若成功凹到则停止任务
//     "start_foldartal_list": [           // 凹开局时希望在开局奖励阶段得到的密文板，可选，默认值 []；仅在主题为 Sami 且模式为 4 时有效；
//       string,                           // 仅当开局拥有列表中所有的密文板时才算凹开局成功；
//       ...                               // 注意，此参数须与 “生活至上分队” 同时使用，其他分队在开局奖励阶段不会获得密文板；
//     ],
//     "use_foldartal": bool,                    // 是否使用密文板，模式 5 下默认值 false，其他模式下默认值 true；仅适用于 Sami 主题
//     "check_collapsal_paradigms": bool,        // 是否检测获取的坍缩范式，模式 5 下默认值 true，其他模式下默认值 false
//     "double_check_collapsal_paradigms": bool, // 是否执行坍缩范式检测防漏措施，默认值 true；
//                                               // 仅在主题为 Sami 且 check_collapsal_paradigms 为 true 时有效
//     "expected_collapsal_paradigms": [         // 希望触发的坍缩范式，默认值 ["目空一些, "睁眼瞎", "图像损坏", "一抹黑"]；
//         string,                               // 仅在主题为 Sami 且模式为 5 时有效
//         ...
//     ],
//     "monthly_squad_auto_iterate": bool,    //是否启动月度小队自动切换
//     "monthly_squad_check_comms": bool,     //是否将月度小队通信也作为切换依据
//     "deep_exploration_auto_iterate": bool, //是否启动深入调查自动切换
//     "collectible_mode_shopping": bool,  // 烧水是否启用购物, 默认值false
//     "collectible_mode_squad": string,   // 烧水时使用的分队, 默认与squad同步, 当squad为空字符串且未指定collectible_mode_squad值时为指挥分队
//     "collectible_mode_start_list": {    // 烧水期望奖励, 默认全false, 键范围：
//         "hot_water": bool,              // [hot_water: 热水, shield: 盾, ingot: 源石锭, hope: 希望, random: 随机奖励, key: 钥匙, dice: 骰子, ideas: 构想, ticket: 票券]
//         ...
//     },
//     "start_with_seed": bool,        // 使用种子刷钱，true 时有效
//                                     // 仅在 Sarkaz 主题，Investment 模式，“点刺成锭分队” or “后勤分队” 时可能为 true
//                                     // 使用固定种子
// }
// 黑流树海扩展参数以 v6.21.1 集成文档为准：blackflow_strategy、blackflow_cultivation_target、
// blackflow_start_animals、blackflow_start_rewards、blackflow_start_utopia、
// blackflow_start_reward_priority、collectible_mode_roles、collectible_mode_core_char_list。
type TaskRoguelike struct {
	Enable                        bool                 `json:"enable"                          v:"boolean" d:"true"                                           dc:"是否启用本任务"`
	Theme                         string               `json:"theme"                           v:"in:Phantom,Mizuki,Sami,Sarkaz,JieGarden,BlackFlow" d:"Phantom"   dc:"主题类型：Phantom/Mizuki/Sami/Sarkaz/JieGarden/BlackFlow（黑流树海）"`
	Mode                          int                  `json:"mode"                            v:"in:0,1,2,3,4,5,6,7,30001,30003" d:"0"                                   dc:"作战模式：0刷分,1刷锭,2已弃用,3开发中,4凹开局,5坍缩范式,6月度小队,7深入调查,30001刷襁褓动物,30003黑流树海刷开局"`
	Squad                         string               `json:"squad"                           v:""                     d:"指挥分队"                            dc:"开局分队"`
	Roles                         string               `json:"roles"                           v:""                     d:"取长补短"                            dc:"职业组合"`
	CoreChar                      string               `json:"core_char"                       v:""                     d:""                                  dc:"核心干员名称"`
	UseSupport                    bool                 `json:"use_support"                     v:""                     d:"false"                              dc:"使用助战"`
	UseNonFriendSupport           bool                 `json:"use_nonfriend_support"           v:""                     d:"false"                              dc:"允许非好友助战"`
	StartsCount                   int                  `json:"starts_count"                    v:"min:0|max:3"          d:"3"                                  dc:"开始探索次数上限，范围0-3，达到后自动停止"`
	Difficulty                    int                  `json:"difficulty"                     v:"min:-1"               d:"0"                                  dc:"指定难度等级，-1为不指定，未解锁时自动选择最高已解锁难度"`
	StopAtFinalBoss               bool                 `json:"stop_at_final_boss"              v:""                     d:"false"                              dc:"最终关卡前停止"`
	StopAtMaxLevel                bool                 `json:"stop_at_max_level"               v:""                     d:"false"                              dc:"满级后停止"`
	InvestmentEnabled             bool                 `json:"investment_enabled"              v:""                     d:"true"                               dc:"启用投资"`
	InvestmentsCount              int                  `json:"investments_count"               v:"min:0"                d:"2147483647"                         dc:"投资次数上限，达到后自动停止"`
	StopWhenInvestmentFull        bool                 `json:"stop_when_investment_full"        v:""                     d:"false"                              dc:"投资满后停止"`
	InvestmentWithMoreScore       bool                 `json:"investment_with_more_score"       v:""                     d:"false"                              dc:"投资后购物"`
	StartWithEliteTwo             bool                 `json:"start_with_elite_two"             v:""                     d:"false"                              dc:"凹精二开局"`
	OnlyStartWithEliteTwo         bool                 `json:"only_start_with_elite_two"         v:""                     d:"false"                              dc:"仅精二开局"`
	RefreshTraderWithDice         bool                 `json:"refresh_trader_with_dice"         v:""                     d:"false"                              dc:"骰子刷新商店"`
	FirstFloorFoldartal           string               `json:"first_floor_foldartal"            v:""                     d:""                                  dc:"首层密文板"`
	StartFoldartalList            []string             `json:"start_foldartal_list"             v:""                     d:""                                  dc:"开局密文板列表"`
	UseFoldartal                  bool                 `json:"use_foldartal"                    v:""                     d:"true"                               dc:"是否使用密文板（Sami主题；模式5建议为false，其余模式建议为true）"`
	CheckCollapsalParadigms       bool                 `json:"check_collapsal_paradigms"        v:""                     d:"false"                              dc:"是否检测坍缩范式（模式5建议为true）"`
	DoubleCheckCollapsalParadigms bool                 `json:"double_check_collapsal_paradigms" v:""                     d:"true"                               dc:"是否执行坍缩范式检测防漏（仅Sami且开启坍缩检测时有效）"`
	ExpectedCollapsalParadigms    []string             `json:"expected_collapsal_paradigms"     v:""                     d:"目空一些,睁眼瞎,图像损坏,一抹黑"     dc:"期望坍缩范式（仅Sami主题模式5有效）"`
	MonthlySquadAutoIterate       bool                 `json:"monthly_squad_auto_iterate"       v:""                     d:"false"                          dc:"月度小队自动迭代"`
	MonthlySquadCheckComms        bool                 `json:"monthly_squad_check_comms"        v:""                     d:"false"                             dc:"检查月度通信"`
	DeepExplorationAutoIterate    bool                 `json:"deep_exploration_auto_iterate"    v:""                     d:"false"                             dc:"深入调查自动迭代"`
	CollectibleModeShopping       bool                 `json:"collectible_mode_shopping"        v:""                     d:"false"                             dc:"烧水购物模式"`
	CollectibleModeSquad          string               `json:"collectible_mode_squad"           v:""                     d:""                                 dc:"烧水分队配置"`
	CollectibleModeStartList      map[string]bool      `json:"collectible_mode_start_list"      v:""                     d:""                                 dc:"烧水期望奖励配置"`
	StartWithSeed                 bool                 `json:"start_with_seed"                  v:""                     d:"false"                             dc:"使用种子刷钱"`
	BlackflowStrategy             string               `json:"blackflow_strategy,omitempty" v:"in:baby_animal,investment,burn_with_investment,burn,ending1,ending2,ending3" dc:"黑流树海路线策略，留空由核心推断；刷开局模式不读取"`
	BlackflowCultivationTarget    string               `json:"blackflow_cultivation_target,omitempty" v:"in:swaddled_cat,swaddled_dog,swaddled_cerberus,swaddled_feathered_serpent" dc:"黑流树海目标襁褓动物"`
	BlackflowStartAnimals         []string             `json:"blackflow_start_animals" v:"" dc:"黑流树海刷开局准备局动物，0至2项，羽蛇最多1项"`
	BlackflowStartRewards         int                  `json:"blackflow_start_rewards" v:"min:0|max:4095" dc:"黑流树海刷开局奖励位掩码，0至4095"`
	BlackflowStartUtopia          bool                 `json:"blackflow_start_utopia" dc:"黑流树海刷开局凹第一层实托邦，需难度6及以上"`
	BlackflowStartRewardPriority  []string             `json:"blackflow_start_reward_priority,omitempty" dc:"凹实托邦时的奖励优先级，不可为空或重复"`
	CollectibleModeRoles          string               `json:"collectible_mode_roles,omitempty" dc:"黑流树海刷开局准备局职业组，留空沿用roles"`
	CollectibleModeCoreCharList   []RoguelikeStartOper `json:"collectible_mode_core_char_list,omitempty" dc:"黑流树海刷开局准备局干员顺位，留空沿用核心干员"`
}

type RoguelikeStartOper struct {
	Name       string `json:"name"`
	UseSupport bool   `json:"use_support"`
}

func (t *TaskRoguelike) Serialize() string {
	params := *t
	// 已保存的旧定时任务也必须遵守新的探索次数上限。
	if params.StartsCount > 3 {
		params.StartsCount = 3
	}
	if params.BlackflowStartAnimals == nil {
		params.BlackflowStartAnimals = []string{}
	}
	b, err := json.Marshal(params)
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *TaskRoguelike) GetType() TaskType {
	return TaskType(Roguelike)
}

// 9. 自动抄作业（Copilot）
// 对应任务参数：
// {
//     "enable": bool,             // 是否启用本任务，可选，默认为 true
//     "filename": string,         // 作业 JSON 的文件路径，绝对、相对路径均可。不支持运行期设置
//     "formation": bool           // 是否进行 “快捷编队”，可选，默认否。不支持运行期设置
// }
type TaskCopilot struct {
	Enable    bool   `json:"enable"`
	Filename  string `json:"filename"`
	Formation bool   `json:"formation"`
}

func (t *TaskCopilot) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *TaskCopilot) GetType() TaskType {
	return TaskType(Copilot)
}

// 10. 自动抄保全作业（SSSCopilot）
// 对应任务参数：
//
//	{
//	    "enable": bool,
//	    "filename": string,
//	    "loop_times": int
//	}
type TaskSSSCopilot struct {
	Enable    bool   `json:"enable"`
	Filename  string `json:"filename"`
	LoopTimes int    `json:"loop_times"`
}

func (t *TaskSSSCopilot) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

// 11. 仓库识别（Depot）
// 对应任务参数：
//
//	{
//	    "enable": bool
//	}
type TaskDepot struct {
	Enable bool `json:"enable"`
}

func (t *TaskDepot) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

// 12. 可选的访问好友任务（Visit）
// 目前没有详细参数描述，这里仅给出一个简单示例
type TaskVisit struct {
	Enable bool `json:"enable"`
}

func (t *TaskVisit) Serialize() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}
