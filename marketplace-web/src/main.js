import { createApp } from 'vue'
import { createPinia } from 'pinia' // For state management

import App from './App.vue'
import router from './router'

// Import global styles
import './assets/main.css'
import 'highlight.js/styles/atom-one-dark.css';

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')

