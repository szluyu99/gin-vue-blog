const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Profile',
  path: '/',
  component: Layout,
  redirect: '/profile',
  meta: {
    title: '个人设置',
    icon: 'healthicons:ui-user-profile-outline',
    order: 7,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'Profile',
      path: '/profile',
      component: () => import('./index.vue'),
      meta: {
        title: '个人设置',
        icon: 'mdi:account',
        keepAlive: true,
      },
    },
  ],
}
