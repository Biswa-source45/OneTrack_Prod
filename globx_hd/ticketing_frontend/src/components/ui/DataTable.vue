<template>
  <div class="bg-white rounded-lg shadow-sm border border-blue-100 overflow-hidden">
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-blue-100">
        <thead class="bg-blue-50">
          <tr>
            <th v-for="col in columns" :key="col.key" class="px-4 py-3 text-left text-xs font-semibold text-blue-800 uppercase tracking-wider whitespace-nowrap">
              {{ col.label }}
            </th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-blue-800 uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-blue-50">
          <tr v-for="row in pagedRows" :key="row[rowKey]"
              :class="[rowClickable ? 'cursor-pointer transition hover:bg-blue-100' : 'hover:bg-blue-50/50']"
              @click="rowClickable ? $emit('row-click', row) : null">
            <td v-for="col in columns" :key="col.key" class="px-4 py-3 text-sm text-blue-900">
              <slot :name="`cell:${col.key}`" :row="row">{{ row[col.key] }}</slot>
            </td>
            <td class="px-4 py-3 text-right">
              <slot name="row-actions" :row="row" />
            </td>
          </tr>
          <tr v-if="!rows.length">
            <td :colspan="columns.length + 1" class="px-4 py-10 text-center text-blue-700">No data</td>
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
