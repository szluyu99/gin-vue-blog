// 相对图片地址 => 完整的图片路径, 用于本地文件上传
// - 如果包含 http 说明是 Web 图片资源
// - 否则是服务器上的图片，需要拼接服务器路径
const SERVER_URL = import.meta.env.VITE_BACKEND_URL

/**
 * 将相对地址转换为完整的图片路径
 * @param {string} imgUrl
 * @returns {string} 完整的图片路径
 */
export function convertImgUrl(imgUrl) {
  if (!imgUrl) {
    return 'http://dummyimage.com/400x400'
  }
  // 网络资源
  if (imgUrl.startsWith('http')) {
    return imgUrl
  }
  // 服务器资源
  return `${SERVER_URL}/${imgUrl}`
}

export * from './local'
export * from './http'
