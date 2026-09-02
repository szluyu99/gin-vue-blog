/*
按需加载 MathJax

tex-mml-chtml.js 压缩后约 1MB, 以前在 index.html 里无条件引入, 连首页、归档、
友链这些没有公式的页面也要拉一遍, 走 jsdelivr 还经常很慢。改成正文里真的出现
公式时才动态插入 script。
公式的分隔符与 public/js/mathjax.js 里的 MathJax 配置保持一致。
*/

const MATH_PATTERN = /\$\$[\s\S]+?\$\$|\$[^$\n]+\$|\\\([\s\S]+?\\\)|\\\[[\s\S]+?\\\]/
const MATHJAX_URL = 'https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js'

let loadPromise = null

export function hasMath(text) {
  return !!text && MATH_PATTERN.test(text)
}

// 只插入一次, 后续调用复用同一个 Promise
function loadMathJax() {
  if (loadPromise) {
    return loadPromise
  }

  loadPromise = new Promise((resolve, reject) => {
    const script = document.createElement('script')
    // MathJax 靠这个 id 找到自己的 script 标签, 不能拼错
    script.id = 'MathJax-script'
    script.async = true
    script.src = MATHJAX_URL
    script.onload = () => resolve()
    script.onerror = () => {
      loadPromise = null // 允许下次重试
      reject(new Error('MathJax 加载失败'))
    }
    document.head.appendChild(script)
  })

  return loadPromise
}

// 正文里有公式才加载并渲染; 加载失败只打日志, 不影响文章正常显示
export async function typesetMath(text) {
  if (!hasMath(text)) {
    return
  }

  try {
    await loadMathJax()
    await window.MathJax?.startup?.promise
    window.MathJax?.typeset?.()
  }
  catch (err) {
    console.error(err)
  }
}
