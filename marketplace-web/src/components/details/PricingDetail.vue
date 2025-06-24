<template>
    <div class="flex items-center justify-center min-h-full">
        <div class="max-w-4xl mx-auto flex flex-col gap-8">
            <div class="text-center">
                <h1 class="text-3xl font-bold">Pricing for Version {{ version?.version_string }}</h1>
            </div>
            <div class="dashboard-card text-center">
                <h2 class="text-xl font-semibold text-gray-500 dark:text-gray-400">Price per Call</h2>
                <template v-if="version">
                    <p v-if="version.price_per_call > 0" class="my-4">
                        <span class="text-5xl font-extrabold text-blue-600">${{ version.price_per_call }}</span>
                    </p>
                    <p v-else class="my-4 text-5xl font-extrabold text-green-500">Free</p>
                </template>
                <p v-else class="my-4 text-gray-500">Pricing not available.</p>
            </div>
        </div>
    </div>
</template>

<script setup>
defineProps({
    version: {
        type: Object,
        default: null
    }
});
</script>

<style scoped></style>