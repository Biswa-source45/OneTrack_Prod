<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header with Back Button -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-full px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center space-x-4">
            <button
              @click="goBack"
              class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              Back to Tickets
            </button>
            <div class="h-6 w-px bg-gray-300"></div>
            <div v-if="ticket">
              <h1 class="text-xl font-semibold text-gray-900">{{ ticket.subject || 'Ticket Details' }}</h1>
              <p class="text-sm text-gray-500">{{ ticket.ticket_id }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content Layout -->
    <div class="flex max-w-full">
      <!-- Left Sidebar - Properties Panel -->
      <div class="w-80 bg-white border-r border-gray-200 min-h-screen">
        <div v-if="loading" class="flex justify-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
        
        <div v-else-if="error" class="p-4 text-center">
          <div class="text-red-600 mb-4">
            <svg class="mx-auto h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <p class="text-sm text-gray-600 mb-4">{{ error }}</p>
          <button
            @click="loadTicketData"
            class="text-sm px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Try Again
          </button>
        </div>

        <div v-else-if="ticket" class="divide-y divide-gray-200">
          <!-- Ticket Properties (No edit icon for engineers) -->
          <PropertySection 
            title="Ticket Properties" 
            :show-edit-icon="false"
          >
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-600">Ticket ID</span>
                <span class="text-sm font-medium text-gray-900">{{ ticket.ticket_id }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-600">Subject</span>
                <span class="text-sm font-medium text-gray-900">{{ ticket.subject || 'No Subject' }}</span>
              </div>
            </div>
          </PropertySection>

          <!-- Contact Info -->
          <PropertySection title="Contact Info">
            <div class="space-y-3">
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Contact Name</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ getContactName(ticket) }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Account Name</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ getAccountName(ticket) }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Contact Email</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ ticket.contact?.email || 'N/A' }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Mobile Number</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ ticket.contact?.mobile || 'N/A' }}</p>
              </div>
            </div>
          </PropertySection>

          <!-- Key Information -->
          <PropertySection title="Key Information">
            <div class="space-y-3">
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Product</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ ticket.product?.product_name || 'N/A' }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Status</span>
                <div class="mt-1">
                  <!-- Status badge with click to change -->
                  <button
                    @click="openStatusChangeModal"
                    :class="`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium hover:opacity-80 transition-opacity cursor-pointer ${getStatusBadgeClass(ticket.ticket_status)}`"
                  >
                    {{ ticket.ticket_status }}
                  </button>
                </div>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Assigned Engineer</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ getEngineerName(ticket) }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Created</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ formatDateTime(ticket.created_at) }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Last Updated</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ formatDateTime(ticket.updated_at) }}</p>
              </div>
            </div>
          </PropertySection>

          <!-- Ticket Details -->
          <PropertySection title="Ticket Details">
            <div class="space-y-3">
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Description</span>
                <p class="text-sm text-gray-900 mt-1 whitespace-pre-wrap">{{ ticket.ticket_details || 'No description provided' }}</p>
              </div>
            </div>
          </PropertySection>

          <!-- Additional Information -->
          <PropertySection title="Additional Information">
            <div class="space-y-3">
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Priority</span>
                <div class="mt-1">
                  <!-- Read-only priority display for engineers -->
                  <span :class="`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityBadgeClass(ticket.priority)}`">
                    {{ ticket.priority || 'Medium' }}
                  </span>
                </div>
              </div>
            </div>
          </PropertySection>
        </div>
      </div>

      <!-- Right Content Area -->
      <div class="flex-1 bg-white">
        <div v-if="ticket" class="h-full">
          <!-- Tab Navigation (All tabs like manager) -->
          <div class="border-b border-gray-200 bg-white">
            <nav class="flex space-x-8 px-6" aria-label="Tabs">
              <button
                v-for="tab in tabs"
                :key="tab.id"
                @click="activeTab = tab.id"
                :class="[
                  'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-colors uppercase tracking-wide',
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
          <div class="h-full">
            <!-- Conversation Tab (Same as manager) -->
            <ConversationTab
              v-if="activeTab === 'conversation'"
              :ticket-id="ticketId"
              :ticket="ticket"
              @comment-added="handleCommentAdded"
              @close-ticket="handleCloseTicket"
            />

            <!-- Calls Tab (Same as manager) -->
            <CallsTab
              v-if="activeTab === 'calls'"
              :ticket-id="ticketId"
              :ticket="ticket"
              @call-scheduled="handleCallScheduled"
              @show-notification="showNotification"
            />

            <!-- Approvals Tab -->
            <ApprovalsTab
              v-if="activeTab === 'approvals'"
              :ticket-id="ticketId"
              :ticket="ticket"
              @approval-updated="handleApprovalUpdated"
              @show-notification="showNotification"
            />

            <!-- Attachments Tab -->
            <EngineerAttachmentsTab
              v-if="activeTab === 'attachments'"
              :ticket-id="ticketId"
              :ticket="ticket"
            />

            <!-- History Tab (Same as manager) -->
            <HistoryTab
              v-if="activeTab === 'history'"
              :ticket-id="ticketId"
              :ticket="ticket"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Notification Toast -->
    <NotificationToast
      :show="notification.show"
      :type="notification.type"
      :title="notification.title"
      :message="notification.message"
      @close="closeNotification"
    />

    <!-- Status Change Modal -->
    <StatusChangeModal
      :open="showStatusChangeModal"
      :current-status="ticket?.ticket_status || ''"
      :status-options="statusOptions"
      :pre-selected-status="preSelectedStatusForModal"
      @close="showStatusChangeModal = false; preSelectedStatusForModal = ''"
      @change="handleStatusChange"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchTicketFullDetails } from '../../api/tickets'
import { changeEngineerTicketStatus } from '../../api/engineer'
import { formatDateTime } from '../../utils/date'
import ConversationTab from '../shared/tabs/ConversationTab.vue'
import CallsTab from '../shared/tabs/CallsTab.vue'
import ApprovalsTab from '../shared/tabs/ApprovalsTab.vue'
import HistoryTab from '../shared/tabs/HistoryTab.vue'
import EngineerAttachmentsTab from '../shared/tabs/EngineerAttachmentsTab.vue'
import PropertySection from '../shared/PropertySection.vue'
import InlineEditDropdown from '../shared/InlineEditDropdown.vue'
import NotificationToast from '../shared/NotificationToast.vue'
import StatusChangeModal from '../shared/StatusChangeModal.vue'

const route = useRoute()
const router = useRouter()

// Reactive state
const ticket = ref(null)
const loading = ref(false)
const error = ref('')
const activeTab = ref('conversation')

// Notification state
const notification = ref({
  show: false,
  type: 'success',
  title: '',
  message: ''
})

// Status change modal state
const showStatusChangeModal = ref(false)
const preSelectedStatusForModal = ref('')

// Computed properties
const ticketId = computed(() => route.params.id)

// Status options for engineers (cannot set to CLOSED)
const statusOptions = [
  { value: 'OPEN', label: 'OPEN' },
  { value: 'IN PROGRESS', label: 'IN PROGRESS' },
  { value: 'RESOLVED', label: 'RESOLVED' }
]

// Engineer tabs (same as manager)
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
    id: 'approvals',
    name: 'Approvals',
    count: ticket.value?.counts?.approvals
  },
  {
    id: 'attachments',
    name: 'Attachments',
    count: ticket.value?.counts?.attachments
  },
  {
    id: 'history',
    name: 'History',
    count: ticket.value?.counts?.activities
  }
])

// Badge styling functions (same as manager)
const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'OPEN': return 'bg-blue-100 text-blue-800'
    case 'IN PROGRESS': return 'bg-yellow-100 text-yellow-800'
    case 'RESOLVED': return 'bg-green-100 text-green-800'
    case 'CLOSED': return 'bg-gray-100 text-gray-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

const getPriorityBadgeClass = (priority) => {
  switch (priority) {
    case 'High': return 'bg-red-100 text-red-800'
    case 'Medium': return 'bg-yellow-100 text-yellow-800'
    case 'Low': return 'bg-green-100 text-green-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

// Helper functions
const getContactName = (ticket) => {
  if (!ticket?.contact) return 'N/A'
  return `${ticket.contact.first_name} ${ticket.contact.last_name || ''}`.trim()
}

const getAccountName = (ticket) => {
  return ticket?.contact?.account?.account_name || ticket?.account?.account_name || 'N/A'
}

const getEngineerName = (ticket) => {
  if (!ticket?.engineer) return 'Not Assigned'
  return `${ticket.engineer.first_name} ${ticket.engineer.last_name || ''}`.trim()
}

// Load ticket data
const loadTicketData = async () => {
  if (!ticketId.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetchTicketFullDetails(ticketId.value)
    ticket.value = response.ticket
    
    // Update counts if available
    if (response.counts) {
      ticket.value.counts = response.counts
    }
  } catch (err) {
    console.error('Failed to load ticket details:', err)
    error.value = 'Failed to load ticket details. Please try again.'
  } finally {
    loading.value = false
  }
}

// Open status change modal
const openStatusChangeModal = () => {
  preSelectedStatusForModal.value = '' // Clear any pre-selected status
  showStatusChangeModal.value = true
}

// Open status change modal with pre-selected status
const openStatusChangeModalWithStatus = (preSelectedStatus) => {
  preSelectedStatusForModal.value = preSelectedStatus
  showStatusChangeModal.value = true
}

// Handle status change with remarks (engineers can change status)
const handleStatusChange = async (changeData) => {
  if (!ticket.value || !changeData.status) return
  
  try {
    await changeEngineerTicketStatus(ticket.value.id, changeData.status, changeData.remarks)
    ticket.value.ticket_status = changeData.status
    showNotification('success', 'Status Updated', 'Ticket status has been updated successfully.')
    
    // Close modal and refresh ticket data
    showStatusChangeModal.value = false
    preSelectedStatusForModal.value = '' // Clear pre-selected status
    await loadTicketData()
  } catch (err) {
    console.error('Failed to update status:', err)
    showNotification('error', 'Update Failed', 'Failed to update ticket status. Please try again.')
    throw err // Re-throw to let modal handle the error
  }
}

// Event handlers (same as manager)
const handleCommentAdded = () => {
  loadTicketData() // Refresh to get updated counts
}

const handleCallScheduled = () => {
  loadTicketData() // Refresh to get updated counts
}

const handleApprovalUpdated = () => {
  loadTicketData() // Refresh to get updated counts
}

const handleCloseTicket = () => {
  loadTicketData() // Refresh ticket data
}

// Navigation
const goBack = () => {
  router.push('/engineer/tickets')
}

// Notification handling
const showNotification = (type, title, message) => {
  notification.value = {
    show: true,
    type,
    title,
    message
  }
}

const closeNotification = () => {
  notification.value.show = false
}

// Watchers
watch(() => route.params.id, (newId) => {
  if (newId) {
    loadTicketData()
  }
}, { immediate: true })

// Lifecycle
onMounted(() => {
  loadTicketData()
})
</script>
