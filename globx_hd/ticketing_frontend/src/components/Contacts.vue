<template>
  <div class="p-8">
    <PageHeader title="Contacts">
      <template #actions>
        <Button @click="openCreate=true"><span class="mr-2">+</span>Add Contact</Button>
      </template>
    </PageHeader>

    <DataTable :columns="columns" :rows="rows" rowKey="id">
      <template #cell:account.account_name="{ row }">
        {{ row.account && row.account.account_name ? row.account.account_name : '' }}
      </template>
      <template #row-actions="{ row }">
        <DropdownMenu>
          <template #button="{ toggle }">
            <button @click="toggle" class="inline-flex items-center justify-center w-8 h-8 rounded hover:bg-blue-100">&#8942;</button>
          </template>
          <div class="py-1">
            <button class="w-full text-left px-3 py-2 text-sm" @click="startEdit(row)">Edit</button>
            <button class="w-full text-left px-3 py-2 text-sm text-red-700 hover:bg-red-50" @click="askDelete(row)">Delete</button>
          </div>
        </DropdownMenu>
      </template>
    </DataTable>

    <ConfirmDialog :open="confirmOpen" title="Delete contact"
  :message="'Delete contact \'' + (selected?.first_name || '') + ' ' + (selected?.last_name || '') + '\'?'"
  @cancel="confirmOpen=false" @confirm="confirmDelete" />

    <Modal :open="openCreate || openEdit" :title="openEdit ? 'Edit Contact' : 'Add Contact'" @close="closeForm">
      <FormLayout @submit="save">
        <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
          <p class="text-sm text-red-700">{{ errorMessage }}</p>
        </div>
        <FormField label="Contact Type">
          <select v-model="form.contact_type" class="w-full border border-blue-200 rounded px-3 py-2">
            <option disabled value="">Select contact type</option>
            <option value="Govt.">Govt.</option>
            <option value="Private">Private</option>
            <option value="Individual">Individual</option>
          </select>
        </FormField>
        <FormField v-if="form.contact_type === 'Govt.' || form.contact_type === 'Private'" label="Account">
          <select v-model.number="form.account_id" class="w-full border border-blue-200 rounded px-3 py-2">
            <option disabled :value="0">Select account</option>
            <option v-for="a in accounts" :key="a.id" :value="a.id">{{ a.account_name }}</option>
          </select>
        </FormField>
        <FormField label="First Name"><input v-model="form.first_name" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Last Name"><input v-model="form.last_name" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Designation">
          <select v-model.number="form.designation_id" class="w-full border border-blue-200 rounded px-3 py-2">
            <option disabled :value="0">Select designation</option>
            <option v-for="d in designations" :key="d.id" :value="d.id">{{ d.designation_name }}</option>
          </select>
        </FormField>
        <FormField label="Department"><input v-model="form.department" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField v-if="!openEdit" label="Location"><input v-model="form.location" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Email (Optional)"><input type="email" v-model="form.email" placeholder="Enter email address" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField label="Mobile"><input v-model="form.mobile" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
        <FormField v-if="!openEdit" label="Password (Optional)"><input type="password" v-model="form.password" placeholder="Enter password (min 6 characters)" class="w-full border border-blue-200 rounded px-3 py-2" /></FormField>
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
import { fetchContacts, createContact, updateContact, deleteContact, fetchContactDesignations, fetchAccounts } from '../api/contacts';

const rows = ref([]);
const accounts = ref([]);
const designations = ref([]);
const columns = [
  { key: 'first_name', label: 'First Name' },
  { key: 'last_name', label: 'Last Name' },
  { key: 'contact_type', label: 'Contact Type' },
  { key: 'account.account_name', label: 'Account' },
  { key: 'department', label: 'Department' },
  { key: 'mobile', label: 'Mobile' },
];

const openCreate = ref(false); const openEdit = ref(false); const editingId = ref(null);
const errorMessage = ref('');
const form = reactive({ contact_type: '', account_id: 0, designation_id: 0, department: '', location: '', first_name: '', last_name: '', email: '', mobile: '', password: '' });

onMounted(async () => {
  accounts.value = await fetchAccounts();
  designations.value = await fetchContactDesignations();
  await load();
});

async function load(){
  const resp = await fetchContacts();
  console.log('Contacts API response:', resp);
  rows.value = resp;
}
function closeForm(){ openCreate.value=false; openEdit.value=false; editingId.value=null; errorMessage.value=''; Object.assign(form, { contact_type:'', account_id:0, designation_id:0, department:'', location:'', first_name:'', last_name:'', email:'', mobile:'', password:'' }); }
function startEdit(row){ openEdit.value=true; editingId.value=row.id; Object.assign(form, { contact_type: row.contact_type || '', account_id: row.account_id || 0, designation_id: row.designation_id, department: row.department || '', location: row.location || '', first_name: row.first_name, last_name: row.last_name, email: row.email, mobile: row.mobile, password: '' }); }
const submitting = ref(false);
async function save(){
  if (submitting.value) return;
  
  // Client-side validation
  if (!form.contact_type) {
    errorMessage.value = 'Contact Type is required.';
    return;
  }
  if ((form.contact_type === 'Govt.' || form.contact_type === 'Private') && !form.account_id) {
    errorMessage.value = 'Account is required for Govt. and Private contacts.';
    return;
  }
  if (!form.first_name || !form.first_name.trim()) {
    errorMessage.value = 'First Name is required.';
    return;
  }
  if (!form.designation_id) {
    errorMessage.value = 'Designation is required.';
    return;
  }
  if (!form.mobile || !form.mobile.trim()) {
    errorMessage.value = 'Mobile number is required.';
    return;
  }
  if (form.email && form.email.trim()) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(form.email.trim())) {
      errorMessage.value = 'Please enter a valid email address.';
      return;
    }
  }
  if (!openEdit.value && form.password && form.password.length < 6) {
    errorMessage.value = 'Password must be at least 6 characters long.';
    return;
  }
  if (openEdit.value && form.password && form.password.length > 0 && form.password.length < 6) {
    errorMessage.value = 'Password must be at least 6 characters long.';
    return;
  }

  submitting.value = true;
  errorMessage.value = '';
  try {
    if(openEdit.value){
      const payload = { contact_type: form.contact_type, designation_id: form.designation_id, department: form.department, location: form.location, first_name: form.first_name, last_name: form.last_name, email: form.email, mobile: form.mobile };
      // Only include account_id for Govt. and Private contacts
      if (form.contact_type === 'Govt.' || form.contact_type === 'Private') {
        payload.account_id = form.account_id;
      }
      if(form.password) payload.password = form.password;
      await updateContact(editingId.value, payload);
    } else {
      const payload = { contact_type: form.contact_type, designation_id: form.designation_id, department: form.department, location: form.location, first_name: form.first_name, last_name: form.last_name, email: form.email, mobile: form.mobile, password: form.password };
      // Only include account_id for Govt. and Private contacts
      if (form.contact_type === 'Govt.' || form.contact_type === 'Private') {
        payload.account_id = form.account_id;
      }
      await createContact(payload);
    }
    closeForm(); await load();
  } catch (error) {
    console.error('Failed to save contact:', error);
    errorMessage.value = error.response?.data?.error || 'Failed to save contact. Please try again.';
  } finally {
    submitting.value = false;
  }
}

const confirmOpen = ref(false); const selected = ref(null);
function askDelete(row){ selected.value=row; confirmOpen.value=true; }
async function confirmDelete(){ confirmOpen.value=false; if(!selected.value) return; await deleteContact(selected.value.id); await load(); }
</script>
