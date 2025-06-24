<template>
  <aside>
    <div class="p-6 bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-xl shadow-lg lg:sticky top-28">
      <h2 class="text-xl font-bold mb-6">Filters</h2>

      <div>
        <h3 class="font-semibold mb-3">Categories</h3>
        <ul class="space-y-2 text-gray-600 dark:text-gray-300">
          <li>
            <a href="#" @click.prevent="$emit('clear-filters')" class="hover:text-blue-600 dark:hover:text-blue-400">All</a>
          </li>
          <li v-for="category in categories" :key="category.slug">
            <a href="#" @click.prevent="$emit('filter-change', category.slug)" class="hover:text-blue-600 dark:hover:text-blue-400">{{ category.name }}</a>
          </li>
        </ul>
      </div>

      <div class="mt-8 border-t border-gray-200 dark:border-gray-700 pt-6">
        <button @click="$emit('clear-filters')" class="w-full bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 font-semibold py-2.5 rounded-lg transition-colors">Clear Filters</button>
      </div>
    </div>
  </aside>
</template>

<script setup>
defineProps({
  categories: {
    type: Array,
    required: true,
  }
});

defineEmits(['filter-change', 'clear-filters']);
</script>