<template>
  <form @submit.prevent="submitForm" class="space-y-4">
    <div>
      <label for="group-name" class="block text-sm font-medium text-gray-300">Group Name</label>
      <input
        type="text"
        id="group-name"
        v-model="form.label"
        class="form-input mt-1"
        required
      />
    </div>
    <div class="flex justify-end">
      <button
        type="submit"
        :disabled="isLoading"
        class="inline-flex items-center justify-center rounded-lg bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-2 focus:ring-offset-gray-900 disabled:opacity-50"
      >
        <span v-if="isLoading">Saving...</span>
        <span v-else>Save Changes</span>
      </button>
    </div>
  </form>
</template>

<script setup>
import { ref, watch } from 'vue';

const props = defineProps({
  group: {
    type: Object,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['submit']);

const form = ref({
  id: props.group?.id,
  label: props.group?.label || '',
});

watch(() => props.group, (newGroup) => {
  if (newGroup) {
    form.value.id = newGroup.id;
    form.value.label = newGroup.label;
  }
}, { immediate: true });


const submitForm = () => {
  if (!props.isLoading) {
    emit('submit', { ...form.value });
  }
};
</script>

<style scoped>
.form-input {
  @apply block w-full bg-gray-700 border border-gray-600 rounded-md shadow-sm py-2 px-3 text-sm text-white;
  @apply focus:outline-none focus:ring-sky-700 focus:border-sky-700 transition;
}
</style>
