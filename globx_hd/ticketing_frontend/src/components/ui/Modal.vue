<template>
  <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
    <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="$emit('close')"></div>
    <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-2xl mx-4 p-6 border border-gray-100">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-neutral-dark">{{ title }}</h3>
        <button class="w-8 h-8 rounded-full hover:bg-teal-50 text-neutral-medium hover:text-brand-teal transition-all duration-200" @click="$emit('close')">✕</button>
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

