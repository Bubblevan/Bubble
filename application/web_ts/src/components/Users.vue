<template>
    <Layout>
      <div class="users-container">
        <h1>用户管理</h1>
        <el-input
          v-model="searchQuery"
          placeholder="输入用户名进行搜索"
          @keyup.enter="fetchUsers"
        />
        <el-button type="primary" @click="fetchUsers">搜索</el-button>
        <el-button type="success" @click="showAddUserDialog">添加用户</el-button>
        <el-table :data="users" style="width: 100%">
          <el-table-column prop="username" label="用户名" width="180" />
          <el-table-column prop="email" label="电子邮件" />
          <el-table-column prop="fullName" label="全名" />
          <el-table-column prop="roleText" label="角色" />
          <el-table-column label="操作">
            <template v-slot="scope">
              <el-button
                type="primary"
                size="mini"
                @click="editUser(scope.row)"
              >编辑</el-button>
              <el-button
                type="danger"
                size="mini"
                @click="deleteUser(scope.row)"
              >删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <el-dialog :visible.sync="addUserDialogVisible" title="添加用户">
        <el-form :model="newUser">
          <el-form-item label="用户名">
            <el-input v-model="newUser.username" />
          </el-form-item>
          <el-form-item label="电子邮件">
            <el-input v-model="newUser.email" />
          </el-form-item>
          <el-form-item label="全名">
            <el-input v-model="newUser.fullName" />
          </el-form-item>
          <el-form-item label="角色">
            <el-select v-model="newUser.role">
              <el-option label="游客" :value="false" />
              <el-option label="主办方" :value="true" />
            </el-select>
          </el-form-item>
        </el-form>
        <span slot="footer" class="dialog-footer">
          <el-button @click="addUserDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addUser">添加</el-button>
        </span>
      </el-dialog>
    </Layout>
  </template>
  
  <script>
  import axios from 'axios';
  import Layout from '../layout/Layout.vue';
  export default {
    components: {
      Layout
    },
    data() {
      return {
        users: [],
        searchQuery: '',
        addUserDialogVisible: false,
        newUser: {
          username: '',
          email: '',
          fullName: '',
          role: false,
        },
      };
    },
    methods: {
      fetchUsers() {
        // 替换为实际的 API 调用
        axios.get(`/api/users?search=${this.searchQuery}`).then((response) => {
          this.users = response.data;
        });
      },
      showAddUserDialog() {
        this.addUserDialogVisible = true;
      },
      addUser() {
        // 替换为实际的 API 调用
        axios.post('/api/users', this.newUser).then(() => {
          this.addUserDialogVisible = false;
          this.fetchUsers();
        });
      },
      editUser(user) {
        this.$router.push({ name: 'UsersEdit', params: { id: user.id } });
      },
      deleteUser(user) {
        // 替换为实际的 API 调用
        axios.delete(`/api/users/${user.id}`).then(() => {
          this.fetchUsers();
        });
      },
    },
    created() {
      this.fetchUsers();
    },
  };
  </script>
  
  <style>
  .users-container {
    padding: 20px;
  }
  </style>
  