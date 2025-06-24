<template>
    <div v-if="apis.length > 0" class="grid gap-6 sm:grid-cols-2 xl:grid-cols-3">
        <ApiCard v-for="api in apis" :key="api.id" :api="api" />
    </div>
    <div v-else class="text-center py-16 bg-gray-50 dark:bg-gray-800/50 rounded-xl">
        <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-300">No APIs Found</h3>
        <p class="text-gray-500 dark:text-gray-400 mt-2">Try adjusting your search or filter criteria.</p>
    </div>
</template>

<script setup>
import ApiCard from '@/components/api/ApiCard.vue';

defineProps({
  apis: {
    type: Array,
    required: true,
  },
});
</script>