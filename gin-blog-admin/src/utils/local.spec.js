import { beforeEach, describe, expect, it, vi } from 'vitest'
import { getLocal, removeLocal, setLocal } from '@/utils'

describe('local storage 工具', () => {
  beforeEach(() => {
    window.localStorage.clear()
    vi.useRealTimers()
  })

  it('存取一个对象', () => {
    setLocal('loginInfo', { username: 'admin', password: '123456' })

    expect(getLocal('loginInfo')).toEqual({ username: 'admin', password: '123456' })
    // 存的是加密后的内容, 不能明文出现
    expect(window.localStorage.getItem('loginInfo')).not.toContain('admin')
  })

  it('支持中文等非 Latin1 字符', () => {
    // btoa 直接吃中文会抛 InvalidCharacterError
    setLocal('loginInfo', { username: '管理员', password: '密码' })

    expect(getLocal('loginInfo')).toEqual({ username: '管理员', password: '密码' })
  })

  it('过期后取不到, 并且清掉存储', () => {
    setLocal('foo', 'bar', 1)
    vi.useFakeTimers()
    vi.setSystemTime(Date.now() + 2000)

    expect(getLocal('foo')).toBeNull()
    expect(window.localStorage.getItem('foo')).toBeNull()
  })

  it('expire 传 0 表示不过期', () => {
    setLocal('foo', 'bar', 0)
    vi.useFakeTimers()
    vi.setSystemTime(Date.now() + 10 * 365 * 24 * 3600 * 1000)

    expect(getLocal('foo')).toBe('bar')
  })

  it('数据被改坏时返回 null 而不是抛错', () => {
    window.localStorage.setItem('foo', 'not-a-valid-base64-@@@')

    expect(getLocal('foo')).toBeNull()
    expect(window.localStorage.getItem('foo')).toBeNull()
  })

  it('没有存过的 key 返回 null', () => {
    expect(getLocal('nothing')).toBeNull()
  })

  it('removeLocal 删除数据', () => {
    setLocal('foo', 'bar')
    removeLocal('foo')

    expect(getLocal('foo')).toBeNull()
  })
})
