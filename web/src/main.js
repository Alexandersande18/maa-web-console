import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// import Vant from 'vant'
// import 'vant/lib/index.css'


// // Vuetify
// import 'vuetify/styles'
// import { createVuetify } from 'vuetify'
// import * as components from 'vuetify/components'
// import * as directives from 'vuetify/directives'
// const vuetify = createVuetify({
//     components,
//     directives,
//   })

// 保存原始的 fetch 方法
const originalFetch = window.fetch

// 重写全局 fetch
window.fetch = async (...args) => {
  try {
    const response = await originalFetch(...args)
    // 当返回 401 状态码时，说明 JWT 无效或过期
    if (response.status === 401) {
      // 移除本地保存的 JWT
      localStorage.removeItem('jwtToken')
      // 跳转到登录页
      router.push('/login')
      return Promise.reject(new Error('Unauthorized'))
    }
    return response
  } catch (err) {
    return Promise.reject(err)
  }
}
  

createApp(App).
use(router).
// use(Vant).
// use(vuetify).
mount('#app')
