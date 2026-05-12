<template>
  <header class="bg-blue-800 text-white px-4 py-4 flex items-center justify-between shadow relative">
    <div class="flex items-center gap-3">
      <button @click="$emit('toggleSidebar')" class="flex items-center justify-center w-10 h-10 rounded-full bg-blue-700 hover:bg-blue-600 focus:outline-none">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <g>
            <rect x="4" y="7" width="16" height="2" rx="1" fill="currentColor" />
            <rect x="4" y="11" width="16" height="2" rx="1" fill="currentColor" />
            <rect x="4" y="15" width="16" height="2" rx="1" fill="currentColor" />
          </g>
        </svg>
      </button>
      <div class="font-bold text-xl">GlobX Ticketing</div>
    </div>
    <div class="flex items-center gap-4">
      <!-- Notification Bell -->
      <NotificationBell v-if="user" />
      
      <!-- User Info and Logout -->
      <div class="flex items-center gap-4">
        <span v-if="user" class="text-sm">Welcome, {{ user.first_name || user.email }}</span>
        <button v-if="user" @click="logout" class="bg-blue-600 px-3 py-1 rounded hover:bg-blue-700 text-sm">Logout</button>
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
