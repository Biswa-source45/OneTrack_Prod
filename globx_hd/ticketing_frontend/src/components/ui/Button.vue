<template>
  <button
    :class="[
      'inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed',
      sizeClasses,
      variantClasses
    ]"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>

<script setup>
import { computed } from 'vue';
const props = defineProps({
  variant: { type: String, default: 'primary' },
  size: { type: String, default: 'md' },
});
const sizeClasses = computed(() => {
  if (props.size === 'sm') return 'px-3 py-1 text-sm';
  if (props.size === 'lg') return 'px-5 py-3 text-base';
  return 'px-4 py-2 text-sm';
});
const variantClasses = computed(() => {
  switch (props.variant) {
    case 'secondary':
      return 'bg-white text-blue-800 border border-blue-200 hover:bg-blue-50';
    case 'danger':
      return 'bg-red-600 text-white hover:bg-red-700';
    default:
      return 'bg-blue-600 text-white hover:bg-blue-700';
  }
});
</script>
