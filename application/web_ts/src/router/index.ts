import { createRouter, createWebHistory } from 'vue-router';
import Login from '../components/Login.vue';
import Register from '../components/Register.vue';
import Main from '../components/Main.vue';
import Events from '../components/Events.vue';
import CreateEvents from '../components/CreateEvents.vue';
import Profile from '../components/Profile.vue';
import Purchase from '../components/Purchase.vue'

const routes = [
  { path: '/', component: Login },
  { path: '/register', component: Register },
  { path: '/main', component: Main },
  { path: '/events', name: 'Events', component: Events },
  { path: '/events/create', name: 'CreateEvents', component: CreateEvents },
  { path: '/profile', name: 'Profile', component: Profile },
  { path: '/orders', name: 'Purchase', component: Purchase },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
