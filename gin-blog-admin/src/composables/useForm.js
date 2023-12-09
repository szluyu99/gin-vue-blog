import { ref } from 'vue'

/**
 * 可复用的表单对象
 * @param {any} initForm 表单初始值
 */
export function useForm(initForm = {}) {
  const formRef = ref(null)
  const formModel = ref({ ...initForm })

  const validation = async () => {
    try {
      await formRef.value?.validate()
      return true
    }
    catch (error) {
      return false
    }
  }

  const rules = {
    required: {
      required: true,
      message: '此为必填项',
      trigger: ['blur', 'change'],
    },
  }
  return { formRef, formModel, validation, rules }
}
