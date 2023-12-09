export const config = {
  _APP_ID: '101878726',
  QQ_REDIRECT_URI: 'https://www.talkxj.com/oauth/login/qq',
  WEIBO_APP_ID: '4039197195',
  WEIBO_REDIRECT_URI: 'https://www.talkxj.com/oauth/login/weibo',
  TENCENT_CAPTCHA: '2088053498',
}

// 登录方式选项
export const loginTypeOptions = [
  { label: '邮箱', value: 1 },
  { label: 'QQ', value: 2 },
  { label: '微博', value: 3 },
]

export const loginTypeMap = {
  1: { name: '邮箱', tag: 'success' },
  2: { name: 'QQ', tag: 'info' },
  3: { name: '微博', tag: 'warning' },
}

// 文章类型选项
export const articleTypeOptions = [
  { label: '原创', value: 1 },
  { label: '转载', value: 2 },
  { label: '翻译', value: 3 },
]

export const articleTypeMap = {
  1: { name: '原创', tag: 'error' },
  2: { name: '转载', tag: 'success' },
  3: { name: '翻译', tag: 'warning' },
}

// 评论类型选项
export const commentTypeOptions = [
  { label: '文章', value: 1 },
  { label: '友链', value: 2 },
  { label: '说说', value: 3 },
]

export const commentTypeMap = {
  1: { name: '文章', tag: 'info' },
  2: { name: '友链', tag: 'warning' },
  3: { name: '说说', tag: 'error' },
}
