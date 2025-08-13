<template>
  <section class="!space-y-3">
    <header>
      <h1 class="text-3xl font-bold">Admin Dashboard</h1>
      <p class="text-gray-600 dark:text-gray-400 !my-1">
        Marketplace performance and revenue overview.
      </p>
    </header>

    <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-3">
      <StatCard title="Total Revenue" :value="formattedTotalRevenue" icon="credit-card" color="green" />
      <StatCard title="Total Users" :value="totalUsers" icon="users" color="blue" />
      <StatCard title="Published APIs" :value="publishedApis" icon="package"color="purple" />
      <StatCard title="Transactions (24H)" :value="transactions24H" icon="clock" color="yellow" />
    </div>

    <DashboardCard title="Recent Transactions" >
      <TransactionTable :rows="recentTransactions" />
    </DashboardCard>
  </section>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue';
import Chart from 'chart.js/auto';
import StatCard from '../StatCard.vue';
import DashboardCard from './DashboardCard.vue';
import TransactionTable from './TransactionTable.vue';
import { formatCurrency } from '@/helpers/currency';
import apiClient from '@/services/apiClient';

const totalRevenue = ref(0);
const totalUsers = ref(0);
const publishedApis = ref(0);
const transactions24H = ref(0);
const formattedTotalRevenue = computed(() => formatCurrency(totalRevenue.value));

const recentTransactions = ref([]);
// [
//   { id: 1, date: '22 Jul 2025, 18:05', user: 'user_a@mail.com', method: 'QRIS', amount: 100000 },
//   { id: 2, date: '22 Jul 2025, 17:45', user: 'user_b@mail.com', method: 'Credit Card', amount: 250000 },
//   { id: 3, date: '22 Jul 2025, 16:30', user: 'user_c@mail.com', method: 'QRIS', amount: 50000 },
// ];

const fetchBillingInfo = async () => {
  try {
   const response = await apiClient.get('/billing/info?wallet_type=admin');
    const data = await response.json();
    // Proses data billing jika diperlukan
    console.log('User Billing Info:', data);
    totalRevenue.value = data.balance || 0;
    console.log('User Balance:', totalRevenue.value);
  } catch (e) {
    console.error('Failed to fetch billing info', e);
    totalRevenue.value = 0;
  }
};

const fetchAdminData = async () => {
  try {
  const response = await apiClient.get('/api/admin/data');
  const data = await response.json();
  console.log('Admin Data:', data);

  totalUsers.value = data?.user_count ?? 0;
  publishedApis.value = data?.api_count ?? 0;
  transactions24H.value = data?.transaction_24h ?? 0;

  const formatDate = (raw) => {
    if (!raw) return '-';
    // Try direct parse first
    let d = new Date(raw);
    if (isNaN(d)) {
    // Handle common "YYYY-MM-DD HH:MM:SS" -> "YYYY-MM-DDTHH:MM:SS"
    d = new Date(raw.replace(' ', 'T'));
    }
    return isNaN(d) ? raw : d.toLocaleString();
  };

  recentTransactions.value = (data?.recent_topups || []).map(t => ({
    id: t.id,
    date: formatDate(t.created_at),
    user: t.user?.email || '-',
    amount: t.amount ?? 0,
  }));
  } catch (e) {
  console.error('Failed to fetch admin data', e);
  }
};

let chart = null;
onMounted(() => {
  fetchBillingInfo();
  fetchAdminData();
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
          backgroundColor: 'rgba(22,163,74,0.1)',
          fill: true,
          tension: 0.4
        }]
      },
      options: { responsive: true, maintainAspectRatio: false }
    });
  });
});
</script>