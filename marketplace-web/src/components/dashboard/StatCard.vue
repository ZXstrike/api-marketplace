<template>
  <div class="dashboard-card">
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">{{ title }}</h3>
      <i :data-feather="icon" :class="iconColorClass"></i>
    </div>
    <p class="text-2xl font-bold !mt-2">{{ value }}</p>
  </div>
</template>

<script setup>
import { computed, onMounted, nextTick } from 'vue';

const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  value: {
    type: [String, Number],
    required: true,
  },
  icon: {
    type: String,
    required: true,
  },
  color: {
    type: String,
    default: 'blue', // warna default
  }
});

const iconColorClass = computed(() => {
    const colorMap = {
        green: 'text-green-500',
        blue: 'text-blue-500',
        yellow: 'text-yellow-500',
        red: 'text-red-500',
    };
    return colorMap[props.color] || 'text-gray-500';
});

// Pastikan ikon Feather di-render saat komponen dimuat
onMounted(() => {
    nextTick(() => {
        if(window.feather) {
            window.feather.replace();
        }
    });
});
</script>

<style scoped>
</style>
