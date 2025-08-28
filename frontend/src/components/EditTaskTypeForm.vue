<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <div>
      <label for="label" class="block text-sm font-medium text-gray-300 mb-1">Task Type Label</label>
      <input
        type="text"
        id="label"
        v-model="editedLabel"
        required
        class="form-input"
      />
    </div>
    <button
      type="submit"
      :disabled="isLoading"
      class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-sky-600 hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-sky-500"
    >
      {{ isLoading ? 'Updating...' : 'Update Task Type' }}
    </button>
  </form>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useToastStore } from '../stores/toast'; // New import

const emit = defineEmits(['submit']);
const props = defineProps({
  initialLabel: {
    type: String,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});

const toastStore = useToastStore(); // New initialization

const editedLabel = ref(props.initialLabel);

// Watch for changes in initialLabel prop to update editedLabel
watch(() => props.initialLabel, (newVal) => {
  editedLabel.value = newVal;
});

const handleSubmit = () => {
  if (editedLabel.value.trim()) {
    emit('submit', { label: editedLabel.value.trim() });
  } else {
    toastStore.addToast('Task type label cannot be empty.', 'info');
  }
};
</script>

<style scoped>
.form-input {
  @apply block w-full bg-gray-700 border border-gray-600 rounded-md shadow-sm py-2 px-3 text-sm text-white;
  @apply focus:outline-none focus:ring-sky-500 focus:border-sky-500 transition;
}
</style>
