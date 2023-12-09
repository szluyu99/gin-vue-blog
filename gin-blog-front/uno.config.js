import {
  defineConfig,
  presetIcons,
  presetTypography,
  presetUno,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

export default defineConfig({
  exclude: ['node_modules', '.git', '.github', '.husky', '.vscode', 'build', 'dist', 'mock', 'public', 'types', './stats.html'],
  shortcuts: [
    ['f-c-c', 'flex justify-center items-center'],
  ],
  presets: [
    presetUno(),
    presetIcons({ warn: true, autoInstall: true, cdn: 'https://esm.sh/' }),
    presetTypography(),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
})
