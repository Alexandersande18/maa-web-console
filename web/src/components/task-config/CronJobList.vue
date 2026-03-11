<template>
  <div class="box" id="cronjob-section">
    <div class="cronjob-scroll">
      <template v-if="cronJobs.length">
        <div
          class="cronjob-item"
          v-for="job in cronJobs"
          :key="job.id"
          :class="{ selected: currentCronJob && currentCronJob.id === job.id }"
          @click="$emit('select-cron-job', job)"
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
          <div class="item-actions">
            <van-button size="small" type="primary" @click.stop="$emit('edit-cron-job', job)">配 置</van-button>
            <van-button size="small" type="danger" @click.stop="$emit('delete-cron-job', job)">删 除</van-button>
          </div>
        </div>
      </template>
      <div v-else class="empty-cronjob">没有定时任务</div>
    </div>
  </div>
</template>

<script>
export default {
  name: "CronJobList",
  props: {
    cronJobs: {
      type: Array,
      required: true
    },
    currentCronJob: {
      type: Object,
      default: null
    },
    formatTimestamp: {
      type: Function,
      required: true
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
  max-height: 300px;
  overflow-y: auto;
  width: 100%;
  padding-right: 5px;
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

.config-summary {
  clear: both;
  font-size: 12px;
  color: #666;
  word-break: break-all;
}
</style>

