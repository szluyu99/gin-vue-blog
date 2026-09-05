import { Icon } from '@iconify/vue'
import dayjs from 'dayjs'
import { NIcon } from 'naive-ui'
import { h } from 'vue'

export * from './http'
export * from './local'
export * from './naiveTool'

// 图片占位图: 内联 SVG, 不走网络
// 原来这里返回 https://dummyimage.com/400x400 —— 图片本身挂了之后, 兜底又依赖一次外网
export const IMG_PLACEHOLDER = 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="80" height="80"><rect width="80" height="80" fill="%23f3f4f6"/><circle cx="30" cy="30" r="7" fill="%23d1d5db"/><path d="M16 60l16-20 11 13 9-9 12 16z" fill="%23d1d5db"/></svg>'

// 相对图片地址 => 可访问的图片路径, 用于本地文件上传
// 如果包含 http 说明是 Web 图片资源, 原样返回
// 否则是后端服务器上的图片, 返回根相对路径, 由 vite dev proxy / nginx 转发到后端
// 不能拼 VITE_SERVER_URL: 那里写的是 localhost, 从别的机器访问页面时
// localhost 指的是浏览器所在的机器, 图片必然裂开
export function convertImgUrl(imgUrl) {
  if (!imgUrl) {
    return IMG_PLACEHOLDER
  }
  // 网络资源
  if (imgUrl.startsWith('http')) {
    return imgUrl
  }
  return `/${imgUrl.replace(/^\/+/, '')}`
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

/**
 * 安全解析 JSON 字符串, 解析失败或非字符串时返回 fallback
 * 后端部分操作日志的 request_param 是空串, 上传接口被网关拦截时也可能返回 HTML
 */
export function parseJson(str, fallback = null) {
  if (typeof str !== 'string' || !str.trim()) {
    return fallback
  }
  try {
    return JSON.parse(str)
  }
  catch {
    return fallback
  }
}

/**
 * 格式化 JSON 字符串用于展示, 无法解析时原样返回, 不抛异常
 */
const PARSE_FAILED = Symbol('parseFailed')
export function formatJson(str) {
  const parsed = parseJson(str, PARSE_FAILED)
  return parsed === PARSE_FAILED ? (str ?? '') : JSON.stringify(parsed, null, 2)
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
