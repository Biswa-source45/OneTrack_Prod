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
          
          <!-- Action Buttons -->
          <div class="flex items-center space-x-3">
            <!-- Empty space - triple-dot menu moved to Conversation tab -->
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
          <!-- Ticket Properties -->
          <PropertySection 
            title="Ticket Properties" 
            :show-edit-icon="true"
            @edit="openEditModal"
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
                <span class="text-xs text-gray-500 uppercase tracking-wide">Ticket Owner</span>
                <div class="mt-1">
                  <InlineEditDropdown
                    :model-value="ticket.assigned_engineer ? String(ticket.assigned_engineer) : ''"
                    :options="usersWithFullName"
                    option-value="id"
                    option-label="full_name"
                    empty-label="Unassigned"
                    :loading="usersLoading"
                    @change="handleTicketOwnerChange"
                  />
                </div>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Status</span>
                <button
                  @click="openStatusChangeModal"
                  :class="`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium hover:opacity-80 transition-opacity cursor-pointer ${getStatusBadgeClass(ticket.ticket_status)}`"
                >
                  {{ ticket.ticket_status }}
                </button>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Date Created</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ formatDateTime(ticket.created_at) }}</p>
              </div>
            </div>
          </PropertySection>

          <!-- Ticket Information -->
          <PropertySection title="Ticket Information">
            <div class="space-y-3">
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Product Name</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ ticket.product?.product_name || 'N/A' }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Ticket Details</span>
                <p class="text-sm text-gray-900 mt-1 whitespace-pre-wrap">{{ ticket.ticket_details || 'No details provided' }}</p>
              </div>
            </div>
          </PropertySection>

          <!-- Additional Information -->
          <PropertySection title="Additional Information">
            <div class="space-y-3">
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Priority</span>
                <div class="mt-1">
                  <InlineEditDropdown
                    :model-value="ticket.priority"
                    :options="priorityOptions"
                    option-value="value"
                    option-label="label"
                    empty-label="Select Priority"
                    :allow-empty="false"
                    :display-class="`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityBadgeClass(ticket.priority)}`"
                    @change="handlePriorityChange"
                  />
                </div>
              </div>
            </div>
          </PropertySection>
        </div>
      </div>

      <!-- Right Content Area -->
      <div class="flex-1 bg-white">
        <div v-if="ticket" class="h-full">
          <!-- Tab Navigation -->
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
            <!-- Conversation Tab -->
            <ConversationTab
              v-if="activeTab === 'conversation'"
              :ticket-id="ticketId"
              :ticket="ticket"
              @comment-added="handleCommentAdded"
              @close-ticket="handleCloseTicket"
            />

            <!-- Calls Tab -->
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
            <ManagerAttachmentsTab
              v-if="activeTab === 'attachments'"
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
      @close="notification.show = false"
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
import { fetchTicketFullDetails, updateTicket, fetchUsers, changeTicketStatus } from '../../api/tickets'
import { formatDateTime } from '../../utils/date'
import { formatUserName } from '../../utils/user'
import ConversationTab from './tabs/ConversationTab.vue'
import CallsTab from './tabs/CallsTab.vue'
import ApprovalsTab from './tabs/ApprovalsTab.vue'
import ManagerAttachmentsTab from './tabs/ManagerAttachmentsTab.vue'
import PropertySection from './PropertySection.vue'
import InlineEditDropdown from './InlineEditDropdown.vue'
import NotificationToast from './NotificationToast.vue'
import StatusChangeModal from './StatusChangeModal.vue'

const route = useRoute()
const router = useRouter()

const ticketId = computed(() => route.params.id)
const ticket = ref(null)
const loading = ref(false)
const error = ref('')
const activeTab = ref('conversation')

// Data for inline editing
const users = ref([])
const usersLoading = ref(false)
const statusOptions = [
  { value: 'OPEN', label: 'OPEN' },
  { value: 'IN PROGRESS', label: 'IN PROGRESS' },
  { value: 'RESOLVED', label: 'RESOLVED' },
  { value: 'CLOSED', label: 'CLOSED' }
]
const priorityOptions = [
  { value: 'High', label: 'High' },
  { value: 'Medium', label: 'Medium' },
  { value: 'Low', label: 'Low' }
]

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

// Computed property to add full_name to all users
const usersWithFullName = computed(() => {
  console.log('🔧 DEBUG: Ticket Owner dropdown - All users:', users.value.length);
  return users.value.map(user => ({
    ...user,
    full_name: formatUserName(user)
  }))
})

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
  }
])


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

// Load users for ticket owner dropdown
const loadUsers = async () => {
  usersLoading.value = true
  try {
    const response = await fetchUsers()
    users.value = Array.isArray(response) ? response : (response?.users || [])
  } catch (err) {
    console.error('Failed to load users:', err)
  } finally {
    usersLoading.value = false
  }
}

// Navigation
const goBack = () => {
  router.go(-1) // Go back to previous page
}

// Event handlers
const handleCommentAdded = () => {
  // Refresh ticket data to update counts
  loadTicketData()
}

const handleCloseTicket = () => {
  // Open status change modal with pre-selected status
  if (ticket.value?.ticket_status === 'CLOSED') {
    // Reopen ticket - pre-select OPEN status
    openStatusChangeModalWithStatus('OPEN')
  } else {
    // Close ticket - pre-select CLOSED status
    openStatusChangeModalWithStatus('CLOSED')
  }
}

// Open status change modal with pre-selected status
const openStatusChangeModalWithStatus = (preSelectedStatus) => {
  // Set the pre-selected status in the modal
  preSelectedStatusForModal.value = preSelectedStatus
  showStatusChangeModal.value = true
}

const handleCallScheduled = () => {
  // Refresh ticket data to update counts
  loadTicketData()
}

const handleApprovalUpdated = () => {
  // Refresh ticket data to get updated counts
  loadTicketData()
}

const openEditModal = () => {
  // Navigate to edit page instead of opening modal
  const currentPath = router.currentRoute.value.path
  const basePath = currentPath.includes('/manager/') ? '/manager' : '/engineer'
  router.push(`${basePath}/tickets/${ticketId.value}/edit`)
}

// Notification helper
const showNotification = (type, title, message = '') => {
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

// Inline edit handlers
const handleTicketOwnerChange = async (newOwnerId) => {
  try {
    // Handle empty string or null/undefined as unassigned (null)
    // parseInt("") returns NaN, so we need explicit check
    let assignedEngineerId = null
    if (newOwnerId && newOwnerId !== '') {
      assignedEngineerId = parseInt(newOwnerId)
    }
    
    const updatePayload = {
      assigned_engineer: assignedEngineerId
    }
    await updateTicket(ticket.value.id, updatePayload)
    
    // Show success notification
    const ownerName = newOwnerId 
      ? usersWithFullName.value.find(u => u.id == newOwnerId)?.full_name || 'Unknown User'
      : 'Unassigned'
    showNotification('success', 'Ticket Owner Updated', `Assigned to: ${ownerName}`)
    
    // Refresh ticket data to update display and history
    await loadTicketData()
  } catch (err) {
    console.error('Failed to update ticket owner:', err)
    showNotification('error', 'Update Failed', 'Could not update ticket owner')
  }
}

// Open status change modal
const openStatusChangeModal = () => {
  preSelectedStatusForModal.value = '' // Clear any pre-selected status
  showStatusChangeModal.value = true
}

// Handle status change with remarks
const handleStatusChange = async (changeData) => {
  try {
    // Use the dedicated status change API with remarks
    await changeTicketStatus(ticket.value.id, changeData.status, changeData.remarks)
    
    // Show success notification
    showNotification('success', 'Status Updated', `Changed to: ${changeData.status}`)
    
    // Close modal and refresh ticket data
    showStatusChangeModal.value = false
    preSelectedStatusForModal.value = '' // Clear pre-selected status
    await loadTicketData()
  } catch (err) {
    console.error('Failed to update ticket status:', err)
    showNotification('error', 'Update Failed', 'Could not update ticket status')
    throw err // Re-throw to let modal handle the error
  }
}

const handlePriorityChange = async (newPriority) => {
  try {
    const updatePayload = {
      priority: newPriority
    }
    await updateTicket(ticket.value.id, updatePayload)
    
    // Show success notification
    showNotification('success', 'Priority Updated', `Changed to: ${newPriority}`)
    
    // Refresh ticket data to update display and history
    await loadTicketData()
  } catch (err) {
    console.error('Failed to update ticket priority:', err)
    showNotification('error', 'Update Failed', 'Could not update ticket priority')
  }
}


// Utility functions
const getStatusBadgeClass = (status) => {
  const classes = {
    'OPEN': 'bg-blue-100 text-blue-800',
    'IN PROGRESS': 'bg-yellow-100 text-yellow-800',
    'RESOLVED': 'bg-green-100 text-green-800',
    'CLOSED': 'bg-gray-100 text-gray-800'
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

const getAssignedEngineer = (ticket) => {
  if (!ticket?.engineer) return 'Unassigned'
  return `${ticket.engineer.first_name} ${ticket.engineer.last_name || ''}`.trim()
}

// Lifecycle
onMounted(() => {
  loadTicketData()
  loadUsers()
})

// Watch for route changes
watch(() => route.params.id, () => {
  if (route.params.id) {
    loadTicketData()
  }
})
</script>

<style scoped>
/* Custom scrollbar for page content */
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
