import antfu from '@antfu/eslint-config'

export default antfu({
  unocss: true,
  // public 下是原样拷贝到 dist 的静态资源, 依赖 MathJax 等外部全局变量, 不参与 lint
  ignores: ['public'],
  rules: {
    'no-console': 'warn',
    'curly': 'off',
    'style/brace-style': 'off',
    'node/prefer-global/process': 'off',
    'unused-imports/no-unused-imports': 'off',
  },
})
