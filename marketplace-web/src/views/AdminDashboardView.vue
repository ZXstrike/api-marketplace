<template>
  <div class="w-full flex justify-center">
    <main class="container mx-auto px-6 py-8 md:py-12">
      <div class="lg:flex lg:space-x-8">

        <!-- Menggunakan komponen Sidebar Admin yang baru -->
        <AdminSidebar :active-section="activeSection" @navigate="handleNavigation" />

        <!-- Konten Dasbor Admin -->
        <main class="flex-1 mt-8 lg:mt-0">
          <section v-if="activeSection === 'revenue'" class="space-y-8">
            <div>
              <h1 class="text-3xl font-bold">Admin Dashboard</h1>
              <p class="text-gray-600 dark:text-gray-400 mt-1">Marketplace performance and revenue overview.</p>
            </div>

            <!-- Kartu Statistik Utama -->
            <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-6">
              <StatCard title="Total Revenue" :value="formatCurrency(totalRevenue)" icon="credit-card" color="green" />
              <StatCard title="Total Users" value="1,245" icon="users" color="blue" />
              <StatCard title="Published APIs" value="82" icon="package" color="yellow" />
              <StatCard title="Transactions (24h)" value="5,120" icon="activity" color="red" />
            </div>

            <!-- Grafik Pendapatan -->
            <div class="dashboard-card">
              <h3 class="text-xl font-bold mb-4">Revenue (Last 30 Days)</h3>
              <div class="h-80"><canvas id="adminRevenueChart"></canvas></div>
            </div>

            <!-- Tabel Transaksi Terbaru -->
            <div class="dashboard-card">
              <h3 class="text-xl font-bold mb-4">Recent Top-Ups</h3>
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
                    <tr v-for="tx in recentTransactions" :key="tx.id"
                      class="border-b border-gray-100 dark:border-gray-700/50">
                      <td class="p-3">{{ tx.date }}</td>
                      <td class="p-3 font-medium">{{ tx.user }}</td>
                      <td class="p-3">{{ tx.method }}</td>
                      <td class="p-3 text-right font-semibold text-green-600">{{ formatCurrency(tx.amount) }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </section>
          <!-- Placeholder for other admin sections -->
          <section v-if="activeSection === 'users'">
             <h1 class="text-3xl font-bold">User Management</h1>
             <p class="text-center py-16 text-gray-500">User management page coming soon.</p>
          </section>
          <section v-if="activeSection === 'apis'">
             <h1 class="text-3xl font-bold">API Management</h1>
             <p class="text-center py-16 text-gray-500">API management page coming soon.</p>
          </section>
        </main>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';
import StatCard from '@/components/dashboard/StatCard.vue';
import AdminSidebar from '@/components/dashboard/admin/Sidebar.vue'; // Impor komponen baru
import Chart from 'chart.js/auto';

const activeSection = ref('revenue');
let revenueChartInstance = null;

// --- Data Mockup (akan diganti dengan panggilan API) ---
const totalRevenue = ref(157305000); // dalam IDR
const recentTransactions = ref([
  { id: 1, date: '22 Jul 2025, 18:05', user: 'user_a@mail.com', method: 'QRIS', amount: 100000 },
  { id: 2, date: '22 Jul 2025, 17:45', user: 'user_b@mail.com', method: 'Credit Card', amount: 250000 },
  { id: 3, date: '22 Jul 2025, 16:30', user: 'user_c@mail.com', method: 'QRIS', amount: 50000 },
]);

// --- Metode ---
const handleNavigation = (section) => {
  activeSection.value = section;
};

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value);
};

const renderRevenueChart = () => {
  const ctx = document.getElementById('adminRevenueChart');
  if (!ctx) return;
  if (revenueChartInstance) revenueChartInstance.destroy();

  revenueChartInstance = new Chart(ctx, {
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
};

onMounted(() => {
  nextTick(() => {
    renderRevenueChart();
  });
});
</script>
