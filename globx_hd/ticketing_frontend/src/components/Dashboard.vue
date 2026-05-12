<template>
  <div class="flex flex-col items-center justify-center h-full">
    <h1 class="text-3xl font-bold text-blue-800 mb-4">Hello, {{ userName }}</h1>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';
const auth = useAuthStore();
const userName = auth.user?.first_name || auth.user?.name || auth.user?.email || 'User';
</script>
