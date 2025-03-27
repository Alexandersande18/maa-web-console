<template>
  <div>
    <h3>日志</h3>
    <div class="box">
      <div class="log-wrapper">
        <div class="log-container" ref="logContainer">
          <div
            v-for="(log, idx) in logs"
            :key="idx"
            :class="['log-entry', getLogLevelClass(log.level)]"
          >
            {{ formatDisplay(log) }}
          </div>
        </div>
        <div class="buttons">
          <!-- 使用 Vant 的 van-button 组件替换原生 button -->
          <van-button type="primary" size="small" @click="clearLogs">清空日志</van-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "LogDisplay",
  data() {
    return {
      logs: [],
      ws: null,
    };
  },
  created() {
    const token = localStorage.getItem("jwtToken") || "";
    if (!token) {
      console.warn("未找到 jwtToken，日志流将无法连接。");
      return;
    }

    const wsUrl = `ws://${window.location.host}/maasst/logs-ws?token=${token}`;
    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log("WebSocket 连接已建立");
    };

    this.ws.onmessage = (event) => {
      if (!event.data) return;
      try {
        const data = JSON.parse(event.data);
        if (data.event === "heartbeat") return;
        this.logs.push({
          timestamp: data.timestamp,
          message: data.message,
          level: data.level || "info",
        });
        if (this.logs.length > 300) {
          this.logs.shift();
        }
      } catch (e) {
        this.logs.push({
          timestamp: Date.now() / 1000,
          message: event.data,
          level: "info",
        });
      }
      this.scrollToBottom();
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket 发生错误:", error);
    };

    this.ws.onclose = (event) => {
      console.log("WebSocket 连接关闭:", event);
    };
  },
  beforeUnmount() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  },
  methods: {
    scrollToBottom() {
      this.$nextTick(() => {
        if (this.$refs.logContainer) {
          this.$refs.logContainer.scrollTop = this.$refs.logContainer.scrollHeight;
        }
      });
    },
    clearLogs() {
      this.logs = [];
    },
    formatDisplay(log) {
      const ts = log.timestamp ? new Date(log.timestamp * 1000) : new Date();
      const timeStr = ts.toLocaleTimeString("zh-CN", {
        hour12: false,
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
      });
      return `[${timeStr}] ${log.message}`;
    },
    getLogLevelClass(level) {
      const lvl = level ? level.toLowerCase() : "info";
      switch (lvl) {
        case "debug":
          return "log-level-debug";
        case "info":
          return "log-level-info";
        case "warning":
        case "warn":
          return "log-level-warning";
        case "error":
          return "log-level-error";
        default:
          return "";
      }
    },
  },
};
</script>

<style scoped>
.box {
  background-color: var(--white);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px var(--shadow-color);
}

.log-wrapper {
  margin-top: 10px;
}

.log-container {
  border: 1px solid var(--border-color);
  background: var(--white);
  height: 200px;
  padding: 5px;
  overflow-y: auto;
  font-family: monospace;
  font-size: 12px;
  line-height: 1.5;
}

.log-entry {
  padding: 2px 4px;
  margin: 2px 0;
  border-radius: 3px;
}

.log-level-debug {
  color: #999;
}
.log-level-info {
  color: #28a745;
}
.log-level-warning {
  color: #ffc107;
}
.log-level-error {
  color: #dc3545;
}

.buttons {
  text-align: right;
  margin-top: 10px;
}
</style>