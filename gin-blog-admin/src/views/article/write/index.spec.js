import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import ArticleWrite from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getCategoryOption: vi.fn().mockResolvedValue({ code: 0, data: [] }),
    getTagOption: vi.fn().mockResolvedValue({ code: 0, data: [] }),
    getArticleById: vi.fn(),
    saveOrUpdateArticle: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

let routeParams = {}
vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRoute: () => ({ params: routeParams, query: {} }),
  useRouter: () => ({ push: vi.fn() }),
}))

vi.mock('md-editor-v3', () => ({ MdEditor: { template: '<div />' } }))

function mountPage() {
  return mount(ArticleWrite, {
    global: {
      stubs: {
        CommonPage: { template: '<div><slot name="action" /><slot /></div>' },
        UploadOne: true,
        CrudModal: true,
      },
    },
  })
}

describe('写文章', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    routeParams = {}
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn(), warning: vi.fn() }
    window.$loadingBar = { start: vi.fn(), finish: vi.fn(), error: vi.fn() }
  })

  it('新建时 tag_names 是数组, watcher 不会抛 TypeError', async () => {
    // 回归 A5: 新建分支重置表单时漏了 tag_names,
    // watch(() => formModel.tag_names) 紧接着 newVal.includes(...) 就会炸
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getCategoryOption).toHaveBeenCalled())

    expect(wrapper.vm.formModel.tag_names).toEqual([])
    expect(wrapper.vm.formModel.category_name).toBe('')

    // 触发 watcher, 不应抛错
    wrapper.vm.formModel.tag_names = ['Go']
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.formModel.tag_names).toEqual(['Go'])
  })

  it('编辑导入的草稿(无分类无标签)不会崩', async () => {
    // 同 A1 的根因: category 为 null
    routeParams = { id: '4' }
    api.getArticleById.mockResolvedValue({
      code: 0,
      data: { id: 4, title: 'e2e-import', content: '# x', category: null, tags: [], status: 3 },
    })

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.formModel.title).toBe('e2e-import'))

    expect(wrapper.vm.formModel.tag_names).toEqual([])
    expect(wrapper.vm.formModel.category_name).toBe('')
  })
})
