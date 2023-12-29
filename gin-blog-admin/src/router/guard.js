import { useAuthStore } from '@/store'

export function setupRouterGuard(router) {
  createPageLoadingGuard(router)
  createPermissionGuard(router)
  createPageTitleGuard(router)
}

/**
 * 根据导航设置顶部加载条的状态
 */
function createPageLoadingGuard(router) {
  router.beforeEach(() => window.$loadingBar?.start())
  router.afterEach(() => setTimeout(() => window.$loadingBar?.finish(), 200))
  // 在导航期间每次发生未捕获的错误时都会调用该处理程序
  router.onError(() => window.$loadingBar?.error())
}

// 无需 Token 也能访问的页面
const WHITE_LIST = ['/login', '/404']
/**
 * 根据有无 Token 判断能否访问页面
 */
function createPermissionGuard(router) {
  // const base = import.meta.env.VITE_BASE_URL
  // 路由前置守卫: 根据有没有 Token 判断前往哪个页面
  router.beforeEach(async (to) => {
    if (WHITE_LIST.includes(to.path)) {
      return true
    }

    const { accessToken } = useAuthStore()

    // 没有 Token
    if (!accessToken) {
      window.$message.error('没有 Token，请重新登录！')
      // TODO: 重定向到登录页, 并且携带 redirect 参数, 登录后自动重定向到原本的目标页面
      // if (to.path !== `${base}/`) {
      //   return { path: `${base}/login`, query: { ...to.query, redirect: to.path } }
      // }
      // else {
      //   return { path: `${base}/login`, query: { ...to.query } }
      // }
      return { name: 'Login' }
    }

    // 有 Token 的时候无需访问登录页面
    if (to.name === 'Login') {
      return { path: '/' }
    }

    // 能在路由中找到, 则正常访问
    if (router.getRoutes().find(e => e.name === to.name)) {
      return true
    }

    // TODO: 刷新 Token
    // await refreshAccessToken()

    // TODO: 判断是无权限还是 404
    // return { name: '404', query: { path: to.fullPath } }
    return { name: '404' }
  })
}

/**
 * 根据路由元信息设置页面标题
 */
function createPageTitleGuard(router) {
  const baseTitle = import.meta.env.VITE_TITLE
  router.afterEach((to) => {
    const pageTitle = to.meta?.title
    document.title = pageTitle ? `${pageTitle} | ${baseTitle}` : baseTitle
  })
}
