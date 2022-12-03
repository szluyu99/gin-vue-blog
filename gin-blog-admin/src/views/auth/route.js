const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Auth',
  path: '/auth',
  component: Layout,
  redirect: '/auth/menu',
  meta: {
    title: '权限管理',
    icon: 'cib:adguard',
    order: 3,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'MenuList',
      path: 'menu',
      component: () => import('./menu/index.vue'),
      meta: {
        title: '菜单管理',
        icon: 'ic:twotone-menu-book',
        keepAlive: true,
      },
    },
    {
      name: 'ResourceList',
      path: 'resource',
      component: () => import('./resource/index.vue'),
      meta: {
        title: '接口管理',
        icon: 'mdi:api',
        keepAlive: true,
      },
    },
    {
      name: 'RoleList',
      path: 'role',
      component: () => import('./role/index.vue'),
      meta: {
        title: '角色管理',
        icon: 'carbon:user-role',
        keepAlive: true,
      },
    },
  ],
}
