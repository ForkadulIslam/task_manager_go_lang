<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-white">My Tasks</h1>
      <button @click="toggleFilters" class="inline-flex items-center justify-center rounded-lg bg-gray-700 px-4 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 focus:ring-offset-gray-900">
        <svg class="-ml-0.5 mr-1.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
          <path fill-rule="evenodd" d="M3 3a1 1 0 011-1h12a1 1 0 011 1v3a1 1 0 01-.293.707L12 11.414V15a1 1 0 01-.293.707l-2 2A1 1 0 018 17v-5.586L3.293 6.707A1 1 0 013 6V3z" clip-rule="evenodd" />
        </svg>
        {{ showFilters ? 'Hide Filters' : 'Show Filters' }}
      </button>
    </div>

    <!-- Filter Section -->
    <div v-if="showFilters" class="bg-gray-900 p-4 rounded-lg shadow-xl mb-4">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-2 mb-4 items-end">
        <div>
          <label for="fromDate" class="block text-sm font-medium text-gray-300 mb-1">From Date</label>
          <input type="date" id="fromDate" v-model="filter.fromDate" class="form-input">
        </div>
        <div>
          <label for="toDate" class="block text-sm font-medium text-gray-300 mb-1">To Date</label>
          <input type="date" id="toDate" v-model="filter.toDate" class="form-input">
        </div>
        <div>
          <label for="status" class="block text-sm font-medium text-gray-300 mb-1">Status</label>
          <select id="status" v-model="filter.status" class="form-input">
            <option value="">All Statuses</option>
            <option>Pending</option>
            <option>In Progress</option>
            <option>In Review</option>
            <option>Completed</option>
          </select>
        </div>
        <div>
          <label for="taskType" class="block text-sm font-medium text-gray-300 mb-1">Task Type</label>
          <select id="taskType" v-model="filter.taskTypeId" class="form-input">
            <option value="">All Types</option>
            <option v-for="type in metaStore.taskTypes" :key="type.ID" :value="type.ID">{{ type.Label }}</option>
          </select>
        </div>
        <div class="self-end">
          <button @click="applyFilters" class="w-full inline-flex items-center justify-center rounded-lg bg-sky-600 px-6 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-2 focus:ring-offset-gray-900">
            Apply Filters
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="tasksStore.isLoading" class="text-center text-gray-400">
      <p>Loading tasks...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="tasksStore.error" class="bg-red-900 border border-red-700 text-red-200 px-4 py-3 rounded-md">
      <p>{{ tasksStore.error }}</p>
    </div>

    <!-- Data Table -->
    <div v-else class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
        <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
          <table class="min-w-full divide-y divide-gray-700">
            <thead class="bg-gray-800">
              <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-6">Task</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">Priority</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">Status</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">Due Date</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">Creator</th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                  <span class="sr-only">View</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-800 bg-gray-900">
              <tr v-for="task in tasksStore.tasks" :key="task.ID" class="hover:bg-gray-800/50">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-6">{{ task.Label }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-300">
                  <span :class="getPriorityClass(task.Priority)" class="px-2 py-1 text-xs font-semibold rounded-full">
                    {{ task.Priority }}
                  </span>
                </td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-300">
                  <span :class="getStatusClass(task.Status)" class="px-2 py-1 text-xs font-semibold rounded-full">
                    {{ task.Status }}
                  </span>
                </td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-400">{{ formatDate(task.DueDate) }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-400">{{ task.Creator ? task.Creator.username : 'N/A' }}</td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                  <div class="flex items-center justify-end space-x-2">
                    <button @click="openViewModal(task.ID)" type="button" class="inline-flex items-center gap-x-1.5 rounded-md bg-gray-700 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-gray-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-600">
                      <svg class="-ml-0.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path d="M10 12.5a2.5 2.5 0 100-5 2.5 2.5 0 000 5z" />
                        <path fill-rule="evenodd" d="M.664 10.59a1.651 1.651 0 010-1.18l.879-1.148a1.65 1.65 0 011.956-.543l1.223.345a1.65 1.65 0 011.215 1.088l.28.84a1.65 1.65 0 01-1.67 1.977l-1.257-.355a1.65 1.65 0 01-1.552-1.223l-.212-.635zM19.336 10.59a1.651 1.651 0 010-1.18l-.879-1.148a1.65 1.65 0 01-1.956-.543l-1.223.345a1.65 1.65 0 01-1.215 1.088l-.28.84a1.65 1.65 0 011.67 1.977l1.257-.355a1.65 1.65 0 011.552-1.223l.212-.635zM10 18a8 8 0 100-16 8 8 0 000 16z" clip-rule="evenodd" />
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

    <!-- Modals -->
    <Modal :show="showViewModal" @close="showViewModal = false">
      <template #header><h2>Task Details</h2></template>
      <template #body>
        <TaskDetailsView :task="selectedTask" @commentSubmitted="handleCommentSubmitted" />
      </template>
    </Modal>

  </div>
</template>

<script setup>
import { onMounted, ref, reactive } from 'vue';
import { useTasksStore } from '../stores/tasks';
import { useMetaStore } from '../stores/meta';
import { useToastStore } from '../stores/toast'; // New import
import Modal from '../components/Modal.vue';
import TaskDetailsView from '../components/TaskDetailsView.vue';

const tasksStore = useTasksStore();
const toastStore = useToastStore(); // New initialization

// Modal visibility state
const showViewModal = ref(false);
const selectedTask = ref(null);

const showFilters = ref(false); // Default to hidden

const toggleFilters = () => {
  showFilters.value = !showFilters.value;
};

// Filter state
const filter = reactive({
  fromDate: '',
  toDate: '',
  status: '',
  taskTypeId: '',
});

const metaStore = useMetaStore(); // Import metaStore

const applyFilters = () => {
  const filterData = {};
  if (filter.fromDate) filterData.from_date = filter.fromDate;
  if (filter.toDate) filterData.to_date = filter.toDate;
  if (filter.status) filterData.status = filter.status;
  if (filter.taskTypeId) filterData.task_type_id = filter.taskTypeId;

  tasksStore.fetchMyTasks(filterData);
};

const openViewModal = async (taskId) => {
  const taskFromList = tasksStore.tasks.find(t => t.ID === taskId);
  if (taskFromList) {
    selectedTask.value = taskFromList; // Set to partial task from list
  } else {
    selectedTask.value = { ID: taskId, Label: 'Loading...' }; // Or a minimal placeholder
  }

  showViewModal.value = true; // Open modal

  try {
    const fullTask = await tasksStore.fetchTaskById(taskId);
    selectedTask.value = fullTask; // Update with full task
  } catch (error) {
    console.error("Failed to fetch full task details:", error);
    toastStore.addToast('Failed to load task details. Please try again.', 'error'); // Changed to addToast
    showViewModal.value = false; // Close modal on error
  }
};

const handleCommentSubmitted = async () => {
  // Re-fetch the full task details to update comments
  if (selectedTask.value && selectedTask.value.ID) {
    try {
      const fullTask = await tasksStore.fetchTaskById(selectedTask.value.ID);
      selectedTask.value = fullTask;
    } catch (error) {
      console.error("Failed to re-fetch task details after comment submission:", error);
      // Handle error, maybe keep the old comments or show a message
    }
  }
};

onMounted(async () => {
  await metaStore.fetchMeta(); // Fetch meta data for task types
  tasksStore.fetchMyTasks(); // Initial fetch without filters
});

const getPriorityClass = (priority) => {
  switch (priority) {
    case 'High': return 'bg-red-500/20 text-red-300';
    case 'Escalation': return 'bg-red-700/40 text-red-200';
    case 'Medium': return 'bg-yellow-500/20 text-yellow-300';
    case 'Normal':
    default:
      return 'bg-sky-500/20 text-sky-300';
  }
};

const getStatusClass = (status) => {
  switch (status) {
    case 'Completed': return 'bg-green-500/20 text-green-300';
    case 'In Progress': return 'bg-blue-500/20 text-blue-300';
    case 'In Review': return 'bg-purple-500/20 text-purple-300';
    case 'Pending':
    default:
      return 'bg-gray-500/20 text-gray-300';
  }
};

const formatDate = (dateString) => {
  if (!dateString) return 'N/A';
  const options = { year: 'numeric', month: 'short', day: 'numeric' };
  return new Date(dateString).toLocaleDateString(undefined, options);
};
</script>

<style scoped>
.form-input {
  @apply block w-full bg-gray-700 border border-gray-600 rounded-md shadow-sm py-2 px-3 text-sm text-white;
  @apply focus:outline-none focus:ring-sky-700 focus:border-sky-700 transition;
}
</style>