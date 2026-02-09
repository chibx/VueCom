// import { PrimeVueResolver } from '@primevue/auto-import-resolver'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'

const GO_HOST = 'localhost'
const GO_PORT = 2500

// https://vite.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        // vueDevTools(),
        tailwindcss(),
        Components({
            directoryAsNamespace: true,
        }),
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
            '~': fileURLToPath(new URL('./src', import.meta.url)),
            '@@': fileURLToPath(new URL('.', import.meta.url)),
        },
    },
    server: {
        proxy: {
            '/api': {
                target: `http://${GO_HOST}:${GO_PORT}/`,
                changeOrigin: true,
            },
        },
    },
})
