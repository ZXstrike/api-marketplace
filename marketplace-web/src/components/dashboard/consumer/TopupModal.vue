<template>
  <!-- Overlay Background -->
  <div v-if="isOpen" @click.self="$emit('close')"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">

    <!-- Dialog Content -->
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl w-full max-w-lg p-8">

      <!-- Dialog Header -->
      <div class="flex justify-between items-center !mb-3">
        <h2 class="text-2xl font-bold">Top Up Your Balance</h2>
        <button @click="$emit('close')" class="p-1 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
          <i data-feather="x" class="w-6 h-6"></i>
        </button>
      </div>

      <!-- Form Content -->
      <form @submit.prevent="handleTopup">
        <div class="space-y-6">
          <!-- Top-up Amount Selection -->
          <div>
            <label class="form-label !mb-2">Select Amount (in IDR)</label>
            <div class="mt-2 grid grid-cols-3 gap-4">
              <button v-for="amount in presetAmounts" :key="amount" type="button" @click="selectedAmount = amount"
                :class="[
                  'py-3 px-4 rounded-lg border text-center font-semibold transition-all duration-200 ease-in-out',
                  selectedAmount === amount
                    ? 'bg-blue-600 text-white border-blue-600 ring-2 ring-blue-500 ring-offset-2 dark:ring-offset-gray-800'
                    : 'bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-800 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-600'
                ]">
                Rp {{ formatCurrency(amount) }}
              </button>
              <input type="text" inputmode="numeric" v-model.number="selectedAmount"
                @input="selectedAmount = $event.target.value.replace(/\D/g, '')" placeholder="Other Amount"
                class="form-input text-center col-span-3">
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="pt-4 border-t border-gray-200 dark:border-gray-700">
            <button type="submit" class="w-full btn-primary text-lg">
              Proceed to Payment
            </button>
          </div>
        </div>
      </form>

    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import apiClient from '@/services/apiClient.js';

defineProps({
  isOpen: {
    type: Boolean,
    required: true,
  }
});

const emit = defineEmits(['close', 'topup-success']);

// Preset amounts in IDR
const presetAmounts = [50000, 100000, 250000];
const selectedAmount = ref(50000);

// Function to format numbers into Indonesian currency format
const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID').format(value);
};

const handleTopup = async () => {
  if (selectedAmount.value <= 0) {
    alert("Please enter a valid amount.");
    return;
  }
  console.log(`Starting top-up of Rp ${formatCurrency(selectedAmount.value)}`);
  try {
    const response = await apiClient.put('/billing/update-balance', {
      amount: selectedAmount.value
    });
    if (response.status !== 200) {
      throw new Error("Failed to top up balance");
    }
    // Assuming the response contains the updated balance
    const data = await response.json();
    console.log('Top-up successful:', data);
    // Emit success event to parent component
    
    console.log(`Top-up of Rp ${selectedAmount.value} successful!`);
    alert(`Top-up of Rp ${formatCurrency(selectedAmount.value)} successful!`);
  } catch (error) {
    console.error("Top-up failed:", error);
    alert("Top-up failed. Please try again.");
  }
  emit('topup-success');
  emit('close');
};
</script>

<style scoped></style>
