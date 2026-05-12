<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <button
              @click="goBack"
              class="p-2 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
            </button>
            <div>
              <h1 class="text-2xl font-bold text-gray-900">Task Details</h1>
              <p v-if="task" class="text-gray-600 mt-1">{{ task.subject }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
        <div class="flex">
          <svg class="w-5 h-5 text-red-400 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error loading task</h3>
            <p class="text-sm text-red-700 mt-1">{{ error }}</p>
          </div>
        </div>
        <div class="mt-4">
          <button
            @click="loadTaskData"
            class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 transition-colors text-sm"
          >
            Try Again
          </button>
        </div>
      </div>

      <!-- Main Content -->
      <div v-else-if="task" class="flex gap-8">
        <!-- Left Sidebar -->
        <div class="w-80 space-y-6">
          <!-- Task Properties -->
          <PropertySection 
            title="Task Properties" 
            :show-edit-icon="true"
            @edit="openEditModal"
          >
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-600">Subject</span>
                <span class="text-sm font-medium text-gray-900">{{ task.subject || 'No Subject' }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-600">Created By</span>
                <span class="text-sm font-medium text-gray-900">{{ getCreatorName(task) }}</span>
              </div>
            </div>
          </PropertySection>

          <!-- Key Information -->
          <PropertySection title="Key Information">
            <div class="space-y-3">
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Status</span>
                <div class="mt-1">
                  <InlineEditDropdown
                    :model-value="task.task_status"
                    :options="statusOptions"
                    option-value="value"
                    option-label="label"
                    empty-label="Select Status"
                    :allow-empty="false"
                    :display-class="`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusBadgeClass(task.task_status)}`"
                    @change="handleStatusChange"
                  />
                </div>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Priority</span>
                <div class="mt-1">
                  <InlineEditDropdown
                    :model-value="task.priority"
                    :options="priorityOptions"
                    option-value="value"
                    option-label="label"
                    empty-label="Select Priority"
                    :allow-empty="false"
                    :display-class="`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityBadgeClass(task.priority)}`"
                    @change="handlePriorityChange"
                  />
                </div>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Due Date</span>
                <div class="mt-1">
                  <InlineEditDropdown
                    :model-value="task.due_date && task.due_date !== null ? formatDateForInput(task.due_date) : ''"
                    :options="[]"
                    empty-label="No due date"
                    :allow-empty="true"
                    display-class="text-sm font-medium text-gray-900"
                    @change="handleDueDateChange"
                    :is-date-picker="true"
                  />
                </div>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Created By</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ getCreatorName(task) }}</p>
              </div>
            </div>
          </PropertySection>
        </div>

        <!-- Right Content Area -->
        <div class="flex-1 bg-white rounded-lg shadow-sm">
          <!-- Comments Section -->
          <TaskCommentsTab
            :task-id="taskId"
            :task="task"
            @comment-added="handleCommentAdded"
            @mark-completed="handleMarkCompleted"
          />
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchTaskFullDetails, updateTask } from '../../api/tasks'
import { fetchUsers } from '../../api/tickets'
import { formatDateTime } from '../../utils/date'
import { formatUserName } from '../../utils/user'
import PropertySection from './PropertySection.vue'
import InlineEditDropdown from './InlineEditDropdown.vue'
import NotificationToast from './NotificationToast.vue'
import TaskCommentsTab from './tabs/TaskCommentsTab.vue'

const route = useRoute()
const router = useRouter()

const taskId = computed(() => route.params.id)
const task = ref(null)
const loading = ref(false)
const error = ref('')

// Data for inline editing
const users = ref([])
const usersLoading = ref(false)
const statusOptions = [
  { value: 'Not Started', label: 'Not Started' },
  { value: 'In Progress', label: 'In Progress' },
  { value: 'Completed', label: 'Completed' },
  { value: 'Deferred', label: 'Deferred' },
  { value: 'Waiting on someone else', label: 'Waiting on someone else' },
  { value: 'Canceled', label: 'Canceled' }
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

// Load task data
const loadTaskData = async () => {
  if (!taskId.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetchTaskFullDetails(taskId.value)
    task.value = response.task
    
    // Update counts if available
    if (response.counts) {
      task.value.counts = response.counts
    }
  } catch (err) {
    console.error('Failed to load task details:', err)
    error.value = 'Failed to load task details. Please try again.'
  } finally {
    loading.value = false
  }
}

// Load users for task owner dropdown
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
  router.go(-1)
}

// Event handlers
const handleCommentAdded = () => {
  loadTaskData()
}

const handleMarkCompleted = () => {
  handleStatusChange('Completed')
}

const openEditModal = () => {
  router.push(`/manager/tasks/${taskId.value}/edit`)
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
const handleStatusChange = async (newStatus) => {
  try {
    const updatePayload = {
      task_status: newStatus
    }
    await updateTask(task.value.id, updatePayload)
    
    // Show success notification
    showNotification('success', 'Status Updated', `Changed to: ${newStatus}`)
    
    // Refresh task data to update display and history
    await loadTaskData()
  } catch (err) {
    console.error('Failed to update task status:', err)
    showNotification('error', 'Update Failed', 'Could not update task status')
  }
}

const handlePriorityChange = async (newPriority) => {
  try {
    const updatePayload = {
      priority: newPriority
    }
    await updateTask(task.value.id, updatePayload)
    
    // Show success notification
    showNotification('success', 'Priority Updated', `Changed to: ${newPriority}`)
    
    // Refresh task data to update display and history
    await loadTaskData()
  } catch (err) {
    console.error('Failed to update task priority:', err)
    showNotification('error', 'Update Failed', 'Could not update task priority')
  }
}

const handleDueDateChange = async (newDueDate) => {
  try {
    const updatePayload = {
      due_date: newDueDate || null
    }
    await updateTask(task.value.id, updatePayload)
    
    // Show success notification
    const displayDate = newDueDate ? formatDateTime(newDueDate) : 'No due date'
    showNotification('success', 'Due Date Updated', `Changed to: ${displayDate}`)
    
    // Refresh task data to update display and history
    await loadTaskData()
  } catch (err) {
    console.error('Failed to update task due date:', err)
    showNotification('error', 'Update Failed', 'Could not update task due date')
  }
}

// Utility functions
const getStatusBadgeClass = (status) => {
  const classes = {
    'Not Started': 'bg-gray-100 text-gray-800',
    'In Progress': 'bg-blue-100 text-blue-800',
    'Completed': 'bg-green-100 text-green-800',
    'Deferred': 'bg-yellow-100 text-yellow-800',
    'Waiting on someone else': 'bg-purple-100 text-purple-800',
    'Canceled': 'bg-red-100 text-red-800'
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

const getCreatorName = (task) => {
  if (!task?.creator) return 'Unknown'
  return formatUserName(task.creator)
}

const formatDateForInput = (dateStr) => {
  if (!dateStr || dateStr === null || dateStr === undefined) return ''
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return ''
    return date.toISOString().split('T')[0]
  } catch (error) {
    console.error('Error formatting date:', error, dateStr)
    return ''
  }
}

// Lifecycle
onMounted(() => {
  loadTaskData()
  loadUsers()
})

// Watch for route changes
watch(() => route.params.id, () => {
  if (route.params.id) {
    loadTaskData()
  }
})
</script>
