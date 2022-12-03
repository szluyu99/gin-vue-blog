module.exports = {
  printWidth: 100, // 超过最大值换行
  tabWidth: 2, // 缩进字节数
  useTabs: false, // 使用空格缩进
  semi: false, // 句尾不添加分号
  singleQuote: true, // 单引号
  endOfLine: 'lf',
  quoteProps: 'as-needed', // 对象的 key 是否用引号括起来  "as-needed"|"consistent"|"preserve"
  bracketSpacing: true, // 是否在对象的{}内部两侧加空格 true - { foo: bar } false - {foo: bar}.
  // arrowParens: 'avoid', // 箭头函数只有一个参数时是否需要小括号
}
