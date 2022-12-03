const proxyConfigMappings = {
  dev: {
    prefix: '/api',
    target: 'http://localhost:8765/api/',
  },
  test: {
    prefix: '/api',
    target: 'http://localhost:8765/api/',
  },
  prod: {
    prefix: '/api',
    target: 'http://localhost:8765/api/',
  },
}

export function getProxyConfig(envType = 'dev') {
  return proxyConfigMappings[envType]
}
