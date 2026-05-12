<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="closeModal">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-6xl max-h-[90vh] overflow-hidden" @click.stop>
      <!-- Modal Header -->
      <div class="bg-blue-600 text-white px-6 py-4 flex items-center justify-between">
        <div class="flex-1">
          <h2 class="text-xl font-semibold">{{ ticket?.subject || 'Ticket Details' }}</h2>
          <p class="text-blue-100 text-sm mt-1">{{ ticket?.ticket_id }}</p>
        </div>
        <button @click="closeModal" class="text-white hover:text-blue-200 transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Ticket Summary -->
      <div class="bg-blue-50 px-6 py-4 border-b border-blue-100">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 text-sm">
          <div>
            <span class="font-medium text-blue-700">Status:</span>
            <span class="ml-2 px-2 py-1 rounded-full text-xs font-medium" :class="getStatusBadgeClass(ticket?.ticket_status)">
              {{ ticket?.ticket_status }}
            </span>
          </div>
          <div>
            <span class="font-medium text-blue-700">Priority:</span>
            <span class="ml-2 px-2 py-1 rounded-full text-xs font-medium" :class="getPriorityBadgeClass(ticket?.priority)">
              {{ ticket?.priority }}
            </span>
          </div>
          <div>
            <span class="font-medium text-blue-700">Contact:</span>
            <span class="ml-2 text-blue-900">{{ getContactName(ticket) }}</span>
          </div>
          <div>
            <span class="font-medium text-blue-700">Account:</span>
            <span class="ml-2 text-blue-900">{{ getAccountName(ticket) }}</span>
          </div>
        </div>
      </div>

      <!-- Tab Navigation -->
      <div class="bg-white border-b border-gray-200">
        <nav class="flex space-x-8 px-6" aria-label="Tabs">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-colors',
              activeTab === tab.id
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            {{ tab.name }}
            <span v-if="tab.count !== undefined" class="ml-2 bg-gray-100 text-gray-600 py-0.5 px-2 rounded-full text-xs">
              {{ tab.count }}
            </span>
          </button>
        </nav>
      </div>

      <!-- Tab Content -->
      <div class="flex-1 overflow-hidden">
        <!-- Conversation Tab -->
        <ConversationTab
          v-if="activeTab === 'conversation'"
          :ticket-id="ticketId"
          :ticket="ticket"
          @comment-added="handleCommentAdded"
        />

        <!-- Calls Tab -->
        <CallsTab
          v-if="activeTab === 'calls'"
          :ticket-id="ticketId"
          :ticket="ticket"
          @call-scheduled="handleCallScheduled"
        />

        <!-- History Tab -->
        <HistoryTab
          v-if="activeTab === 'history'"
          :ticket-id="ticketId"
          :ticket="ticket"
        />
      </div>

      <!-- Modal Footer -->
      <div class="bg-gray-50 px-6 py-4 flex justify-end space-x-3 border-t border-gray-200">
        <button
          @click="closeModal"
          class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
        >
          Close
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { fetchTicketFullDetails } from '../../api/tickets'
import ConversationTab from './tabs/ConversationTab.vue'
import CallsTab from './tabs/CallsTab.vue'
import HistoryTab from './tabs/HistoryTab.vue'

const props = defineProps({
  ticketId: {
    type: [String, Number],
    required: true
  },
  initialTab: {
    type: String,
    default: 'conversation'
  }
})

const emit = defineEmits(['close', 'ticket-updated'])

const ticket = ref(null)
const loading = ref(false)
const error = ref('')
const activeTab = ref(props.initialTab)

const tabs = computed(() => [
  {
    id: 'conversation',
    name: 'Conversation',
    count: ticket.value?.counts?.comments
  },
  {
    id: 'calls',
    name: 'Calls',
    count: ticket.value?.counts?.calls
  },
  {
    id: 'history',
    name: 'History',
    count: ticket.value?.counts?.activities
  }
])

// Load ticket data
const loadTicketData = async () => {
  if (!props.ticketId) return
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetchTicketFullDetails(props.ticketId)
    ticket.value = response.ticket
    
    // Update counts if available
    if (response.counts) {
      ticket.value.counts = response.counts
    }
  } catch (err) {
    console.error('Failed to load ticket details:', err)
    error.value = 'Failed to load ticket details'
  } finally {
    loading.value = false
  }
}

// Event handlers
const closeModal = () => {
  emit('close')
}

const handleCommentAdded = () => {
  // Refresh ticket data to update counts
  loadTicketData()
  emit('ticket-updated')
}

const handleCallScheduled = () => {
  // Refresh ticket data to update counts
  loadTicketData()
  emit('ticket-updated')
}

// Utility functions
const getStatusBadgeClass = (status) => {
  const classes = {
    'OPEN': 'bg-blue-100 text-blue-800',
    'IN_PROGRESS': 'bg-yellow-100 text-yellow-800',
    'RESOLVED': 'bg-green-100 text-green-800',
    'CLOSED': 'bg-gray-100 text-gray-800',
    'ESCALATED': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getPriorityBadgeClass = (priority) => {
  const classes = {
    'High': 'bg-red-100 text-red-800',
    'Medium': 'bg-yellow-100 text-yellow-800',
    'Low': 'bg-green-100 text-green-800'
  }
  return classes[priority] || 'bg-gray-100 text-gray-800'
}

const getContactName = (ticket) => {
  if (!ticket?.contact) return 'Unknown Contact'
  return `${ticket.contact.first_name} ${ticket.contact.last_name || ''}`.trim()
}

const getAccountName = (ticket) => {
  return ticket?.contact?.account?.account_name || ticket?.account?.account_name || 'Unknown Account'
}

// Lifecycle
onMounted(() => {
  loadTicketData()
})

// Watch for ticket ID changes
watch(() => props.ticketId, () => {
  loadTicketData()
})
</script>

<style scoped>
/* Custom scrollbar for modal content */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
