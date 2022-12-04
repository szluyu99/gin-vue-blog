export function createPageLoadingGuard(router) {
  // 路由前置守卫: 开始加载效果
  router.beforeEach(() => {
    window.$loadingBar?.start()
  })

  // 路由后置钩子: 加载效果结束
  router.afterEach(() => {
    setTimeout(() => {
      window.$loadingBar?.finish()
    }, 200)
  })

  // 在导航期间每次发生未捕获的错误时都会调用该处理程序
  router.onError(() => {
    window.$loadingBar?.error()
  })
}
