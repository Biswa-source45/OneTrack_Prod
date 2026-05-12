<template>
  <div>
    <label v-if="label" class="block text-sm font-medium text-blue-900 mb-1">{{ label }}</label>
    <slot />
    <p v-if="error" class="text-sm text-red-600 mt-1">{{ error }}</p>
  </div>
</template>

<script setup>
const props = defineProps({ label: String, error: String });
</script>
