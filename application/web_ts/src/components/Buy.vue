<template>
    <Layout>
      <div class="buy-container">
        <h1>购买门票</h1>
        <el-form :model="ticket" label-width="120px">
          <el-form-item label="活动名称">
            <el-input v-model="ticket.EventName" disabled></el-input>
          </el-form-item>
          <el-form-item label="场地">
            <el-input v-model="ticket.Venue" disabled></el-input>
          </el-form-item>
          <el-form-item label="活动日期">
            <el-input v-model="formattedDate" disabled></el-input>
          </el-form-item>
          <el-form-item label="价格">
            <el-input v-model="ticket.Price" disabled></el-input>
          </el-form-item>
          <el-form-item label="剩余票数">
            <el-input v-model="ticket.Num" disabled></el-input>
          </el-form-item>
          <el-form-item label="购买数量">
            <el-input v-model="quantity"></el-input>
          </el-form-item>
        </el-form>
        <el-button type="primary" @click="buyTicket">购买</el-button>
      </div>
    </Layout>
  </template>
  
  <script>
  import Layout from '../layout/Layout.vue';
  import axios from 'axios';
  import { format } from 'date-fns';
  
  export default {
    components: {
      Layout
    },
    data() {
      return {
        ticket: {},
        quantity: 1,
        formattedDate: ''
      };
    },
    created() {
      this.fetchTicket();
    },
    methods: {
      async fetchTicket() {
        const id = this.$route.params.id;
        try {
          const response = await axios.get(`http://localhost:8080/tickets/${id}`);
          this.ticket = response.data.data;
          this.formattedDate = format(new Date(this.ticket.EventDate), 'yyyy-MM-dd');
        } catch (error) {
          this.$message.error('获取门票详情失败');
        }
      },
      async buyTicket() {
        if (this.quantity > this.ticket.Num) {
          this.$message.error('购买数量超过剩余票数');
          return;
        }
  
        const order = {
          userID: 1, // 应从实际登录的用户获取
          ticketID: this.ticket.id,
          quantity: this.quantity,
          totalPrice: this.quantity * this.ticket.Price,
          orderDate: new Date()
        };
  
        try {
          await axios.post('http://localhost:8080/orders', order);
          this.$message.success('购买成功');
          this.$router.push('/events');
        } catch (error) {
          this.$message.error('购买失败');
        }
      }
    }
  };
  </script>
  
  <style scoped>
  .buy-container {
    padding: 20px;
  }
  </style>
  