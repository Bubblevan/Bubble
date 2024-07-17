<template>
    <Layout>
      <div class="purchase-container">
        <h1>我的购买</h1>
        <el-table :data="orders">
          <el-table-column prop="ticketName" label="活动名称"></el-table-column>
          <el-table-column prop="venue" label="场地"></el-table-column>
          <el-table-column prop="eventDate" label="活动日期"></el-table-column>
          <el-table-column prop="quantity" label="购买数量"></el-table-column>
          <el-table-column prop="totalPrice" label="总价格"></el-table-column>
        </el-table>
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
        orders: []
      };
    },
    created() {
      this.fetchOrders();
    },
    methods: {
      async fetchOrders() {
        try {
          const token = sessionStorage.getItem('token');
          const userID = sessionStorage.getItem('userID');
          const response = await axios.get(`http://localhost:8080/orders/${userID}`, {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          });
          this.orders = response.data.data;
        } catch (error) {
          this.$message.error('获取购买历史失败');
        }
      }
    }
  };
  </script>
  
  <style scoped>
  .purchase-container {
    max-width: 800px;
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
  
  .el-table {
    width: 100%;
  }
  </style>
  