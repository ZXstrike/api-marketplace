<template>
    <aside class="lg:w-1/4 xl:w-1/5">
        <nav class="space-y-2">
            <a href="#" @click.prevent="$emit('navigate', 'revenue')"
                :class="{ 'active': activeSection === 'revenue' }" class="sidebar-link">
                <i data-feather="bar-chart-2" class="w-5 h-5 mr-3"></i>
                <span>Revenue</span>
            </a>
            <a href="#" @click.prevent="$emit('navigate', 'users')" :class="{ 'active': activeSection === 'users' }"
                class="sidebar-link">
                <i data-feather="users" class="w-5 h-5 mr-3"></i>
                <span>Users</span>
            </a>
            <a href="#" @click.prevent="$emit('navigate', 'apis')" :class="{ 'active': activeSection === 'apis' }"
                class="sidebar-link">
                <i data-feather="package" class="w-5 h-5 mr-3"></i>
                <span>APIs</span>
            </a>
            <a href="#" @click.prevent="$emit('navigate', 'settings')" :class="{ 'active': activeSection === 'settings' }"
                class="sidebar-link">
                <i data-feather="settings" class="w-5 h-5 mr-3"></i>
                <span>Settings</span>
            </a>
        </nav>
    </aside>
</template>

<script setup>
defineProps({
    activeSection: {
        type: String,
        required: true
    }
});
defineEmits(['navigate']);
</script>