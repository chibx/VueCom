import { createRouter, createWebHistory } from 'vue-router'
import adminSlug from './admin-slug'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    ...adminSlug,
    {
      path: '/:path(.*)',
      component: () => import('@/pages/NotFound.vue'),
    },
  ],
})

export default router
