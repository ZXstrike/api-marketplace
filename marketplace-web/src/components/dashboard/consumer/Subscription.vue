<template>
  <section class="flex flex-col gap-2">
    <h1 class="text-3xl font-bold">My Subscriptions</h1>
    <p class="text-gray-600 dark:text-gray-400 -mt-4">Manage your API subscriptions and their associated keys.</p>

    <div v-if="!subscriptions.length" class="dashboard-card text-center">
      <p class="text-gray-500">You are not subscribed to any APIs yet.</p>
    </div>

    <div class="flex flex-col gap-3">
      <div v-for="sub in subscriptions" :key="sub.id" class="dashboard-card">
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center">
          <div>
            <h2 class="text-xl font-bold">{{ sub.apiName }}</h2>
            <p class="text-sm text-gray-500 dark:text-gray-400">Subscribed on: {{ sub.subscribedDate }}</p>
          </div>
          <div class="flex flex-col md:flex-row items-start md:items-center gap-2 md:gap-4">
            <router-link :to="`/api-details/${sub.apiId}`"
              class="mt-3 md:mt-0 text-blue-500 hover:text-blue-700 font-semibold flex items-center space-x-2 bg-blue-100/50 dark:bg-blue-900/30 hover:bg-blue-100 dark:hover:bg-blue-900/50 px-4 py-2 rounded-lg">
              <i data-feather="info" class="w-4 h-4"></i>
              <span>See Details</span>
            </router-link>
            <button @click="$emit('unsubscribe', sub.id)"
              class="mt-3 md:mt-0 text-red-500 hover:text-red-700 font-semibold flex items-center space-x-2 bg-red-100/50 dark:bg-red-900/30 hover:bg-red-100 dark:hover:bg-red-900/50 px-4 py-2 rounded-lg">
              <i data-feather="trash-2" class="w-4 h-4"></i>
              <span>Unsubscribe</span>
            </button>
          </div>
        </div>
        <div class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
          <h3 class="font-semibold mb-2">API Key</h3>
          <div class="flex items-center space-x-4 p-3 bg-gray-100 dark:bg-gray-900/50 rounded-lg">
            <p class="font-mono text-gray-700 dark:text-gray-300 flex-1">{{ sub.apiKeyMasked }}</p>
            <button @click="emit('copy-key', sub.apiKey)" title="Copy"
              class="p-2 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md"><i data-feather="copy"
                class="w-5 h-5"></i></button>
            <button @click="$emit('regenerate-key', sub.id)" title="Regenerate"
              class="p-2 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md text-blue-500"><i data-feather="refresh-cw"
                class="w-5 h-5"></i></button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted, nextTick } from 'vue';

const props = defineProps({
  subscriptions: {
    type: Array,
    required: true
  }
});


const emit = defineEmits(['unsubscribe', 'regenerate-key', 'copy-key']);

onMounted(() => {
  if (window.feather) {
    nextTick(() => {
      window.feather.replace();
    });
  }
  console.log('Subscription component mounted');
  console.log('Subscriptions:', props.subscriptions);
});
</script>