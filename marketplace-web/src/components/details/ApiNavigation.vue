<template>
  <div class="flex justify-between items-center border-b border-gray-200 dark:border-gray-700">
    <nav class="-mb-px flex space-x-6" aria-label="Tabs">
      <a href="#" @click.prevent="updateTab('endpoints')" :class="{'active': activeTab === 'endpoints'}" class="tab-link">Endpoints</a>
      <a href="#" @click.prevent="updateTab('pricing')" :class="{'active': activeTab === 'pricing'}" class="tab-link">Pricing</a>
    </nav>
    
    <div v-if="versions && versions.length > 0">
      <select :value="selectedVersionId" @change="updateVersion" class="form-input text-sm">
        <option v-for="version in versions" :key="version.id" :value="version.id">
          {{ version.version_string }}
        </option>
      </select>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  versions: { type: Array, default: () => [] },
  selectedVersionId: { type: String, default: null },
  activeTab: { type: String, required: true }
});

const emit = defineEmits(['update:selectedVersionId', 'update:activeTab']);

const updateVersion = (event) => {
  emit('update:selectedVersionId', event.target.value);
};

const updateTab = (tabName) => {
  emit('update:activeTab', tabName);
};
</script>

<style scoped>
</style>