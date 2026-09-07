import fs from 'node:fs'
import { createRequire } from 'node:module'
import {
  defineConfig,
  presetIcons,
  presetTypography,
  presetUno,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

const require = createRequire(import.meta.url)

// presetIcons 默认的自动加载在当前 pnpm 布局下解析不到 @iconify/json,
// 结果是所有 i-xxx:yyy 类名图标静默不生成, 故显式指定图标集的读取方式。
function loadCollection(name) {
  return () => JSON.parse(
    fs.readFileSync(require.resolve(`@iconify/json/json/${name}.json`), 'utf-8'),
  )
}

// 源码中实际用到的图标集
const iconCollections = Object.fromEntries([
  'akar-icons',
  'ant-design',
  'carbon',
  'ep',
  'fa',
  'fa-solid',
  'fluent-emoji-flat',
  'ic',
  'icon-park',
  'material-symbols',
  'mdi',
  'simple-icons',
  'uiw',
].map(name => [name, loadCollection(name)]))

export default defineConfig({
  shortcuts: [
    ['f-c-c', 'flex justify-center items-center'],
  ],
  // 语义色都指向 CSS 变量, 变量在 styles/index.css 里按 html.dark 切换,
  // 这样切主题只改根元素 class, 不用给每个类名都写一遍 dark: 变体
  theme: {
    colors: {
      'primary': '#49b1f5',
      // 强调色: 置顶标记、hover 高亮等处用, 原来散落在各处硬编码
      'accent': '#ff7242',
      'surface': 'var(--c-surface)',
      'surface-soft': 'var(--c-surface-soft)',
      'main': 'var(--c-text)',
      'muted': 'var(--c-text-muted)',
      'line': 'var(--c-border)',
      'divider': 'var(--c-divider)',
    },
  },
  presets: [
    presetUno(),
    /*
      extraProperties 里的 display 不能省: presetIcons 生成的规则只有 width/height +
      mask, 没有 display。图标 span 落在非 flex 的父元素里时仍是 inline,
      inline 元素不吃 width/height, 盒子塌成 0x0 —— 元素还在 DOM 里、计算样式也是 24px,
      但既看不见也点不到。移动端顶栏的主题/搜索/菜单三个按钮就是这么消失的
      (按钮里只有一个图标 span, button 本身不是 flex, 于是整个按钮 0x0)。
    */
    presetIcons({
      warn: true,
      collections: iconCollections,
      extraProperties: { display: 'inline-block' },
    }),
    presetTypography(),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
})
