import '@unocss/reset/tailwind.css'
import 'uno.css'
import './style.scss'

import { createApp } from 'vue'
import App from './App.vue'
import { setupNaiveDiscreteApi } from './utils'
import { setupRouter } from '@/router'
import { setupStore } from '@/store'

async function setupApp() {
  const app = createApp(App)
  setupStore(app)
  setupNaiveDiscreteApi()
  await setupRouter(app)

  // 解决 naive-ui 与 unocss 的样式冲突
  const meta = document.createElement('meta')
  meta.name = 'naive-ui-style'
  document.head.appendChild(meta)
  app.mount('#app')
}

// 每次刷新页面都会执行 main.js 中的内容
setupApp()
