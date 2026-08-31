import antfu from '@antfu/eslint-config'

export default antfu({
  unocss: true,
  // public 下是原样拷贝到 dist 的静态资源, 依赖外部全局变量, 不参与 lint
  ignores: ['public'],
  rules: {
    'no-console': 'warn',
    'curly': 'off',
    // antfu v9 起 brace-style 由 @stylistic 提供, 旧的 @typescript-eslint/brace-style 已移除
    'style/brace-style': 'off',
    'unused-imports/no-unused-imports': 'off',
    'node/prefer-global/process': 'off',
  },
}, {
  // naive-ui 的 discrete api 由 setupNaiveDiscreteApi() 挂到 window 上
  // (见 src/utils/naiveTool.js), 模板中直接以裸变量形式使用
  languageOptions: {
    globals: {
      $message: 'readonly',
      $dialog: 'readonly',
      $notification: 'readonly',
      $loadingBar: 'readonly',
    },
  },
})
