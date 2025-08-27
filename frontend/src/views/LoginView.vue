<template>
  <div class="bg-gray-900 font-sans text-gray-100 w-full h-screen flex items-center justify-center">
    <div class="bg-gray-800 border border-gray-700 w-full max-w-md p-8 rounded-lg shadow-lg">
      
      <!-- Logo Placeholder -->
      <div class="text-center mb-6">
        <svg class="mx-auto h-12 w-auto text-sky-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>

      <h2 class="text-2xl font-bold text-center mb-1">Sign in to your account</h2>
      <p class="text-center text-sm text-gray-400 mb-8">Welcome back!</p>

      <form @submit.prevent="handleLogin" class="space-y-6">
        <!-- Error Message -->
        <div v-if="errorMessage" class="bg-red-900 border border-red-700 text-red-200 px-4 py-3 rounded-md text-sm">
          <p>{{ errorMessage }}</p>
        </div>

        <div>
          <label for="username" class="block text-sm font-medium text-gray-400 mb-2">Username</label>
          <div class="relative">
            <span class="absolute inset-y-0 left-0 flex items-center pl-3">
              <svg class="h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
              </svg>
            </span>
            <input 
              v-model="username"
              type="text" 
              id="username" 
              name="username" 
              class="block w-full pl-10 pr-3 py-2 bg-gray-700 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-sky-400 focus:border-transparent transition"
              placeholder="your_username"
              required
            >
          </div>
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-400 mb-2">Password</label>
          <div class="relative">
            <span class="absolute inset-y-0 left-0 flex items-center pl-3">
              <svg class="h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
              </svg>
            </span>
            <input 
              v-model="password"
              type="password" 
              id="password" 
              name="password" 
              class="block w-full pl-10 pr-3 py-2 bg-gray-700 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-sky-400 focus:border-transparent transition"
              placeholder="••••••••"
              required
            >
          </div>
        </div>

        <div>
          <button 
            type="submit" 
            :disabled="isLoading"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-sky-500 hover:bg-sky-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-sky-400 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isLoading">Signing in...</span>
            <span v-else>Sign in</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const username = ref('');
const password = ref('');
const errorMessage = ref('');
const isLoading = ref(false);

const handleLogin = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    await authStore.login({ 
      username: username.value, 
      password: password.value 
    });
    // Redirect to dashboard on successful login
    router.push('/dashboard');

  } catch (error) {
    // Basic error handling. We can make this more specific based on API responses.
    errorMessage.value = 'Login failed. Please check your credentials.';
    console.error(error);
  } finally {
    isLoading.value = false;
  }
};
</script>
