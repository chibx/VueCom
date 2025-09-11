import Aura from '@primeuix/themes/aura'
import PrimeVue from 'primevue/config'

export default defineNuxtPlugin((nuxt) => {
    nuxt.vueApp.use(PrimeVue, {
        theme: {
            preset: Aura,
        },
    })
})
