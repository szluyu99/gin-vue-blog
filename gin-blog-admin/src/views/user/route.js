const Layout = () => import('@/layout/index.vue')

export default {
  name: 'User',
  path: '/user',
  component: Layout,
  redirect: '/user/list',
  meta: {
    title: '用户管理',
    icon: 'ph:user-list-bold',
    order: 5,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'UserList',
      path: 'list',
      component: () => import('./list/index.vue'),
      meta: {
        title: '用户列表',
        icon: 'mdi:account',
        keepAlive: true,
      },
    },
    {
      name: 'OnlineUserList',
      path: 'online',
      component: () => import('./online/index.vue'),
      meta: {
        title: '在线用户',
        icon: 'ic:outline-online-prediction',
        keepAlive: true,
      },
    },
  ],
}
