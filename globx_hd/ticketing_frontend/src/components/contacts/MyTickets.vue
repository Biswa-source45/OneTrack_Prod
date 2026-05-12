<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-7xl mx-auto mt-8">
    <h1 class="text-2xl font-bold text-blue-800 mb-6">My Tickets</h1>
    
    <!-- Tickets Card Layout -->
    <div class="space-y-4">
      <div v-for="ticket in tickets" :key="ticket.id" 
           class="bg-white border border-blue-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer"
           @click="openDetail(ticket)">
        
        <!-- Card Content -->
        <div class="p-4">
          <div class="flex items-start justify-between mb-3">
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
              
              <!-- Product Name -->
              <span>{{ ticket.product?.product_name || 'Unknown Product' }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Issue Name -->
              <span>{{ ticket.issue?.issue_name || 'N/A' }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Date Raised -->
              <span>{{ formatDate(ticket.created_at) }}</span>
            </div>
            
            <!-- Right side: Status and Priority -->
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
            </div>
          </div>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-if="!tickets.length" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-blue-900">No tickets found</h3>
        <p class="mt-1 text-sm text-blue-500">Get started by creating a new ticket.</p>
      </div>
    </div>

    <!-- Error Modal -->
    <Modal v-if="error" @close="error = ''">
      <div class="text-red-700 font-semibold">{{ error }}</div>
    </Modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { fetchMyTickets } from '../../api/tickets';
import Modal from '../ui/Modal.vue';
import { formatDateIST } from '../../utils/date';

const router = useRouter();
const formatDate = formatDateIST;
const auth = useAuthStore();
const tickets = ref([]);
const error = ref('');

// Badge styling functions (matching Manager implementation exactly)
function getStatusBadgeClass(status) {
  switch (status) {
    case 'OPEN': return 'bg-blue-100 text-blue-800';
    case 'MEETING LOCKED IN WITH OEM': return 'bg-indigo-100 text-indigo-800';
    case 'PARTS ORDERED': return 'bg-purple-100 text-purple-800';
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

// Load tickets on component mount
onMounted(async () => {
  try {
    const result = await fetchMyTickets(auth.user?.id);
    tickets.value = Array.isArray(result.tickets) ? result.tickets : [];
  } catch (err) {
    error.value = err.message || 'Failed to fetch tickets.';
  }
});

// Navigate to ticket detail page (matching existing contact route structure)
function openDetail(ticket) {
  router.push(`/contacts/my-tickets/${ticket.id}`);
}
</script>
