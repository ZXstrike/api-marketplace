<template>
  <header class="bg-white/80 dark:bg-gray-900/80 backdrop-blur-lg sticky top-0 z-50 border-b border-gray-200 dark:border-gray-800 justify-items-center">
    <nav class="container mx-auto px-6 py-4 flex justify-between items-center">
      <router-link to="/" class="flex items-center space-x-2">
        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-blue-600"><path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"></path></svg>
        <span class="text-2xl font-bold text-gray-800 dark:text-white">Go-API Mart</span>
      </router-link>
      
      <!-- Menus for Unauthenticated Users -->
      <div v-if="!authStore.isLoggedIn" class=" flex items-center gap-6">
          <router-link to="/browse" class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 font-medium">APIs</router-link>
          <a href="/#pricing" class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 font-medium">Pricing</a>
          <a href="/#features" class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 font-medium">Features</a>
      </div>

       <!-- Menus for Authenticated Users -->
      <div v-if="authStore.isLoggedIn" class=" lg:flex items-center gap-6">
            <router-link v-if="authStore.user.role !== 'provider'" to="/browse" class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 font-medium">APIs</router-link>
          <router-link to="/dashboard" class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 font-medium">Dashboard</router-link>
          <a href="#" class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 font-medium">Docs</a>
      </div>


      <!-- Auth Buttons for Unauthenticated Users -->
      <div v-if="!authStore.isLoggedIn" class=" flex items-center gap-4">
          <router-link to="/login" class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 font-semibold">Log In</router-link>
          <router-link to="/register" class="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-5 py-2.5 rounded-lg shadow-md">Sign Up</router-link>
      </div>

      <!-- Profile Dropdown for Authenticated Users -->
      <div v-if="authStore.isLoggedIn" class=" lg:flex items-center gap-4">
        <!-- <button class="text-gray-500 dark:text-gray-400 hover:text-blue-600">
            <i data-feather="bell" class="w-6 h-6"></i>
        </button> -->
        <div class="relative">
            <button @click="isProfileOpen = !isProfileOpen" class="flex items-center gap-3">
                <img class="w-9 h-9 rounded-full object-cover" :src="authStore.user.avatarUrl" alt="User avatar">
                <span class="font-semibold text-gray-700 dark:text-gray-200">{{ authStore.user.name }}</span>
                <i data-feather="chevron-down" class="w-4 h-4 text-gray-500 transition-transform" :class="{'rotate-180': isProfileOpen}"></i>
            </button>
            <!-- Profile Dropdown Menu -->
            <div v-if="isProfileOpen" @click="isProfileOpen = false" class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-xl z-50 border dark:border-gray-700 py-1">
                <router-link to="/dashboard" class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">Dashboard</router-link>
                <router-link to="/settings" class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">Settings</router-link>
                <div class="border-t border-gray-200 dark:border-gray-700 my-1"></div>
                <a href="#" @click.prevent="handleLogout" class="block w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700">Log Out</a>
            </div>
        </div>
      </div>
      
    </nav>
  </header>
</template>

<script setup>
import { ref, watch } from 'vue';
import { RouterLink, useRouter } from 'vue-router';
// This is where you would import your actual Pinia store
import { useAuthStore } from '@/store/auth';

const authStore = useAuthStore();
const router = useRouter();
const isProfileOpen = ref(false);


function handleLogout() {
  authStore.logout();
  router.push('/');
}
</script>
