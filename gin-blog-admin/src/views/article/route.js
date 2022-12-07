const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Article',
  path: '/article',
  component: Layout,
  redirect: '/article/list',
  meta: {
    title: '文章管理',
    icon: 'ic:twotone-article',
    order: 2,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'ArticleList',
      path: 'list',
      component: () => import('./list/index.vue'),
      meta: {
        title: '文章列表',
        icon: 'material-symbols:format-list-bulleted',
        // role: ['admin'],
        // requireAuth: true,
        keepAlive: true,
      },
    },
    {
      name: 'ArticleWrite',
      path: 'write',
      component: () => import('./write/index.vue'),
      meta: {
        title: '发布文章',
        icon: 'icon-park-outline:write',
        // role: ['admin'],
        // requireAuth: true,
        keepAlive: true,
      },
    },
    {
      name: 'ArticleEdit',
      path: 'write/:id',
      component: () => import('./write/index.vue'),
      isHidden: true,
      meta: {
        title: '编辑文章',
        icon: 'icon-park-outline:write',
        // role: ['admin'],
        // requireAuth: true,
        // keepAlive: true,
      },
    },
    {
      name: 'CategoryList',
      path: 'category-list',
      component: () => import('./category/index.vue'),
      meta: {
        title: '分类管理',
        icon: 'tabler:category',
        // role: ['admin'],
        // requireAuth: true,
        keepAlive: true,
      },
    },
    {
      name: 'TagList',
      path: 'tag-list',
      component: () => import('./tag/index.vue'),
      meta: {
        title: '标签管理',
        icon: 'tabler:tag',
        keepAlive: true,
      },
    },
  ],
}
