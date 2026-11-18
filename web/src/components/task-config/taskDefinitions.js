export const TASK_TEMPLATE = {
  taskRecruit: {
    label: "自动公招",
    enable: false,
    config: {
      refresh: true,
      select: "3,4,5",
      confirm: "3,4,5",
      first_tags: "",
      extra_tags_mode: 0,
      times: 4,
      set_time: true,
      expedite: false,
      expedite_times: 0,
      skip_robot: true,
      recruitment_time: '{"3":540,"4":540,"5":540,"6":540}',
      report_to_penguin: false,
      penguin_id: "",
      report_to_yituliu: false,
      yituliu_id: "",
      server: "CN"
    },
    configSummary: ""
  },
  taskInfrast: {
    label: "基建换班",
    enable: false,
    config: {
      mode: 0,
      facility: "Mfg,Trade,Power,Control,Reception,Office,Dorm",
      drones: "Money",
      threshold: 0.3,
      replenish: false,
      dorm_notstationed_enabled: true,
      dorm_trust_enabled: true,
      filename: "",
      plan_index: 0
    },
    configSummary: ""
  },
  taskFight: {
    label: "刷理智",
    enable: false,
    config: {
      stage: "",
      medicine: 0,
      expiring_medicine: 0,
      stone: 0,
      times: 2147483647,
      series: 10,
      drops: "{}",
      report_to_penguin: false,
      penguin_id: "",
      server: "CN",
      client_type: "Official",
      DrGrandet: false
    },
    configSummary: ""
  },
  taskMall: {
    label: "信用商店购物",
    enable: false,
    config: {
      shopping: true,
      visit_friends: true,
      buy_first: "招聘许可,龙门币",
      blacklist: "",
      force_shopping_if_credit_full: true,
      only_buy_discount: false,
      reserve_max_credit: false
    },
    configSummary: ""
  },
  taskAward: {
    label: "领取奖励",
    enable: false,
    config: {
      award: true,
      mail: true,
      recruit: true,
      orundum: true,
      mining: true,
      specialaccess: true
    },
    configSummary: ""
  },
  taskRoguelike: {
    label: "刷肉鸽",
    enable: false,
    config: {
      theme: "JieGarden",
      mode: 0,
      squad: "指挥分队",
      roles: "取长补短",
      core_char: "",
      use_support: false,
      use_nonfriend_support: false,
      starts_count: 3,
      difficulty: 18,
      stop_at_final_boss: false,
      stop_at_max_level: false,
      investment_enabled: true,
      investments_count: 2147483647,
      stop_when_investment_full: false,
      investment_with_more_score: false,
      start_with_elite_two: false,
      only_start_with_elite_two: false,
      refresh_trader_with_dice: false,
      first_floor_foldartal: "",
      start_foldartal_list: "",
      use_foldartal: true,
      check_collapsal_paradigms: false,
      double_check_collapsal_paradigms: true,
      expected_collapsal_paradigms: "目空一些,睁眼瞎,图像损坏,一抹黑",
      monthly_squad_auto_iterate: true,
      monthly_squad_check_comms: false,
      deep_exploration_auto_iterate: true,
      collectible_mode_shopping: false,
      collectible_mode_squad: "",
      collectible_mode_start_list: "{}",
      start_with_seed: false,
      blackflow_strategy: "",
      blackflow_cultivation_target: "swaddled_cat",
      blackflow_start_animals: "",
      blackflow_start_rewards: 0,
      blackflow_start_utopia: false,
      blackflow_start_reward_priority: "unnumbered_item,survey_advance,space_lease,regression_compensation,forest_ride,nest_parasitism",
      collectible_mode_roles: "",
      collectible_mode_core_char_list: "[]"
    },
    configSummary: ""
  },
  taskActivity: {
    label: "活动刷取",
    enable: false,
    config: { activity_id: "spa-regulations", current_page: false },
    configSummary: ""
  },
  taskReclamation: {
    label: "生息演算",
    enable: false,
    config: {
      theme: "Tales",
      mode: 0,
      tools_to_craft: "荧光棒",
      increment_mode: 0,
      num_craft_batches: 1,
      clear_store: false,
      current_page: false
    },
    configSummary: ""
  }
};

export const FIELD_CONFIG = {
  taskRecruit: {
    refresh: { label: "如果是三星Tag则继续尝试刷新" },
    select: {
      label: "会自动点击的点击标签等级",
      hint: "多个标签用逗号分隔，例如:3,4"
    },
    confirm: {
      label: "会自动确认的标签等级",
      hint: "多个标签用逗号分隔"
    },
    first_tags: {
      label: "偏好标签",
      hint: "仅在标签等级为3时有效，多个标签用逗号分隔"
    },
    extra_tags_mode: {
      label: "额外标签模式",
      type: "select",
      options: {
        "0": "默认行为（0）",
        "1": "选3个标签（1）",
        "2": "组合高星标签（2）"
      }
    },
    times: { label: "执行招募的次数", hint: "0表示仅公招计算" },
    set_time: { label: "设置招募时间" },
    expedite: { label: "使用加急许可" },
    expedite_times: {
      label: "加急次数",
      hint: "仅在使用加急许可时有效"
    },
    skip_robot: { label: "跳过机器人标签" },
    recruitment_time: {
      label: "招募时限（分钟）",
      hint: 'JSON 格式，例如: {"3":540, "4":540, "5":540, "6":540}'
    },
    report_to_penguin: { label: "汇报企鹅物流数据" },
    penguin_id: { label: "企鹅物流数据ID" },
    report_to_yituliu: { label: "汇报一图流数据" },
    yituliu_id: { label: "一图流数据ID" },
    server: {
      label: "服务器",
      type: "select",
      options: { CN: "国服", US: "美服", JP: "日服", KR: "韩服" }
    }
  },
  taskInfrast: {
    mode: {
      label: "模式",
      type: "select",
      options: { "0": "智能换班（0）", "10000": "自定义方案（10000）" }
    },
    facility: {
      label: "换班设施",
      hint: "逗号分隔，例如：Mfg,Trade"
    },
    drones: {
      label: "无人机用途",
      type: "select",
      options: {
        _NotUse: "不使用",
        Money: "龙门币",
        SyntheticJade: "合成玉"
      }
    },
    threshold: { label: "心情阈值", min: 0, max: 1, step: 0.1, integer: false },
    replenish: { label: "贸易站自动补货" },
    dorm_notstationed_enabled: { label: "启用未进驻宿舍" },
    dorm_trust_enabled: { label: "填入信赖未满干员" },
    filename: { label: "配置文件路径" },
    plan_index: { label: "方案序号" }
  },
  taskFight: {
    stage: {
      label: "关卡名称",
      hint: "关卡名，为空则识别当前/上次的关卡。支持全部主线关卡，如 1-7、S3-2等，可在关卡结尾输入Normal/Hard表示需要切换标准与磨难难度，剿灭作战输入Annihilation"
    },
    medicine: { label: "最大使用理智药数量" },
    expiring_medicine: { label: "最大使用临期理智药数量" },
    stone: { label: "最大吃石头数量" },
    times: {
      label: "执行次数",
      min: -1,
      hint: "2147483647表示无限，刷到理智不够"
    },
    series: {
      label: "连战次数",
      min: 0,
      max: 10,
      hint: "0~10，0为自动选择当前可用的最大连战次数"
    },
    drops: { label: "掉落限制", hint: "JSON格式，例如：{\"30011\":5}" },
    report_to_penguin: { label: "上报企鹅物流" },
    penguin_id: { label: "企鹅物流ID" },
    server: {
      label: "服务器",
      type: "select",
      options: { CN: "国服", US: "美服", JP: "日服", KR: "韩服" }
    },
    client_type: {
      label: "客户端",
      type: "select",
      options: { Official: "官服", Bilibili: "B服", txwy: "渠道服" }
    },
    DrGrandet: { label: "智能碎石" }
  },
  taskMall: {
    shopping: { label: "执行购物" },
    visit_friends: { label: "访问好友" },
    buy_first: {
      label: "优先购买",
      hint: "逗号分隔商品"
    },
    blacklist: { label: "黑名单", hint: "逗号分隔商品" },
    force_shopping_if_credit_full: {
      label: "信用溢出时无视黑名单"
    },
    only_buy_discount: { label: "仅购折扣" },
    reserve_max_credit: { label: "保留信用" }
  },
  taskAward: {
    award: { label: "任务奖励" },
    mail: { label: "邮件奖励" },
    recruit: { label: "公招奖励" },
    orundum: { label: "合成玉" },
    mining: { label: "开采许可" },
    specialaccess: { label: "月卡奖励" }
  },
  taskRoguelike: {
    theme: {
      label: "主题",
      type: "select",
      options: {
        Phantom: "傀影与猩红孤钻",
        Mizuki: "水月与深蓝之树",
        Sami: "探索者的银凇止境",
        Sarkaz: "萨卡兹的无终奇语",
        JieGarden: "界园",
        BlackFlow: "黑流树海"
      }
    },
    mode: {
      label: "作战模式",
      type: "select",
      options: {
        "0": "0 - 刷奖励等级（尽可能稳定打更多层数）",
        "1": "1 - 刷锭",
        "2": "2 - 已弃用",
        "3": "3 - 开发中",
        "4": "4 - 凹开局",
        "5": "5 - 坍缩范式",
        "6": "6 - 月度小队",
        "7": "7 - 深入调查"
      }
    },
    squad: { label: "开局分队" },
    roles: { label: "职业组合" },
    core_char: { label: "核心干员名称" },
    use_support: { label: "使用助战" },
    use_nonfriend_support: { label: "允许非好友助战" },
    starts_count: { label: "开始探索次数上限", min: 0, max: 10, hint: "最多10次，达到后自动停止" },
    difficulty: { label: "指定难度等级", min: -1, hint: "-1表示不指定难度；黑流树海凹实托邦需6及以上" },
    stop_at_final_boss: { label: "最终关卡前停止" },
    stop_at_max_level: { label: "满级后停止" },
    investment_enabled: { label: "投资源石锭" },
    investments_count: { label: "源石锭投资次数上限" },
    stop_when_investment_full: { label: "源石锭投资满后停止" },
    investment_with_more_score: { label: "投资后购物" },
    start_with_elite_two: { label: "凹精二开局" },
    only_start_with_elite_two: { label: "仅精二开局" },
    refresh_trader_with_dice: { label: "骰子刷新商店" },
    first_floor_foldartal: { label: "首层密文板" },
    start_foldartal_list: { label: "开局密文板列表", hint: "逗号分隔，例如：板A,板B" },
    use_foldartal: { label: "使用密文板" },
    check_collapsal_paradigms: { label: "检测坍缩范式" },
    double_check_collapsal_paradigms: { label: "坍缩范式防漏检测" },
    expected_collapsal_paradigms: { label: "期望坍缩范式", hint: "逗号分隔" },
    monthly_squad_auto_iterate: { label: "月度小队自动切换" },
    monthly_squad_check_comms: { label: "检查月度通信" },
    deep_exploration_auto_iterate: { label: "深入调查自动切换" },
    collectible_mode_shopping: { label: "烧水模式购物" },
    collectible_mode_squad: { label: "准备局分队", hint: "刷开局使用，留空沿用开局分队；机械师快速飞三层可使用特勤分队或堡垒战术分队" },
    collectible_mode_start_list: {
      label: "烧水期望奖励",
      hint: 'JSON 格式，例如：{"hot_water":true,"ingot":true}'
    },
    start_with_seed: { label: "使用种子刷钱" },
    blackflow_strategy: {
      label: "黑流树海路线策略",
      options: {
        "": "自动（按作战模式和投资开关推断）",
        investment: "第一层投资后退出",
        burn_with_investment: "投资后快速到第三层重开",
        burn: "快速到第三层重开",
        baby_animal: "培育襁褓动物",
        ending1: "第一结局",
        ending2: "第二结局",
        ending3: "第三结局"
      },
      hint: "刷等级模式可指定路线；结局是否开放以游戏内为准"
    },
    blackflow_cultivation_target: {
      label: "目标襁褓动物",
      options: {
        swaddled_cat: "襁褓猫",
        swaddled_dog: "襁褓狗",
        swaddled_cerberus: "襁褓三头犬",
        swaddled_feathered_serpent: "襁褓羽蛇"
      }
    },
    blackflow_start_animals: {
      label: "准备局襁褓动物",
      hint: "逗号分隔，最多2项，可重复，羽蛇最多1项；可填 swaddled_cat（猫）、swaddled_dog（狗）、swaddled_cerberus（三头犬）、swaddled_feathered_serpent（羽蛇）；留空只要求到第三层"
    },
    blackflow_start_rewards: {
      label: "期望开局奖励",
      min: 0,
      max: 4095,
      hint: "命中任一即停，多个数值相加：1未编号物、2调查预付款、4空间租赁、8退行补偿、16林间代步、32巢寄生；带羽蛇可加64天马、128骏鹰、256巨龙、512九头蛇、1024白泽、2048金乌；0不限制"
    },
    blackflow_start_utopia: { label: "凹第一层实托邦", hint: "目标局难度需6及以上；未开启时命中奖励后保留该局，不选择奖励" },
    blackflow_start_reward_priority: {
      label: "实托邦开局奖励优先级",
      hint: "逗号分隔，不可重复：unnumbered_item、survey_advance、space_lease、regression_compensation、forest_ride、nest_parasitism、pegasus、hippogriff、dragon、hydra、bai_ze、golden_crow"
    },
    collectible_mode_roles: { label: "准备局职业组", hint: "留空沿用职业组合；机械师快速飞三层使用稳扎稳打或坚不可摧" },
    collectible_mode_core_char_list: {
      label: "准备局开局干员",
      hint: 'JSON数组，例如：[{"name":"机械师","use_support":false}]；留空沿用核心干员，快速飞三层需要自己的满天赋机械师'
    }
  },
  taskActivity: {
    activity_id: {
      label: "活动",
      options: { "spa-regulations": "卫戍协议三期 · 规程" },
      hint: "使用卫戍协议默认主题。已有模拟进度请先在游戏中处理。"
    },
    current_page: {
      label: "从当前游戏页面开始",
      hint: "手动进入有「极速模拟」的页面后，单独手动运行。"
    }
  },
  taskReclamation: {
    theme: {
      label: "主题",
      options: { Tales: "沙洲遗闻", RelaunchAnchor: "重启锚点" }
    },
    mode: { label: "刷取模式" },
    tools_to_craft: { label: "制造道具", hint: "用逗号分隔，留空默认荧光棒。" },
    increment_mode: { label: "制造加量方式", options: { 0: "连点", 1: "长按" } },
    num_craft_batches: { label: "每次最大制造轮数", min: 1, max: 2147483647 },
    clear_store: { label: "刷完后兑换生息演算商店" },
    current_page: {
      label: "从当前游戏页面开始",
      hint: "手动进入对应主题或关卡后，单独手动运行。"
    }
  }
};

const RELAUNCH_ANCHOR_STAGES = {
  16: {
    name: "RA-1", tokens: 159, duration: "2 分 10 秒", operator: "无需指定干员",
    steps: ["推进主线并通关 RA-1。"],
    flow: "通过精耕细作、建设、交付资源和结算循环获取收益。"
  },
  48: {
    name: "RA-4", tokens: 417, duration: "1 分 40 秒", operator: "维什戴尔，可使用助战",
    steps: [
      "推进主线并通关 RA-4。",
      "解锁「策略筹划经营」。",
      "使用自己的维什戴尔时，先手动配队。前五位选费用低于维什戴尔的干员，最后把维什戴尔编入第六位。确认招募后，放弃本次建设并回到开始界面。",
      "使用助战时，确认维什戴尔出现在狙击助战首页，可优先使用挚友的助战。按上述要求编入前五位，第六位临时借用维什戴尔，确认招募再放弃建设。回到开始界面时确认前五位已保存，第六位为空。"
    ],
    flow: "用筹划经营提供的赤金解锁区域，再由维什戴尔击败首领。"
  },
  32: {
    name: "RA-15", tokens: 500, duration: "3 分钟", operator: "圣聆初雪，可使用助战",
    steps: [
      "推进主线并通关 RA-15。",
      "使用自己的圣聆初雪时，先手动配队。前五位编入任意先锋，无练度要求。最后把圣聆初雪编入第六位，保存编队后退出关卡。",
      "使用助战时，确认圣聆初雪出现在术士助战首页，可优先使用挚友的助战。"
    ],
    flow: "由圣聆初雪完成击杀 60 名敌人的目标。"
  }
};

// In-game prerequisites follow the MAA Toolbox guidance for Reclamation Algorithm.
export function getReclamationStageGuide(config) {
  const stage = config.theme === "RelaunchAnchor" && RELAUNCH_ANCHOR_STAGES[Number(config.mode)];
  if (!stage) return null;
  return [
    { label: "收益参考", text: `每轮约 ${stage.tokens} 代币和统筹点数，耗时约 ${stage.duration}。` },
    { label: "干员要求", text: stage.operator },
    { label: "前置步骤", steps: stage.steps },
    { label: "启动位置", text: `选择从当前游戏页面开始时，在大地图打开 ${stage.name}，看到右下角的「开启建设」后再启动。默认由 Web 任务定位关卡。` },
    { label: "注意", text: "若已解锁开局携带额外物品的科技，请先拆除会提供这些物品的基地设施，例如食品供给站、饮品供给站和兽栏。" },
    { label: "任务流程", text: stage.flow }
  ];
}

export function getNumberFieldOptions(taskKey, key) {
  if (typeof TASK_TEMPLATE[taskKey]?.config[key] !== "number") return null;
  const field = FIELD_CONFIG[taskKey]?.[key] || {};
  return {
    min: field.min ?? 0,
    max: field.max ?? Number.MAX_SAFE_INTEGER,
    step: field.step ?? 1,
    integer: field.integer ?? true
  };
}

export function getFieldOptions(taskKey, key, config, fieldConfig = FIELD_CONFIG[taskKey]) {
  if (taskKey === "taskReclamation" && key === "mode") {
    return config.theme === "RelaunchAnchor"
      ? Object.fromEntries(Object.entries(RELAUNCH_ANCHOR_STAGES).map(([mode, stage]) => [mode, stage.name]))
      : { 0: "无存档刷点数", 1: "有存档制造刷点数" };
  }
  if (taskKey === "taskRoguelike" && key === "mode" && config.theme === "BlackFlow") {
    return { 0: "刷奖励等级", 1: "刷源石锭", 30001: "刷襁褓动物", 30003: "刷开局" };
  }
  const options = fieldConfig?.[key]?.options;
  if (taskKey === "taskRoguelike" && key === "mode") {
    return Object.fromEntries(Object.entries(options).filter(([mode]) => Number(mode) < 30000));
  }
  return options;
}

export function isTaskFieldVisible(taskKey, key, config) {
  if (taskKey === "taskReclamation") {
    if (key === "clear_store") return config.theme === "Tales" && Number(config.mode) === 0;
    if (["tools_to_craft", "increment_mode", "num_craft_batches"].includes(key)) {
      return config.theme === "Tales" && Number(config.mode) === 1;
    }
    return true;
  }
  if (taskKey !== "taskRoguelike") return true;
  const blackflow = config.theme === "BlackFlow";
  const startRun = Number(config.mode) === 30003;
  if (key === "blackflow_strategy") return blackflow && Number(config.mode) === 0;
  if (key === "blackflow_cultivation_target") {
    return blackflow && (Number(config.mode) === 30001 || (Number(config.mode) === 0 && config.blackflow_strategy === "baby_animal"));
  }
  if (key.startsWith("blackflow_start_")) {
    return blackflow && startRun && (key !== "blackflow_start_reward_priority" || config.blackflow_start_utopia);
  }
  if (["collectible_mode_squad", "collectible_mode_roles", "collectible_mode_core_char_list"].includes(key)) return blackflow && startRun;
  return !blackflow || !["start_with_elite_two", "only_start_with_elite_two", "refresh_trader_with_dice", "monthly_squad_auto_iterate", "monthly_squad_check_comms", "deep_exploration_auto_iterate", "start_with_seed"].includes(key);
}

export function createDefaultTasks() {
  return JSON.parse(JSON.stringify(TASK_TEMPLATE));
}

