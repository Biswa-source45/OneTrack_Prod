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
          <!-- Task Properties (No edit icon for engineers) -->
          <PropertySection 
            title="Task Properties" 
            :show-edit-icon="false"
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
                  <!-- Inline editable status for engineers -->
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
                  <!-- Read-only priority display for engineers -->
                  <span :class="`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityBadgeClass(task.priority)}`">
                    {{ task.priority || 'Medium' }}
                  </span>
                </div>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Due Date</span>
                <div class="mt-1">
                  <!-- Read-only due date display for engineers -->
                  <span class="text-sm font-medium text-gray-900">
                    {{ task.due_date ? formatDate(task.due_date) : 'No due date' }}
                  </span>
                </div>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Created By</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ getCreatorName(task) }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Assigned To</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ getAssignedUserName(task) }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Created</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ formatDate(task.created_at) }}</p>
              </div>
              <div>
                <span class="text-xs text-gray-500 uppercase tracking-wide">Last Updated</span>
                <p class="text-sm font-medium text-gray-900 mt-1">{{ formatDate(task.updated_at) }}</p>
              </div>
            </div>
          </PropertySection>

          <!-- Task Description -->
          <PropertySection title="Task Description">
            <div class="space-y-3">
              <div>
                <p class="text-sm text-gray-900 whitespace-pre-wrap">{{ task.description || 'No description provided' }}</p>
              </div>
            </div>
          </PropertySection>
        </div>

        <!-- Right Content Area -->
        <div class="flex-1 bg-white rounded-lg shadow-sm">
          <!-- Comments Section (Same as manager) -->
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
import { fetchEngineerTask, changeEngineerTaskStatus } from '../../api/engineer'
import { formatDateIST } from '../../utils/date'
import { formatUserName } from '../../utils/user'
import TaskCommentsTab from '../shared/tabs/TaskCommentsTab.vue'
import PropertySection from '../shared/PropertySection.vue'
import InlineEditDropdown from '../shared/InlineEditDropdown.vue'
import NotificationToast from '../shared/NotificationToast.vue'

const route = useRoute()
const router = useRouter()
const formatDate = formatDateIST

// Reactive state
const task = ref(null)
const loading = ref(false)
const error = ref('')

// Notification state
const notification = ref({
  show: false,
  type: 'success',
  title: '',
  message: ''
})

// Computed properties
const taskId = computed(() => route.params.id)

// Status options for engineers (same as manager)
const statusOptions = [
  { value: 'TODO', label: 'To Do' },
  { value: 'IN PROGRESS', label: 'In Progress' },
  { value: 'COMPLETED', label: 'Completed' },
  { value: 'ON_HOLD', label: 'On Hold' },
  { value: 'CANCELLED', label: 'Cancelled' }
]

// Badge styling functions (same as manager)
const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'TODO': return 'bg-gray-100 text-gray-800'
    case 'IN PROGRESS':
    case 'IN_PROGRESS': return 'bg-blue-100 text-blue-800'
    case 'COMPLETED': return 'bg-green-100 text-green-800'
    case 'ON_HOLD': return 'bg-yellow-100 text-yellow-800'
    case 'CANCELLED': return 'bg-red-100 text-red-800'
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
const getCreatorName = (task) => {
  if (!task?.creator) return 'Unknown'
  return formatUserName(task.creator)
}

const getAssignedUserName = (task) => {
  if (!task?.assigned_user) return 'Unassigned'
  return formatUserName(task.assigned_user)
}

// Load task data
const loadTaskData = async () => {
  if (!taskId.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetchEngineerTask(taskId.value)
    task.value = response.task
  } catch (err) {
    console.error('Failed to load task details:', err)
    error.value = 'Failed to load task details. Please try again.'
  } finally {
    loading.value = false
  }
}

// Handle status change (engineers can change status)
const handleStatusChange = async (newStatus) => {
  if (!task.value || !newStatus) return
  
  try {
    await changeEngineerTaskStatus(task.value.id, newStatus)
    task.value.task_status = newStatus
    showNotification('success', 'Status Updated', 'Task status has been updated successfully.')
  } catch (err) {
    console.error('Failed to update status:', err)
    showNotification('error', 'Update Failed', 'Failed to update task status. Please try again.')
  }
}

// Event handlers (same as manager)
const handleCommentAdded = () => {
  loadTaskData() // Refresh task data
}

const handleMarkCompleted = () => {
  if (task.value) {
    task.value.task_status = 'COMPLETED'
    showNotification('success', 'Task Completed', 'Task has been marked as completed.')
  }
}

// Navigation
const goBack = () => {
  router.push('/engineer/tasks')
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
    loadTaskData()
  }
}, { immediate: true })

// Lifecycle
onMounted(() => {
  loadTaskData()
})
</script>
