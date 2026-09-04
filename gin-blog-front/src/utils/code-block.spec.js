import { beforeEach, describe, expect, it, vi } from 'vitest'

import { addCopyButtons } from './code-block'

function makeArticle(html) {
  const el = document.createElement('div')
  el.innerHTML = html
  document.body.appendChild(el)
  return el
}

const CODE_HTML = '<pre><code>go build ./...</code></pre><pre><code>pnpm dev</code></pre>'

describe('代码块复制按钮', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    // jsdom 默认没有 clipboard
    Object.defineProperty(navigator, 'clipboard', {
      value: { writeText: vi.fn().mockResolvedValue(undefined) },
      configurable: true,
      writable: true,
    })
  })

  it('每个代码块加一个按钮', () => {
    const root = makeArticle(CODE_HTML)

    expect(addCopyButtons(root)).toBe(2)
    expect(root.querySelectorAll('[data-copy-btn]')).toHaveLength(2)
  })

  // 样式必须是内联的: UnoCSS 默认不扫描 .js, 写类名不会生成 CSS
  it('按钮样式与定位走内联样式', () => {
    const root = makeArticle(CODE_HTML)
    addCopyButtons(root)

    const btn = root.querySelector('[data-copy-btn]')
    expect(btn.style.position).toBe('absolute')
    expect(btn.style.opacity).toBe('0')
    expect(btn.className).toBe('')
    // 按钮要定位在代码块内, pre 得是定位父级
    expect(root.querySelector('pre').style.position).toBe('relative')
  })

  it('鼠标移入代码块时按钮才显形', () => {
    const root = makeArticle(CODE_HTML)
    addCopyButtons(root)
    const pre = root.querySelector('pre')
    const btn = pre.querySelector('[data-copy-btn]')

    pre.dispatchEvent(new MouseEvent('mouseenter'))
    expect(btn.style.opacity).toBe('1')
    pre.dispatchEvent(new MouseEvent('mouseleave'))
    expect(btn.style.opacity).toBe('0')
  })

  // 正文重新渲染时会重复调用
  it('重复调用不会加出第二个按钮', () => {
    const root = makeArticle(CODE_HTML)

    addCopyButtons(root)
    expect(addCopyButtons(root)).toBe(0)
    expect(root.querySelectorAll('[data-copy-btn]')).toHaveLength(2)
  })

  it('点击后把代码写进剪贴板并回显状态', async () => {
    const root = makeArticle(CODE_HTML)
    addCopyButtons(root)

    const btn = root.querySelector('[data-copy-btn]')
    btn.click()
    await vi.waitFor(() => expect(btn.textContent).toBe('已复制'))

    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('go build ./...')
  })

  it('复制失败时提示失败而不是假装成功', async () => {
    navigator.clipboard.writeText = vi.fn().mockRejectedValue(new Error('denied'))
    const root = makeArticle(CODE_HTML)
    addCopyButtons(root)

    root.querySelector('[data-copy-btn]').click()
    await vi.waitFor(() => expect(root.querySelector('[data-copy-btn]').textContent).toBe('复制失败'))
  })

  // 非 https 或旧浏览器下没有剪贴板 API, 加了按钮点了也没反应
  it('没有剪贴板 API 时不加按钮', () => {
    Object.defineProperty(navigator, 'clipboard', { value: undefined, configurable: true, writable: true })
    const root = makeArticle(CODE_HTML)

    expect(addCopyButtons(root)).toBe(0)
    expect(root.querySelectorAll('[data-copy-btn]')).toHaveLength(0)
  })

  it('容器为 null 或没有代码块时返回 0', () => {
    expect(addCopyButtons(null)).toBe(0)
    expect(addCopyButtons(makeArticle('<p>只有正文</p>'))).toBe(0)
    // pre 里没有 code 也跳过
    expect(addCopyButtons(makeArticle('<pre>裸 pre</pre>'))).toBe(0)
  })
})
