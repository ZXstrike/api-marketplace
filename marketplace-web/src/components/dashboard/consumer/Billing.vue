<template>
  <section class="flex flex-col gap-3">
    <h1 class="text-3xl font-bold">Wallet</h1>
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-3">
      <div class="dashboard-card lg:col-span-1 text-center flex flex-col justify-center gap-4">
        <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">Current Balance</h3>
        <p class="text-4xl font-bold mt-2 text-green-500">
          {{ new Intl.NumberFormat('id-ID', {
            style: 'currency', currency: 'IDR', minimumFractionDigits: 0
          }).format(userBalance || 0) }}
        </p>
        <p class="text-sm text-gray-500 mt-2">This balance is used for all Pay-As-You-Go API calls.</p>
        <button class="w-full mt-6 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2.5 rounded-lg shadow-md"
          @click="$emit('top-up')">
          Top Up Balance
        </button>
      </div>
      <!-- <div class="dashboard-card lg:col-span-2">
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
                <td class="p-3">Jul 20, 2025</td>
                <td class="p-3">Usage Charge: Geolocation API</td>
                <td class="p-3 text-right text-red-500">-{{ formatCurrency(2500) }}</td>
              </tr>
              <tr class="border-b border-gray-100 dark:border-gray-700/50">
                <td class="p-3">Jul 19, 2025</td>
                <td class="p-3">Usage Charge: Sentiment Analysis</td>
                <td class="p-3 text-right text-red-500">-{{ formatCurrency(5000) }}</td>
              </tr>
              <tr>
                <td class="p-3">Jul 1, 2025</td>
                <td class="p-3 font-semibold">Balance Top-up</td>
                <td class="p-3 text-right text-green-500">+{{ formatCurrency(50000) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div> -->
    </div>
  </section>
</template>

<script setup>
defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  userBalance: {
    type: Number,
    default: 0
  }
});

const emit = defineEmits(['top-up']);

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value);
};

const close = () => {
  emit('close');
};
</script>