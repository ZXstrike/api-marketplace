<template>
  <div class="overflow-x-auto">
    <table class="w-full text-left">
      <thead>
        <tr class="border-b-2 border-gray-200 dark:border-gray-700 text-sm">
          <th class="p-3">Date</th>
          <th class="p-3">User</th>
          <th class="p-3">Method</th>
          <th class="p-3 text-right">Amount</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in rows" :key="row.id"
            class="border-b border-gray-100 dark:border-gray-700/50">
          <td class="p-3">{{ row.date }}</td>
          <td class="p-3 font-medium">{{ row.user }}</td>
          <td class="p-3">{{ row.method }}</td>
          <td class="p-3 text-right font-semibold text-green-600">
            {{ formatCurrency(row.amount) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { formatCurrency } from '@/helpers/currency';
defineProps({ rows: { type: Array, required: true } });
</script>