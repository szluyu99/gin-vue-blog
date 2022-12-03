import type { Router } from 'vue-router'

export function setupRouterGuard(router: Router) {
  router.afterEach((to) => {
    document.title = `${to.meta?.title ?? import.meta.env.VITE_APP_TITLE}`
  })
}
