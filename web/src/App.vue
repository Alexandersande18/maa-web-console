<template>
  <router-view/>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
/* 定义全局 CSS 变量 */
:root {
  --primary-color: #007BFF;
  --primary-hover: #0056b3;
  --danger-color: #dc3545;
  --danger-hover: #c82333;
  --success-color: #35dca5;
  --success-hover: #2ebd94;
  --white: #fff;
  --border-color: #ddd;
  --shadow-color: rgba(0, 0, 0, 0.1);
  --background-color: #f2f2f2;
}

:root:root {
  --van-button-small-padding: 10px;
  --van-button-small-font-size: 14px;
  --van-button-small-height: 35px
}

body {
  margin: 0;
  font-family: Arial, sans-serif;
  background-color: var(--background-color);
  color: #333;
}
</style>