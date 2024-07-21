<template>
    <Layout>
      <div class="users-edit-container">
        <h1>编辑用户</h1>
        <el-form :model="user" label-width="120px">
          <el-form-item label="用户名">
            <el-input v-model="user.username"></el-input>
          </el-form-item>
          <el-form-item label="电子邮件">
            <el-input v-model="user.email"></el-input>
          </el-form-item>
          <el-form-item label="全名">
            <el-input v-model="user.fullName"></el-input>
          </el-form-item>
          <el-form-item label="角色">
            <el-select v-model="user.role">
              <el-option label="游客" :value="false" />
              <el-option label="主办方" :value="true" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="updateUser">保存</el-button>
            <el-button type="danger" @click="deleteUser">删除</el-button>
          </el-form-item>
        </el-form>
      </div>
    </Layout>
  </template>
  
  <script>
  import axios from 'axios';
  import Layout from '../layout/Layout.vue';
  export default {
    components: {
      Layout
    },
    props: ['id'],
    data() {
      return {
        user: {
          username: '',
          email: '',
          fullName: '',
          role: false,
        },
      };
    },
    methods: {
      fetchUser() {
        // 替换为实际的 API 调用
        axios.get(`/api/users/${this.id}`).then((response) => {
          this.user = response.data;
        });
      },
      updateUser() {
        // 替换为实际的 API 调用
        axios.put(`/api/users/${this.id}`, this.user).then(() => {
          this.$router.push({ name: 'Users' });
        });
      },
      deleteUser() {
        // 替换为实际的 API 调用
        axios.delete(`/api/users/${this.id}`).then(() => {
          this.$router.push({ name: 'Users' });
        });
      },
    },
    created() {
      this.fetchUser();
    },
  };
  </script>
  
  <style>
  .users-edit-container {
    padding: 20px;
  }
  </style>
  