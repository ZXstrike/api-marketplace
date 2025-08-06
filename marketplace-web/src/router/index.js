import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/store/auth";
import apiClient from "@/services/apiClient"; // Pastikan ini sesuai dengan path yang benar

const routes = [
  {
    path: "/",
    name: "home",
    component: () => import("@/views/HomeView.vue"),
  },
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/LoginView.vue"),
    meta: {
      hideNavbar: true,
      hideFooter: true,
    },
  },
  {
    path: "/register",
    name: "register",
    component: () => import("@/views/RegisterView.vue"),
    meta: {
      hideNavbar: true,
      hideFooter: true,
    },
  },
  {
    path: "/settings",
    name: "settings",
    component: () => import("@/views/SettingsView.vue"),
    meta: { requiresAuth: true }, // meta field untuk proteksi route
  },
  {
    path: "/browse",
    name: "browse",
    component: () => import("@/views/BrowseApisView.vue"),
  },
  {
    path: "/docs",
    name: "docs",
    component: () => import("@/views/DocsView.vue"),
  },
  {
    path: "/api-details/:id", // Example with a dynamic parameter
    name: "api-details",
    component: () => import("@/views/ApiDetailsView.vue"),
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: () => import("@/views/ConsumerDashboardView.vue"),
    meta: { requiresAuth: true }, // meta field untuk proteksi route
  },
  {
    path: "/store",
    name: "store",
    component: () => import("@/views/ProviderDashboardView.vue"),
    meta: { requiresAuth: true }, // meta field untuk proteksi route
  },
  {
    path: "/admin",
    name: "admin",
    component: () => import("@/views/AdminDashboardView.vue"),
    // meta: { requiresAuth: true }, // meta field untuk proteksi route
    
  },
  {
    path: "/dashboard/create-new-api",
    name: "create-api",
    component: () => import("@/views/CreateNewApiView.vue"),
    meta: { requiresAuth: true }, // meta field untuk proteksi route
  },
  {
    path: "/dashboard/edit-api/:id",
    name: "edit-api",
    component: () => import("@/views/EditAPIView.vue"),
    meta: { requiresAuth: true }, // meta field untuk proteksi route
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();

  // Cek apakah rute tujuan membutuhkan autentikasi
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    // Jika butuh auth, tapi pengguna belum login
    // alihkan ke halaman login
    next({ name: "login" });
  }
  if (to.meta.requiresAuth && authStore.isLoggedIn) {
    try {
      // Coba refresh token sebelum melanjutkan
      await apiClient.refreshToken();
      next(); // Lanjutkan ke rute yang diminta
    } catch (error) {
      console.error("Unexpected error during token refresh:", error);
      authStore.logout(); // Jika ada error lain, logout pengguna
      if (to.meta.requiresAuth) {
        next({ name: "login" }); // Arahkan ke halaman login
      }
    }
  } else {
    try {
      // Coba refresh token sebelum melanjutkan
      await apiClient.refreshToken();
      next(); // Lanjutkan ke rute yang diminta
    } catch (error) {
      console.error("Unexpected error during token refresh:", error);
      authStore.logout(); // Jika ada error lain, logout pengguna
    }
    next();
  }
});

export default router;
