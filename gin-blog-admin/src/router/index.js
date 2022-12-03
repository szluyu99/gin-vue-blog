import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'
import { basicRoutes, EMPTY_ROUTE, NOT_FOUND_ROUTE } from './routes'
import { setupRouterGuard } from './guard'
import { getToken } from '@/utils'
import { useUserStore, usePermissionStore } from '@/store'

const isHash = import.meta.env.VITE_USE_HASH === 'true'
const basePath = import.meta.env.VITE_PUBLIC_PATH // '/blog-admin'
export const router = createRouter({
  history: isHash ? createWebHashHistory(basePath) : createWebHistory(basePath),
  routes: basicRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

export async function setupRouter(app) {
  await addDynamicRoutes()
  setupRouterGuard(router) // 路由导航
  app.use(router)
}

// 用户退出登录时需要重置路由
export async function resetRouter() {
  const basicRouteNames = getRouteNames(basicRoutes)
  // console.log(basicRouteNames)
  router.getRoutes().forEach((route) => {
    const name = route.name
    if (!basicRouteNames.includes(name)) {
      router.removeRoute(name)
    }
  })
}

// 添加动态路由
export async function addDynamicRoutes() {
  const token = getToken()

  // 没有 token 的情况
  if (!token) {
    router.addRoute(EMPTY_ROUTE)
    return
  }

  // 有 token 的情况
  try {
    // 根据权限生成动态路由, 目前是全部加入
    const userStore = useUserStore()
    const permissionStore = usePermissionStore()
    // userId 不存在, 则调用接口根据 token 获取用户信息
    !userStore.userId && (await userStore.getUserInfo())
    // 根据用户角色生成路由, 并添加到当前路由实例中
    // const accessRoutes = permissionStore.generateRoutes(['admin']) // ! 前端控制路由
    const accessRoutes = await permissionStore.generateRoutesFromBack() // ! 后端生成路由
    // 将当前没有的路由添加进去
    accessRoutes.forEach((route) => !router.hasRoute(route.name) && router.addRoute(route))
    // console.log('accessRoutes', accessRoutes)
    // 移除 EMPTY_ROUTE 页面
    router.hasRoute(EMPTY_ROUTE.name) && router.removeRoute(EMPTY_ROUTE.name)
    // 添加 404 页面
    router.addRoute(NOT_FOUND_ROUTE)
    // console.log('route.getRoutes', router.getRoutes())
  } catch (err) {
    console.log(err)
  }
}

export function getRouteNames(routes) {
  return routes.map((route) => getRouteName(route)).flat(1)
}

// 获取 route 下的路由名称
function getRouteName(route) {
  const names = [route.name]
  if (route.children?.length) {
    names.push(...route.children.map((e) => getRouteName(e)).flat(1))
  }
  return names
}
