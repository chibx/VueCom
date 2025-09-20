import { createRouter, createWebHistory } from 'vue-router'
import adminSlug from './admin-slug'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [...adminSlug],
})

export default router
