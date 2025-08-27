<template>
  <Modal :show="show" @close="emit('close')">
    <template #header>
      <h2>Assign Members to Group</h2>
    </template>
    <template #body>
      <div class="space-y-4">
        <!-- Add User Section -->
        <div class="min-h-[200px]"> <!-- Added min-h-[200px] -->
          <h3 class="text-lg font-semibold text-white mb-2">Select Members to Add</h3>
          <MultiSelectCombobox
            v-model="selectedUsersToAdd"
            :items="metaStore.users"
            displayProperty="username"
          />
          <button @click="assignUsersToGroup" class="mt-3 w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-sky-600 hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-sky-500 transition-colors">
            Assign Selected Users
          </button>
        </div>
      </div>
    </template>
    <template #footer>
      <button @click="emit('close')" class="px-4 py-2 text-sm font-medium text-gray-300 rounded-md border border-gray-600 hover:bg-gray-700 transition-colors">Skip & Close</button>
    </template>
  </Modal>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import Modal from './Modal.vue';
import MultiSelectCombobox from './MultiSelectCombobox.vue';
import { useMetaStore } from '../stores/meta';
import apiClient from '../services/api';

const props = defineProps({
  show: {
    type: Boolean,
    required: true,
  },
  groupId: {
    type: Number,
    required: true,
  },
});

const emit = defineEmits(['close', 'usersAssigned']);

const metaStore = useMetaStore();
const selectedUsersToAdd = ref([]);

// Filter users already in the group from the general list of users
const availableUsers = computed(() => {
  if (!metaStore.users) return [];
  // For initial assignment, all users are available unless already assigned (which is handled by backend)
  return metaStore.users;
});

onMounted(async () => {
  await metaStore.fetchMeta();
});

const assignUsersToGroup = async () => {
  if (selectedUsersToAdd.value.length === 0) {
    alert('Please select users to assign.');
    return;
  }

  try {
    const payload = {
      user_ids: selectedUsersToAdd.value.map(user => user.id),
      group_id: props.groupId,
    };

    const response = await apiClient.post('/user-groups', payload);

    console.log('User assignment successful, response:', response);

    if (response.status === 200 || response.status === 201) {
      // No alert here anymore
    } else if (response.status === 207) { // StatusMultiStatus
      alert('Some users could not be assigned: ' + response.data.failed_assignments.join(', '));
    }

    console.log('Code reached after assignment logic.');
    selectedUsersToAdd.value = []; // Clear selection
    emit('usersAssigned'); // Notify parent that users were assigned
    emit('close'); // Close the modal
  } catch (error) {
    console.error('Failed to assign users to group:', error);
    toastStore.showToast('Failed to assign users. Please try again.', 'error');
  }
};
</script>
