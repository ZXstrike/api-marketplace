<template>
  <div class="flex min-h-screen w-full  justify-center">
    <main class="container mx-auto px-6 py-8 md:py-12">
      <div v-if="loading" class="text-center py-16">
        <p>Loading API details...</p>
      </div>

      <div v-if="!loading && api">
        <ApiHeader :api="api" @subscribe="handleSubscribe" class="mb-8" />

        <ApiNavigation
          :versions="api.versions"
          v-model:selectedVersionId="selectedVersionId"
          v-model:activeTab="activeTab"
          class="mb-8"
        />

        <div class="py-3">
          <section v-show="activeTab === 'endpoints'" class="grid grid-cols-1 lg:grid-cols-3 gap-4">
            <aside class="lg:col-span-1">
              <div class="sticky top-28">
                   <EndpointList 
                     :endpoints="endpoints"
                     :version-string="selectedVersion?.version_string"
                     :selected-endpoint="selectedEndpoint"
                     @select="selectedEndpoint = $event"
                   />
              </div>
            </aside>
            <main class="lg:col-span-2">
              <EndpointDetail 
                v-if="selectedEndpoint"
                :endpoint="selectedEndpoint"
                :base-url="dynamicBaseUrl"
              />
              <div v-else class="text-center py-20 text-gray-500">
                <p>Select an endpoint on the left to see its details.</p>
              </div>
            </main>
          </section>

          <section v-show="activeTab === 'pricing'">
            <PricingDetail :version="selectedVersion" />
          </section>
        </div>
      </div>
      
      <div v-if="!loading && !api" class="text-center py-16">
          <p class="text-red-500">Failed to load API details.</p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed, nextTick } from 'vue';
import { useRoute } from 'vue-router';

// Import apiClient dan child components
import apiClient from '@/services/apiClient.js';
import ApiHeader from '@/components/details/ApiHeader.vue';
import ApiNavigation from '@/components/details/ApiNavigation.vue';
import EndpointList from '@/components/details/EndpointList.vue';
import EndpointDetail from '@/components/details/EndpointDetail.vue';
import PricingDetail from '@/components/details/PricingDetail.vue';
import router from '@/router/index';

const route = useRoute();
const api = ref(null);
const endpoints = ref([]);
const loading = ref(true);
const activeTab = ref('endpoints');
const selectedVersionId = ref(null);
const selectedEndpoint = ref(null);

// --- Computed Properties ---
const selectedVersion = computed(() => {
  if (!api.value || !selectedVersionId.value) return null;
  return api.value.versions.find(v => v.id === selectedVersionId.value);
});

const dynamicBaseUrl = computed(() => {
  if (!api.value || !api.value.provider) return '';
  return `https://${api.value.provider.username}.api.zxsttm.tech/${api.value.slug}`;
});


// --- Methods ---
const fetchApiDetails = async (apiId) => {
  try {
    // FIX: Menggunakan apiClient.get, yang sudah menangani base URL dan error dasar
    const response = await apiClient.get(`/api/${apiId}`);
    const data = await response.json(); // Mengambil data JSON dari respons
    api.value = data;

    if (api.value?.versions?.length > 0) {
      api.value.versions.sort((a, b) => b.version_string.localeCompare(a.version_string, undefined, { numeric: true }));
      selectedVersionId.value = api.value.versions[0].id;
    }
  } catch (error) {
    console.error("Error fetching API details:", error);
    api.value = null;
  }
};

const fetchEndpoints = async (versionId) => {
  if (!versionId) return;
  endpoints.value = [];
  selectedEndpoint.value = null;

  try {
    // FIX: Menggunakan apiClient.get untuk mengambil endpoint
    const response = await apiClient.get(`/api/api-endpoints/${versionId}`);
    const data = await response.json();
    endpoints.value = data.reverse(); // Reverse the order to show the latest first

    if (endpoints.value.length > 0) {
      selectedEndpoint.value = endpoints.value[0];
    }
  } catch (error) {
    console.error(`Error fetching endpoints for version ${versionId}:`, error);
    endpoints.value = [];
  }
};

const handleSubscribe = async () => {
  if (!selectedVersionId.value) {
    alert("Please select an API version to subscribe.");
    return;
  }

  try {
    const response = await apiClient.post('/subscriptions/subscribe', {
      api_version_id: selectedVersionId.value
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || `Failed to subscribe. Status: ${response.status}`);
    }

    alert(data.message || "Successfully subscribed!");
    router.push('/dashboard'); // Redirect to dashboard after subscription
    // You might want to refresh the API details here to update the UI
    // await fetchApiDetails(route.params.id);
  } catch (error) {
    console.error("Error subscribing to API:", error);
    alert(`Subscription failed: ${error.message}`);
  }
};

// --- Lifecycle & Watchers ---
onMounted(async () => {
  loading.value = true;
  await fetchApiDetails(route.params.id);
  loading.value = false;
  await nextTick();
  if(window.feather) window.feather.replace();
});

watch(selectedVersionId, (newId, oldId) => {
  if (newId && newId !== oldId) {
    fetchEndpoints(newId);
  }
});

watch(endpoints, async () => {
    await nextTick();
    if (window.hljs) window.hljs.highlightAll();
})
</script>
