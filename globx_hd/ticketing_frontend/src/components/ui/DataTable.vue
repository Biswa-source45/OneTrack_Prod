<template>
  <div class="bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden">
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gradient-to-r from-neutral-light to-white">
          <tr>
            <th v-for="col in columns" :key="col.key" class="px-4 py-3 text-left text-xs font-semibold text-neutral-dark uppercase tracking-wider whitespace-nowrap">
              {{ col.label }}
            </th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-neutral-dark uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100">
          <tr v-for="row in pagedRows" :key="row[rowKey]"
              :class="[rowClickable ? 'cursor-pointer transition-all duration-200 hover:bg-teal-50' : 'hover:bg-gray-50 transition-all duration-200']"
              @click="rowClickable ? $emit('row-click', row) : null">
            <td v-for="col in columns" :key="col.key" class="px-4 py-3 text-sm text-neutral-dark">
              <slot :name="`cell:${col.key}`" :row="row">{{ row[col.key] }}</slot>
            </td>
            <td class="px-4 py-3 text-right">
              <slot name="row-actions" :row="row" />
            </td>
          </tr>
          <tr v-if="!rows.length">
            <td :colspan="columns.length + 1" class="px-4 py-10 text-center text-neutral-medium">No data</td>
          </tr>
        </tbody>
      </table>
    </div>
    <Pagination :page="page" :pageSize="pageSize" :total="rows.length" @update:page="page=$event" @update:pageSize="pageSize=$event" />
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue';
import Pagination from './Pagination.vue';
const props = defineProps({ columns: Array, rows: { type: Array, default: () => [] }, rowKey: { type: String, default: 'id' }, rowClickable: { type: Boolean, default: false } });
const page = ref(1);
const pageSize = ref(10);
const rowClickable = props.rowClickable;
watch(pageSize, () => { page.value = 1; });
const pagedRows = computed(() => {
  const start = (page.value - 1) * pageSize.value;
  return props.rows.slice(start, start + pageSize.value);
});
</script>
