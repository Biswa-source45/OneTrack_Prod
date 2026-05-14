<template>
  <div class="flex flex-col min-h-screen bg-neutral-light">
    <Header :sidebarOpen="sidebarOpen" @toggleSidebar="toggleSidebar" />
    <div class="flex flex-1 relative">
      <Sidebar :sidebarOpen="sidebarOpen" />
      <main class="flex-1 p-8 transition-all duration-300">
        <router-view />
      </main>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import Header from './Header.vue';
import Sidebar from './Sidebar.vue';
import Footer from './Footer.vue';
import { ref } from 'vue';
const sidebarOpen = ref(true);
function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value;
}
</script>
