import { describe, expect, it } from 'vitest'
import { convertImgUrl, stripMarkdown } from '@/utils'
import { hasMath } from '@/utils/mathjax'

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
    // 返回根相对路径, 由 dev proxy / nginx 转发, 不拼 VITE_BACKEND_URL:
    // 拼上 localhost 后从别的机器访问页面时图片会裂
    expect(convertImgUrl('public/uploaded/a.png')).toBe('/public/uploaded/a.png')
    expect(convertImgUrl('/public/uploaded/a.png')).toBe('/public/uploaded/a.png')
  })
})

describe('stripMarkdown', () => {
  it('空值返回空字符串', () => {
    expect(stripMarkdown('')).toBe('')
    expect(stripMarkdown(null)).toBe('')
    expect(stripMarkdown(undefined)).toBe('')
  })

  it('去掉标题、列表、引用记号', () => {
    expect(stripMarkdown('## 标题\n- 第一项\n- 第二项\n> 引用')).toBe('标题 第一项 第二项 引用')
    expect(stripMarkdown('1. 一\n2. 二')).toBe('一 二')
  })

  it('去掉代码块与分隔线, 保留行内代码内容', () => {
    expect(stripMarkdown('前\n```js\nconsole.log(1)\n```\n后')).toBe('前 后')
    expect(stripMarkdown('用 `pnpm dev` 启动')).toBe('用 pnpm dev 启动')
    expect(stripMarkdown('上\n\n---\n\n下')).toBe('上 下')
  })

  it('图片整体去掉, 链接只保留文字', () => {
    expect(stripMarkdown('![封面](/a.png) 正文')).toBe('正文')
    expect(stripMarkdown('见 [文档](https://a.com/doc)')).toBe('见 文档')
  })

  it('去掉强调记号和内联 HTML', () => {
    expect(stripMarkdown('**加粗** *斜体* ~~删除~~')).toBe('加粗 斜体 删除')
    expect(stripMarkdown('<div class="x">内容</div>')).toBe('内容')
  })

  it('折叠多余空白', () => {
    expect(stripMarkdown('  一段\n\n\n文字   ')).toBe('一段 文字')
  })
})

describe('hasMath', () => {
  it('没有公式时返回 false', () => {
    expect(hasMath('')).toBe(false)
    expect(hasMath(null)).toBe(false)
    expect(hasMath('价格是 $100')).toBe(false)
    expect(hasMath('普通正文')).toBe(false)
  })

  it('两个 $ 之间的普通文本会被当成公式', () => {
    // 与 MathJax 自身的行内分隔符配置一致, 最多多加载一次脚本, 不影响正文渲染
    expect(hasMath('价格是 $100 到 $200')).toBe(true)
  })

  it('识别行内与块级公式', () => {
    expect(hasMath('这里有 $a+b$ 公式')).toBe(true)
    expect(hasMath('$$\n\\sum_{i=1}^n i\n$$')).toBe(true)
    expect(hasMath('用 \\(x^2\\) 表示')).toBe(true)
    expect(hasMath('\\[ x^2 \\]')).toBe(true)
  })
})
