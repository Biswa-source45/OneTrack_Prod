<template>
  <button
    :class="[
      'inline-flex items-center justify-center rounded-lg font-medium transition-all duration-200 focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed',
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
      return 'bg-white text-brand-teal border border-teal-200 hover:bg-teal-50 hover:shadow-md';
    case 'danger':
      return 'bg-red-600 text-white hover:bg-red-700 shadow-md hover:shadow-lg hover:-translate-y-0.5';
    default:
      return 'bg-gradient-to-r from-brand-teal to-brand-cyan text-white shadow-md hover:shadow-lg hover:-translate-y-0.5';
  }
});
</script>
