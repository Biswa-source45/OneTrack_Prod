<template>
  <form @submit.prevent="$emit('submit')" class="bg-white rounded-lg shadow-sm border border-blue-100 p-6 max-w-3xl">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <slot />
    </div>
    <div class="flex justify-end gap-2 mt-6">
      <slot name="actions" />
    </div>
  </form>
</template>

<script setup>
</script>
