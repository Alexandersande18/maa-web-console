<template>
  <div>
    <h2>任务配置</h2>
    <div class="box" id="task-config-box">
      <div id="task-configurations">
        <div
          class="task-item"
          v-for="(task, key) in tasks"
          :key="key"
          :class="{ selected: task.selected }"
          @click="toggleTaskSelection(key)"
        >
          <strong>{{ task.label }}</strong>
          <button class="config-button" @click.stop="openConfig(key)">配置</button>
          <!-- 仅当配置保存后才显示 JSON 预览 -->
          <div class="config-summary" v-if="task.configSummary">
            {{ task.configSummary }}
          </div>
        </div>
      </div>
      <div class="action-buttons">
        <button class="btn btn-success" @click="startTasks">开始任务</button>
        <button class="btn btn-danger" @click="stopTasks">停止任务</button>
      </div>

      <!-- 引入任务配置模态框 -->
      <ConfigModal
        v-if="showModal"
        :taskKey="currentTaskKey"
        :taskConfig="tasks[currentTaskKey].config"
        :fieldConfig="fieldConfig[currentTaskKey] || {}"
        @save="saveTaskConfig"
        @close="closeModal"
      />
    </div>
  </div>
</template>

<script>
import ConfigModal from "./ConfigModal.vue";

export default {
  name: "TaskConfig",
  components: { ConfigModal },
  data() { 
    return {
      tasks: {
        taskRecruit: {
          label: "自动公招",
          selected: false,
          config: {
            // 默认配置
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
          // 默认不显示预览
          configSummary: ""
        },
        taskInfrast: {
          label: "基建换班",
          selected: false,
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
          selected: false,
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
          selected: false,
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
          selected: false,
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
      // 字段中文映射及其他配置，预留了每个任务字段的中文 label 及提示
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
          replenish: { label: "贸易站源石碎片是否自动补货" },
          dorm_notstationed_enabled: { label: "是否启用宿舍未进驻选项" },
          dorm_trust_enabled: { label: "是否将宿舍剩余位置填入信赖未满干员" },
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
          stone: { label: "最大吃源石数量" },
          times: {
            label: "执行次数",
            hint: "2147483647表示无限，刷到理智不够"
          },
          series: { label: "连战次数" },
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
            label: "是否在信用溢出时无视黑名单"
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
      currentTaskKey: ""
    };
  },
  methods: {
    toggleTaskSelection(key) {
      this.tasks[key].selected = !this.tasks[key].selected;
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
      // 只在用户保存配置后更新预览
      this.tasks[this.currentTaskKey].configSummary = JSON.stringify(newConfig);
      this.closeModal();
    },
    /**
     * 根据不同任务，对部分配置进行转换（如字符串拆分、JSON解析等）
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
          if (
            typeof config.first_tags === "string" &&
            config.first_tags.trim() !== ""
          ) {
            config.first_tags = config.first_tags
              .split(",")
              .map((s) => s.trim())
              .filter((s) => s !== "");
          } else {
            config.first_tags = [];
          }
          if (
            typeof config.recruitment_time === "string" &&
            config.recruitment_time.trim() !== ""
          ) {
            try {
              config.recruitment_time = JSON.parse(config.recruitment_time);
            } catch (e) {
              config.recruitment_time = {};
            }
          }
          break;
        case "taskInfrast":
          if (
            typeof config.facility === "string" &&
            config.facility.trim() !== ""
          ) {
            config.facility = config.facility
              .split(",")
              .map((s) => s.trim())
              .filter((s) => s !== "");
          }
          break;
        case "taskFight":
          if (
            typeof config.drops === "string" &&
            config.drops.trim() !== ""
          ) {
            try {
              config.drops = JSON.parse(config.drops);
            } catch (e) {
              config.drops = {};
            }
          }
          break;
        case "taskMall":
          if (
            typeof config.buy_first === "string" &&
            config.buy_first.trim() !== ""
          ) {
            config.buy_first = config.buy_first
              .split(",")
              .map((s) => s.trim())
              .filter((s) => s !== "");
          } else {
            config.buy_first = [];
          }
          if (
            typeof config.blacklist === "string" &&
            config.blacklist.trim() !== ""
          ) {
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
    async startTasks() {
      const body = {};
      for (const key in this.tasks) {
        const task = this.tasks[key];
        if (task.selected) {
          let config = { ...task.config };
          config = this.processTaskConfig(key, config);
          // 启动任务时设置启用标记
          config.enable = true;
          body[key] = config;
        }
      }
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
      const body = {};
      for (const key in this.tasks) {
        const task = this.tasks[key];
        if (task.selected) {
          let config = { ...task.config };
          config = this.processTaskConfig(key, config);
          config.enable = false;
          body[key] = config;
        }
      }
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
  background-color: #007bff;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 8px 14px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.config-button:hover {
  background-color: #0056b3;
}

.config-summary {
  clear: both;
  margin-top: 10px;
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
</style>