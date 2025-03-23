<template>
  <div>
    <h2>截图展示</h2>
    
    <div class="box" id="screenshot-box">
      <div class="center-button">
        <button class="btn btn-primary" @click="takeScreenshot">截图</button>
      </div>
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
          alert("截图失败")
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

/* 复用全局 btn 样式 */
.btn {
  padding: 8px 13px;
  margin: 0px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease;
  font-size: 14px;
  color: var(--white);
}

.btn-primary {
  background-color: var(--primary-color);
}

.btn-primary:hover {
  background-color: var(--primary-hover);
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