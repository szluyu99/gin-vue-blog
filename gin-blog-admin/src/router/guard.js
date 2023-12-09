import { getToken } from '@/utils'

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

/**
 * 根据有无 Token 判断能否访问页面
 */
function createPermissionGuard(router) {
  // 路由前置守卫: 根据有没有 Token 判断前往哪个页面
  router.beforeEach(async (to) => {
    const token = getToken()
    if (!token) {
      // 登录 和 404 不需要 Token
      if (['/login', '/404'].includes(to.path)) {
        return true
      }
      // 重定向到登录页, 并且携带 redirect 参数, 登录后自动重定向到原本的目标页面
      return { path: 'login', query: { ...to.query, redirect: to.path } }
    }
    if (to.path === '/login') {
      return { path: '/' }
    }

    // TODO: 刷新 Token
    // await refreshAccessToken()

    return true
  })
}

const baseTitle = import.meta.env.VITE_TITLE

/**
 * 根据路由元信息设置页面标题
 */
function createPageTitleGuard(router) {
  router.afterEach((to) => {
    const pageTitle = to.meta?.title
    document.title = pageTitle ? `${pageTitle} | ${baseTitle}` : baseTitle
  })
}
