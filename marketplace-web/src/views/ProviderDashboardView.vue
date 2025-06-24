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
            :totalSubs="totalSubscribers" 
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

const totalSubscribers = computed(() => {
  return myApis.value.reduce((total, api) => total + api.subscribers, 0);
});

// Fetch data when the component is mounted
onMounted(async () => {
  try {
    const response = await apiClient.get('/store/apis');
    // Map API response to the data structure expected by child 
    const data = await response.json();
    myApis.value = data.map(api => ({
      id: api.id,
      name: api.name,
      subscribers: api.subs_count,
      status: 'Published', // Placeholder as API doesn't provide this
      monthlyRevenue: 0,   // Placeholder as API doesn't provide this
    }));
  } catch (error) {
    console.error("Failed to fetch APIs:", error);
    // Optionally, set an error state to show a message to the user
  }
});

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