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
      series: 6,
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
    threshold: { label: "心情阈值" },
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
      hint: "2147483647表示无限，刷到理智不够"
    },
    series: {
      label: "连战次数",
      hint: "0~6, 0为自动切换为当前可用的最大次数, 如当前理智不够6次, 则选择最低可用次数"
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
  }
};

export function createDefaultTasks() {
  return JSON.parse(JSON.stringify(TASK_TEMPLATE));
}

