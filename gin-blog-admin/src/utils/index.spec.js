import { describe, expect, it } from 'vitest'
import { convertImgUrl, formatDate, formatJson, IMG_PLACEHOLDER, parseJson } from '@/utils'

describe('convertImgUrl', () => {
  // 占位图从 dummyimage.com 换成内联 SVG: 图片本身挂了之后, 兜底不该再依赖一次外网
  it('空值返回内联占位图, 不请求外网', () => {
    expect(convertImgUrl('')).toBe(IMG_PLACEHOLDER)
    expect(convertImgUrl(undefined)).toBe(IMG_PLACEHOLDER)
    expect(convertImgUrl(null)).toBe(IMG_PLACEHOLDER)
    expect(IMG_PLACEHOLDER.startsWith('data:image/svg+xml')).toBe(true)
  })

  it('网络地址原样返回', () => {
    expect(convertImgUrl('https://cdn.test/a.png')).toBe('https://cdn.test/a.png')
    expect(convertImgUrl('http://cdn.test/a.png')).toBe('http://cdn.test/a.png')
  })

  it('相对地址拼上服务器地址', () => {
    // 返回根相对路径, 由 dev proxy / nginx 转发, 不拼 VITE_SERVER_URL:
    // 拼上 localhost 后从别的机器访问页面时图片会裂
    expect(convertImgUrl('article/a.png')).toBe('/article/a.png')
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

describe('parseJson', () => {
  it('正常 JSON 正常解析', () => {
    expect(parseJson('{"a":1}')).toEqual({ a: 1 })
    expect(parseJson('[1,2]')).toEqual([1, 2])
  })

  it('空串和非字符串返回 fallback', () => {
    // 操作日志里 DELETE /menu/:id 这类请求的 request_param 就是空串
    expect(parseJson('')).toBeNull()
    expect(parseJson('   ')).toBeNull()
    expect(parseJson(undefined)).toBeNull()
    expect(parseJson(null, 'x')).toBe('x')
  })

  it('非法 JSON 不抛异常', () => {
    // 上传接口被网关拦截时返回的是 HTML
    expect(parseJson('<html>413</html>')).toBeNull()
    expect(parseJson('{a:1}', {})).toEqual({})
  })
})

describe('formatJson', () => {
  it('格式化合法 JSON', () => {
    expect(formatJson('{"a":1}')).toBe('{\n  "a": 1\n}')
  })

  it('无法解析时原样返回, 不抛异常', () => {
    expect(formatJson('')).toBe('')
    expect(formatJson(undefined)).toBe('')
    expect(formatJson('not json')).toBe('not json')
  })
})
