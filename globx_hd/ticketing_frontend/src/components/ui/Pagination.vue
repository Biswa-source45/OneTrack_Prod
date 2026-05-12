<template>
  <div class="flex items-center justify-between mt-4">
    <div class="text-sm text-blue-700">Page {{ page }} of {{ totalPages }}</div>
    <div class="flex items-center gap-2">
      <button class="px-3 py-1 rounded border border-blue-200 text-blue-800 disabled:opacity-50" :disabled="page===1" @click="$emit('update:page', page-1)">Prev</button>
      <button class="px-3 py-1 rounded border border-blue-200 text-blue-800 disabled:opacity-50" :disabled="page===totalPages" @click="$emit('update:page', page+1)">Next</button>
      <select class="ml-2 border border-blue-200 rounded px-2 py-1 text-sm" :value="pageSize" @change="$emit('update:pageSize', parseInt($event.target.value))">
        <option :value="10">10</option>
        <option :value="20">20</option>
        <option :value="50">50</option>
      </select>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
const props = defineProps({ page: Number, pageSize: Number, total: Number });
const totalPages = computed(() => Math.max(1, Math.ceil((props.total || 0) / (props.pageSize || 10))));
</script>
