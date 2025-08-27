<template>
  <form @submit.prevent="handleSubmit">
    <div class="space-y-4">
      <div>
        <label for="label" class="block text-sm font-medium text-gray-300 mb-1">Group Name</label>
        <input v-model="form.label" type="text" id="label" class="form-input" :class="{ 'border-red-500': v$.label.$error }">
        <div v-if="v$.label.$error" class="text-red-400 text-xs mt-1">{{ v$.label.$errors[0].$message }}</div>
      </div>

      <div class="pt-2">
        <button type="submit" :disabled="props.isLoading" class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-sky-600 hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-sky-500 transition-colors">
          Create Group
        </button>
      </div>
    </div>
  </form>
</template>

<script setup>
import { reactive, onMounted } from 'vue';
import { useVuelidate } from '@vuelidate/core';
import { required, minLength, maxLength, helpers } from '@vuelidate/validators';
import { useMetaStore } from '../stores/meta';

const emit = defineEmits(['submit']);

const props = defineProps({
  isLoading: {
    type: Boolean,
    default: false,
  },
});

const metaStore = useMetaStore();

const form = reactive({
  label: '',
});

const rules = {
  label: {
    required: helpers.withMessage('Group name is required', required),
    minLength: minLength(3),
    maxLength: maxLength(100),
  },
};

const v$ = useVuelidate(rules, form);

onMounted(async () => {
  await metaStore.fetchMeta();
});

const handleSubmit = async () => {
  const isValid = await v$.value.$validate();
  if (!isValid) return;

  emit('submit', {
    label: form.label,
  });
};
</script>

<style scoped>
.form-input {
  @apply block w-full bg-gray-700 border border-gray-600 rounded-md shadow-sm py-2 px-3 text-sm text-white;
  @apply focus:outline-none focus:ring-sky-500 focus:border-sky-500 transition;
}
</style>
