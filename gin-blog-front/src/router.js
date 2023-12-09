import { createRouter, createWebHistory } from 'vue-router'

import NProgress from 'nprogress'
import './styles/nprogress.css'

const basicRoutes = [
  {
    name: 'Home',
    path: '/',
    component: () => import('@/views/home/index.vue'),
  },
  {
    name: 'Article',
    path: '/article/:id',
    component: () => import('@/views/article/detail/index.vue'),
  },
  {
    name: 'Archive',
    path: '/archives',
    component: () => import('@/views/discover/archive/index.vue'),
    meta: {
      title: '归档',
    },
  },
  {
    name: 'Category',
    path: '/categories',
    component: () => import('@/views/discover/category/index.vue'),
    meta: {
      title: '分类',
    },
  },
  {
    name: 'CategoryArticles',
    path: '/categories/:categoryId',
    component: () => import('@/views/article/list/index.vue'),
    meta: {
      title: '分类',
    },
  },
  {
    name: 'Tag',
    path: '/tags',
    component: () => import('@/views/discover/tag/index.vue'),
    meta: {
      title: '标签',
    },
  },
  {
    name: 'TagArticles',
    path: '/tags/:tagId',
    component: () => import('@/views/article/list/index.vue'),
    meta: {
      title: '标签',
    },
  },
  {
    name: 'Album',
    path: '/albums',
    component: () => import('@/views/entertainment/album/index.vue'),
    meta: {
      title: '相册',
    },
  },
  {
    name: 'Link',
    path: '/links',
    component: () => import('@/views/link/index.vue'),
    meta: {
      title: '友情链接',
    },
  },
  {
    name: 'About',
    path: '/about',
    component: () => import('@/views/about/index.vue'),
    meta: {
      title: '关于我',
    },
  },
  {
    name: 'MessageBoard',
    path: '/message',
    component: () => import('@/views/message/index.vue'),
    meta: {
      title: '留言',
    },
  },
  {
    name: 'User',
    path: '/user',
    component: () => import('@/views/user/index.vue'),
    meta: {
      title: '个人中心',
    },
  },
  {
    name: '404',
    path: '/404',
    component: () => import('@/views/error-page/404.vue'),
  },
  // 无匹配路由跳转 404
  {
    name: 'NotFound',
    path: '/:pathMatch(.*)*',
    redirect: '/404',
    isHidden: true,
  },
]

export const router = createRouter({
  history: createWebHistory('/'),
  routes: basicRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

router.afterEach((to) => {
  document.title = `${to.meta?.title ?? import.meta.env.VITE_APP_TITLE}`
})

NProgress.configure({ showSpinner: false })

router.beforeEach((to, from, next) => {
  NProgress.start()
  for (let i = 0; i < 5; i++) NProgress.inc()
  setTimeout(() => NProgress.done(), 300)
  next()
})
