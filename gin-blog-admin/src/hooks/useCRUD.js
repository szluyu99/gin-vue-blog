import { $$, $computed, $ref } from 'vue/macros'

const ACTIONS = {
  view: '查看',
  edit: '编辑',
  add: '新增',
}

/**
 * @param {*} name 表单标题
 * @param {*} iniForm 初始表单内容
 * @param {*} doCreate 增加操作
 * @param {*} doDelete 删除操作
 * @param {*} doUpdate 修改操作
 * @param {*} refresh 查找(刷新)操作
 */
export default function ({ name, initForm = {}, doCreate, doDelete, doUpdate, refresh }) {
  let modalVisible = $ref(false) // 弹框显示
  /** 操作: add - 新增, edit - 删除, view - 查看 */
  let modalAction = $ref('')
  let modalLoading = $ref(false)
  const modalTitle = $computed(() => ACTIONS[modalAction] + name) // 弹窗标题

  let modalForm = $ref({ ...initForm })
  const modalFormRef = $ref(null)

  /** 新增 */
  function handleAdd() {
    modalAction = 'add'
    modalVisible = true
    modalForm = { ...initForm } // 使用初始表单
  }

  /** 修改 */
  function handleEdit(row) {
    modalAction = 'edit'
    modalVisible = true
    modalForm = { ...row }
  }

  /** 查看 */
  function handleView(row) {
    modalAction = 'view'
    modalVisible = true
    modalForm = { ...row }
  }

  /** 保存 */
  function handleSave() {
    if (!['edit', 'add'].includes(modalAction)) {
      modalVisible = false
      return
    }
    modalFormRef?.validate(async (err) => {
      if (!err) {
        const actions = {
          add: {
            api: () => doCreate(modalForm),
            cb: () => window.$message.success('新增成功'),
          },
          edit: {
            api: () => doUpdate(modalForm),
            cb: () => window.$message.success('编辑成功'),
          },
        }
        const action = actions[modalAction]

        try {
          modalLoading = true
          const data = await action.api()
          action.cb()
          modalLoading = modalVisible = false
          data && refresh(data)
        }
        catch (error) {
          modalLoading = false
        }
      }
    })
  }

  /**
   * 删除 (传入数组为批量删除, 传入单个 id 为普通删除)
   * @param {*} ids 主键数组
   * @param {boolean} needConfirm 是否需要确认窗口
   */
  async function handleDelete(ids, needConfirm = true) {
    if (!ids || (Array.isArray(ids) && !ids.length)) {
      window.$message.info('请选择要删除的数据')
      return
    }

    // 调用删除接口
    const callDeleteAPI = async () => {
      try {
        modalLoading = true
        const data = await doDelete(JSON.stringify(ids))
        // 针对软删除的情况做判断
        if (data?.code === 0)
          window.$message.success('删除成功')
        modalLoading = false
        refresh(data)
      }
      catch (error) {
        modalLoading = false
      }
    }

    needConfirm
      ? window.$dialog.confirm({
        content: '确定删除？',
        confirm: () => callDeleteAPI(),
      })
      : callDeleteAPI()
  }

  return $$({
    modalVisible,
    modalAction,
    modalTitle,
    modalLoading,
    handleAdd,
    handleDelete,
    handleEdit,
    handleView,
    handleSave,
    modalForm,
    modalFormRef,
  })
}
