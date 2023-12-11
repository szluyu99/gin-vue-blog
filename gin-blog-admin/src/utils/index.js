import { h } from 'vue'
import { Icon } from '@iconify/vue'
import { NIcon } from 'naive-ui'
import dayjs from 'dayjs'

export * from './http'
export * from './local'
export * from './naiveTool'

// 相对图片地址 => 完整的图片路径, 用于本地文件上传
// 如果包含 http 说明是 Web 图片资源
// 否则是服务器上的图片，需要拼接服务器路径
const SERVER_URL = import.meta.env.VITE_SERVER_URL
export function convertImgUrl(imgUrl) {
  if (!imgUrl) {
    return 'http://dummyimage.com/400x400'
  }
  // 网络资源
  if (imgUrl.startsWith('http')) {
    return imgUrl
  }
  return `${SERVER_URL}/${imgUrl}`
}

/**
 * 格式化时间
 */
export function formatDate(date = undefined, format = 'YYYY-MM-DD') {
  return dayjs(date).format(format)
}

/**
 * 使用 NIcon 渲染图标
 */
export function renderIcon(icon, props = { size: 12 }) {
  return () => h(NIcon, props, { default: () => h(Icon, { icon }) })
}

// 前端导出, 传入文件内容和文件名称
export function downloadFile(content, fileName) {
  const aEle = document.createElement('a') // 创建下载链接
  aEle.download = fileName // 设置下载的名称
  aEle.style.display = 'none'// 隐藏的可下载链接
  // 字符内容转变成 blob 地址
  const blob = new Blob([content])
  aEle.href = URL.createObjectURL(blob)
  // 绑定点击时间
  document.body.appendChild(aEle)
  aEle.click()
  // 然后移除
  document.body.removeChild(aEle)
}
