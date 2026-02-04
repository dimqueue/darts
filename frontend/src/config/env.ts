declare global {
    interface Window {
        __ENV__?: {
            VITE_API_URL?: string;
            VITE_USE_MOCK_API?: string;
        };
    }
}

const runtimeEnv = window.__ENV__ || {};

export const env = {
    VITE_API_URL:
        runtimeEnv.VITE_API_URL || import.meta.env.VITE_API_URL || 'http://localhost:8080',
    VITE_USE_MOCK_API: runtimeEnv.VITE_USE_MOCK_API || import.meta.env.VITE_USE_MOCK_API || 'false',
} as const;

export interface Config {
    apiUrl: string;
    useMockApi: boolean;
}

export const config: Config = {
    apiUrl: env.VITE_API_URL,
    useMockApi: env.VITE_USE_MOCK_API === 'true',
};

export default config;
