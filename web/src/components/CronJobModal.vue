<template>
    <div class="modal" @click="handleBackgroundClick">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ modalTitle }}</h3>
          <span class="modal-close" @click="closeModal">&times;</span>
        </div>
        <div id="modal-body">
          <div class="form-group">
            <label>任务名称</label>
            <input type="text" v-model="localCronJob.name" />
          </div>
          <div class="form-group time-group">
            <label>选择执行时间</label>
            <div class="time-picker">
              <select v-model.number="localCronJob.hour">
                <option v-for="hour in hours" :key="hour" :value="hour">
                  {{ hour < 10 ? "0" + hour : hour }}
                </option>
              </select>
              <span class="time-separator">:</span>
              <select v-model.number="localCronJob.minute">
                <option v-for="minute in minutes" :key="minute" :value="minute">
                  {{ minute < 10 ? "0" + minute : minute }}
                </option>
              </select>
            </div>
          </div>
          <div class="form-group checkbox-label">
            <label>
              <input type="checkbox" v-model="localIsTemplate" />
              模板任务
            </label>
          </div>
        </div>
        <div class="modal-footer">
          <button class="save-button" @click="submitCronJob">
            {{ actionButtonLabel }}
          </button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: "CronJobModal",
    props: {
      // 任务数据：添加时为空对象或者默认值，编辑时包含完整数据（如 id、name、hour、minute、isTemplate）
      cronJobData: {
        type: Object,
        default: () => ({ name: "", hour: 0, minute: 0, isTemplate: 0 })
      },
      // 是否处于编辑模式
      isEditing: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        // 本地数据副本，可供双向绑定使用
        localCronJob: { ...this.cronJobData },
        localIsTemplate: this.cronJobData.isTemplate == 1
      };
    },
    computed: {
      // 根据当前模式动态显示标题
      modalTitle() {
        return this.isEditing ? "编辑定时任务" : "添加定时任务";
      },
      // 根据当前模式动态设置提交按钮文字
      actionButtonLabel() {
        return this.isEditing ? "保存" : "添加";
      },
      // 生成 0 到 23 的数组
      hours() {
        return Array.from({ length: 24 }, (_, i) => i);
      },
      // 生成 0 到 59 的数组
      minutes() {
        return Array.from({ length: 60 }, (_, i) => i);
      }
    },
    watch: {
      // 当本地模板值更改时，同步更新 localCronJob 对象中的 isTemplate 字段
      localIsTemplate(val) {
        this.localCronJob.isTemplate = val ? 1 : 0;
      },
      // 当外部传入的数据发生变化时，同步更新本地数据（如关闭后再次打开时更新数据）
      cronJobData: {
        handler(newVal) {
          this.localCronJob = { ...newVal };
          this.localIsTemplate = newVal.isTemplate == 1;
        },
        immediate: true,
        deep: true
      }
    },
    methods: {
      submitCronJob() {
        // 根据模式发出不同事件
        if (this.isEditing) {
          this.$emit("update", this.localCronJob);
        } else {
          this.$emit("add", this.localCronJob);
        }
        this.closeModal();
      },
      handleBackgroundClick() {
        this.closeModal();
      },
      closeModal() {
        this.$emit("close");
      }
    }
  };
  </script>
  
  <style scoped>
  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }
  .modal {
    position: fixed;
    z-index: 1000;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    overflow: auto;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .modal-content {
    background-color: #fff;
    margin: 15px;
    padding: 25px 30px;
    border-radius: 10px;
    width: calc(100% - 40px);
    max-width: 400px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    animation: fadeIn 0.3s ease-out;
    box-sizing: border-box;
    position: relative;
  }
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #eee;
    margin-bottom: 20px;
    padding-bottom: 10px;
  }
  .modal-header h3 {
    margin: 0;
    font-size: 22px;
    color: #333;
  }
  .modal-close {
    font-size: 26px;
    cursor: pointer;
    color: #999;
    transition: color 0.3s;
  }
  .modal-close:hover {
    color: #666;
  }
  #modal-body {
    max-height: 300px;
    overflow-y: auto;
    margin-bottom: 20px;
  }
  .form-group {
    margin-bottom: 15px;
  }
  .form-group label {
    display: block;
    margin-bottom: 5px;
    font-weight: 600;
    color: #444;
  }
  /* 优化后的下拉选择样式 */
  .time-picker {
    display: flex;
    align-items: center;
  }
  .time-picker select {
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 16px;
    background-color: #fff;
    background-image: url("data:image/svg+xml,%3Csvg width='12' height='12' viewBox='0 0 20 20' xmlns='http://www.w3.org/2000/svg'%3E%3Cpolygon points='0,0 20,0 10,10' fill='%23999'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 10px center;
    background-size: 12px;
    margin-right: 8px;
  }
  /* 去除最后一个下拉框右侧外边距 */
  .time-picker select:last-child {
    margin-right: 0;
  }
  .time-separator {
    font-size: 16px;
    margin: 0 4px;
    color: #555;
  }
  .checkbox-label label {
    font-weight: 600;
    display: flex;
    align-items: center;
  }
  .checkbox-label input[type="checkbox"] {
    margin-right: 8px;
  }
  input[type="text"] {
    width: 100%;
    box-sizing: border-box;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 16px;
  }
  .modal-footer {
    text-align: right;
  }
  .save-button {
    background-color: #007bff;
    color: #fff;
    border: none;
    padding: 10px 20px;
    font-size: 14px;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s;
  }
  .save-button:hover {
    background-color: #0056b3;
  }
  </style>