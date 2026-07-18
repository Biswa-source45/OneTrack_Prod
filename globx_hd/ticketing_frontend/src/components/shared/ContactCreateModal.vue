<template>
  <Modal :open="open" title="Add Contact" @close="$emit('close')">
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
      <FormField label="First Name">
        <input v-model="form.first_name" class="w-full border border-blue-200 rounded px-3 py-2" />
      </FormField>
      <FormField label="Last Name">
        <input v-model="form.last_name" class="w-full border border-blue-200 rounded px-3 py-2" />
      </FormField>
      <FormField label="Designation">
        <select v-model.number="form.designation_id" class="w-full border border-blue-200 rounded px-3 py-2">
          <option disabled :value="0">Select designation</option>
          <option v-for="d in designations" :key="d.id" :value="d.id">{{ d.designation_name }}</option>
        </select>
      </FormField>
      <FormField label="Department">
        <input v-model="form.department" class="w-full border border-blue-200 rounded px-3 py-2" />
      </FormField>
      <FormField label="Location">
        <input v-model="form.location" class="w-full border border-blue-200 rounded px-3 py-2" />
      </FormField>
      <FormField label="Email (Optional)">
        <input type="email" v-model="form.email" placeholder="Enter email address" class="w-full border border-blue-200 rounded px-3 py-2" />
      </FormField>
      <FormField label="Mobile">
        <input v-model="form.mobile" class="w-full border border-blue-200 rounded px-3 py-2" />
      </FormField>
      <FormField label="Password (Optional)">
        <input type="password" v-model="form.password" placeholder="Enter password (min 6 characters)" class="w-full border border-blue-200 rounded px-3 py-2" />
      </FormField>
      <template #actions>
        <Button variant="secondary" type="button" @click="$emit('close')">Cancel</Button>
        <Button type="submit" :disabled="submitting">{{ submitting ? 'Saving...' : 'Save' }}</Button>
      </template>
    </FormLayout>
  </Modal>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue';
import Modal from '../ui/Modal.vue';
import FormLayout from '../ui/FormLayout.vue';
import FormField from '../ui/FormField.vue';
import Button from '../ui/Button.vue';
import { createContact, fetchContactDesignations, fetchAccounts } from '../../api/contacts';

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  initialName: {
    type: String,
    default: ''
  },
  preSelectedAccount: {
    type: Object,
    default: null
  },
  accountsList: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['close', 'created']);

const accounts = ref([]);
const designations = ref([]);
const submitting = ref(false);
const errorMessage = ref('');

const form = reactive({
  contact_type: '',
  account_id: 0,
  designation_id: 0,
  department: '',
  location: '',
  first_name: '',
  last_name: '',
  email: '',
  mobile: '',
  password: ''
});

// Load required data
onMounted(async () => {
  try {
    // Load designations (accounts will be handled by watcher)
    const designationsData = await fetchContactDesignations();
    designations.value = designationsData;
  } catch (error) {
    console.error('Failed to load form data:', error);
    errorMessage.value = 'Failed to load form data. Please try again.';
  }
});

// Watch for accountsList prop changes and update local accounts
watch(() => props.accountsList, (newAccountsList) => {
  if (newAccountsList && newAccountsList.length > 0) {
    console.log('🔄 ContactCreateModal: Received accounts list:', newAccountsList.length, 'accounts');
    accounts.value = newAccountsList;
  }
}, { immediate: true });

// Fallback: Load accounts if not provided via prop
watch(() => props.open, async (isOpen) => {
  if (isOpen && accounts.value.length === 0) {
    try {
      const accountsData = await fetchAccounts();
      accounts.value = Array.isArray(accountsData) ? accountsData : (accountsData?.accounts || []);
    } catch (error) {
      console.error('Failed to load accounts:', error);
    }
  }
});

// Parse initial name and set pre-selected account when modal opens
watch(() => props.open, (isOpen) => {
  if (isOpen) {
    // Parse initial name
    if (props.initialName) {
      const nameParts = props.initialName.trim().split(' ');
      if (nameParts.length >= 2) {
        form.first_name = nameParts[0];
        form.last_name = nameParts.slice(1).join(' ');
      } else {
        form.first_name = props.initialName;
      }
    }
    
    // Set pre-selected account
    if (props.preSelectedAccount) {
      form.account_id = props.preSelectedAccount.id;
      form.contact_type = 'Private'; // Default to Private when account is pre-selected
    }
  }
});

// Watch for preSelectedAccount changes (in case it updates after modal is open)
watch(() => props.preSelectedAccount, (newAccount) => {
  if (newAccount && props.open) {
    console.log('🎯 ContactCreateModal: Pre-selecting account:', newAccount.account_name);
    form.account_id = newAccount.id;
    if (!form.contact_type) {
      form.contact_type = 'Private'; // Default to Private when account is pre-selected
    }
  }
}, { immediate: true });

// Reset form when modal closes
watch(() => props.open, (isOpen) => {
  if (!isOpen) {
    Object.assign(form, {
      contact_type: '',
      account_id: 0,
      designation_id: 0,
      department: '',
      location: '',
      first_name: '',
      last_name: '',
      email: '',
      mobile: '',
      password: ''
    });
    errorMessage.value = '';
  }
});

const save = async () => {
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
  if (form.password && form.password.length < 6) {
    errorMessage.value = 'Password must be at least 6 characters long.';
    return;
  }

  submitting.value = true;
  errorMessage.value = '';
  
  try {
    const payload = {
      contact_type: form.contact_type,
      designation_id: form.designation_id,
      department: form.department,
      location: form.location,
      first_name: form.first_name,
      last_name: form.last_name,
      email: form.email,
      mobile: form.mobile,
      password: form.password
    };
    
    // Only include account_id for Govt. and Private contacts
    if (form.contact_type === 'Govt.' || form.contact_type === 'Private') {
      payload.account_id = form.account_id;
    }
    
    const newContact = await createContact(payload);
    emit('created', newContact);
    emit('close');
  } catch (error) {
    console.error('Failed to create contact:', error);
    errorMessage.value = error.response?.data?.error || 'Failed to create contact. Please try again.';
  } finally {
    submitting.value = false;
  }
};
</script>
