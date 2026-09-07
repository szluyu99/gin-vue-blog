import { createGenerator } from 'unocss'
import { describe, expect, it } from 'vitest'

import config from '~/uno.config'

/*
图标类必须带 display

presetIcons 生成的规则只有 width / height / mask, 没有 display。图标 span 落在
非 flex 的父元素里时仍是 inline, 而 inline 元素不吃 width/height, 盒子塌成 0x0:
元素在 DOM 里、计算样式也是 24px, 但既看不见也点不到。
移动端顶栏的主题 / 搜索 / 菜单三个按钮就这么消失过 —— 手机上等于没有导航入口。
*/
describe('unocss 图标规则', () => {
  it('生成的图标类带 inline-block', async () => {
    const uno = await createGenerator(config)

    // 这三个正是移动端顶栏用到的
    const { css } = await uno.generate('i-ic:sharp-menu i-ic:round-search i-mdi:weather-night', {
      preflights: false,
    })

    expect(css).toContain('i-ic\\:sharp-menu')
    // 每个图标类都要有 display, 不然在非 flex 父元素里是 0x0
    const rules = css.split('}').filter(r => r.includes('i-ic\\:') || r.includes('i-mdi\\:'))
    expect(rules).toHaveLength(3)
    rules.forEach(rule => expect(rule).toContain('display:inline-block'))
  })
})
