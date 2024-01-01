const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Home',
  path: '/',
  component: Layout,
  redirect: '/home',
  meta: {
    order: 0,
  },
  isCatalogue: true,
  children: [
    {
      name: 'Home',
      path: 'home',
      component: () => import('./index.vue'),
      meta: {
        title: '首页',
        icon: 'ic:sharp-home',
        order: 0,
      },
    },
  ],
}
