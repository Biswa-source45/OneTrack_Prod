<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-7xl mx-auto mt-8">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-blue-800">{{ pageTitle }}</h1>
      <button 
        v-if="statusFilter" 
        @click="clearFilter" 
        class="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
      >
        Clear Filter
      </button>
    </div>
    <div v-if="showSuccess" class="text-green-600 mt-2 text-center font-semibold">
  {{ successMessage }}
</div>
    <!-- Tickets Card Layout -->
    <div class="space-y-4">
      <div v-for="ticket in filteredTickets" :key="ticket.id" 
           class="bg-white border border-blue-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer"
           @click="openDetail(ticket)">
        
        <!-- Card Content -->
        <div class="p-4">
          <div class="mb-3">
            <div class="flex-1 min-w-0">
              <h3 class="text-lg font-semibold text-blue-900 truncate">
                {{ ticket.subject || 'No Subject' }}
              </h3>
            </div>
          </div>
          
          <div class="flex flex-wrap items-center justify-between gap-2">
            <!-- Left side: Ticket details in a single line -->
            <div class="flex flex-wrap items-center gap-1 text-sm text-blue-900">
              <!-- Ticket ID -->
              <span class="font-mono">{{ ticket.ticket_id }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Contact Name -->
              <span>{{ ticket.contact ? `${ticket.contact.first_name} ${ticket.contact.last_name}` : 'Unknown Contact' }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Account Name -->
              <span>{{ ticket.contact?.account?.account_name || ticket.account?.account_name || 'Unknown Account' }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Date Raised -->
              <span>{{ formatDate(ticket.created_at) }}</span>
            </div>
            
            <!-- Right side: Status, Priority, and Delete -->
            <div class="flex items-center gap-3">
              <!-- Status Badge -->
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" 
                    :class="getStatusBadgeClass(ticket.ticket_status)">
                {{ ticket.ticket_status }}
              </span>
              
              <!-- Priority Badge -->
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="getPriorityBadgeClass(ticket.priority)">
                {{ ticket.priority || 'Medium' }}
              </span>
              
              <!-- Delete Button -->
              <button 
                @click.stop="openDelete(ticket)"
                class="p-1.5 rounded-full hover:bg-red-100 transition-colors group"
                title="Delete Ticket"
              >
                <svg class="w-4 h-4 text-red-600 group-hover:text-red-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1-1H9a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-if="!filteredTickets.length" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-blue-900">No tickets found</h3>
        <p class="mt-1 text-sm text-blue-500">Get started by creating a new ticket.</p>
      </div>
    </div>
    <!-- Edit Details Modal -->
    <Modal :open="showEdit" title="Edit Ticket Details" @close="closeEdit">
      <div class="mb-3">
        <label class="block mb-1 text-blue-800 font-medium">Product</label>
        <select v-model="editForm.product_id" class="w-full border rounded px-3 py-2">
          <option v-for="prod in products" :key="prod.id" :value="String(prod.id)">{{ prod.product_name }}</option>
        </select>
      </div>
      <div class="mb-3">
        <label class="block mb-1 text-blue-800 font-medium">Subject</label>
        <input v-model="editForm.subject" type="text" class="w-full border rounded px-3 py-2" />
      </div>
      <div class="mb-3">
        <label class="block mb-1 text-blue-800 font-medium">Priority</label>
        <select v-model="editForm.priority" class="w-full border rounded px-3 py-2">
          <option value="High">High</option>
          <option value="Medium">Medium</option>
          <option value="Low">Low</option>
        </select>
      </div>
      <div class="mb-3">
        <label class="block mb-1 text-blue-800 font-medium">Details</label>
        <textarea v-model="editForm.ticket_details" rows="3" class="w-full border rounded px-3 py-2"></textarea>
      </div>
      <div class="flex gap-2 justify-end">
        <Button size="sm" @click="saveEdit">Save</Button>
        <Button size="sm" variant="secondary" @click="closeEdit">Cancel</Button>
      </div>
      <div v-if="editError" class="text-red-600 mt-2">{{ editError }}</div>
    </Modal>
    <!-- Assign Engineer Modal -->
    <Modal :open="showAssign" title="Assign Engineer" @close="closeAssign">
  <div class="mb-3">
    <label class="block mb-1 text-blue-800 font-medium">Engineer</label>
    <select v-model="selectedEngineer" class="w-full border rounded px-3 py-2">
      <option value="">No engineers assigned</option>
      <option v-for="eng in engineers" :key="eng.id" :value="String(eng.id)">
        {{ formatUserName(eng) }}
      </option>
    </select>
  </div>
  <div class="flex gap-2 justify-end">
    <Button size="sm" @click="assignEngineerToTicket">Save</Button>
    <Button size="sm" variant="secondary" @click="closeAssign">Cancel</Button>
  </div>
  <div v-if="assignError" class="text-red-600 mt-2">{{ assignError }}</div>
</Modal>
    <!-- Ticket Details Modal -->
    <Modal :open="showDetail" title="Ticket Details" @close="closeDetail">
      <div v-if="detailTicket" class="space-y-2">
                <div><span class="font-semibold text-gray-700">Ticket ID:</span> {{ detailTicket.ticket_id }}</div>
        <div><span class="font-semibold text-gray-700">Product:</span> {{ detailTicket.product?.product_name || '-' }}</div>
        <div><span class="font-semibold text-gray-700">Issue:</span> {{ detailTicket.issue?.issue_name || 'N/A' }}</div>
        <div><span class="font-semibold text-gray-700">Subject:</span> {{ detailTicket.subject }}</div>
        <div><span class="font-semibold text-gray-700">Priority:</span> <span :class="priorityClass(detailTicket.priority)">{{ detailTicket.priority }}</span></div>
        <div><span class="font-semibold text-gray-700">Details:</span> {{ detailTicket.ticket_details }}</div>
        <div><span class="font-semibold text-gray-700">Status:</span> {{ detailTicket.ticket_status }}</div>
        <div><span class="font-semibold text-gray-700">Created:</span> {{ formatDate(detailTicket.created_at) }}</div>
        <div><span class="font-semibold text-gray-700">Assigned Engineer:</span>
          <span v-if="detailTicket.engineer && detailTicket.engineer.first_name">
            {{ detailTicket.engineer.first_name }} {{ detailTicket.engineer.last_name }}<span v-if="detailTicket.engineer.phone"> ({{ detailTicket.engineer.phone }})</span>
          </span>
          <span v-else>No engineers assigned</span>
        </div>
        <div><span class="font-semibold text-gray-700">Contact Details:</span> {{ detailTicket.contact?.first_name }} {{ detailTicket.contact?.last_name }}<span v-if="detailTicket.contact?.mobile"> ({{ detailTicket.contact.mobile }})</span></div>
        <div><span class="font-semibold text-gray-700">Account Detail:</span> {{ detailTicket.account?.account_name }}</div>
      </div>
    </Modal>
    <!-- Change Status Modal -->
    <Modal :open="showStatus" @close="closeStatus">
      <div class="p-4 w-full max-w-md">
        <h2 class="text-lg font-semibold text-blue-800 mb-2">Change Ticket Status</h2>
        <select v-model="selectedStatus" class="w-full border rounded px-3 py-2 mb-4">
  <option disabled value="">Please select status</option>
          <option v-for="status in statusOptions" :key="status" :value="status">{{ status }}</option>
        </select>
        <div class="flex gap-2 justify-end">
          <Button size="sm" @click="changeStatusOfTicket">Change</Button>
          <Button size="sm" variant="secondary" @click="closeStatus">Cancel</Button>
        </div>
        <div v-if="statusError" class="text-red-600 mt-2">{{ statusError }}</div>
      </div>
    </Modal>

    <!-- Delete Confirmation Modal -->
    <Modal :open="showDelete" title="Delete Ticket" @close="closeDelete">
      <div class="p-4 w-full max-w-md">
        <div class="flex items-center mb-4">
          <svg class="w-12 h-12 text-red-600 mr-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.314 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <div>
            <h3 class="text-lg font-semibold text-red-800 mb-1">Delete Ticket</h3>
            <p class="text-sm text-gray-600">Are you sure you want to delete ticket <span class="font-mono font-semibold text-red-700">{{ selectedTicket?.ticket_id }}</span>?</p>
            <!-- <p class="text-xs text-gray-500 mt-1">This action cannot be undone.</p> -->
          </div>
        </div>
        <div class="flex gap-2 justify-end">
          <Button size="sm" variant="danger" @click="deleteTicketFromList">Delete</Button>
          <Button size="sm" variant="secondary" @click="closeDelete">Cancel</Button>
        </div>
        <div v-if="deleteError" class="text-red-600 mt-2 text-sm">{{ deleteError }}</div>
      </div>
    </Modal>

  </div>
</template>
<script setup>
import { fetchTicketDetail } from '@/api/tickets.js';
import { formatDateIST } from '@/utils/date';
import { ref, onMounted, computed, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import Button from '@/components/ui/Button.vue';
import Modal from '@/components/ui/Modal.vue';
import {
  fetchAllTickets,
  fetchEngineers,
  fetchUsers,
  assignEngineer,
  changeTicketStatus,
  updateTicket,
  fetchProducts,
  deleteTicket,
} from '@/api/tickets.js';
import { formatUserName } from '@/utils/user';

const router = useRouter();
const route = useRoute();
const formatDate = formatDateIST;
const tickets = ref([]);
const statusFilter = ref(null);
const engineers = ref([]);
const products = ref([]);
const showEdit = ref(false);
const showAssign = ref(false);
const showStatus = ref(false);
const showDetail = ref(false);
const showDelete = ref(false);
const detailTicket = ref(null);
const selectedTicket = ref(null);
const selectedEngineer = ref('');
const selectedStatus = ref('');
const assignError = ref('');
const statusError = ref('');
const deleteError = ref('');
const editError = ref('');
const editForm = ref({ product_id: '', subject: '', priority: '', ticket_details: '' });
const statusOptions = [
  'OPEN',
  'IN PROGRESS',
  'RESOLVED',
  'CLOSED'
];


function statusClass(status) {
  switch (status) {
    case 'OPEN': return 'text-blue-700 font-semibold';
    case 'IN PROGRESS':
    case 'IN_PROGRESS': return 'text-yellow-600 font-semibold';
    case 'RESOLVED': return 'text-green-600 font-semibold';
    case 'CLOSED': return 'text-gray-400 font-semibold';
    default: return '';
  }
}

function priorityClass(priority) {
  switch (priority) {
    case 'High': return 'text-red-600 font-semibold';
    case 'Medium': return 'text-yellow-600 font-semibold';
    case 'Low': return 'text-green-600 font-semibold';
    default: return 'text-gray-600';
  }
}

function getStatusBadgeClass(status) {
  switch (status) {
    case 'OPEN': return 'bg-blue-100 text-blue-800';
    case 'MEETING LOCKED IN WITH OEM': return 'bg-indigo-100 text-indigo-800';
    case 'PARTS ORDERED': return 'bg-purple-100 text-purple-800';
    case 'IN PROGRESS':
    case 'IN_PROGRESS': return 'bg-yellow-100 text-yellow-800';
    case 'RESOLVED': return 'bg-green-100 text-green-800';
    case 'CLOSED': return 'bg-gray-100 text-gray-800';
    case 'ON HOLD': return 'bg-orange-100 text-orange-800';
    case 'ESCALATED': return 'bg-red-100 text-red-800';
    default: return 'bg-gray-100 text-gray-800';
  }
}

function getPriorityBadgeClass(priority) {
  switch (priority) {
    case 'High': return 'bg-red-100 text-red-800';
    case 'Medium': return 'bg-yellow-100 text-yellow-800';
    case 'Low': return 'bg-green-100 text-green-800';
    default: return 'bg-gray-100 text-gray-800';
  }
}

// Computed property for page title based on filter
const pageTitle = computed(() => {
  if (statusFilter.value) {
    return `${statusFilter.value} Tickets`
  }
  return 'All Tickets'
})

// Computed property for filtered tickets
const filteredTickets = computed(() => {
  if (!statusFilter.value) {
    return tickets.value
  }
  const filterVal = statusFilter.value.toUpperCase();
  return tickets.value.filter(ticket => {
    const status = (ticket.ticket_status || '').toUpperCase();
    if (filterVal === 'IN PROGRESS') {
      return status === 'IN PROGRESS' || status === 'IN_PROGRESS';
    }
    return status === filterVal;
  })
})

// Clear filter and navigate back to all tickets
const clearFilter = () => {
  router.push('/manager/tickets')
}

async function loadTickets() {
  const res = await fetchAllTickets();
  console.log('fetchAllTickets response:', res);
  tickets.value = res.tickets || [];
  console.log('tickets.value after load:', tickets.value);
}
async function loadEngineers() {
  const res = await fetchUsers();
  engineers.value = Array.isArray(res) ? res : (res.users || []);
}
async function loadProducts() {
  const res = await fetchProducts();
  products.value = res.products || res;
}

async function openEdit(ticket) {
  selectedTicket.value = ticket;
  editError.value = '';
  await loadProducts();
  editForm.value = {
    product_id: ticket.product_id ? String(ticket.product_id) : '',
    subject: ticket.subject || '',
    priority: ticket.priority || 'Medium',
    ticket_details: ticket.ticket_details
  };
  showEdit.value = true;
}
function closeEdit() {
  showEdit.value = false;
  selectedTicket.value = null;
}
async function saveEdit() {
  if (!editForm.value.product_id || !editForm.value.subject || !editForm.value.ticket_details) {
    editError.value = 'All fields are required.';
    return;
  }
  try {
    await updateTicket(selectedTicket.value.id, {
      product_id: Number(editForm.value.product_id),
      subject: editForm.value.subject,
      priority: editForm.value.priority,
      ticket_details: editForm.value.ticket_details
    });
    showEdit.value = false;
    await loadTickets();
  } catch (e) {
    editError.value = e.response?.data?.error || 'Failed to update ticket.';
  }
}

async function openAssign(ticket) {
  selectedTicket.value = ticket;
  assignError.value = '';
  await loadEngineers();
  selectedEngineer.value = ticket.assigned_engineer ? String(ticket.assigned_engineer) : '';
  showAssign.value = true;
}
function closeAssign() {
  showAssign.value = false;
  selectedTicket.value = null;
  selectedEngineer.value = '';
}
const showSuccess = ref(false);
const successMessage = ref('');

async function assignEngineerToTicket() {
  // Allow unassigning engineer if empty string
  let engineerId = selectedEngineer.value ? Number(selectedEngineer.value) : null;
  try {
    await assignEngineer(selectedTicket.value.id, engineerId);
    showAssign.value = false;
    await loadTickets();
    successMessage.value = 'Engineer assigned successfully!';
    showSuccess.value = true;
    setTimeout(() => { showSuccess.value = false; }, 2000);
  } catch (e) {
    assignError.value = e.response?.data?.error || 'Failed to assign engineer.';
  }
}
function openStatus(ticket) {
  selectedTicket.value = ticket;
  // Fallback: ensure selectedStatus is never undefined
  selectedStatus.value = ticket.ticket_status || '';
  statusError.value = '';
  showStatus.value = true;
}
function closeStatus() {
  showStatus.value = false;
  selectedTicket.value = null;
  selectedStatus.value = '';
}
async function changeStatusOfTicket() {
  if (!selectedStatus.value) {
    statusError.value = 'Please select a status.';
    return;
  }
  try {
    await changeTicketStatus(selectedTicket.value.id, selectedStatus.value);
    showStatus.value = false;
    await loadTickets();
    successMessage.value = 'Status changed successfully!';
    showSuccess.value = true;
    setTimeout(() => { showSuccess.value = false; }, 2000);
  } catch (e) {
    statusError.value = e.response?.data?.error || 'Failed to change status.';
  }
}

// Delete ticket functions
function openDelete(ticket) {
  selectedTicket.value = ticket;
  deleteError.value = '';
  showDelete.value = true;
}

function closeDelete() {
  showDelete.value = false;
  selectedTicket.value = null;
  deleteError.value = '';
}

async function deleteTicketFromList() {
  if (!selectedTicket.value) {
    deleteError.value = 'No ticket selected for deletion.';
    return;
  }
  
  try {
    // Call backend API to delete ticket
    const response = await deleteTicket(selectedTicket.value.id);
    
    console.log(`🗑️ TICKET DELETED: ${selectedTicket.value.ticket_id} (ID: ${selectedTicket.value.id})`);
    
    // Close modal
    showDelete.value = false;
    
    // Refresh ticket list from backend
    await loadTickets();
    
    // Show success message
    const message = response.message || `Ticket ${selectedTicket.value.ticket_id} deleted successfully!`;
    successMessage.value = message;
    showSuccess.value = true;
    setTimeout(() => { showSuccess.value = false; }, 3000);
    
    selectedTicket.value = null;
    
  } catch (e) {
    console.error('Delete ticket error:', e);
    deleteError.value = e.response?.data?.error || 'Failed to delete ticket.';
  }
}

// Watch for route query changes
watch(() => route.query.status, (newStatus) => {
  statusFilter.value = newStatus || null
}, { immediate: true })

onMounted(() => {
  loadTickets();
});

async function openDetail(ticket) {
  // Navigate to the dedicated ticket detail page
  router.push(`/manager/tickets/${ticket.id}`);
}

function closeDetail() {
  showDetail.value = false;
  detailTicket.value = null;
}

</script>
<style scoped>
table {
  border-collapse: collapse;
}
th, td {
  border-bottom: 1px solid #e0e7ff;
}
</style>
