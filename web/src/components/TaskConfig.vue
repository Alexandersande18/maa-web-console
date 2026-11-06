<template>
  <div>
    <h3>任务配置</h3>

    <CronJobList
      :cronJobs="cronJobs"
      :currentCronJob="currentCronJob"
      :formatTimestamp="formatTimestamp"
      @select-cron-job="selectCronJob"
      @edit-cron-job="openEditCronJobModal"
      @delete-cron-job="confirmAndDelete"
    />

    <TaskListPanel
      :tasks="tasks"
      :currentCronJob="currentCronJob"
      @toggle-task="toggleTaskSelection"
      @open-config="openConfig"
      @start-tasks="startTasks"
      @stop-tasks="stopTasks"
      @add-cron-job="openAddCronJobFromTasks"
      @save-cron-job="updateCronJob"
    />

    <ConfigModal
      v-if="showModal"
      :taskKey="currentTaskKey"
      :taskConfig="tasks[currentTaskKey].config"
      :fieldConfig="fieldConfig[currentTaskKey] || {}"
      @save="saveTaskConfig"
      @close="closeModal"
    />

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
import CronJobList from "./task-config/CronJobList.vue";
import TaskListPanel from "./task-config/TaskListPanel.vue";
import { createDefaultTasks, FIELD_CONFIG } from "./task-config/taskDefinitions";
import {
  buildCronJobPayloadFromTasks,
  buildEnabledTasksPayload,
  formatTimestamp,
  syncTasksFromCronJob
} from "./task-config/taskConfigUtils";

export default {
  name: "TaskConfig",
  components: { ConfigModal, CronJobModal, CronJobList, TaskListPanel },
  data() {
    return {
      tasks: createDefaultTasks(),
      fieldConfig: FIELD_CONFIG,
      showModal: false,
      currentTaskKey: "",
      cronJobs: [],
      currentCronJob: null,
      showCronJobModal: false,
      cronJobModalData: null,
      isEditingCronJob: false
    };
  },
  mounted() {
    this.fetchCronJobs();
  },
  methods: {
    formatTimestamp,
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
    async startTasks() {
      const jwtToken = localStorage.getItem("jwtToken");
      try {
        const body = buildEnabledTasksPayload(this.tasks, true);
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
          alert(result.data?.message || "任务已加入队列，将按顺序执行");
        } else {
          alert(`操作失败：${result.message || "未知业务错误"}`);
        }
      } catch (e) {
        alert("请求异常：" + e.message);
      }
    },
    async stopTasks() {
      const body = {};
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
      syncTasksFromCronJob(this.tasks, cronJob);
    },
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
    closeCronJobModal() {
      this.showCronJobModal = false;
      this.cronJobModalData = null;
      this.isEditingCronJob = false;
    },
    handleAddCronJob(data) {
      const payload = {
        name: data.name,
        hour: data.hour,
        minute: data.minute,
        isTemplate: data.isTemplate,
        ...buildCronJobPayloadFromTasks(this.tasks)
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
        .then((response) => response.json())
        .then((result) => {
          if (result.data && result.data.id !== undefined) {
            alert("添加定时任务成功");
            this.fetchCronJobs();
          } else {
            alert("添加定时任务失败");
          }
        })
        .catch((e) => {
          alert("添加定时任务异常：" + e.message);
        });
      this.closeCronJobModal();
    },
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
        .then((response) => response.json())
        .then((result) => {
          if (result.code === 0) {
            alert("定时任务修改成功");
            this.fetchCronJobs();
            this.closeCronJobModal();
          } else {
            alert("定时任务修改失败：" + (result.message || "未知错误"));
          }
        })
        .catch((e) => {
          alert("更新定时任务异常：" + e.message);
        });
    },
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
        ...buildCronJobPayloadFromTasks(this.tasks)
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
    confirmAndDelete(job) {
      if (confirm("确定要删除该定时任务吗？")) {
        const jwtToken = localStorage.getItem("jwtToken");
        fetch(`/maasst/cronjob?id=${job.id}`, {
          method: "DELETE",
          headers: { Authorization: jwtToken }
        })
          .then((response) => response.json())
          .then((result) => {
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
          .catch((e) => {
            alert("删除定时任务异常：" + e.message);
          });
      }
    }
  }
};
</script>

