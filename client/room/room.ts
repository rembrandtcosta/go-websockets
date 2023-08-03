import './assets/main.css'

import { createApp } from 'vue'
import Room from './Room.vue'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/room/:room', component: Room }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})


const app = createApp(Room)

app.use(router)

app.mount('#app')
