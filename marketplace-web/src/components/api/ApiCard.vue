<template>
  <div class="card p-0 overflow-hidden flex flex-col h-full">
    <div class="p-6">
      <div class="flex items-center justify-between">
        <h3 class="text-xl font-bold text-gray-900 dark:text-white truncate pr-2">{{ api.name }}</h3>
        <span 
          v-if="api.categories && api.categories.length > 0"
          :class="categoryColor(api.categories[0].name)"
          class="text-xs font-semibold px-2.5 py-1 rounded-full flex-shrink-0">
          {{ api.categories[0].name }}
        </span>
      </div>
      <p class="mt-3 text-gray-600 dark:text-gray-400 h-9 line-clamp-3">{{ api.description }}</p>
    </div>
    
    <div class="border-t border-gray-200 dark:border-gray-700 py-3 px-4 bg-gray-50/50 dark:bg-gray-800/20">
      <div class="flex items-center text-sm gap-2">
          <img :src="fullProfilePictureUrl" alt="Provider" class="w-8 h-8 rounded-full">
          <div>
         <span class="font-semibold">{{ api.provider.username }}</span>
          </div>
      </div>
    </div>

    <div class="mt-auto pb-3 px-3">
        <router-link :to="`/api-details/${api.id}`" class="block w-full text-center bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 font-semibold py-2 rounded-lg transition-colors">
            View Details
        </router-link>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { RouterLink } from 'vue-router';

const props = defineProps({
  api: {
    type: Object,
    required: true,
  }
});

// A computed property that builds the full, absolute URL for the image
const fullProfilePictureUrl = computed(() => {
  // Get the base URL from the environment variable.
  // Vite exposes these on `import.meta.env`
  const baseUrl = import.meta.env.VITE_API_BASE_URL;

  // Check that the URL parts exist to avoid errors
  if (baseUrl && props.api?.provider?.profile_picture_url) {
    // Combine the base URL and the relative path from the API
    // Result: "http://127.0.0.1:8081/files/profiles/image.png"
    return `${baseUrl}${props.api.provider.profile_picture_url}`;
  }

  // Optional: Provide a fallback to a default image in your public folder
  // if the API doesn't provide a URL.
  return '/default-avatar.png';
});

const categoryColor = (categoryName) => {
  // Return a default color if categoryName is not provided
  if (!categoryName) {
    return 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300';
  }

  // A palette of colors for dynamic assignment for categories
  const colorPalette = [
    'bg-blue-100 dark:bg-blue-900/50 text-blue-600 dark:text-blue-300',
    'bg-purple-100 dark:bg-purple-900/50 text-purple-600 dark:text-purple-300',
    'bg-green-100 dark:bg-green-900/50 text-green-600 dark:text-green-300',
    'bg-red-100 dark:bg-red-900/50 text-red-600 dark:text-red-300',
    'bg-yellow-100 dark:bg-yellow-900/50 text-yellow-600 dark:text-yellow-300',
    'bg-indigo-100 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-300',
    'bg-pink-100 dark:bg-pink-900/50 text-pink-600 dark:text-pink-300',
    'bg-orange-100 dark:bg-orange-900/50 text-orange-600 dark:text-orange-300',
    'bg-teal-100 dark:bg-teal-900/50 text-teal-600 dark:text-teal-300',
    'bg-cyan-100 dark:bg-cyan-900/50 text-cyan-600 dark:text-cyan-300',
  ];

  // Simple hash function to get a consistent color for a category name
  let hash = 0;
  for (let i = 0; i < categoryName.length; i++) {
    const char = categoryName.charCodeAt(i);
    hash = ((hash << 5) - hash) + char;
    hash = hash | 0; // Convert to 32bit integer
  }
  
  const index = Math.abs(hash % colorPalette.length);
  return colorPalette[index];
};
</script>