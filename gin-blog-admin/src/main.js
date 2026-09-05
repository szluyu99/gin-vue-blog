import { createApp } from 'vue'
import App from './App.vue'

import { setupRouter } from './router'
import { setupStore } from './store'
import { setupMock, setupNaiveDiscreteApi, setupNaiveUnocss, setupThemeVars } from './utils'
import '@unocss/reset/tailwind.css'
import 'uno.css'

async function bootstrap() {
  await setupMock() // mock 模式下需要在发出请求前装上适配器
  const app = createApp(App)
  setupStore(app) // 优先级最高
  setupNaiveUnocss()
  setupThemeVars() // 让 uno.config.js 里那几组语义色真的有值
  setupNaiveDiscreteApi()
  await setupRouter(app)
  app.mount('#app')
}

bootstrap()
