<template>
<div>
    <h2 class="text-lg font-semibold pb-4">API Endpoints ({{ versionString }})</h2>
    <nav v-if="endpoints.length > 0" class="flex flex-col gap-2">
        <a v-for="endpoint in endpoints" :key="endpoint.id" 
                href="#" @click.prevent="$emit('select', endpoint)"
                :class="{'active': selectedEndpoint && selectedEndpoint.id === endpoint.id}" 
                class="endpoint-link flex items-center gap-3">
            <span :class="getMethodColor(endpoint.http_method)" class="http-method">{{ endpoint.http_method }}</span>
            <span class="font-mono">{{ endpoint.path }}</span>
        </a>
    </nav>
    <div v-else class="text-sm text-gray-500">
            No endpoints found for this version.
    </div>
</div>
</template>

<script setup>
defineProps({
  endpoints: { type: Array, required: true },
  selectedEndpoint: { type: Object, default: null },
  versionString: { type: String, default: '' }
});

defineEmits(['select']);

const getMethodColor = (method) => {
  const colors = {
    'GET': 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
    'POST': 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
    'PUT': 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
    'DELETE': 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
  };
  return colors[method.toUpperCase()] || 'bg-gray-100 text-gray-800';
};
</script>

<style scoped>
</style>