import { getToken } from '@/utils'

// 没有 Token 也可访问的白名单页面
const WHITE_LIST = ['/login', '/404']

// 根据有无 Token 判断能否访问页面
export function createPermissionGuard(router) {
  // 路由前置守卫: 根据有没有 Token 判断前往哪个页面
  router.beforeEach(async (to) => {
    const token = getToken()
    // 没有 Token 的情况
    if (!token) {
      // 前往白名单中的页面可以放行
      if (WHITE_LIST.includes(to.path))
        return true
      // 重定向到登录页, 并且携带 redirect 参数, 登录后自动重定向到原本的目标页面
      return { path: 'login', query: { ...to.query, redirect: to.path } }
    }
    // 有 Token 的情况
    if (to.path === '/login')
      return { path: '/' }

    // TODO: 刷新 Token
    // await refreshAccessToken()

    return true
  })
}
