import { describe, expect, it } from 'vitest'
import { convertImgUrl, formatDate } from '@/utils'

describe('convertImgUrl', () => {
  it('空值返回占位图', () => {
    expect(convertImgUrl('')).toBe('http://dummyimage.com/400x400')
    expect(convertImgUrl(undefined)).toBe('http://dummyimage.com/400x400')
    expect(convertImgUrl(null)).toBe('http://dummyimage.com/400x400')
  })

  it('网络地址原样返回', () => {
    expect(convertImgUrl('https://cdn.test/a.png')).toBe('https://cdn.test/a.png')
    expect(convertImgUrl('http://cdn.test/a.png')).toBe('http://cdn.test/a.png')
  })

  it('相对地址拼上服务器地址', () => {
    // VITE_SERVER_URL 在 vitest.config.js 里固定为 http://test-server
    expect(convertImgUrl('article/a.png')).toBe('http://test-server/article/a.png')
  })
})

describe('formatDate', () => {
  it('默认格式为 YYYY-MM-DD', () => {
    expect(formatDate('2024-03-05T10:20:30Z')).toMatch(/^\d{4}-\d{2}-\d{2}$/)
  })

  it('支持自定义格式', () => {
    expect(formatDate('2024-03-05 10:20:30', 'YYYY/MM/DD HH:mm:ss')).toBe('2024/03/05 10:20:30')
  })
})
