<template>
  <Layout>
    <h1>欢迎来到主页面</h1>
    <!-- <div class="token-display">
        <h2>Token:</h2>
        <p>{{ token }}</p>
      </div> -->
  </Layout>
</template>

<script>
import Layout from '../layout/Layout.vue';

export default {
  components: {
    Layout
  }
};
// const token = sessionStorage.getItem('token');
// console.log(token, '真的执行了打印');
</script>

<style scoped>

</style>
