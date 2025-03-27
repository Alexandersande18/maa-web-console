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
  --van-nav-bar-text-color: #ff4d4f;
  --van-button-small-padding: 10px;
  --van-button-small-font-size: 14px;
  --van-button-small-height: 35px;
  --van-dialog-small-screen-width: 90%;
  --van-cell-vertical-padding: 10px;
  --van-cell-group-inset-padding: 0px;
  --van-cell-group-title-color: rgba(0, 0, 0, 0.575);
  --van-cell-group-title-font-size: 15px;
  /* --van-cell-horizontal-padding: */
  --van-dialog-width: 360px;
  /* --van-dialog-background:var(--background-color); */

}

h3 {
  margin-bottom: 5px; /* 根据需求调整具体数值 */
}

body {
  margin: 0;
  font-family: Arial, sans-serif;
  background-color: var(--background-color);
  color: #333;
}
</style>