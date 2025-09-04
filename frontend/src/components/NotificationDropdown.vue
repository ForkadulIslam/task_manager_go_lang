<template>
  <div class="absolute right-0 z-20 w-64 mt-2 overflow-hidden bg-white rounded-md shadow-lg sm:w-80 dark:bg-gray-800">
    <div class="py-2">
      <div v-if="notifications.length === 0" class="px-4 py-2 text-sm text-gray-700 dark:text-gray-200">No new notifications</div>
      <div v-for="notification in notifications" :key="notification.id" @click="handleNotificationClick(notification)" class="flex items-center px-4 py-3 -mx-2 border-b hover:bg-gray-100 dark:hover:bg-gray-700 dark:border-gray-700 cursor-pointer">
        <div class="mx-3">
          <p class="text-sm font-medium text-gray-800 dark:text-white" :class="{ 'font-bold': !notification.is_read }">{{ notification.message }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ new Date(notification.created_at).toLocaleString() }}</p>
        </div>
      </div>
    </div>
    <a href="#" @click.prevent="handleMarkAllAsRead" class="block py-2 font-bold text-center text-white bg-gray-800 dark:bg-gray-700 hover:underline">Mark all as read</a>
  </div>
</template>

<script setup>
import { useNotificationStore } from '@/stores/notification';
import { storeToRefs } from 'pinia';

const store = useNotificationStore();
const { notifications } = storeToRefs(store);
const { markAsRead, markAllAsRead } = store;

const handleNotificationClick = async (notification) => {
  await markAsRead(notification.id);
  emit('close');
};

const handleMarkAllAsRead = async () => {
  await markAllAsRead();
};

const emit = defineEmits(['close']);
</script>
