<template>
  <div class="w-full">
    <!-- Main Content -->
    <main>
      <!-- Hero Section -->
      <section class=" relative py-20 md:py-36 bg-white dark:bg-gray-900 flex items-center justify-center">
        <!-- FIX: Re-instated the standard container pattern -->
        <div class="container mx-auto px-6 text-center relative flex flex-col items-center justify-center">
          <h1 class="text-4xl md:text-6xl font-extrabold text-gray-900 dark:text-white leading-tight">
            Flexible APIs, Fair Pricing.
            <br>
            <span class="bg-clip-text text-transparent bg-gradient-to-r from-blue-500 to-teal-400">Pay Only For What You
              Use.</span>
          </h1>
          <p class="pt-6 max-w-3xl mx-auto text-lg md:text-xl text-gray-600 dark:text-gray-300">
            Discover, integrate, and scale with our Golang-powered API marketplace. No subscriptions, no hidden fees.
            Just pure, consumption-based pricing.
          </p>
          <div class="pt-10 flex justify-center gap-4">
            <router-link to="/browse"
              class="bg-blue-600 hover:bg-blue-700 text-white font-bold text-lg px-8 py-4 rounded-lg shadow-lg">Explore
              APIs</router-link>
          </div>
        </div>
      </section>

      <!-- Featured APIs Section -->
      <section id="apis" class="py-20 bg-gray-50 dark:bg-gray-900/50 justify-items-center">
        <!-- FIX: Re-instated the standard container pattern -->
        <div class="container mx-auto px-6 justify-items-center">
          <h2 class="text-3xl font-bold text-gray-800 dark:text-white sm:text-4xl text-center">Explore Our API Catalog
          </h2>
          <p class="pt-4 pb-5 text-lg text-gray-600 dark:text-gray-300 text-center max-w-2xl mx-auto ">From data
            processing to machine learning, find the perfect API to power your next project.</p>
          <div class="mt-12 grid gap-8 grid-cols-1 md:grid-cols-2 lg:grid-cols-4 ">
            <!-- API Card examples -->
            <div
              class="bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-xl shadow-lg p-6 text-center">
              <h3 class="text-xl font-bold">Geolocation API</h3>
              <p class="mt-4 text-gray-600 dark:text-gray-400">Convert IP addresses to physical locations.</p>
            </div>
            <div
              class="bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-xl shadow-lg p-6 text-center">
              <h3 class="text-xl font-bold">Sentiment Analysis</h3>
              <p class="mt-4 text-gray-600 dark:text-gray-400">Analyze text to determine emotional tone.</p>
            </div>
            <div
              class="bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-xl shadow-lg p-6 text-center">
              <h3 class="text-xl font-bold">Currency Exchange</h3>
              <p class="mt-4 text-gray-600 dark:text-gray-400">Get real-time foreign exchange rates.</p>
            </div>
            <div
              class="bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-xl shadow-lg p-6 text-center">
              <h3 class="text-xl font-bold">Image Resizer</h3>
              <p class="mt-4 text-gray-600 dark:text-gray-400">Optimize and resize images on-the-fly.</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Features Section -->
      <section id="features" class="py-20 bg-white dark:bg-gray-900 justify-items-center">
        <!-- FIX: Re-instated the standard container pattern -->
        <div class="container mx-auto px-6">
          <h2 class="text-3xl font-bold text-gray-800 dark:text-white sm:text-4xl text-center">Built for Modern
            Developers</h2>
          <div class="mt-12 grid gap-8 grid-cols-1 md:grid-cols-2 lg:grid-cols-4 py-12">
            <div class="flex flex-col items-center text-center gap-4">
              <div
                class="w-16 h-16 bg-blue-100 dark:bg-blue-900/50 rounded-full flex items-center justify-center shrink-0">
                <i data-feather="dollar-sign" class="w-8 h-8 text-blue-600 dark:text-blue-400"></i></div>
              <h3 class="text-xl font-bold">Pay-As-You-Go</h3>
              <p class="text-gray-600 dark:text-gray-400">No fixed costs. Your bill reflects your exact usage.</p>
            </div>
            <div class="flex flex-col items-center text-center gap-4">
              <div
                class="w-16 h-16 bg-green-100 dark:bg-green-900/50 rounded-full flex items-center justify-center shrink-0">
                <i data-feather="shield" class="w-8 h-8 text-green-600 dark:text-green-400"></i></div>
              <h3 class="text-xl font-bold">Secure Authentication</h3>
              <p class="text-gray-600 dark:text-gray-400">JWT and API Key authentication methods.</p>
            </div>
            <div class="flex flex-col items-center text-center gap-4">
              <div
                class="w-16 h-16 bg-yellow-100 dark:bg-yellow-900/50 rounded-full flex items-center justify-center shrink-0">
                <i data-feather="bar-chart-2" class="w-8 h-8 text-yellow-600 dark:text-yellow-400"></i></div>
              <h3 class="text-xl font-bold">Dynamic Rate Limiting</h3>
              <p class="text-gray-600 dark:text-gray-400">Prevents abuse and ensures stable performance.</p>
            </div>
            <div class="flex flex-col items-center text-center gap-4">
              <div
                class="w-16 h-16 bg-red-100 dark:bg-red-900/50 rounded-full flex items-center justify-center shrink-0">
                <i data-feather="activity" class="w-8 h-8 text-red-600 dark:text-red-400"></i></div>
              <h3 class="text-xl font-bold">Real-time Monitoring</h3>
              <p class="text-gray-600 dark:text-gray-400">Track usage and manage spending from your dashboard.</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Pricing Section -->
      <section id="pricing" class="py-20 bg-gray-50 dark:bg-gray-900/50 flex justify-center items-center">
        <!-- FIX: Re-instated the standard container pattern -->
        <div class="container mx-auto px-6 flex flex-col items-center gap-12">
          <h2 class="text-3xl font-bold text-gray-800 dark:text-white sm:text-4xl text-center">Simple, Transparent
            Pricing</h2>
          <div
            class="max-w-lg w-full bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-8 border border-gray-200 dark:border-gray-700 ">
            <h3 class="text-2xl font-bold text-center">Pay-As-You-Go</h3>
            <div class="text-center my-8"><span class="text-5xl font-extrabold">$0</span><span
                class="text-xl font-medium text-gray-500">/ mo</span></div>
            <ul class="flex flex-col gap-2">
              <li class="flex items-center space-x-3"><i data-feather="check-circle"
                  class="w-6 h-6 text-green-500"></i><span>Access to all APIs</span></li>
              <li class="flex items-center space-x-3"><i data-feather="check-circle"
                  class="w-6 h-6 text-green-500"></i><span>Per-request billing</span></li>
              <li class="flex items-center space-x-3"><i data-feather="check-circle"
                  class="w-6 h-6 text-green-500"></i><span>Usage monitoring dashboard</span></li>
            </ul>
            <div class="h-5"></div>
            <router-link to="/auth"
              class="block w-full text-center bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 rounded-lg shadow-lg text-lg">Create
              Free Account</router-link>
          </div>
        </div>
      </section>
    </main>

  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router';
// No specific script logic is needed for this static landing page,
// but the setup block is here for consistency and future use.

</script>

<style scoped>
/* You can add component-specific styles here if needed */
</style>
