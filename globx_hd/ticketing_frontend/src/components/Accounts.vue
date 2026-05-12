<template>
  <div class="p-8">
    <PageHeader title="Accounts">
      <template #actions>
        <Button @click="openCreate=true"><span class="mr-2">+</span>Add Account</Button>
      </template>
    </PageHeader>

    <DataTable :columns="columns" :rows="rows" rowKey="id">
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

    <ConfirmDialog :open="confirmOpen" title="Delete account" :message="`Delete account “${selected?.account_name}”?`" @cancel="confirmOpen=false" @confirm="confirmDelete" />

    <Modal :open="openCreate || openEdit" :title="openEdit ? 'Edit Account' : 'Add Account'" @close="closeForm">
      <FormLayout @submit="save">
        <FormField label="Account Name"><input v-model="form.account_name" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Account Owner"><input v-model="form.account_owner" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Address"><input v-model="form.address" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <template #actions>
          <Button variant="secondary" type="button" @click="closeForm">Cancel</Button>
          <Button type="submit">Save</Button>
        </template>
      </FormLayout>
      <p v-if="openEdit" class="mt-2 text-xs text-blue-600">Customer Code: {{ form.customer_code }}</p>
    </Modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import PageHeader from './ui/PageHeader.vue';
import Button from './ui/Button.vue';
import DataTable from './ui/DataTable.vue';
import DropdownMenu from './ui/DropdownMenu.vue';
import ConfirmDialog from './ui/ConfirmDialog.vue';
import Modal from './ui/Modal.vue';
import FormLayout from './ui/FormLayout.vue';
import FormField from './ui/FormField.vue';
import { fetchAccounts, createAccount, updateAccount, deleteAccount } from '../api/auth';

const rows = ref([]);
const columns = [
  { key: 'account_name', label: 'Account Name' },
  { key: 'account_owner', label: 'Owner' },
  { key: 'customer_code', label: 'Customer Code' },
  { key: 'address', label: 'Address' },
];

const openCreate = ref(false); const openEdit = ref(false); const editingId = ref(null);
const form = reactive({ account_name: '', account_owner: '', address: '', customer_code: '' });

onMounted(load);
async function load(){ rows.value = await fetchAccounts(); }
function closeForm(){ openCreate.value=false; openEdit.value=false; editingId.value=null; Object.assign(form, { account_name:'', account_owner:'', address:'', customer_code:'' }); }
function startEdit(row){ openEdit.value=true; editingId.value=row.id; Object.assign(form, { account_name: row.account_name, account_owner: row.account_owner, address: row.address, customer_code: row.customer_code }); }
const submitting = ref(false);
async function save(){
  if (submitting.value) return;
  submitting.value = true;
  try {
    if(openEdit.value){ await updateAccount(editingId.value, { account_name: form.account_name, account_owner: form.account_owner, address: form.address }); }
    else { await createAccount({ account_name: form.account_name, account_owner: form.account_owner, address: form.address }); }
    closeForm(); await load();
  } finally {
    submitting.value = false;
  }
}

const confirmOpen = ref(false); const selected = ref(null);
function askDelete(row){ selected.value=row; confirmOpen.value=true; }
async function confirmDelete(){ confirmOpen.value=false; if(!selected.value) return; await deleteAccount(selected.value.id); await load(); }
</script>
