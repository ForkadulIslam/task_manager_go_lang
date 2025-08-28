<template>
  <form @submit.prevent="handleSubmit">
    <div class="space-y-5">
      <!-- Task Title -->
      <div>
        <label for="label" class="block text-sm font-medium text-gray-300 mb-1">Task Title</label>
        <input v-model="form.label" type="text" id="label" class="form-input" :class="{ 'border-red-500': v$.label.$error }">
        <div v-if="v$.label.$error" class="text-red-400 text-xs mt-1">{{ v$.label.$errors[0].$message }}</div>
      </div>

      <!-- Horizontal Layout: Task Type & Priority -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="taskTypeId" class="block text-sm font-medium text-gray-300 mb-1">Task Type</label>
          <select v-model="form.taskTypeId" id="taskTypeId" class="form-input" :class="{ 'border-red-500': v$.taskTypeId.$error }">
            <option disabled value="">Please select one</option>
            <option v-for="type in metaStore.taskTypes" :key="type.ID" :value="type.ID">{{ type.Label }}</option>
          </select>
          <div v-if="v$.taskTypeId.$error" class="text-red-400 text-xs mt-1">{{ v$.taskTypeId.$errors[0].$message }}</div>
        </div>
        <div>
          <label for="priority" class="block text-sm font-medium text-gray-300 mb-1">Priority</label>
          <select v-model="form.priority" id="priority" class="form-input" :class="{ 'border-red-500': v$.priority.$error }">
            <option>Normal</option>
            <option>Medium</option>
            <option>High</option>
            <option>Escalation</option>
          </select>
          <div v-if="v$.priority.$error" class="text-red-400 text-xs mt-1">{{ v$.priority.$errors[0].$message }}</div>
        </div>
      </div>

      <!-- Horizontal Layout: Dates -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="startDate" class="block text-sm font-medium text-gray-300 mb-1">Start Date</label>
          <input v-model="form.startDate" type="date" id="startDate" class="form-input" :class="{ 'border-red-500': v$.startDate.$error }">
          <div v-if="v$.startDate.$error" class="text-red-400 text-xs mt-1">{{ v$.startDate.$errors[0].$message }}</div>
        </div>
        <div>
          <label for="dueDate" class="block text-sm font-medium text-gray-300 mb-1">Due Date</label>
          <input v-model="form.dueDate" type="date" id="dueDate" class="form-input" :class="{ 'border-red-500': v$.dueDate.$error }">
          <div v-if="v$.dueDate.$error" class="text-red-400 text-xs mt-1">{{ v$.dueDate.$errors[0].$message }}</div>
        </div>
      </div>

      <!-- Horizontal Layout: Assignees -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-1">Assign To Users</label>
          <MultiSelectCombobox v-model="form.assignedToUsers" :items="metaStore.users" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-1">Assign To Groups</label>
          <MultiSelectCombobox v-model="form.assignedToGroups" :items="metaStore.groups" displayProperty="label" />
        </div>
      </div>

      <!-- Follow-up Users & Groups -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-1">Follow-up Users</label>
          <MultiSelectCombobox v-model="form.followUpUsers" :items="metaStore.users" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-1">Follow-up Groups</label>
          <MultiSelectCombobox v-model="form.followUpGroups" :items="metaStore.groups" displayProperty="label" />
        </div>
      </div>

      <!-- Description -->
      <div>
        <label for="description" class="block text-sm font-medium text-gray-300 mb-1">Description</label>
        <textarea v-model="form.description" id="description" rows="3" class="form-input"></textarea>
      </div>


      <div>
        <label for="attachment" class="block text-sm font-medium text-gray-300 mb-1">Attachment</label>
        <div class="flex items-center space-x-2">
          <input type="file" id="attachment" @change="handleFileChange" :disabled="isUploading" class="form-input file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-sky-50 file:text-sky-700 hover:file:bg-sky-100" />
          <span v-if="isUploading" class="text-gray-400 text-sm">Uploading...</span>
          <span v-else-if="form.attachmentPath" class="text-green-400 text-sm">Uploaded: {{ form.attachmentPath.split('/').pop() }}</span>
        </div>
      </div>

      <!-- Submit Button -->
      <div class="pt-2">
        <button type="submit" class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-sky-600 hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-sky-500 transition-colors">
          Create Task
        </button>
      </div>
    </div>
  </form>
</template>

<script setup>
import { reactive, onMounted, ref } from 'vue';
import { useMetaStore } from '../stores/meta';
import { useVuelidate } from '@vuelidate/core';
import { required, minLength, helpers } from '@vuelidate/validators';
import MultiSelectCombobox from './MultiSelectCombobox.vue';
import apiClient from '../services/api'; // Import apiClient
import { useToastStore } from '../stores/toast'; // New import

const emit = defineEmits(['submit']);

const metaStore = useMetaStore();
const toastStore = useToastStore(); // New initialization

const isUploading = ref(false); // New ref declaration

const defaultFormState = {
  label: '',
  taskTypeId: '',
  priority: 'Normal',
  startDate: '',
  dueDate: '',
  description: '',
  assignedToUsers: [],
  assignedToGroups: [],
  followUpUsers: [],
  followUpGroups: [],
  attachment: null, // This will store the file object temporarily
  attachmentPath: '', // This will store the path returned by the backend
};

const form = reactive({ ...defaultFormState });

const rules = {
  label: { required, minLength: minLength(3) },
  taskTypeId: { required: helpers.withMessage('Task type is required', required) },
  priority: { required },
  startDate: { required: helpers.withMessage('Start date is required', required) },
  dueDate: { 
    // Due date is optional, but if provided, must be >= start date
    // Custom validator for gtefield
    isAfterStartDate: helpers.withMessage(
      'Due date must be on or after start date',
      (value) => !value || !form.startDate || new Date(value) >= new Date(form.startDate)
    ),
  },
};

const v$ = useVuelidate(rules, form);

onMounted(async () => {
  await metaStore.fetchMeta(); // Ensure meta data is fetched
});

const uploadAttachment = async (file) => {
  isUploading.value = true;
  try {
    const formData = new FormData();
    formData.append('attachment', file);

    const response = await apiClient.post('/upload-attachment', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    form.attachmentPath = response.data.filePath; // Assuming backend returns { filePath: "..." }
    toastStore.addToast('File uploaded successfully!', 'success');
  } catch (error) {
    console.error('File upload failed:', error);
    toastStore.addToast('File upload failed. Please try again.', 'error');
    form.attachment = null; // Clear selected file on error
    form.attachmentPath = '';
  } finally {
    isUploading.value = false;
  }
};

const handleFileChange = (event) => {
  const file = event.target.files[0];
  if (file) {
    form.attachment = file;
    uploadAttachment(file);
  } else {
    form.attachment = null;
    form.attachmentPath = '';
  }
};

const handleSubmit = async () => {
  const isValid = await v$.value.$validate();
  if (!isValid) return;

  const payload = {
    label: form.label,
    task_type_id: form.taskTypeId,
    priority: form.priority,
    start_date: form.startDate,
    due_date: form.dueDate || null,
    description: form.description,
    attachment: form.attachmentPath,
    assigned_to_users: form.assignedToUsers.map(user => user.id),
    assigned_to_groups: form.assignedToGroups.map(group => group.id),
    follow_up_users: form.followUpUsers.map(user => user.id),
    follow_up_groups: form.followUpGroups.map(group => group.id),
  };

  emit('submit', payload);
};
</script>

<style scoped>
.form-input {
  @apply block w-full bg-gray-700 border border-gray-600 rounded-md shadow-sm py-2 px-3 text-sm text-white;
  @apply focus:outline-none focus:ring-sky-500 focus:border-sky-500 transition;
}
</style>
