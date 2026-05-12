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
            <h1 class="text-2xl font-bold text-gray-900">Create New Task</h1>
            <p class="text-gray-600 mt-1">Create and assign a new task</p>
          </div>
        </div>
      </div>

      <!-- Create Form -->
      <div class="bg-white rounded-lg shadow-sm p-6">
        <form @submit.prevent="createTask" class="space-y-6">
          <!-- Subject -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Subject *</label>
            <input
              v-model="newTask.subject"
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
              v-model="newTask.due_date"
              type="date"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <!-- Task Owner -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Task Owner</label>
            <select
              v-model="newTask.assigned_to"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Unassigned</option>
              <option v-for="user in users" :key="user.id" :value="user.id">
                {{ formatUserName(user) }}
              </option>
            </select>
          </div>

          <!-- Priority -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Priority *</label>
            <select
              v-model="newTask.priority"
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
              v-model="newTask.description"
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
              :disabled="creating"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition-colors"
            >
              {{ creating ? 'Creating...' : 'Create Task' }}
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
              <h3 class="text-lg font-medium text-gray-900">Task Created</h3>
              <p class="text-sm text-gray-600 mt-1">Your task has been created successfully.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { createTask as apiCreateTask } from '../../api/tasks'
import { fetchUsers } from '../../api/tickets'
import { formatUserName } from '../../utils/user'

const router = useRouter()

// State
const users = ref([])
const creating = ref(false)
const showSuccess = ref(false)

// New task form
const newTask = ref({
  subject: '',
  description: '',
  due_date: '',
  assigned_to: '',
  priority: 'Medium'
})

// Load users
const loadUsers = async () => {
  try {
    const response = await fetchUsers()
    users.value = Array.isArray(response) ? response : (response?.users || [])
  } catch (err) {
    console.error('Failed to load users:', err)
  }
}

// Create task
const createTask = async () => {
  creating.value = true
  
  try {
    const taskData = {
      ...newTask.value,
      assigned_to: newTask.value.assigned_to || null,
      due_date: newTask.value.due_date || null
    }
    
    await apiCreateTask(taskData)
    
    // Show success message
    showSuccess.value = true
    
    // Redirect to tasks list after 2 seconds
    setTimeout(() => {
      router.push('/manager/tasks')
    }, 2000)
    
  } catch (err) {
    console.error('Failed to create task:', err)
    alert('Failed to create task. Please try again.')
  } finally {
    creating.value = false
  }
}

// Navigation
const goBack = () => {
  router.push('/manager/tasks')
}

// Lifecycle
onMounted(() => {
  loadUsers()
})
</script>
