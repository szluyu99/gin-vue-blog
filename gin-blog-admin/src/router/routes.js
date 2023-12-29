// const Layout = () => import('@/layout/index.vue')

// 基础路由, 无需权限, 总是会注册到最终路由中
export const basicRoutes = [
  // TODO: 如何去除这个代码？
  // {
  //   name: '',
  //   path: '/',
  //   component: Layout,
  //   redirect: '/home', // 默认跳转首页
  //   isHidden: true,
  //   meta: { order: 0 },
  // },
  {
    name: 'Login',
    path: '/login',
    component: () => import('@/views/Login.vue'),
    isHidden: true,
    meta: {
      title: '登录页',
    },
  },
  {
    name: '404',
    path: '/404',
    component: () => import('@/views/error-page/404.vue'),
    isHidden: true,
    meta: {
      title: '错误页',
    },
  },
]

// 前端控制路由: 加载 views 下每个模块的 routes.js 文件
const routeModules = import.meta.glob('@/views/**/route.js', { eager: true })
const asyncRoutes = []
Object.keys(routeModules).forEach((key) => {
  asyncRoutes.push(routeModules[key].default)
})

// 加载 views 下每个模块的 index.vue 文件
const vueModules = import.meta.glob('@/views/**/index.vue')

export {
  asyncRoutes,
  vueModules,
}
