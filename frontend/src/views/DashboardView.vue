<template>
  <div>
    <h1 class="text-3xl font-bold text-white">Dashboard</h1>
    <p class="text-gray-400 mt-2">Welcome to your Task Manager dashboard.</p>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mt-8">
      <div class="bg-gray-800 rounded-lg p-6 shadow-lg transform hover:scale-105 transition-transform duration-300">
        <div class="flex items-center">
          <div class="bg-yellow-500 rounded-full p-3">
            <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-gray-400">Pending</p>
            <p class="text-2xl font-bold text-white">{{ taskCounts.Pending || 0 }}</p>
          </div>
        </div>
      </div>
      <div class="bg-gray-800 rounded-lg p-6 shadow-lg transform hover:scale-105 transition-transform duration-300">
        <div class="flex items-center">
          <div class="bg-blue-500 rounded-full p-3">
            <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-gray-400">In Progress</p>
            <p class="text-2xl font-bold text-white">{{ taskCounts['In Progress'] || 0 }}</p>
          </div>
        </div>
      </div>
      <div class="bg-gray-800 rounded-lg p-6 shadow-lg transform hover:scale-105 transition-transform duration-300">
        <div class="flex items-center">
          <div class="bg-purple-500 rounded-full p-3">
            <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"></path>
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-gray-400">In Review</p>
            <p class="text-2xl font-bold text-white">{{ taskCounts['In Review'] || 0 }}</p>
          </div>
        </div>
      </div>
      <div class="bg-gray-800 rounded-lg p-6 shadow-lg transform hover:scale-105 transition-transform duration-300">
        <div class="flex items-center">
          <div class="bg-green-500 rounded-full p-3">
            <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-gray-400">Completed</p>
            <p class="text-2xl font-bold text-white">{{ taskCounts.Completed || 0 }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue';
import { useTasksStore } from '@/stores/tasks';
import { storeToRefs } from 'pinia';

const tasksStore = useTasksStore();
const { tasks } = storeToRefs(tasksStore);

onMounted(() => {
  tasksStore.fetchMyTasks();
});

const taskCounts = computed(() => {
  return tasks.value.reduce((acc, task) => {
    acc[task.Status] = (acc[task.Status] || 0) + 1;
    return acc;
  }, {});
});
</script>
