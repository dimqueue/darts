import RealApiClient from './realApi';
import MockApiClient from './mockApi';

const USE_MOCK_API = import.meta.env.VITE_USE_MOCK_API === 'true';

const api = USE_MOCK_API ? new MockApiClient() : new RealApiClient();

export default api;