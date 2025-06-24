<template>
  <div class="min-h-screen flex flex-col items-center justify-center auth-bg p-4">
    <router-link to="/" class="flex items-center space-x-2 mb-8">
      <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none"
        stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-blue-600">
        <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"></path>
      </svg>
      <span class="text-4xl font-extrabold text-center text-gray-800 ">Go-API Mart</span>
    </router-link>

    <div
      class="w-full max-w-md bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl shadow-2xl p-8 overflow-hidden relative">
      <div class="form-container">
        <h2 class="text-2xl font-bold text-center text-gray-900 dark:text-white">Welcome Back</h2>
        <p class="text-center text-gray-600 dark:text-gray-400 mt-2">Sign in to continue to your dashboard.</p>

        <form @submit.prevent="handleLogin" class="mt-8 flex flex-col gap-3">
          <div v-if="errorMessage"
            class="p-3 bg-red-100 dark:bg-red-900/30 border border-red-400 dark:border-red-500/50 text-red-700 dark:text-red-300 rounded-md text-sm">
            {{ errorMessage }}
          </div>
          <div>
            <label for="login-email" class="form-label">Email Address</label>
            <div>
              <input id="login-email" v-model="loginData.email" type="email" autocomplete="email" required
                class="w-full form-input mt-1" placeholder="you@example.com">
            </div>
          </div>

          <div>
            <div class="flex justify-between items-center">
              <label for="login-password" class="form-label">Password</label>
              <a href="#" class="text-sm font-medium text-blue-600 hover:text-blue-500">Forgot password?</a>
            </div>
            <div class="relative mt-1">
              <input id="login-password" v-model="loginData.password" type="password" autocomplete="current-password"
                required class="form-input w-full mt-0 pr-10" placeholder="••••••••">
              <button type="button"
                class="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 dark:text-gray-400 hover:text-blue-600"
                @click="togglePasswordVisibility('login-password')">
                <i data-feather="eye-off" class="w-5 h-5"></i>
              </button>
            </div>
          </div>

          <div>
            <button type="submit"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-lg font-bold text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
              Log In
            </button>
          </div>
        </form>
        <p class="pt-4 text-center text-sm text-gray-600 dark:text-gray-400">
          Don't have an account?
          <router-link to="/register" class="font-medium text-blue-600 hover:text-blue-500">
            Sign up here.
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { RouterLink, useRouter } from 'vue-router';
import { useAuthStore } from '@/store/auth';

onMounted(() => {
  if (window.feather) {
    window.feather.replace();
  }
});

const router = useRouter();
const authStore = useAuthStore();

// Form fields
const loginData = ref({
  email: '',
  password: ''
});

const errorMessage = ref(null);
async function handleLogin() {
  errorMessage.value = null;
  try {
    // 1. Call the login API
    const loginResponse = await fetch(`${import.meta.env.VITE_API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(loginData.value),
    });

    if (!loginResponse.ok) {
      const errorData = await loginResponse.json();
      throw new Error(errorData.message || 'Invalid email or password.');
    }

    const { token } = await loginResponse.json();

    // 2. Get user profile with the token
    const profileResponse = await fetch(`${import.meta.env.VITE_API_BASE_URL}/user/me`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (!profileResponse.ok) {
      throw new Error('Failed to fetch user profile.');
    }

    const userData = await profileResponse.json();
    console.log('User data:', userData);

    // Construct avatarUrl and determine role from userData
    const avatarUrl = userData.profile_picture_url
      ? `${import.meta.env.VITE_API_BASE_URL}${userData.profile_picture_url}`
      : null;

    // Map 'store_owner' role to 'provider', default to 'consumer'
    const isProvider = userData.roles?.some(role => role.name === 'store_owner');
    const role = isProvider ? 'provider' : 'consumer';

    console.log('Role determined:', role);

    const user = {
      name: userData.username,
      avatarUrl: avatarUrl,
      role: role
    };

    console.log('User profile:', user);

    // 3. Call the store action to update the state
    authStore.login(user, token);

    // 4. Navigate to the dashboard
    router.push('/dashboard');

  } catch (e) {
    errorMessage.value = e.message || 'Login failed. Please check your credentials.';
    console.error(e);
  }
}

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

<style scoped></style>