<template>
    <section class="flex flex-col gap-3">
        <div>
            <h1 class="text-3xl font-bold">Welcome back, {{ authStore.user?.name || 'User' }}!</h1>
            <p class=" text-gray-600 dark:text-gray-400  ">Here's a summary of your account activity.</p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-3">
            <div class="dashboard-card">
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">Current Balance</h3><i
                        data-feather="credit-card" class="text-green-500"></i>
                </div>
                <p class="text-2xl font-bold mt-2">{{ new Intl.NumberFormat('id-ID', {
                    style: 'currency', currency:
                        'IDR', minimumFractionDigits: 0
                }).format(userBalance || 0) }}</p>
            </div>
            <div class="dashboard-card">
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">Monthly Cost</h3><i
                        data-feather="trending-up" class="text-blue-500"></i>
                </div>
                <p class="text-2xl font-bold mt-2">{{ new Intl.NumberFormat('id-ID', {
                    style: 'currency', currency:
                        'IDR', minimumFractionDigits: 0
                }).format(monthlyCost || 0) }}</p>
            </div>
            <div class="dashboard-card">
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">Active Subscriptions</h3><i
                        data-feather="grid" class="text-yellow-500"></i>
                </div>
                <p class="text-2xl font-bold mt-2">{{ subscriptions}}</p>
            </div>
            <div class="dashboard-card">
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-medium text-gray-500 dark:text-gray-400">Requests (Month)</h3><i
                        data-feather="bar-chart-2" class="text-red-500"></i>
                </div>
                <p class="text-2xl font-bold mt-2">{{ new Intl.NumberFormat('id-ID', {
                    style: 'decimal',
                    minimumFractionDigits: 0
                }).format(requestsThisMonth || 0) }}</p>
            </div>
        </div>

        <div class="dashboard-card">
            <h3 class="text-xl font-bold mb-4">Daily Usage (Last 7 Days)</h3>
            <div class="h-80"><canvas id="usageChart"></canvas></div>
        </div>
    </section>
</template>

<script setup>
import { onMounted, onBeforeUnmount, nextTick, watch } from 'vue';
import { useAuthStore } from '@/store/auth';
import Chart from 'chart.js/auto';

// Get the auth store
const authStore = useAuthStore();

// Define props to receive data from the parent
const props = defineProps({
    userBalance: {
        type: Number,
        default: 0
    },
    monthlyCost: {
        type: Number,
        default: 0
    },
    requestsThisMonth: {
        type: Number,
        default: 0
    },
    subscriptions: {
        type: Number,
        default: 0
    }, 
    weeklyUsage: {
        type: Array,
        default: () => []
    }
});

let usageChartInstance = null;

const renderUsageChart = () => {
    const ctx = document.getElementById('usageChart');
    if (!ctx) return;

    if (usageChartInstance) {
        usageChartInstance.destroy();
    }

    const labels = props.weeklyUsage.map(entry => entry.date);
    const data = props.weeklyUsage.map(entry => entry.count);

    usageChartInstance = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'API Requests',
                data: data,
                borderColor: '#3b82f6',
                backgroundColor: 'rgba(59, 130, 246, 0.1)',
                tension: 0.3
            }]
        },
        options: { responsive: true, maintainAspectRatio: false }
    });
};

watch(() => props.weeklyUsage, (newVal) => {
    if (newVal && newVal.length > 0) {
        renderUsageChart();
    }
}, { deep: true });

onMounted(() => {
    nextTick(() => {
        renderUsageChart();
        // Initialize feather icons
        if (window.feather) {
            window.feather.replace();
        }
    });
});

onBeforeUnmount(() => {
    if (usageChartInstance) {
        usageChartInstance.destroy();
    }
});
</script>