<template>
    <section class="flex flex-col gap-3">
        <div class="flex flex-col md:flex-row justify-between md:items-center gap-4">
            <div>
                <h1 class="text-3xl font-bold">My API Products</h1>
                <p class="text-gray-600 dark:text-gray-400 mt-1">Manage, update, and create new API products.</p>
            </div>
            <router-link to="/dashboard/create-new-api"
                class="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-5 py-2.5 rounded-lg shadow-md flex items-center justify-center space-x-2">
                <i data-feather="plus" class="w-5 h-5"></i>
                <span>Create New API</span>
            </router-link>
        </div>

        <div class="dashboard-card overflow-x-auto">
            <table class="w-full text-left">
                <thead class="border-b-2 border-gray-200 dark:border-gray-700">
                    <tr class="text-sm text-gray-600 dark:text-gray-400">
                        <th class="p-3">API Name</th>
                        <th class="p-3">Status</th>
                        <th class="p-3">Subscribers</th>
                        <th class="p-3">Revenue (Mo)</th>
                        <th class="p-3">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="api in apis" :key="api.id" class="border-b border-gray-100 dark:border-gray-700/50">
                        <td class="p-3 font-semibold">{{ api.name }}</td>
                        <td class="p-3">
                            <span
                                :class="api.status === 'Published' ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300' : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'"
                                class="text-xs font-medium mr-2 px-2.5 py-0.5 rounded-full">
                                {{ api.status }}
                            </span>
                        </td>
                        <td class="p-3 font-medium">{{ api.subscribers }}</td>
                        <td class="p-3 font-medium text-green-600">${{ api.monthlyRevenue.toFixed(2) }}</td>
                        <td class="p-3 flex space-x-2">
                            <button @click="emitEditAPI(api.id)" class="p-2 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md"><i
                                    data-feather="edit-2" class="w-4 h-4"></i></button>
                            <button class="p-2 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md"><i
                                    data-feather="bar-chart-2" class="w-4 h-4"></i></button>
                            <button @click="emitDelete(api.id)"
                                class="p-2 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md text-red-500"><i
                                    data-feather="trash-2" class="w-4 h-4"></i></button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </section>
</template>

<script setup>
import { onMounted, onUpdated } from 'vue';
import { RouterLink } from 'vue-router';

defineProps({
    apis: {
        type: Array,
        required: true
    }
});

const emit = defineEmits(['delete-api', 'edit-api']);

const emitDelete = (id) => {
    if (confirm('Are you sure you want to delete this API? This action cannot be undone.')) {
        emit('delete-api', id);
    }
};

const emitEditAPI = (api) => {
    emit('edit-api', api);
};

onMounted(() => {
    if (window.feather) {
        window.feather.replace();
    }
});

onUpdated(() => {
    if (window.feather) {
        window.feather.replace();
    }
});
</script>

<style scoped></style>