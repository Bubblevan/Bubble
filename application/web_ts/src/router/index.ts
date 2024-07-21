import { createRouter, createWebHistory } from 'vue-router';
import Login from '../components/Login.vue';
import Register from '../components/Register.vue';
import Main from '../components/Main.vue';
import Events from '../components/Events.vue';
import CreateEvents from '../components/CreateEvents.vue';
import Profile from '../components/Profile.vue';
import Purchased from '../components/Purchase.vue'
import Users from '../components/Users.vue'
import UsersEdit from '../components/UserEdit.vue'
import Buy from '../components/Buy.vue'

const routes = [
  { path: '/', component: Login },
  { path: '/register', component: Register },
  { path: '/main', component: Main },
  { path: '/events', name: 'Events', component: Events },
  { path: '/events/create', name: 'CreateEvents', component: CreateEvents },
  { path: '/users', name: 'Users', component: Users },
  { path: '/users/edit:id', name: 'UsersEdit', component: UsersEdit, props: true },
  { path: '/profile', name: 'Profile', component: Profile },
  { path: '/orders', name: 'Purchased', component: Purchased },
  { path: '/events/buy/:id', name: 'Buy', component: Buy, props: true },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
