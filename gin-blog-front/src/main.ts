// vuetify 按需引入
import 'vuetify/styles'
import { createVuetify } from 'vuetify'

// 自定义样式
import '@/styles/reset.css'
import '@/styles/index.scss'
import '@/styles/animate.scss'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

// vue
import { createApp } from 'vue'
import { setupRouter } from './router'
import { setupStore } from './store'
import App from './App.vue'

const app = createApp(App)
setupStore(app)
setupRouter(app)
app.use(createVuetify())
app.mount('#app')
