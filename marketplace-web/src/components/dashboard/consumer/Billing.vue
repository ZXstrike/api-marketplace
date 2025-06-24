<template>
  <section class="space-y-6">
    <h1 class="text-3xl font-bold">Billing & Usage</h1>
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="dashboard-card lg:col-span-1">
        <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">Current Balance</h3>
        <p class="text-4xl font-bold mt-2 text-green-500">$42.50</p>
        <p class="text-sm text-gray-500 mt-2">This balance is used for all Pay-As-You-Go API calls.</p>
        <button
          class="w-full mt-6 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2.5 rounded-lg shadow-md">Top-Up
          Balance</button>
      </div>
      <div class="dashboard-card lg:col-span-2">
        <h3 class="text-xl font-bold mb-4">Transaction History</h3>
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead>
              <tr class="border-b-2 border-gray-200 dark:border-gray-700">
                <th class="p-3">Date</th>
                <th class="p-3">Description</th>
                <th class="p-3 text-right">Amount</th>
              </tr>
            </thead>
            <tbody>
              <tr class="border-b border-gray-100 dark:border-gray-700/50">
                <td class="p-3">Jun 20, 2025</td>
                <td class="p-3">Usage Charge: Geolocation API</td>
                <td class="p-3 text-right text-red-500">-$2.50</td>
              </tr>
              <tr class="border-b border-gray-100 dark:border-gray-700/50">
                <td class="p-3">Jun 19, 2025</td>
                <td class="p-3">Usage Charge: Sentiment Analysis</td>
                <td class="p-3 text-right text-red-500">-$5.00</td>
              </tr>
              <tr>
                <td class="p-3">Jun 1, 2025</td>
                <td class="p-3 font-semibold">Balance Top-up</td>
                <td class="p-3 text-right text-green-500">+$50.00</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
// No script logic needed for this static component yet
</script>