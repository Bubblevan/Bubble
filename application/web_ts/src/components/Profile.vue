<template>
    <Layout>
      <div class="profile-container">
        <h1>个人信息</h1>
        <el-form :model="user" label-width="120px">
          <el-form-item label="用户名">
            <el-input v-model="user.username" disabled></el-input>
          </el-form-item>
          <el-form-item label="电子邮件">
            <el-input v-model="user.email" disabled></el-input>
          </el-form-item>
          <el-form-item label="全名">
            <el-input v-model="user.fullName" disabled></el-input>
          </el-form-item>
          <el-form-item label="角色">
            <el-input v-model="user.roleText" disabled></el-input>
          </el-form-item>
        </el-form>
      </div>
    </Layout>
  </template>
  
  <script>
  import Layout from '../layout/Layout.vue';
  import axios from 'axios';
  
  export default {
    components: {
      Layout
    },
    data() {
      return {
        user: {
          username: '',
          email: '',
          fullName: '',
          roleText: ''
        }
      };
    },
    created() {
      this.fetchProfile();
    },
    methods: {
      async fetchProfile() {
        try {
          const token = sessionStorage.getItem('token');
          console.log(token);
          const response = await axios.get('http://localhost:8080/profile', {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          });
          this.user = response.data.data;
          this.user.roleText = this.user.role ? '亚运主办方' : '游客';
        } catch (error) {
          this.$message.error('获取个人信息失败');
        }
      }
    }
  };
  </script>
  
  <style scoped>
  .profile-container {
    max-width: 600px;
    margin: 50px auto;
    padding: 20px;
    background: #fff;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    border-radius: 8px;
  }
  
  h1 {
    text-align: center;
    margin-bottom: 20px;
  }
  
  .el-form-item {
    margin-bottom: 20px;
  }
  </style>
  