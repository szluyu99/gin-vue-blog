import { createHtmlPlugin } from 'vite-plugin-html'

export function configHtmlPlugin(viteEnv, isBuild) {
  const { VITE_TITLE } = viteEnv

  // https://github.com/vbenjs/vite-plugin-html
  return createHtmlPlugin({
    minify: isBuild,
    // entry: 'src/main.js', // index.html 中无需添加 script 标签
    inject: {
      // 需要注入 index.html ejs 模版的数据
      data: {
        title: VITE_TITLE,
      },
    },
  })
}
