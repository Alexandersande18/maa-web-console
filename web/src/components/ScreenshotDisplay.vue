<template>
  <div>
    <h3>截图展示</h3>
    
    <div class="box" id="screenshot-box">
      <van-row justify="center">
        <!-- 使用 Vant 的 van-button，并设置 size 属性为 small -->
        <van-button type="primary" size="small" @click="takeScreenshot">
          截 图
        </van-button>
      </van-row>
      <div id="screenshot-display">
        <img
          id="screenshot-img"
          class="screenshot-img"
          :src="screenshotSrc"
          alt="截图内容"
          v-if="screenshotSrc"
        />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ScreenshotDisplay',
  data() {
    return {
      screenshotSrc: ''
    }
  },
  methods: {
    async takeScreenshot() {
      const jwtToken = localStorage.getItem('jwtToken')
      try {
        const response = await fetch('/maasst/screenshot', {
          method: 'POST',
          headers: { 
            'Content-Type': 'application/json',
            'Authorization': jwtToken
          },
          body: JSON.stringify({})
        })
        const result = await response.json()
        if (result && result.data && result.data.imageData) {
          this.screenshotSrc = "data:image/png;base64," + result.data.imageData
        } else {
          alert(result?.message || "截图失败")
        }
      } catch(e) {
        alert("截图请求出错：" + e)
      }
    }
  }
}
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

.center-button {
  text-align: center;
  margin-bottom: 15px;
}

#screenshot-display {
  text-align: center;
  margin-top: 20px;
}

.screenshot-img {
  max-width: 100%;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  box-shadow: 0 2px 4px var(--shadow-color);
}
</style>
