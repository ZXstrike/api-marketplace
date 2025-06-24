<template>
  <aside class="lg:w-1/4 xl:w-1/5">
    <nav class="space-y-2">
      <a href="#" 
         @click.prevent="$emit('navigate', 'overview')" 
         :class="{ 'active': activeSection === 'overview' }"
         class="sidebar-link">
        <i data-feather="home" class="w-5 h-5 mr-3"></i>
        <span>Overview</span>
      </a>
      <a href="#" 
         @click.prevent="$emit('navigate', 'subscriptions')"
         :class="{ 'active': activeSection === 'subscriptions' }" 
         class="sidebar-link">
        <i data-feather="key" class="w-5 h-5 mr-3"></i>
        <span>My Subscriptions</span>
      </a>
      <a href="#" 
         @click.prevent="$emit('navigate', 'billing')" 
         :class="{ 'active': activeSection === 'billing' }"
         class="sidebar-link">
        <i data-feather="credit-card" class="w-5 h-5 mr-3"></i>
        <span>Billing & Usage</span>
      </a>
    </nav>
  </aside>
</template>

<script setup>
import { onMounted, nextTick } from 'vue';

// Define props and emits for component communication
defineProps({
  activeSection: {
    type: String,
    required: true
  }
});

defineEmits(['navigate']);

// If you are using a library like feather-icons, you might need to initialize it
// when the component mounts to ensure icons are rendered.
onMounted(() => {
  if (window.feather) {
    window.feather.replace();
  }
});
</script>

<style scoped>

</style>