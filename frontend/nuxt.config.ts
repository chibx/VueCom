import { PrimeVueResolver } from '@primevue/auto-import-resolver'
import tailwindcss from '@tailwindcss/vite'
import Components from 'unplugin-vue-components/vite'

const { GO_PORT, GO_HOST } = process.env

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    ssr: false,
    devtools: { enabled: false },
    nitro: {
        prerender: {
            crawlLinks: false,
            autoSubfolderIndex: false,
        },
    },
    modules: ['@pinia/nuxt', '@nuxtjs/i18n', '@formkit/auto-animate/nuxt'],
    i18n: {
        defaultLocale: 'en',
    },
    vite: {
        plugins: [
            tailwindcss(),
            Components({
                resolvers: [PrimeVueResolver({ components: { prefix: 'Prime' } })],
            }),
        ],
        server: {
            proxy: {
                '/api': `http://${GO_HOST}:${GO_PORT}`,
            },
        },
    },
    typescript: {
        tsConfig: {
            include: ['./app/components.d.ts'],
        },
    },
})
