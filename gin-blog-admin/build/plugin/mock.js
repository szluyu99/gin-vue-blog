import { viteMockServe } from 'vite-plugin-mock'

export function configMockPlugin(isBuild) {
  // https://github.com/vbenjs/vite-plugin-mock
  return viteMockServe({
    mockPath: 'mock/api', // 模拟 .js 文件的存储文件夹
    localEnabled: !isBuild, // 是否启用本地 xxxx.js 文件
    prodEnabled: isBuild, // 打包是否启用 mock 功能
    // 这样可以控制关闭 mock 的时候不让 mock 打包到最终代码内
    injectCode: `
      import { setupProdMockServer } from '../mock';
      setupProdMockServer();
    `,
  })
}
