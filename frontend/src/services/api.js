import axios from 'axios';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Later, we will add an interceptor here to automatically add
// the JWT token to every outgoing request.

export default apiClient;
