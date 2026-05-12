<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-blue-100">
    <div class="bg-white shadow-lg rounded-lg p-8 w-full max-w-md">
      <h2 class="text-2xl font-bold mb-6 text-center text-blue-800">User Login</h2>
      <form @submit.prevent="onSubmit" class="space-y-4">
        <div>
          <label class="block text-gray-700 mb-1">Username</label>
          <input v-model="form.username" type="text" required class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400" />
        </div>
        <div>
          <label class="block text-gray-700 mb-1">Password</label>
          <input v-model="form.password" type="password" required class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400" />
        </div>
        <button type="submit" class="w-full bg-blue-700 text-white py-2 rounded hover:bg-blue-800 transition">Login</button>
      </form>
      <div v-if="error" class="text-red-600 mt-4 text-center">{{ error }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { loginUser } from '../api/auth';
import { decodeJWT } from '../utils/jwt';
import { debugJWT } from '../utils/debug-jwt';

const router = useRouter();
const auth = useAuthStore();
const form = ref({ username: '', password: '' });
const error = ref('');

async function onSubmit() {
  error.value = '';
  try {
    const { data } = await loginUser({ username: form.value.username, password: form.value.password });
    const claims = decodeJWT(data.access_token);
    console.log('Decoded JWT claims:', claims); // Debug line
    const userType = data.role ? data.role.toLowerCase() : 'admin';
    auth.setAuth(data.access_token, userType, data.first_login, data.user);
    if (data.first_login) {
      router.push('/reset-password/user');
    } else {
      // Role-based redirect after login
      if (userType === 'admin') {
        router.push('/dashboard');
      } else if (userType === 'contact') {
        router.push('/contacts');
      } else if (userType === 'engineer') {
        router.push('/engineer/tickets');
      } else {
        router.push('/dashboard'); // fallback
      }
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Login failed';
  }
}
</script>
