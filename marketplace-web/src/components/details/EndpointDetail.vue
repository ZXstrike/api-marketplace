<template>
  <div v-if="endpoint" class="grid grid-cols-1 md:grid-cols-2 gap-8">
    <div class="space-y-6">
      <h1 class="text-2xl font-bold font-mono">{{ endpoint.http_method }} {{ endpoint.path }}</h1>
      <h2 class="text-xl font-semibold">Documentation</h2>
      <p class="text-gray-600 dark:text-gray-400">{{ documentationText }}</p>
    </div>

    <div class="sticky top-28 space-y-8">
    <div class="bg-gray-800 dark:bg-black rounded-xl overflow-hidden shadow-2xl">
      <div class="bg-gray-900 p-3 border-b border-gray-700">
        <span class="text-sm font-semibold text-gray-300">REQUEST EXAMPLE</span>
      </div>
      <pre class="p-4 overflow-x-auto"><code class="language-bash hljs">curl --request {{ endpoint.http_method }} \
--url '{{ baseUrl }}{{ endpoint.path }}' \
--header 'Authorization: Bearer YOUR_API_KEY'</code></pre>
    </div>
      
    <div v-if="formattedOutput" class="bg-gray-800 dark:bg-black rounded-xl overflow-hidden shadow-2xl">
      <div class="bg-gray-900 p-3 border-b border-gray-700">
        <span class="text-sm font-semibold text-gray-300">EXAMPLE OUTPUT</span>
      </div>
      <pre class="p-4 overflow-x-auto"><code class="language-json hljs">{{ formattedOutput }}</code></pre>
    </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  endpoint: { type: Object, required: true },
  baseUrl: { type: String, required: true }
});

const isJsonOutput = computed(() => {
    return props.endpoint?.documentation?.startsWith('Output : ');
});

const documentationText = computed(() => {
    return isJsonOutput.value ? 'See example output.' : props.endpoint.documentation;
});

const formattedOutput = computed(() => {
  if (!isJsonOutput.value) return '';
  try {
    const docJson = JSON.parse(props.endpoint.documentation.replace('Output : ', ''));
    return JSON.stringify(docJson, null, 2);
  } catch (e) {
    // If it's not valid JSON after all, don't show the output block
    return ''; 
  }
});
</script>