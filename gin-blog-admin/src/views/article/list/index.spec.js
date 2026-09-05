import { mount } from '@vue/test-utils'
import { NDataTable } from 'naive-ui'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import ArticleList from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getArticles: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    getCategoryOption: vi.fn().mockResolvedValue({ code: 0, data: [] }),
    getTagOption: vi.fn().mockResolvedValue({ code: 0, data: [] }),
    deleteArticle: vi.fn().mockResolvedValue({ code: 0 }),
    softDeleteArticle: vi.fn().mockResolvedValue({ code: 0 }),
    updateArticleTop: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

// 只替换两个组合式函数: 整个模块换掉会让 src/router 里的 createRouter 变成 undefined
vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRoute: () => ({ query: {}, params: {} }),
  useRouter: () => ({ push: vi.fn() }),
}))

function mountPage() {
  return mount(ArticleList, {
    global: {
      stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } },
    },
  })
}

// 导入生成的草稿没有分类, 后端返回 category: null
const draft = { id: 4, title: 'e2e-import', category: null, tags: [], status: 3, is_delete: 0 }
const normal = { id: 3, title: '有分类', category: { id: 1, name: '项目' }, tags: [], status: 1, is_delete: 0 }

describe('文章列表', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn(), warning: vi.fn() }
    window.$dialog = { confirm: ({ confirm }) => confirm() }
  })

  it('分类为 null 的文章不会让表格渲染崩掉', async () => {
    // 回归 A1: 原来是 row.category.name, 遇到导入的草稿会抛在 render 里
    api.getArticles.mockResolvedValue({
      code: 0,
      data: { page_data: [normal, draft], total: 2 },
    })

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.findComponent(NDataTable).props('data')).toHaveLength(2))
    await wrapper.vm.$nextTick()

    const text = wrapper.text()
    expect(text).toContain('项目')
    expect(text).toContain('无') // 没有分类时展示「无」
  })

  it('删除文章会把结果 return 出去, 让 useCRUD 能拿到业务码', async () => {
    // 回归 A3: updateOrDeleteArticles 少了 return, useCRUD 里 await 到 undefined,
    // 删成功不弹提示, 失败也捕获不到
    api.getArticles.mockResolvedValue({ code: 0, data: { page_data: [normal], total: 1 } })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getArticles).toHaveBeenCalled())

    const ret = wrapper.vm.updateOrDeleteArticles(JSON.stringify([3]))
    expect(ret).toBeInstanceOf(Promise)
    await expect(ret).resolves.toMatchObject({ code: 0 })
  })

  // 乐观更新: 失败后必须回滚, 否则开关显示已置顶而后端没变
  // (菜单/接口/用户三个页面的同类开关都在 catch 里回滚, 只有这里漏了)
  it('置顶失败后回滚开关状态', async () => {
    api.updateArticleTop.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    const row = { id: 3, is_top: false }

    await wrapper.vm.handleUpdateTop(row)

    expect(row.is_top).toBe(false)
    expect(window.$message.success).not.toHaveBeenCalled()
  })

  it('置顶成功后保留新状态并提示', async () => {
    api.updateArticleTop.mockResolvedValue({ code: 0 })
    const wrapper = mountPage()
    const row = { id: 3, is_top: false }

    await wrapper.vm.handleUpdateTop(row)

    expect(row.is_top).toBe(true)
    expect(window.$message.success).toHaveBeenCalledWith('已成功置顶')
  })

  // 相邻的置顶/审核/删除都有提示, 只有恢复原来既不 catch 也不提示
  it('从回收站恢复会提示成功', async () => {
    api.softDeleteArticle.mockResolvedValue({ code: 0 })
    const wrapper = mountPage()

    await wrapper.vm.handleRestore({ id: 3 })

    expect(api.softDeleteArticle).toHaveBeenCalledWith([3], false)
    expect(window.$message.success).toHaveBeenCalledWith('已恢复该文章')
  })

  it('恢复失败不抛出未捕获的 rejection', async () => {
    api.softDeleteArticle.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()

    await expect(wrapper.vm.handleRestore({ id: 3 })).resolves.toBeUndefined()
    expect(window.$message.success).not.toHaveBeenCalled()
  })

  it('导入放行 .md 与 .markdown, 拦住其他后缀', async () => {
    // 回归 A20: 后端 F9 之后同时接受两种后缀, 前端只放行 .md
    api.getArticles.mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    const wrapper = mountPage()

    expect(wrapper.vm.beforeUpload({ file: { name: 'a.md' } })).toBe(true)
    expect(wrapper.vm.beforeUpload({ file: { name: 'a.markdown' } })).toBe(true)
    expect(wrapper.vm.beforeUpload({ file: { name: 'a.MD' } })).toBe(true)
    expect(wrapper.vm.beforeUpload({ file: { name: 'a.txt' } })).toBe(false)
  })
})
