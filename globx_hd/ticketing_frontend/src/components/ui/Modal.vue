<template>
  <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
    <div class="absolute inset-0 bg-black/40" @click="$emit('close')"></div>
    <div class="relative bg-white rounded-lg shadow-lg w-full max-w-2xl mx-4 p-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-blue-900">{{ title }}</h3>
        <button class="w-8 h-8 rounded hover:bg-blue-50" @click="$emit('close')">✕</button>
      </div>
      <div>
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({ open: Boolean, title: String });
</script>

