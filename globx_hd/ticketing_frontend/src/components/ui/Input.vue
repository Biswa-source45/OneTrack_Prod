<template>
  <input
    :type="type"
    :value="modelValue"
    @input="$emit('update:modelValue', $event.target.value)"
    :placeholder="placeholder"
    :disabled="disabled"
    :required="required"
    :class="[
      'w-full px-3 py-2 border rounded-lg transition-all duration-200',
      'focus:outline-none focus:ring-2 focus:ring-brand-teal focus:border-transparent',
      'disabled:bg-gray-100 disabled:cursor-not-allowed',
      error ? 'border-red-500' : 'border-gray-300',
      inputClass
    ]"
    v-bind="$attrs"
  />
</template>

<script setup>
defineProps({
  type: { type: String, default: 'text' },
  modelValue: [String, Number],
  placeholder: String,
  disabled: Boolean,
  required: Boolean,
  error: Boolean,
  inputClass: String
});

defineEmits(['update:modelValue']);
</script>
