import { createApp } from 'vue'
import { createPinia } from 'pinia' // For state management

import App from './App.vue'
import router from './router'

// Import global styles
import './assets/main.css'

// Import third-party scripts that need to run after Vue is ready
// We will call these from App.vue instead to ensure they run after each route change

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')

