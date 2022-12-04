import { getProxyConfig } from '../../setting' // vite.config.js 中加载的资源无法使用别名路径

// https://cn.vitejs.dev/config/server-options.html#server-proxy
/**
 * @param {*} isUseProxy 是否使用代理
 * @param {*} proxyType 当前模式
 */
export function createViteProxy(isUseProxy = true, proxyType) {
  if (!isUseProxy)
    return undefined

  const proxyConfig = getProxyConfig(proxyType)
  return {
    [proxyConfig.prefix]: {
      target: proxyConfig.target,
      changeOrigin: true,
      rewrite: path => path.replace(new RegExp(`^${proxyConfig.prefix}`), ''),
    },
  }
}
