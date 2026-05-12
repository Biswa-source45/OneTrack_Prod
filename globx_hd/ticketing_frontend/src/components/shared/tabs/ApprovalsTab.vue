<template>
  <div class="min-h-[600px] flex flex-col">
    <!-- Approvals List -->
    <div class="flex-1 overflow-y-auto p-6">
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-lg font-medium text-gray-900">Approval Requests</h3>
        <button
          @click="showAddApprovalModal = true"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <span>Add Approval</span>
        </button>
      </div>

      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
      
      <div v-else-if="approvals.length === 0" class="text-center py-8 text-gray-500">
        <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-lg font-medium">No approval requests</p>
        <p class="text-sm">Create approval requests to get approvals from managers.</p>
      </div>

      <!-- Approval Items -->
      <div v-else class="space-y-4">
        <div v-for="approval in approvals" :key="approval.id" class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <div class="flex items-center space-x-2">
                  <div 
                    :class="[
                      'w-3 h-3 rounded-full',
                      getApprovalStatusColor(approval.status)
                    ]"
                  ></div>
                  <span class="font-medium text-gray-900">{{ approval.subject }}</span>
                  <span 
                    :class="[
                      'px-2 py-1 text-xs font-medium rounded-full',
                      getApprovalStatusBadgeClass(approval.status)
                    ]"
                  >
                    {{ approval.status }}
                  </span>
                </div>
              </div>
              
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-600 mb-3">
                <div>
                  <span class="font-medium">Requested by:</span>
                  {{ getUserName(approval.requester) }}
                </div>
                <div>
                  <span class="font-medium">Approver:</span>
                  {{ getUserName(approval.approver) }}
                </div>
                <div>
                  <span class="font-medium">Requested on:</span>
                  {{ formatDateTime(approval.created_at) }}
                </div>
                <div v-if="approval.status !== 'PENDING'">
                  <span class="font-medium">{{ approval.status === 'APPROVED' ? 'Approved' : 'Rejected' }} on:</span>
                  {{ formatDateTime(approval.updated_at) }}
                </div>
              </div>
              
              <!-- Show remarks for approved/rejected approvals -->
              <div v-if="approval.remarks" class="mt-2 p-3 bg-gray-50 border-l-4 border-gray-400 rounded text-sm">
                <div class="font-medium text-gray-800 mb-1">{{ approval.status === 'APPROVED' ? 'Approval' : 'Rejection' }} Remarks:</div>
                <div class="text-gray-700">{{ approval.remarks }}</div>
              </div>
            </div>
            
            <!-- Approval Actions (only for designated approver) -->
            <div v-if="canApproveReject(approval)" class="flex space-x-2 ml-4">
              <button
                @click="approveRequest(approval)"
                class="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
              >
                Approve
              </button>
              <button
                @click="rejectRequest(approval)"
                class="px-3 py-1 text-sm bg-red-600 text-white rounded hover:bg-red-700 transition-colors"
              >
                Reject
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Approval Modal -->
    <div v-if="showAddApprovalModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="closeAddApprovalModal">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md" @click.stop>
        <div class="p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900">Request Approval</h3>
            <button @click="closeAddApprovalModal" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="submitApproval" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Approver *</label>
              <select
                v-model="approvalForm.approver_id"
                class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              >
                <option value="">Select Manager</option>
                <option v-for="manager in managers" :key="manager.id" :value="manager.id">
                  {{ manager.first_name }} {{ manager.last_name }}
                </option>
              </select>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Subject *</label>
              <input
                type="text"
                v-model="approvalForm.subject"
                class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Brief description of what needs approval"
                required
              >
            </div>
            
            <div class="flex justify-end space-x-3 pt-4">
              <button
                type="button"
                @click="closeAddApprovalModal"
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                {{ submitting ? 'Requesting...' : 'Request Approval' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Approval Action Modal (Approve/Reject with Remarks) -->
    <StatusChangeModal
      :open="showApprovalActionModal"
      :current-status="actionType === 'approve' ? 'PENDING' : 'PENDING'"
      :status-options="[]"
      :pre-selected-status="actionType === 'approve' ? 'APPROVED' : 'REJECTED'"
      :hide-status-dropdown="true"
      :title="actionType === 'approve' ? 'Approve Request' : 'Reject Request'"
      @close="closeApprovalActionModal"
      @change="handleApprovalAction"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../../../stores/auth'
import { ticketApprovals, fetchUsers } from '../../../api/tickets'
import { formatDateTime } from '../../../utils/date'
import StatusChangeModal from '../StatusChangeModal.vue'

const auth = useAuthStore()

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

const emit = defineEmits(['approval-updated', 'show-notification'])

// State
const approvals = ref([])
const managers = ref([])
const loading = ref(false)
const submitting = ref(false)
const error = ref('')

// Modal states
const showAddApprovalModal = ref(false)
const showApprovalActionModal = ref(false)
const actionType = ref('') // 'approve' or 'reject'
const selectedApproval = ref(null)

// Form data
const approvalForm = ref({
  approver_id: '',
  subject: ''
})

// Get current user ID
const currentUserId = computed(() => {
  return auth.user?.id
})

// Load approvals
const loadApprovals = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await ticketApprovals.fetch(props.ticketId, { limit: 100 })
    approvals.value = response.approvals || []
  } catch (err) {
    console.error('Failed to load approvals:', err)
    error.value = 'Failed to load approvals'
  } finally {
    loading.value = false
  }
}

// Load managers for dropdown
const loadManagers = async () => {
  try {
    const response = await fetchUsers()
    // Backend returns users array directly, not wrapped in an object
    const users = Array.isArray(response) ? response : (response.users || [])
    // Filter users to show only managers (role_id = 2)
    managers.value = users.filter(user => user.role_id === 2)
  } catch (err) {
    console.error('Failed to load managers:', err)
  }
}

// Submit approval request
const submitApproval = async () => {
  submitting.value = true
  error.value = ''
  
  try {
    const response = await ticketApprovals.create(props.ticketId, {
      approver_id: parseInt(approvalForm.value.approver_id),
      subject: approvalForm.value.subject
    })
    
    // Add new approval to list
    approvals.value.unshift(response.approval)
    
    closeAddApprovalModal()
    emit('approval-updated')
    emit('show-notification', 'success', 'Approval Requested', 'Approval request has been sent successfully.')
  } catch (err) {
    console.error('Failed to create approval:', err)
    error.value = 'Failed to create approval request'
    emit('show-notification', 'error', 'Request Failed', 'Could not create approval request')
  } finally {
    submitting.value = false
  }
}

// Check if current user can approve/reject
const canApproveReject = (approval) => {
  return approval.status === 'PENDING' && approval.approver_id === currentUserId.value
}

// Approve request
const approveRequest = (approval) => {
  selectedApproval.value = approval
  actionType.value = 'approve'
  showApprovalActionModal.value = true
}

// Reject request
const rejectRequest = (approval) => {
  selectedApproval.value = approval
  actionType.value = 'reject'
  showApprovalActionModal.value = true
}

// Handle approval action (approve/reject)
const handleApprovalAction = async (changeData) => {
  if (!selectedApproval.value) return
  
  // Store action type before modal closes and clears it
  const currentActionType = actionType.value
  
  try {
    let response
    if (currentActionType === 'approve') {
      response = await ticketApprovals.approve(props.ticketId, selectedApproval.value.id, changeData.remarks)
    } else {
      response = await ticketApprovals.reject(props.ticketId, selectedApproval.value.id, changeData.remarks)
    }
    
    // Update approval in list
    const index = approvals.value.findIndex(a => a.id === selectedApproval.value.id)
    if (index !== -1) {
      approvals.value[index] = response.approval
    }
    
    closeApprovalActionModal()
    emit('approval-updated')
    
    const actionText = currentActionType === 'approve' ? 'approved' : 'rejected'
    emit('show-notification', 'success', `Request ${actionText.charAt(0).toUpperCase() + actionText.slice(1)}`, `Approval request has been ${actionText} successfully.`)
  } catch (err) {
    console.error(`Failed to ${currentActionType} approval:`, err)
    emit('show-notification', 'error', 'Action Failed', `Could not ${currentActionType} approval request`)
    throw err // Re-throw to let modal handle the error
  }
}

// Form helpers
const closeAddApprovalModal = () => {
  showAddApprovalModal.value = false
  approvalForm.value = {
    approver_id: '',
    subject: ''
  }
}

const closeApprovalActionModal = () => {
  showApprovalActionModal.value = false
  selectedApproval.value = null
  actionType.value = ''
}

// Utility functions
const getUserName = (user) => {
  if (!user) return 'Unknown User'
  return `${user.first_name} ${user.last_name || ''}`.trim()
}

const getApprovalStatusColor = (status) => {
  const colors = {
    'PENDING': 'bg-yellow-500',
    'APPROVED': 'bg-green-500',
    'REJECTED': 'bg-red-500'
  }
  return colors[status] || 'bg-gray-500'
}

const getApprovalStatusBadgeClass = (status) => {
  const classes = {
    'PENDING': 'bg-yellow-100 text-yellow-800',
    'APPROVED': 'bg-green-100 text-green-800',
    'REJECTED': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

// Lifecycle
onMounted(() => {
  loadApprovals()
  loadManagers()
})
</script>
