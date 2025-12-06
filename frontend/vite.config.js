import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), '')

    return {
        base: env.VITE_BASE_PATH || '/',
        plugins: [react()],
        resolve: {
            alias: {
                '@': path.resolve(__dirname, 'src')
            }
        },
        server: {
            proxy: {
                '/api': 'http://localhost:8080',
                '/auth': 'http://localhost:8080'
            }
        }
    }
})