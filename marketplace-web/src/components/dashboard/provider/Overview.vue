<template>
    <section class="flex flex-col gap-3">
        <div>
            <h1 class="text-3xl font-bold">Welcome back, {{ authStore.user?.name || 'Seller' }}!</h1>
            <p class="text-gray-600 dark:text-gray-400">An overview of your API performance and earnings.</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-3">
            <div class="dashboard-card">
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">Total Revenue</h3><i
                        data-feather="dollar-sign" class="text-green-500"></i>
                </div>
                <p class="text-3xl font-bold mt-2">$1,230.75</p>
            </div>
            <div class="dashboard-card">
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">Total Subscribers</h3><i
                        data-feather="users" class="text-blue-500"></i>
                </div>
                <p class="text-3xl font-bold mt-2">{{ totalSubs }}</p>
            </div>
            <div class="dashboard-card">
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">APIs Published</h3><i
                        data-feather="package" class="text-yellow-500"></i>
                </div>
                <p class="text-3xl font-bold mt-2">{{ publishedApiCount }}</p>
            </div>
            <div class="dashboard-card">
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">Requests (Month)</h3><i
                        data-feather="activity" class="text-red-500"></i>
                </div>
                <p class="text-3xl font-bold mt-2">8.2M</p>
            </div>
        </div>
        <div class="dashboard-card">
            <h3 class="text-xl font-bold mb-4">Revenue (Last 30 Days)</h3>
            <div class="h-80"><canvas id="revenueChart"></canvas></div>
        </div>
    </section>
</template>

<script setup>
import { onMounted, onBeforeUnmount, computed, nextTick } from 'vue';
import { useAuthStore } from '@/store/auth';
import Chart from 'chart.js/auto';

// Get the auth store
const authStore = useAuthStore();

const props = defineProps({
    apis: {
        type: Array,
        required: true
    },
    totalSubs:{
        type: Number,
        default: 0
    },
});

let revenueChartInstance = null;

const publishedApiCount = computed(() => {
    return props.apis.filter(api => api.status === 'Published').length;
});

const renderRevenueChart = () => {
    const ctx = document.getElementById('revenueChart');
    if (!ctx) return;

    // Destroy the old chart instance if it exists
    if (revenueChartInstance) {
        revenueChartInstance.destroy();
    }

    revenueChartInstance = new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['Week 1', 'Week 2', 'Week 3', 'This Week'],
            datasets: [{
                label: 'Revenue ($)',
                data: [250, 310, 280, 450],
                borderColor: '#3b82f6',
                backgroundColor: 'rgba(59, 130, 246, 0.1)',
                tension: 0.3
            }]
        },
        options: { responsive: true, maintainAspectRatio: false }
    });
};

onMounted(() => {
    nextTick(() => {
        renderRevenueChart();
        // Initialize feather icons
        if (window.feather) {
            window.feather.replace();
        }
    });
});

onBeforeUnmount(() => {
    if (revenueChartInstance) {
        revenueChartInstance.destroy();
    }
});
</script>

<style scoped></style>