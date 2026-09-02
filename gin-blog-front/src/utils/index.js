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

/**
 * 去掉 Markdown 记号, 用于列表页展示纯文本摘要
 *
 * 以前的做法是 marked(md) 转成 HTML 再用正则剥标签, 为了这一个用途把整个
 * marked (gzip 12KB) 打进了首页 chunk, 这里用正则直接处理。
 * @param {string} md
 * @returns {string} 纯文本
 */
export function stripMarkdown(md) {
  if (!md) {
    return ''
  }

  return md
    .replace(/```[\s\S]*?```/g, ' ') // 代码块
    .replace(/`([^`]*)`/g, '$1') // 行内代码
    .replace(/!\[[^\]]*\]\([^)]*\)/g, ' ') // 图片
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1') // 链接: 只留文字
    .replace(/^\s{0,3}>+\s?/gm, '') // 引用
    .replace(/^\s{0,3}#{1,6}\s+/gm, '') // 标题
    .replace(/^\s{0,3}([-*+]|\d+\.)\s+/gm, '') // 列表
    .replace(/^\s{0,3}([-*_]\s?){3,}$/gm, ' ') // 分隔线
    .replace(/(\*\*|__)(.*?)\1/g, '$2') // 加粗
    .replace(/(\*|_)(.*?)\1/g, '$2') // 斜体
    .replace(/~~(.*?)~~/g, '$1') // 删除线
    .replace(/<[^>]+>/g, '') // 内联 HTML
    .replace(/\s+/g, ' ') // 折叠空白
    .trim()
}

/**
 * 取一句"一言", 失败返回 null (由调用方决定兜底文案)
 *
 * 首页的 HomeBanner 和 TalkingCarousel 都要用, 以前各自 fetch 一次, 同一个接口请求了两遍;
 * 这里缓存 Promise, 一次页面加载只请求一次。这个接口在部分网络下不通, 所以内部吃掉异常。
 * @returns {Promise<string | null>}
 */
let sentencePromise = null
export function getOneSentence() {
  if (!sentencePromise) {
    sentencePromise = fetch('https://v1.hitokoto.cn?c=i')
      .then(res => res.json())
      .then(data => data?.hitokoto || null)
      .catch(() => null)
  }
  return sentencePromise
}

export * from './http'
