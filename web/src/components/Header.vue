<template>
  <div class="header">
    <h1>🉑🍬🉑🧂</h1>
    <button v-if="showLogout" id="logout-button" @click="$emit('logout')">登出</button>
  </div>
</template>

<script>
export default {
  name: 'AppHeader',
  props: {
    showLogout: {
      type: Boolean,
      default: false,
    },
  },
}
</script>

<style scoped>
.header {
  position: fixed;          /* 固定在页面顶部 */
  top: 0;
  left: 0;
  width: 100%;              /* 撑满整个宽度 */
  background-color: var(--white);
  z-index: 1000;
  padding: 10px 20px;       /* 增大左右内边距，确保内容与窗口边缘有间距 */
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: center;
}

.header h1 {
  margin: 0;
  font-size: 1.5em;
}

/* 调整登出按钮，确保其距离右边有足够距离，以免被滚动条覆盖 */
#logout-button {
  position: absolute;
  right: 50px;              /* 往内移，避免与滚动条重叠 */
  background-color: var(--danger-color);
  padding: 10px 15px;
  color: var(--white);
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

#logout-button:hover {
  background-color: var(--danger-hover);
}
</style>