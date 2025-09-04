<template>
  <div class="min-h-screen bg-gray-900 font-sans text-gray-100 flex">
    <!-- Sidebar -->
    <aside class="w-64 bg-gray-800 border-r border-gray-700 flex-shrink-0">
      <div class="p-4">
        <h1 class="text-2xl font-bold text-white">Task Manager</h1>
        <nav class="mt-10 space-y-2">
          <!-- Navigation links will go here -->
          <RouterLink to="/dashboard" class="block py-2.5 px-4 rounded transition-colors hover:bg-gray-700">Dashboard</RouterLink>
          <RouterLink to="/tasks" class="block py-2.5 px-4 rounded transition-colors hover:bg-gray-700">Generate Tasks</RouterLink>
          <RouterLink to="/my-tasks" class="block py-2.5 px-4 rounded transition-colors hover:bg-gray-700">My Tasks</RouterLink>
          <RouterLink to="/groups" class="block py-2.5 px-4 rounded transition-colors hover:bg-gray-700">Groups</RouterLink>
          <RouterLink to="/task-types" class="block py-2.5 px-4 rounded transition-colors hover:bg-gray-700">Manage Task Types</RouterLink>
        </nav>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col">
      <!-- Header -->
      <header class="bg-gray-800 border-b border-gray-700 p-4 flex justify-end items-center">
        <NotificationIcon class="mr-4" />
        <div>
          <span class="mr-4">Welcome, {{ authStore.user?.username }}!</span>
          <button @click="handleSyncUsers" :disabled="isSyncing" class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded transition-colors mr-4">
            <span v-if="isSyncing">Syncing...</span>
            <span v-else>Sync Users</span>
          </button>
          <button @click="handleLogout" class="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded transition-colors">Logout</button>
        </div>
      </header>

      <!-- Content Area -->
      <main class="flex-1 p-6">
        <RouterView />
      </main>
    </div>
    <ToastContainer /> <!-- New component -->
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import ToastContainer from '../components/ToastContainer.vue'; // New import
import NotificationIcon from '../components/NotificationIcon.vue';
import apiClient from '../services/api';
import { useToastStore } from '../stores/toast';

const authStore = useAuthStore();
const router = useRouter();
const toastStore = useToastStore();

const isSyncing = ref(false);

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};

const handleSyncUsers = async () => {
  isSyncing.value = true;
  try {
    await apiClient.get('/sync-user');
    toastStore.addToast('Users synced successfully!', 'success');
  } catch (error) {
    console.error('Failed to sync users:', error);
    toastStore.addToast('Failed to sync users. Please try again.', 'error');
  } finally {
    isSyncing.value = false;
  }
};
</script>

<style scoped>
.router-link-exact-active {
  @apply bg-sky-500 text-white;
}
</style>
