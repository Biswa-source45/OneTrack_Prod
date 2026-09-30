<template>
  <div class="mt-2">
    <button v-if="!url" type="button" @click="load" :disabled="loading" class="text-sm text-blue-700 hover:underline disabled:opacity-60">
      {{ loading ? 'Loading screenshot…' : 'View screenshot' }}
    </button>
    <a v-else :href="url" target="_blank" rel="noopener" title="Open full size">
      <img :src="url" alt="Attached screenshot" class="max-h-72 rounded border border-blue-200" />
    </a>
    <p v-if="error" class="text-xs text-red-600">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue';
import { fetchFeedbackImageUrl } from '../../api/feedback';

const props = defineProps({ id: { type: Number, required: true } });
const url = ref(null);
const loading = ref(false);
const error = ref('');

async function load() {
  loading.value = true;
  error.value = '';
  try {
    url.value = await fetchFeedbackImageUrl(props.id);
  } catch {
    error.value = 'Could not load the screenshot.';
  } finally {
    loading.value = false;
  }
}

onBeforeUnmount(() => url.value && URL.revokeObjectURL(url.value));
</script>
