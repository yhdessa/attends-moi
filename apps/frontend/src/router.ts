import { createRouter, createWebHistory } from 'vue-router'
import BoardList from './views/BoardList.vue'
import BoardView from './views/BoardView.vue'

const routes = [
  { path: '/', component: BoardList },
  { path: '/boards/:id', component: BoardView, props: true },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
