window.MathJax = {
  tex: {
    // 行内公式选择符
    inlineMath: [
      ['$', '$'],
      ['\\(', '\\)'],
    ],
    // 段内公式选择符
    displayMath: [
      ['$$', '$$'],
      ['\\[', '\\]'],
    ],
  },
  startup: {
    ready() {
      MathJax.startup.defaultReady()
    },
  },
}
