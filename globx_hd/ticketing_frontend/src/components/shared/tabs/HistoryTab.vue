<template>
  <div class="min-h-[600px] overflow-y-auto p-6">
    <div class="flex justify-between items-center mb-6">
      <h3 class="text-lg font-medium text-gray-900">Activity History</h3>
      <div class="flex space-x-2">
        <button
          @click="viewMode = 'timeline'"
          :class="[
            'px-3 py-1 text-sm rounded-md transition-colors',
            viewMode === 'timeline' 
              ? 'bg-blue-600 text-white' 
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          ]"
        >
          Timeline
        </button>
        <button
          @click="viewMode = 'list'"
          :class="[
            'px-3 py-1 text-sm rounded-md transition-colors',
            viewMode === 'list' 
              ? 'bg-blue-600 text-white' 
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          ]"
        >
          List
        </button>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>
    
    <div v-else-if="activities.length === 0" class="text-center py-8 text-gray-500">
      <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
      </svg>
      <p class="text-lg font-medium">No activity history</p>
      <p class="text-sm">Activity will appear here as actions are performed on this ticket.</p>
    </div>

    <!-- Timeline View -->
    <div v-else-if="viewMode === 'timeline'" class="space-y-6">
      <div v-for="(dayActivities, date) in groupedActivities" :key="date" class="relative">
        <!-- Date Header -->
        <div class="sticky top-0 bg-white z-10 pb-2">
          <div class="flex items-center">
            <div class="bg-blue-600 text-white px-3 py-1 rounded-full text-sm font-medium">
              {{ formatDateHeader(date) }}
            </div>
            <div class="flex-1 h-px bg-gray-200 ml-4"></div>
          </div>
        </div>

        <!-- Activities for this date -->
        <div class="relative pl-8 space-y-4">
          <!-- Timeline line -->
          <div class="absolute left-4 top-0 bottom-0 w-px bg-gray-200"></div>
          
          <div v-for="(activity, index) in dayActivities" :key="activity.id" class="relative">
            <!-- Timeline dot -->
            <div 
              :class="[
                'absolute -left-2 w-4 h-4 rounded-full border-2 border-white',
                getActivityColor(activity.activity_type)
              ]"
            ></div>
            
            <!-- Activity content -->
            <div class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm ml-4">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <div class="flex items-center space-x-2 mb-2">
                    <div class="flex items-center space-x-2">
                      <svg :class="['w-4 h-4', getActivityIconColor(activity.activity_type)]" fill="currentColor" viewBox="0 0 20 20">
                        <path v-html="getActivityIcon(activity.activity_type)"></path>
                      </svg>
                      <span class="font-medium text-gray-900">{{ activity.description }}</span>
                    </div>
                    <span 
                      :class="[
                        'px-2 py-1 text-xs font-medium rounded-full',
                        getActivityTypeBadgeClass(activity.activity_type)
                      ]"
                    >
                      {{ formatActivityType(activity.activity_type) }}
                    </span>
                  </div>
                  
                  <div class="text-sm text-gray-600 mb-2">
                    <span class="font-medium">{{ getUserName(activity.user) }}</span>
                    <span class="mx-2">•</span>
                    <span>{{ formatTime(activity.created_at) }}</span>
                  </div>
                  
                  <!-- Show old/new values for field changes -->
                  <div v-if="activity.old_value || activity.new_value" class="mt-2 p-2 bg-gray-50 rounded text-sm">
                    <div v-if="activity.old_value" class="text-red-600">
                      <span class="font-medium">From:</span> {{ activity.old_value }}
                    </div>
                    <div v-if="activity.new_value" class="text-green-600">
                      <span class="font-medium">To:</span> {{ activity.new_value }}
                    </div>
                  </div>
                  
                  <!-- Show remarks for status changes -->
                  <div v-if="activity.remarks" class="mt-2 p-3 bg-blue-50 border-l-4 border-blue-400 rounded text-sm">
                    <div class="font-medium text-blue-800 mb-1">Remarks:</div>
                    <div class="text-blue-700">{{ activity.remarks }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- List View -->
    <div v-else class="space-y-3">
      <div v-for="activity in activities" :key="activity.id" class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
        <div class="flex items-start space-x-3">
          <div 
            :class="[
              'w-8 h-8 rounded-full flex items-center justify-center',
              getActivityColor(activity.activity_type)
            ]"
          >
            <svg :class="['w-4 h-4 text-white']" fill="currentColor" viewBox="0 0 20 20">
              <path v-html="getActivityIcon(activity.activity_type)"></path>
            </svg>
          </div>
          
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between mb-1">
              <p class="font-medium text-gray-900">{{ activity.description }}</p>
              <span 
                :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  getActivityTypeBadgeClass(activity.activity_type)
                ]"
              >
                {{ formatActivityType(activity.activity_type) }}
              </span>
            </div>
            
            <div class="text-sm text-gray-600 mb-2">
              <span class="font-medium">{{ getUserName(activity.user) }}</span>
              <span class="mx-2">•</span>
              <span>{{ formatDateTime(activity.created_at) }}</span>
            </div>
            
            <!-- Show old/new values for field changes -->
            <div v-if="activity.old_value || activity.new_value" class="mt-2 p-2 bg-gray-50 rounded text-sm">
              <div v-if="activity.old_value" class="text-red-600">
                <span class="font-medium">From:</span> {{ activity.old_value }}
              </div>
              <div v-if="activity.new_value" class="text-green-600">
                <span class="font-medium">To:</span> {{ activity.new_value }}
              </div>
            </div>
            
            <!-- Show remarks for status changes -->
            <div v-if="activity.remarks" class="mt-2 p-3 bg-blue-50 border-l-4 border-blue-400 rounded text-sm">
              <div class="font-medium text-blue-800 mb-1">Remarks:</div>
              <div class="text-blue-700">{{ activity.remarks }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Load More Button -->
    <div v-if="hasMore && !loading" class="text-center mt-6">
      <button
        @click="loadMore"
        class="px-4 py-2 text-sm text-blue-600 hover:text-blue-800 transition-colors"
      >
        Load More Activities
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { ticketActivities } from '../../../api/tickets'
import { formatDate, formatDateTime, formatTime } from '../../../utils/date'
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

// State
const activities = ref([])
const loading = ref(false)
const error = ref('')
const viewMode = ref('timeline')
const hasMore = ref(true)
const currentOffset = ref(0)
const limit = 50

// Load activities
const loadActivities = async (append = false) => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await ticketActivities.fetch(props.ticketId, {
      limit,
      offset: append ? currentOffset.value : 0
    })
    
    if (append) {
      activities.value.push(...(response.activities || []))
    } else {
      activities.value = response.activities || []
      currentOffset.value = 0
    }
    
    currentOffset.value += (response.activities || []).length
    hasMore.value = (response.activities || []).length === limit
  } catch (err) {
    console.error('Failed to load activities:', err)
    error.value = 'Failed to load activities'
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  loadActivities(true)
}

// Group activities by date for timeline view
const groupedActivities = computed(() => {
  const groups = {}
  
  activities.value.forEach(activity => {
    const date = new Date(activity.created_at).toDateString()
    if (!groups[date]) {
      groups[date] = []
    }
    groups[date].push(activity)
  })
  
  return groups
})

// Utility functions
const getUserName = (user) => {
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

const getActivityIconColor = (type) => {
  return 'text-gray-600'
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
  loadActivities()
})

// Watch for ticket changes to refresh activities
watch(() => props.ticket?.updated_at, (newValue, oldValue) => {
  if (newValue && newValue !== oldValue) {
    loadActivities()
  }
}, { deep: true })
</script>
