<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <div>
      <h1 class="text-2xl font-bold text-blue-800">Feedback</h1>
      <p class="text-sm text-blue-600 mt-1">
        Found a bug or something confusing? Tell us what happened.
        Press <kbd class="px-1.5 py-0.5 rounded border border-blue-200 bg-white text-xs font-mono">Ctrl</kbd> + <kbd class="px-1.5 py-0.5 rounded border border-blue-200 bg-white text-xs font-mono">I</kbd>
        on any page to screenshot it and land here with the image attached.
      </p>
    </div>

    <form @submit.prevent="submit" class="bg-white rounded-lg shadow-md p-6 space-y-4" novalidate>
      <div v-if="success" class="rounded border border-green-200 bg-green-50 px-4 py-2 text-sm text-green-800" role="status">{{ success }}</div>
      <div v-if="errors.submit" class="rounded border border-red-200 bg-red-50 px-4 py-2 text-sm text-red-700" role="alert">{{ errors.submit }}</div>

      <div>
        <label for="fb-title" class="block text-sm font-medium text-gray-700 mb-1">Title</label>
        <input id="fb-title" v-model="title" maxlength="200" placeholder="Short summary of the issue"
               class="w-full border rounded px-3 py-2" :class="errors.title ? 'border-red-400' : 'border-blue-200'" />
        <p v-if="errors.title" class="text-xs text-red-600 mt-1">{{ errors.title }}</p>
      </div>

      <div>
        <label for="fb-desc" class="block text-sm font-medium text-gray-700 mb-1">What happened?</label>
        <textarea id="fb-desc" v-model="description" rows="5" placeholder="What were you doing, what did you expect, and what happened instead?"
                  class="w-full border rounded px-3 py-2" :class="errors.description ? 'border-red-400' : 'border-blue-200'"></textarea>
        <p v-if="errors.description" class="text-xs text-red-600 mt-1">{{ errors.description }}</p>
      </div>

      <div>
        <span class="block text-sm font-medium text-gray-700 mb-1">Screenshot <span class="text-gray-400 font-normal">(optional · JPEG, PNG or WebP · max 5 MB)</span></span>
        <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/webp" class="hidden" data-testid="feedback-file"
               @change="e => e.target.files[0] && setImage(e.target.files[0], false)" />
        <div v-if="previewUrl" class="relative inline-block">
          <img :src="previewUrl" alt="Attached screenshot preview" class="max-h-48 rounded border border-blue-200" />
          <button type="button" @click="clearImage" aria-label="Remove screenshot"
                  class="absolute -top-2 -right-2 w-6 h-6 rounded-full bg-gray-800 text-white text-sm leading-6 text-center hover:bg-gray-700">×</button>
          <span v-if="fromCapture" class="absolute bottom-1 left-1 text-[11px] px-1.5 py-0.5 rounded bg-white/90 border border-blue-200 text-blue-800">Captured automatically</span>
        </div>
        <button v-else type="button" @click="fileInput.click()"
                class="px-3 py-2 rounded border border-dashed border-blue-300 text-sm text-blue-700 hover:bg-blue-50">
          Attach a screenshot
        </button>
        <p v-if="errors.image" class="text-xs text-red-600 mt-1">{{ errors.image }}</p>
      </div>

      <div class="flex justify-end">
        <button type="submit" :disabled="submitting"
                class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-60">
          {{ submitting ? 'Submitting…' : 'Submit Feedback' }}
        </button>
      </div>
    </form>

    <div class="bg-white rounded-lg shadow-md p-6">
      <h2 class="text-lg font-semibold text-blue-800 mb-3">My feedback</h2>
      <div v-if="loading" class="text-sm text-gray-500 py-4">Loading…</div>
      <div v-else-if="loadError" class="text-sm text-red-600 py-4">{{ loadError }}</div>
      <div v-else-if="mine.length === 0" class="text-sm text-gray-500 py-4">Nothing submitted yet.</div>
      <ul v-else class="divide-y divide-gray-100">
        <li v-for="f in mine" :key="f.id" class="py-3">
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <p class="font-medium text-gray-900 truncate">{{ f.title }}</p>
              <p class="text-xs text-gray-500">#{{ f.id }} · {{ formatDateTime(f.created_at) }}</p>
            </div>
            <span class="shrink-0 px-2.5 py-0.5 rounded-full text-xs font-medium border" :class="FEEDBACK_STATUS_CLASSES[f.status]">{{ f.status }}</span>
          </div>
          <p class="text-sm text-gray-700 mt-1 whitespace-pre-line line-clamp-3">{{ f.description }}</p>
          <FeedbackImage v-if="f.image_path" :id="f.id" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted, onBeforeUnmount } from 'vue';
import { createFeedback, fetchMyFeedback, FEEDBACK_STATUS_CLASSES } from '../../api/feedback';
import { feedbackDraft } from '../../utils/feedbackCapture';
import { formatDateTime } from '../../utils/date';
import FeedbackImage from './FeedbackImage.vue';

const MAX_IMAGE_BYTES = 5 * 1024 * 1024;
const IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/webp'];

const title = ref('');
const description = ref('');
const image = ref(null);
const previewUrl = ref(null);
const fromCapture = ref(false);
const pageUrl = ref('');
const fileInput = ref(null);
const errors = reactive({});
const success = ref('');
const submitting = ref(false);

const mine = ref([]);
const loading = ref(true);
const loadError = ref('');

function setImage(file, captured) {
  errors.image = '';
  if (!IMAGE_TYPES.includes(file.type)) {
    errors.image = 'Only JPEG, PNG or WebP images are allowed.';
    return;
  }
  if (file.size > MAX_IMAGE_BYTES) {
    errors.image = `Image is ${(file.size / 1048576).toFixed(1)} MB — it must be 5 MB or smaller.`;
    return;
  }
  clearImage();
  image.value = file;
  previewUrl.value = URL.createObjectURL(file);
  fromCapture.value = captured;
}

function clearImage() {
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value);
  image.value = null;
  previewUrl.value = null;
  fromCapture.value = false;
  if (fileInput.value) fileInput.value.value = '';
}

// Ctrl+I drops a draft here (also when we're already on this page).
watch(feedbackDraft, (draft) => {
  if (!draft) return;
  feedbackDraft.value = null;
  success.value = '';
  title.value = draft.title;
  pageUrl.value = draft.pageUrl;
  if (draft.file) setImage(draft.file, true);
  else errors.image = 'Could not capture the screen automatically — you can attach a screenshot manually.';
}, { immediate: true });

async function loadMine() {
  loading.value = true;
  loadError.value = '';
  try {
    mine.value = await fetchMyFeedback();
  } catch (e) {
    loadError.value = e.response?.data?.error || 'Could not load your feedback.';
  } finally {
    loading.value = false;
  }
}

async function submit() {
  success.value = '';
  errors.title = title.value.trim() ? '' : 'Give it a short title.';
  errors.description = description.value.trim() ? '' : 'Describe what happened.';
  errors.submit = '';
  // errors.image doesn't block: a rejected file is never attached, the note just says so.
  if (errors.title || errors.description) return;

  submitting.value = true;
  try {
    await createFeedback({ title: title.value.trim(), description: description.value.trim(), pageUrl: pageUrl.value, image: image.value });
    success.value = 'Thanks! Your feedback was sent to the Super Admin.';
    title.value = '';
    description.value = '';
    pageUrl.value = '';
    clearImage();
    errors.image = '';
    loadMine();
  } catch (e) {
    errors.submit = e.response?.data?.error || 'Could not submit feedback. Please try again.';
  } finally {
    submitting.value = false;
  }
}

onMounted(loadMine);
onBeforeUnmount(clearImage);
</script>
