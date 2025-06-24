<template>
    <div class="min-h-screen flex flex-col items-center justify-center auth-bg p-4">
        <router-link to="/" class="flex items-center space-x-2 mb-8">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none"
                     stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                     class="text-blue-600">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"></path>
            </svg>
            <span class="text-4xl font-extrabold text-center text-gray-800">Go-API Mart</span>
        </router-link>

        <div
                class="w-full max-w-md bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl shadow-2xl p-8 overflow-hidden relative">
            <div class="form-container">
                <h2 class="text-2xl font-bold text-center text-gray-900 dark:text-white">Create Account</h2>
                <p class="text-center text-gray-600 dark:text-gray-400 mt-2">Get started with the best APIs on the market.</p>

                <form @submit.prevent="handleRegister" class="mt-8 flex flex-col gap-3">
                    <div v-if="errorMessage"
                             class="p-3 bg-red-100 dark:bg-red-900/30 border border-red-400 dark:border-red-500/50 text-red-700 dark:text-red-300 rounded-md text-sm">
                        {{ errorMessage }}
                    </div>
                    <div>
                        <label for="register-username" class="form-label">Username</label>
                        <input id="register-username" v-model="registerData.username" type="text" required
                                     class="w-full form-input mt-1" placeholder="your-username">
                    </div>
                    <div>
                        <label for="register-email" class="form-label">Email Address</label>
                        <input id="register-email" v-model="registerData.email" type="email" autocomplete="email"
                                     required class="w-full form-input mt-1" placeholder="you@example.com">
                    </div>
                    <div>
                        <label for="register-password" class="form-label">Password</label>
                        <div class="relative mt-1">
                            <input id="register-password" v-model="registerData.password" type="password"
                                         autocomplete="new-password" required class="w-full form-input mt-0 pr-10"
                                         placeholder="Create a strong password">
                            <button type="button"
                                            class="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 dark:text-gray-400 hover:text-blue-600"
                                            @click="togglePasswordVisibility('register-password')">
                                <i data-feather="eye-off" class="w-5 h-5"></i>
                            </button>
                        </div>
                    </div>
                    <div>
                        <label for="confirm-password" class="form-label">Confirm Password</label>
                        <div class="relative mt-1">
                            <input id="confirm-password" v-model="registerData.confirmPassword" type="password"
                                         autocomplete="new-password" required class="w-full form-input mt-0 pr-10"
                                         placeholder="Re-type your password">
                            <button type="button"
                                            class="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 dark:text-gray-400 hover:text-blue-600"
                                            @click="togglePasswordVisibility('confirm-password')">
                                <i data-feather="eye-off" class="w-5 h-5"></i>
                            </button>
                        </div>
                    </div>
                    <div>
                        <button type="submit"
                                        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-lg font-bold text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                            Create Account
                        </button>
                    </div>
                </form>
                <p class="pt-4 text-center text-sm text-gray-600 dark:text-gray-400">
                    Already have an account?
                    <router-link to="/login" class="font-medium text-blue-600 hover:text-blue-500">
                        Log in here.
                    </router-link>
                </p>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { RouterLink, useRouter } from 'vue-router';

onMounted(() => {
    if (window.feather) {
        window.feather.replace();
    }
});

const router = useRouter();

const registerData = ref({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
});

const errorMessage = ref(null);
const handleRegister = async () => {
    errorMessage.value = null;
    if (registerData.value.password !== registerData.value.confirmPassword) {
        errorMessage.value = "Passwords do not match.";
        return;
    }

    try {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: registerData.value.username,
                email: registerData.value.email,
                password: registerData.value.password,
            }),
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to register user.');
        }

        // On success, redirect to the login page
        await router.push('/login');
    } catch (error) {
        errorMessage.value = error.message;
        console.error('Registration failed:', error);
    }
};

const togglePasswordVisibility = (fieldId) => {
        const input = document.getElementById(fieldId);
        if (!input) return;

        const button = input.nextElementSibling;

        if (input.type === "password") {
                input.type = "text";
                if (button) button.innerHTML = '<i data-feather="eye" class="w-5 h-5"></i>';
        } else {
                input.type = "password";
                if (button) button.innerHTML = '<i data-feather="eye-off" class="w-5 h-5"></i>';
        }
        // Re-render the icon
        if (window.feather) {
                window.feather.replace();
        }
};
</script>

<style scoped>
</style>