<template>
  <div class="min-h-[600px] flex flex-col">
    <!-- Header (Read-only, no action menu) -->
    <div class="flex justify-between items-center p-6 border-b border-gray-200">
      <h3 class="text-lg font-medium text-gray-900">Conversations</h3>
    </div>

    <!-- Comments List (Original Issue + Resolutions Only) -->
    <div class="flex-1 overflow-y-auto p-6 space-y-4">
      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <div v-else-if="error" class="text-center py-8 text-red-600">
        <p>{{ error }}</p>
        <button @click="loadComments" class="mt-2 text-blue-600 hover:text-blue-800">
          Try again
        </button>
      </div>

      <div v-else-if="!filteredComments.length" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.955 8.955 0 01-4.126-.98L3 20l1.98-5.874A8.955 8.955 0 013 12c0-4.418 3.582-8 8-8s8 3.582 8 8z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No conversations yet</h3>
        <p class="mt-1 text-sm text-gray-500">Updates and resolutions will appear here.</p>
      </div>

      <!-- Original Issue (First Message) -->
      <div v-if="ticket && ticket.ticket_details" class="bg-blue-50 border border-blue-200 rounded-lg p-4 shadow-sm mb-6">
        <div class="flex items-start space-x-3">
          <!-- User Avatar -->
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
              {{ getContactInitials(ticket.contact) }}
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

      <!-- Resolution Comments Only -->
      <div v-for="comment in filteredComments" :key="comment.id" class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
        <div class="flex items-start space-x-3">
          <!-- User Avatar -->
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
              {{ getUserInitials(comment.user) }}
            </div>
          </div>
          
          <!-- Comment Content -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center space-x-2">
                <p class="text-sm font-medium text-gray-900">{{ getUserName(comment.user) }}</p>
                <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                  Resolution
                </span>
              </div>
              <span class="text-xs text-gray-500">{{ formatDate(comment.created_at) }}</span>
            </div>
            
            <div class="prose prose-sm max-w-none">
              <p class="text-gray-700 whitespace-pre-wrap">{{ comment.content }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { fetchTicketComments } from '../../../api/tickets'
import { formatDateTime } from '../../../utils/date'
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

// Reactive state
const comments = ref([])
const loading = ref(false)
const error = ref('')

// Format date helper
const formatDate = formatDateTime

// Filter comments to show only resolutions (type === 'resolution')
const filteredComments = computed(() => {
  return comments.value.filter(comment => comment.type === 'resolution')
})

// Helper functions
const getUserName = (user) => {
  if (!user) return 'Unknown User'
  return formatUserName(user)
}

const getUserInitials = (user) => {
  if (!user) return 'U'
  const firstName = user.first_name || ''
  const lastName = user.last_name || ''
  return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase() || 'U'
}

const getContactName = (contact) => {
  if (!contact) return 'Unknown Contact'
  return `${contact.first_name} ${contact.last_name || ''}`.trim()
}

const getContactInitials = (contact) => {
  if (!contact) return 'C'
  const firstName = contact.first_name || ''
  const lastName = contact.last_name || ''
  return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase() || 'C'
}

// Load comments
const loadComments = async () => {
  if (!props.ticketId) return
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetchTicketComments(props.ticketId)
    comments.value = response.comments || []
  } catch (err) {
    console.error('Failed to load comments:', err)
    error.value = 'Failed to load conversations'
  } finally {
    loading.value = false
  }
}

// Auto-refresh comments every 30 seconds
let refreshInterval = null

onMounted(() => {
  loadComments()
  
  // Set up auto-refresh
  refreshInterval = setInterval(() => {
    loadComments()
  }, 30000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
