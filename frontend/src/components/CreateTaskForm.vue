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

      <!-- Description -->
      <div>
        <label for="description" class="block text-sm font-medium text-gray-300 mb-1">Description</label>
        <textarea v-model="form.description" id="description" rows="3" class="form-input"></textarea>
      </div>

      <!-- Submit Button -->
      <div class="pt-2">
        <button type="submit" class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-sky-600 hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-sky-500 transition-colors">Create Task</button>
      </div>
    </div>
  </form>
</template>

<script setup>
import { reactive, onMounted } from 'vue';
import { useMetaStore } from '../stores/meta';
import { useVuelidate } from '@vuelidate/core';
import { required, minLength, helpers } from '@vuelidate/validators';
import MultiSelectCombobox from './MultiSelectCombobox.vue';

const emit = defineEmits(['submit']);

const metaStore = useMetaStore();

const form = reactive({
  label: '',
  taskTypeId: '',
  priority: 'Normal',
  description: '',
  assignedToUsers: [],
  assignedToGroups: [],
});

const rules = {
  label: { required, minLength: minLength(3) },
  taskTypeId: { required: helpers.withMessage('Task type is required', required) },
  priority: { required },
};

const v$ = useVuelidate(rules, form);

onMounted(() => {
  metaStore.fetchMeta();
});

const handleSubmit = async () => {
  const isValid = await v$.value.$validate();
  if (!isValid) return;

  const payload = {
    ...form,
    assigned_to_users: form.assignedToUsers.map(user => user.ID),
    assigned_to_groups: form.assignedToGroups.map(group => group.ID),
  };
  
  delete payload.assignedToUsers;
  delete payload.assignedToGroups;

  emit('submit', payload);
};
</script>

<style scoped>
.form-input {
  @apply block w-full bg-gray-700 border border-gray-600 rounded-md shadow-sm py-2 px-3 text-sm text-white;
  @apply focus:outline-none focus:ring-sky-500 focus:border-sky-500 transition;
}
</style>
