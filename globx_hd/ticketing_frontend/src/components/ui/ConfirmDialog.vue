<template>
  <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
    <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="$emit('cancel')"></div>
    <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md mx-4 p-6 border border-gray-100">
      <h3 class="text-lg font-semibold text-neutral-dark mb-2">{{ title }}</h3>
      <p class="text-sm text-neutral-medium mb-4">{{ message }}</p>
      <div class="flex justify-end gap-2">
        <button class="px-4 py-2 rounded-lg border border-teal-200 text-brand-teal hover:bg-teal-50 transition-all duration-200" @click="$emit('cancel')">Cancel</button>
        <button class="px-4 py-2 rounded-lg bg-red-600 text-white hover:bg-red-700 shadow-md hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200" @click="$emit('confirm')">Confirm</button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({ open: Boolean, title: String, message: String });
</script>
