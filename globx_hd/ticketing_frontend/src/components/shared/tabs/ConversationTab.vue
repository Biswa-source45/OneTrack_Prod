<template>
  <div class="min-h-[600px] flex flex-col">
    <!-- Header with Triple-dot Menu -->
    <div class="flex justify-between items-center p-6 border-b border-gray-200">
      <h3 class="text-lg font-medium text-gray-900">Conversation & Activity</h3>
      
      <!-- Triple-dot Action Menu (only for authorized users) -->
      <div v-if="canAddComments" class="relative">
        <button
          @click="showActionMenu = !showActionMenu"
          class="p-2 rounded-full hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
          title="More actions"
        >
          <svg class="w-5 h-5 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
            <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
          </svg>
        </button>
        
        <!-- Dropdown Menu -->
        <div
          v-if="showActionMenu"
          class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none z-50"
          @click.stop
        >
          <div class="py-1">
            <button
              @click="handleAddComment"
              class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
            >
              <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.955 8.955 0 01-4.126-.98L3 20l1.98-5.874A8.955 8.955 0 013 12c0-4.418 3.582-8 8-8s8 3.582 8 8z" />
              </svg>
              Add Comment
            </button>
            <button
              @click="handleAddResolution"
              class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
            >
              <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Add Resolution
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Comments List -->
    <div class="flex-1 overflow-y-auto p-6 space-y-4">
      <!-- Error Message -->
      <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
        <div class="flex items-center">
          <svg class="w-5 h-5 text-red-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-sm text-red-700">{{ error }}</p>
          <button @click="error = ''" class="ml-auto text-red-400 hover:text-red-600">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
      
      <div v-else>
        <!-- Original Ticket Details (First Message) -->
        <div v-if="ticket && ticket.ticket_details" class="bg-blue-50 border border-blue-200 rounded-lg p-4 shadow-sm">
          <div class="flex items-start space-x-3">
            <!-- User Avatar -->
            <div class="flex-shrink-0">
              <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
                {{ getUserInitials(ticket.contact) }}
              </div>
            </div>
            
            <!-- Original Issue Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center space-x-2">
                  <p class="text-sm font-medium text-gray-900">{{ getContactName(ticket.contact) }}</p>
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                    Original Issue
                  </span>
                </div>
                <span class="text-xs text-gray-500">{{ formatDate(ticket.created_at) }}</span>
              </div>
              
              <div class="prose prose-sm max-w-none">
                <p class="text-gray-700 whitespace-pre-wrap">{{ ticket.ticket_details }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- No Data Message -->
        <div v-if="unifiedTimeline.length === 0" class="text-center py-8 text-gray-500">
          <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.955 8.955 0 01-4.126-.98L3 20l1.98-5.874A8.955 8.955 0 013 12c0-4.418 3.582-8 8-8s8 3.582 8 8z" />
          </svg>
          <p class="text-lg font-medium">No activity yet</p>
          <p class="text-sm">Use the menu above to add comments or resolutions.</p>
        </div>

        <!-- Unified Timeline (Comments + Activities) -->
        <div v-else class="space-y-6">
          <div v-for="(dayItems, date) in groupedTimeline" :key="date">
            <!-- Date Header -->
            <div class="sticky top-0 bg-gray-50 z-10 pb-2 mb-4">
              <div class="flex items-center">
                <div class="bg-blue-600 text-white px-3 py-1 rounded-full text-sm font-medium">
                  {{ formatDateHeader(date) }}
                </div>
                <div class="flex-1 h-px bg-gray-200 ml-4"></div>
              </div>
            </div>

            <!-- Timeline Items for this date -->
            <div class="space-y-4">
              <!-- Comment Item -->
              <div v-for="item in dayItems" :key="item.itemType + '-' + item.id">
                <div v-if="item.itemType === 'comment'" class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
                  <div class="flex items-start space-x-3">
                    <!-- User Avatar -->
                    <div class="flex-shrink-0">
                      <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
                        {{ getUserInitials(item.user) }}
                      </div>
                    </div>
                    
                    <!-- Comment Content -->
                    <div class="flex-1 min-w-0">
                      <div class="flex items-center justify-between mb-2">
                        <div class="flex items-center space-x-2">
                          <p class="text-sm font-medium text-gray-900">{{ getUserName(item.user) }}</p>
                          <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium" 
                                :class="item.type === 'resolution' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'">
                            {{ item.type === 'resolution' ? 'Resolution' : 'Comment' }}
                          </span>
                          <span v-if="item.is_internal" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-yellow-100 text-yellow-800">
                            Internal
                          </span>
                        </div>
                        <div class="flex items-center space-x-2">
                          <span class="text-xs text-gray-500">{{ formatTime(item.created_at) }}</span>
                          <div class="flex space-x-1">
                            <button
                              v-if="canEditComment(item)"
                              @click="editComment(item)"
                              class="text-gray-400 hover:text-gray-600 p-1 rounded"
                              title="Edit comment"
                            >
                              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                              </svg>
                            </button>
                            <button
                              v-if="canEditComment(item)"
                              @click="deleteComment(item)"
                              class="text-gray-400 hover:text-red-600 p-1 rounded"
                              title="Delete comment"
                            >
                              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </div>
                      
                      <!-- Edit Mode -->
                      <div v-if="editingComment && editingComment.id === item.id" class="space-y-2">
                        <textarea
                          v-model="editContent"
                          class="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                          rows="3"
                        ></textarea>
                        <div class="flex space-x-2">
                          <button
                            @click="cancelEdit"
                            class="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 transition-colors"
                          >
                            Cancel
                          </button>
                          <button
                            @click="saveEdit"
                            :disabled="!editContent.trim() || saving"
                            class="px-4 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                          >
                            {{ saving ? 'Saving...' : 'Save' }}
                          </button>
                        </div>
                      </div>
                      
                      <div v-else class="prose prose-sm max-w-none">
                        <p class="text-gray-700 whitespace-pre-wrap">{{ item.content }}</p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Activity Item -->
                <div v-else-if="item.itemType === 'activity'" class="bg-white border-l-4 rounded-lg p-4 shadow-sm" :class="getActivityBorderClass(item.activity_type)">
                  <div class="flex items-start space-x-3">
                    <!-- Activity Icon -->
                    <div class="flex-shrink-0">
                      <div class="w-8 h-8 rounded-full flex items-center justify-center" :class="getActivityColor(item.activity_type)">
                        <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                          <path v-html="getActivityIcon(item.activity_type)"></path>
                        </svg>
                      </div>
                    </div>
                    
                    <!-- Activity Content -->
                    <div class="flex-1 min-w-0">
                      <div class="flex items-center justify-between mb-2">
                        <div class="flex items-center space-x-2">
                          <p class="text-sm font-medium text-gray-900">{{ item.description }}</p>
                          <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium" :class="getActivityTypeBadgeClass(item.activity_type)">
                            {{ formatActivityType(item.activity_type) }}
                          </span>
                        </div>
                        <span class="text-xs text-gray-500">{{ formatTime(item.created_at) }}</span>
                      </div>
                      
                      <div class="text-sm text-gray-600 mb-2">
                        <span class="font-medium">{{ getActivityUserName(item.user) }}</span>
                      </div>
                      
                      <!-- Show old/new values for field changes -->
                      <div v-if="item.old_value || item.new_value" class="mt-2 p-2 bg-gray-50 rounded text-sm">
                        <div v-if="item.old_value" class="text-red-600">
                          <span class="font-medium">From:</span> {{ item.old_value }}
                        </div>
                        <div v-if="item.new_value" class="text-green-600">
                          <span class="font-medium">To:</span> {{ item.new_value }}
                        </div>
                      </div>
                      
                      <!-- Show remarks for status changes -->
                      <div v-if="item.remarks" class="mt-2 p-3 bg-blue-50 border-l-4 border-blue-400 rounded text-sm">
                        <div class="font-medium text-blue-800 mb-1">Remarks:</div>
                        <div class="text-blue-700">{{ item.remarks }}</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Add Comment/Resolution Form -->
      <div v-if="showAddForm" class="bg-gray-50 border border-gray-200 rounded-lg p-4 shadow-sm">
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h4 class="text-sm font-medium text-gray-900">
              {{ addFormType === 'resolution' ? 'Add Resolution' : 'Add Comment' }}
            </h4>
            <button
              @click="cancelAddForm"
              class="text-gray-400 hover:text-gray-600 p-1 rounded"
              title="Cancel"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <textarea
            v-model="newEntryContent"
            class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            :placeholder="`Enter your ${addFormType}...`"
            required
          ></textarea>
          
          <div class="flex items-center justify-end">
            <div class="flex space-x-2">
              <button
                @click="cancelAddForm"
                class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                @click="submitNewEntry"
                :disabled="!newEntryContent.trim() || submitting"
                class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {{ submitting ? 'Adding...' : `Add ${addFormType === 'resolution' ? 'Resolution' : 'Comment'}` }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Close Ticket Section (Only for Managers) -->
    <div v-if="canCloseTicket" class="border-t border-gray-200 p-6 bg-gray-50">
      <div class="flex items-center justify-end">
        <button
          @click="closeTicket"
          :class="[
            'px-6 py-2 text-white rounded-lg transition-colors',
            ticket?.ticket_status === 'CLOSED' 
              ? 'bg-green-600 hover:bg-green-700' 
              : 'bg-gray-600 hover:bg-gray-700'
          ]"
        >
          {{ ticket?.ticket_status === 'CLOSED' ? 'Reopen Ticket' : 'Close Ticket' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { ticketComments, ticketActivities } from '../../../api/tickets'
import { formatDate, formatDateTime, formatTime } from '../../../utils/date'
import { useAuthStore } from '../../../stores/auth'
import { decodeJWT } from '../../../utils/jwt'
import { formatUserName } from '../../../utils/user'

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

const emit = defineEmits(['comment-added', 'close-ticket'])

// Auth store
const auth = useAuthStore()

// Check if current user can add comments/resolutions
const canAddComments = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return false
  
  const claims = decodeJWT(token)
  return claims?.type === 'user'
})

// Check if current user can close/reopen tickets (only managers, not engineers)
const canCloseTicket = computed(() => {
  return auth.userType === 'manager'
})

// State
const comments = ref([])
const activities = ref([])
const loading = ref(false)
const submitting = ref(false)
const saving = ref(false)
const showActionMenu = ref(false)
const showAddForm = ref(false)
const addFormType = ref('comment') // 'comment' or 'resolution'
const newEntryContent = ref('')
const isInternal = ref(false)
const error = ref('')

// New comment form
const newComment = ref({
  content: '',
  type: 'comment',
  is_internal: false
})

// Edit comment state
const editingComment = ref(null)
const editContent = ref('')

// Computed properties
const sortedComments = computed(() => {
  return [...comments.value].sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
})

// Unified timeline combining comments and activities
const unifiedTimeline = computed(() => {
  const timeline = []
  
  // Add comments with type marker
  comments.value.forEach(comment => {
    timeline.push({
      ...comment,
      itemType: 'comment',
      timestamp: new Date(comment.created_at)
    })
  })
  
  // Add activities with type marker
  activities.value.forEach(activity => {
    timeline.push({
      ...activity,
      itemType: 'activity',
      timestamp: new Date(activity.created_at)
    })
  })
  
  // Sort by timestamp (oldest first)
  return timeline.sort((a, b) => a.timestamp - b.timestamp)
})

// Group timeline by date for better organization
const groupedTimeline = computed(() => {
  const groups = {}
  
  unifiedTimeline.value.forEach(item => {
    const date = item.timestamp.toDateString()
    if (!groups[date]) {
      groups[date] = []
    }
    groups[date].push(item)
  })
  
  return groups
})

// Load comments
const loadComments = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await ticketComments.fetch(props.ticketId, { limit: 100 })
    comments.value = response.comments || []
  } catch (err) {
    console.error('Failed to load comments:', err)
    error.value = 'Failed to load comments'
  } finally {
    loading.value = false
  }
}

// Load activities
const loadActivities = async () => {
  try {
    const response = await ticketActivities.fetch(props.ticketId, { limit: 100 })
    activities.value = response.activities || []
  } catch (err) {
    console.error('Failed to load activities:', err)
    // Don't set error here as it's not critical - activities are supplementary
  }
}

// Load both comments and activities
const loadAllData = async () => {
  loading.value = true
  error.value = ''
  
  try {
    await Promise.all([loadComments(), loadActivities()])
  } catch (err) {
    console.error('Failed to load data:', err)
    error.value = 'Failed to load conversation data'
  } finally {
    loading.value = false
  }
}

// Action menu handlers
const handleAddComment = () => {
  showActionMenu.value = false
  
  if (!canAddComments.value) {
    error.value = 'Only standard users may add comments or resolutions'
    return
  }
  
  addFormType.value = 'comment'
  showAddForm.value = true
  newEntryContent.value = ''
  isInternal.value = false
}

const handleAddResolution = () => {
  showActionMenu.value = false
  
  if (!canAddComments.value) {
    error.value = 'Only standard users may add comments or resolutions'
    return
  }
  
  addFormType.value = 'resolution'
  showAddForm.value = true
  newEntryContent.value = ''
  isInternal.value = false
}

// Click outside handler for action menu
const handleClickOutside = (event) => {
  if (showActionMenu.value && !event.target.closest('.relative')) {
    showActionMenu.value = false
  }
}

// Add form handlers
const cancelAddForm = () => {
  showAddForm.value = false
  newEntryContent.value = ''
  isInternal.value = false
  addFormType.value = 'comment'
}

const submitNewEntry = async () => {
  if (!newEntryContent.value.trim()) return
  
  // Check permissions before submitting
  if (!canAddComments.value) {
    error.value = 'Only standard users may add comments or resolutions'
    return
  }
  
  submitting.value = true
  error.value = ''
  
  try {
    const response = await ticketComments.create(props.ticketId, {
      content: newEntryContent.value.trim(),
      type: addFormType.value,
      is_internal: isInternal.value
    })
    
    // Add new comment to list
    comments.value.push(response.comment)
    
    // Reset form
    cancelAddForm()
    
    // Reload activities to capture the new COMMENT_ADDED activity
    await loadActivities()
    
    // Emit event to parent
    emit('comment-added')
  } catch (err) {
    console.error('Failed to add entry:', err)
    
    // Handle specific error cases
    if (err.response?.status === 401) {
      error.value = 'Unauthorized: Please check your login status and permissions'
    } else {
      error.value = `Failed to add ${addFormType.value}: ${err.response?.data?.error || err.message}`
    }
  } finally {
    submitting.value = false
  }
}

// Close ticket function
const closeTicket = async () => {
  // Emit event to parent to handle status change using the same mechanism as inline editor
  emit('close-ticket')
}

// Edit comment
const editComment = (comment) => {
  editingComment.value = comment
  editContent.value = comment.content
}

const cancelEdit = () => {
  editingComment.value = null
  editContent.value = ''
}

const saveEdit = async () => {
  if (!editContent.value.trim()) return
  
  saving.value = true
  
  try {
    const response = await ticketComments.update(
      props.ticketId, 
      editingComment.value.id, 
      { content: editContent.value.trim() }
    )
    
    // Update comment in list
    const index = comments.value.findIndex(c => c.id === editingComment.value.id)
    if (index !== -1) {
      comments.value[index] = response.comment
    }
    
    cancelEdit()
  } catch (err) {
    console.error('Failed to update comment:', err)
    error.value = 'Failed to update comment'
  } finally {
    saving.value = false
  }
}

// Delete comment
const deleteComment = async (comment) => {
  if (!confirm('Are you sure you want to delete this comment?')) return
  
  try {
    await ticketComments.delete(props.ticketId, comment.id)
    
    // Remove from list
    const index = comments.value.findIndex(c => c.id === comment.id)
    if (index !== -1) {
      comments.value.splice(index, 1)
    }
    
    emit('comment-added') // Refresh counts
  } catch (err) {
    console.error('Failed to delete comment:', err)
    error.value = 'Failed to delete comment'
  }
}

// Utility functions
const getUserName = (user) => {
  if (!user) return 'Unknown User'
  
  const firstName = user.first_name || ''
  const lastName = user.last_name || ''
  
  // Check if first_name already contains both names (common data issue)
  if (firstName.includes(' ') || !lastName) {
    return firstName.trim()  // Use only first_name if it already contains full name
  }
  
  return `${firstName} ${lastName}`.trim()
}

const getUserInitials = (user) => {
  if (!user) return '?'
  const firstName = user.first_name || ''
  const lastName = user.last_name || ''
  return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase() || '?'
}

const canEditComment = (comment) => {
  // For now, allow editing own comments (would need user context in real app)
  return true // This should check if current user owns the comment
}

const getContactName = (contact) => {
  if (!contact) return 'Unknown User'
  return `${contact.first_name} ${contact.last_name || ''}`.trim()
}

// Activity-specific utility functions
const getActivityUserName = (user) => {
  if (!user) return 'System'
  return formatUserName(user)
}

const formatActivityType = (type) => {
  return type.replace(/_/g, ' ').toLowerCase().replace(/\b\w/g, l => l.toUpperCase())
}

const formatDateHeader = (dateString) => {
  const date = new Date(dateString)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)
  
  if (date.toDateString() === today.toDateString()) {
    return 'Today'
  } else if (date.toDateString() === yesterday.toDateString()) {
    return 'Yesterday'
  } else {
    return formatDate(date)
  }
}

const getActivityColor = (type) => {
  const colors = {
    'TICKET_CREATED': 'bg-blue-500',
    'STATUS_CHANGED': 'bg-yellow-500',
    'ASSIGNED': 'bg-green-500',
    'UNASSIGNED': 'bg-red-500',
    'PRIORITY_CHANGED': 'bg-orange-500',
    'COMMENT_ADDED': 'bg-blue-400',
    'RESOLUTION_ADDED': 'bg-green-400',
    'CALL_SCHEDULED': 'bg-purple-500',
    'CALL_COMPLETED': 'bg-purple-600',
    'CALL_CANCELLED': 'bg-red-400',
    'APPROVAL_REQUESTED': 'bg-amber-500',
    'APPROVAL_APPROVED': 'bg-emerald-500',
    'APPROVAL_REJECTED': 'bg-rose-500',
    'TICKET_UPDATED': 'bg-gray-500',
    'PRODUCT_CHANGED': 'bg-indigo-500',
    'SUBJECT_CHANGED': 'bg-teal-500'
  }
  return colors[type] || 'bg-gray-400'
}

const getActivityBorderClass = (type) => {
  const colors = {
    'TICKET_CREATED': 'border-blue-500',
    'STATUS_CHANGED': 'border-yellow-500',
    'ASSIGNED': 'border-green-500',
    'UNASSIGNED': 'border-red-500',
    'PRIORITY_CHANGED': 'border-orange-500',
    'COMMENT_ADDED': 'border-blue-400',
    'RESOLUTION_ADDED': 'border-green-400',
    'CALL_SCHEDULED': 'border-purple-500',
    'CALL_COMPLETED': 'border-purple-600',
    'CALL_CANCELLED': 'border-red-400',
    'APPROVAL_REQUESTED': 'border-amber-500',
    'APPROVAL_APPROVED': 'border-emerald-500',
    'APPROVAL_REJECTED': 'border-rose-500',
    'TICKET_UPDATED': 'border-gray-500',
    'PRODUCT_CHANGED': 'border-indigo-500',
    'SUBJECT_CHANGED': 'border-teal-500'
  }
  return colors[type] || 'border-gray-400'
}

const getActivityTypeBadgeClass = (type) => {
  const classes = {
    'TICKET_CREATED': 'bg-blue-100 text-blue-800',
    'STATUS_CHANGED': 'bg-yellow-100 text-yellow-800',
    'ASSIGNED': 'bg-green-100 text-green-800',
    'UNASSIGNED': 'bg-red-100 text-red-800',
    'PRIORITY_CHANGED': 'bg-orange-100 text-orange-800',
    'COMMENT_ADDED': 'bg-blue-100 text-blue-800',
    'RESOLUTION_ADDED': 'bg-green-100 text-green-800',
    'CALL_SCHEDULED': 'bg-purple-100 text-purple-800',
    'CALL_COMPLETED': 'bg-purple-100 text-purple-800',
    'CALL_CANCELLED': 'bg-red-100 text-red-800',
    'APPROVAL_REQUESTED': 'bg-amber-100 text-amber-800',
    'APPROVAL_APPROVED': 'bg-emerald-100 text-emerald-800',
    'APPROVAL_REJECTED': 'bg-rose-100 text-rose-800',
    'TICKET_UPDATED': 'bg-gray-100 text-gray-800',
    'PRODUCT_CHANGED': 'bg-indigo-100 text-indigo-800',
    'SUBJECT_CHANGED': 'bg-teal-100 text-teal-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getActivityIcon = (type) => {
  const icons = {
    'TICKET_CREATED': 'M12 6v6m0 0v6m0-6h6m-6 0H6',
    'STATUS_CHANGED': 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15',
    'ASSIGNED': 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z',
    'UNASSIGNED': 'M13 7a4 4 0 11-8 0 4 4 0 018 0zM9 14a6 6 0 00-6 6v1h12v-1a6 6 0 00-6-6zM21 12h-6',
    'PRIORITY_CHANGED': 'M7 11l5-6 5 6M7 21l5-6 5 6',
    'COMMENT_ADDED': 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.955 8.955 0 01-4.126-.98L3 20l1.98-5.874A8.955 8.955 0 013 12c0-4.418 3.582-8 8-8s8 3.582 8 8z',
    'RESOLUTION_ADDED': 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    'CALL_SCHEDULED': 'M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z',
    'CALL_COMPLETED': 'M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z',
    'CALL_CANCELLED': 'M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728L5.636 5.636m12.728 12.728L18.364 5.636M5.636 18.364l12.728-12.728',
    'APPROVAL_REQUESTED': 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z M15 12a3 3 0 11-6 0 3 3 0 016 0z',
    'APPROVAL_APPROVED': 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    'APPROVAL_REJECTED': 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
    'TICKET_UPDATED': 'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z',
    'PRODUCT_CHANGED': 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4',
    'SUBJECT_CHANGED': 'M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z'
  }
  return icons[type] || 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
}

// Lifecycle
onMounted(() => {
  loadAllData()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.prose p {
  margin: 0;
}
</style>
