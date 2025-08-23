<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-white">Tasks</h1>
      <button @click="openCreateModal" class="bg-sky-500 hover:bg-sky-600 text-white font-bold py-2 px-4 rounded-lg transition-colors">
        + New Task
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
    <div v-else class="bg-gray-800 rounded-lg shadow-lg overflow-hidden">
      <table class="min-w-full text-left">
        <thead class="bg-gray-700">
          <tr>
            <th class="p-4 font-semibold">Task</th>
            <th class="p-4 font-semibold">Priority</th>
            <th class="p-4 font-semibold">Status</th>
            <th class="p-4 font-semibold">Due Date</th>
            <th class="p-4 font-semibold">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="task in tasksStore.tasks" :key="task.ID" class="border-b border-gray-700 hover:bg-gray-700/50 transition-colors">
            <td class="p-4">{{ task.Label }}</td>
            <td class="p-4">
              <span :class="getPriorityClass(task.Priority)" class="px-2 py-1 text-xs font-semibold rounded-full">
                {{ task.Priority }}
              </span>
            </td>
            <td class="p-4">
              <span :class="getStatusClass(task.Status)" class="px-2 py-1 text-xs font-semibold rounded-full">
                {{ task.Status }}
              </span>
            </td>
            <td class="p-4 text-gray-400">{{ formatDate(task.DueDate) }}</td>
            <td class="p-4 space-x-2">
              <button @click="openViewModal(task.ID)" class="bg-gray-600 hover:bg-gray-500 text-white py-1 px-3 rounded-lg text-sm transition-colors">View</button>
              <button @click="openUpdateModal(task)" class="bg-blue-600 hover:bg-blue-500 text-white py-1 px-3 rounded-lg text-sm transition-colors">Edit</button>
            </td>
          </tr>
        </tbody>
      </table>
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

    <Modal :show="showUpdateModal" @close="showUpdateModal = false">
      <template #header><h2>Update Task</h2></template>
      <template #body><p>Form for updating task "{{ selectedTask?.Label }}" will go here.</p></template>
    </Modal>

  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useTasksStore } from '../stores/tasks';
import Modal from '../components/Modal.vue';
import CreateTaskForm from '../components/CreateTaskForm.vue';
import TaskDetailsView from '../components/TaskDetailsView.vue';

const tasksStore = useTasksStore();

// Modal visibility state
const showCreateModal = ref(false);
const showViewModal = ref(false);
const showUpdateModal = ref(false);
const selectedTask = ref(null);

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

const openUpdateModal = (task) => {
  selectedTask.value = task;
  showUpdateModal.value = true;
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
