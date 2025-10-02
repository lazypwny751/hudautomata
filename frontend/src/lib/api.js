import ky from 'ky';

// Production'da Caddy proxy kullanıyoruz, development'ta localhost:8080
const API_URL = import.meta.env.VITE_API_URL || (
  import.meta.env.MODE === 'production' ? '' : 'http://localhost:8080'
);

// Create API client
const api = ky.create({
  prefixUrl: `${API_URL}/api/v1`,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = localStorage.getItem('token');
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      },
    ],
  },
});

// Auth API
export const authAPI = {
  login: (credentials) => api.post('auth/login', { json: credentials }).json(),
  logout: () => api.post('auth/logout').json(),
  getMe: () => api.get('auth/me').json(),
};

// Users API
export const usersAPI = {
  list: (params = {}) => api.get('users', { searchParams: params }).json(),
  create: (data) => api.post('users', { json: data }).json(),
  get: (id) => api.get(`users/${id}`).json(),
  update: (id, data) => api.put(`users/${id}`, { json: data }).json(),
  delete: (id) => api.delete(`users/${id}`).json(),
  getByRFID: (cardId) => api.get(`users/rfid/${cardId}`).json(),
  getBalance: (id) => api.get(`users/${id}/balance`).json(),
  getTransactions: (id) => api.get(`users/${id}/transactions`).json(),
};

// Transactions API
export const transactionsAPI = {
  list: (params = {}) => api.get('transactions', { searchParams: params }).json(),
  create: (data) => api.post('transactions', { json: data }).json(),
  get: (id) => api.get(`transactions/${id}`).json(),
};

// Automation API
export const automationAPI = {
  scan: (data) => api.post('automation/scan', { json: data }).json(),
  checkBalance: (rfidCardId) => api.post('automation/check-balance', { json: { rfid_card_id: rfidCardId } }).json(),
  getHistory: (params = {}) => api.get('automation/history', { searchParams: params }).json(),
};

// Dashboard API
export const dashboardAPI = {
  getStats: () => api.get('dashboard/stats').json(),
  getCharts: (params = {}) => api.get('dashboard/charts', { searchParams: params }).json(),
  getRecent: () => api.get('dashboard/recent').json(),
};

// Admins API
export const adminsAPI = {
  list: () => api.get('admins').json(),
  create: (data) => api.post('admins', { json: data }).json(),
  get: (id) => api.get(`admins/${id}`).json(),
  delete: (id) => api.delete(`admins/${id}`).json(),
};

// Logs API
export const logsAPI = {
  list: (params = {}) => api.get('logs', { searchParams: params }).json(),
};

export default api;
