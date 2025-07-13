import { createRouter, createWebHistory } from 'vue-router'

import { basicRoutes } from './routes'
import { setupRouterGuard } from './guard'

import { useAuthStore, usePermissionStore, useUserStore } from '@/store'

export const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_PUBLIC_PATH), // '/admin'
  routes: basicRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

/**
 * 初始化路由
 */
export async function setupRouter(app) {
  await addDynamicRoutes()
  setupRouterGuard(router)
  app.use(router)
}

/**
 * ! 添加动态路由: 根据配置由前端或后端生成路由
 */
export async function addDynamicRoutes() {
  const authStore = useAuthStore()

  if (!authStore.token) {
    authStore.toLogin()
    return
  }

  // 有 token 的情况
  try {
    const userStore = useUserStore()
    const permissionStore = usePermissionStore()

    // userId 不存在, 则调用接口根据 token 获取用户信息
    if (!userStore.userId) {
      await userStore.getUserInfo()
    }

    // 根据环境变量中的值决定前端生成路由还是后端路由
    const accessRoutes = JSON.parse(import.meta.env.VITE_BACK_ROUTER)
      ? await permissionStore.generateRoutesBack()
      : permissionStore.generateRoutesFront(['admin'])
  console.log(accessRoutes)
    // 检查重复路由名称并添加路由
    const addedNames = new Set()

    const addRoute = (route) => {
      let routeName = route.name

      // 处理重复的路由名称
      if (addedNames.has(routeName) || router.hasRoute(routeName)) {
        routeName = `${routeName}-${Date.now()}`
        route.name = routeName
      }

      addedNames.add(routeName)

      try {
        router.addRoute(route)
      }
      catch (e) {
        console.warn(`添加路由 ${routeName} 失败:`, e.message)
      }
    }

    // 添加路由
    accessRoutes.forEach(addRoute)
  }
  catch (err) {
    console.error('addDynamicRoutes Error: ', err)
  }
}

/**
 * 重置路由
 */
export async function resetRouter() {
  router.getRoutes().forEach((route) => {
    const name = route.name
    if (!basicRoutes.some(e => e.name === name) && router.hasRoute(name)) {
      router.removeRoute(name)
    }
  })
}
