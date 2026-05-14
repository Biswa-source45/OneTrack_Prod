<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-neutral-light to-white">
    <div class="bg-white shadow-2xl rounded-2xl p-8 w-full max-w-md border border-gray-100">
      <div class="flex justify-center mb-6">
        <img src="/logo-globx-5.png" alt="GlobX" class="h-16" />
      </div>
      <h2 class="text-2xl font-bold mb-6 text-center text-neutral-dark">User Login</h2>
      <form @submit.prevent="onSubmit" class="space-y-4">
        <div>
          <label class="block text-gray-700 mb-1">Username</label>
          <input v-model="form.username" type="text" required class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-brand-teal focus:border-transparent transition-all" />
        </div>
        <div>
          <label class="block text-gray-700 mb-1">Password</label>
          <input v-model="form.password" type="password" required class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-brand-teal focus:border-transparent transition-all" />
        </div>
        <button type="submit" class="w-full bg-gradient-to-r from-brand-teal to-brand-cyan text-white py-3 rounded-lg shadow-md hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200 font-medium">Login</button>
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
