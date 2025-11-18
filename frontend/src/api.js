const USE_MOCK = __USE_MOCK_API__;

let ApiClient;

if (USE_MOCK) {
    const { default: MockApiWrapper } = await import('./api-mock.js');
    ApiClient = MockApiWrapper;
} else {
    const { default: RealApiClient } = await import('./api-real.js');
    ApiClient = RealApiClient;
}

export const api = new ApiClient();
export default api;
