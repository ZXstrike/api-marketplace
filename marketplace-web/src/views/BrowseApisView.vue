<template>
  <div class="flex min-h-screen w-full justify-center">
    <main class="container mx-auto px-6 py-8 md:py-12">
      <div class="lg:flex lg:gap-8">

        <FiltersSidebar 
          class="lg:w-1/4 xl:w-1/5 mb-8 lg:mb-0"
          :categories="allCategories" 
          @filter-change="applyCategoryFilter"
          @clear-filters="clearAllFilters"
        />

        <div class="flex flex-1 flex-col gap-6">

          <SearchAndSort v-model:search-query="searchQuery" />

          <div v-if="isLoading" class="text-center py-16">Loading...</div>
          <ApiGrid v-else :apis="apiList" />

          <Pagination 
            v-if="!isLoading && (apiList.length > 0 || currentPage > 1)"
            :current-page="currentPage" 
            :is-last-page="isLastPage"
            @page-changed="goToPage" 
          />

        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import apiClient from '@/services/apiClient.js'; // <-- Import apiClient
import FiltersSidebar from '@/components/api/FilterSidebar.vue';
import SearchAndSort from '@/components/api/SearchAndSort.vue';
import ApiGrid from '@/components/api/ApiGrid.vue';
import Pagination from '@/components/api/Pagination.vue';

// --- Reactive State ---
const apiList = ref([]);
const allCategories = ref([]);
const searchQuery = ref('');
const activeCategory = ref('');
const currentPage = ref(1);
const isLoading = ref(true);
const isLastPage = ref(false);
const itemsPerPage = 10;

// --- API Fetching ---
const fetchApis = async () => {
  isLoading.value = true;
  
  try {
    // FIX: Menggunakan apiClient untuk memanggil API
    // Membuat objek URLSearchParams untuk mengelola query parameters
    const params = new URLSearchParams({
        page: currentPage.value,
        length: itemsPerPage,
    });

    if (searchQuery.value) {
      params.append('search', searchQuery.value);
    }
    if (activeCategory.value) {
      params.append('category', activeCategory.value);
    }

    // Menggunakan apiClient.get dengan endpoint yang sudah termasuk query string
    const response = await apiClient.get(`/api/all?${params.toString()}`);
    
    const data = await response.json();
    apiList.value = data;
    
    // Jika API mengembalikan item kurang dari yang diminta, kita anggap itu halaman terakhir
    isLastPage.value = data.length < itemsPerPage;

  } catch (error) {
    console.error("Failed to fetch APIs:", error);
    apiList.value = []; // Kosongkan daftar jika terjadi error
  } finally {
    isLoading.value = false;
  }
};

const fetchCategories = async () => {
  try {
    // FIX: Menggunakan apiClient untuk mengambil kategori
    const response = await apiClient.get('/api/categories');
    allCategories.value = await response.json();
  } catch (error) {
    console.error("Error fetching categories:", error);
    allCategories.value = [];
  }
};

// --- Watchers untuk memicu pengambilan ulang API ---
watch([searchQuery, activeCategory], () => {
  currentPage.value = 1; // Reset ke halaman pertama saat ada pencarian/filter baru
  fetchApis();
});

watch(currentPage, fetchApis);

// --- Lifecycle Hooks ---
onMounted(() => {
  fetchApis();
  fetchCategories();
});

// --- Event Handlers ---
function applyCategoryFilter(categorySlug) {
  activeCategory.value = categorySlug;
}

function clearAllFilters() {
  activeCategory.value = '';
  searchQuery.value = ''; // Juga bersihkan pencarian
}

function goToPage(pageNumber) {
  if (pageNumber < 1) return;
  currentPage.value = pageNumber;
}
</script>
