<template>
  <div id="login-page">
    <h2>登录</h2>
    <form @submit.prevent="login">
      <input
        type="text"
        v-model="token"
        placeholder="请输入 token"
        name="token"
        autocomplete="on"
      />
      <button type="submit">登录</button>
    </form>
    <div id="login-message">{{ loginMessage }}</div>
  </div>
</template>

<script>
export default {
  name: 'LoginView',
  data() {
    return {
      token: '',
      loginMessage: ''
    }
  },
  methods: {
    async login() {
      const trimmedToken = this.token.trim()
      if (!trimmedToken) {
        this.loginMessage = '请输入 token'
        return
      }
      try {
        const response = await fetch('/maasst/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ Token: trimmedToken })
        })
        const result = await response.json()
        if (result && result.data && result.data.token) {
          localStorage.setItem('jwtToken', result.data.token)
          this.$router.push('/main')
        } else {
          this.loginMessage = '登录失败，请检查 token'
        }
      } catch (e) {
        this.loginMessage = `请求出错：${e}`
      }
    }
  }
}
</script>

<style scoped>
#login-page {
  max-width: 500px;
  margin: 50px auto;
  padding: 20px;
  background-color: var(--white);
  border-radius: 6px;
  box-shadow: 0 2px 4px var(--shadow-color);
}
#login-message {
  color: red;
  text-align: center;
}
input {
  width: 100%;
  padding: 8px;
  margin-top: 5px;
  margin-bottom: 10px;
  box-sizing: border-box;
  border: 1px solid var(--border-color);
  border-radius: 4px;
}
button {
  padding: 10px 15px;
  margin: 5px;
  border: none;
  background-color: var(--primary-color);
  color: var(--white);
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}
button:hover {
  background-color: var(--primary-hover);
}
</style>