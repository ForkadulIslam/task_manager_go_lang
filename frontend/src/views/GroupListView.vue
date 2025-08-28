<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-white">Groups Management</h1>
      <button @click="openCreateModal" class="inline-flex items-center justify-center rounded-lg bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-2 focus:ring-offset-gray-900">
        <svg class="-ml-0.5 mr-1.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
          <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
        </svg>
        New Group
      </button>
    </div>

    <!-- Search Input -->
    <div class="mb-4">
      <input
        type="text"
        v-model="searchQuery"
        placeholder="Search groups..."
        class="form-input w-full"
      />
    </div>

    <!-- Loading State -->
    <div v-if="metaStore.isLoading" class="text-center text-gray-400">
      <p>Loading groups...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="metaStore.error" class="bg-red-900 border border-red-700 text-red-200 px-4 py-3 rounded-md">
      <p>{{ metaStore.error }}</p>
    </div>

    <!-- Data Table -->
    <div v-else class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
        <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
          <table class="min-w-full divide-y divide-gray-700">
            <thead class="bg-gray-800">
              <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-6 cursor-pointer" @click="sortBy('label')">
                  Group Name
                  <span v-if="sortColumn === 'label'">{{ sortDirection === 'asc' ? ' ▲' : ' ▼' }}</span>
                </th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white cursor-pointer" @click="sortBy('users.length')">
                  Members
                  <span v-if="sortColumn === 'users.length'">{{ sortDirection === 'asc' ? ' ▲' : ' ▼' }}</span>
                </th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white cursor-pointer" @click="sortBy('created_by')">
                  Creator
                  <span v-if="sortColumn === 'created_by'">{{ sortDirection === 'asc' ? ' ▲' : ' ▼' }}</span>
                </th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white cursor-pointer" @click="sortBy('created_at')">
                  Created At
                  <span v-if="sortColumn === 'created_at'">{{ sortDirection === 'asc' ? ' ▲' : ' ▼' }}</span>
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                  <span class="sr-only">Edit</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-800 bg-gray-900">
              <tr v-for="group in paginatedGroups" :key="group.id" :class="{'highlighted-row': group.id === highlightedGroupId}" class="even:bg-gray-800/30 hover:bg-gray-800/50">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-6">{{ group.label }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-300">
                  {{ displayUsers(group.users) }}
                  <button v-if="group.users && group.users.length > 0" @click="openMembersModal(group)" class="ml-2 text-sky-400 hover:text-sky-300 text-xs font-medium">More</button>
                </td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-400">{{ displayCreator(group.created_by) }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-400">{{ formatDate(group.created_at) }}</td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                  <div class="flex items-center justify-end space-x-2">
                    <button v-if="canEditGroup(group)" @click="openEditModal(group.id)" type="button" class="inline-flex items-center gap-x-1.5 rounded-md bg-blue-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
                      <svg class="-ml-0.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path d="M2.695 14.763l-1.262 3.154a.5.5 0 00.65.65l3.155-1.262a4 4 0 001.343-.885L17.5 5.5a2.121 2.121 0 00-3-3L3.58 13.42a4 4 0 00-.885 1.343z" />
                      </svg>
                    </button>
                    <button v-if="canEditGroup(group)" @click="handleDeleteGroup(group.id)" type="button" class="inline-flex items-center gap-x-1.5 rounded-md bg-red-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600">
                      <svg class="-ml-0.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd" d="M8.75 1A2.75 2.75 0 006 3.75v.443c-.795.077-1.584.176-2.368.298a.75.75 0 10.232 1.482l.175-.027c.572-.089 1.14-.19 1.706-.297H15.25a.75.75 0 000-1.5H6.75A.75.75 0 017.5 4.5h7.75a.75.75 0 000-1.5H8.75zM10 6a.75.75 0 01.75.75v6.5a.75.75 0 01-1.5 0v-6.5A.75.75 0 0110 6zm-3 0a.75.75 0 01.75.75v6.5a.75.75 0 01-1.5 0v-6.5A.75.75 0 017 6zm6 0a.75.75 0 01.75.75v6.5a.75.75 0 01-1.5 0v-6.5A.75.75 0 0113 6z" clip-rule="evenodd" />
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Pagination Controls -->
    <div class="flex justify-between items-center mt-4 text-sm text-gray-400">
      <div>Page {{ currentPage }} of {{ totalPages }}</div>
      <div class="flex space-x-2">
        <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="px-3 py-1 rounded-md bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Previous
        </button>
        <button
          @click="nextPage"
          :disabled="currentPage === totalPages"
          class="px-3 py-1 rounded-md bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next
        </button>
      </div>
    </div>

    
    <Modal :show="showCreateModal" @close="showCreateModal = false">
      <template #header><h2>Create New Group</h2></template>
      <template #body>
        <CreateGroupForm v-if="showCreateModal" @submit="handleCreateGroup" :isLoading="isCreatingGroup" />
        <div v-if="isCreatingGroup" class="text-center text-gray-400 mt-4">
          <p>Creating group...</p>
        </div>
      </template>
    </Modal>

    <Modal :show="showEditModal" @close="showEditModal = false">
      <template #header><h2>Edit Group</h2></template>
      <template #body>
        <EditGroupForm 
          v-if="showEditModal && selectedGroupForEdit"
          :group="selectedGroupForEdit" 
          @submit="handleUpdateGroup" 
          :isLoading="isUpdatingGroup"
        />
        <div v-if="isUpdatingGroup" class="text-center text-gray-400 mt-4">
          <p>Updating group...</p>
        </div>
      </template>
    </Modal>

    <!-- Group Members Modal -->
    <GroupMembersModal
      v-if="showMembersModal"
      :show="showMembersModal"
      :groupData="selectedGroupForMembers"
      @close="showMembersModal = false"
      @membersUpdated="handleMembersUpdated"
    />

    <!-- Assign Users to Group Modal -->
    <AssignUsersToGroupModal
      v-if="showAssignUsersModal"
      :show="showAssignUsersModal"
      :groupId="newlyCreatedGroupId"
      @close="showAssignUsersModal = false; metaStore.fetchMeta(); handleHighlightClear();"
      @usersAssigned="showAssignUsersModal = false; metaStore.fetchMeta(); handleHighlightClear();"
    />
  </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue';
import { useMetaStore } from '../stores/meta';
import Modal from '../components/Modal.vue';
import GroupMembersModal from '../components/GroupMembersModal.vue'; // Import the new modal component
import CreateGroupForm from '../components/CreateGroupForm.vue';
import EditGroupForm from '../components/EditGroupForm.vue';
import apiClient from '../services/api';
import AssignUsersToGroupModal from '../components/AssignUsersToGroupModal.vue';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';

const metaStore = useMetaStore();
const authStore = useAuthStore();
const toastStore = useToastStore();

const handleDeleteGroup = async (groupId) => {
  if (confirm('Are you sure you want to delete this group and all its associated users? This action cannot be undone.')) {
    try {
      await metaStore.deleteGroup(groupId);
      toastStore.addToast('Group deleted successfully!', 'success');
    } catch (error) {
      console.error("Failed to delete group:", error);
      toastStore.addToast(`Failed to delete group: ${error.message}`, 'error');
    }
  }
};

const showCreateModal = ref(false);
const showEditModal = ref(false);
const selectedGroupForEdit = ref(null);
const isUpdatingGroup = ref(false);
const showMembersModal = ref(false); // New state for members modal
const selectedGroupForMembers = ref(null);
const isCreatingGroup = ref(false); // New loading state
const showAssignUsersModal = ref(false); // New state for assign users modal
const newlyCreatedGroupId = ref(null); // To store the ID of the newly created group
const highlightedGroupId = ref(null); // To store the ID of the group to be highlighted

// Pagination state
const currentPage = ref(1);
const itemsPerPage = ref(10);

// Search state
const searchQuery = ref('');

// Sorting state
const sortColumn = ref('created_at'); // Changed to created_at
const sortDirection = ref('desc'); // Changed to desc

const openCreateModal = () => {
  showCreateModal.value = true;
};

const handleCreateGroup = async (formData) => {
  isCreatingGroup.value = true; // Set loading state
  try {
    const response = await apiClient.post('/groups', formData); // Call backend API
    newlyCreatedGroupId.value = response.data.data.id; // Assuming backend returns the ID of the new group
    highlightedGroupId.value = response.data.data.id; // Set the newly created group for highlighting
    showCreateModal.value = false; // Close create group modal
    showAssignUsersModal.value = true; // Open assign users modal
  } catch (error) {
    console.error('Failed to create group:', error);
    toastStore.addToast('Failed to create group. Please try again.', 'error');
  } finally {
    isCreatingGroup.value = false; // Reset loading state
  }
};

const handleUpdateGroup = async (formData) => {
  isUpdatingGroup.value = true;
  try {
    await apiClient.put(`/groups/${formData.id}`, { label: formData.label });
    toastStore.addToast('Group updated successfully!', 'success');
    showEditModal.value = false;
    metaStore.fetchMeta();
  } catch (error) {
    console.error('Failed to update group:', error);
    toastStore.addToast(`Failed to update group: ${error.message}`, 'error');
  } finally {
    isUpdatingGroup.value = false;
  }
};

const openEditModal = (groupId) => {
  const group = metaStore.groups.find(g => g.id === groupId);
  if (group) {
    selectedGroupForEdit.value = { ...group };
    showEditModal.value = true;
  } else {
    console.error('Group not found for editing');
    toastStore.addToast('Could not find group details to edit.', 'error');
  }
};

const openMembersModal = (group) => {
  selectedGroupForMembers.value = group;
  showMembersModal.value = true;
};

const handleMembersUpdated = () => {
  // When members are updated in the modal, re-fetch meta data to refresh the group list
  metaStore.fetchMeta();
  // Optionally, if the modal is still open, update its data as well
  // This might require re-fetching the specific group data if metaStore.groups doesn't update selectedGroupForMembers
  // For simplicity, we just re-fetch all meta data.
};

const formatDate = (dateString) => {
  if (!dateString) return 'N/A';
  const options = { year: 'numeric', month: 'short', day: 'numeric' };
  return new Date(dateString).toLocaleDateString(undefined, options);
};

const displayUsers = (users) => {
  if (!users || users.length === 0) {
    return 'No members';
  }
  const usernames = users.map(user => user.username);
  const maxDisplay = 3; // Display up to 3 usernames
  if (usernames.length <= maxDisplay) {
    return usernames.join(', ');
  } else {
    return usernames.slice(0, maxDisplay).join(', ') + ` (+${usernames.length - maxDisplay} more)`;
  }
};

// Computed property for filtered groups
const filteredGroups = computed(() => {
  let groups = metaStore.groups;
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    groups = groups.filter(group =>
      group.label.toLowerCase().includes(query) ||
      (group.users && group.users.some(user => user.username.toLowerCase().includes(query)))
    );
  }
  return groups;
});

// Computed property for sorted and filtered groups
const sortedAndFilteredGroups = computed(() => {
  const groups = [...filteredGroups.value]; // Create a shallow copy to avoid mutating the original array
  if (sortColumn.value) {
    groups.sort((a, b) => {
      let valA = a[sortColumn.value];
      let valB = b[sortColumn.value];

      // Handle nested properties for sorting (e.g., 'users.length')
      if (sortColumn.value.includes('.')) {
        const path = sortColumn.value.split('.');
        valA = path.reduce((o, i) => (o ? o[i] : undefined), a);
        valB = path.reduce((o, i) => (o ? o[i] : undefined), b);

        // Special handling for 'users.length' to treat null/undefined users as empty array
        if (sortColumn.value === 'users.length') {
          valA = a.users ? a.users.length : 0;
          valB = b.users ? b.users.length : 0;
        }
      }

      if (typeof valA === 'string' && typeof valB === 'string') {
        return sortDirection.value === 'asc' ? valA.localeCompare(valB) : valB.localeCompare(valA);
      } else {
        // For numbers or other comparable types
        if (valA < valB) return sortDirection.value === 'asc' ? -1 : 1;
        if (valA > valB) return sortDirection.value === 'asc' ? 1 : -1;
        return 0;
      }
    });
  }
  return groups;
});

// Computed property for paginated groups
const paginatedGroups = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value;
  const end = start + itemsPerPage.value;
  return sortedAndFilteredGroups.value.slice(start, end);
});

// Computed property for total pages
const totalPages = computed(() => {
  return Math.ceil(filteredGroups.value.length / itemsPerPage.value);
});

// Methods for pagination
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
};

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
};

// Method for sorting
const sortBy = (column) => {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortColumn.value = column;
    sortDirection.value = 'asc';
  }
  currentPage.value = 1; // Reset to first page on sort change
};

const displayCreator = (creatorId) => {
  const creator = metaStore.users.find(user => user.id === creatorId);
  return creator ? creator.username : 'N/A';
};

const handleHighlightClear = () => {
  setTimeout(() => {
    highlightedGroupId.value = null;
  }, 5000);
};

const canEditGroup = (group) => {
  return authStore.user && group && authStore.user.id === group.created_by;
};

onMounted(() => {
  metaStore.fetchMeta();
});
</script>

<style scoped>
.form-input {
  @apply block w-full bg-gray-700 border border-gray-600 rounded-md shadow-sm py-2 px-3 text-sm text-white;
  @apply focus:outline-none focus:ring-sky-700 focus:border-sky-700 transition;
}

.highlighted-row {
  animation: highlight 2s ease-out;
}

@keyframes highlight {
  0% { background-color: #34d399; } /* Tailwind green-400 */
  100% { background-color: transparent; }
}
</style>
