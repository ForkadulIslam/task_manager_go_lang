<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60" @click.self="close">
    <div class="bg-gray-800 rounded-lg shadow-lg w-full max-w-2xl border border-gray-700">
      <header class="p-4 border-b border-gray-700 flex justify-between items-center">
        <slot name="header">
          <h2 class="text-xl font-bold">Modal Title</h2>
        </slot>
        <button @click="close" class="text-gray-400 hover:text-white transition-colors">&times;</button>
      </header>
      <section class="p-6">
        <slot name="body">
          <p>This is the default body.</p>
        </slot>
      </section>
      <footer class="p-4 border-t border-gray-700">
        <slot name="footer">
          <!-- Default footer can be empty or have a close button -->
        </slot>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue';

const props = defineProps({
  show: {
    type: Boolean,
    required: true,
  },
});

const emit = defineEmits(['close']);

const close = () => {
  emit('close');
};

const handleKeydown = (e) => {
  if (e.key === 'Escape' && props.show) {
    close();
  }
};

onMounted(() => {
  document.addEventListener('keydown', handleKeydown);
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown);
});
</script>
