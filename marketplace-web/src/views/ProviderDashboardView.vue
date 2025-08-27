<template>
  <div class="w-full flex justify-center">
    <main class="container mx-auto px-6 py-8 md:py-12">
      <div class="lg:flex lg:space-x-8">

        <DashboardSidebar 
          :active-section="activeSection" 
          @navigate="handleNavigation"
        />

        <main class="flex-1 container mx-auto px-6">
          <DashboardOverview 
            v-if="activeSection === 'overview'" 
            :apis="myApis"
            :totalSubs="totalSubs" 
            :user-balance="userBalance"
            :revenue="revenue"
            :monthlyRequest="monthlyRequest"
            :publishedApiCount="publishedApiCount"
            :total-subs="totalSubs"
            :weekly-revenue="weeklyRevenue"
          />
          <MyApis 
            v-if="activeSection === 'apis'" 
            :apis="myApis"
            @delete-api="handleDeleteApi" 
            @edit-api="handleEditApi"
          />
          <DashboardAnalytics v-if="activeSection === 'analytics'" />
        </main>

      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import apiClient from '@/services/apiClient'; // Assuming a pre-configured axios instance
import DashboardSidebar from '@/components/dashboard/provider/Sidebar.vue';
import DashboardOverview from '@/components/dashboard/provider/Overview.vue';
import MyApis from '@/components/dashboard/provider/MyApis.vue';
import DashboardAnalytics from '@/components/dashboard/provider/Analytics.vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const activeSection = ref('overview');
const myApis = ref([]);
const userBalance = ref(0);

const revenue = ref(0); // Total revenue in the last 30 days
const monthlyRequest = ref(0); // Total requests in the last 30 days
const publishedApiCount = ref(0);
const totalSubs = ref(0);

const weeklyRevenue = ref([0]);

// Fetch data when the component is mounted
onMounted(async () => {
  try {
    const response = await apiClient.get('/store/apis');
    // Map API response to the data structure expected by child 
    const data = await response.json();

    console.log('API Response Data:', data);
    myApis.value = data.map(api => ({
      id: api.id,
      name: api.name,
      subscribers: api.subs_count,
      status: 'Published', // Placeholder as API doesn't provide this
      monthlyRevenue: api.total_revenue,   // Placeholder as API doesn't provide this
    }));

    console.log('Fetched APIs:', myApis.value);

    fetchUsernBillingInfo();
    fetchProviderOverview();
  } catch (error) {
    console.error("Failed to fetch APIs:", error);
    // Optionally, set an error state to show a message to the user
  }
});


const fetchProviderOverview = async () => {
  try {
    const response = await apiClient.get('/api/provider/overview');
    const data = await response.json();
    // Process overview data if needed
    console.log('Provider Overview:', data);
    revenue.value = data.total_revenue || 0;
    totalSubs.value = data.active_subscriber_count || 0;
    publishedApiCount.value = data.published_apis_count || 0;
    monthlyRequest.value = data.requests_last_30_days || 0;
    weeklyRevenue.value = data.requests_last_4_weeks || [0];
    console.log('Weekly Revenue:', weeklyRevenue.value);
  } catch (error) {
    console.error('Failed to fetch provider overview:', error);
  }
};


const fetchUsernBillingInfo = async () => {
  try {
    const response = await apiClient.get('/billing/info?wallet_type=provider');
    const data = await response.json();
    // Proses data billing jika diperlukan
    console.log('User Billing Info:', data);
    userBalance.value = data.balance || 0;
    console.log('User Balance:', userBalance.value);
  } catch (error) {
    console.error('Failed to fetch user billing info:', error);
  }
};



// Function to handle the 'navigate' event from the sidebar
const handleNavigation = (section) => {
  activeSection.value = section;
};

// Function to handle the 'delete-api' event from the MyApis component
const handleDeleteApi = (id) => {
  console.log(`Deleting API with ID: ${id} from parent component.`);
  // After successful API call, you would remove the item from the 'myApis' ref array
  const index = myApis.value.findIndex(api => api.id === id);
  if (index !== -1) {
    myApis.value.splice(index, 1);
  }
  // TODO: Call API endpoint DELETE /api/delete/:id
};

const handleEditApi = (api) => {
  console.log(`Editing API: ${api.name}`);
  // Navigate to the edit page or open a modal
  router.push(`/dashboard/edit-api/${api}`);
};

</script>

<style scoped>
/* Scoped styles specific to the main dashboard layout can go here */
</style>