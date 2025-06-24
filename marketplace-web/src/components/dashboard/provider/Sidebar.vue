<template>
  <aside class="lg:w-1/4 xl:w-1/5">
    <nav class="space-y-2">
      <a href="#" @click.prevent="$emit('navigate', 'overview')" :class="{'active': activeSection === 'overview'}" class="sidebar-link">
        <i data-feather="pie-chart" class="w-5 h-5 mr-3"></i>
        <span>Overview</span>
      </a>
      <a href="#" @click.prevent="$emit('navigate', 'apis')" :class="{'active': activeSection === 'apis'}" class="sidebar-link">
        <i data-feather="package" class="w-5 h-5 mr-3"></i>
        <span>My APIs</span>
      </a>
      <a href="#" @click.prevent="$emit('navigate', 'analytics')" :class="{'active': activeSection === 'analytics'}" class="sidebar-link">
          <i data-feather="bar-chart-2" class="w-5 h-5 mr-3"></i>
          <span>Analytics</span>
      </a>
      <a href="#" class="sidebar-link">
          <i data-feather="settings" class="w-5 h-5 mr-3"></i>
          <span>Store Settings</span>
      </a>
    </nav>
  </aside>
</template>

<script setup>
// The component now expects an 'activeSection' prop...
defineProps({
  activeSection: {
    type: String,
    required: true
  }
});

// ...and declares that it will emit a 'navigate' event.
defineEmits(['navigate']);
</script>

<style scoped>
</style>