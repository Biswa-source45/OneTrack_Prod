<template>
  <div class="p-8">
    <PageHeader title="User Roles">
      <template #actions>
        <Button @click="openCreate=true"><span class="mr-2">+</span>Add Role</Button>
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

    <ConfirmDialog :open="confirmOpen" title="Delete role" :message="`Delete “${selected?.role_name}”?`" @cancel="confirmOpen=false" @confirm="confirmDelete" />

    <div v-if="openCreate || openEdit" class="p-6 bg-white rounded-lg shadow border border-blue-100 max-w-lg">
      <h3 class="text-lg font-semibold mb-4">{{ openEdit ? 'Edit Role' : 'Add Role' }}</h3>
      <div class="space-y-3">
        <input v-model="roleName" type="text" placeholder="Role name" class="w-full border border-blue-200 rounded px-3 py-2" />
        <div class="flex justify-end gap-2">
          <Button variant="secondary" @click="closeForm">Cancel</Button>
          <Button @click="saveRole">Save</Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import PageHeader from '../ui/PageHeader.vue';
import Button from '../ui/Button.vue';
import DataTable from '../ui/DataTable.vue';
import DropdownMenu from '../ui/DropdownMenu.vue';
import ConfirmDialog from '../ui/ConfirmDialog.vue';
import { fetchUserRoles, createUserRole } from '../../api/auth';
import api from '../../api/api';

const rows = ref([]);
const columns = [ { key: 'role_name', label: 'Role' }, { key: 'created_at', label: 'Created' } ];
const openCreate = ref(false);
const openEdit = ref(false);
const editingId = ref(null);
const roleName = ref('');

onMounted(load);
async function load(){ rows.value = await fetchUserRoles(); }

function closeForm(){ openCreate.value=false; openEdit.value=false; roleName.value=''; editingId.value=null; }
function startEdit(row){ openEdit.value=true; roleName.value=row.role_name; editingId.value=row.id; }
async function saveRole(){
  if(openEdit.value){
    await api.put(`/roles/${editingId.value}`, { role_name: roleName.value });
  } else {
    await createUserRole({ role_name: roleName.value });
  }
  closeForm();
  await load();
}

const confirmOpen = ref(false); const selected = ref(null);
function askDelete(row){ selected.value=row; confirmOpen.value=true; }
async function confirmDelete(){ confirmOpen.value=false; if(!selected.value) return; await api.delete(`/roles/${selected.value.id}`); await load(); }
</script>
