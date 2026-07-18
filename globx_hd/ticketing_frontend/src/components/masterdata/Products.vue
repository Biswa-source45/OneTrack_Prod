<template>
  <div class="p-8">
    <PageHeader title="Products">
      <template #actions>
        <Button @click="goCreate"><span class="mr-2">+</span>Add Product</Button>
      </template>
    </PageHeader>

    <DataTable :columns="columns" :rows="rows" rowKey="id">
      <template #row-actions="{ row }">
        <DropdownMenu>
          <template #button="{ toggle }">
            <button @click="toggle" class="inline-flex items-center justify-center w-8 h-8 rounded hover:bg-blue-100">
              ⋯
            </button>
          </template>
          <div class="py-1">
            <button class="w-full text-left px-3 py-2 text-sm hover:bg-blue-50" @click="goEdit(row)">Edit</button>
            <button class="w-full text-left px-3 py-2 text-sm text-red-700 hover:bg-red-50" @click="askDelete(row)">Delete</button>
          </div>
        </DropdownMenu>
      </template>
    </DataTable>

    <ConfirmDialog :open="confirmOpen" title="Delete product" :message="`Are you sure you want to delete “${selected?.product_name}”?`" @cancel="confirmOpen=false" @confirm="confirmDelete" />
  </div>
  
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import PageHeader from '../ui/PageHeader.vue';
import Button from '../ui/Button.vue';
import DataTable from '../ui/DataTable.vue';
import DropdownMenu from '../ui/DropdownMenu.vue';
import ConfirmDialog from '../ui/ConfirmDialog.vue';
import { fetchProducts, deleteProduct } from '../../api/auth';
import { useAuthStore } from '../../stores/auth';

const router = useRouter();
const authStore = useAuthStore();
const rows = ref([]);
const loading = ref(false);
const columns = [
  { key: 'product_name', label: 'Product Name' },
  { key: 'product_description', label: 'Description' },
  { key: 'created_at', label: 'Created' },
];

onMounted(load);
async function load() {
  loading.value = true;
  try {
    const data = await fetchProducts();
    rows.value = data;
  } finally {
    loading.value = false;
  }
}

function goCreate() {
  const prefix = authStore.userType === 'manager' ? '/manager' : '';
  router.push(`${prefix}/master-data/products/new`);
}
function goEdit(row) {
  const prefix = authStore.userType === 'manager' ? '/manager' : '';
  router.push(`${prefix}/master-data/products/${row.id}/edit`);
}

const confirmOpen = ref(false);
const selected = ref(null);
function askDelete(row) {
  selected.value = row;
  confirmOpen.value = true;
}
async function confirmDelete() {
  confirmOpen.value = false;
  if (!selected.value) return;
  await deleteProduct(selected.value.id).catch(() => {});
  await load();
}
</script>
