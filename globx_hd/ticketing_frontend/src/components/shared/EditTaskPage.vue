<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
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
            <h1 class="text-2xl font-bold text-gray-900">Edit Task</h1>
            <p class="text-gray-600 mt-1">Update task details</p>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
        <div class="flex">
          <svg class="w-5 h-5 text-red-400 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error loading task</h3>
            <p class="text-sm text-red-700 mt-1">{{ error }}</p>
          </div>
        </div>
      </div>

      <!-- Edit Form -->
      <div v-else-if="task" class="bg-white rounded-lg shadow-sm p-6">
        <form @submit.prevent="saveTask" class="space-y-6">
          <!-- Subject -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Subject *</label>
            <input
              v-model="form.subject"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter task subject"
            />
          </div>

          <!-- Due Date -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Due Date</label>
            <input
              v-model="form.due_date"
              type="date"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <!-- Task Owner -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Task Owner</label>
            <select
              v-model="form.assigned_to"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Unassigned</option>
              <option v-for="user in users" :key="user.id" :value="user.id">
                {{ formatUserName(user) }}
              </option>
            </select>
          </div>

          <!-- Status -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status *</label>
            <select
              v-model="form.task_status"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="Not Started">Not Started</option>
              <option value="In Progress">In Progress</option>
              <option value="Completed">Completed</option>
              <option value="Deferred">Deferred</option>
              <option value="Waiting on someone else">Waiting on someone else</option>
              <option value="Canceled">Canceled</option>
            </select>
          </div>

          <!-- Priority -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Priority *</label>
            <select
              v-model="form.priority"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="High">High</option>
              <option value="Medium">Medium</option>
              <option value="Low">Low</option>
            </select>
          </div>

          <!-- Description -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea
              v-model="form.description"
              rows="4"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter task description"
            ></textarea>
          </div>

          <!-- Actions -->
          <div class="flex justify-end space-x-3 pt-4">
            <button
              type="button"
              @click="goBack"
              class="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="saving"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition-colors"
            >
              {{ saving ? 'Saving...' : 'Save Changes' }}
            </button>
          </div>
        </form>
      </div>

      <!-- Success Message -->
      <div v-if="showSuccess" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 max-w-sm mx-4">
          <div class="flex items-center space-x-3">
            <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <div>
              <h3 class="text-lg font-medium text-gray-900">Task Updated</h3>
              <p class="text-sm text-gray-600 mt-1">Your changes have been saved successfully.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchTask, updateTask } from '../../api/tasks'
import { fetchUsers } from '../../api/tickets'
import { formatUserName } from '../../utils/user'

const route = useRoute()
const router = useRouter()

// State
const task = ref(null)
const users = ref([])
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const showSuccess = ref(false)

// Form data
const form = ref({
  subject: '',
  description: '',
  due_date: '',
  assigned_to: '',
  task_status: '',
  priority: ''
})

// Load task data
const loadTask = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetchTask(route.params.id)
    task.value = response.task
    
    // Populate form
    form.value = {
      subject: task.value.subject || '',
      description: task.value.description || '',
      due_date: task.value.due_date ? formatDateForInput(task.value.due_date) : '',
      assigned_to: task.value.assigned_to || '',
      task_status: task.value.task_status || '',
      priority: task.value.priority || ''
    }
  } catch (err) {
    console.error('Failed to load task:', err)
    error.value = 'Failed to load task details. Please try again.'
  } finally {
    loading.value = false
  }
}

// Load users
const loadUsers = async () => {
  try {
    const response = await fetchUsers()
    users.value = Array.isArray(response) ? response : (response?.users || [])
  } catch (err) {
    console.error('Failed to load users:', err)
  }
}

// Save task
const saveTask = async () => {
  saving.value = true
  error.value = ''
  
  try {
    const updateData = {
      subject: form.value.subject,
      description: form.value.description,
      due_date: form.value.due_date || null,
      assigned_to: form.value.assigned_to || null,
      task_status: form.value.task_status,
      priority: form.value.priority
    }
    
    await updateTask(task.value.id, updateData)
    
    // Show success message
    showSuccess.value = true
    
    // Redirect back after 2 seconds
    setTimeout(() => {
      goBack()
    }, 2000)
    
  } catch (err) {
    console.error('Failed to save task:', err)
    error.value = 'Failed to save task. Please try again.'
  } finally {
    saving.value = false
  }
}

// Navigation
const goBack = () => {
  router.push(`/manager/tasks/${route.params.id}`)
}

// Utility functions
const formatDateForInput = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toISOString().split('T')[0]
}

// Lifecycle
onMounted(() => {
  loadTask()
  loadUsers()
})
</script>
