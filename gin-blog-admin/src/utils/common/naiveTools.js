import { isNullOrUndef } from '@/utils'

export function setupMessage(NMessage) {
  let loadingMessage = null
  class Message {
    /**
     * 规则：
     * * loading message 只显示一个，新的 message 会替换正在显示的 loading message
     * * loading message 不会自动清除，除非被替换成非 loading message，非 loading message 默认 2 秒后自动清除
     */
    removeMessage(message = loadingMessage, duration = 2000) {
      setTimeout(() => {
        if (message) {
          message.destroy()
          message = null
        }
      }, duration)
    }

    showMessage(type, content, option = {}) {
      if (loadingMessage && loadingMessage.type === 'loading') {
        // 如果存在则替换正在显示的 loading message
        loadingMessage.type = type
        loadingMessage.content = content

        if (type !== 'loading') {
          // 非 loading message 需设置自动清除
          this.removeMessage(loadingMessage, option.duration)
        }
      }
      else {
        // 不存在正在显示的 loading 则新建一个 message
        // 如果新建的 message 是 loading message 则将 message 赋值存储下来
        const message = NMessage[type](content, option)
        if (type === 'loading')
          loadingMessage = message
      }
    }

    loading(content) {
      this.showMessage('loading', content, { duration: 0 })
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

export function setupDialog(NDialog) {
  NDialog.confirm = function (option = {}) {
    const showIcon = !isNullOrUndef(option.title)
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
