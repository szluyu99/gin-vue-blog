import { unref } from 'vue'
import { router } from '@/router'

/**
 * @description 跳转到登录页面
 */
export function toLogin() {
  const currentRoute = unref(router.currentRoute)
  // 跳转回去时记录 redirect 到 query 上
  const needRedirect
    = !currentRoute.meta.requireAuth && !['/404', '/login'].includes(currentRoute.path)
  router.replace({
    path: '/login',
    query: needRedirect ? { ...currentRoute.query, redirect: currentRoute.path } : {},
  })
}

/**
 * @description 跳转到 404 页面
 */
export function to404() {
  router.replace({ path: '/404' })
}
