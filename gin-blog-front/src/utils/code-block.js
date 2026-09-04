/**
 * 把文本写进剪贴板
 *
 * navigator.clipboard 只在安全上下文(https / localhost)里存在, 用局域网 IP
 * 走 http 打开时是 undefined —— 联调时很常见, 所以退回 execCommand。
 * execCommand('copy') 已废弃但所有浏览器都还支持, 且不要求安全上下文。
 *
 * @param {string} text
 * @returns {Promise<void>} 复制失败时 reject
 */
async function writeClipboard(text) {
  if (navigator.clipboard) {
    try {
      return await navigator.clipboard.writeText(text)
    }
    catch (err) {
      // 安全上下文里也可能被权限策略拒掉(比如页面没聚焦), 别直接判死, 继续走兜底
      console.warn('clipboard API 写入失败, 退回 execCommand', err)
    }
  }

  const ta = document.createElement('textarea')
  ta.value = text
  // 不能用 display:none / visibility:hidden, 那样选不中也就复制不了
  ta.style.cssText = 'position:fixed;left:-9999px;top:0;opacity:0'
  ta.setAttribute('readonly', '')
  document.body.appendChild(ta)
  try {
    ta.select()
    ta.setSelectionRange(0, ta.value.length)
    if (!document.execCommand?.('copy')) {
      throw new Error('execCommand copy 失败')
    }
  }
  finally {
    ta.remove()
  }
}

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
  if (!root) {
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
      'opacity:0.45',
      'transition:opacity .2s',
    ].join(';')
    // 常态半透明: 之前默认全透明, 不知道有这功能的人永远不会把鼠标放上去
    pre.addEventListener('mouseenter', () => (btn.style.opacity = '1'))
    pre.addEventListener('mouseleave', () => (btn.style.opacity = '0.45'))

    btn.addEventListener('click', async () => {
      try {
        await writeClipboard(code.textContent ?? '')
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
