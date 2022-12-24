import { createRouter, createWebHistory } from 'vue-router'
import { EMPTY_ROUTE, NOT_FOUND_ROUTE, basicRoutes } from './routes'
import { setupRouterGuard } from './guard'
import { getToken } from '@/utils'
import { usePermissionStore, useUserStore } from '@/store'

export const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_PUBLIC_PATH), // '/admin'
  routes: basicRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

export async function setupRouter(app) {
  await addDynamicRoutes() // 每次刷新时都添加动态路由
  setupRouterGuard(router) // 路由守卫
  app.use(router)
}

// 用户退出登录时需要重置路由
export async function resetRouter() {
  const basicRouteNames = getRouteNames(basicRoutes)
  router.getRoutes().forEach((route) => {
    const name = route.name
    if (!basicRouteNames.includes(name))
      router.removeRoute(name)
  })
}

// ! 添加动态路由: 可以由前端生成或后端生成路由
export async function addDynamicRoutes() {
  const token = getToken()

  // 没有 token 的情况
  if (!token) {
    router.addRoute(EMPTY_ROUTE)
    return
  }

  // 有 token 的情况
  try {
    // 根据权限生成动态路由
    const userStore = useUserStore()
    const permissionStore = usePermissionStore()
    // userId 不存在, 则调用接口根据 token 获取用户信息
    !userStore.userId && (await userStore.getUserInfo())

    // 根据环境变量中的值决定前端生成路由还是后端路由
    const accessRoutes = (JSON.parse(import.meta.env.VITE_BACK_ROUTER))
      ? await permissionStore.generateRoutesBack() // ! 后端生成路由
      : permissionStore.generateRoutesFront(['admin']) // ! 前端生成路由 (根据角色), 待完善

    // 将当前没有的路由添加进去
    accessRoutes.forEach(route => !router.hasRoute(route.name) && router.addRoute(route))
    // 移除 EMPTY_ROUTE 页面
    router.hasRoute(EMPTY_ROUTE.name) && router.removeRoute(EMPTY_ROUTE.name)
    // 添加 404 页面
    router.addRoute(NOT_FOUND_ROUTE)
  }
  catch (err) {
    console.error('addDynamicRoutes Error: ', err)
  }
}

export function getRouteNames(routes) {
  return routes.map(route => getRouteName(route)).flat(1)
}

// 获取 route 下的路由名称
function getRouteName(route) {
  const names = [route.name]
  if (route.children?.length)
    names.push(...route.children.map(e => getRouteName(e)).flat(1))
  return names
}
