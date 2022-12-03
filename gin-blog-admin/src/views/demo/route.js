const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Demo',
  path: '/demo',
  component: Layout,
  redirect: '/demo/crud-table',
  isHidden: true,
  meta: {
    title: '示例页面',
    customIcon: 'logo',
    // // role: ['admin'],
    // requireAuth: true,
    order: 9,
  },
  children: [
    {
      name: 'CrudTable',
      path: 'crud-table',
      component: () => import('./table/index.vue'),
      meta: {
        title: 'CRUD表格',
        icon: 'ic:baseline-table-view',
        // role: ['admin'],
        // requireAuth: true,
        keepAlive: true, // 切换到其他标签时, 缓存该标签的内容
      },
    },
    {
      name: 'MDEditor',
      path: 'md-editor',
      component: () => import('./editor/md-editor.vue'),
      meta: {
        title: 'MD编辑器',
        icon: 'ri:markdown-line',
        // role: ['admin'],
        // requireAuth: true,
        keepAlive: true, // 切换到其他标签时, 缓存该标签的内容
      },
    },
    {
      name: 'RichTextEditor',
      path: 'rich-text',
      component: () => import('./editor/rich-text.vue'),
      meta: {
        title: '富文本编辑器',
        icon: 'ic:sharp-text-rotation-none',
        // role: ['admin'],
        // requireAuth: true,
        keepAlive: true, // 切换到其他标签时, 缓存该标签的内容
      },
    },
  ],
}
