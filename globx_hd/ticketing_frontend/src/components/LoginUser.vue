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
          <div class="relative">
            <input 
              v-model="form.password" 
              :type="showPassword ? 'text' : 'password'" 
              required 
              class="w-full pl-3 pr-10 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400" 
            />
            <button 
              type="button" 
              @click="showPassword = !showPassword" 
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500 hover:text-gray-700 focus:outline-none"
            >
              <svg v-if="showPassword" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
              <svg v-else class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
            </button>
          </div>
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
const showPassword = ref(false);
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
