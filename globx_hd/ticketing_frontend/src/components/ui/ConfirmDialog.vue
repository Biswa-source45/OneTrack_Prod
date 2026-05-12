<template>
  <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
    <div class="absolute inset-0 bg-black/40" @click="$emit('cancel')"></div>
    <div class="relative bg-white rounded-lg shadow-lg w-full max-w-md mx-4 p-6">
      <h3 class="text-lg font-semibold text-blue-900 mb-2">{{ title }}</h3>
      <p class="text-sm text-blue-700 mb-4">{{ message }}</p>
      <div class="flex justify-end gap-2">
        <button class="px-4 py-2 rounded border border-blue-200 text-blue-800 hover:bg-blue-50" @click="$emit('cancel')">Cancel</button>
        <button class="px-4 py-2 rounded bg-red-600 text-white hover:bg-red-700" @click="$emit('confirm')">Confirm</button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({ open: Boolean, title: String, message: String });
</script>
