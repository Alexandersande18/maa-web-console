<template>
  <van-nav-bar title="🉑🍬🉑🧂"
  :class="{ 'custom-navbar': showLogout }"
  right-text="登出"
  fixed="true"
  @click-right="$emit('logout')"
  >
  </van-nav-bar>
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

/* 方法 2：使用自定义 class */
::v-deep .custom-navbar {
  --van-nav-bar-text-color: #ff4d4f;
  --van-nav-bar-text-font-size: 14px;
}

/* 添加 hover 效果 */
::v-deep .van-nav-bar__right-text:hover {
  opacity: 0.8;
  cursor: pointer;
}



</style>