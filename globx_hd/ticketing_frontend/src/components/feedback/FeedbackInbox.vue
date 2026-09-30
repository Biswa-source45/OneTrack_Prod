<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-7xl mx-auto">
    <div class="flex flex-wrap items-center justify-between gap-3 mb-6">
      <h1 class="text-2xl font-bold text-blue-800">Feedback Inbox</h1>
      <div class="flex gap-2">
        <select v-model="statusFilter" @change="load" aria-label="Filter by status" class="border rounded px-3 py-2 text-sm">
          <option value="">All statuses</option>
          <option v-for="s in FEEDBACK_STATUSES" :key="s" :value="s">{{ s }}</option>
        </select>
        <button @click="load" class="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700">Refresh</button>
      </div>
    </div>

    <div v-if="error" class="mb-4 rounded border border-red-200 bg-red-50 px-4 py-2 text-sm text-red-700" role="alert">{{ error }}</div>
    <div v-if="loading" class="text-center py-8 text-gray-500">Loading…</div>
    <div v-else-if="items.length === 0" class="text-center py-12 text-sm text-gray-500">No feedback{{ statusFilter ? ` with status ${statusFilter}` : '' }}.</div>

    <div v-else class="space-y-4">
      <div v-for="f in items" :key="f.id" class="border border-gray-200 rounded-lg p-4" :data-testid="`feedback-${f.id}`">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div class="min-w-0">
            <h3 class="text-lg font-semibold text-gray-900">{{ f.title }}</h3>
            <p class="text-sm text-gray-500">
              {{ f.reporter_name }} <span v-if="f.reporter_email">&lt;{{ f.reporter_email }}&gt;</span>
              · {{ f.reporter_type === 'contact' ? 'Customer' : 'Staff' }}
            </p>
          </div>
          <select :value="f.status" @change="changeStatus(f, $event.target.value)" :disabled="saving === f.id"
                  :aria-label="`Status of feedback #${f.id}`"
                  class="px-2.5 py-1 rounded-full text-xs font-medium border" :class="FEEDBACK_STATUS_CLASSES[f.status]">
            <option v-for="s in FEEDBACK_STATUSES" :key="s" :value="s">{{ s }}</option>
          </select>
        </div>
        <p class="text-sm text-gray-700 mt-3 whitespace-pre-line">{{ f.description }}</p>
        <FeedbackImage v-if="f.image_path" :id="f.id" />
        <div class="flex flex-wrap justify-between gap-2 text-xs text-gray-500 border-t pt-3 mt-3">
          <span>#{{ f.id }}<span v-if="f.page_url"> · Page: {{ f.page_url }}</span></span>
          <span>{{ formatDateTime(f.created_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { fetchAllFeedback, updateFeedbackStatus, FEEDBACK_STATUSES, FEEDBACK_STATUS_CLASSES } from '../../api/feedback';
import { formatDateTime } from '../../utils/date';
import FeedbackImage from './FeedbackImage.vue';

const items = ref([]);
const statusFilter = ref('');
const loading = ref(true);
const saving = ref(null);
const error = ref('');

async function load() {
  loading.value = true;
  error.value = '';
  try {
    items.value = await fetchAllFeedback(statusFilter.value);
  } catch (e) {
    error.value = e.response?.data?.error || 'Could not load feedback.';
  } finally {
    loading.value = false;
  }
}

async function changeStatus(f, status) {
  const old = f.status;
  f.status = status;
  saving.value = f.id;
  error.value = '';
  try {
    await updateFeedbackStatus(f.id, status);
  } catch (e) {
    f.status = old;
    error.value = e.response?.data?.error || `Could not update #${f.id}.`;
  } finally {
    saving.value = null;
  }
}

onMounted(load);
</script>
