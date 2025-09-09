import type { RouteRecordRaw } from 'vue-router'

export default [
  {
    path: '/:admin',
    children: [
      {
        name: 'dasboard',
        path: '',
        component: () => import('@/pages/HomePage.vue'),
      },
      {
        name: 'sales',
        path: 'sales',
        component: () => import('@/pages/sales/SalesPage.vue'),
      },
      {
        name: 'analytics',
        path: 'analytics',
        component: () => import('@/pages/analytics/AnalyticsPage.vue'),
      },
    ],
  },
  {
    path: '/:path(.*)',
    component: () => import('@/pages/NotFound.vue'),
  },
] as RouteRecordRaw[]
