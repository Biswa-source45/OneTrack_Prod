<template>
  <div class="flex items-center justify-between mb-6">
    <div>
      <h2 class="text-2xl font-bold text-blue-800">{{ title }}</h2>
      <p v-if="subtitle" class="text-sm text-blue-600 mt-1">{{ subtitle }}</p>
    </div>
    <div class="flex items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
const props = defineProps({ title: String, subtitle: String });
</script>
