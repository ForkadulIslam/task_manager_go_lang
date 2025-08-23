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

  async function createTask(taskData) {
    // Note: This function assumes taskData is already validated
    try {
      await apiClient.post('/tasks', taskData);
      // Refresh the tasks list to show the new task
      await fetchTasks(); 
    } catch (e) {
      console.error("Failed to create task:", e);
      // Re-throw the error to be handled by the component
      throw new Error('Task creation failed on the server.');
    }
  }

  async function fetchTaskById(id) {
    isLoading.value = true; // Use the same loading state for simplicity
    error.value = null;
    try {
      const response = await apiClient.get(`/tasks/${id}`);
      return response.data.data; // Return the full task object
    } catch (e) {
      error.value = `Failed to fetch task with ID ${id}.`;
      console.error(e);
      throw e; // Re-throw to allow component to handle
    } finally {
      isLoading.value = false;
    }
  }

  return {
    tasks,
    isLoading,
    error,
    fetchTasks,
    createTask,
    fetchTaskById, // Add this
  };
});
