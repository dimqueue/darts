/**
 * Environment configuration helper
 *
 * Reads from:
 * - window.__ENV__ (runtime injection in Docker)
 * - import.meta.env (Vite build-time variables)
 *
 * Runtime values take precedence over build-time values
 */

const runtimeEnv = window.__ENV__ || {};

export const env = {
    VITE_API_URL: runtimeEnv.VITE_API_URL || import.meta.env.VITE_API_URL || 'http://localhost:8080',
    VITE_USE_MOCK_API: runtimeEnv.VITE_USE_MOCK_API || import.meta.env.VITE_USE_MOCK_API || 'false',
};

export const config = {
    apiUrl: env.VITE_API_URL,
    useMockApi: env.VITE_USE_MOCK_API === 'true',
};

export default config;