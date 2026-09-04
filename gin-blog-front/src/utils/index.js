// 相对图片地址 => 可访问的图片路径, 用于本地文件上传
// - 如果包含 http 说明是 Web 图片资源, 原样返回
// - 否则是后端服务器上的图片, 返回根相对路径, 由 vite dev proxy / nginx 转发到后端
//   不能拼 VITE_BACKEND_URL: 那里写的是 localhost, 从别的机器访问页面时
//   localhost 指的是浏览器所在的机器, 图片必然裂开
/**
 * 将相对地址转换为可访问的图片路径
 * @param {string} imgUrl
 * @returns {string} 图片路径
 */
export function convertImgUrl(imgUrl) {
  if (!imgUrl) {
    return 'https://dummyimage.com/400x400'
  }
  // 网络资源
  if (imgUrl.startsWith('http')) {
    return imgUrl
  }
  // 后端服务器资源: /public/uploaded/xxx.jpg
  return `/${imgUrl.replace(/^\/+/, '')}`
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
 * 兜底文案: 一言接口不通时随机取一句, 而不是每次都显示同一句
 */
export const FALLBACK_SENTENCES = [
  '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。',
  '书山有路勤为径，学海无涯苦作舟。',
  '纸上得来终觉浅，绝知此事要躬行。',
  '不积跬步，无以至千里；不积小流，无以成江海。',
  '路漫漫其修远兮，吾将上下而求索。',
  '业不可不勤，勤则百弊自去。',
  '博观而约取，厚积而薄发。',
]

export function getRandomSentence() {
  return FALLBACK_SENTENCES[Math.floor(Math.random() * FALLBACK_SENTENCES.length)]
}

/**
 * 取一句"一言", 接口不通时从内置文案里随机取一句
 *
 * 首页的 HomeBanner 和 TalkingCarousel 都要用, 以前各自 fetch 一次, 同一个接口请求了两遍;
 * 这里缓存 Promise, 一次页面加载只请求一次。这个接口在部分网络下不通, 所以内部吃掉异常。
 * @returns {Promise<string>}
 */
let sentencePromise = null
export function getOneSentence() {
  if (!sentencePromise) {
    sentencePromise = fetch('https://v1.hitokoto.cn?c=i')
      .then(res => res.json())
      .then(data => data?.hitokoto || getRandomSentence())
      .catch(() => getRandomSentence())
  }
  return sentencePromise
}

export * from './http'
