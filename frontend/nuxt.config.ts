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
    css: ['~~/assets/main.css'],
    modules: ['@pinia/nuxt', '@nuxtjs/i18n', '@formkit/auto-animate/nuxt'],
    i18n: {
        defaultLocale: 'en',
    },
    $development: {
        modules: ['@primevue/nuxt-module'],
    },
    $production: {
        vite: {
            plugins: [
                Components({
                    resolvers: [PrimeVueResolver({ components: { prefix: 'Prime' } })],
                }),
            ],
        },
    },
    primevue: {
        usePrimeVue: false,
        loadStyles: false,
        components: {
            prefix: 'Prime',
        },
    },
    vite: {
        plugins: [tailwindcss()],
        server: {
            proxy: {
                '/api': `http://${GO_HOST}:${GO_PORT}`,
            },
        },
    },
    typescript: {
        tsConfig: {
            // include: ['./app/components.d.ts'],
        },
    },
})
