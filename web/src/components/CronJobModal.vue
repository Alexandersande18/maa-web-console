<template>
  <van-dialog
    :show="visible"
    :title="modalTitle"
    show-cancel-button
    :confirm-button-text="actionButtonLabel"
    cancel-button-text="取消"
    @confirm="submitCronJob"
    @cancel="closeModal"
  >
    <van-cell-group inset>
      <!-- 任务名称输入 -->
      <van-cell size="normal">
        <van-field
          v-model="localCronJob.name"
          label="任务名称"
          placeholder="请输入任务名称"
        />

        <van-field name="switch" label="模板任务">
          <template #input>
            <van-checkbox v-model="localIsTemplate" shape="square" />
          </template>
        </van-field>

        <van-time-picker
          v-model="currentTime"
          title="选择时间"
          :columns-type="['hour', 'minute']"
          :show-toolbar="false"
        />
      </van-cell>
    </van-cell-group>
  </van-dialog>
</template>

<script>
export default {
  name: "CronJobModal",
  props: {
    // 任务数据：添加时为空对象或者默认值，编辑时包含完整数据（例如 id、name、hour、minute、isTemplate）
    cronJobData: {
      type: Object,
      default: () => ({ name: "", hour: 0, minute: 0, isTemplate: 0 }),
    },
    // 是否处于编辑模式
    isEditing: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      // 控制 Dialog 的显示
      visible: true,
      // 本地数据副本，供双向绑定使用
      localCronJob: { ...this.cronJobData },
      localIsTemplate: this.cronJobData.isTemplate == 1,
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
    // 绑定 vant 的 time-picker，格式为 ["HH", "MM"]
    currentTime: {
      get() {
        const hh = this.localCronJob.hour;
        const mm = this.localCronJob.minute;
        return [String(hh).padStart(2, "0"), String(mm).padStart(2, "0")];
      },
      set(newTime) {
        // newTime 为 ["08", "30"] 格式
        this.localCronJob.hour = Number(newTime[0]);
        this.localCronJob.minute = Number(newTime[1]);
      },
    },
  },
  watch: {
    // 当本地模板值更改时，同步更新 localCronJob 对象中的 isTemplate 字段
    localIsTemplate(val) {
      this.localCronJob.isTemplate = val ? 1 : 0;
    },
    // 当外部传入数据发生变化时，同步更新本地数据（例如关闭后再次打开时更新数据）
    cronJobData: {
      handler(newVal) {
        this.localCronJob = { ...newVal };
        this.localIsTemplate = newVal.isTemplate == 1;
      },
      immediate: true,
      deep: true,
    },
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
    closeModal() {
      this.visible = false;
      this.$emit("close");
    },
  },
};
</script>