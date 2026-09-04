/**
 * 给正文里的代码块加「复制」按钮
 *
 * 高亮由 highlight.js 处理, 它只管上色, 复制得自己加。
 * 做成工具函数而不是组件: 正文是 v-html 渲染出来的, 里面塞不进 Vue 组件。
 *
 * @param {ParentNode | null} root 正文容器, 传 null 直接返回 0
 * @returns {number} 这次加上按钮的代码块数量
 */
export function addCopyButtons(root) {
  // 没有剪贴板 API(旧浏览器 / 非 https)时不加按钮, 免得点了没反应
  if (!root || !navigator.clipboard) {
    return 0
  }

  let count = 0
  root.querySelectorAll('pre').forEach((pre) => {
    // 正文重新渲染时会重复调用, 已经加过的跳过
    if (pre.querySelector('[data-copy-btn]')) {
      return
    }
    const code = pre.querySelector('code')
    if (!code) {
      return
    }

    // 按钮绝对定位在代码块右上角, pre 需要是定位父级
    pre.style.position = 'relative'

    const btn = document.createElement('button')
    btn.dataset.copyBtn = 'true'
    btn.type = 'button'
    btn.textContent = '复制'
    // 用内联样式而不是 UnoCSS 类: UnoCSS 默认只扫描 .vue/.ts 等模板文件, 不扫 .js,
    // 写在这里的类名不会生成 CSS(实测过), 按钮会变成没定位没样式的裸按钮
    btn.style.cssText = [
      'position:absolute',
      'right:8px',
      'top:8px',
      'padding:2px 8px',
      'border-radius:4px',
      'font-size:12px',
      'line-height:1.5',
      'color:#fff',
      'background:rgba(0,0,0,.35)',
      'opacity:0',
      'transition:opacity .2s',
    ].join(';')
    // 鼠标移到代码块上才显形, 不挡代码
    pre.addEventListener('mouseenter', () => (btn.style.opacity = '1'))
    pre.addEventListener('mouseleave', () => (btn.style.opacity = '0'))

    btn.addEventListener('click', async () => {
      try {
        await navigator.clipboard.writeText(code.textContent ?? '')
        btn.textContent = '已复制'
      }
      catch (err) {
        console.error(err)
        btn.textContent = '复制失败'
      }
      setTimeout(() => (btn.textContent = '复制'), 1500)
    })

    pre.appendChild(btn)
    count++
  })

  return count
}
