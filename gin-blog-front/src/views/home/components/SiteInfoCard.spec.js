import { mount } from '@vue/test-utils'
import dayjs from 'dayjs'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { useAppStore } from '@/store'
import SiteInfoCard from './SiteInfoCard.vue'

function mountCard() {
  return mount(SiteInfoCard, {
    global: { stubs: { RouterLink: { template: '<a><slot /></a>' } } },
  })
}

function fillStore(overrides = {}) {
  const store = useAppStore()
  store.blogInfo = {
    article_count: 13,
    category_count: 3,
    tag_count: 6,
    view_count: 233,
    blog_config: { website_notice: '欢迎来到我的博客' },
    ...overrides.blogInfo,
  }
  store.blog_config = {
    website_author: '博主',
    website_intro: '一句简介',
    website_avatar: '/avatar.png',
    website_createtime: dayjs().subtract(400, 'day').toISOString(),
    qq: '123',
    github: 'https://github.com/x',
    gitee: 'https://gitee.com/x',
    ...overrides.blogConfig,
  }
  return store
}

describe('侧边信息卡', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  // 原来是三张卡片各自 hidden lg:block, 窄屏整列消失; 合并后不再有隐藏类
  it('一张卡里同时有博主、计数、公告和网站资讯', async () => {
    fillStore()
    const wrapper = mountCard()
    await wrapper.vm.$nextTick()

    const text = wrapper.text()
    expect(text).toContain('博主')
    expect(text).toContain('一句简介')
    expect(text).toContain('13')
    expect(text).toContain('欢迎来到我的博客')
    expect(text).toContain('总访问量')
    expect(text).toContain('233')
    expect(wrapper.classes()).not.toContain('hidden')
  })

  // 三个计数都是入口, 指向归档 / 分类 / 标签
  it('计数是可点的链接', async () => {
    fillStore()
    const wrapper = mountCard()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.counts.map(c => c.path)).toEqual(['/archives', '/categories', '/tags'])
    expect(wrapper.vm.counts.map(c => c.value)).toEqual([13, 3, 6])
  })

  // 没配公告就不要留一个空标题
  it('公告为空时整段不渲染', async () => {
    fillStore({ blogInfo: { blog_config: {} } })
    const wrapper = mountCard()
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).not.toContain('公告')
  })

  /*
    运行时间必须用 asDays 算总天数

    原来若用 duration.format('D 天'), D 是"天"这个分量, 建站超过一个月就归零重算:
    400 天会显示成 "5 天"。
  */
  it('运行时间按总天数算, 不会每月归零', async () => {
    fillStore()
    const wrapper = mountCard()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.runTime).toMatch(/^(399|400) 天/)
  })

  // 配置还没拉回来时不能显示 NaN
  it('建站时间缺失时显示占位', async () => {
    fillStore({ blogConfig: { website_createtime: '' } })
    const wrapper = mountCard()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.runTime).toBe('-')
    expect(wrapper.text()).not.toContain('NaN')
  })

  // 建站时间是接口回来的, 到了之后要立刻重算而不是等下一个 30 秒
  it('建站时间后到也会重算', async () => {
    fillStore({ blogConfig: { website_createtime: '' } })
    const wrapper = mountCard()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.runTime).toBe('-')

    const store = useAppStore()
    store.blog_config = { ...store.blog_config, website_createtime: dayjs().subtract(10, 'day').toISOString() }
    await vi.waitFor(() => expect(wrapper.vm.runTime).toMatch(/^(9|10) 天/))
  })
})
