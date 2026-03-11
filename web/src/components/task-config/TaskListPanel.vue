<template>
  <div class="box" id="task-config-box">
    <div id="task-configurations">
      <div
        class="task-item"
        v-for="(task, key) in tasks"
        :key="key"
        :class="{ selected: task.enable }"
        @click="$emit('toggle-task', key)"
      >
        <strong>{{ task.label }}</strong>
        <van-button size="small" type="primary" class="config-button" @click.stop="$emit('open-config', key)">配 置</van-button>
        <div class="config-summary" v-if="task.configSummary">
          {{ task.configSummary }}
        </div>
      </div>
    </div>
    <van-row justify="center" wrap="true">
      <van-space wrap>
        <van-button size="small" type="success" @click="$emit('start-tasks')">开 始</van-button>
        <van-button size="small" type="danger" @click="$emit('stop-tasks')">停 止</van-button>
        <van-button size="small" type="primary" @click="$emit('add-cron-job')">添加为定时任务</van-button>
        <van-button v-if="currentCronJob" size="small" type="warning" @click="$emit('save-cron-job')">保存定时任务</van-button>
      </van-space>
    </van-row>
  </div>
</template>

<script>
export default {
  name: "TaskListPanel",
  props: {
    tasks: {
      type: Object,
      required: true
    },
    currentCronJob: {
      type: Object,
      default: null
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
</style>

