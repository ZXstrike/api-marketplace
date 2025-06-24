import { useAuthStore } from "@/store/auth";

const BASE_URL = import.meta.env.VITE_API_BASE_URL; // Ganti dengan URL backend Anda

/**
 * Fungsi pembungkus kustom untuk fetch yang menangani token dan refresh otomatis.
 * @param {string} endpoint - Endpoint yang akan dipanggil (misal: '/api/all').
 * @param {object} options - Opsi konfigurasi untuk fetch (method, body, dll.).
 * @returns {Promise<Response>} - Mengembalikan respons dari fetch.
 */
async function apiFetch(endpoint, options = {}) {
  const authStore = useAuthStore();

  // If a token exists, refresh it before making the API call.
  if (authStore.token) {
    try {
      console.log("Attempting to refresh token before request...");
      await refreshToken();
      console.log("Token refreshed successfully.");
    } catch (refreshError) {
      console.error("Failed to refresh token, logging out:", refreshError);
      // If refresh fails, logout the user and stop the request.
      authStore.logout();
      window.location.hash = "/login";
      return Promise.reject(refreshError);
    }
  }

  // Siapkan header awal
  const headers = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  // Tambahkan token (yang mungkin baru saja di-refresh)
  if (authStore.token) {
    headers["Authorization"] = `Bearer ${authStore.token}`;
  }

  // Gabungkan konfigurasi
  const config = {
    ...options,
    headers,
  };

  // Lakukan panggilan fetch
  const response = await fetch(`${BASE_URL}${endpoint}`, config);

  // If the request still fails with 401, it means the refresh might have failed
  // or the user is truly unauthorized. In this case, we log them out.
  if (response.status === 401) {
    console.error("Request failed with 401 after token refresh. Logging out.");
    authStore.logout();
    window.location.hash = "/login";
    // We throw an error to let the calling code know the request failed.
    return Promise.reject(new Error("Unauthorized after token refresh"));
  }

  console.log("Response received:", response);

  return response;
}

/**
 * Fungsi terpisah untuk menangani logika refresh token.
 */
async function refreshToken() {
  const authStore = useAuthStore();
  // Di sini Anda perlu refresh token yang disimpan (jika ada)
  // Untuk contoh ini, kita anggap refresh token bisa didapat dari suatu tempat

  // Panggil endpoint /auth/refresh
  console.log("Refreshing token...");
  console.log("Current token:", authStore.token);

  const response = await fetch(`${BASE_URL}/auth/refresh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authStore.token}`, // Gunakan token saat ini untuk otentikasi
    },
    body: JSON.stringify({
      // Tambahkan refresh token jika diperlukan
      token: authStore.token, // Misalnya, jika Anda menyimpan refresh token di store
    }),
    // Anda mungkin perlu mengirim refresh token di sini
  });

  if (!response.ok) {
    throw new Error("Could not refresh token.");
  }

  const data = await response.json();
  const newAccessToken = data.token; // Sesuaikan dengan struktur respons Anda

  console.log("New access token received:", newAccessToken);
  // Update token di store
  authStore.renewToken(newAccessToken);  

  return newAccessToken;
}

// Ekspor metode-metode yang mudah digunakan (get, post, dll.)
export default {
refreshToken,
  get: (endpoint, options = {}) =>
    apiFetch(endpoint, { ...options, method: "GET" }),
  post: (endpoint, body, options = {}) =>
    apiFetch(endpoint, {
      ...options,
      method: "POST",
      body: JSON.stringify(body),
    }),
  put: (endpoint, body, options = {}) =>
    apiFetch(endpoint, {
      ...options,
      method: "PUT",
      body: JSON.stringify(body),
    }),
  delete: (endpoint, options = {}) =>
    apiFetch(endpoint, { ...options, method: "DELETE" }),
};
