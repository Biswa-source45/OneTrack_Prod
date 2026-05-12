<template>
  <div class="min-h-[600px] flex flex-col">
    <!-- Calls List -->
    <div class="flex-1 overflow-y-auto p-6">
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-lg font-medium text-gray-900">Call Logs</h3>
        <!-- No Add Call button for customers -->
      </div>

      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
      
      <div v-else-if="calls.length === 0" class="text-center py-8 text-gray-500">
        <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
        </svg>
        <p class="text-lg font-medium">No calls logged</p>
        <p class="text-sm">Call logs for inbound and outbound communications will appear here.</p>
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
                  <span 
                    :class="[
                      'px-2 py-1 text-xs font-medium rounded-full ml-1',
                      getCallStatusBadgeClass(call.status)
                    ]"
                  >
                    {{ call.status }}
                  </span>
                </div>
              </div>
              
              <div v-if="call.description" class="mb-3">
                <span class="font-medium text-gray-700">Description:</span>
                <p class="text-gray-600 mt-1 whitespace-pre-wrap">{{ call.description }}</p>
              </div>
            </div>
            
            <!-- No Call Actions for customers -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { fetchTicketCalls } from '../../../api/tickets'
import { formatDateTime } from '../../../utils/date'

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
const calls = ref([])
const loading = ref(false)
const error = ref('')

// Helper functions (matching manager implementation exactly)
const getUserName = (user) => {
  if (!user) return 'Unknown User'
  return user.first_name || 'Unknown User'
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

// Load calls
const loadCalls = async () => {
  if (!props.ticketId) return
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetchTicketCalls(props.ticketId)
    calls.value = response.calls || []
  } catch (err) {
    console.error('Failed to load calls:', err)
    error.value = 'Failed to load calls'
  } finally {
    loading.value = false
  }
}

// Auto-refresh calls every 30 seconds
let refreshInterval = null

onMounted(() => {
  loadCalls()
  
  // Set up auto-refresh
  refreshInterval = setInterval(() => {
    loadCalls()
  }, 30000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
