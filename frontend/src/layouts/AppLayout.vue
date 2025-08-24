<template>
  <div class="min-h-screen bg-gray-900 font-sans text-gray-100 flex">
    <!-- Sidebar -->
    <aside class="w-64 bg-gray-800 border-r border-gray-700 flex-shrink-0">
      <div class="p-4">
        <h1 class="text-2xl font-bold text-white">Task Manager</h1>
        <nav class="mt-10 space-y-2">
          <!-- Navigation links will go here -->
          <RouterLink to="/dashboard" class="block py-2.5 px-4 rounded transition-colors hover:bg-gray-700">Dashboard</RouterLink>
          <RouterLink to="/tasks" class="block py-2.5 px-4 rounded transition-colors hover:bg-gray-700">Tasks</RouterLink>
          <a href="#" class="block py-2.5 px-4 rounded transition-colors hover:bg-gray-700">Groups</a>
        </nav>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col">
      <!-- Header -->
      <header class="bg-gray-800 border-b border-gray-700 p-4 flex justify-end">
        <div>
          <span class="mr-4">Welcome, {{ authStore.user?.username }}!</span>
          <button @click="handleLogout" class="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded transition-colors">Logout</button>
        </div>
      </header>

      <!-- Content Area -->
      <main class="flex-1 p-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};
</script>

<style scoped>
.router-link-exact-active {
  @apply bg-sky-500 text-white;
}
</style>
