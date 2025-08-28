import { defineStore } from 'pinia';
import { ref } from 'vue';
import apiClient from '../services/api';

export const useMetaStore = defineStore('meta', () => {
  const taskTypes = ref([]);
  const users = ref([]);
  const groups = ref([]);
  const isLoading = ref(false);
  const error = ref(null);

  async function fetchMeta() {
    isLoading.value = true;
    error.value = null;
    try {
      // Fetch all metadata in parallel
      const [typesRes, usersRes, groupsRes] = await Promise.all([
        apiClient.get('/task-types'),
        apiClient.get('/users'), // Note: Endpoint to get all users needs to be confirmed
        apiClient.get('/groups'),
      ]);

      taskTypes.value = typesRes.data.data || [];
      users.value = usersRes.data.data || [];
      groups.value = groupsRes.data.data || [];

    } catch (e) {
      error.value = 'Failed to fetch form metadata.';
      console.error(e);
      // Ensure all refs are arrays even on error
      taskTypes.value = [];
      users.value = [];
      groups.value = [];
    } finally {
      isLoading.value = false;
    }
  }

  async function deleteGroup(groupId) {
    try {
      await apiClient.delete(`/groups/${groupId}`);
      // After successful deletion, refresh the groups list
      await fetchMeta();
    } catch (e) {
      console.error(`Failed to delete group ${groupId}:`, e);
      throw new Error('Group deletion failed on the server.');
    }
  }

  async function createTaskType(taskTypeData) {
    try {
      await apiClient.post('/task-types', taskTypeData);
      await fetchMeta(); // Refresh all meta data, including task types
    } catch (e) {
      console.error("Failed to create task type:", e);
      throw new Error('Task type creation failed on the server.');
    }
  }

  async function updateTaskType(id, taskTypeData) {
    try {
      await apiClient.put(`/task-types/${id}`, taskTypeData);
      await fetchMeta(); // Refresh all meta data, including task types
    } catch (e) {
      console.error(`Failed to update task type ${id}:`, e);
      throw new Error('Task type update failed on the server.');
    }
  }

  async function deleteTaskType(id) {
    try {
      await apiClient.delete(`/task-types/${id}`);
      await fetchMeta(); // Refresh all meta data, including task types
    } catch (e) {
      console.error(`Failed to delete task type ${id}:`, e);
      throw new Error('Task type deletion failed on the server.');
    }
  }

  return {
    taskTypes,
    users,
    groups,
    isLoading,
    error,
    fetchMeta,
    deleteGroup,
    createTaskType,
    updateTaskType,
    deleteTaskType,
  };
});