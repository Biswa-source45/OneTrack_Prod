<template>
  <div class="flex flex-col min-h-screen bg-blue-50">
    <Header :sidebarOpen="sidebarOpen" @toggleSidebar="toggleSidebar" />
    <div class="flex flex-1 relative">
      <Sidebar :sidebarOpen="sidebarOpen" />
      <main class="flex-1 p-8 transition-all duration-300">
        <router-view />
      </main>
    </div>
    <Footer />
    <div v-if="capturing" :id="CAPTURE_OVERLAY_ID" class="fixed inset-0 z-[100] bg-black/40 flex items-center justify-center" role="status" aria-live="polite">
      <div class="bg-white rounded-lg shadow-lg px-6 py-4 flex items-center gap-3 text-blue-800 font-medium">
        <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24" fill="none"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" /><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" /></svg>
        Capturing screenshot…
      </div>
    </div>
  </div>
</template>

<script setup>
import Header from './Header.vue';
import Sidebar from './Sidebar.vue';
import Footer from './Footer.vue';
import { ref, nextTick, onMounted, onBeforeUnmount } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { feedbackDraft, capturePageScreenshot, pageNameFromPath, CAPTURE_OVERLAY_ID } from '../utils/feedbackCapture';
const sidebarOpen = ref(true);
function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value;
}

// Ctrl+I anywhere: screenshot the page and open Feedback with it attached.
// Super Admin only receives feedback, so the shortcut is inert for them.
const auth = useAuthStore();
const route = useRoute();
const router = useRouter();
const capturing = ref(false);
async function onKeyDown(e) {
  // Exclude Shift/Alt/Meta: Ctrl+Shift+I is DevTools.
  if (!e.ctrlKey || e.shiftKey || e.altKey || e.metaKey || e.key.toLowerCase() !== 'i') return;
  if (auth.userType === 'superadmin') return;
  e.preventDefault();
  if (capturing.value) return;
  capturing.value = true;
  const draft = { file: null, title: `Issue on: ${pageNameFromPath(route.path)}`, pageUrl: route.fullPath };
  try {
    await nextTick(); // let the overlay render before the (slow) DOM clone starts
    draft.file = await capturePageScreenshot();
  } catch (err) {
    console.error('Screenshot capture failed:', err); // form still opens; user can attach manually
  } finally {
    capturing.value = false;
  }
  feedbackDraft.value = draft;
  router.push('/feedback');
}
onMounted(() => window.addEventListener('keydown', onKeyDown));
onBeforeUnmount(() => window.removeEventListener('keydown', onKeyDown));
</script>
