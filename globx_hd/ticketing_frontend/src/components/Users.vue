<template>
  <div class="p-8">
    <PageHeader title="Users">
      <template #actions>
        <Button @click="openCreate=true"><span class="mr-2">+</span>Add User</Button>
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

    <ConfirmDialog :open="confirmOpen" title="Delete user" :message="`Delete user “${selected?.username}”?`" @cancel="confirmOpen=false" @confirm="confirmDelete" />

    <Modal :open="openCreate || openEdit" :title="openEdit ? 'Edit User' : 'Add User'" @close="closeForm">
      <FormLayout @submit="save">
        <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
          <p class="text-sm text-red-700">{{ errorMessage }}</p>
        </div>
        <FormField label="Employee ID"><input v-model="form.employee_id" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Username"><input v-model="form.username" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField v-if="!openEdit" label="Password">
          <input 
            type="password" 
            v-model="form.password" 
            minlength="6"
            required
            placeholder="Minimum 6 characters"
            class="w-full border border-blue-200 rounded px-3 py-2" 
          />
        </FormField>
        <FormField label="First Name"><input v-model="form.first_name" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Last Name"><input v-model="form.last_name" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Email"><input type="email" v-model="form.email" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Phone"><input v-model="form.phone" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Designation">
          <select v-model.number="form.designation_id" class="w-full border border-blue-200 rounded px-3 py-2">
            <option disabled :value="0">Select designation</option>
            <option v-for="d in userDesignations" :key="d.id" :value="d.id">{{ d.designation_name }}</option>
          </select>
        </FormField>
        <FormField label="Role">
          <select v-model.number="form.role_id" class="w-full border border-blue-200 rounded px-3 py-2">
            <option disabled :value="0">Select role</option>
            <option v-for="r in roles" :key="r.id" :value="r.id">{{ r.role_name }}</option>
          </select>
        </FormField>
        <!-- Password Reset Section (Only in Edit mode) -->
        <div v-if="openEdit" class="md:col-span-2 mt-4 pt-4 border-t border-gray-100">
          <div v-if="!showPasswordReset">
            <Button type="button" variant="secondary" @click="showPasswordReset = true" class="w-full justify-center">
              Reset User Password
            </Button>
          </div>
          <div v-else class="rounded-lg border border-amber-200 bg-amber-50 p-4">
            <p class="mb-3 text-sm font-medium text-neutral-dark">Set a new password for this user.</p>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2 items-end">
              <FormField label="New Password">
                <input 
                  v-model="passwordReset.newPassword" 
                  type="password" 
                  minlength="6" 
                  placeholder="Minimum 6 characters" 
                  class="w-full border border-amber-300 rounded px-3 py-2 bg-white" 
                />
              </FormField>
              <div class="flex space-x-2 h-10">
                <Button 
                  type="button" 
                  @click="performPasswordReset" 
                  :disabled="passwordReset.newPassword.length < 6"
                  class="flex-1 justify-center"
                >
                  Apply
                </Button>
                <Button 
                  type="button" 
                  variant="secondary" 
                  @click="showPasswordReset = false; passwordReset.newPassword = '';"
                  class="flex-1 justify-center"
                >
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        </div>
        <template #actions>
          <Button variant="secondary" type="button" @click="closeForm">Cancel</Button>
          <Button type="submit">Save</Button>
        </template>
      </FormLayout>
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
import { fetchUsers, createUser, updateUser, deleteUser, fetchUserDesignations, fetchUserRoles, resetManagedUserPassword } from '../api/auth';

const rows = ref([]);
const roles = ref([]); const userDesignations = ref([]);
const columns = [
  { key: 'employee_id', label: 'Employee ID' },
  { key: 'username', label: 'Username' },
  { key: 'first_name', label: 'First Name' },
  { key: 'last_name', label: 'Last Name' },
  { key: 'email', label: 'Email' },
  { key: 'phone', label: 'Phone' },
];

const openCreate = ref(false); const openEdit = ref(false); const editingId = ref(null);
const form = reactive({ employee_id: '', username: '', password: '', first_name: '', last_name: '', email: '', phone: '', designation_id: 0, role_id: 0 });
const errorMessage = ref('');

const showPasswordReset = ref(false);
const passwordReset = reactive({ newPassword: '' });

onMounted(async () => {
  roles.value = await fetchUserRoles();
  userDesignations.value = await fetchUserDesignations();
  await load();
});

async function load(){ rows.value = await fetchUsers(); }
function closeForm(){ 
  openCreate.value=false; 
  openEdit.value=false; 
  editingId.value=null; 
  errorMessage.value='';
  showPasswordReset.value = false;
  passwordReset.newPassword = '';
  Object.assign(form, { employee_id:'', username:'', password:'', first_name:'', last_name:'', email:'', phone:'', designation_id:0, role_id:0 }); 
}
function startEdit(row){ 
  openEdit.value=true; 
  editingId.value=row.id; 
  showPasswordReset.value = false;
  passwordReset.newPassword = '';
  Object.assign(form, { employee_id: row.employee_id, username: row.username, password:'', first_name: row.first_name, last_name: row.last_name, email: row.email, phone: row.phone, designation_id: row.designation_id, role_id: row.role_id }); 
}

async function performPasswordReset() {
  try {
    errorMessage.value = '';
    await resetManagedUserPassword(editingId.value, { password: passwordReset.newPassword });
    passwordReset.newPassword = '';
    showPasswordReset.value = false;
    alert('Password reset successfully!');
  } catch (error) {
    console.error('Error resetting password:', error);
    if (error.response?.data?.error) {
      errorMessage.value = error.response.data.error;
    } else {
      errorMessage.value = 'Failed to reset password. Please try again.';
    }
  }
}
async function save(){
  try {
    errorMessage.value = ''; // Clear previous errors
    
    if(openEdit.value){
      const payload = { employee_id: form.employee_id, username: form.username, first_name: form.first_name, last_name: form.last_name, email: form.email, phone: form.phone, designation_id: form.designation_id, role_id: form.role_id };
      if(form.password) payload.password = form.password;
      await updateUser(editingId.value, payload);
    } else {
      // Validate password length before sending
      if (!form.password || form.password.length < 6) {
        errorMessage.value = 'Password must be at least 6 characters long';
        return;
      }
      
      await createUser({ employee_id: form.employee_id, username: form.username, password: form.password, first_name: form.first_name, last_name: form.last_name, email: form.email, phone: form.phone, designation_id: form.designation_id, role_id: form.role_id });
    }
    closeForm(); 
    await load();
  } catch (error) {
    console.error('Error saving user:', error);
    // Extract error message from response
    if (error.response?.data?.error) {
      errorMessage.value = error.response.data.error;
    } else {
      errorMessage.value = 'An error occurred while saving the user. Please try again.';
    }
  }
}

const confirmOpen = ref(false); const selected = ref(null);
function askDelete(row){ selected.value=row; confirmOpen.value=true; }
async function confirmDelete(){ confirmOpen.value=false; if(!selected.value) return; await deleteUser(selected.value.id); await load(); }
</script>
