<template>
  <span
    :class="[
      'inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium shadow-sm transition-all duration-200',
      variantClasses
    ]"
  >
    <span v-if="icon" class="w-2 h-2 rounded-full" :class="dotClass"></span>
    <slot />
  </span>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  variant: { 
    type: String, 
    default: 'default',
    validator: (value) => [
      'default', 'open', 'in_progress', 'completed', 'closed', 
      'pending', 'approved', 'rejected', 'high', 'urgent', 'low', 'medium',
      'success', 'warning', 'error', 'info'
    ].includes(value)
  },
  icon: { type: Boolean, default: false }
});

const variantClasses = computed(() => {
  const variants = {
    // Status variants
    'open': 'bg-brand-teal text-white',
    'in_progress': 'bg-brand-cyan text-white',
    'completed': 'bg-emerald-500 text-white',
    'closed': 'bg-gray-500 text-white',
    'pending': 'bg-amber-500 text-white',
    'approved': 'bg-emerald-500 text-white',
    'rejected': 'bg-red-500 text-white',
    
    // Priority variants
    'urgent': 'bg-red-600 text-white animate-pulse',
    'high': 'bg-orange-500 text-white',
    'medium': 'bg-amber-500 text-white',
    'low': 'bg-gray-400 text-white',
    
    // Semantic variants
    'success': 'bg-emerald-500 text-white',
    'warning': 'bg-amber-500 text-white',
    'error': 'bg-red-500 text-white',
    'info': 'bg-brand-cyan text-white',
    
    // Default
    'default': 'bg-gray-200 text-gray-800'
  };
  
  return variants[props.variant] || variants.default;
});

const dotClass = computed(() => {
  if (props.variant === 'urgent') return 'bg-white animate-pulse';
  return 'bg-white/80';
});
</script>
