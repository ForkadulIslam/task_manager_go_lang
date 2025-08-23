import { defineStore } from 'pinia';
import apiClient from '../services/api';
import { ref } from 'vue';

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || null);
  const user = ref(null); // We can store user info here later

  function setToken(newToken) {
    localStorage.setItem('token', newToken);
    token.value = newToken;
    // Set the authorization header for all future requests
    apiClient.defaults.headers.common['Authorization'] = `Bearer ${newToken}`;
  }

  function clearToken() {
    localStorage.removeItem('token');
    token.value = null;
    delete apiClient.defaults.headers.common['Authorization'];
  }

  async function login(credentials) {
    try {
      const response = await apiClient.post('/login', credentials);
      if (response.data && response.data.token) {
        setToken(response.data.token);
        return true;
      } else {
        // Handle cases where the token is not in the expected place
        throw new Error('Login response did not contain a token.');
      }
    } catch (error) {
      console.error('Login failed:', error);
      // We can add more sophisticated error handling here
      // to show messages to the user.
      throw error; // Re-throw the error to be caught in the component
    }
  }

  function logout() {
    clearToken();
    // We might want to redirect the user to the login page
  }

  // Set auth header on initial load if token exists
  if (token.value) {
    apiClient.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
  }

  return {
    token,
    user,
    login,
    logout,
    isAuthenticated: !!token.value, // Simple getter-like property
  };
});
