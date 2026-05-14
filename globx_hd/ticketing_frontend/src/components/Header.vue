<template>
  <header class="bg-white text-gray-800 px-4 py-4 flex items-center justify-between shadow-md relative">
    <div class="flex items-center gap-3">
      <button @click="$emit('toggleSidebar')" class="flex items-center justify-center w-10 h-10 rounded-full bg-brand-teal hover:bg-teal-600 focus:outline-none transition-all duration-200">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <g>
            <rect x="4" y="7" width="16" height="2" rx="1" fill="currentColor" />
            <rect x="4" y="11" width="16" height="2" rx="1" fill="currentColor" />
            <rect x="4" y="15" width="16" height="2" rx="1" fill="currentColor" />
          </g>
        </svg>
      </button>
      <img src="/logo-globx-5.png" alt="GlobX" class="h-10" />
    </div>
    <div class="flex items-center gap-4">
      <!-- Notification Bell -->
      <NotificationBell v-if="user" />
      
      <!-- User Info and Logout -->
      <div class="flex items-center gap-4">
        <span v-if="user" class="text-sm text-neutral-medium font-medium">Welcome, {{ user.first_name || user.email }}</span>
        <button v-if="user" @click="logout" class="bg-gradient-to-r from-brand-teal to-brand-cyan text-white px-4 py-2 rounded-lg hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200 text-sm font-medium">Logout</button>
      </div>
    </div>
  </header>
  
</template>

<script setup>
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import { onMounted } from 'vue';
import { decodeJWT } from '../utils/jwt';
import NotificationBell from './shared/NotificationBell.vue';
const auth = useAuthStore();
const router = useRouter();
const user = auth.user;
const userType = auth.userType;
function logout() {
  auth.clearAuth();
  if (userType === 'contact') {
    router.push('/login/contact');
  } else {
    router.push('/login/user');
  }
}
onMounted(() => {
  if (auth.token) {
    const claims = decodeJWT(auth.token);
    console.log('Decoded JWT claims in Header:', claims);
  } else {
    console.log('No JWT token found in Header.');
  }
});
</script>
