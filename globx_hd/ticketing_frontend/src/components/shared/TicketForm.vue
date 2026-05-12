<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-4xl mx-auto mt-8">
    <h2 class="text-2xl font-bold text-blue-800 mb-6">
      {{ isEditMode ? 'Edit Ticket Properties' : 'Raise Ticket on Behalf of Customer' }}
    </h2>
    
    
    <form @submit.prevent="onSubmit" class="space-y-6">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        
        <!-- Account Selection -->
        <div class="md:col-span-2">
          <FormField label="Account" required>
            <FuzzySearchDropdown
              v-model="form.selectedAccount"
              :items="accounts"
              :search-keys="['account_name', 'account_owner', 'customer_code']"
              :display-key="'account_name'"
              placeholder="Search accounts by name, owner, or customer code..."
              :allow-create="true"
              @select="selectAccount"
              @create-new="openAccountCreateModal"
            >
              <template #item="{ item }">
                <div class="flex items-center justify-between">
                  <div>
                    <div class="font-medium text-blue-900">{{ item.account_name }}</div>
                    <div class="text-sm text-gray-600">{{ item.account_owner || 'No owner' }}</div>
                    <div class="text-sm text-gray-500">Code: {{ item.customer_code }}</div>
                  </div>
                </div>
              </template>
            </FuzzySearchDropdown>
          </FormField>
        </div>

        <!-- Account Owner (Auto-filled) - Only show if account is selected -->
        <FormField v-if="form.selectedAccount" label="Account Owner">
          <input 
            v-model="form.accountOwner" 
            type="text" 
            readonly
            class="w-full px-3 py-2 border rounded bg-gray-50 text-gray-700"
            placeholder="Will be filled automatically"
          />
        </FormField>

        <!-- Contact Selection - Only show if account is selected -->
        <FormField v-if="form.selectedAccount" label="Contact" required>
          <FuzzySearchDropdown
            v-model="form.selectedContact"
            :items="filteredContactsWithFullName"
            :search-keys="['full_name', 'first_name', 'last_name', 'email', 'mobile']"
            :display-key="'full_name'"
            placeholder="Search contacts within this account..."
            :allow-create="true"
            @select="selectContact"
            @create-new="openContactCreateModal"
          >
            <template #item="{ item }">
              <div class="flex items-center justify-between">
                <div>
                  <div class="font-medium text-blue-900">{{ item.first_name }} {{ item.last_name }}</div>
                  <div class="text-sm text-gray-600">{{ item.email }}</div>
                  <div class="text-sm text-gray-500">{{ item.mobile || 'No phone' }}</div>
                </div>
              </div>
            </template>
          </FuzzySearchDropdown>
        </FormField>

        <!-- Email (Auto-filled) -->
        <FormField label="Email">
          <input 
            v-model="form.email" 
            type="email" 
            readonly
            class="w-full px-3 py-2 border rounded bg-gray-50 text-gray-700"
            placeholder="Will be filled automatically"
          />
        </FormField>

        <!-- Phone (Auto-filled) -->
        <FormField label="Phone">
          <input 
            v-model="form.phone" 
            type="text" 
            readonly
            class="w-full px-3 py-2 border rounded bg-gray-50 text-gray-700"
            placeholder="Will be filled automatically"
          />
        </FormField>

        <!-- Product Name -->
        <FormField label="Product Name" required>
          <FuzzySearchDropdown
            v-model="form.selectedProduct"
            :items="products"
            :search-keys="['product_name']"
            :display-key="'product_name'"
            placeholder="Search products by name..."
            :allow-create="true"
            @select="selectProduct"
            @create-new="openProductCreateModal"
          >
            <template #item="{ item }">
              <div class="font-medium text-blue-900">{{ item.product_name }}</div>
            </template>
          </FuzzySearchDropdown>
        </FormField>

        <!-- Subject -->
        <FormField label="Subject" required>
          <input 
            v-model="form.subject" 
            type="text" 
            class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
            placeholder="Enter ticket subject..."
          />
        </FormField>

        <!-- Status (only in create mode) -->
        <FormField v-if="!isEditMode" label="Status" required>
          <select 
            v-model="form.status" 
            class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
          >
            <option value="" disabled>Select Status</option>
            <option value="OPEN">OPEN</option>
            <option value="IN PROGRESS">IN PROGRESS</option>
            <option value="RESOLVED">RESOLVED</option>
            <option value="CLOSED">CLOSED</option>
          </select>
        </FormField>

        <!-- Ticket Owner (hidden for engineers) -->
        <FormField v-if="!hideTicketOwner" label="Ticket Owner">
          <select 
            v-model="form.assignedEngineer" 
            class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
          >
            <option value="">No Assignment</option>
            <option v-for="user in engineerUsers" :key="user.id" :value="user.id">
              {{ formatUserName(user) }}
            </option>
          </select>
        </FormField>

        <!-- Priority (only in create mode) -->
        <FormField v-if="!isEditMode" label="Priority" required>
          <select 
            v-model="form.priority" 
            class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
          >
            <option value="" disabled>Select Priority</option>
            <option value="High">High</option>
            <option value="Medium">Medium</option>
            <option value="Low">Low</option>
          </select>
        </FormField>

        <!-- Channel -->
        <FormField label="Channel" required>
          <select 
            v-model="form.channel" 
            class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
          >
            <option value="" disabled>Select Channel</option>
            <option value="Phone">Phone</option>
            <option value="Mail">Mail</option>
          </select>
        </FormField>

        <!-- Description -->
        <div class="md:col-span-2">
          <FormField label="Description" required>
            <textarea 
              v-model="form.description" 
              rows="4" 
              class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
              placeholder="Describe the issue in detail..."
            ></textarea>
          </FormField>
        </div>

        <!-- File Upload Section (only for create mode) -->
        <div v-if="!isEditMode" class="md:col-span-2">
          <FormField label="Attach Files">
            <FileUpload 
              ref="fileUpload"
              :ticket-id="createdTicketId"
              :show-upload-button="false"
              @files-selected="onFilesSelected"
              @upload-success="onUploadSuccess"
              @upload-error="onUploadError"
            />
          </FormField>
        </div>
      </div>

      <!-- Submit Button -->
      <div class="flex justify-end space-x-3">
        <button
          v-if="isEditMode"
          type="button"
          @click="$emit('cancel')"
          class="px-6 py-2 border border-gray-300 rounded text-gray-700 hover:bg-gray-50 transition"
        >
          Cancel
        </button>
        <Button 
          type="submit" 
          class="bg-blue-700 text-white px-6 py-2 rounded hover:bg-blue-800 transition"
          :disabled="submitting"
        >
          {{ submitting ? 'Saving...' : (isEditMode ? 'Save Changes' : 'Create Ticket') }}
        </Button>
      </div>
    </form>

    <!-- Success Message -->
    <div v-if="showSuccess" class="mt-4 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
      {{ isEditMode ? 'Ticket updated successfully!' : 'Ticket created successfully!' }}
    </div>

    <!-- Error Modal -->
    <Modal v-if="error" @close="error = ''">
      <div class="text-red-700 font-semibold">{{ error }}</div>
    </Modal>

    <!-- Account Create Modal -->
    <AccountCreateModal 
      :open="showAccountCreateModal"
      :initial-name="accountCreateInitialName"
      @close="showAccountCreateModal = false"
      @created="onAccountCreated"
    />

    <!-- Contact Create Modal -->
    <ContactCreateModal 
      :open="showContactCreateModal"
      :initial-name="contactCreateInitialName"
      :pre-selected-account="form.selectedAccount"
      :accounts-list="accounts"
      @close="showContactCreateModal = false"
      @created="onContactCreated"
    />

    <!-- Product Create Modal -->
    <ProductCreateModal 
      :open="showProductCreateModal"
      :initial-name="productCreateInitialName"
      @close="showProductCreateModal = false"
      @created="onProductCreated"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed, nextTick } from 'vue';
import { useRoute } from 'vue-router';
import { fetchContacts } from '../../api/contacts';
import { fetchAccounts } from '../../api/auth';
import { fetchProducts, managerCreateTicket, engineerCreateTicket, fetchUsers, updateTicket } from '../../api/tickets';
import { formatUserName } from '../../utils/user';
import { useAuthStore } from '../../stores/auth';
import FormField from '../ui/FormField.vue';
import Button from '../ui/Button.vue';
import Modal from '../ui/Modal.vue';
import FileUpload from './FileUpload.vue';
import FuzzySearchDropdown from '../ui/FuzzySearchDropdown.vue';
import ContactCreateModal from './ContactCreateModal.vue';
import ProductCreateModal from './ProductCreateModal.vue';
import AccountCreateModal from './AccountCreateModal.vue';

const props = defineProps({
  isEditMode: {
    type: Boolean,
    default: false
  },
  hideTicketOwner: {
    type: Boolean,
    default: false
  },
  ticketData: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['success', 'cancel']);

const route = useRoute();
const auth = useAuthStore();

const form = ref({
  selectedAccount: null,
  selectedContact: null,
  selectedProduct: null,
  accountName: '',
  accountOwner: '',
  email: '',
  phone: '',
  productId: '',
  subject: '',
  status: '',
  assignedEngineer: '',
  priority: '',
  channel: '',
  description: ''
});

const accounts = ref([]);
const contacts = ref([]);
const products = ref([]);
const users = ref([]);
const showSuccess = ref(false);
const error = ref('');
const submitting = ref(false);

// Modal states
const showAccountCreateModal = ref(false);
const showContactCreateModal = ref(false);
const showProductCreateModal = ref(false);
const accountCreateInitialName = ref('');
const contactCreateInitialName = ref('');
const productCreateInitialName = ref('');

// Computed property to filter users to show only Engineers (role_id 3 and designation_id 3)
const engineerUsers = computed(() => {
  const engineers = users.value.filter(user => 
    user.role_id === 3 && user.designation_id === 3
  );
  console.log('🔧 DEBUG: Filtered engineers:', engineers.length, 'out of', users.value.length, 'total users');
  return engineers;
});

// Computed property to add full_name to contacts for fuzzy search display
const contactsWithFullName = computed(() => {
  return contacts.value.map(contact => ({
    ...contact,
    full_name: `${contact.first_name} ${contact.last_name}`.trim()
  }));
});

// Computed property to filter contacts by selected account
const filteredContactsWithFullName = computed(() => {
  if (!form.value.selectedAccount) {
    return [];
  }
  
  return contactsWithFullName.value.filter(contact => 
    contact.account_id === form.value.selectedAccount.id
  );
});

// File upload variables
const fileUpload = ref(null);
const createdTicketId = ref(null);
const selectedFiles = ref([]);

// Load initial data
onMounted(async () => {
  try {
    const [accountsData, contactsData, productsData, usersData] = await Promise.all([
      fetchAccounts(),
      fetchContacts(),
      fetchProducts(),
      fetchUsers()
    ]);
    
    console.log('Raw API responses:', { accountsData, contactsData, productsData, usersData });
    
    // Handle different API response structures
    accounts.value = Array.isArray(accountsData) ? accountsData : (accountsData?.accounts || []);
    contacts.value = Array.isArray(contactsData) ? contactsData : (contactsData?.contacts || []);
    products.value = Array.isArray(productsData) ? productsData : (productsData?.products || []);
    users.value = Array.isArray(usersData) ? usersData : (usersData?.users || []);
    
    console.log('Processed data:', { 
      accounts: accounts.value.length,
      contacts: contacts.value.length, 
      products: products.value.length, 
      users: users.value.length 
    });

    // If in edit mode, populate form with existing data
    if (props.isEditMode && props.ticketData) {
      populateFormWithTicketData();
    } else if (!props.isEditMode && Object.keys(route.query).length > 0) {
      // Handle pre-fill from query params (e.g. from Dumped Queries)
      populateFormFromQuery();
    }
  } catch (err) {
    console.error('Failed to load form data:', err);
    error.value = 'Failed to load form data. Please refresh the page.';
  }
});

// Populate form from URL query parameters
const populateFormFromQuery = async () => {
  const query = route.query;
  console.log('🔄 Populating form from query params:', query);

  if (query.subject) {
    form.value.subject = String(query.subject);
  }

  if (query.description) {
    form.value.description = String(query.description);
  }

  // Try to match contact by email
  if (query.email) {
    const emailToFind = String(query.email).toLowerCase();
    const foundContact = contactsWithFullName.value.find(c => 
      c.email && c.email.toLowerCase() === emailToFind
    );

    if (foundContact) {
      console.log('✅ Found matching contact by email:', foundContact);
      
      // If contact belongs to an account, select the account first
      if (foundContact.account_id) {
        const foundAccount = accounts.value.find(a => a.id === foundContact.account_id);
        if (foundAccount) {
          console.log('✅ Found associated account:', foundAccount);
          selectAccount(foundAccount);
          
          // Wait for reactivity (though selectAccount is synchronous in state update,
          // it might trigger watchers if any. Here we just set values directly after).
          await nextTick();
        }
      }
      
      // Select the contact (overriding any clearing done by selectAccount)
      selectContact(foundContact);
    } else {
      console.warn('⚠️ No contact found for email:', emailToFind);
      // Optional: Pre-fill email/phone fields even if not found? 
      // The current form logic relies on selectedContact for email/phone (they are readonly).
      // If we wanted to allow creating a new contact easily, we might want to trigger the create modal?
      // For now, we just leave it blank if not found, as per standard behavior.
    }
  }
};

// Populate form when in edit mode
const populateFormWithTicketData = async () => {
  const ticket = props.ticketData;
  console.log('🔄 Populating form with ticket data:', ticket);
  console.log('Available accounts:', accounts.value.length);
  console.log('Available contacts:', contacts.value.length);
  console.log('Available products:', products.value.length);
  console.log('Available users:', users.value.length);
  
  // Use preloaded Contact from ticket data (backend sends it with Preload)
  if (ticket.Contact) {
    const contact = ticket.Contact;
    console.log('✅ Using preloaded contact:', contact);
    
    // Add full_name for fuzzy search display
    const contactWithFullName = {
      ...contact,
      full_name: `${contact.first_name} ${contact.last_name}`.trim()
    };
    
    // Set the contact
    form.value.selectedContact = contactWithFullName;
    form.value.email = contact.email || '';
    form.value.phone = contact.mobile || '';
    
    // Use preloaded Account from Contact or ticket
    const account = ticket.Contact.Account || ticket.Account;
    if (account) {
      console.log('✅ Using preloaded account:', account);
      form.value.selectedAccount = account;
      form.value.accountOwner = account.account_owner || '';
    } else {
      console.log('⚠️ No account found in preloaded data');
    }
  } else {
    console.warn('❌ Contact not found in ticket data');
  }
  
  // Use preloaded Product from ticket data
  if (ticket.Product) {
    console.log('✅ Using preloaded product:', ticket.Product);
    form.value.selectedProduct = ticket.Product;
    form.value.productId = ticket.Product.id ? String(ticket.Product.id) : '';
  } else {
    console.warn('❌ Product not found in ticket data');
  }
  
  // Set other form fields
  form.value.subject = ticket.subject || '';
  form.value.status = ticket.ticket_status || '';
  form.value.assignedEngineer = ticket.assigned_engineer ? String(ticket.assigned_engineer) : '';
  form.value.priority = ticket.priority || '';
  form.value.channel = ticket.channel || '';
  form.value.description = ticket.ticket_details || '';
  
  console.log('✅ Form populated successfully:', form.value);
  console.log('📋 Account:', form.value.selectedAccount);
  console.log('📋 Contact:', form.value.selectedContact);
  console.log('📋 Product:', form.value.selectedProduct);
  
  // Force reactivity update with nextTick to ensure FuzzySearchDropdown watchers trigger
  await nextTick();
  
  // Additional tick to ensure DOM is fully updated
  await nextTick();
};

// Watch for changes in ticketData prop
watch(() => props.ticketData, (newData) => {
  if (props.isEditMode && newData && accounts.value.length > 0 && contacts.value.length > 0 && products.value.length > 0 && users.value.length > 0) {
    populateFormWithTicketData();
  }
});

// Watch for data loading completion in edit mode
watch([() => accounts.value.length, () => contacts.value.length, () => products.value.length, () => users.value.length], () => {
  if (props.isEditMode && props.ticketData && accounts.value.length > 0 && contacts.value.length > 0 && products.value.length > 0 && users.value.length > 0) {
    populateFormWithTicketData();
  }
});

// Account selection
const selectAccount = (account) => {
  form.value.selectedAccount = account;
  form.value.accountOwner = account.account_owner || '';
  
  // Clear contact selection when account changes
  form.value.selectedContact = null;
  form.value.email = '';
  form.value.phone = '';
};

// Contact selection
const selectContact = (contact) => {
  form.value.selectedContact = contact;
  form.value.email = contact.email || '';
  form.value.phone = contact.mobile || '';
};

// Product selection
const selectProduct = (product) => {
  form.value.selectedProduct = product;
  form.value.productId = product.id;
};

// Form submission
const onSubmit = async () => {
  // Validate account and contact selection (both create and edit modes)
  if (!form.value.selectedAccount) {
    error.value = 'Please select an account';
    return;
  }
  
  if (!form.value.selectedContact) {
    error.value = 'Please select a contact';
    return;
  }
  
  if (!form.value.selectedProduct) {
    error.value = 'Please select a product';
    return;
  }

  // Validate required fields for backend
  if (!form.value.subject?.trim()) {
    error.value = 'Please enter a subject';
    return;
  }

  if (!form.value.description?.trim()) {
    error.value = 'Please enter ticket details';
    return;
  }

  // Priority is only required in create mode
  if (!props.isEditMode && !form.value.priority?.trim()) {
    error.value = 'Please select a priority';
    return;
  }

  submitting.value = true;
  error.value = '';

  try {
    if (props.isEditMode) {
      // Update existing ticket - include contact_id for contact changes
      const updatePayload = {
        contact_id: form.value.selectedContact.id,
        product_id: form.value.selectedProduct?.id || (form.value.productId ? parseInt(form.value.productId) : null),
        subject: form.value.subject,
        assigned_engineer: form.value.assignedEngineer ? parseInt(form.value.assignedEngineer) : null,
        channel: form.value.channel,
        ticket_details: form.value.description
      };
      console.log('Update payload:', updatePayload);
      await updateTicket(props.ticketData.id, updatePayload);
      showSuccess.value = true;
      setTimeout(() => {
        emit('success');
      }, 1500);
    } else {
      // Create new ticket - ensure all required fields are present
      const createPayload = {
        contact_id: form.value.selectedContact.id,
        product_id: form.value.selectedProduct?.id || parseInt(form.value.productId),
        subject: form.value.subject || 'No Subject',
        ticket_status: form.value.status || 'OPEN',
        assigned_engineer: form.value.assignedEngineer ? parseInt(form.value.assignedEngineer) : null,
        priority: form.value.priority || 'Medium',
        channel: form.value.channel || 'Manager',
        ticket_details: form.value.description || 'No details provided'
      };
      console.log('Create payload:', createPayload);
      
      // Use appropriate API based on user role
      const response = auth.userType === 'engineer' 
        ? await engineerCreateTicket(createPayload)
        : await managerCreateTicket(createPayload);
      const ticketId = response.ticket.ticket_id;
      createdTicketId.value = ticketId;
      
      // Upload files if any are selected
      if (selectedFiles.value.length > 0) {
        console.log('🔥 DEBUG: Uploading files for manager ticket...');
        await fileUpload.value.uploadFiles(ticketId);
      }
      
      showSuccess.value = true;
      emit('success', response.ticket); // Emit success with ticket data
      
      // Reset form after successful creation
      setTimeout(() => {
        form.value = {
          selectedAccount: null,
          selectedContact: null,
          selectedProduct: null,
          accountName: '',
          accountOwner: '',
          email: '',
          phone: '',
          productId: '',
          subject: '',
          status: '',
          assignedEngineer: '',
          priority: '',
          channel: '',
          description: ''
        };
        createdTicketId.value = null;
        selectedFiles.value = [];
        if (fileUpload.value) {
          fileUpload.value.clearFiles();
        }
        showSuccess.value = false;
      }, 2000);
    }
  } catch (err) {
    console.error('Failed to save ticket:', err);
    console.error('Error response:', err.response?.data);
    console.error('Error status:', err.response?.status);
    
    let errorMessage = 'Unknown error occurred';
    if (err.response?.data?.error) {
      errorMessage = err.response.data.error;
    } else if (err.response?.data?.message) {
      errorMessage = err.response.data.message;
    } else if (err.message) {
      errorMessage = err.message;
    }
    
    error.value = `Failed to ${props.isEditMode ? 'update' : 'create'} ticket: ${errorMessage}`;
  } finally {
    submitting.value = false;
  }
};

// File upload event handlers
const onFilesSelected = (files) => {
  selectedFiles.value = files;
};

const onUploadSuccess = (response) => {
  console.log('Files uploaded successfully:', response);
  if (response.errors && response.errors.length > 0) {
    error.value = 'Some files failed to upload: ' + response.errors.join(', ');
  }
};

const onUploadError = (errorMessage) => {
  error.value = errorMessage;
};

// Modal handlers
const openAccountCreateModal = (searchQuery) => {
  if (showAccountCreateModal.value) {
    return; // Prevent opening modal if already open
  }
  accountCreateInitialName.value = searchQuery;
  showAccountCreateModal.value = true;
};

const openContactCreateModal = (searchQuery) => {
  contactCreateInitialName.value = searchQuery;
  showContactCreateModal.value = true;
};

const openProductCreateModal = (searchQuery) => {
  productCreateInitialName.value = searchQuery;
  showProductCreateModal.value = true;
};

const onAccountCreated = async (newAccount) => {
  try {
    // Refresh accounts list
    const accountsData = await fetchAccounts();
    accounts.value = Array.isArray(accountsData) ? accountsData : (accountsData?.accounts || []);
    
    // Find the account in the refreshed list to ensure it has all properties
    const foundAccount = accounts.value.find(acc => acc.id === newAccount.id || acc.account_name === newAccount.account_name);
    
    // Select the newly created account
    if (foundAccount) {
      selectAccount(foundAccount);
    } else {
      selectAccount(newAccount);
    }
  } catch (error) {
    console.error('Failed to refresh accounts after creation:', error);
    // Still select the account even if refresh fails
    selectAccount(newAccount);
  }
};

const onContactCreated = async (newContact) => {
  // Refresh contacts list
  const contactsData = await fetchContacts();
  contacts.value = Array.isArray(contactsData) ? contactsData : (contactsData?.contacts || []);
  
  // Select the newly created contact
  const contactWithFullName = {
    ...newContact,
    full_name: `${newContact.first_name} ${newContact.last_name}`.trim()
  };
  selectContact(contactWithFullName);
};

const onProductCreated = async (newProduct) => {
  // Refresh products list
  const productsData = await fetchProducts();
  products.value = Array.isArray(productsData) ? productsData : (productsData?.products || []);
  
  // Select the newly created product
  selectProduct(newProduct);
};

</script>
