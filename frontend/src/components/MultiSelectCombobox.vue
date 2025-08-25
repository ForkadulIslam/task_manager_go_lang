<template>
  <Combobox as="div" :modelValue="modelValue" @update:modelValue="value => emit('update:modelValue', value)" multiple>
    <div class="relative">
            <div class="relative w-full cursor-default overflow-hidden rounded-md bg-gray-700 text-left shadow-sm focus-within:border-sky-500 focus-within:ring-1 focus-within:ring-sky-500 sm:text-sm">
        <div class="flex flex-wrap items-center gap-1 py-1 pl-3 pr-10">
          <span v-for="item in modelValue" :key="item.id" class="px-2 py-0.5 rounded-full text-xs font-medium bg-sky-800 text-sky-200">
            {{ item[props.displayProperty] }}
          </span>
          <ComboboxInput
            class="flex-grow border-none py-1 pl-0 pr-0 text-sm leading-5 text-white bg-gray-700 focus:ring-0"
            :displayValue="(item) => query === '' ? '' : item[props.displayProperty]"
            @change="query = $event.target.value"
          />
        </div>
        <ComboboxButton class="absolute inset-y-0 right-0 flex items-center pr-2">
          <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path fill-rule="evenodd" d="M10 3a.75.75 0 01.53.22l3.5 3.5a.75.75 0 01-1.06 1.06L10 4.81 7.53 7.28a.75.75 0 01-1.06-1.06l3.5-3.5a.75.75 0 0110 3zm-3.72 9.28a.75.75 0 011.06 0L10 15.19l2.47-2.47a.75.75 0 111.06 1.06l-3.5 3.5a.75.75 0 01-1.06 0l-3.5-3.5a.75.75 0 010-1.06z" clip-rule="evenodd" />
          </svg>
        </ComboboxButton>
      </div>

      <transition leave-active-class="transition ease-in duration-100" leave-from-class="opacity-100" leave-to-class="opacity-0">
        <ComboboxOptions class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-gray-700 py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
          <div v-if="filteredItems.length === 0 && query !== ''" class="relative cursor-default select-none py-2 px-4 text-gray-400">
            Nothing found.
          </div>
          <ComboboxOption v-for="item in filteredItems" :key="item.id" :value="item" as="template" v-slot="{ active, selected }">
            <li :class="['relative cursor-default select-none py-2 pl-3 pr-9', active ? 'bg-sky-600 text-white' : 'text-gray-200']">
              <span :class="['block truncate', selected && 'font-semibold']">{{ item[props.displayProperty] }}</span>
              <span v-if="selected" :class="['absolute inset-y-0 right-0 flex items-center pr-4', active ? 'text-white' : 'text-sky-600']">
                <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clip-rule="evenodd" />
                </svg>
              </span>
            </li>
          </ComboboxOption>
        </ComboboxOptions>
      </transition>
    </div>
  </Combobox>
</template>

<script setup>
import { ref, computed } from 'vue';
import {
  Combobox,
  ComboboxButton,
  ComboboxInput,
  ComboboxOptions,
  ComboboxOption,
} from '@headlessui/vue';

const props = defineProps({
  items: {
    type: Array,
    required: true,
  },
  modelValue: {
    type: Array,
    required: true,
  },
  displayProperty: {
    type: String,
    default: 'username', // Default to 'username' for backward compatibility
  },
});

const emit = defineEmits(['update:modelValue']);

const query = ref('');
const filteredItems = computed(() =>
  query.value === ''
    ? props.items
    : props.items.filter((item) => {
        return item[props.displayProperty]
          .toLowerCase()
          .includes(query.value.toLowerCase());
      })
);
</script>
