<template>
  <div class="min-h-[600px] flex flex-col">
    <!-- Calls List -->
    <div class="flex-1 overflow-y-auto p-6">
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-lg font-medium text-gray-900">Call Logs</h3>
        <button
          @click="navigateToAddCall"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <span>Add Call</span>
        </button>
      </div>

      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
      
      <div v-else-if="calls.length === 0" class="text-center py-8 text-gray-500">
        <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
        </svg>
        <p class="text-lg font-medium">No calls logged</p>
        <p class="text-sm">Add call logs for inbound and outbound communications.</p>
      </div>

      <!-- Call Items -->
      <div v-else class="space-y-4">
        <div v-for="call in calls" :key="call.id" class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <div class="flex items-center space-x-2">
                  <div 
                    :class="[
                      'w-3 h-3 rounded-full',
                      getCallStatusColor(call.status)
                    ]"
                  ></div>
                  <span class="font-medium text-gray-900">{{ call.subject || 'Call Log' }}</span>
                  <span 
                    :class="[
                      'px-2 py-1 text-xs font-medium rounded-full',
                      getDirectionBadgeClass(call.direction)
                    ]"
                  >
                    {{ call.direction || 'Unknown' }}
                  </span>
                </div>
              </div>
              
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-600 mb-3">
                <div v-if="call.start_time">
                  <span class="font-medium">Start Time:</span>
                  {{ formatDateTime(call.start_time) }}
                </div>
                <div>
                  <span class="font-medium">Logged by:</span>
                  {{ getUserName(call.user) }}
                </div>
                <div>
                  <span class="font-medium">Status:</span>
                  <!-- Inline editing for managers only -->
                  <div v-if="isManager" class="ml-1 inline-block">
                    <InlineEditDropdown
                      :model-value="call.status"
                      :options="callStatusOptions"
                      option-value="value"
                      option-label="label"
                      empty-label="Select Status"
                      :allow-empty="false"
                      :display-class="`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getCallStatusBadgeClass(call.status)}`"
                      @change="(newStatus) => handleCallStatusChange(call, newStatus)"
                    />
                  </div>
                  <!-- Read-only display for engineers -->
                  <span 
                    v-else
                    :class="[
                      'px-2 py-1 text-xs font-medium rounded-full ml-1',
                      getCallStatusBadgeClass(call.status)
                    ]"
                  >
                    {{ call.status }}
                  </span>
                </div>
                <div v-if="call.oem_ticket_id">
                  <span class="font-medium">OEM Ticket ID:</span>
                  <span class="ml-1 text-blue-600">{{ call.oem_ticket_id }}</span>
                </div>
                <div v-if="call.due_date">
                  <span class="font-medium">Due Date:</span>
                  {{ formatDate(call.due_date) }}
                </div>
              </div>
              
              <div v-if="call.description" class="mb-3">
                <span class="font-medium text-gray-700">Description:</span>
                <p class="text-gray-600 mt-1 whitespace-pre-wrap">{{ call.description }}</p>
              </div>
              
              <div v-if="call.mail_content" class="mb-3">
                <span class="font-medium text-gray-700">Mail Content:</span>
                <div class="mt-1 p-3 bg-gray-50 border border-gray-200 rounded text-sm text-gray-700 whitespace-pre-wrap font-mono max-h-48 overflow-y-auto">{{ call.mail_content }}</div>
              </div>
              
              <div v-if="call.close_remarks" class="mb-3">
                <span class="font-medium text-gray-700">Close Remarks:</span>
                <div class="mt-1 p-3 bg-blue-50 border-l-4 border-blue-400 rounded text-sm">
                  <div class="text-blue-700">{{ call.close_remarks }}</div>
                </div>
              </div>
              
              <div v-if="call.attachments && call.attachments.length > 0" class="mb-3">
                <span class="font-medium text-gray-700">Attachments ({{ call.attachments.length }}):</span>
                <div class="mt-2 space-y-2">
                  <div
                    v-for="attachment in call.attachments"
                    :key="attachment.id"
                    class="flex items-center justify-between p-2 bg-gray-50 border border-gray-200 rounded hover:bg-gray-100 transition-colors"
                  >
                    <div class="flex items-center space-x-2 flex-1 min-w-0">
                      <svg class="w-4 h-4 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                      </svg>
                      <span class="text-sm text-gray-700 truncate">{{ attachment.original_filename }}</span>
                      <span class="text-xs text-gray-500">({{ formatFileSize(attachment.file_size) }})</span>
                    </div>
                    <button
                      @click="downloadAttachment(attachment)"
                      class="ml-2 text-blue-600 hover:text-blue-800 flex-shrink-0"
                      title="Download"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Call Actions -->
            <div class="flex flex-col space-y-2 ml-4">
              <button
                v-if="call.status === 'Open' || call.status === 'In Progress'"
                @click="cancelCall(call)"
                class="px-3 py-1 text-sm bg-red-600 text-white rounded hover:bg-red-700 transition-colors"
              >
                Cancel
              </button>
              <button
                v-if="isManager && call.status !== 'Completed'"
                @click="openCloseCallModal(call)"
                class="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
              >
                Close Call
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Schedule Call Modal -->
    <div v-if="showScheduleForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="closeScheduleForm">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md max-h-[80vh] overflow-y-auto" @click.stop>
        <div class="p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900">
              {{ editingCall ? 'Edit Call' : 'Schedule New Call' }}
            </h3>
            <button @click="closeScheduleForm" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="submitCall" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Call Type *</label>
              <select
                v-model="callForm.call_type"
                class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              >
                <option value="">Select call type</option>
                <option value="OEM">OEM</option>
                <option value="Customer">Customer</option>
                <option value="Internal">Internal</option>
                <option value="Vendor">Vendor</option>
                <option value="Support">Support</option>
              </select>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Party Name *</label>
              <input
                type="text"
                v-model="callForm.party_name"
                class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="e.g., Dell Support, John Doe"
                required
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Contact Information</label>
              <input
                type="text"
                v-model="callForm.party_contact"
                class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Phone number or email"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Scheduled Date & Time *</label>
              <input
                type="datetime-local"
                v-model="callForm.scheduled_at"
                class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Purpose</label>
              <textarea
                v-model="callForm.purpose"
                class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                rows="3"
                placeholder="Purpose of the call..."
              ></textarea>
            </div>
            
            <div class="flex justify-end space-x-3 pt-4">
              <button
                type="button"
                @click="closeScheduleForm"
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                {{ submitting ? 'Saving...' : (editingCall ? 'Update Call' : 'Schedule Call') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Complete Call Modal -->
    <div v-if="showCompleteForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="closeCompleteForm">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md" @click.stop>
        <div class="p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900">Complete Call</h3>
            <button @click="closeCompleteForm" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="submitComplete" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Call Notes</label>
              <textarea
                v-model="completeForm.notes"
                class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                rows="4"
                placeholder="Add notes about the call outcome, next steps, etc."
              ></textarea>
            </div>
            
            <div class="flex justify-end space-x-3 pt-4">
              <button
                type="button"
                @click="closeCompleteForm"
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700 disabled:opacity-50"
              >
                {{ submitting ? 'Completing...' : 'Complete Call' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
  
  <!-- Close Call Modal (Remarks) -->
  <StatusChangeModal
    :open="showCloseCallModal"
    :current-status="selectedCall?.status || ''"
    :status-options="[]"
    :hide-status-dropdown="true"
    :title="'Close Call'"
    @close="closeCloseCallModal"
    @change="handleCloseCall"
  />
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../../stores/auth'
import { ticketCalls } from '../../../api/tickets'
import { formatDate, formatDateTime } from '../../../utils/date'
import InlineEditDropdown from '../InlineEditDropdown.vue'
import StatusChangeModal from '../StatusChangeModal.vue'

const router = useRouter()
const auth = useAuthStore()

// Check if current user is manager (can edit call status)
const isManager = computed(() => {
  return auth.userType === 'manager'
})

const props = defineProps({
  ticketId: {
    type: [String, Number],
    required: true
  },
  ticket: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['call-scheduled', 'show-notification'])

// State
const calls = ref([])
const loading = ref(false)
const submitting = ref(false)
const error = ref('')

// Form states
const showScheduleForm = ref(false)
const editingCall = ref(null)
const showCloseCallModal = ref(false)
const selectedCall = ref(null)
const completingCall = ref(null)

// Form data
const callForm = ref({
  call_type: '',
  party_name: '',
  party_contact: '',
  scheduled_at: '',
  purpose: ''
})

// Call status options for inline editing
const callStatusOptions = [
  { value: 'Open', label: 'Open' },
  { value: 'In Progress', label: 'In Progress' },
  { value: 'Completed', label: 'Completed' }
]

const completeForm = ref({
  notes: ''
})

// Load calls
const loadCalls = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await ticketCalls.fetch(props.ticketId, { limit: 100 })
    calls.value = response.calls || []
  } catch (err) {
    console.error('Failed to load calls:', err)
    error.value = 'Failed to load calls'
  } finally {
    loading.value = false
  }
}

// Navigate to Add Call page
const navigateToAddCall = () => {
  const currentPath = router.currentRoute.value.path
  const basePath = currentPath.includes('/manager/') ? '/manager' : '/engineer'
  router.push(`${basePath}/tickets/${props.ticketId}/add-call`)
}

// Submit call form
const submitCall = async () => {
  submitting.value = true
  error.value = ''
  
  try {
    const callData = {
      ...callForm.value,
      scheduled_at: new Date(callForm.value.scheduled_at).toISOString()
    }
    
    let response
    if (editingCall.value) {
      response = await ticketCalls.update(props.ticketId, editingCall.value.id, callData)
    } else {
      response = await ticketCalls.schedule(props.ticketId, callData)
    }
    
    // Update or add call to list
    if (editingCall.value) {
      const index = calls.value.findIndex(c => c.id === editingCall.value.id)
      if (index !== -1) {
        calls.value[index] = response.call
      }
    } else {
      calls.value.unshift(response.call)
    }
    
    closeScheduleForm()
    emit('call-scheduled')
  } catch (err) {
    console.error('Failed to save call:', err)
    error.value = 'Failed to save call'
  } finally {
    submitting.value = false
  }
}

// Complete call
const completeCall = (call) => {
  completingCall.value = call
  completeForm.value.notes = call.notes || ''
  showCompleteForm.value = true
}

const submitComplete = async () => {
  submitting.value = true
  error.value = ''
  
  try {
    const response = await ticketCalls.complete(
      props.ticketId, 
      completingCall.value.id, 
      completeForm.value.notes
    )
    
    // Update call in list
    const index = calls.value.findIndex(c => c.id === completingCall.value.id)
    if (index !== -1) {
      calls.value[index] = response.call
    }
    
    closeCompleteForm()
    emit('call-scheduled')
  } catch (err) {
    console.error('Failed to complete call:', err)
    error.value = 'Failed to complete call'
  } finally {
    submitting.value = false
  }
}

// Cancel call
const cancelCall = async (call) => {
  if (!confirm(`Are you sure you want to cancel this call: ${call.subject}?`)) {
    return
  }
  
  try {
    await ticketCalls.cancel(props.ticketId, call.id)
    emit('show-notification', 'success', 'Call Cancelled', 'Call has been cancelled successfully')
    loadCalls()
    
    // Emit call-scheduled to refresh parent data
    emit('call-scheduled')
  } catch (err) {
    console.error('Failed to cancel call:', err)
    emit('show-notification', 'error', 'Cancellation Failed', 'Could not cancel call')
  }
}

// Close call modal functions
const openCloseCallModal = (call) => {
  selectedCall.value = call
  showCloseCallModal.value = true
}

const closeCloseCallModal = () => {
  showCloseCallModal.value = false
  selectedCall.value = null
}

const handleCloseCall = async (data) => {
  if (!selectedCall.value) return
  
  try {
    await ticketCalls.close(props.ticketId, selectedCall.value.id, data.remarks)
    emit('show-notification', 'success', 'Call Closed', 'Call has been closed successfully')
    closeCloseCallModal()
    loadCalls()
    
    // Emit call-scheduled to refresh parent data
    emit('call-scheduled')
  } catch (err) {
    console.error('Failed to close call:', err)
    emit('show-notification', 'error', 'Close Failed', err.response?.data?.error || 'Could not close call')
    throw err // Re-throw to keep modal open
  }
}

// Edit call
const editCall = (call) => {
  editingCall.value = call
  callForm.value = {
    call_type: call.call_type,
    party_name: call.party_name,
    party_contact: call.party_contact || '',
    scheduled_at: new Date(call.scheduled_at).toISOString().slice(0, 16),
    purpose: call.purpose || ''
  }
  showScheduleForm.value = true
}

// Form helpers
const closeScheduleForm = () => {
  showScheduleForm.value = false
  editingCall.value = null
  callForm.value = {
    call_type: '',
    party_name: '',
    party_contact: '',
    scheduled_at: '',
    purpose: ''
  }
}

const closeCompleteForm = () => {
  showCompleteForm.value = false
  completingCall.value = null
  completeForm.value.notes = ''
}

// Handle call status change
const handleCallStatusChange = async (call, newStatus) => {
  if (!call || !newStatus || call.status === newStatus) return
  
  try {
    // Update call status via API
    const updatePayload = { status: newStatus }
    const response = await ticketCalls.update(props.ticketId, call.id, updatePayload)
    
    // Update call in local list
    const index = calls.value.findIndex(c => c.id === call.id)
    if (index !== -1) {
      calls.value[index] = response.call || { ...calls.value[index], status: newStatus }
    }
    
    // Show success notification
    emit('show-notification', 'success', 'Call Status Updated', `Changed to: ${newStatus}`)
    
    // Emit call-scheduled to refresh parent data
    emit('call-scheduled')
  } catch (err) {
    console.error('Failed to update call status:', err)
    emit('show-notification', 'error', 'Update Failed', 'Could not update call status')
  }
}

// Utility functions
const getUserName = (user) => {
  if (!user) return 'Unknown User'
  return user.first_name || 'Unknown User'
}

const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const downloadAttachment = async (attachment) => {
  try {
    const response = await ticketCalls.downloadAttachment(attachment.id)
    
    // Create blob and download
    const blob = new Blob([response.data], { type: attachment.mime_type || 'application/octet-stream' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = attachment.original_filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (err) {
    console.error('Failed to download attachment:', err)
    emit('show-notification', 'error', 'Download Failed', 'Could not download attachment')
  }
}

const getCallStatusColor = (status) => {
  const colors = {
    'Open': 'bg-blue-500',
    'In Progress': 'bg-yellow-500',
    'Completed': 'bg-green-500',
    'CANCELLED': 'bg-red-500',
    'Cancelled': 'bg-red-500',
    // Legacy support
    'SCHEDULED': 'bg-blue-500',
    'Scheduled': 'bg-blue-500',
    'COMPLETED': 'bg-green-500'
  }
  return colors[status] || 'bg-gray-500'
}

const getCallStatusBadgeClass = (status) => {
  const classes = {
    'Open': 'bg-blue-100 text-blue-800',
    'In Progress': 'bg-yellow-100 text-yellow-800',
    'Completed': 'bg-green-100 text-green-800',
    'CANCELLED': 'bg-red-100 text-red-800',
    'Cancelled': 'bg-red-100 text-red-800',
    // Legacy support
    'SCHEDULED': 'bg-blue-100 text-blue-800',
    'Scheduled': 'bg-blue-100 text-blue-800',
    'COMPLETED': 'bg-green-100 text-green-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getCallTypeBadgeClass = (type) => {
  const classes = {
    'OEM': 'bg-purple-100 text-purple-800',
    'Customer': 'bg-blue-100 text-blue-800',
    'Internal': 'bg-gray-100 text-gray-800',
    'Vendor': 'bg-orange-100 text-orange-800',
    'Support': 'bg-green-100 text-green-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getDirectionBadgeClass = (direction) => {
  const classes = {
    'Inbound': 'bg-green-100 text-green-800',
    'Outbound': 'bg-blue-100 text-blue-800'
  }
  return classes[direction] || 'bg-gray-100 text-gray-800'
}

// Lifecycle
onMounted(() => {
  loadCalls()
})
</script>
