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

      taskTypes.value = typesRes.data.data;
      users.value = usersRes.data.data;
      groups.value = groupsRes.data.data;

    } catch (e) {
      error.value = 'Failed to fetch form metadata.';
      console.error(e);
      // Handle cases where one or more requests fail
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

  return {
    taskTypes,
    users,
    groups,
    isLoading,
    error,
    fetchMeta,
    deleteGroup,
  };
});
