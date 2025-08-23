import { defineStore } from 'pinia';
import { ref } from 'vue';
import apiClient from '../services/api';

export const useTasksStore = defineStore('tasks', () => {
  const tasks = ref([]);
  const isLoading = ref(false);
  const error = ref(null);

  async function fetchTasks() {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await apiClient.get('/tasks');
      tasks.value = response.data.data; // Adjust based on your API response structure
    } catch (e) {
      error.value = 'Failed to fetch tasks.';
      console.error(e);
    } finally {
      isLoading.value = false;
    }
  }

  return {
    tasks,
    isLoading,
    error,
    fetchTasks,
  };
});
