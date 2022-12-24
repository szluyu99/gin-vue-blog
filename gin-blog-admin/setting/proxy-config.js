const proxyConfigMappings = {
  // 开发时调用本地 API
  dev: {
    prefix: '/api',
    target: 'http://localhost:8765/api/',
  },
  // 这里写线上接口, 打包后调用线上 API
  prod: {
    prefix: '/api',
    // target: 'https://hahacode.cn/api/',
    target: 'http://localhost:8765/api/',
  },
  test: {
    prefix: '/api',
    target: 'http://localhost:8765/api/',
  },
}

export function getProxyConfig(envType = 'dev') {
  return proxyConfigMappings[envType]
}
