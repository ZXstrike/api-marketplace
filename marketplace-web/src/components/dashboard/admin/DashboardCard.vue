<template>
  <div class="dashboard-card">
    <h3 class="text-xl font-bold mb-4">{{ title }}</h3>
    <slot />
  </div>
</template>
<script setup>
defineProps({ title: { type: String, required: true } });
</script>