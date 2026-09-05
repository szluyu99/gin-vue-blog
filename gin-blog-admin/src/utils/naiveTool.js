import * as NaiveUI from 'naive-ui'
import { computed } from 'vue'
import themes from '@/assets/themes'
import { useThemeStore } from '@/store'

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
    // themes.js 导出的键名是 naiveThemeOverrides(App.vue 用的就是这个),
    // 原来写 themes.themeOverrides 取到 undefined, 离散 API 全部回落到 naive 默认绿
    themeOverrides: themes.naiveThemeOverrides,
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

/**
 * 把 naiveThemeOverrides.common 里的色值写成 :root 上的 CSS 变量
 *
 * uno.config.js 声明了 primary / info / success / warning / error 五组语义色, 都指向
 * CSS 变量, 但除 --primary-color 外没人给这些变量赋值 —— 实测 --info-color、
 * --success-color、--warning-color 全是空的, 于是 text-info、bg-success 这类类名
 * 静默失效(颜色回落到继承值, 看起来就是"没生效"), 20 个 token 里 19 个是死的。
 *
 * 命名对齐 uno.config.js: primaryColorHover -> --primary-color-hover;
 * naive 没有 Active 这一档, 用它的 Suppl 顶上。
 */
export function setupThemeVars(el = document.documentElement) {
  const common = themes.naiveThemeOverrides?.common ?? {}
  for (const [key, value] of Object.entries(common)) {
    const name = key
      .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
      .toLowerCase()
      .replace(/-suppl$/, '-active')
    el.style.setProperty(`--${name}`, value)
  }
}
