<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-blue-100">
    <div class="bg-white shadow-lg rounded-lg p-8 w-full max-w-md">
      <h2 class="text-2xl font-bold mb-6 text-center text-blue-800">Reset Password (Manager)</h2>
      <form @submit.prevent="onSubmit" class="space-y-4">
        <div>
          <label class="block text-gray-700 mb-1">Username or Email</label>
          <input v-model="form.username_or_email" type="text" required class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400" />
        </div>
        <div>
          <label class="block text-gray-700 mb-1">New Password</label>
          <input v-model="form.new_password" type="password" required class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400" />
        </div>
        <button type="submit" class="w-full bg-blue-700 text-white py-2 rounded hover:bg-blue-800 transition">Reset Password</button>
      </form>
      <div v-if="success" class="text-green-600 mt-4 text-center">Password reset successful! Please login again.</div>
      <div v-if="error" class="text-red-600 mt-4 text-center">{{ error }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { resetUserPassword } from '../api/auth';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const auth = useAuthStore();
const form = ref({ username_or_email: '', new_password: '' });
const error = ref('');
const success = ref(false);

async function onSubmit() {
  error.value = '';
  success.value = false;
  try {
    await resetUserPassword(form.value);
    success.value = true;
    // Clear authentication to force fresh login with updated firstLogin status
    auth.clearAuth();
    setTimeout(() => router.push('/login/user'), 1500);
  } catch (err) {
    error.value = err.response?.data?.error || 'Reset failed';
  }
}
</script>
