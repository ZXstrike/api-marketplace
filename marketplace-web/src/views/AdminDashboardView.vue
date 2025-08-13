<template>
  <div class="w-full flex justify-center">
    <main class="container mx-auto px-6 py-8 md:py-12">
      <div class="lg:flex lg:!space-x-8">
        <!-- sidebar -->
        <AdminSidebar :active-section="activeSection" @navigate="handleNavigation" />

        <!-- dynamic content -->
        <main class="flex-1 mt-8 lg:mt-0">
          <component :is="currentPage" />
        </main>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue';
import AdminSidebar  from '@/components/dashboard/admin/Sidebar.vue';
import RevenueView   from '@/components/dashboard/admin/RevenueView.vue';
import UsersView     from '@/components/dashboard/admin/UserView.vue';
import ApisView      from '@/components/dashboard/admin/ApisView.vue';

const activeSection = ref('revenue');1

const pages = {
  revenue: RevenueView,
  users:   UsersView,
  apis:    ApisView,
};

const currentPage = computed(() => pages[activeSection.value] || RevenueView);

const handleNavigation = (section) => (activeSection.value = section);
</script>