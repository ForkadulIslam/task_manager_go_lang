<template>
  <div class="relative">
    <button @click="toggleDropdown" class="relative z-10 block p-2 text-gray-700 bg-white border border-transparent rounded-md dark:text-white focus:border-blue-500 focus:ring-blue-500 dark:focus:ring-opacity-40 dark:focus:ring-blue-300 focus:ring-opacity-40 dark:bg-gray-800 focus:outline-none">
      <svg class="w-5 h-5 text-gray-800 dark:text-white" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
        <path d="M19 13.586V10c0-3.217-2.185-5.927-5.145-6.742C13.562 2.52 12.846 2 12 2s-1.562.52-1.855 1.258C7.185 4.073 5 6.783 5 10v3.586l-1.707 1.707A.996.996 0 0 0 3 16v2a1 1 0 0 0 1 1h16a1 1 0 0 0 1-1v-2a.996.996 0 0 0-.293-.707L19 13.586zM19 17H5v-1l1.293-1.293A.996.996 0 0 0 7 14v-4c0-2.757 2.243-5 5-5s5 2.243 5 5v4c0 .266.105.52.293.707L19 16v1z" />
      </svg>
      <span v-if="unreadCount > 0" class="absolute top-0 right-0 inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-red-100 transform translate-x-1/2 -translate-y-1/2 bg-red-600 rounded-full">
        {{ unreadCount }}
      </span>
    </button>
    <NotificationDropdown v-if="isDropdownOpen" @close="closeDropdown" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useNotificationStore } from '@/stores/notification';
import { storeToRefs } from 'pinia';
import NotificationDropdown from './NotificationDropdown.vue';

const store = useNotificationStore();
const { unreadCount } = storeToRefs(store);
const { fetchNotifications } = store;

const isDropdownOpen = ref(false);

const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value;
  if (isDropdownOpen.value) {
    fetchNotifications();
  }
};

const closeDropdown = () => {
  isDropdownOpen.value = false;
};

onMounted(() => {
  fetchNotifications();
  setInterval(fetchNotifications, 60000); // Refresh every minute
});
</script>
