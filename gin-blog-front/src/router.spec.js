import { describe, expect, it } from 'vitest'
import { router } from './router'

describe('路由滚动行为', () => {
  const scrollBehavior = to => router.options.scrollBehavior(to, {}, null)

  it('普通跳转回到页面顶部', () => {
    expect(scrollBehavior({ query: {} })).toEqual({ left: 0, top: 0 })
  })

  // 回归: 返回 top: 0 会把评论组件刚定位到的位置拉回顶部
  it('带 ?comment= 时不接管滚动, 交给评论组件定位', () => {
    expect(scrollBehavior({ query: { comment: '14' } })).toBe(false)
  })
})
