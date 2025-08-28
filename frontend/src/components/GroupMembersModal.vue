<template>
  <Modal :show="show" @close="emit('close')">
    <template #header>
      <h2>Group Members ({{ groupLabel }})</h2>
    </template>
    <template #body>
      <div class="space-y-4">
        <!-- Current Members List -->
        <div>
          <h3 class="text-lg font-semibold text-white mb-2">Current Members</h3>
          <ul v-if="currentMembers.length > 0" class="space-y-2">
            <li v-for="member in currentMembers" :key="member.id" class="flex items-center justify-between bg-gray-700 p-2 rounded-md">
              <span class="text-gray-200">{{ member.username }}</span>
              <button v-if="isGroupCreator" @click="confirmRemoveUser(member)" class="text-red-400 hover:text-red-300 text-sm">Remove</button>
            </li>
          </ul>
          <p v-else class="text-gray-400">No members in this group.</p>
        </div>

        <!-- Add User Section -->
        <div v-if="isGroupCreator" class="border-t border-gray-700 pt-4">
          <h3 class="text-lg font-semibold text-white mb-2">Add Members</h3>
          <MultiSelectCombobox
            v-model="selectedUsersToAdd"
            :items="availableUsers"
            displayProperty="username"
          />
          <button @click="addUsersToGroup" class="mt-3 w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-sky-600 hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-sky-500 transition-colors">
            Add Selected Users
          </button>
        </div>
      </div>
    </template>
    <template #footer>
      <button @click="emit('close')" class="px-4 py-2 text-sm font-medium text-gray-300 rounded-md border border-gray-600 hover:bg-gray-700 transition-colors">Close</button>
    </template>
  </Modal>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import Modal from './Modal.vue';
import MultiSelectCombobox from './MultiSelectCombobox.vue';
import { useMetaStore } from '../stores/meta';
import apiClient from '../services/api';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast'; // New import

const props = defineProps({
  show: {
    type: Boolean,
    required: true,
  },
  groupData: {
    type: Object,
    required: true,
  },
});

const emit = defineEmits(['close', 'membersUpdated']);

const metaStore = useMetaStore();
const authStore = useAuthStore();
const toastStore = useToastStore(); // New initialization

const isGroupCreator = computed(() => {
  return authStore.user && props.groupData && authStore.user.id === props.groupData.created_by;
});
const currentMembers = ref([]);
const selectedUsersToAdd = ref([]);

const groupLabel = computed(() => props.groupData?.label || 'N/A');

// Filter users already in the group from the general list of users
const availableUsers = computed(() => {
  if (!metaStore.users || !currentMembers.value) return [];
  const currentMemberIds = new Set(currentMembers.value.map(m => m.id)); // Assuming currentMembers already has 'id'
  const filtered = metaStore.users.filter(user => !currentMemberIds.has(user.ID)); // Use user.ID here
  return filtered;
});

// Map availableUsers to have 'id' (lowercase) property for MultiSelectCombobox
const mappedAvailableUsers = computed(() => {
  return availableUsers.value.map(user => ({
    id: user.ID,
    username: user.Username,
    // Copy other properties if needed
    ...user
  }));
});


// Watch for changes in groupData to update currentMembers
watch(() => props.groupData, (newData) => {
  if (newData && newData.users) {
    // newData.users are UserResponse objects from backend, which have 'id' and 'association_id'
    currentMembers.value = newData.users.map(user => ({
      id: user.id, // Use user.id (lowercase)
      username: user.username,
      association_id: user.association_id, // Use user.association_id (lowercase)
      // Copy other properties if needed
      ...user
    }));
  } else {
    currentMembers.value = [];
  }
}, { immediate: true, deep: true });

// Fetch meta data when the component is mounted (or re-mounted if modal is toggled)
onMounted(async () => {
  await metaStore.fetchMeta();
});

const addUsersToGroup = async () => {
  if (selectedUsersToAdd.value.length === 0) {
    toastStore.addToast('Please select users to add.', 'info');
    return;
  }

  try {
    const payload = {
      user_ids: selectedUsersToAdd.value.map(user => user.id),
      group_id: props.groupData.id,
    };

    const response = await apiClient.post('/user-groups', payload);

    if (response.status === 200 || response.status === 201) {
      toastStore.addToast('Users added successfully!', 'success');
      // Optimistically add the users to the local list
      selectedUsersToAdd.value.forEach(user => {
        // Ensure the user object has the necessary properties (id, username)
        // The user object from selectedUsersToAdd should already have these
        currentMembers.value.push({ id: user.id, username: user.username });
      });
    } else if (response.status === 207) { // StatusMultiStatus
      toastStore.addToast('Some users could not be added: ' + response.data.failed_assignments.join(', '), 'warning');
    }

    selectedUsersToAdd.value = []; // Clear selection
    emit('membersUpdated'); // Notify parent to refresh group data
  } catch (error) {
    console.error('Failed to add users to group:', error);
    toastStore.addToast('Failed to add users. Please try again.', 'error');
  }
};

const confirmRemoveUser = async (user) => {
  if (confirm(`Are you sure you want to remove ${user.username} from this group?`)) {
    try {
      if (!user.association_id) {
        toastStore.addToast('Error: User association ID not found. Cannot remove user.', 'error');
        console.error('Attempted to remove user without association_id:', user);
        return;
      }
      await apiClient.delete(`/user-groups/${user.association_id}`);
      toastStore.addToast('User removed successfully!', 'success');
      // Optimistically remove the user from the local list
      currentMembers.value = currentMembers.value.filter(member => member.id !== user.id);
      emit('membersUpdated'); // Notify parent to refresh group data
    } catch (error) {
      console.error('Failed to remove user from group:', error);
      toastStore.addToast('Failed to remove user. Please try again.', 'error');
    }
  }
};

</script>
