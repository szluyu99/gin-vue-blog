// vue
import { createApp } from 'vue'
import App from './App.vue'
import { router } from './router'

import { pinia } from './store'
import { setupMock } from './utils/http'
// custom style
import './styles/index.css'

import './styles/common.css'

import './styles/animate.css'
// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

async function bootstrap() {
  await setupMock() // mock 模式下需要在发出请求前装上适配器

  const app = createApp(App)
  app.use(router)
  app.use(pinia)
  app.mount('#app')
}

bootstrap()
