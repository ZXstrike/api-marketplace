<template>
  <div class="flex justify-center">
    <nav class="inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
      <button 
        @click="$emit('page-changed', currentPage - 1)" 
        :disabled="currentPage <= 1"
        class="relative inline-flex items-center px-4 py-2 rounded-l-md border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed">
        Previous
      </button>
      
      <span class="relative inline-flex items-center px-4 py-2 border-y border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800 text-sm font-medium text-gray-700 dark:text-gray-300">
        Page {{ currentPage }}
      </span>
      
      <button 
        @click="$emit('page-changed', currentPage + 1)" 
        :disabled="isLastPage"
        class="relative inline-flex items-center px-4 py-2 rounded-r-md border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed">
        Next
      </button>
    </nav>
  </div>
</template>

<script setup>
defineProps({
    currentPage: {
        type: Number,
        required: true,
    },
    isLastPage: {
        type: Boolean,
        default: false,
    }
});

defineEmits(['page-changed']);
</script>