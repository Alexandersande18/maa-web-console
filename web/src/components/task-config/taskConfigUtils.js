const TASK_KEYS_WITH_COMPLEX_SERIALIZATION = [
  "taskRecruit",
  "taskInfrast",
  "taskFight",
  "taskMall"
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
    const reversed = reverseProcessTaskConfig(key, cronJob[key]);
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

