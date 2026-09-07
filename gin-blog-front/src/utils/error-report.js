/**
 * 前端错误上报
 *
 * 未捕获异常原来只进浏览器 console, 博主不在场就等于没发生。这里把
 * window.onerror 和 unhandledrejection 收集起来 POST 到后端一张表,
 * 后台「前端错误」页能看到, 不引入 Sentry 这类外部服务。
 *
 * 只报未捕获的: 业务里已经 try/catch 过的错误是预期内的, 不该占位。
 */

// 同一条错误在单次会话里只报一次: 上报接口有按 IP 配额,
// 一个死循环里的报错会把配额吃光, 真正需要看的错误反而被丢掉
const reported = new Set()

// 上报路径不能走 api.js: 那里的响应拦截器出错时会弹提示,
// 而上报本身失败绝不能再产生用户可见的动静
const REPORT_URL = `${import.meta.env.VITE_API}/front/error/report`

function post(payload) {
  // 用 fetch 而不是 axios: 绕开项目里的拦截器(错误提示 / token 处理),
  // keepalive 让页面正在卸载时的上报也有机会发出去
  return fetch(REPORT_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
    keepalive: true,
  }).catch(() => {
    // 静默: 上报失败再报错就成环了
  })
}

/**
 * 上报一条错误, 同一条在本次会话里只发一次
 * @param {string} source front | admin
 * @param {string} message 错误信息
 * @param {string} stack 调用栈
 */
export function reportError(source, message, stack = '') {
  if (!message) {
    return
  }

  // 去重键带上栈顶: message 相同但来源不同的问题仍然分开
  const key = `${message}|${stack.split('\n')[0] ?? ''}`
  if (reported.has(key)) {
    return
  }
  reported.add(key)

  return post({
    source,
    message: String(message),
    stack: String(stack),
    url: window.location.href,
  })
}

/**
 * 挂上全局错误监听
 *
 * mock 模式(GitHub Pages 演示站)不上报: 那里没有后端, 请求必然失败。
 * @param {string} source front | admin
 * @returns {Function} 取消监听, 供测试与热更新清理用(监听器留在 window 上会重复上报)
 */
export function setupErrorReport(source) {
  if (import.meta.env.VITE_USE_MOCK === 'true') {
    return () => {}
  }

  // 同步错误与资源加载错误
  const onError = (event) => {
    // 资源加载失败(img/script)没有 error 对象, event.message 也是空的, 跳过
    if (!event.error && !event.message) {
      return
    }
    reportError(source, event.message || String(event.error), event.error?.stack ?? '')
  }

  // 没有 catch 的 Promise
  const onRejection = (event) => {
    const reason = event.reason
    const message = reason?.message ?? String(reason ?? 'unhandledrejection')
    reportError(source, message, reason?.stack ?? '')
  }

  window.addEventListener('error', onError)
  window.addEventListener('unhandledrejection', onRejection)

  return () => {
    window.removeEventListener('error', onError)
    window.removeEventListener('unhandledrejection', onRejection)
  }
}
