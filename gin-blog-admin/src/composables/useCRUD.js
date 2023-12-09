import { computed, ref } from 'vue'
import { useForm } from './useForm'

const ACTIONS = {
  view: '查看',
  edit: '编辑',
  add: '新增',
}

/**
 * @typedef {object} FormObject
 * @property {string} name - 名称
 * @property {object} initForm - 初始表单对象
 * @property {Function} doCreate - 执行创建操作的函数
 * @property {Function} doDelete - 执行删除操作的函数
 * @property {Function} doUpdate - 执行更新操作的函数
 * @property {Function} refresh - 刷新函数
 */

/**
 * 可复用的 CRUD 操作
 * @param {FormObject} options
 */
export function useCRUD({ name, initForm = {}, doCreate, doDelete, doUpdate, refresh }) {
  const modalVisible = ref(false) // 弹框显示
  /** @type {'add' | 'edit' | 'view'} 弹窗操作类型 */
  const modalAction = ref('')
  /** 弹窗加载状态 */
  const modalLoading = ref(false)
  /** 弹窗标题 */
  const modalTitle = computed(() => ACTIONS[modalAction.value] + name) // 弹窗标题

  const { formModel: modalForm, formRef: modalFormRef, validation } = useForm(initForm)

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
  async function handleSave() {
    if (!['edit', 'add'].includes(modalAction.value)) {
      modalVisible.value = false
      return
    }

    if (!(await validation())) {
      return false
    }
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
      console.error(error)
      modalLoading.value = false
    }
  }

  /**
   * 删除 (传入数组为批量删除, 传入单个 id 为普通删除)
   * @param {Array} ids 主键数组
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

        // TODO: 优化逻辑
        let data
        if (typeof ids === 'number' || typeof ids === 'string') {
          data = await doDelete(ids)
        }
        else {
          data = await doDelete(JSON.stringify(ids))
        }

        // 针对软删除的情况做判断
        if (data?.code === 0) {
          window.$message.success('删除成功')
        }
        modalLoading.value = false
        refresh(data)
      }
      catch (error) {
        console.error(error)
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
