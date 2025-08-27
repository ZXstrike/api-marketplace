<template>
  <div class="w-full flex justify-center">
    <main class="container mx-auto px-6 py-8 md:py-12">
      <div class="lg:flex lg:space-x-8 h-screen">

        <DashboardSidebar :active-section="activeSection" @navigate="handleNavigation" />

        <main class="flex-1 container mx-auto px-6">
          <DashboardOverview v-if="activeSection === 'overview'" :subscriptions="subscriptionsCount"
            :user-balance="userBalance" :monthlyCost="monthlyCost" :requests-this-month="requestsThisMonth"
            :weeklyUsage="weeklyUsage" />
          <SubscriptionManager v-if="activeSection === 'subscriptions'" :subscriptions="subscriptions"
            :loading="loadingSubscriptions" @unsubscribe="handleUnsubscribe" @regenerate-key="handleRegenerateKey"
            @copy-key="handleCopyKey" />
          <BillingUsage v-if="activeSection === 'billing'" :user-balance="userBalance" @top-up="isModalOpen = true" />
        </main>

        <TopupModal :is-open="isModalOpen" @close="isModalOpen = false" @topup-success="handleSuccess" />
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
import TopupModal from '@/components/dashboard/consumer/TopupModal.vue';

const activeSection = ref('overview');
const subscriptions = ref([]);
const subscriptionsCount = ref(0);
const loadingSubscriptions = ref(true);
const weeklyUsage = ref([]);
const userBalance = ref(0);
const monthlyCost = ref(0);
const requestsThisMonth = ref(0);
const isModalOpen = ref(false);

const handleSuccess = () => {
  isModalOpen.value = false;
  // Refresh data after successful top-up
  fetchUsernBillingInfo();
  fetchSubscriptions();
};

// --- Logika API ---

const maskKey = (key) => {
  if (!key) {
    return 'No key generated';
  }
  return key.replace(/(.{13})(.*)(.{6})/, '$1***************************$3');
};

const fetchUsernBillingInfo = async () => {
  try {
    const response = await apiClient.get('/billing/info?wallet_type=consumer');
    const data = await response.json();
    // Proses data billing jika diperlukan
    console.log('User Billing Info:', data);
    userBalance.value = data.balance || 0;
    console.log('User Balance:', userBalance.value);
  } catch (error) {
    console.error('Failed to fetch user billing info:', error);
  }
};

const fetchUserOverview = async () => {
  try {
    const response = await apiClient.get('/api/consumer/overview');
    const data = await response.json();
    console.log('User Overview:', data);
    subscriptionsCount.value = data.active_subscriptions_count || 0;
    console.log('Subscriptions Count:', subscriptionsCount.value);
    monthlyCost.value = data.total_monthly_cost || 0;
    weeklyUsage.value = data.requests_last_7_days || [];
    requestsThisMonth.value = data.requests_last_30_days || 0;
  } catch (error) {
    console.error('Failed to fetch user overview:', error);
  }
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
  if (confirm('Are you sure you want to unsubscribe from this API?')) {
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
  if (confirm('Are you sure you want to regenerate the API key?')) {
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
  if (!navigator.clipboard) {
    const textArea = document.createElement("textarea");
    textArea.value = key;

    // Avoid scrolling to bottom
    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.position = "fixed";

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
      const successful = document.execCommand('copy');
      if (successful) {
        alert('API key copied to clipboard.');
      } else {
        alert('Failed to copy API key.');
      }
    } catch (err) {
      console.error('Fallback: Oops, unable to copy', err);
      alert('Failed to copy API key.');
    }

    document.body.removeChild(textArea);
    return;
  }

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
  fetchUsernBillingInfo();
  fetchUserOverview();
});

// Panggil ulang jika tab diubah (opsional, untuk memastikan data selalu baru)
watch(activeSection, (newSection) => {
  if (newSection === 'subscriptions') {
    fetchSubscriptions();
  }
});
</script>
