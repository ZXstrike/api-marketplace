import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useAuthStore = defineStore('auth', () => {
    // --- STATE ---
    // Initialize state from localStorage to persist the login session.
    const user = ref(JSON.parse(localStorage.getItem('user')));
    const token = ref(localStorage.getItem('authToken'));

    // --- GETTERS ---
    // Computed properties automatically update when their dependencies change.
    const isLoggedIn = computed(() => !!token.value && !!user.value);
    const userRole = computed(() => user.value?.role || null);

    // --- ACTIONS ---

    /**
     * Handles the login process by saving user data and the token.
     * @param {object} userData - User data from the login form.
     * @param {string} apiToken - The token received from the backend API.
     */
    function login(userData, apiToken) {
        // Set a default avatar if one isn't provided.
        if (!userData.avatarUrl) {
            const initial = userData.name ? userData.name.charAt(0).toUpperCase() : 'U';
            userData.avatarUrl = `https://placehold.co/100x100/E2E8F0/4A5568?text=${initial}`;
        }

        console.log('Login successful:', userData, apiToken);

        // Update reactive state.
        user.value = userData;
        token.value = apiToken;

        // Persist data to localStorage.
        localStorage.setItem('user', JSON.stringify(userData));
        localStorage.setItem('authToken', apiToken);
    }

    /**
     * Renews the authentication token in the state and localStorage.
     * @param {string} newToken - The new authentication token.
     */
    function renewToken(newToken) {
        token.value = newToken;
        localStorage.setItem('authToken', newToken);
    }

    /**
     * Upgrades the current user's role to 'provider' in the state and localStorage.
     */
    function upgradeToProvider() {
        if (user.value) {
            user.value.role = 'provider';
            localStorage.setItem('user', JSON.stringify(user.value));
            console.log('User role upgraded to provider in store.');
        }
    }

    /**
     * Handles the logout process by clearing user data and the token.
     */
    function logout() {
        // Reset reactive state.
        user.value = null;
        token.value = null;

        // Clear data from localStorage.
        localStorage.removeItem('user');
        localStorage.removeItem('authToken');
    }

    return {
        // State
        user,
        token,
        // Getters
        isLoggedIn,
        userRole,
        // Actions
        login,
        logout,
        renewToken,
        upgradeToProvider,
    };
});
