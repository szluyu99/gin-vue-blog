const baseTitle = import.meta.env.VITE_TITLE

// 路由后置钩子: 更改页面标题
export function createPageTitleGuard(router) {
  router.afterEach((to) => {
    const pageTitle = to.meta?.title
    document.title = pageTitle ? `${pageTitle} | ${baseTitle}` : baseTitle
  })
}
