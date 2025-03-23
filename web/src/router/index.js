// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import MainView from '../views/MainView.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: {
      title: "Login"
    }
  },
  {
    path: '/main',
    name: 'Main',
    component: MainView,
    meta: { 
      requiresAuth: true,
      title: "🍬🧂🍬🧂🍬🧂" 
    } // 设置元数据，标记此路由需要认证
  },
  { path: '/', redirect: '/login' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('jwtToken')
  if (to.meta.title) {
    document.title = to.meta.title;
  }
  // 如果路由需要认证，但没有 token，则跳转到登录页
  if (to.matched.some(record => record.meta.requiresAuth) && !token) {
    next({ path: '/login' })
  }
  // 如果访问登录页时已经拥有 token，则自动跳转到主页面
  else if (to.path === '/login' && token) {
    next({ path: '/main' })
  }
  else {
    next() // 正常跳转
  }
})

export default router