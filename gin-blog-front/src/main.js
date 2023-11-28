// 自定义样式
import './styles/index.scss'
import './styles/animate.scss'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

// vue
import { createApp } from 'vue'
import { router } from './router'
import { pinia } from './store'
import App from './App.vue'

const app = createApp(App)
app.use(router).use(pinia)
app.mount('#app')
