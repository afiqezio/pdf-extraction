import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: { enabled: true },
    // Application configuration
    app: {
        head: {
            title: 'Workbench Platform',
            meta: [
                { charset: 'utf-8' },
                { name: 'viewport', content: 'width=device-width, initial-scale=1' },
                { name: 'description', content: 'Modern project management platform' }
            ],
            link: [
                { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
            ]
        }
    },
    // Runtime configuration
    runtimeConfig: {
        public: {
            apiUrl: process.env.API_URL || 'http://localhost:8081'
        }
    },
    vite: {
        plugins: [
            tailwindcss(),
        ],
    },
    css: ['~/assets/css/main.css'],
    // SSR - disabled for now for simplicity
    ssr: false,
    modules: ['@nuxt/fonts', '@nuxt/icon', '@nuxt/image', '@pinia/nuxt',
        '@vueuse/nuxt',],
    // Auto imports
    imports: {
        dirs: ['stores']
    }
})