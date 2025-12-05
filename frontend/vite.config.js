import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), '')

    return {
        base: process.env.VITE_BASE_PATH || env.VITE_BASE_PATH || '/',
        plugins: [react()],
        server: {
            proxy: {
                '/api': 'http://localhost:8080',
                '/auth': 'http://localhost:8080'
            }
        },
        define: {
            '__USE_MOCK_API__': (process.env.VITE_USE_MOCK_API || env.VITE_USE_MOCK_API) === 'true' ? 'true' : 'false'
        }
    }
})