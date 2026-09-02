import { describe, expect, it } from 'vitest'
import { convertImgUrl } from '@/utils'

describe('convertImgUrl', () => {
  it('没有图片时返回占位图', () => {
    expect(convertImgUrl('')).toBe('http://dummyimage.com/400x400')
    expect(convertImgUrl(null)).toBe('http://dummyimage.com/400x400')
    expect(convertImgUrl(undefined)).toBe('http://dummyimage.com/400x400')
  })

  it('网络图片原样返回', () => {
    expect(convertImgUrl('http://a.com/1.png')).toBe('http://a.com/1.png')
    expect(convertImgUrl('https://a.com/1.png')).toBe('https://a.com/1.png')
  })

  it('相对路径拼接后端地址', () => {
    // VITE_BACKEND_URL 在 vitest.config.js 里固定为 http://test-server
    expect(convertImgUrl('public/uploaded/a.png')).toBe('http://test-server/public/uploaded/a.png')
  })
})
