<template>
  <div class="w-full flex justify-center">
    <main class="container mx-auto px-6 py-8 md:py-12">
      <div class="lg:flex lg:space-x-8 h-screen">

        <DashboardSidebar 
          :active-section="activeSection" 
          @navigate="handleNavigation"
        />

        <main class="flex-1 container mx-auto px-6">
          <DashboardOverview 
            v-if="activeSection === 'overview'" 
            :subscriptions="subscriptions" 
          />
          <SubscriptionManager 
            v-if="activeSection === 'subscriptions'" 
            :subscriptions="subscriptions"
            :loading="loadingSubscriptions"
            @unsubscribe="handleUnsubscribe" 
            @regenerate-key="handleRegenerateKey"
            @copy-key="handleCopyKey"
          />
          <BillingUsage v-if="activeSection === 'billing'" />
        </main>

      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import apiClient from '@/services/apiClient.js';

// Import komponen anak
import DashboardSidebar from '@/components/dashboard/consumer/Sidebar.vue';
import DashboardOverview from '@/components/dashboard/consumer/Overview.vue';
import SubscriptionManager from '@/components/dashboard/consumer/Subscription.vue';
import BillingUsage from '@/components/dashboard/consumer/Billing.vue';

const activeSection = ref('overview');
const subscriptions = ref([]);
const loadingSubscriptions = ref(true);

// --- Logika API ---

const maskKey = (key) => {
  if (!key) {
    return 'No key generated';
  }
  return key.replace(/(.{13})(.*)(.{6})/, '$1***************************$3');
};

const fetchSubscriptions = async () => {
  loadingSubscriptions.value = true;
  try {
    const response = await apiClient.get('/subscriptions/get-by-user');
    const data = await response.json();

    // Lakukan pengecekan untuk memastikan data adalah sebuah array sebelum digunakan.
    if (Array.isArray(data)) {
      subscriptions.value = data.map(sub => {
        // Memetakan data dari respons API ke format yang dibutuhkan komponen
        return {
          id: sub.id,
          apiName: sub.api_version?.api?.name || 'Unknown API',
          apiId: sub.api_version?.api?.id || 'Unknown ID',
          // 'created_at' is not available in the new JSON structure.
          // Providing a placeholder.
          subscribedDate: 'N/A',
          apiKey: sub.api_keys?.[0]?.key_prefix || 'No key generated',
          apiKeyMasked: maskKey(sub.api_keys?.[0]?.key_prefix) || 'No key generated'
        };
      });
    } else {
      // Jika respons bukan array (misalnya null atau objek error), set ke array kosong.
      subscriptions.value = [];
    }

  } catch (error) {
    console.error("Failed to fetch subscriptions:", error);
    subscriptions.value = []; // Pastikan selalu array kosong jika terjadi error
  } finally {
    loadingSubscriptions.value = false;
  }
};

const handleUnsubscribe = async (subscriptionId) => {
    if(confirm('Are you sure you want to unsubscribe from this API?')) {
        try {
            await apiClient.post('/subscriptions/unsubscribe', { subscription_id: subscriptionId });
            alert('Successfully unsubscribed.');
            // Ambil ulang daftar langganan untuk memperbarui UI
            await fetchSubscriptions();
        } catch (error) {
            console.error('Failed to unsubscribe:', error);
            alert('Unsubscription failed.');
        }
    }
};

const handleRegenerateKey = async (subscriptionId) => {
    if(confirm('Are you sure you want to regenerate the API key?')) {
        try {
            const response = await apiClient.post('/api-keys/create', { subscription_id: subscriptionId });
            alert('API key regenerated successfully.');
            await fetchSubscriptions();
        } catch (error) {
            console.error('Failed to regenerate API key:', error);
            alert('Failed to regenerate API key.');
        }
    }
};

const handleCopyKey = (key) => {
    navigator.clipboard.writeText(key).then(() => {
        alert('API key copied to clipboard.');
    }).catch(err => {
        console.error('Failed to copy API key:', err);
        alert('Failed to copy API key.');
    });
};

// --- Logika Navigasi ---
const handleNavigation = (section) => {
  activeSection.value = section;
};

// --- Lifecycle & Watchers ---
onMounted(() => {
    fetchSubscriptions();
});

// Panggil ulang jika tab diubah (opsional, untuk memastikan data selalu baru)
watch(activeSection, (newSection) => {
    if (newSection === 'subscriptions') {
        fetchSubscriptions();
    }
});
</script>
