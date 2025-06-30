<template>
  <div class="w-full flex justify-center">
    <main class="container justify-items-center mx-auto px-6 py-8 md:py-12">
      <div class=" md:w-3xl mx-auto">
        <h1 class="text-3xl font-bold">Settings</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">Manage your account settings and personal information.</p>

        <div v-if="loading" class="text-center py-16">Loading your profile...</div>

        <div v-if="!loading && user" class="flex flex-col gap-8 mt-8">
          <!-- FIX: Profile picture and info are now combined in one card -->
          <div class="dashboard-card">
            <h2 class="section-heading">Profile Information</h2>
            <form @submit.prevent="handleProfileUpdate">
              <div class="flex flex-col md:flex-row gap-8">
                <!-- Profile Picture Section -->
                <div class="flex-shrink-0 flex flex-col items-center">
                  <img :src="avatarPreview || user.profile_picture_url || 'https://placehold.co/128x128'"
                    alt="Current Avatar" class="w-32 h-32 rounded-full object-cover shadow-md">
                  <input type="file" @change="handleFileSelect" ref="fileInput" class="hidden">
                  <button type="button" @click="$refs.fileInput.click()" class="btn-secondary mt-4">Choose
                    Image</button>
                  <button v-if="selectedFile" @click="handleAvatarUpdate" :disabled="isUploading"
                    class="btn-primary mt-2 text-sm">
                    {{ isUploading ? 'Uploading...' : 'Upload' }}
                  </button>
                </div>

                <!-- Profile Fields Section -->
                <div class="flex-grow space-y-4">
                  <div>
                    <label for="username" class="form-label">Username</label>
                    <input type="text" id="username" v-model="profileData.username" class="form-input" required>
                  </div>
                  <div>
                    <label for="email" class="form-label">Email Address</label>
                    <input type="email" id="email" v-model="profileData.email"
                      class="form-input bg-gray-100 dark:bg-gray-800" disabled>
                  </div>
                  <div>
                    <label for="description" class="form-label">Description / Bio</label>
                    <textarea id="description" rows="3" v-model="profileData.description" class="form-input"
                      placeholder="Tell us a little about yourself."></textarea>
                  </div>
                </div>
              </div>

              <div class="text-right mt-6 border-t border-gray-200 dark:border-gray-700 pt-4">
                <button type="submit" :disabled="isSavingProfile" class="btn-primary">
                  {{ isSavingProfile ? 'Saving...' : 'Save Profile Changes' }}
                </button>
              </div>
            </form>
          </div>

          <!-- NEW: Become a Provider Card -->
          <div v-if="authStore.userRole === 'consumer'" class="dashboard-card">
            <h2 class="section-heading">Seller Account</h2>
            <div class="flex flex-col md:flex-row items-center justify-between gap-4">
                <p class="text-gray-600 dark:text-gray-400">Ready to share your APIs with the world? Open a seller account to start publishing.</p>
                <button @click="handleBecomeProvider" :disabled="isUpgrading" class="btn-primary flex-shrink-0">
                    {{ isUpgrading ? 'Creating Store...' : 'Become a Provider' }}
                </button>
            </div>
          </div>

          <!-- Bagian Ubah Password -->
          <div class="dashboard-card">
            <h2 class="section-heading">Change Password</h2>
            <form @submit.prevent="handleChangePassword" class="space-y-4">
              <div>
                <label for="current-password" class="form-label">Current Password</label>
                <input type="password" id="current-password" v-model="passwordData.old_password" class="form-input"
                  required>
              </div>
              <div>
                <label for="new-password" class="form-label">New Password</label>
                <input type="password" id="new-password" v-model="passwordData.new_password" class="form-input"
                  required>
              </div>
              <div class="text-right mt-6 border-t border-gray-200 dark:border-gray-700 pt-4">
                <button type="submit" :disabled="isChangingPassword" class="btn-primary">
                  {{ isChangingPassword ? 'Saving...' : 'Change Password' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useAuthStore } from '@/store/auth';
import apiClient from '@/services/apiClient.js';

const authStore = useAuthStore();
const loading = ref(true);
const isSavingProfile = ref(false);
const isUploading = ref(false);
const isChangingPassword = ref(false);
const isUpgrading = ref(false);

const user = ref(null);
const fileInput = ref(null);
const selectedFile = ref(null);
const avatarPreview = ref(null);

const profileData = reactive({
  username: '',
  email: '',
  description: ''
});
const passwordData = reactive({
  old_password: '',
  new_password: ''
});

const fetchUserData = async () => {
  loading.value = true;
  try {
    const response = await apiClient.get('/user/me');
    user.value = await response.json();

    Object.assign(profileData, {
      username: user.value.username,
      email: user.value.email,
      description: user.value.description,
    });

  } catch (error) {
    console.error("Failed to fetch user data:", error);
    alert("Could not load your profile.");
  } finally {
    loading.value = false;
  }
};

const handleProfileUpdate = async () => {
  isSavingProfile.value = true;
  try {
    await apiClient.put('/user/update', {
      username: profileData.username,
      description: profileData.description
    });
    alert('Profile updated successfully!');
    authStore.user.name = profileData.username;
  } catch (error) {
    console.error("Error updating profile:", error);
    alert("Failed to update profile.");
  } finally {
    isSavingProfile.value = false;
  }
};

const handleFileSelect = (event) => {
  const file = event.target.files[0];
  if (file) {
    selectedFile.value = file;
    const reader = new FileReader();
    reader.onload = (e) => {
      avatarPreview.value = e.target.result;
    };
    reader.readAsDataURL(file);
  }
};

const handleAvatarUpdate = async () => {
  if (!selectedFile.value) return;
  isUploading.value = true;

  const formData = new FormData();
  formData.append('profile_picture', selectedFile.value);

  try {
    const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/user/update-profile-picture`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
      },
      body: formData,
    });

    if (!response.ok) throw new Error('Failed to upload image.');

    alert('Profile picture updated!');
    await fetchUserData();
    selectedFile.value = null;
    avatarPreview.value = null;

  } catch (error) {
    console.error("Error uploading avatar:", error);
    alert("Failed to upload image.");
  } finally {
    isUploading.value = false;
  }
};

const handleChangePassword = async () => {
  isChangingPassword.value = true;
  try {
    const response = await apiClient.put('/user/change-password', passwordData);
    if (!response.ok) throw new Error('Failed to change password.');
    const data = await response.json();
    if (!data.success) throw new Error(data.message || 'Unknown error');
    passwordData.old_password = '';
    passwordData.new_password = '';
  } catch (error) {
    console.error("Error changing password:", error);
    alert("Failed to change password. Please check your current password.");
  } finally {
    isChangingPassword.value = false;
  }
};


const handleBecomeProvider = async () => {
    if (!confirm("This will create a seller store for your account. Do you want to proceed?")) {
        return;
    }
    isUpgrading.value = true;
    try {
        // Asumsi API tidak memerlukan body, sesuaikan jika perlu
        await apiClient.post('/store/create', {});
        
        // Panggil action di store untuk update role
        authStore.upgradeToProvider();

        alert('Congratulations! Your seller store has been created.');
        
        // Arahkan ke dasbor, yang sekarang akan menampilkan dasbor provider
        router.push('/dashboard');

    } catch (error) {
        console.error("Error creating store:", error);
        alert("Failed to create your seller store. Please try again.");
    } finally {
        isUpgrading.value = false;
    }
};

onMounted(() => {
  fetchUserData();
});
</script>

<style scoped></style>
