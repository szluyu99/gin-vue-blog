// https://cn.vitejs.dev/config/server-options.html#server-proxy
import type { ProxyOptions } from 'vite'
import { getProxyConfig } from '../../settings/proxy-config'

/**
 * @param {*} isUseProxy 是否使用代理
 * @param {*} proxyType 当前模式
 */
export function createViteProxy(isUseProxy = true, proxyType: ProxyType) {
  if (!isUseProxy)
    return undefined

  const proxyConfig = getProxyConfig(proxyType)
  const proxy: Record<string, string | ProxyOptions> = {
    [proxyConfig.prefix]: {
      target: proxyConfig.target,
      changeOrigin: true,
      rewrite: (path: string) => path.replace(new RegExp(`^${proxyConfig.prefix}`), ''),
    },
  }
  return proxy
}
