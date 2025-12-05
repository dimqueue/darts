import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), '')

    const useMockApi = mode === 'pages' ||
        (process.env.VITE_USE_MOCK_API || env.VITE_USE_MOCK_API) === 'true'

    const apiEntry = useMockApi
        ? path.resolve(__dirname, 'src/api/index.mock.js')
        : path.resolve(__dirname, 'src/api/index.prod.js')

    return {
        base: process.env.VITE_BASE_PATH || env.VITE_BASE_PATH || '/',
        plugins: [react()],
        resolve: {
            alias: {
                '@/api': apiEntry
            }
        },
        server: {
            proxy: {
                '/api': 'http://localhost:8080',
                '/auth': 'http://localhost:8080'
            }
        },
        define: {
            '__USE_MOCK_API__': useMockApi ? 'true' : 'false'
        }
    }
})