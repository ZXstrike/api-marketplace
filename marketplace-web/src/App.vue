<template>
  <div id="app-wrapper" class="bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100 antialiased min-h-screen">
    <Navbar v-if="!$route.meta.hideNavbar" />

    <main class="">
      <RouterView :key="$route.fullPath" />
    </main>

    <Footer v-if="!$route.meta.hideFooter" />
  </div>
</template>

<script setup>
import { onMounted, watch, nextTick } from 'vue';
import { RouterView, useRoute } from 'vue-router';
import Navbar from '@/components/core/Navbar.vue';
import Footer from '@/components/core/Footer.vue';

// Force dark theme on application start
onMounted(() => {
  document.documentElement.classList.add('dark');
  if (window.feather) {
    window.feather.replace();
  }
});

const route = useRoute();

// For feather icons on route change
watch(() => route.path, async () => {
  await nextTick();
  if (window.feather) {
    window.feather.replace();
    console.log('Feather icons updated');
  }
});


</script>

<style>
/* You can place global styles that are not part of Tailwind here */
</style>