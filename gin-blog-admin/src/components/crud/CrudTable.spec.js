import { mount } from '@vue/test-utils'
import { NDataTable } from 'naive-ui'
import { describe, expect, it, vi } from 'vitest'
import CrudTable from './CrudTable.vue'

const columns = [{ title: '标题', key: 'title' }]

function makeGetData() {
  return vi.fn().mockResolvedValue({
    code: 0,
    data: { page_data: [{ id: 1, title: 'a' }], total: 30 },
  })
}

function mountTable(getData) {
  return mount(CrudTable, {
    props: { columns, getData },
    global: { stubs: { NButton: true, NSpace: true } },
  })
}

describe('crudTable', () => {
  it('翻页只发一次请求', async () => {
    // 回归 A6: naive-ui 的 mergedOnUpdatePage 会先调 pagination.onChange
    // 再触发 onUpdate:page, 两处都发请求的话每页会请求两次
    const getData = makeGetData()
    const wrapper = mountTable(getData)
    await wrapper.vm.handleQuery()
    expect(getData).toHaveBeenCalledTimes(1)

    const table = wrapper.findComponent(NDataTable)
    table.vm.$emit('update:page', 2)
    await wrapper.vm.$nextTick()

    expect(getData).toHaveBeenCalledTimes(2)
    expect(getData.mock.calls[1][0]).toMatchObject({ page_num: 2, page_size: 10 })
  })

  it('pagination 上不再挂 onChange', async () => {
    // onChange 和 @update:page 同时存在就是双请求的根源
    const wrapper = mountTable(makeGetData())
    const table = wrapper.findComponent(NDataTable)
    expect(table.props('pagination').onChange).toBeUndefined()
  })

  it('请求失败时清空数据并归零总数', async () => {
    const getData = vi.fn().mockRejectedValue(new Error('boom'))
    const wrapper = mountTable(getData)
    await wrapper.vm.handleQuery()
    expect(wrapper.vm.tableData).toEqual([])
  })
})
