import RealApiClient from './realApi';
import MockApiClient from './mockApi';
import config from '../config/env';

const api = config.useMockApi ? new MockApiClient() : new RealApiClient();

if (import.meta.env.DEV) {
    console.log(`API: ${config.useMockApi ? 'Mock' : 'Real'} (${config.apiUrl})`);
}

export default api;
