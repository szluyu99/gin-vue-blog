// custom style
import './styles/index.css'
import './styles/common.css'
import './styles/animate.css'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

// vue
import { createApp } from 'vue'

import { router } from './router'
import { pinia } from './store'
import App from './App.vue'

const app = createApp(App)
app.use(router)
app.use(pinia)
app.mount('#app')
