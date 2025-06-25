<template>
  <div v-if="api" class="flex flex-col md:flex-row justify-between md:items-center gap-4">
    <div>
      <div class="flex items-center gap-4">
        <img 
          :src="api.icon_url || 'https://placehold.co/64x64/3B82F6/FFFFFF?text=API'" 
          alt="API Logo" 
          class="w-16 h-16 rounded-lg shadow-md object-cover"
        >
        <div>
          <h1 class="text-4xl font-bold">{{ api.name }}</h1>
          <p class="text-gray-500 dark:text-gray-400">
            By <a href="#" class="font-semibold text-blue-600 hover:underline">{{ api.provider.username }}</a>
          </p>
        </div>
      </div>
    </div>
    <div class="flex-shrink-0">
      <!-- This assumes a reactive `isLoggedIn` boolean is available in the component's scope -->
      <button v-if="authStore.isLoggedIn" @click="$emit('subscribe')" class="w-full md:w-auto bg-blue-600 hover:bg-blue-700 text-white font-bold px-8 py-3 rounded-lg shadow-lg flex items-center justify-center space-x-2">
      <i data-feather="key" class="w-5 h-5"></i>
      <span>Subscribe to this API</span>
      </button>
      <router-link v-else to="/login" class="w-full md:w-auto bg-gray-500 hover:bg-gray-600 text-white font-bold px-8 py-3 rounded-lg shadow-lg flex items-center justify-center space-x-2">
        <i data-feather="log-in" class="w-5 h-5"></i>
        <span>Sign In to Subscribe</span>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/store/auth';

const authStore = useAuthStore();

defineProps({
  api: {
    type: Object,
    required: true,
  }
});

defineEmits(['subscribe']);
</script>