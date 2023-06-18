import { computed, ref } from 'vue'

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
  const modalVisible = ref(false) // 弹框显示
  /** 操作: add - 新增, edit - 删除, view - 查看 */
  const modalAction = ref('')
  const modalLoading = ref(false)
  const modalTitle = computed(() => ACTIONS[modalAction.value] + name) // 弹窗标题

  const modalForm = ref({ ...initForm })
  const modalFormRef = ref(null)

  /** 新增 */
  function handleAdd() {
    modalAction.value = 'add'
    modalVisible.value = true
    modalForm.value = { ...initForm } // 使用初始表单
  }

  /** 修改 */
  function handleEdit(row) {
    modalAction.value = 'edit'
    modalVisible.value = true
    modalForm.value = { ...row }
  }

  /** 查看 */
  function handleView(row) {
    modalAction.value = 'view'
    modalVisible.value = true
    modalForm.value = { ...row }
  }

  /** 保存 */
  function handleSave() {
    if (!['edit', 'add'].includes(modalAction.value)) {
      modalVisible.value = false
      return
    }
    modalFormRef.value?.validate(async (err) => {
      if (!err) {
        const actions = {
          add: {
            api: () => doCreate(modalForm.value),
            cb: () => window.$message.success('新增成功'),
          },
          edit: {
            api: () => doUpdate(modalForm.value),
            cb: () => window.$message.success('编辑成功'),
          },
        }
        const action = actions[modalAction.value]

        try {
          modalLoading.value = true
          const data = await action.api()
          action.cb()
          modalLoading.value = modalVisible.value = false
          data && refresh(data)
        }
        catch (error) {
          modalLoading.value = false
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
        modalLoading.value = true
        const data = await doDelete(JSON.stringify(ids))
        // 针对软删除的情况做判断
        if (data?.code === 0)
          window.$message.success('删除成功')
        modalLoading.value = false
        refresh(data)
      }
      catch (error) {
        modalLoading.value = false
      }
    }

    if (needConfirm) {
      window.$dialog.confirm({
        content: '确定删除？',
        confirm: () => callDeleteAPI(),
      })
    }
    else {
      callDeleteAPI()
    }
  }

  return {
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
  }
}
