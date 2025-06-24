<template>
  <div id="app-wrapper" class="bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100 antialiased min-h-screen">
    <Navbar v-if="!$route.meta.hideNavbar" />

    <main>
      <RouterView :key="$route.fullPath" />
    </main>

    <Footer v-if="!$route.meta.hideFooter" />
  </div>
</template>

<script setup>
import { onUpdated } from 'vue';
import { RouterView, useRoute } from 'vue-router'; // Import useRoute
import Navbar from '@/components/core/Navbar.vue';
import Footer from '@/components/core/Footer.vue';

// Although $route is available in the template, if you needed to access it
// in the script section for other logic, you would use this:
// const route = useRoute();

// This function re-initializes Feather Icons every time the view changes.
// This is crucial because Vue re-renders parts of the DOM, and the icons need to be redrawn.
onUpdated(() => {
  if (window.feather) {
    window.feather.replace();
    console.log('Feather icons updated');
  }
});
</script>

<style>
/* You can place global styles that are not part of Tailwind here */
</style>