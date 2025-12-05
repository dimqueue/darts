import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), '')

    return {
        base: env.VITE_BASE_PATH || '/',
        plugins: [react()],
        server: {
            proxy: {
                '/api': 'http://localhost:8080',
                '/auth': 'http://localhost:8080'
            }
        },
        define: {
            // Use environment variable to determine mock API usage
            '__USE_MOCK_API__': env.VITE_USE_MOCK_API === 'true' ? 'true' : 'false'
        }
    }
})