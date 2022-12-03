import { request } from '@/utils'

export default {
  // TODO:
  refreshToken: () => request.post('/auth/refreshToken', null, { noNeedTip: true }),
  report: () => request.post('/report'), // 上报用户信息
}

export * from './home'
export * from './article'
export * from './category'
export * from './tag'
export * from './auth'
export * from './user'
export * from './comment'
export * from './message'
export * from './upload'
export * from './profile'
export * from './link'
export * from './setting'
export * from './log'
