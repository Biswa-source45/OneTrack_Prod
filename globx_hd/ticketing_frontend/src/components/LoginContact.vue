<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-blue-100">
    <div class="bg-white shadow-lg rounded-lg p-8 w-full max-w-md">
      <h2 class="text-2xl font-bold mb-6 text-center text-blue-800">Customer Login</h2>
      <form @submit.prevent="onSubmit" class="space-y-4">
        <div>
          <label class="block text-gray-700 mb-1">Email</label>
          <input v-model="form.username" type="email" required class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400" />
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
import { loginContact } from '../api/auth';

const router = useRouter();
const auth = useAuthStore();
const form = ref({ username: '', password: '' });
const error = ref('');

async function onSubmit() {
  error.value = '';
  try {
    const { data } = await loginContact({ username: form.value.username, password: form.value.password });
    console.log('LoginContact API response:', data);
  auth.setAuth(data.access_token, 'contact', data.first_login, data.contact);
    console.log('auth.user after setAuth:', auth.user);
    if (data.first_login) {
      router.push('/reset-password/contact');
    } else {
      router.push('/contacts');
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Login failed';
  }
}
</script>
