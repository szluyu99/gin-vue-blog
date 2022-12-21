// 重置样式
import '@/styles/reset.css'
import '@/styles/global.scss'
import 'uno.css'
import 'virtual:svg-icons-register' // vite-plugin-svg-icons

import { createApp } from 'vue'
import App from './App.vue'
import { setupRouter } from '@/router'
import { setupStore } from '@/store'

async function setupApp() {
  const app = createApp(App)
  setupStore(app)
  await setupRouter(app)
  app.mount('#app')
}

// 每次刷新页面都会执行 main.js 中的内容
setupApp()
