<template>
  <section class="!space-y-3">
    <header>
      <h1 class="text-3xl font-bold">Admin Dashboard</h1>
      <p class="text-gray-600 dark:text-gray-400 !my-1">
        Marketplace performance and revenue overview.
      </p>
    </header>

    <!-- stat cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-3">
      <StatCard title="Total Revenue" :value="formatCurrency(totalRevenue)"
                icon="credit-card" color="green" />
      <StatCard title="Total Users" value="1,245" icon="users" color="blue" />
      <StatCard title="Published APIs" value="82" icon="package" color="yellow" />
      <StatCard title="Transactions (24h)" value="5,120" icon="activity" color="red" />
    </div>

    <!-- revenue chart -->
    <DashboardCard title="Revenue (Last 30 Days)">
      <div class="h-80"><canvas id="adminRevenueChart"></canvas></div>
    </DashboardCard>

    <!-- recent top-ups table -->
    <DashboardCard title="Recent Top-Ups">
      <TransactionTable :rows="recentTransactions" />
    </DashboardCard>
  </section>
</template>

<script setup>
import { onMounted, nextTick } from 'vue';
import Chart from 'chart.js/auto';
import StatCard        from '../StatCard.vue';
import DashboardCard   from './DashboardCard.vue';
import TransactionTable from './TransactionTable.vue';
import { formatCurrency } from '@/helpers/currency';

const totalRevenue = 157305000;
const recentTransactions = [
  { id: 1, date: '22 Jul 2025, 18:05', user: 'user_a@mail.com', method: 'QRIS', amount: 100000 },
  { id: 2, date: '22 Jul 2025, 17:45', user: 'user_b@mail.com', method: 'Credit Card', amount: 250000 },
  { id: 3, date: '22 Jul 2025, 16:30', user: 'user_c@mail.com', method: 'QRIS', amount: 50000 },
];

let chart = null;
onMounted(() => {
  nextTick(() => {
    const ctx = document.getElementById('adminRevenueChart');
    if (!ctx) return;
    if (chart) chart.destroy();
    chart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: ['Week 1', 'Week 2', 'Week 3', 'This Week'],
        datasets: [{
          label: 'Revenue (IDR)',
          data: [35000000, 42000000, 38000000, 48500000],
          borderColor: '#16A34A',
          backgroundColor: 'rgba(22, 163, 74, 0.1)',
          fill: true,
          tension: 0.4
        }]
      },
      options: { responsive: true, maintainAspectRatio: false }
    });
  });
});
</script>