import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'http://localhost:8080', // Assuming your Go backend runs on port 8080
  headers: {
    'Content-Type': 'application/json',
  },
});

// Later, we will add an interceptor here to automatically add
// the JWT token to every outgoing request.

export default apiClient;
