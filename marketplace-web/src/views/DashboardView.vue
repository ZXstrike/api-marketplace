<template>
    <!-- Dynamically render the correct dashboard based on user role -->
    <component :is="activeDashboard" />
</template>

<script setup>
import { computed } from 'vue';
import { useAuthStore } from '@/store/auth';

// FIX: Removed TheHeader and TheFooter imports as they are now handled globally in App.vue
import ConsumerDashboardView from '@/views/ConsumerDashboardView.vue';
import ProviderDashboardView from '@/views/ProviderDashboardView.vue';

const authStore = useAuthStore();

// This computed property determines which dashboard component to show.
const activeDashboard = computed(() => {
  if (authStore.userRole === 'provider') {
    return ProviderDashboardView;
  }
  return ConsumerDashboardView;
});
</script>
