import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'
import UploadOne from './UploadOne.vue'

describe('uploadOne', () => {
  beforeEach(() => setActivePinia(createPinia()))

  it('preview 后到时也能显示出来', async () => {
    // 回归 FE1: 父组件的 avatar 来自 onMounted 里的异步 getUserInfo,
    // 一定晚于本组件创建。不 watch props 就会永远停在初始值
    const wrapper = mount(UploadOne, { props: { preview: '' } })
    expect(wrapper.find('img').exists()).toBe(false)

    await wrapper.setProps({ preview: 'public/uploaded/a.jpg' })

    const img = wrapper.find('img')
    expect(img.exists()).toBe(true)
    // FE4: 根相对路径, 不能拼 localhost
    expect(img.attributes('src')).toBe('/public/uploaded/a.jpg')
  })

  it('http 开头的头像原样展示', async () => {
    const wrapper = mount(UploadOne, { props: { preview: 'https://cdn.test/a.png' } })
    expect(wrapper.find('img').attributes('src')).toBe('https://cdn.test/a.png')
  })
})
