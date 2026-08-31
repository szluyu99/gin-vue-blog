import antfu from '@antfu/eslint-config'

export default antfu({
  unocss: true,
  rules: {
    'no-console': 'warn',
    'curly': 'off',
    'style/brace-style': 'off',
    'node/prefer-global/process': 'off',
    'unused-imports/no-unused-imports': 'off',
  },
})
