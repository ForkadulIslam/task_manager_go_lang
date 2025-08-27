<template>
  <div :class="['toast', `toast--${type}`]">
    <div class="toast__icon">
      <svg v-if="type === 'success'" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      <svg v-if="type === 'error'" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      <svg v-if="type === 'info'" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
    </div>
    <div class="toast__content">
      <p class="toast__message">{{ message }}</p>
    </div>
    <button @click="$emit('close')" class="toast__close-button">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
    </button>
  </div>
</template>

<script setup>
defineProps({
  message: {
    type: String,
    required: true,
  },
  type: {
    type: String,
    default: 'info', // 'info', 'success', 'error'
  },
});
defineEmits(['close']);
</script>

<style scoped>
.toast {
  display: flex;
  align-items: center;
  padding: 1rem;
  margin-bottom: 1rem;
  border-radius: 0.5rem;
  color: white;
  background: linear-gradient(145deg, #2d3748, #1a202c);
  border: 1px solid #4a5568;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05), 0 0 20px rgba(0, 0, 0, 0.2);
  position: relative;
  overflow: hidden;
  width: 350px;
}

.toast::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><path d="M10 10 Q 50 10, 90 90" stroke="%23ffffff10" stroke-width="1" fill="none" /><path d="M20 50 Q 70 20, 80 80" stroke="%23ffffff10" stroke-width="1" fill="none" /><path d="M5 80 Q 40 10, 95 50" stroke="%23ffffff10" stroke-width="1" fill="none" /></svg>');
  opacity: 0.5;
  z-index: 0;
}

.toast--success {
  border-left: 4px solid #48bb78;
  box-shadow: 0 0 15px #48bb7880;
}

.toast--error {
  border-left: 4px solid #f56565;
  box-shadow: 0 0 15px #f5656580;
}

.toast--info {
  border-left: 4px solid #4299e1;
  box-shadow: 0 0 15px #4299e180;
}

.toast__icon {
  flex-shrink: 0;
  margin-right: 0.75rem;
  z-index: 1;
}

.toast--success .toast__icon { color: #48bb78; }
.toast--error .toast__icon { color: #f56565; }
.toast--info .toast__icon { color: #4299e1; }

.toast__content {
  flex-grow: 1;
  z-index: 1;
}

.toast__message {
  font-size: 0.875rem;
  font-weight: 500;
}

.toast__close-button {
  margin-left: 1rem;
  color: #a0aec0;
  background: transparent;
  border: none;
  cursor: pointer;
  z-index: 1;
}
.toast__close-button:hover {
  color: white;
}
</style>
