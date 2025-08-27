<template>
  <div
    :class="[
      'rounded-full flex items-center justify-center font-bold',
      avatarSizeClass,
      backgroundColorClass,
      textColorClass
    ]"
  >
    {{ initials }}
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  username: {
    type: String,
    required: true,
  },
  size: {
    type: String,
    default: 'md', // 'sm', 'md', 'lg'
  },
});

const initials = computed(() => {
  if (!props.username) return '?';
  const parts = props.username.split(' ');
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  return props.username.substring(0, 2).toUpperCase();
});

const colors = [
  'bg-red-500', 'bg-pink-500', 'bg-purple-500', 'bg-indigo-500',
  'bg-blue-500', 'bg-cyan-500', 'bg-teal-500', 'bg-green-500',
  'bg-lime-500', 'bg-yellow-500', 'bg-amber-500', 'bg-orange-500',
  'bg-brown-500', 'bg-gray-500', 'bg-blue-gray-500'
];

const generateColorIndex = (username) => {
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 5) - hash);
  }
  return Math.abs(hash % colors.length);
};

const backgroundColorClass = computed(() => {
  return colors[generateColorIndex(props.username)];
});

const textColorClass = computed(() => {
  return 'text-white';
});

const avatarSizeClass = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'w-6 h-6 text-xs';
    case 'lg':
      return 'w-10 h-10 text-base';
    case 'md':
    default:
      return 'w-8 h-8 text-sm';
  }
});
</script>

<style scoped>
</style>