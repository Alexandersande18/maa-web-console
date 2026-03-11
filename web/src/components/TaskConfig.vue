<template>
  <div>
    <h3>任务配置</h3>
    <!-- 定时任务列表区域（无标题），滚动显示 -->
    <div class="box" id="cronjob-section">
      <div class="cronjob-scroll">
        <template v-if="cronJobs.length">
          <div
            class="cronjob-item"
            v-for="job in cronJobs"
            :key="job.id"
            :class="{ selected: currentCronJob && currentCronJob.id === job.id }"
            @click="selectCronJob(job)"
          >
            <div class="item-content">
              <div>
                <strong>{{ job.name }}</strong>
                <span v-if="job.isTemplate === 1" class="badge badge-template">模板任务</span>
                <span v-else class="badge badge-normal">常规任务</span>
              </div>
              <div class="config-summary">
                执行时间: {{ job.hour }}:
                {{ job.minute < 10 ? ('0' + job.minute) : job.minute }}
              </div>
              <div class="config-summary">
                上次执行: {{ formatTimestamp(job.lastExecutedAt) }}
              </div>
            </div>
            <!-- 每个表项右侧增加“配置”和“删除”按钮 -->
            <div class="item-actions">
              <van-button size="small" type="primary" @click.stop="openEditCronJobModal(job)">配 置</van-button>
              <van-button size="small" type="danger" @click.stop="confirmAndDelete(job)">删 除</van-button>
            </div>
          </div>
        </template>
        <div v-else class="empty-cronjob">没有定时任务</div>
      </div>
    </div>

    <!-- 原有任务配置区域 -->
    <div class="box" id="task-config-box">
      <div id="task-configurations">
        <div
          class="task-item"
          v-for="(task, key) in tasks"
          :key="key"
          :class="{ selected: task.enable }"
          @click="toggleTaskSelection(key)"
        >
          <strong>{{ task.label }}</strong>
          <van-button size="small" type="primary" class="config-button" @click.stop="openConfig(key)">配 置</van-button>
          <!-- 配置保存后显示 JSON 预览 -->
          <div class="config-summary" v-if="task.configSummary">
            {{ task.configSummary }}
          </div>
        </div>
      </div>
      <van-row justify="center" wrap="true">
        <van-space wrap>
          <van-button size="small" type="success" @click="startTasks">开 始</van-button>
          <van-button size="small" type="danger" @click="stopTasks">停 止</van-button>
          <!-- 添加定时任务按钮放到这组按钮的右侧 -->
          <van-button size="small" type="primary" @click="openAddCronJobFromTasks">添加为定时任务</van-button>
          <van-button v-if="currentCronJob" size="small" type="warning" @click="updateCronJob">保存定时任务</van-button>
        </van-space>
      </van-row>

      <!-- 任务配置模态框 -->
      <ConfigModal
        v-if="showModal"
        :taskKey="currentTaskKey"
        :taskConfig="tasks[currentTaskKey].config"
        :fieldConfig="fieldConfig[currentTaskKey] || {}"
        @save="saveTaskConfig"
        @close="closeModal"
      />
    </div>

    <!-- 定时任务模态框，用于添加或编辑（配置）定时任务 -->
    <CronJobModal
      v-if="showCronJobModal"
      :cronJobData="cronJobModalData"
      :isEditing="isEditingCronJob"
      @add="handleAddCronJob"
      @update="updateCronJobFromModal"
      @close="closeCronJobModal"
    />
  </div>
</template>

<script>
import ConfigModal from "./ConfigModal.vue";
import CronJobModal from "./CronJobModal.vue";

export default {
  name: "TaskConfig",
  components: { ConfigModal, CronJobModal },
  data() {
    return {
      tasks: {
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
      },
      fieldConfig: {
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
            hint: "JSON 格式，例如: {\"3\":540, \"4\":540, \"5\":540, \"6\":540}"
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
          series: { label: "连战次数", hint:"0~6, 0为自动切换为当前可用的最大次数, 如当前理智不够6次, 则选择最低可用次数", },
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
      },
      showModal: false,
      currentTaskKey: "",
      // 定时任务数据
      cronJobs: [],
      currentCronJob: null,
      showCronJobModal: false,
      // 用于定时任务模态框（添加或编辑）的数据及编辑标识
      cronJobModalData: null,
      isEditingCronJob: false
    };
  },
  mounted() {
    this.fetchCronJobs();
  },
  methods: {
    // 将服务器返回的 UTC+0 时间转换为 UTC+8 (北京) 时间进行显示
    formatTimestamp(timestamp) {
      if (!timestamp) return "未执行";
      // 如果返回的时间字符串没有时区标记，则追加 "Z" 表示UTC时间
      if (!timestamp.endsWith("Z")) {
        timestamp += "Z";
      }
      const date = new Date(timestamp);
      // 通过 toLocaleString 并指定 timeZone 为 Asia/Shanghai 来展示 UTC+8 时间
      return date.toLocaleString("zh-CN", { timeZone: "Asia/Shanghai" });
    },
    toggleTaskSelection(key) {
      this.tasks[key].enable = !this.tasks[key].enable;
    },
    openConfig(key) {
      this.currentTaskKey = key;
      this.showModal = true;
    },
    closeModal() {
      this.showModal = false;
    },
    saveTaskConfig(newConfig) {
      this.tasks[this.currentTaskKey].config = newConfig;
      this.tasks[this.currentTaskKey].configSummary = JSON.stringify(newConfig);
      this.closeModal();
    },
    /**
     * 配置转换方法
     */
    processTaskConfig(key, config) {
      switch (key) {
        case "taskRecruit":
          if (typeof config.select === "string") {
            config.select = config.select
              .split(",")
              .map((s) => parseInt(s.trim()))
              .filter((n) => !isNaN(n));
          }
          if (typeof config.confirm === "string") {
            config.confirm = config.confirm
              .split(",")
              .map((s) => parseInt(s.trim()))
              .filter((n) => !isNaN(n));
          }
          if (typeof config.first_tags === "string" && config.first_tags.trim() !== "") {
            config.first_tags = config.first_tags
              .split(",")
              .map((s) => s.trim())
              .filter((s) => s !== "");
          } else {
            config.first_tags = [];
          }
          if (typeof config.recruitment_time === "string" && config.recruitment_time.trim() !== "") {
            try {
              config.recruitment_time = JSON.parse(config.recruitment_time);
            } catch (e) {
              config.recruitment_time = {};
            }
          }
          break;
        case "taskInfrast":
          if (typeof config.facility === "string" && config.facility.trim() !== "") {
            config.facility = config.facility
              .split(",")
              .map((s) => s.trim())
              .filter((s) => s !== "");
          }
          break;
        case "taskFight":
          if (typeof config.drops === "string" && config.drops.trim() !== "") {
            try {
              config.drops = JSON.parse(config.drops);
            } catch (e) {
              config.drops = {};
            }
          }
          break;
        case "taskMall":
          if (typeof config.buy_first === "string" && config.buy_first.trim() !== "") {
            config.buy_first = config.buy_first
              .split(",")
              .map((s) => s.trim())
              .filter((s) => s !== "");
          } else {
            config.buy_first = [];
          }
          if (typeof config.blacklist === "string" && config.blacklist.trim() !== "") {
            config.blacklist = config.blacklist
              .split(",")
              .map((s) => s.trim())
              .filter((s) => s !== "");
          } else {
            config.blacklist = [];
          }
          break;
        default:
          break;
      }
      return config;
    },
    reverseProcessTaskConfig(key, config) {
      let cfg = { ...config };
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
    },
    // 获取每个启用任务的 payload（减少 startTasks 与 stopTasks 重复）
    getTasksPayload(enableFlag) {
      const payload = {};
      for (const key in this.tasks) {
        const task = this.tasks[key];
        if (task.enable) {
          let config = { ...task.config };
          config = this.processTaskConfig(key, config);
          config.enable = enableFlag;
          payload[key] = config;
        }
      }
      return payload;
    },
    async startTasks() {
      const body = this.getTasksPayload(true);
      const jwtToken = localStorage.getItem("jwtToken");
      try {
        const response = await fetch("/maasst/linkstart", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: jwtToken
          },
          body: JSON.stringify(body)
        });
        const result = await response.json();
        if (result.code === 0) {
          alert("Linkstart 启动成功");
        } else {
          alert(`操作失败：${result.message || "未知业务错误"}`);
        }
      } catch (e) {
        alert("请求异常：" + e.message);
      }
    },
    async stopTasks() {
      const body = this.getTasksPayload(false);
      const jwtToken = localStorage.getItem("jwtToken");
      try {
        const response = await fetch("/maasst/stoptask", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: jwtToken
          },
          body: JSON.stringify(body)
        });
        const result = await response.json();
        if (result.code === 0) {
          alert("停止任务成功");
        } else {
          alert(`停止任务失败：${result.message || "未知错误"}`);
        }
      } catch (e) {
        alert("请求异常：" + e.message);
      }
    },
    async fetchCronJobs() {
      const jwtToken = localStorage.getItem("jwtToken");
      try {
        const response = await fetch("/maasst/cronjob", {
          headers: { Authorization: jwtToken }
        });
        const result = await response.json();
        if (result.data && result.data.cronJobList) {
          this.cronJobs = result.data.cronJobList;
        }
      } catch (e) {
        alert("请求定时任务列表异常：" + e.message);
      }
    },
    selectCronJob(cronJob) {
      this.currentCronJob = {
        ...cronJob,
        isTemplateBool: cronJob.isTemplate === 1
      };

      // 用数组简化除了 taskAward 之外的重复代码
      const keysToProcess = ["taskRecruit", "taskInfrast", "taskFight", "taskMall"];
      keysToProcess.forEach((key) => {
        if (cronJob[key]) {
          const reversedConfig = this.reverseProcessTaskConfig(key, cronJob[key]);
          this.tasks[key].config = { ...reversedConfig };
          this.tasks[key].configSummary = JSON.stringify(
            this.processTaskConfig(key, { ...cronJob[key] })
          );
          if (cronJob[key].enable !== undefined) {
            this.tasks[key].enable = cronJob[key].enable;
          }
        }
      });
      if (cronJob.taskAward) {
        this.tasks.taskAward.config = { ...cronJob.taskAward };
        this.tasks.taskAward.configSummary = JSON.stringify(cronJob.taskAward);
        if (cronJob.taskAward.enable !== undefined) {
          this.tasks.taskAward.enable = cronJob.taskAward.enable;
        }
      }
    },
    // --------------------
    // 新增：从任务配置区域点击“添加为定时任务”，打开模态框（添加模式）
    openAddCronJobFromTasks() {
      this.isEditingCronJob = false;
      this.cronJobModalData = {
        name: "",
        hour: 0,
        minute: 0,
        isTemplate: 0
      };
      this.showCronJobModal = true;
    },
    // --------------------
    // 新增：点击表项中“配置”按钮，打开模态框（编辑模式）
    openEditCronJobModal(job) {
      this.isEditingCronJob = true;
      this.cronJobModalData = {
        id: job.id,
        name: job.name,
        hour: job.hour,
        minute: job.minute,
        isTemplate: job.isTemplate
      };
      this.showCronJobModal = true;
    },
    // 关闭定时任务模态框
    closeCronJobModal() {
      this.showCronJobModal = false;
      this.cronJobModalData = null;
      this.isEditingCronJob = false;
    },
    // --------------------
    // 新增：处理添加定时任务（数据来源由模态框返回）
    handleAddCronJob(data) {
      const payload = {
        name: data.name,
        hour: data.hour,
        minute: data.minute,
        isTemplate: data.isTemplate,
        taskRecruit: this.processTaskConfig("taskRecruit", { ...this.tasks.taskRecruit.config, enable: this.tasks.taskRecruit.enable }),
        taskInfrast: this.processTaskConfig("taskInfrast", { ...this.tasks.taskInfrast.config, enable: this.tasks.taskInfrast.enable }),
        taskFight: this.processTaskConfig("taskFight", { ...this.tasks.taskFight.config, enable: this.tasks.taskFight.enable }),
        taskMall: this.processTaskConfig("taskMall", { ...this.tasks.taskMall.config, enable: this.tasks.taskMall.enable }),
        taskAward: this.processTaskConfig("taskAward", { ...this.tasks.taskAward.config, enable: this.tasks.taskAward.enable })
      };
      const jwtToken = localStorage.getItem("jwtToken");
      fetch("/maasst/cronjob", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: jwtToken
        },
        body: JSON.stringify(payload)
      })
        .then(response => response.json())
        .then(result => {
          if (result.data && result.data.id !== undefined) {
            alert("添加定时任务成功");
            this.fetchCronJobs();
          } else {
            alert("添加定时任务失败");
          }
        })
        .catch(e => {
          alert("添加定时任务异常：" + e.message);
        });
      this.closeCronJobModal();
    },
    // --------------------
    // 编辑模式：点击“修改”按钮，更新定时任务（仅修改定时任务基本信息）
    updateCronJobFromModal(updatedCronJob) {
      if (!updatedCronJob || !updatedCronJob.id) {
        alert("定时任务不存在");
        return;
      }
      const payload = {
        id: updatedCronJob.id,
        name: updatedCronJob.name,
        hour: updatedCronJob.hour,
        minute: updatedCronJob.minute,
        isTemplate: updatedCronJob.isTemplate ? 1 : 0
      };
      const jwtToken = localStorage.getItem("jwtToken");
      fetch("/maasst/cronjob", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: jwtToken
        },
        body: JSON.stringify(payload)
      })
        .then(response => response.json())
        .then(result => {
          if (result.code === 0) {
            alert("定时任务修改成功");
            this.fetchCronJobs();
            this.closeCronJobModal();
          } else {
            alert("定时任务修改失败：" + (result.message || "未知错误"));
          }
        })
        .catch(e => {
          alert("更新定时任务异常：" + e.message);
        });
    },
    // --------------------
    // “保存定时任务”方法，将任务配置区域的各项配置同步更新到定时任务
    async updateCronJob() {
      if (!this.currentCronJob || !this.currentCronJob.id) {
        alert("没有选中的定时任务");
        return;
      }
      this.currentCronJob.isTemplate = this.currentCronJob.isTemplateBool ? 1 : 0;
      const payload = {
        id: this.currentCronJob.id,
        name: this.currentCronJob.name,
        hour: this.currentCronJob.hour,
        minute: this.currentCronJob.minute,
        isTemplate: this.currentCronJob.isTemplate,
        taskRecruit: this.processTaskConfig("taskRecruit", { ...this.tasks.taskRecruit.config, enable: this.tasks.taskRecruit.enable }),
        taskInfrast: this.processTaskConfig("taskInfrast", { ...this.tasks.taskInfrast.config, enable: this.tasks.taskInfrast.enable }),
        taskFight: this.processTaskConfig("taskFight", { ...this.tasks.taskFight.config, enable: this.tasks.taskFight.enable }),
        taskMall: this.processTaskConfig("taskMall", { ...this.tasks.taskMall.config, enable: this.tasks.taskMall.enable }),
        taskAward: this.processTaskConfig("taskAward", { ...this.tasks.taskAward.config, enable: this.tasks.taskAward.enable })
      };
      const jwtToken = localStorage.getItem("jwtToken");
      try {
        const response = await fetch("/maasst/cronjob", {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: jwtToken
          },
          body: JSON.stringify(payload)
        });
        const result = await response.json();
        if (result.code === 0) {
          alert("定时任务更新成功");
          this.fetchCronJobs();
        } else {
          alert("定时任务更新失败：" + (result.message || "未知错误"));
        }
      } catch (e) {
        alert("更新定时任务异常：" + e.message);
      }
    },
    // --------------------
    // 删除定时任务：点击删除按钮后确认
    confirmAndDelete(job) {
      if (confirm("确定要删除该定时任务吗？")) {
        const jwtToken = localStorage.getItem("jwtToken");
        fetch(`/maasst/cronjob?id=${job.id}`, {
          method: "DELETE",
          headers: { Authorization: jwtToken }
        })
          .then(response => response.json())
          .then(result => {
            if (result.code === 0) {
              alert("定时任务删除成功");
              if (this.currentCronJob && this.currentCronJob.id === job.id) {
                this.currentCronJob = null;
              }
              this.fetchCronJobs();
            } else {
              alert("删除失败：" + (result.message || "未知错误"));
            }
          })
          .catch(e => {
            alert("删除定时任务异常：" + e.message);
          });
      }
    }
  }
};
</script>

<style scoped>
.box {
  background-color: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.cronjob-scroll {
  max-height: 300px; /* 可根据实际需求调整高度 */
  overflow-y: auto;
  width: 100%;
  padding-right: 5px; /* 保留一点右边间隙，显示滚动条 */
}

.empty-cronjob {
  text-align: center;
  color: #999;
  padding: 10px;
}

.cronjob-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 10px;
  margin: 5px 0;
  cursor: pointer;
  transition: background-color 0.3s;
}

.cronjob-item:hover {
  background-color: #f0f0f0;
}

.cronjob-item.selected {
  background-color: #e6f7ff;
  border-color: #1890ff;
}

.item-content {
  flex: 1;
}

.item-actions {
  display: flex;
  align-items: center;
}

.item-actions button {
  margin-left: 5px;
}

.badge {
  display: inline-block;
  margin-left: 5px;
  padding: 2px 6px;
  font-size: 12px;
  border-radius: 4px;
  color: #fff;
}

.badge-template {
  background-color: #007bff;
}

.badge-normal {
  background-color: #6c757d;
}

.task-item {
  background-color: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 10px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: background-color 0.3s, border-color 0.3s;
  overflow: hidden;
}

.task-item:hover {
  background-color: #f9f9f9;
}

.task-item.selected {
  background-color: #e6f7ff;
  border-color: #1890ff;
}

.config-button {
  float: right;
}

.config-button:hover {
  background-color: #0056b3;
}

.config-summary {
  clear: both;
  font-size: 12px;
  color: #666;
  word-break: break-all;
}

.action-buttons {
  text-align: center;
  margin-top: 15px;
}

.btn {
  padding: 10px 15px;
  margin: 5px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  color: #fff;
}

.btn-success {
  background-color: #35dca5;
}

.btn-success:hover {
  background-color: #2ebd94;
}

.btn-danger {
  background-color: #dc3545;
}

.btn-danger:hover {
  background-color: #c82333;
}

.btn-info {
  background-color: #17a2b8;
}

.btn-info:hover {
  background-color: #138496;
}

.btn-add-cron {
  background-color: #007bff;
}

.btn-add-cron:hover {
  background-color: #0056b3;
}
</style>