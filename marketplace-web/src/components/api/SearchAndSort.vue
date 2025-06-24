<template>
  <div class="flex flex-col md:flex-row justify-between items-center gap-3">
    <div class="relative w-full md:w-1/3 lg:w-full">
      <input 
        type="search" 
        placeholder="Search APIs..." 
        class="form-input w-full py-3 pl-4 pr-10"
        :value="searchQuery"
        @input="$emit('update:searchQuery', $event.target.value)"
      >
      <i data-feather="search" class="absolute top-1/2 right-4 -translate-y-1/2 text-gray-400"></i>
    </div>
  </div>
</template>

<script setup>
defineProps({
  searchQuery: String,
});

defineEmits(['update:searchQuery']);
</script>