const proxyConfigMappings: Record<ProxyType, ProxyConfig> = {
  // 开发环境调用的接口
  dev: {
    prefix: '/api/front',
    target: 'http://localhost:5678/api/front',
  },
  // 生产环境调用的接口
  prod: {
    prefix: '/api/front',
    target: 'http://localhost:5678/api/front',
    // target: 'https://hahacode.cn/api/front',
  },
  test: {
    prefix: '/api/front',
    target: 'http://localhost:5678/api/front',
  },
}

export function getProxyConfig(envType: ProxyType = 'dev'): ProxyConfig {
  return proxyConfigMappings[envType]
}
