<template>
    <div class="w-full flex justify-center">
        <main class="container justify-items-center mx-auto px-6 py-8 md:py-12">
            <div class=" md:w-3xl mx-auto">
                <router-link to="/dashboard" class="text-blue-600 hover:underline flex items-center mb-4">
                    <i data-feather="arrow-left" class="w-4 h-4 mr-2"></i>
                    Back to Dashboard
                </router-link>
                <h1 class="text-3xl font-bold">Create a New API Product</h1>
                <p class="text-gray-600 dark:text-gray-400 mt-1">Follow the steps below to list your API on the
                    marketplace.</p>
            </div>

            <form v-if="currentStep === 1" @submit.prevent="handleCreateApi" class="space-y-8 md:w-3xl mx-auto mt-8">
                <div class="dashboard-card">
                    <h2 class="section-heading">Step 1: API Details & Pricing</h2>
                    <div class="space-y-6">
                        <div>
                            <label for="api-name" class="form-label">API Name</label>
                            <input type="text" id="api-name" v-model="apiData.name"
                                placeholder="e.g., My Awesome Geolocation API" class="form-input" required>
                        </div>
                        <div>
                            <label for="api-description" class="form-label">Description</label>
                            <textarea id="api-description" rows="3" v-model="apiData.description" class="form-input"
                                placeholder="A short, clear description of what your API does." required></textarea>
                        </div>
                        <div>
                            <label for="base-url" class="form-label">Your Backend Base URL</label>
                            <input type="url" id="base-url" v-model="apiData.base_url"
                                placeholder="https://your-api.com/base/path" class="form-input" required>
                        </div>
                        <!-- FIX: Added price_per_call input field -->
                        <div>
                            <label for="price-per-call" class="form-label">Price per Call (in USD)</label>
                            <input type="number" id="price-per-call" v-model.number="apiData.price_per_call"
                                placeholder="e.g., 0.001" step="0.0001" min="0" class="form-input" required>
                            <p class="text-xs text-gray-500 mt-1">Set to 0 for a free API.</p>
                        </div>
                        <div>
                            <label class="form-label">Categories</label>
                            <div v-if="availableCategories.length > 0"
                                class="mt-2 grid grid-cols-2 md:grid-cols-3 gap-x-4 gap-y-2 border p-4 rounded-md">
                                <label v-for="category in availableCategories" :key="category.id"
                                    class="flex items-center cursor-pointer">
                                    <input type="checkbox" :value="category.slug" v-model="apiData.categories"
                                        class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500">
                                    <span class="ml-2 text-sm">{{ category.name }}</span>
                                </label>
                            </div>
                            <div v-else class="mt-2 text-sm text-gray-500">Loading categories...</div>
                        </div>
                    </div>
                </div>
                <div class="flex justify-end">
                    <button type="submit" :disabled="isLoading" class="btn-primary">
                        {{ isLoading ? 'Saving...' : 'Save & Add Endpoints' }}
                    </button>
                </div>
            </form>

            <!-- Step 2 Form -->
            <form v-if="currentStep === 2" @submit.prevent="handleCreateEndpoints"
                class="space-y-8 md:w-3xl mx-auto mt-8">
                <div class="dashboard-card">
                    <h2 class="section-heading">Step 2: Define Endpoints for {{ apiData.name }}</h2>
                    <div id="endpoints-list" class="space-y-4 mb-6">
                        <div v-for="(endpoint, index) in endpoints" :key="index"
                            class="p-3 bg-gray-50 dark:bg-gray-800/50 rounded-lg border dark:border-gray-700 flex justify-between items-center">
                            <div>
                                <span :class="getMethodColor(endpoint.http_method)"
                                    class="font-mono text-sm font-bold px-2 py-1 rounded-md">{{ endpoint.http_method
                                    }}</span>
                                <span class="ml-4 font-mono">{{ endpoint.path }}</span>
                            </div>
                            <button type="button" @click="removeEndpoint(index)"
                                class="text-red-500 p-1 rounded-full hover:bg-red-100 dark:hover:bg-red-900/50">
                                <i data-feather="x" class="w-5 h-5"></i>
                            </button>
                        </div>
                        <div v-if="endpoints.length === 0"
                            class="text-center py-8 text-gray-500 border-2 border-dashed rounded-lg">
                            <p>No endpoints added yet.</p>
                        </div>
                    </div>
                    <div class="p-4 border border-dashed rounded-lg">
                        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                            <div><label class="form-label">Method</label><select v-model="newEndpoint.http_method"
                                    class="form-input">
                                    <option>GET</option>
                                    <option>POST</option>
                                    <option>PUT</option>
                                    <option>DELETE</option>
                                </select></div>
                            <div><label class="form-label">Path</label><input type="text" v-model="newEndpoint.path"
                                    placeholder="/users/{id}" class="form-input"></div>
                        </div>
                        <div class="mt-4"><label class="form-label">Documentation (JSON Output)</label><textarea
                                v-model="newEndpoint.documentation" rows="4" class="form-input font-mono"
                                placeholder='Example: {"message": "Success"}'></textarea></div>
                        <div class="text-right mt-4">
                            <button type="button" @click="addEndpoint" class="btn-secondary">Add Endpoint</button>
                        </div>
                    </div>
                </div>
                <div class="flex justify-between">
                    <button type="button" @click="currentStep = 1" class="btn-secondary">Back to Details</button>
                    <button type="submit" :disabled="isLoading || endpoints.length === 0" class="btn-primary">
                        {{ isLoading ? 'Publishing...' : 'Publish API' }}
                    </button>
                </div>
            </form>
        </main>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import apiClient from '@/services/apiClient.js';

const router = useRouter();
const isLoading = ref(false);
const currentStep = ref(1);

const availableCategories = ref([]);
const apiData = reactive({
    name: '',
    description: '',
    base_url: '',
    price_per_call: 0, // FIX: Added price_per_call to the reactive state
    categories: [],
});

const newlyCreatedApiVersionId = ref(null);
const endpoints = ref([]);
const newEndpoint = ref({
    http_method: 'GET',
    path: '',
    documentation: '',
});

// --- Methods ---

const fetchCategories = async () => {
    try {
        const response = await apiClient.get('/api/categories');
        const data = await response.json();
        if (Array.isArray(data)) {
            availableCategories.value = data;
        }
    } catch (error) {
        console.error("Error fetching categories:", error);
        alert("Could not load categories. Please try again later.");
    }
};

const handleCreateApi = async () => {
    isLoading.value = true;
    try {
        const response = await apiClient.post('/api/create', apiData);
        const data = await response.json();

        if (!response.ok) throw new Error(data.message || 'Failed to create API.');

        console.log("API created successfully:", data);

        if (!data || !data.api_id) {
            throw new Error("API created, but did not return the required API Version ID.");
        }

        newlyCreatedApiVersionId.value = data.api_id;
        alert(data.message || "API details saved successfully!");
        currentStep.value = 2;

    } catch (error) {
        console.error("Error creating API:", error);
        alert(`Error: ${error.message}`);
    } finally {
        isLoading.value = false;
    }
};

const addEndpoint = () => {
    if (!newEndpoint.value.path) {
        alert("Endpoint path cannot be empty.");
        return;
    }
    endpoints.value.push({ ...newEndpoint.value });
    newEndpoint.value = { http_method: 'GET', path: '', documentation: '' };
};

const removeEndpoint = (index) => {
    endpoints.value.splice(index, 1);
};

const handleCreateEndpoints = async () => {
    isLoading.value = true;
    try {
        const payload = {
            api_version_id: newlyCreatedApiVersionId.value,
            endpoints: endpoints.value,
        };
        const response = await apiClient.post('/api/create-endpoint', payload);
        const data = await response.json();

        if (!response.ok) throw new Error(data.message || 'Failed to create endpoints.');

        alert(data.message || "API and endpoints published successfully!");
        router.push('/dashboard');

    } catch (error) {
        console.error("Error creating endpoints:", error);
        alert(`Error: ${error.message}`);
    } finally {
        isLoading.value = false;
    }
};

const getMethodColor = (method) => {
    const colors = { 'GET': 'bg-blue-100 text-blue-800', 'POST': 'bg-green-100 text-green-800', 'PUT': 'bg-yellow-100 text-yellow-800', 'DELETE': 'bg-red-100 text-red-800' };
    return colors[method.toUpperCase()] || 'bg-gray-100 text-gray-800';
};

onMounted(() => {
    fetchCategories();
});
</script>

<style scoped></style>
