<template>
  <div class="w-full flex justify-center">
    <main class="container mx-auto px-6 py-8 md:py-12">
      <div v-if="!loading && api" class="max-w-3xl mx-auto">
        <router-link to="/dashboard/provider" class="text-blue-600 hover:underline flex items-center mb-4">
          <i data-feather="arrow-left" class="w-4 h-4 mr-2"></i>
          Back to Dashboard
        </router-link>
        <h1 class="text-3xl font-bold">Edit API: {{ api.name }}</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">Update the details, pricing, and endpoints for your API.</p>
      </div>
      
      <div v-if="loading" class="text-center py-16">Loading API data...</div>

      <form v-if="!loading && api" @submit.prevent="handleUpdateApi" class="space-y-8 max-w-3xl mx-auto mt-8">
        <!-- Bagian Detail & Harga API -->
        <div class="dashboard-card">
          <h2 class="section-heading">API Details & Pricing</h2>
          <div class="space-y-6">
            <div>
              <label for="api-name" class="form-label">API Name</label>
              <input type="text" id="api-name" v-model="apiData.name" class="form-input" required>
            </div>
            <div>
              <label for="api-description" class="form-label">Description</label>
              <textarea id="api-description" rows="3" v-model="apiData.description" class="form-input" required></textarea>
            </div>
            <div>
              <label for="base-url" class="form-label">Your Backend Base URL</label>
              <input type="url" id="base-url" v-model="apiData.base_url" class="form-input" required>
            </div>
            <div>
              <label for="price-per-call" class="form-label">Price per Call (in USD)</label>
              <input type="number" id="price-per-call" v-model.number="apiData.price_per_call" step="0.0001" min="0" class="form-input" required>
            </div>
            <div>
              <label class="form-label">Categories</label>
              <div v-if="availableCategories.length > 0" class="mt-2 grid grid-cols-2 md:grid-cols-3 gap-x-4 gap-y-2 border p-4 rounded-md">
                <label v-for="category in availableCategories" :key="category.id" class="flex items-center cursor-pointer">
                  <input type="checkbox" :value="category.slug" v-model="apiData.categories" class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500">
                  <span class="ml-2 text-sm">{{ category.name }}</span>
                </label>
              </div>
            </div>
          </div>
        </div>
        
        <div class="flex justify-end">
          <button type="submit" :disabled="isSaving" class="btn-primary">
            {{ isSaving ? 'Saving...' : 'Save API Details' }}
          </button>
        </div>
      </form>

      <!-- Bagian Manajemen Endpoint -->
      <div v-if="!loading && api" class="space-y-8 max-w-3xl mx-auto mt-8">
          <div class="dashboard-card">
            <h2 class="section-heading">Manage Endpoints for {{ selectedVersion?.version_string }}</h2>
            <!-- Daftar Endpoint yang Sudah Ada -->
            <div class="space-y-4">
                 <div v-for="endpoint in endpoints" :key="endpoint.id" class="p-3 bg-gray-50 dark:bg-gray-800/50 rounded-lg border dark:border-gray-700">
                    <div class="flex justify-between items-center">
                        <p class="font-mono text-sm"><strong>{{ endpoint.http_method }}</strong> {{ endpoint.path }}</p>
                        <button @click="deleteEndpoint(endpoint.id)" class="text-red-500 p-1 rounded-full hover:bg-red-100"><i data-feather="trash-2" class="w-5 h-5"></i></button>
                    </div>
                 </div>
                 <div v-if="endpoints.length === 0" class="text-center py-8 text-gray-500"><p>No endpoints defined for this version.</p></div>
            </div>
            <!-- Form Tambah Endpoint Baru -->
            <div class="mt-6 pt-6 border-t border-dashed">
                <h3 class="font-semibold mb-4">Add New Endpoint</h3>
                <form @submit.prevent="addEndpoint" class="space-y-4">
                     <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                        <div><label class="form-label">Method</label><select v-model="newEndpoint.http_method" class="form-input"><option>GET</option><option>POST</option><option>PUT</option><option>DELETE</option></select></div>
                        <div><label class="form-label">Path</label><input type="text" v-model="newEndpoint.path" placeholder="/users/{id}" class="form-input"></div>
                    </div>
                    <div><label class="form-label">Documentation</label><textarea v-model="newEndpoint.documentation" rows="3" class="form-input font-mono"></textarea></div>
                    <div class="text-right">
                        <button type="submit" class="btn-secondary">Add Endpoint</button>
                    </div>
                </form>
            </div>
          </div>
      </div>

    </main>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import apiClient from '@/services/apiClient.js';

const route = useRoute();

const loading = ref(true);
const isSaving = ref(false);
const api = ref(null);
const availableCategories = ref([]);
const endpoints = ref([]);
const newEndpoint = ref({ http_method: 'GET', path: '', documentation: '' });

// Gunakan reactive untuk data form agar mudah di-update
const apiData = reactive({
    name: '',
    description: '',
    base_url: '',
    price_per_call: 0,
    categories: [],
});

// --- Computed Properties ---
const selectedVersion = computed(() => {
    // Asumsikan kita mengedit versi pertama untuk saat ini
    return api.value?.versions?.[0];
});

// --- Methods ---
const fetchInitialData = async () => {
  const apiId = route.params.id;
  try {
    // Ambil detail API
    const apiResponse = await apiClient.get(`/api/${apiId}`);
    api.value = await apiResponse.json();

    // Isi form dengan data yang ada
    Object.assign(apiData, {
        name: api.value.name,
        description: api.value.description,
        base_url: api.value.base_url,
        price_per_call: api.value.versions?.[0]?.price_per_call || 0,
        categories: api.value.categories?.map(c => c.slug) || [],
    });
    
    // Ambil semua kategori yang tersedia
    const catResponse = await apiClient.get('/api/categories');
    availableCategories.value = await catResponse.json();

    // Ambil endpoint untuk versi yang ada
    if (selectedVersion.value) {
        const endpointResponse = await apiClient.get(`/api/api-endpoints/${selectedVersion.value.id}`);
        endpoints.value = await endpointResponse.json();
    }

  } catch (error) {
    console.error("Failed to load API data:", error);
    alert("Could not load data for editing.");
  } finally {
    loading.value = false;
  }
};

const handleUpdateApi = async () => {
    isSaving.value = true;
    try {
        const apiId = route.params.id;
        // Panggil endpoint PUT untuk update
        await apiClient.put(`/api/update/${apiId}`, apiData);
        alert('API details updated successfully!');
    } catch (error) {
        console.error("Error updating API:", error);
        alert("Failed to update API details.");
    } finally {
        isSaving.value = false;
    }
};

const addEndpoint = async () => {
    if (!newEndpoint.value.path || !selectedVersion.value) return;
    try {
        const payload = {
            api_version_id: selectedVersion.value.id,
            endpoints: [newEndpoint.value],
        };
        await apiClient.post('/api/create-endpoint', payload);
        // Refresh daftar endpoint
        const endpointResponse = await apiClient.get(`/api/api-endpoints/${selectedVersion.value.id}`);
        endpoints.value = await endpointResponse.json();
        // Reset form
        newEndpoint.value = { http_method: 'GET', path: '', documentation: '' };
    } catch (error) {
        console.error("Error adding endpoint:", error);
        alert("Failed to add endpoint.");
    }
};

const deleteEndpoint = async (endpointId) => {
    if (confirm('Are you sure you want to delete this endpoint?')) {
        try {
            await apiClient.delete(`/api/delete-endpoint/${endpointId}`);
            // Hapus dari daftar secara lokal untuk update UI instan
            endpoints.value = endpoints.value.filter(ep => ep.id !== endpointId);
        } catch (error) {
            console.error("Error deleting endpoint:", error);
            alert("Failed to delete endpoint.");
        }
    }
};

onMounted(() => {
    fetchInitialData();
});
</script>

<style scoped>
</style>
