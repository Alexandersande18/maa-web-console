import { TASK_TEMPLATE } from "./taskDefinitions.js";

const TASK_KEYS_WITH_COMPLEX_SERIALIZATION = [
  "taskRecruit",
  "taskInfrast",
  "taskFight",
  "taskMall",
  "taskRoguelike"
];

export function formatTimestamp(timestamp) {
  if (!timestamp) return "未执行";
  let ts = timestamp;
  if (!ts.endsWith("Z")) {
    ts += "Z";
  }
  const date = new Date(ts);
  return date.toLocaleString("zh-CN", { timeZone: "Asia/Shanghai" });
}

export function processTaskConfig(key, config) {
  const cfg = { ...config };
  switch (key) {
    case "taskRecruit":
      if (typeof cfg.select === "string") {
        cfg.select = cfg.select
          .split(",")
          .map((s) => parseInt(s.trim(), 10))
          .filter((n) => !isNaN(n));
      }
      if (typeof cfg.confirm === "string") {
        cfg.confirm = cfg.confirm
          .split(",")
          .map((s) => parseInt(s.trim(), 10))
          .filter((n) => !isNaN(n));
      }
      if (typeof cfg.first_tags === "string" && cfg.first_tags.trim() !== "") {
        cfg.first_tags = cfg.first_tags
          .split(",")
          .map((s) => s.trim())
          .filter((s) => s !== "");
      } else {
        cfg.first_tags = [];
      }
      if (typeof cfg.recruitment_time === "string" && cfg.recruitment_time.trim() !== "") {
        try {
          cfg.recruitment_time = JSON.parse(cfg.recruitment_time);
        } catch (e) {
          cfg.recruitment_time = {};
        }
      }
      break;
    case "taskInfrast":
      if (typeof cfg.facility === "string" && cfg.facility.trim() !== "") {
        cfg.facility = cfg.facility
          .split(",")
          .map((s) => s.trim())
          .filter((s) => s !== "");
      }
      break;
    case "taskFight":
      cfg.series = Math.max(0, Math.min(10, Number(cfg.series ?? 10)));
      if (typeof cfg.drops === "string" && cfg.drops.trim() !== "") {
        try {
          cfg.drops = JSON.parse(cfg.drops);
        } catch (e) {
          cfg.drops = {};
        }
      }
      break;
    case "taskMall":
      if (typeof cfg.buy_first === "string" && cfg.buy_first.trim() !== "") {
        cfg.buy_first = cfg.buy_first
          .split(",")
          .map((s) => s.trim())
          .filter((s) => s !== "");
      } else {
        cfg.buy_first = [];
      }
      if (typeof cfg.blacklist === "string" && cfg.blacklist.trim() !== "") {
        cfg.blacklist = cfg.blacklist
          .split(",")
          .map((s) => s.trim())
          .filter((s) => s !== "");
      } else {
        cfg.blacklist = [];
      }
      break;
    case "taskRoguelike":
      cfg.mode = Number(cfg.mode ?? 0);
      cfg.starts_count = Math.max(0, Math.min(3, Number(cfg.starts_count ?? 3)));
      for (const key of ["blackflow_start_animals", "blackflow_start_reward_priority"]) {
        if (typeof cfg[key] === "string") {
          cfg[key] = cfg[key].split(",").map((s) => s.trim()).filter(Boolean);
        }
      }
      cfg.blackflow_start_rewards = Number(cfg.blackflow_start_rewards ?? 0);
      if (typeof cfg.collectible_mode_core_char_list === "string") {
        cfg.collectible_mode_core_char_list = JSON.parse(cfg.collectible_mode_core_char_list || "[]");
      }
      if (!cfg.collectible_mode_core_char_list?.length) delete cfg.collectible_mode_core_char_list;
      if (!cfg.collectible_mode_roles) delete cfg.collectible_mode_roles;
      if (cfg.theme === "BlackFlow") {
        if (![0, 1, 30001, 30003].includes(cfg.mode)) cfg.mode = 0;
        cfg.start_with_elite_two = false;
        cfg.only_start_with_elite_two = false;
        if (cfg.mode !== 0) delete cfg.blackflow_strategy;
        if (cfg.mode === 30003) {
          const animals = cfg.blackflow_start_animals ?? [];
          const allowed = ["swaddled_cat", "swaddled_dog", "swaddled_cerberus", "swaddled_feathered_serpent"];
          if (!Array.isArray(animals) || animals.length > 2 || animals.some((s) => !allowed.includes(s)) || animals.filter((s) => s === "swaddled_feathered_serpent").length > 1) {
            throw new Error("准备局襁褓动物最多2项，羽蛇最多1项，请使用提示中的动物代码");
          }
          cfg.blackflow_start_animals = animals;
          if (!Number.isInteger(cfg.blackflow_start_rewards) || cfg.blackflow_start_rewards < 0 || cfg.blackflow_start_rewards > 4095) throw new Error("期望开局奖励须为0至4095的整数");
          if (cfg.blackflow_start_utopia) {
            const priority = cfg.blackflow_start_reward_priority;
            const rewards = ["unnumbered_item", "survey_advance", "space_lease", "regression_compensation", "forest_ride", "nest_parasitism", "pegasus", "hippogriff", "dragon", "hydra", "bai_ze", "golden_crow"];
            if (!Array.isArray(priority) || !priority.length || new Set(priority).size !== priority.length || priority.some((s) => !rewards.includes(s))) throw new Error("实托邦奖励优先级须填写有效奖励代码，不可为空或重复");
            if (Number(cfg.difficulty) < 6) throw new Error("凹第一层实托邦需要指定难度6及以上");
          }
        }
      } else if (cfg.mode >= 30000) {
        cfg.mode = 0;
      }
      if (typeof cfg.start_foldartal_list === "string" && cfg.start_foldartal_list.trim() !== "") {
        cfg.start_foldartal_list = cfg.start_foldartal_list
          .split(",")
          .map((s) => s.trim())
          .filter((s) => s !== "");
      } else if (!Array.isArray(cfg.start_foldartal_list)) {
        cfg.start_foldartal_list = [];
      }
      if (
        typeof cfg.expected_collapsal_paradigms === "string" &&
        cfg.expected_collapsal_paradigms.trim() !== ""
      ) {
        cfg.expected_collapsal_paradigms = cfg.expected_collapsal_paradigms
          .split(",")
          .map((s) => s.trim())
          .filter((s) => s !== "");
      } else if (!Array.isArray(cfg.expected_collapsal_paradigms)) {
        cfg.expected_collapsal_paradigms = [];
      }
      if (
        typeof cfg.collectible_mode_start_list === "string" &&
        cfg.collectible_mode_start_list.trim() !== ""
      ) {
        try {
          cfg.collectible_mode_start_list = JSON.parse(cfg.collectible_mode_start_list);
        } catch (e) {
          cfg.collectible_mode_start_list = {};
        }
      }
      break;
    default:
      break;
  }
  return cfg;
}

export function reverseProcessTaskConfig(key, config) {
  const cfg = { ...config };
  switch (key) {
    case "taskRecruit":
      if (Array.isArray(cfg.select)) {
        cfg.select = cfg.select.join(",");
      }
      if (Array.isArray(cfg.confirm)) {
        cfg.confirm = cfg.confirm.join(",");
      }
      if (Array.isArray(cfg.first_tags)) {
        cfg.first_tags = cfg.first_tags.join(",");
      }
      if (cfg.recruitment_time && typeof cfg.recruitment_time === "object") {
        try {
          cfg.recruitment_time = JSON.stringify(cfg.recruitment_time);
        } catch (e) {
          cfg.recruitment_time = "";
        }
      }
      break;
    case "taskInfrast":
      if (Array.isArray(cfg.facility)) {
        cfg.facility = cfg.facility.join(",");
      }
      break;
    case "taskFight":
      if (cfg.drops && typeof cfg.drops === "object") {
        try {
          cfg.drops = JSON.stringify(cfg.drops);
        } catch (e) {
          cfg.drops = "";
        }
      }
      break;
    case "taskMall":
      if (Array.isArray(cfg.buy_first)) {
        cfg.buy_first = cfg.buy_first.join(",");
      }
      if (Array.isArray(cfg.blacklist)) {
        cfg.blacklist = cfg.blacklist.join(",");
      }
      break;
    case "taskRoguelike":
      cfg.starts_count = Math.max(0, Math.min(3, Number(cfg.starts_count ?? 3)));
      for (const key of ["blackflow_start_animals", "blackflow_start_reward_priority"]) {
        if (Array.isArray(cfg[key])) cfg[key] = cfg[key].join(",");
      }
      if (Array.isArray(cfg.collectible_mode_core_char_list)) {
        cfg.collectible_mode_core_char_list = JSON.stringify(cfg.collectible_mode_core_char_list);
      }
      if (Array.isArray(cfg.start_foldartal_list)) {
        cfg.start_foldartal_list = cfg.start_foldartal_list.join(",");
      }
      if (Array.isArray(cfg.expected_collapsal_paradigms)) {
        cfg.expected_collapsal_paradigms = cfg.expected_collapsal_paradigms.join(",");
      }
      if (cfg.collectible_mode_start_list && typeof cfg.collectible_mode_start_list === "object") {
        try {
          cfg.collectible_mode_start_list = JSON.stringify(cfg.collectible_mode_start_list);
        } catch (e) {
          cfg.collectible_mode_start_list = "{}";
        }
      }
      break;
    default:
      break;
  }
  return cfg;
}

export function buildEnabledTasksPayload(tasks, enableFlag) {
  const payload = {};
  for (const key in tasks) {
    const task = tasks[key];
    if (!task.enable) continue;
    const config = processTaskConfig(key, { ...task.config });
    payload[key] = { ...config, enable: enableFlag };
  }
  return payload;
}

export function buildCronJobPayloadFromTasks(tasks) {
  const payload = {};
  Object.keys(tasks).forEach((key) => {
    payload[key] = processTaskConfig(key, {
      ...tasks[key].config,
      enable: tasks[key].enable
    });
  });
  return payload;
}

export function syncTasksFromCronJob(tasks, cronJob) {
  const nextTasks = tasks;
  TASK_KEYS_WITH_COMPLEX_SERIALIZATION.forEach((key) => {
    if (!cronJob[key]) return;
    const reversed = reverseProcessTaskConfig(key, { ...TASK_TEMPLATE[key].config, ...cronJob[key] });
    nextTasks[key].config = { ...reversed };
    nextTasks[key].configSummary = JSON.stringify(
      processTaskConfig(key, { ...cronJob[key] })
    );
    if (cronJob[key].enable !== undefined) {
      nextTasks[key].enable = cronJob[key].enable;
    }
  });

  if (cronJob.taskAward) {
    nextTasks.taskAward.config = { ...cronJob.taskAward };
    nextTasks.taskAward.configSummary = JSON.stringify(cronJob.taskAward);
    if (cronJob.taskAward.enable !== undefined) {
      nextTasks.taskAward.enable = cronJob.taskAward.enable;
    }
  }

  return nextTasks;
}

