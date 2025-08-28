<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-white">Manage Task Types</h1>
      <button @click="openCreateModal" class="inline-flex items-center justify-center rounded-lg bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-2 focus:ring-offset-gray-900">
        <svg class="-ml-0.5 mr-1.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
          <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
        </svg>
        New Task Type
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="metaStore.isLoading" class="text-center text-gray-400">
      <p>Loading task types...</p>
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
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-6">Label</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">Created At</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">Updated At</th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                  <span class="sr-only">Actions</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-800 bg-gray-900">
              <tr v-for="taskType in metaStore.taskTypes" :key="taskType.ID" class="even:bg-gray-800/30 hover:bg-gray-800/50">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-6">{{ taskType.Label }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-400">{{ formatDate(taskType.CreatedAt) }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-400">{{ formatDate(taskType.UpdatedAt) }}</td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                  <div class="flex items-center justify-end space-x-2">
                    <!-- Edit Button -->
                    <button @click="openEditModal(taskType)" type="button" class="inline-flex items-center gap-x-1.5 rounded-md bg-blue-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
                      <svg class="-ml-0.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path d="M2.695 14.763l-1.262 3.154a.5.5 0 00.65.65l3.155-1.262a4 4 0 001.343-.885L17.5 5.5a2.121 2.121 0 00-3-3L3.58 13.42a4 4 0 00-.885 1.343z" />
                      </svg>
                    </button>
                    <!-- Delete Button -->
                    <button @click="handleDeleteTaskType(taskType.ID)" type="button" class="inline-flex items-center gap-x-1.5 rounded-md bg-red-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600">
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
  </div>

  <!-- Modals -->
  <Modal :show="showCreateModal" @close="showCreateModal = false">
    <template #header><h2>Create New Task Type</h2></template>
    <template #body>
      <CreateTaskTypeForm @submit="handleCreateTaskType" :isLoading="isCreatingTaskType" />
    </template>
  </Modal>

  <!-- Edit Modal -->
  <Modal :show="showEditModal" @close="showEditModal = false">
    <template #header><h2>Edit Task Type</h2></template>
    <template #body v-if="selectedTaskType">
      <EditTaskTypeForm
        :initialLabel="selectedTaskType.Label"
        @submit="handleUpdateTaskType"
        :isLoading="isUpdatingTaskType"
      />
    </template>
  </Modal>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useMetaStore } from '../stores/meta';
import { useToastStore } from '../stores/toast';
import Modal from '../components/Modal.vue';
import CreateTaskTypeForm from '../components/CreateTaskTypeForm.vue';
import EditTaskTypeForm from '../components/EditTaskTypeForm.vue'; // New import

const metaStore = useMetaStore();
const toastStore = useToastStore();

const showCreateModal = ref(false);
const isCreatingTaskType = ref(false);

const showEditModal = ref(false); // New ref
const selectedTaskType = ref(null); // New ref
const isUpdatingTaskType = ref(false); // New ref

const openCreateModal = () => {
  showCreateModal.value = true;
};

const handleCreateTaskType = async (formData) => {
  isCreatingTaskType.value = true;
  try {
    await metaStore.createTaskType(formData);
    toastStore.addToast('Task type created successfully!', 'success');
    showCreateModal.value = false;
  } catch (error) {
    console.error("Failed to create task type:", error);
    toastStore.addToast(`Failed to create task type: ${error.message}`, 'error');
  } finally {
    isCreatingTaskType.value = false;
  }
};

const openEditModal = (taskType) => {
  selectedTaskType.value = { ...taskType }; // Create a copy to avoid direct mutation
  showEditModal.value = true;
};

const handleUpdateTaskType = async (formData) => {
  if (!selectedTaskType.value) return;
  isUpdatingTaskType.value = true;
  try {
    await metaStore.updateTaskType(selectedTaskType.value.ID, formData);
    toastStore.addToast('Task type updated successfully!', 'success');
    showEditModal.value = false;
    selectedTaskType.value = null;
  } catch (error) {
    console.error("Failed to update task type:", error);
    toastStore.addToast(`Failed to update task type: ${error.message}`, 'error');
  } finally {
    isUpdatingTaskType.value = false;
  }
};

const handleDeleteTaskType = async (id) => {
  if (confirm('Are you sure you want to delete this task type? This action cannot be undone.')) {
    try {
      await metaStore.deleteTaskType(id);
      toastStore.addToast('Task type deleted successfully!', 'success');
    } catch (error) {
      console.error("Failed to delete task type:", error);
      toastStore.addToast(`Failed to delete task type: ${error.message}`, 'error');
    }
  }
};

onMounted(() => {
  // metaStore.fetchMeta() already fetches task types, users, and groups.
  // We only need task types here, but fetching all is fine for now.
  // If performance becomes an issue, a dedicated fetchTaskTypes could be added to metaStore.
  metaStore.fetchMeta();
});

const formatDate = (dateString) => {
  if (!dateString) return 'N/A';
  const options = { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' };
  return new Date(dateString).toLocaleDateString(undefined, options);
};
</script>

<style scoped>
/* Add any specific styles for this view here */
</style>
