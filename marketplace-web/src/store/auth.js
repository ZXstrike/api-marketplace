import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useAuthStore = defineStore('auth', () => {
    // --- STATE ---
    // Initialize state from localStorage to persist login
    const user = ref(JSON.parse(localStorage.getItem('user')));
    const token = ref(localStorage.getItem('authToken'));

    // --- GETTERS ---
    // This is defined ONCE. It automatically updates when `user` or `token` change.
    const isLoggedIn = computed(() => !!token.value && !!user.value);
    const userRole = computed(() => user.value?.role || null);

    // --- ACTIONS ---

    /**
     * Handles the login process.
     * @param {object} userData - User data from the login form.
     * @param {string} apiToken - The token received from the backend API.
     */
    function login(userData, apiToken) {
        // Ensure the avatar URL is set, defaulting to a placeholder if not provided.
        if (!userData.avatarUrl) {
            userData.avatarUrl = `https://placehold.co/100x100/E2E8F0/4A5568?text=${userData.name ? userData.name.charAt(0) : 'U'}`;
        }

        console.log('Login successful:', userData, apiToken);
        
        // Store user data and token in localStorage
        localStorage.setItem('user', JSON.stringify(userData));
        localStorage.setItem('authToken', apiToken);

        // Update the reactive state. This is all you need to do.
        // The `isLoggedIn` computed property will update automatically.
        user.value = userData;
        token.value = apiToken;
    }

    function renewToken(newToken) {
        // Update the token in localStorage and the reactive state
        localStorage.setItem('authToken', newToken);
        token.value = newToken;
    }

    /**
     * Handles the logout process.
     */
    function logout() {
        // Clear data from localStorage
        localStorage.removeItem('user');
        localStorage.removeItem('authToken');

        // Reset the state. `isLoggedIn` will become false automatically.
        user.value = null;
        token.value = null;
    }

    return { user, token, isLoggedIn, userRole, login, logout, renewToken };
});
