<template>
  <div>
    <Header :showLogout="true" @logout="logout" />
    <div id="main-page">
      <TaskConfig />
      <LogDisplay />
      <ScreenshotDisplay />
    </div>
  </div>
</template>

<script>
import Header from '../components/Header.vue'
import TaskConfig from '../components/TaskConfig.vue'
import LogDisplay from '../components/LogDisplay.vue'
import ScreenshotDisplay from '../components/ScreenshotDisplay.vue'

export default {
  name: 'MainView',
  components: {
    Header,
    TaskConfig,
    LogDisplay,
    ScreenshotDisplay
  },
  methods: {
    logout() {
      localStorage.removeItem('jwtToken')
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
/* 增加上边距，避免内容被固定的 Header 遮挡 */
#main-page {
  max-width: 700px;
  margin: 50px auto 20px;
  padding: 10px;
}
</style>