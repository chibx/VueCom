// import Aura from '@primeuix/themes/aura'
// import PrimeVue from 'primevue/config'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'

import App from './app.vue'
import router from './router'
import '@@/assets/main.css'

const app = createApp(App)

const i18n = createI18n({ legacy: false })

app.use(createPinia())
app.use(i18n)
app.use(router)
// app.use(PrimeVue, {
//   theme: {
//     preset: Aura,
//   },
// })

app.mount('#app')
