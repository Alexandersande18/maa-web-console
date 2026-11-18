<!-- src/components/ConfigModal.vue -->
<template>
  <!-- 点击最外层背景区域也触发退出并保存 -->
  <div class="modal" @click="handleModalBackgroundClick">
    <!-- 在 modal-content 中阻止事件冒泡，防止点击内容时退出 -->
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>{{ modalTitle }}</h3>
        <!-- 关闭按钮点击时自动保存并关闭 -->
        <span class="modal-close" @click="closeAndSave">&times;</span>
      </div>

      <div id="modal-body">
        <!-- 遍历过滤后的配置项 -->
        <div class="form-group" v-for="(value, key) in filteredTaskConfig" :key="key">
          <!-- 如果配置字段为布尔类型则复选框放在文字前面 -->
          <template v-if="isBoolean(value)">
            <label class="checkbox-label">
              <input type="checkbox" v-model="localConfig[key]" />
              <span>{{ (fieldConfig[key] && fieldConfig[key].label) || key }}</span>
            </label>
            <small v-if="fieldConfig[key] && fieldConfig[key].hint" class="hint-text">
              {{ fieldConfig[key].hint }}
            </small>
          </template>
          <template v-else>
            <label :id="'config-label-' + key">
              {{ (fieldConfig[key] && fieldConfig[key].label) || key }}
            </label>
            <small v-if="fieldConfig[key] && fieldConfig[key].hint" class="hint-text">
              {{ fieldConfig[key].hint }}
            </small>
            <div v-if="getOptions(key)" class="input-wrapper">
              <select
                v-model="localConfig[key]"
                :aria-labelledby="'config-label-' + key"
                :aria-describedby="key === 'mode' && reclamationStageGuide ? 'reclamation-stage-guide' : undefined"
              >
                <option
                  v-for="(optText, optValue) in getOptions(key)"
                  :key="optValue"
                  :value="typeof value === 'number' ? Number(optValue) : optValue"
                >
                  {{ optText }}
                </option>
              </select>
            </div>
            <div v-else-if="getNumberOptions(key)" class="input-wrapper">
              <van-stepper
                v-model="localConfig[key]"
                v-bind="getNumberOptions(key)"
                :name="key"
                :aria-labelledby="'config-label-' + key"
                class="number-input"
                button-size="40"
                @blur="handleNumberInput(key)"
              />
            </div>
            <div
              v-else-if="['facility', 'buy_first', 'blacklist', 'select', 'confirm', 'first_tags', 'start_foldartal_list', 'expected_collapsal_paradigms'].includes(key)"
              class="input-wrapper"
            >
              <textarea v-model="localConfig[key]" rows="2"></textarea>
            </div>
            <div v-else-if="['drops', 'recruitment_time', 'collectible_mode_start_list', 'collectible_mode_core_char_list'].includes(key)" class="input-wrapper">
              <textarea v-model="localConfig[key]" rows="4"></textarea>
            </div>
            <div v-else class="input-wrapper">
              <input type="text" v-model="localConfig[key]" />
            </div>
            <aside v-if="key === 'mode' && reclamationStageGuide" id="reclamation-stage-guide" class="stage-guide" aria-live="polite">
              <div v-for="section in reclamationStageGuide" :key="section.label">
                <strong>{{ section.label }}</strong>
                <ol v-if="section.steps">
                  <li v-for="step in section.steps" :key="step">{{ step }}</li>
                </ol>
                <p v-else>{{ section.text }}</p>
              </div>
            </aside>
          </template>
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="save-button" @click="closeAndSave">保存</button>
      </div>
    </div>
  </div>
</template>

<script>
import { TASK_TEMPLATE, getFieldOptions, getNumberFieldOptions, getReclamationStageGuide, isTaskFieldVisible } from './task-config/taskDefinitions'
import { normalizeNumberField, processTaskConfig } from './task-config/taskConfigUtils'

export default {
  name: 'ConfigModal',
  props: {
    taskKey: String,
    taskConfig: Object,
    fieldConfig: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      // 使用浅拷贝避免直接修改 prop 数据
      localConfig: { ...this.taskConfig }
    }
  },
  computed: {
    // 过滤掉不需要在弹窗中显示的字段，比如 'enable'
    filteredTaskConfig() {
      const filteredKeys = new Set(['enable', 'stone', 'report_to_penguin', 'penguin_id', 'server', 'client_type', 'DrGrandet', 'first_floor_foldartal', 'start_foldartal_list', 'use_foldartal', 'check_collapsal_paradigms', 'double_check_collapsal_paradigms', 'expected_collapsal_paradigms', 'collectible_mode_shopping', 'collectible_mode_squad', 'collectible_mode_start_list']);
      let result = {};
      for (const key of Object.keys(this.localConfig)) {
        if ((!filteredKeys.has(key) || (key === 'collectible_mode_squad' && this.localConfig.theme === 'BlackFlow')) && isTaskFieldVisible(this.taskKey, key, this.localConfig)) {
          result[key] = this.localConfig[key];
        }
      }
      return result;
    },
    modalTitle() {
      return "配置 " + (TASK_TEMPLATE[this.taskKey]?.label || this.taskKey)
    },
    reclamationStageGuide() {
      return this.taskKey === 'taskReclamation' ? getReclamationStageGuide(this.localConfig) : null
    }
  },
  watch: {
    'localConfig.theme'() {
      const options = this.getOptions('mode')
      if (options && !Object.hasOwn(options, this.localConfig.mode)) this.localConfig.mode = Number(Object.keys(options)[0])
    }
  },
  methods: {
    getOptions(key) {
      return getFieldOptions(this.taskKey, key, this.localConfig, this.fieldConfig)
    },
    isBoolean(val) {
      return typeof val === 'boolean'
    },
    getNumberOptions(key) {
      return getNumberFieldOptions(this.taskKey, key)
    },
    handleNumberInput(key) {
      this.localConfig[key] = normalizeNumberField(this.taskKey, key, this.localConfig[key])
    },
    closeAndSave() {
      for (const key of Object.keys(this.localConfig)) {
        if (this.getNumberOptions(key)) this.handleNumberInput(key)
      }
      try {
        processTaskConfig(this.taskKey, this.localConfig)
      } catch (e) {
        alert('配置错误：' + e.message)
        return
      }
      // 自动保存当前配置并通知关闭
      this.$emit('save', this.localConfig)
      this.$emit('close')
    },
    handleModalBackgroundClick() {
      // 点击弹窗外部背景区域，同样保存并关闭
      this.closeAndSave()
    }
  }
}
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

/* 将最外层 modal 设置为 flex 布局，方便使内容居中 */
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
  margin: 15px; /* 固定的上下左右间距，保证左右各 20px */
  padding: 25px 30px;
  border-radius: 10px;
  /* 使用 calc(100% - 40px) 确保宽度为 (viewport 宽度 - 2×20px)，同时最大宽度为 600px */
  width: calc(100% - 40px);
  max-width: 600px;
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
  max-height: min(500px, calc(100dvh - 160px));
  overflow-y: auto;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  font-weight: 600;
  color: #444;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  margin-right: 8px;
  transform: scale(1.0);
}

.form-group label {
  font-weight: 600;
  color: #444;
  display: block;
  margin-bottom: 5px;
}

.hint-text {
  display: block;
  font-size: 12px;
  color: #888;
  margin-bottom: 8px;
}

.stage-guide {
  display: grid;
  gap: 10px;
  margin-top: 12px;
  padding: 12px;
  border: 1px solid #e4e9ef;
  border-radius: 6px;
  background: #f7f9fc;
  color: #374151;
  font-size: 14px;
  line-height: 1.55;
}

.stage-guide strong {
  color: #252e3a;
}

.stage-guide p,
.stage-guide ol {
  margin: 2px 0 0;
}

.stage-guide ol {
  padding-inline-start: 1.5em;
  list-style: decimal;
}

.stage-guide li + li {
  margin-top: 4px;
}

.input-wrapper input[type="text"],
.input-wrapper textarea,
.input-wrapper select {
  width: 100%;
  box-sizing: border-box;
  padding: 10px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  line-height: 1.4;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.input-wrapper input[type="text"]:focus,
.input-wrapper textarea:focus,
.input-wrapper select:focus {
  border-color: var(--primary-color);
  outline: none;
  box-shadow: 0 0 4px rgba(0, 123, 255, 0.3);
}

.modal-footer {
  text-align: right;
}

.number-input {
  display: flex;
  align-items: center;
  gap: 8px;
}

.number-input :deep(.van-stepper__input) {
  flex: 1;
  min-width: 0;
  width: 0;
  height: 40px;
  margin: 0;
  padding: 8px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--white);
  font-size: 16px;
}

.number-input :deep(button) {
  flex-shrink: 0;
}

.save-button {
  background-color: var(--primary-color);
  color: #fff;
  border: none;
  padding: 10px 20px;
  font-size: 14px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s, box-shadow 0.3s;
}

.save-button:hover {
  background-color: var(--primary-hover, #0056b3);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
}
</style>
