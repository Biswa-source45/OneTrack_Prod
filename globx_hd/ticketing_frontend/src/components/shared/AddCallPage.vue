<template>
  <div class="min-h-screen bg-gray-50 py-6">
    <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="mb-8">
        <div class="flex items-center space-x-4 mb-4">
          <button
            @click="goBack"
            class="p-2 rounded-full hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
            title="Go back"
          >
            <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Add Call</h1>
            <!-- <p class="text-sm text-gray-600">Ticket #{{ ticketId }}</p> -->
          </div>
        </div>
      </div>

      <!-- Error Message -->
      <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
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

      <!-- Add Call Form -->
      <div class="bg-white shadow rounded-lg">
        <form @submit.prevent="submitCall" class="p-6 space-y-6">
          <!-- Subject -->
          <div>
            <label for="subject" class="block text-sm font-medium text-gray-700 mb-2">
              Subject <span class="text-red-500">*</span>
            </label>
            <input
              id="subject"
              v-model="form.subject"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="What is this call about?"
            />
          </div>

          <!-- Direction -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Direction <span class="text-red-500">*</span>
            </label>
            <div class="flex space-x-4">
              <label class="flex items-center">
                <input
                  v-model="form.direction"
                  type="radio"
                  value="Inbound"
                  required
                  class="mr-2 text-blue-600 focus:ring-blue-500"
                />
                <span class="text-sm text-gray-700">Inbound</span>
              </label>
              <label class="flex items-center">
                <input
                  v-model="form.direction"
                  type="radio"
                  value="Outbound"
                  required
                  class="mr-2 text-blue-600 focus:ring-blue-500"
                />
                <span class="text-sm text-gray-700">Outbound</span>
              </label>
            </div>
          </div>

          <!-- Call Status -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Call Status <span class="text-red-500">*</span>
            </label>
            <div class="flex space-x-4">
              <label class="flex items-center">
                <input
                  v-model="form.status"
                  type="radio"
                  value="Open"
                  required
                  class="mr-2 text-blue-600 focus:ring-blue-500"
                />
                <span class="text-sm text-gray-700">Open</span>
              </label>
              <label class="flex items-center">
                <input
                  v-model="form.status"
                  type="radio"
                  value="In Progress"
                  required
                  class="mr-2 text-blue-600 focus:ring-blue-500"
                />
                <span class="text-sm text-gray-700">In Progress</span>
              </label>
              <label class="flex items-center">
                <input
                  v-model="form.status"
                  type="radio"
                  value="Completed"
                  required
                  class="mr-2 text-blue-600 focus:ring-blue-500"
                />
                <span class="text-sm text-gray-700">Completed</span>
              </label>
            </div>
          </div>

          <!-- Start Time -->
          <div>
            <label for="startTime" class="block text-sm font-medium text-gray-700 mb-2">
              Start Time
            </label>
            <input
              id="startTime"
              v-model="form.startTime"
              type="datetime-local"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <!-- Description -->
          <div>
            <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
              Description
            </label>
            <textarea
              id="description"
              v-model="form.description"
              rows="4"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
              placeholder="Optional call details or notes..."
            ></textarea>
          </div>

          <!-- OEM Ticket ID -->
          <div>
            <label for="oemTicketId" class="block text-sm font-medium text-gray-700 mb-2">
              OEM Ticket ID
            </label>
            <input
              id="oemTicketId"
              v-model="form.oemTicketId"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="External OEM ticket reference (e.g., DELL-12345)"
            />
          </div>

          <!-- Due Date -->
          <div>
            <label for="dueDate" class="block text-sm font-medium text-gray-700 mb-2">
              Due Date
            </label>
            <input
              id="dueDate"
              v-model="form.dueDate"
              type="date"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <!-- Mail Content -->
          <div>
            <label for="mailContent" class="block text-sm font-medium text-gray-700 mb-2">
              Mail Content
            </label>
            <textarea
              id="mailContent"
              v-model="form.mailContent"
              rows="6"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none font-mono text-sm"
              placeholder="Paste email content here (text or HTML)..."
            ></textarea>
          </div>

          <!-- File Attachments -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Attachments
            </label>
            <div class="space-y-3">
              <!-- File Input -->
              <div class="flex items-center space-x-3">
                <input
                  ref="fileInput"
                  type="file"
                  multiple
                  @change="handleFileSelect"
                  class="hidden"
                  accept="*/*"
                />
                <button
                  type="button"
                  @click="$refs.fileInput.click()"
                  class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors flex items-center space-x-2"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                  </svg>
                  <span>Choose Files</span>
                </button>
                <span class="text-sm text-gray-500">Max 3MB per file</span>
              </div>

              <!-- Selected Files List -->
              <div v-if="selectedFiles.length > 0" class="space-y-2">
                <div
                  v-for="(file, index) in selectedFiles"
                  :key="index"
                  class="flex items-center justify-between p-3 bg-gray-50 border border-gray-200 rounded-lg"
                >
                  <div class="flex items-center space-x-3 flex-1 min-w-0">
                    <svg class="w-5 h-5 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <div class="flex-1 min-w-0">
                      <p class="text-sm font-medium text-gray-900 truncate">{{ file.name }}</p>
                      <p class="text-xs text-gray-500">{{ formatFileSize(file.size) }}</p>
                    </div>
                  </div>
                  <button
                    type="button"
                    @click="removeFile(index)"
                    class="ml-3 text-red-600 hover:text-red-800 flex-shrink-0"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              </div>

              <!-- File Size Warning -->
              <p v-if="fileSizeError" class="text-sm text-red-600">{{ fileSizeError }}</p>
            </div>
          </div>

          <!-- Form Actions -->
          <div class="flex items-center justify-end space-x-3 pt-6 border-t border-gray-200">
            <button
              type="button"
              @click="goBack"
              class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="submitting || !isFormValid"
              class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ submitting ? 'Adding...' : 'Add Call' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ticketCalls } from '../../api/tickets'

const route = useRoute()
const router = useRouter()

// Props
const ticketId = computed(() => route.params.id)

// State
const submitting = ref(false)
const error = ref('')

// Form data
const form = ref({
  subject: '',
  direction: '',
  status: '',
  startTime: '',
  description: '',
  oemTicketId: '',
  dueDate: '',
  mailContent: ''
})

// File handling
const selectedFiles = ref([])
const fileInput = ref(null)
const fileSizeError = ref('')

// Computed
const isFormValid = computed(() => {
  return form.value.subject.trim() && 
         form.value.direction && 
         form.value.status
})

// Methods
const goBack = () => {
  router.back()
}

const submitCall = async () => {
  if (!isFormValid.value) return
  
  submitting.value = true
  error.value = ''
  
  try {
    // Prepare payload
    const payload = {
      subject: form.value.subject.trim(),
      direction: form.value.direction,
      status: form.value.status,
      description: form.value.description.trim(),
      oem_ticket_id: form.value.oemTicketId.trim(),
      mail_content: form.value.mailContent.trim()
    }
    
    // Add start_time if provided
    if (form.value.startTime) {
      payload.start_time = new Date(form.value.startTime).toISOString()
    }
    
    // Add due_date if provided
    if (form.value.dueDate) {
      payload.due_date = new Date(form.value.dueDate).toISOString()
    }
    
    // Create call
    const response = await ticketCalls.schedule(ticketId.value, payload)
    const callId = response.call.id
    
    // Upload attachments if any
    if (selectedFiles.value.length > 0) {
      await uploadAttachments(callId)
    }
    
    // Redirect back to ticket details
    router.push(`/manager/tickets/${ticketId.value}`)
    
  } catch (err) {
    console.error('Failed to add call:', err)
    if (err.response?.status === 401) {
      error.value = 'Unauthorized: Please check your login status and permissions'
    } else {
      error.value = `Failed to add call: ${err.response?.data?.error || err.message}`
    }
  } finally {
    submitting.value = false
  }
}

// File handling methods
const handleFileSelect = (event) => {
  const files = Array.from(event.target.files)
  fileSizeError.value = ''
  
  // Check file sizes (3MB limit)
  const maxSize = 3 * 1024 * 1024 // 3MB
  const oversizedFiles = files.filter(file => file.size > maxSize)
  
  if (oversizedFiles.length > 0) {
    fileSizeError.value = `${oversizedFiles.length} file(s) exceed 3MB limit and were not added`
    const validFiles = files.filter(file => file.size <= maxSize)
    selectedFiles.value = [...selectedFiles.value, ...validFiles]
  } else {
    selectedFiles.value = [...selectedFiles.value, ...files]
  }
  
  // Reset file input
  event.target.value = ''
}

const removeFile = (index) => {
  selectedFiles.value.splice(index, 1)
  fileSizeError.value = ''
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const uploadAttachments = async (callId) => {
  if (selectedFiles.value.length === 0) return
  
  const formData = new FormData()
  selectedFiles.value.forEach(file => {
    formData.append('files', file)
  })
  
  try {
    await ticketCalls.uploadAttachments(ticketId.value, callId, formData)
  } catch (err) {
    console.error('Failed to upload attachments:', err)
    // Don't throw error, just log it - call was created successfully
  }
}

// Set default start time to now
onMounted(() => {
  const now = new Date()
  now.setMinutes(now.getMinutes() - now.getTimezoneOffset()) // Adjust for timezone
  form.value.startTime = now.toISOString().slice(0, 16) // Format for datetime-local
})
</script>
