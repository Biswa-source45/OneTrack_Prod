<template>
  <div class="p-8">
    <PageHeader title="Issues">
      <template #actions>
        <select v-model.number="filterProductId" class="mr-3 border border-blue-200 rounded px-2 py-2 text-sm">
          <option :value="0">All Products</option>
          <option v-for="p in products" :key="p.id" :value="p.id">{{ p.product_name }}</option>
        </select>
        <Button @click="openCreate=true"><span class="mr-2">+</span>Add Issue</Button>
      </template>
    </PageHeader>
    <DataTable :columns="columns" :rows="viewRows" rowKey="id">
      <template #row-actions="{ row }">
        <DropdownMenu>
          <template #button="{ toggle }">
            <button @click="toggle" class="inline-flex items-center justify-center w-8 h-8 rounded hover:bg-blue-100">⋯</button>
          </template>
          <div class="py-1">
            <button class="w-full text-left px-3 py-2 text-sm" @click="startEdit(row)">Edit</button>
            <button class="w-full text-left px-3 py-2 text-sm text-red-700 hover:bg-red-50" @click="askDelete(row)">Delete</button>
          </div>
        </DropdownMenu>
      </template>
    </DataTable>

    <ConfirmDialog :open="confirmOpen" title="Delete issue" :message="`Delete “${selected?.issue_name}”?`" @cancel="confirmOpen=false" @confirm="confirmDelete" />

    <div v-if="openCreate || openEdit" class="p-6 bg-white rounded-lg shadow border border-blue-100 max-w-lg">
      <h3 class="text-lg font-semibold mb-4">{{ openEdit ? 'Edit Issue' : 'Add Issue' }}</h3>
      <div class="space-y-3">
        <select v-model.number="formProductId" class="w-full border border-blue-200 rounded px-3 py-2">
          <option disabled value="0">Select product</option>
          <option v-for="p in products" :key="p.id" :value="p.id">{{ p.product_name }}</option>
        </select>
        <input v-model="name" type="text" placeholder="Issue name" class="w-full border border-blue-200 rounded px-3 py-2" />
        <div class="flex justify-end gap-2">
          <Button variant="secondary" @click="closeForm">Cancel</Button>
          <Button @click="save">Save</Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue';
import PageHeader from '../ui/PageHeader.vue';
import Button from '../ui/Button.vue';
import DataTable from '../ui/DataTable.vue';
import DropdownMenu from '../ui/DropdownMenu.vue';
import ConfirmDialog from '../ui/ConfirmDialog.vue';
import { fetchIssues, createIssue, updateIssue, deleteIssue, fetchProducts } from '../../api/auth';

const rows = ref([]);
const columns = [ { key: 'product_name', label: 'Product' }, { key: 'issue_name', label: 'Issue' }, { key: 'created_at', label: 'Created' } ];
const openCreate = ref(false); const openEdit = ref(false); const editingId = ref(null); const name = ref('');
const products = ref([]); const filterProductId = ref(0); const formProductId = ref(0);

const viewRows = computed(() => {
  const map = new Map(products.value.map(p => [p.id, p.product_name]));
  return rows.value.map(r => ({ ...r, product_name: map.get(r.product_id) || '' }));
});

onMounted(async () => {
  products.value = await fetchProducts();
  await load();
});
watch(filterProductId, load);
async function load(){
  const params = filterProductId.value ? { product_id: filterProductId.value } : undefined;
  rows.value = await fetchIssues(params);
}

function closeForm(){ openCreate.value=false; openEdit.value=false; editingId.value=null; name.value=''; formProductId.value=0; }
function startEdit(row){ openEdit.value=true; editingId.value=row.id; name.value=row.issue_name; formProductId.value=row.product_id; }
async function save(){
  if(!formProductId.value){ return; }
  if(openEdit.value){ await updateIssue(editingId.value, { issue_name: name.value, product_id: formProductId.value }); }
  else { await createIssue({ issue_name: name.value, product_id: formProductId.value }); }
  closeForm(); await load();
}
const confirmOpen = ref(false); const selected = ref(null);
function askDelete(row){ selected.value=row; confirmOpen.value=true; }
async function confirmDelete(){ confirmOpen.value=false; if(!selected.value) return; await deleteIssue(selected.value.id); await load(); }
</script>
