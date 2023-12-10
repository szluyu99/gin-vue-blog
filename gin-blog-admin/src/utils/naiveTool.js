import { computed } from 'vue'
import * as NaiveUI from 'naive-ui'
import { useThemeStore } from '@/store'
import themes from '@/assets/themes'

function setupMessage(NMessage) {
  class Message {
    static instance
    message
    removeTimer

    constructor() {
      if (Message.instance) {
        return Message.instance
      }
      Message.instance = this
      this.message = {}
      this.removeTimer = {}
    }

    destroy(key, duration = 200) {
      setTimeout(() => {
        if (this.message[key]) {
          this.message[key].destroy()
          delete this.message[key]
        }
      }, duration)
    }

    removeMessage(key, duration = 5000) {
      this.removeTimer[key] && clearTimeout(this.removeTimer[key])
      this.removeTimer[key] = setTimeout(() => {
        this.message[key]?.destroy()
      }, duration)
    }

    showMessage(type, content, option = {}) {
      if (Array.isArray(content)) {
        return content.forEach(msg => NMessage[type](msg, option))
      }

      if (!option.key) {
        return NMessage[type](content, option)
      }

      const currentMessage = this.message[option.key]
      if (currentMessage) {
        currentMessage.type = type
        currentMessage.content = content
      }
      else {
        this.message[option.key] = NMessage[type](content, {
          ...option,
          duration: 0,
          onAfterLeave: () => {
            delete this.message[option.key]
          },
        })
      }
      this.removeMessage(option.key, option.duration)
    }

    loading(content, option = { duration: 0 }) {
      this.showMessage('loading', content, option)
    }

    success(content, option = {}) {
      this.showMessage('success', content, option)
    }

    error(content, option = {}) {
      this.showMessage('error', content, option)
    }

    info(content, option = {}) {
      this.showMessage('info', content, option)
    }

    warning(content, option = {}) {
      this.showMessage('warning', content, option)
    }
  }

  return new Message()
}

function setupDialog(NDialog) {
  NDialog.confirm = function (option = {}) {
    const showIcon = !!(option.title)
    return NDialog[option.type || 'warning']({
      showIcon,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: option.confirm,
      onNegativeClick: option.cancel,
      onMaskClick: option.cancel,
      ...option,
    })
  }
  return NDialog
}

/**
 * 挂载 NaiveUI API
 */
export function setupNaiveDiscreteApi() {
  const themeStore = useThemeStore()
  const configProviderProps = computed(() => ({
    theme: themeStore.darkMode ? NaiveUI.darkTheme : undefined,
    themeOverrides: themes.themeOverrides,
  }))
  const { message, dialog, notification, loadingBar } = NaiveUI.createDiscreteApi(
    ['message', 'dialog', 'notification', 'loadingBar'],
    { configProviderProps },
  )

  window.$loadingBar = loadingBar
  window.$notification = notification
  window.$message = setupMessage(message)
  window.$dialog = setupDialog(dialog)
}

/**
 * 解决 naive-ui 和 unocss 样式冲突
 */
export function setupNaiveUnocss() {
  const meta = document.createElement('meta')
  meta.name = 'naive-ui-style'
  document.head.appendChild(meta)
}
