const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Message',
  path: '/message',
  component: Layout,
  redirect: '/message/comment',
  meta: {
    title: '消息管理',
    icon: 'ic:twotone-email',
    order: 3,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'CommentList',
      path: 'comment',
      component: () => import('./comment/index.vue'),
      meta: {
        title: '评论管理',
        icon: 'ic:twotone-comment',
        keepAlive: true,
      },
    },
    {
      name: 'LeaveMsgList',
      path: 'leave-msg',
      component: () => import('./leave-msg/index.vue'),
      meta: {
        title: '留言管理',
        icon: 'ic:twotone-message',
        keepAlive: true,
      },
    },
  ],
}
