import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import Vant from 'vant'
import 'vant/lib/index.css'


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
// window.fetch = async (...args) => {
//   try {
//     const response = await originalFetch(...args)
//     // 当返回 401 状态码时，说明 JWT 无效或过期
//     if (response.status === 401) {
//       // 移除本地保存的 JWT
//       localStorage.removeItem('jwtToken')
//       // 跳转到登录页
//       router.push('/login')
//       return Promise.reject(new Error('Unauthorized'))
//     }
//     return response
//   } catch (err) {
//     return Promise.reject(err)
//   }
// }

window.fetch = async (...args) => {
  // 1. 请求前统一添加 Token（避免手动设置 headers）
  const [url, config = {}] = args;
  const jwtToken = localStorage.getItem('jwtToken');
  
  const newConfig = {
    ...config,
    headers: {
      ...config.headers,
      ...(jwtToken ? { Authorization: jwtToken } : {}),
    },
  };

  try {
    const response = await originalFetch(url, newConfig);

    // 2. 统一处理 HTTP 错误状态码
    if (!response.ok) {
      const error = new Error(`HTTP Error: ${response.status}`);
      error.response = response;
      throw error; // 主动抛出错误，确保进入 catch
    }

    return response;
  } catch (error) {
    // 3. 集中处理 401
    if (error.response?.status === 401) {
      localStorage.removeItem('jwtToken');
      router.push('/login'); // 确保 router 已初始化
      alert('登录已过期，请重新登录'); // 或使用 Vant 轻提示
    } else {
      alert(`请求失败: ${error.message}`);
    }

    throw error; // 继续抛出错误，但组件可选择忽略
  }
};
  
createApp(App).
use(router).
use(Vant).
// use(vuetify).
mount('#app')
