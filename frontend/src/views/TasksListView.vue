<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-white">Generate Tasks</h1>
      <button @click="openCreateModal" class="inline-flex items-center justify-center rounded-lg bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-2 focus:ring-offset-gray-900">
        <svg class="-ml-0.5 mr-1.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
          <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
        </svg>
        New Task
      </button>
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
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                  <span class="sr-only">Edit</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-800 bg-gray-900">
              <tr v-for="task in tasksStore.tasks" :key="task.ID" class="even:bg-gray-800/30 hover:bg-gray-800/50">
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
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                  <div class="flex items-center justify-end space-x-2">
                    <button @click="openViewModal(task.ID)" type="button" class="inline-flex items-center gap-x-1.5 rounded-md bg-gray-700 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-gray-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-600">
                      <svg class="-ml-0.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path d="M10 12.5a2.5 2.5 0 100-5 2.5 2.5 0 000 5z" />
                        <path fill-rule="evenodd" d="M.664 10.59a1.651 1.651 0 010-1.18l.879-1.148a1.65 1.65 0 011.956-.543l1.223.345a1.65 1.65 0 011.215 1.088l.28.84a1.65 1.65 0 01-1.67 1.977l-1.257-.355a1.65 1.65 0 01-1.552-1.223l-.212-.635zM19.336 10.59a1.651 1.651 0 010-1.18l-.879-1.148a1.65 1.65 0 01-1.956-.543l-1.223.345a1.65 1.65 0 01-1.215 1.088l-.28.84a1.65 1.65 0 011.67 1.977l1.257-.355a1.65 1.65 0 011.552-1.223l.212-.635zM10 18a8 8 0 100-16 8 8 0 000 16z" clip-rule="evenodd" />
                      </svg>
                    </button>
                    <button @click="openEditModal(task.ID)" type="button" class="inline-flex items-center gap-x-1.5 rounded-md bg-blue-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
                      <svg class="-ml-0.5 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path d="M2.695 14.763l-1.262 3.154a.5.5 0 00.65.65l3.155-1.262a4 4 0 001.343-.885L17.5 5.5a2.121 2.121 0 00-3-3L3.58 13.42a4 4 0 00-.885 1.343z" />
                      </svg>
                    </button>
                    <button @click="handleDeleteTask(task.ID)" type="button" class="inline-flex items-center gap-x-1.5 rounded-md bg-red-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600">
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

    <!-- Modals -->
    <Modal :show="showCreateModal" @close="showCreateModal = false">
      <template #header><h2>Create New Task</h2></template>
      <template #body>
        <CreateTaskForm @submit="handleCreateTask" />
      </template>
    </Modal>

    <Modal :show="showViewModal" @close="showViewModal = false">
      <template #header><h2>Task Details</h2></template>
      <template #body>
        <TaskDetailsView :task="selectedTask" @commentSubmitted="handleCommentSubmitted" />
      </template>
    </Modal>

    <Modal :show="showEditModal" @close="showEditModal = false">
      <template #header><h2>Edit Task</h2></template>
      <template #body v-if="selectedTask">
        <EditTaskForm 
          :initialTaskData="selectedTask" 
          @submit="handleUpdateTask" 
        />
      </template>
    </Modal>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useTasksStore } from '../stores/tasks';
import Modal from '../components/Modal.vue';
import CreateTaskForm from '../components/CreateTaskForm.vue';
import EditTaskForm from '../components/EditTaskForm.vue';
import TaskDetailsView from '../components/TaskDetailsView.vue';
import { useToastStore } from '../stores/toast';

const tasksStore = useTasksStore();
const toastStore = useToastStore();

// Modal visibility state
const showCreateModal = ref(false);
const showViewModal = ref(false);
const showEditModal = ref(false);
const selectedTask = ref(null);

const handleDeleteTask = async (taskId) => {
  if (confirm('Are you sure you want to delete this task? This action cannot be undone.')) {
    try {
      await tasksStore.deleteTask(taskId);
      toastStore.addToast('Task deleted successfully!', 'success');
    } catch (error) {
      console.error("Failed to delete task:", error);
      toastStore.addToast(`Failed to delete task: ${error.message}`, 'error');
    }
  }
};



const openCreateModal = () => {
  showCreateModal.value = true;
};

const handleCreateTask = async (formData) => {
  try {
    await tasksStore.createTask(formData);
    showCreateModal.value = false; // Close modal on success
  } catch (error) {
    console.error("Failed to create task from view:", error);
    // Optionally, show an error message to the user
    alert(error.message);
  }
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
    alert("Failed to load task details. Please try again.");
    showViewModal.value = false; // Close modal on error
  }
};

const openEditModal = async (taskId) => {
  try {
    const task = await tasksStore.fetchTaskById(taskId);
    selectedTask.value = task;
    showEditModal.value = true;
  } catch (error) {
    console.error("Failed to fetch task for editing:", error);
    alert("Failed to load task for editing.");
  }
};

const handleUpdateTask = async (formData) => {
  if (!selectedTask.value) return;
  try {
    await tasksStore.updateTask(selectedTask.value.ID, formData);
    showEditModal.value = false;
    selectedTask.value = null;
  } catch (error) {
    console.error("Failed to update task:", error);
    alert(error.message);
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

onMounted(() => {
  tasksStore.fetchTasks();
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