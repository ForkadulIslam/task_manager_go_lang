<template>
  <div class="toast-container">
    <transition-group name="toast-animation" tag="div">
      <ToastNotification
        v-for="toast in toastStore.toasts"
        :key="toast.id"
        :message="toast.message"
        :type="toast.type"
        @close="toastStore.removeToast(toast.id)"
      />
    </transition-group>
  </div>
</template>

<script setup>
import { useToastStore } from '../stores/toast';
import ToastNotification from './ToastNotification.vue';

const toastStore = useToastStore();
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: 9999;
}

.toast-animation-enter-active,
.toast-animation-leave-active {
  transition: all 0.5s cubic-bezier(0.68, -0.55, 0.27, 1.55);
}

.toast-animation-enter-from,
.toast-animation-leave-to {
  opacity: 0;
  transform: translateX(100%) scale(0.8);
}

.toast-animation-move {
  transition: transform 0.4s ease;
}
</style>
