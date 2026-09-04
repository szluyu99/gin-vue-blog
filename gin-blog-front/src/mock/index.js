// 无后端 Mock 适配器: 通过替换 axios 的 adapter 拦截所有请求
// 开启方式: .env.development 中设置 VITE_USE_MOCK = true
import {
  about,
  articles,
  blogConfig,
  categories,
  comments,
  commentUsers,
  currentUser,
  links,
  messages,
  pages,
  tags,
} from './data.js'

const DELAY = 200 // 模拟网络延迟 (ms)

// 运行期可变状态, 刷新页面后重置
const state = {
  articles: articles.map(e => ({ ...e })),
  comments: comments.map(e => ({ ...e })),
  messages: messages.map(e => ({ ...e })),
  user: { ...currentUser },
}

function ok(data) {
  return { code: 0, message: 'OK', data }
}

function fail(message, code = 1) {
  return { code, message, data: null }
}

function findTags(ids = []) {
  return ids.map(id => tags.find(t => t.id === id)).filter(Boolean)
}

function findCategory(id) {
  const category = categories.find(c => c.id === id)
  return category ? { ...category, Articles: null } : null
}

// 列表项: 携带分类与标签
function toListItem(a) {
  const { tag_ids, like_count, view_count, ...rest } = a
  return { ...rest, tags: findTags(tag_ids), category: findCategory(a.category_id), user: null }
}

// 简要信息: 上一篇/下一篇/推荐/最新
function toBrief(a, withDate = false) {
  return a
    ? (withDate ? { id: a.id, img: a.img, title: a.title, created_at: a.created_at } : { id: a.id, img: a.img, title: a.title })
    : { id: 0, img: '', title: '' }
}

// 评论: 拼装成后端返回的结构 (user.info)
function toComment(c) {
  const replies = state.comments.filter(e => e.parent_id === c.id)
  return {
    id: c.id,
    created_at: c.created_at,
    content: c.content,
    like_count: c.like_count,
    user_id: c.user_id,
    reply_count: replies.length,
    reply_list: replies.slice(0, 5).map(toReply),
    user: { info: commentUsers[c.user_id] ?? commentUsers[2] },
  }
}

function toReply(c) {
  return {
    id: c.id,
    created_at: c.created_at,
    content: c.content,
    like_count: c.like_count,
    user_id: c.user_id,
    reply_user_id: c.reply_user_id,
    user: { info: commentUsers[c.user_id] ?? commentUsers[2] },
    reply_user: { info: commentUsers[c.reply_user_id] ?? commentUsers[2] },
  }
}

function paginate(list, { page_num = 1, page_size = 10 } = {}) {
  const num = Number(page_num) || 1
  const size = Number(page_size) || 10
  return {
    page_num: num,
    page_size: size,
    total: list.length,
    page_data: list.slice((num - 1) * size, num * size),
  }
}

// 搜索: 与后端一致, 命中关键字用 span 高亮
function highlight(text, keyword) {
  return text.replaceAll(keyword, `<span style='color:#f47466'>${keyword}</span>`)
}

// 点赞/取消点赞: 同时维护当前用户的点赞集合与目标的点赞数
function toggleLike(set, target) {
  const idx = set.indexOf(target.id)
  if (idx > -1) {
    set.splice(idx, 1)
    target.like_count--
  }
  else {
    set.push(target.id)
    target.like_count++
  }
}

// 路由表: [method, 路径正则, 处理函数(params, body, 正则捕获组)]
const handlers = [
  ['POST', /^\/login$/, () => ok({ token: 'mock-token', ...state.user })],
  ['POST', /^\/register$/, () => ok(null)],
  ['GET', /^\/logout$/, () => ok(null)],
  ['GET', /^\/code$/, () => ok(null)],

  ['GET', /^\/front\/about$/, () => ok(about)],
  ['GET', /^\/front\/page$/, () => ok(pages)],
  ['GET', /^\/front\/home$/, () => ok({
    article_count: state.articles.length,
    user_count: Object.keys(commentUsers).length,
    message_count: state.messages.length,
    category_count: categories.length,
    tag_count: tags.length,
    view_count: state.articles.reduce((sum, e) => sum + e.view_count, 0),
    blog_config: blogConfig,
  })],

  ['GET', /^\/front\/article\/list$/, (params) => {
    let list = state.articles
    if (params.category_id) {
      list = list.filter(e => e.category_id === Number(params.category_id))
    }
    if (params.tag_id) {
      list = list.filter(e => e.tag_ids.includes(Number(params.tag_id)))
    }
    // 与后端一致: 返回分页结构, 前台靠 total 渲染分页器
    const page = paginate(list, params)
    return ok({ ...page, page_data: page.page_data.map(toListItem) })
  }],

  ['GET', /^\/front\/article\/archive$/, (params) => {
    const list = state.articles.map(e => ({ id: e.id, title: e.title, created_at: e.created_at }))
    return ok(paginate(list, params))
  }],

  ['GET', /^\/front\/article\/search$/, (params) => {
    const keyword = (params.keyword ?? '').trim()
    if (!keyword) {
      return ok([])
    }
    return ok(state.articles
      .filter(e => e.title.includes(keyword) || e.content.includes(keyword))
      .map(e => ({
        id: e.id,
        title: highlight(e.title, keyword),
        content: highlight(e.content, keyword),
      })))
  }],

  ['GET', /^\/front\/article\/like\/(\d+)$/, (params, body, [id]) => {
    const article = state.articles.find(e => e.id === Number(id))
    if (!article) {
      return fail('文章不存在')
    }
    const set = state.user.article_like_set
    toggleLike(set, article)
    return ok(null)
  }],

  ['GET', /^\/front\/article\/(\d+)$/, (params, body, [id]) => {
    const idx = state.articles.findIndex(e => e.id === Number(id))
    if (idx === -1) {
      return fail('文章不存在')
    }
    const article = state.articles[idx]
    article.view_count++
    return ok({
      ...toListItem(article),
      comment_count: state.comments.filter(e => e.type === 1 && e.topic_id === article.id).length,
      like_count: article.like_count,
      view_count: article.view_count,
      // 列表按时间倒序, 所以下一篇在前, 上一篇在后
      last_article: toBrief(state.articles[idx + 1]),
      next_article: toBrief(state.articles[idx - 1]),
      recommend_articles: state.articles.filter(e => e.id !== article.id).slice(0, 3).map(e => toBrief(e, true)),
      newest_articles: state.articles.slice(0, 5).map(e => toBrief(e, true)),
    })
  }],

  ['GET', /^\/front\/category\/list$/, () => ok(categories.map(c => ({
    ...c,
    Articles: null,
    article_count: state.articles.filter(a => a.category_id === c.id).length,
  })))],

  ['GET', /^\/front\/tag\/list$/, () => ok(tags.map(t => ({
    ...t,
    article_count: state.articles.filter(a => a.tag_ids.includes(t.id)).length,
  })))],

  ['GET', /^\/front\/message\/list$/, () => ok(state.messages)],
  ['POST', /^\/front\/message$/, (params, body) => {
    state.messages.push({ id: Date.now(), created_at: new Date().toISOString(), ...body })
    return ok(null)
  }],

  ['GET', /^\/front\/link\/list$/, () => ok(links)],

  ['GET', /^\/front\/comment\/list$/, (params) => {
    const type = Number(params.type ?? 1)
    const topicId = Number(params.topic_id ?? 0)
    const list = state.comments
      .filter(e => e.parent_id === 0 && e.type === type && e.topic_id === topicId)
      .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    const page = paginate(list, params)
    return ok({ ...page, page_data: page.page_data.map(toComment) })
  }],

  ['GET', /^\/front\/comment\/replies\/(\d+)$/, (params, body, [id]) => {
    const list = state.comments.filter(e => e.parent_id === Number(id))
    return ok(paginate(list, params).page_data.map(toReply))
  }],

  ['POST', /^\/front\/comment$/, (params, body) => {
    state.comments.push({
      id: Date.now(),
      type: Number(body.type ?? 1),
      topic_id: Number(body.topic_id ?? 0),
      parent_id: Number(body.parent_id ?? 0),
      user_id: state.user.id,
      reply_user_id: Number(body.reply_user_id ?? 0),
      content: body.content,
      like_count: 0,
      created_at: new Date().toISOString(),
    })
    return ok(null)
  }],

  ['GET', /^\/front\/comment\/like\/(\d+)$/, (params, body, [id]) => {
    const comment = state.comments.find(e => e.id === Number(id))
    if (!comment) {
      return fail('评论不存在')
    }
    const set = state.user.comment_like_set
    toggleLike(set, comment)
    return ok(null)
  }],

  ['GET', /^\/front\/user\/info$/, () => ok(state.user)],
  ['PUT', /^\/front\/user\/info$/, (params, body) => {
    Object.assign(state.user, body)
    return ok(null)
  }],
]

/**
 * axios 适配器: 命中 mock 路由则返回假数据, 未命中返回 404
 * @param {import('axios').InternalAxiosRequestConfig} config
 */
export function mockAdapter(config) {
  const method = (config.method ?? 'get').toUpperCase()
  // baseURL 为 /api 或 /api/front, 去掉 /api 前缀后与路由表匹配
  const path = `${config.baseURL ?? ''}${config.url ?? ''}`.replace(/^\/api/, '').split('?')[0]

  let body = config.data
  if (typeof body === 'string') {
    try {
      body = JSON.parse(body)
    }
    catch {
      body = {}
    }
  }

  const matched = handlers.find(([m, re]) => m === method && re.test(path))
  const data = matched
    ? matched[2](config.params ?? {}, body ?? {}, path.match(matched[1]).slice(1))
    : fail(`[mock] 未定义的接口: ${method} ${path}`)

  if (!matched) {
    console.warn(`[mock] 未定义的接口: ${method} ${path}`)
  }

  return new Promise(resolve => setTimeout(resolve, DELAY, {
    data,
    status: 200,
    statusText: 'OK',
    headers: {},
    config,
  }))
}
