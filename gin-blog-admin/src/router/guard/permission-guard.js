import { getToken, refreshAccessToken } from '@/utils'

const WHITE_LIST = ['/login', '/404']

// 项目采用前端角色权限, 即根据登录角色来判断过滤路由, 将没有权限的路由删除
export function createPermissionGuard(router) {
  // 路由前置守卫: 根据有没有 token 判断前往哪个页面
  router.beforeEach(async (to) => {
    // console.log('permissionGuard before')
    const token = getToken()
    // 没有 token
    if (!token) {
      // 前往白名单中的页面可以放行
      if (WHITE_LIST.includes(to.path)) return true
      // 重定向到登录页, 并且携带 redirect 参数, 登录后自动重定向到原本的目标页面
      return { path: 'login', query: { ...to.query, redirect: to.path } }
    }
    // 有 token
    if (to.path === '/login') return { path: '/' }
    // 刷新 token
    await refreshAccessToken()
    return true
  })
}
