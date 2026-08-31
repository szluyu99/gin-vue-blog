// 无后端 Mock 适配器: 通过替换 axios 的 adapter 拦截所有请求
// 开启方式: .env.development 中设置 VITE_USE_MOCK = true
import {
  about,
  articles,
  categories,
  comments,
  config,
  links,
  menus,
  messages,
  operationLogs,
  pages,
  resources,
  roles,
  tags,
  users,
} from './data.js'

const DELAY = 200 // 模拟网络延迟 (ms)

// 运行期可变状态, 刷新页面后重置
const state = {
  articles: articles.map(e => ({ ...e })),
  categories: categories.map(e => ({ ...e })),
  tags: tags.map(e => ({ ...e })),
  comments: comments.map(e => ({ ...e })),
  messages: messages.map(e => ({ ...e })),
  links: links.map(e => ({ ...e })),
  logs: operationLogs.map(e => ({ ...e })),
  users: users.map(e => ({ ...e })),
  menus: menus.map(e => ({ ...e })),
  resources: resources.map(e => ({ ...e })),
  roles: roles.map(e => ({ ...e })),
  pages: pages.map(e => ({ ...e })),
  config: { ...config },
  about,
  onlineUserIds: [1],
}

let autoId = 1000
function nextId() {
  return ++autoId
}

function ok(data = null) {
  return { code: 0, message: 'OK', data }
}

function fail(message, code = 1) {
  return { code, message, data: null }
}

function nowISO() {
  return new Date().toISOString()
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

/**
 * 按查询条件过滤: 字符串条件做包含匹配, 其他做相等匹配
 * @param {any[]} list 待过滤列表
 * @param {object} params 查询参数
 * @param {Record<string, (item: any) => any>} fields 支持的查询字段 => 取值函数
 */
function filterList(list, params, fields) {
  return Object.entries(fields).reduce((result, [key, getValue]) => {
    const expected = params[key]
    if (expected === undefined || expected === null || expected === '') {
      return result
    }
    return result.filter((item) => {
      const actual = getValue(item)
      if (typeof expected === 'string' && typeof actual === 'string') {
        return actual.includes(expected)
      }
      if (typeof actual === 'boolean') {
        return actual === (expected === 'true' || expected === true)
      }
      return String(actual) === String(expected)
    })
  }, list)
}

// 新增或修改: 有 id 则合并, 否则追加
function upsert(list, data, extra = {}) {
  if (data.id) {
    const target = list.find(e => e.id === Number(data.id))
    if (!target) {
      return fail('数据不存在')
    }
    Object.assign(target, data, { updated_at: nowISO() }, extra)
    return ok(target)
  }
  const created = { ...data, ...extra, id: nextId(), created_at: nowISO(), updated_at: nowISO() }
  list.unshift(created)
  return ok(created)
}

// 批量删除: body 为 id 数组
function removeByIds(list, ids = []) {
  const set = new Set((Array.isArray(ids) ? ids : [ids]).map(Number))
  for (let i = list.length - 1; i >= 0; i--) {
    if (set.has(list[i].id)) {
      list.splice(i, 1)
    }
  }
  return ok(set.size)
}

// 扁平数据 => 树形结构 (菜单 / 资源)
function buildTree(list, mapper) {
  return list
    .filter(e => e.parent_id === 0)
    .map(e => ({
      ...mapper(e),
      children: list.filter(c => c.parent_id === e.id).map(c => ({ ...mapper(c), children: null })),
    }))
    .map(e => ({ ...e, children: e.children.length ? e.children : null }))
}

function toMenu(m) {
  return { ...m, created_at: '2024-01-01T10:00:00.000Z', updated_at: '2024-01-01T10:00:00.000Z', roles: null }
}

function toResource(r) {
  const { parent_id, ...rest } = r
  return { ...rest, created_at: '2024-01-01T10:00:00.000Z' }
}

// 树形选项 (菜单/资源选择器)
function toOptionTree(list) {
  return buildTree(list, e => ({ key: e.id, label: e.name }))
    .map(e => ({ key: e.key, label: e.label, children: e.children?.map(c => ({ key: c.key, label: c.label, children: null })) ?? null }))
}

function findTags(ids = []) {
  return ids.map(id => state.tags.find(t => t.id === id)).filter(Boolean)
}

function findCategory(id) {
  const category = state.categories.find(c => c.id === id)
  return category ? { ...category, Articles: null } : null
}

// 文章列表项: 携带分类与标签
function toArticle(a) {
  const { tag_ids, ...rest } = a
  return { ...rest, tags: findTags(tag_ids), category: findCategory(a.category_id), user: null }
}

// 评论列表项: 后台展示需要昵称、头像、文章标题
function toComment(c) {
  const user = state.users.find(u => u.id === c.user_id)
  const replyUser = state.users.find(u => u.id === c.reply_user_id)
  const article = state.articles.find(a => a.id === c.topic_id)
  return {
    ...c,
    nickname: user?.info?.nickname ?? '匿名',
    avatar: user?.info?.avatar ?? '',
    reply_nickname: replyUser?.info?.nickname ?? '',
    article_title: article?.title ?? '',
  }
}

// 用户列表项: 携带角色
function toUser(u) {
  const { role_ids, ...rest } = u
  return { ...rest, roles: role_ids.map(id => state.roles.find(r => r.id === id)).filter(Boolean) }
}

// 路由表: [method, 路径正则, 处理函数(params, body, 正则捕获组)]
const handlers = [
  ['POST', /^\/login$/, () => ok({ token: 'mock-token' })],
  ['GET', /^\/logout$/, () => ok()],
  ['POST', /^\/report$/, () => ok()],
  ['GET', /^\/home$/, () => ok({
    article_count: state.articles.length,
    user_count: state.users.length,
    message_count: state.messages.length,
    view_count: state.articles.reduce((sum, e) => sum + e.view_count, 0),
  })],

  // 文章
  ['GET', /^\/article\/list$/, (params) => {
    const list = filterList(state.articles, params, {
      title: e => e.title,
      type: e => e.type,
      status: e => e.status,
      category_id: e => e.category_id,
      is_delete: e => e.is_delete,
      tag_id: e => e.tag_ids.join(','),
    })
    const page = paginate(list, params)
    return ok({ ...page, page_data: page.page_data.map(toArticle) })
  }],
  ['POST', /^\/article\/export$/, () => ok([])],
  ['POST', /^\/article\/import$/, () => ok()],
  ['PUT', /^\/article\/soft-delete$/, (params, body) => {
    state.articles.filter(e => body.ids?.includes(e.id)).forEach(e => e.is_delete = body.is_delete)
    return ok()
  }],
  ['PUT', /^\/article\/top$/, (params, body) => {
    const article = state.articles.find(e => e.id === Number(body.id))
    if (!article) {
      return fail('文章不存在')
    }
    article.is_top = body.is_top
    return ok()
  }],
  ['POST', /^\/article$/, (params, body) => {
    // 分类与标签用名称提交, 不存在则自动创建
    const { category_name, tag_names = [], ...rest } = body
    let category = state.categories.find(e => e.name === category_name)
    if (category_name && !category) {
      category = { id: nextId(), name: category_name, created_at: nowISO(), updated_at: nowISO() }
      state.categories.unshift(category)
    }
    const tag_ids = tag_names.map((name) => {
      let tag = state.tags.find(e => e.name === name)
      if (!tag) {
        tag = { id: nextId(), name, created_at: nowISO(), updated_at: nowISO() }
        state.tags.unshift(tag)
      }
      return tag.id
    })
    return upsert(state.articles, rest, {
      category_id: category?.id ?? 0,
      tag_ids,
      img: rest.img || '',
      like_count: rest.like_count ?? 0,
      view_count: rest.view_count ?? 0,
      comment_count: rest.comment_count ?? 0,
    })
  }],
  ['DELETE', /^\/article$/, (params, body) => removeByIds(state.articles, body)],
  ['GET', /^\/article\/(\d+)$/, (params, body, [id]) => {
    const article = state.articles.find(e => e.id === Number(id))
    return article ? ok(toArticle(article)) : fail('文章不存在')
  }],

  // 分类
  ['GET', /^\/category\/list$/, (params) => {
    const list = filterList(state.categories, params, { name: e => e.name })
      .map(e => ({ ...e, Articles: null, article_count: state.articles.filter(a => a.category_id === e.id).length }))
    return ok(paginate(list, params))
  }],
  ['POST', /^\/category$/, (params, body) => upsert(state.categories, body)],
  ['DELETE', /^\/category$/, (params, body) => removeByIds(state.categories, body)],
  ['GET', /^\/category\/option$/, () => ok(state.categories.map(e => ({ value: e.id, label: e.name })))],

  // 标签
  ['GET', /^\/tag\/list$/, (params) => {
    const list = filterList(state.tags, params, { name: e => e.name })
      .map(e => ({ ...e, article_count: state.articles.filter(a => a.tag_ids.includes(e.id)).length }))
    return ok(paginate(list, params))
  }],
  ['POST', /^\/tag$/, (params, body) => upsert(state.tags, body)],
  ['DELETE', /^\/tag$/, (params, body) => removeByIds(state.tags, body)],
  ['GET', /^\/tag\/option$/, () => ok(state.tags.map(e => ({ value: e.id, label: e.name })))],

  // 留言
  ['GET', /^\/message\/list$/, (params) => {
    const list = filterList(state.messages, params, { nickname: e => e.nickname, is_review: e => e.is_review })
    return ok(paginate(list, params))
  }],
  ['DELETE', /^\/message$/, (params, body) => removeByIds(state.messages, body)],
  ['PUT', /^\/message\/review$/, (params, body) => {
    state.messages.filter(e => body.ids?.includes(e.id)).forEach(e => e.is_review = body.is_review)
    return ok()
  }],

  // 评论
  ['GET', /^\/comment\/list$/, (params) => {
    const list = filterList(state.comments.map(toComment), params, {
      nickname: e => e.nickname,
      type: e => e.type,
      is_review: e => e.is_review,
    })
    return ok(paginate(list, params))
  }],
  ['DELETE', /^\/comment$/, (params, body) => removeByIds(state.comments, body)],
  ['PUT', /^\/comment\/review$/, (params, body) => {
    state.comments.filter(e => body.ids?.includes(e.id)).forEach(e => e.is_review = body.is_review)
    return ok()
  }],

  // 友链
  ['GET', /^\/link\/list$/, params => ok(paginate(filterList(state.links, params, { name: e => e.name }), params))],
  ['POST', /^\/link$/, (params, body) => upsert(state.links, body)],
  ['DELETE', /^\/link$/, (params, body) => removeByIds(state.links, body)],

  // 操作日志
  ['GET', /^\/operation\/log\/list$/, (params) => {
    const keyword = params.keyword
    const list = keyword
      ? state.logs.filter(e => `${e.opt_module}${e.opt_desc}${e.nickname}`.includes(keyword))
      : state.logs
    return ok(paginate(list, params))
  }],
  ['DELETE', /^\/operation\/log$/, (params, body) => removeByIds(state.logs, body)],

  // 用户
  ['GET', /^\/user\/info$/, () => {
    const user = state.users[0]
    return ok({ ...user.info, article_like_set: [], comment_like_set: [] })
  }],
  ['PUT', /^\/user\/current\/password$/, () => ok()],
  ['PUT', /^\/user\/current$/, (params, body) => {
    Object.assign(state.users[0].info, body)
    return ok()
  }],
  ['GET', /^\/user\/list$/, (params) => {
    const list = filterList(state.users, params, {
      username: e => e.username,
      nickname: e => e.info.nickname,
      login_type: e => e.login_type,
    })
    return ok({ ...paginate(list, params), page_data: paginate(list, params).page_data.map(toUser) })
  }],
  ['GET', /^\/user\/online$/, (params) => {
    const keyword = params.keyword
    const list = state.users
      .filter(e => state.onlineUserIds.includes(e.id))
      .filter(e => !keyword || e.info.nickname.includes(keyword) || e.username.includes(keyword))
    return ok(list.map(toUser))
  }],
  ['POST', /^\/user\/offline\/(\d+)$/, (params, body, [id]) => {
    state.onlineUserIds = state.onlineUserIds.filter(e => e !== Number(id))
    return ok()
  }],
  ['PUT', /^\/user\/disable$/, (params, body) => {
    const user = state.users.find(e => e.id === Number(body.id))
    if (!user) {
      return fail('用户不存在')
    }
    user.is_disable = body.is_disable
    return ok()
  }],
  ['PUT', /^\/user$/, (params, body) => {
    const user = state.users.find(e => e.id === Number(body.id))
    if (!user) {
      return fail('用户不存在')
    }
    Object.assign(user, { role_ids: body.role_ids ?? user.role_ids })
    Object.assign(user.info, { nickname: body.nickname ?? user.info.nickname })
    return ok()
  }],

  // 网站设置
  ['GET', /^\/config$/, () => ok(state.config)],
  ['PATCH', /^\/config$/, (params, body) => {
    Object.assign(state.config, body)
    return ok()
  }],
  ['GET', /^\/setting\/about$/, () => ok(state.about)],
  ['PUT', /^\/setting\/about$/, (params, body) => {
    state.about = body.content ?? body.about ?? state.about
    return ok()
  }],

  // 菜单
  ['GET', /^\/menu\/user\/list$/, () => ok(buildTree(state.menus, toMenu))],
  ['GET', /^\/menu\/list$/, (params) => {
    const tree = buildTree(state.menus, toMenu)
    return ok(params.keyword ? tree.filter(e => e.name.includes(params.keyword)) : tree)
  }],
  ['GET', /^\/menu\/option$/, () => ok(toOptionTree(state.menus))],
  ['POST', /^\/menu$/, (params, body) => upsert(state.menus, body, { parent_id: body.parent_id ?? 0 })],
  ['DELETE', /^\/menu\/(\d+)$/, (params, body, [id]) => removeByIds(state.menus, [Number(id)])],

  // 接口资源
  ['GET', /^\/resource\/list$/, (params) => {
    const tree = buildTree(state.resources, toResource)
    return ok(params.keyword ? tree.filter(e => e.name.includes(params.keyword)) : tree)
  }],
  ['GET', /^\/resource\/option$/, () => ok(toOptionTree(state.resources))],
  ['POST', /^\/resource$/, (params, body) => upsert(state.resources, body, { parent_id: body.parent_id ?? 0 })],
  ['PUT', /^\/resource\/anonymous$/, (params, body) => {
    const resource = state.resources.find(e => e.id === Number(body.id))
    if (!resource) {
      return fail('资源不存在')
    }
    resource.is_anonymous = body.is_anonymous
    return ok()
  }],
  ['DELETE', /^\/resource\/(\d+)$/, (params, body, [id]) => removeByIds(state.resources, [Number(id)])],

  // 角色
  ['GET', /^\/role\/list$/, (params) => {
    const list = filterList(state.roles, params, { name: e => e.name, label: e => e.label })
    return ok(paginate(list, params))
  }],
  ['GET', /^\/role\/option$/, () => ok(state.roles.map(e => ({ value: e.id, label: e.name })))],
  ['POST', /^\/role$/, (params, body) => upsert(state.roles, body, {
    resource_ids: body.resource_ids ?? [],
    menu_ids: body.menu_ids ?? [],
  })],
  ['DELETE', /^\/role$/, (params, body) => removeByIds(state.roles, body)],

  // 页面
  ['GET', /^\/page\/list$/, () => ok(state.pages)],
  ['POST', /^\/page$/, (params, body) => upsert(state.pages, body)],
  ['DELETE', /^\/page$/, (params, body) => removeByIds(state.pages, body)],
]

/**
 * axios 适配器: 命中 mock 路由则返回假数据, 未命中返回错误
 * @param {import('axios').InternalAxiosRequestConfig} requestConfig
 */
export function mockAdapter(requestConfig) {
  const method = (requestConfig.method ?? 'get').toUpperCase()
  // baseURL 为 /api, 去掉前缀后与路由表匹配
  const path = `${requestConfig.baseURL ?? ''}${requestConfig.url ?? ''}`.replace(/^\/api/, '').split('?')[0]

  let body = requestConfig.data
  if (typeof body === 'string') {
    try {
      body = JSON.parse(body)
    }
    catch {
      body = {}
    }
  }

  const matched = handlers.find(([m, re]) => m === method && re.test(path))
  if (!matched) {
    console.warn(`[mock] 未定义的接口: ${method} ${path}`)
  }
  const data = matched
    ? matched[2](requestConfig.params ?? {}, body ?? {}, path.match(matched[1]).slice(1))
    : fail(`[mock] 未定义的接口: ${method} ${path}`)

  return new Promise(resolve => setTimeout(resolve, DELAY, {
    data,
    status: 200,
    statusText: 'OK',
    headers: {},
    config: requestConfig,
  }))
}
