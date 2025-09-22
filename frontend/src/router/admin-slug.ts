import type { RouteRecordRaw } from 'vue-router'

export default [
  // {
  //   path: "/:admin",
  //   component: ()=> import("@/layouts/AuthLayout.vue"),
  //   children: [
  //      {
  //       path: "/login",
  //       component: ()=> import("@/pages/LoginPage.vue")
  //      }
  //   ]
  // },
  {
    path: '/:admin/login',
    name: 'login',
    component: () => import('@/pages/login.vue'),
  },
  {
    path: '/:admin',
    component: () => import('@/layouts/Admin.vue'),
    children: [
      {
        name: 'dasboard',
        path: '',
        component: () => import('@/pages/home.vue'),
      },
      {
        name: 'sales',
        path: 'sales',
        component: () => import('@/pages/sales/sales.vue'),
      },
      {
        name: 'analytics',
        path: 'analytics',
        component: () => import('@/pages/analytics/analytics.vue'),
      },
    ],
  },
] as RouteRecordRaw[]
