import fs from 'node:fs'
import { createRequire } from 'node:module'
import { defineConfig, presetIcons, presetUno } from 'unocss'

const require = createRequire(import.meta.url)

// presetIcons 默认的自动加载在当前 pnpm 布局下解析不到 @iconify/json,
// 结果是所有 i-xxx:yyy 类名图标静默不生成, 故显式指定图标集的读取方式。
function loadCollection(name) {
  return () => JSON.parse(
    fs.readFileSync(require.resolve(`@iconify/json/json/${name}.json`), 'utf-8'),
  )
}

// 模板中实际用到的图标集
const iconCollections = Object.fromEntries([
  'ant-design',
  'bxs',
  'fa6-solid',
  'fe',
  'heroicons',
  'ic',
  'ion',
  'line-md',
  'lucide',
  'majesticons',
  'material-symbols',
  'mdi',
  'mi',
  'mingcute',
  'uiw',
].map(name => [name, loadCollection(name)]))

export default defineConfig({
  presets: [
    presetUno(),
    presetIcons({
      prefix: ['i-'],
      scale: 1.2,
      // 加载失败时告警, 避免再次静默不生成
      warn: true,
      collections: iconCollections,
      extraProperties: {
        'display': 'inline-block',
        'vertical-align': 'middle',
      },
    }),
  ],
  theme: {
    colors: {
      primary: 'var(--primary-color)',
      primary_hover: 'var(--primary-color-hover)',
      primary_pressed: 'var(--primary-color-pressed)',
      primary_active: 'var(--primary-color-active)',
      info: 'var(--info-color)',
      info_hover: 'var(--info-color-hover)',
      info_pressed: 'var(--info-color-pressed)',
      info_active: 'var(--info-color-active)',
      success: 'var(--success-color)',
      success_hover: 'var(--success-color-hover)',
      success_pressed: 'var(--success-color-pressed)',
      success_active: 'var(--success-color-active)',
      warning: 'var(--warning-color)',
      warning_hover: 'var(--warning-color-hover)',
      warning_pressed: 'var(--warning-color-pressed)',
      warning_active: 'var(--warning-color-active)',
      error: 'var(--error-color)',
      error_hover: 'var(--error-color-hover)',
      error_pressed: 'var(--error-color-pressed)',
      error_active: 'var(--error-color-active)',
      dark: '#18181c',
    },
  },
})
