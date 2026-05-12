<template>
  <div class="min-h-[600px] flex flex-col">
    <!-- Header -->
    <div class="flex justify-between items-center p-6 border-b border-gray-200">
      <h3 class="text-lg font-medium text-gray-900">Attachments</h3>
    </div>

    <!-- Attachments List -->
    <div class="flex-1 overflow-y-auto p-6">
      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <div v-else-if="error" class="text-center py-8 text-red-600">
        <p>{{ error }}</p>
        <button @click="loadAttachments" class="mt-2 text-blue-600 hover:text-blue-800">
          Try again
        </button>
      </div>

      <!-- Empty State with Upload Option (Manager Only) -->
      <div v-else-if="!attachments.length" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No attachments</h3>
        <p class="mt-1 text-sm text-gray-500 mb-6">Files uploaded to this ticket will appear here.</p>
        
        <!-- Upload Section for Empty State -->
        <div class="max-w-md mx-auto">
          <FileUpload 
            ref="fileUpload"
            :ticket-id="fullTicketId"
            @files-selected="onFilesSelected"
            @upload-success="onUploadSuccess"
            @upload-error="onUploadError"
          />
        </div>
      </div>

      <!-- Attachment Items -->
      <div v-else class="space-y-4">
        <!-- Upload Section (when attachments exist) -->
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
          <h4 class="text-sm font-medium text-blue-900 mb-3">Add More Attachments</h4>
          <FileUpload 
            ref="fileUpload"
            :ticket-id="fullTicketId"
            @files-selected="onFilesSelected"
            @upload-success="onUploadSuccess"
            @upload-error="onUploadError"
          />
        </div>

        <!-- Existing Attachments -->
        <div 
          v-for="attachment in attachments" 
          :key="attachment.id" 
          class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
          @click="openAttachment(attachment)"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1 min-w-0">
              <div class="flex items-center space-x-3 mb-2">
                <!-- File Type Icon -->
                <div class="flex-shrink-0">
                  <div :class="[
                    'w-10 h-10 rounded-lg flex items-center justify-center text-white text-sm font-medium',
                    getFileTypeColor(attachment.mime_type)
                  ]">
                    {{ getFileTypeIcon(attachment.mime_type) }}
                  </div>
                </div>
                
                <!-- File Info -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-center space-x-2">
                    <span class="font-medium text-gray-900 truncate">{{ attachment.original_filename }}</span>
                    <span :class="[
                      'px-2 py-1 text-xs font-medium rounded-full',
                      getFileTypeBadgeClass(attachment.mime_type)
                    ]">
                      {{ getFileTypeLabel(attachment.mime_type) }}
                    </span>
                  </div>
                  <div class="flex items-center space-x-4 text-sm text-gray-500 mt-1">
                    <span>{{ formatFileSize(attachment.file_size) }}</span>
                    <span>•</span>
                    <span>{{ formatDateTime(attachment.uploaded_at) }}</span>
                    <span>•</span>
                    <span>Uploaded by {{ getUploaderName(attachment.contact) }}</span>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Download Button -->
            <div class="flex-shrink-0 ml-4">
              <button
                @click.stop="downloadAttachment(attachment)"
                class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
                :disabled="downloading === attachment.id"
              >
                <svg v-if="downloading === attachment.id" class="animate-spin -ml-1 mr-2 h-4 w-4 text-gray-500" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                {{ downloading === attachment.id ? 'Downloading...' : 'Download' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { getTicketAttachments, downloadAttachment as downloadAttachmentAPI, createDownloadUrl } from '../../../api/attachments'
import { formatDateTime } from '../../../utils/date'
import FileUpload from '../FileUpload.vue'

const props = defineProps({
  ticketId: {
    type: [String, Number],
    required: true
  },
  ticket: {
    type: Object,
    default: null
  }
})

// Reactive state
const attachments = ref([])
const loading = ref(false)
const error = ref('')
const downloading = ref(null)
const selectedFiles = ref([])

// File upload ref
const fileUpload = ref(null)

// Computed full ticket ID
const fullTicketId = computed(() => props.ticket?.ticket_id || props.ticketId)

// Load attachments
const loadAttachments = async () => {
  if (!fullTicketId.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    console.log('🔥 DEBUG: Loading attachments for manager ticket:', fullTicketId.value)
    const response = await getTicketAttachments(fullTicketId.value)
    attachments.value = response.attachments || []
    console.log('✅ DEBUG: Loaded attachments:', attachments.value.length)
  } catch (err) {
    console.error('Failed to load attachments:', err)
    error.value = 'Failed to load attachments. Please try again.'
  } finally {
    loading.value = false
  }
}

// Download attachment
const downloadAttachment = async (attachment) => {
  downloading.value = attachment.id
  
  try {
    const response = await downloadAttachmentAPI(attachment.id)
    createDownloadUrl(response.data, attachment.original_filename)
  } catch (err) {
    console.error('Failed to download attachment:', err)
    error.value = 'Failed to download attachment. Please try again.'
  } finally {
    downloading.value = null
  }
}

// Open attachment in new tab (for preview)
const openAttachment = async (attachment) => {
  try {
    console.log('🔥 DEBUG: Opening attachment for preview:', attachment.original_filename)
    
    // Download the file through our authenticated API
    const response = await downloadAttachmentAPI(attachment.id)
    
    // Create a blob URL for preview
    const blob = response.data
    const blobUrl = window.URL.createObjectURL(blob)
    
    // Open in new tab
    window.open(blobUrl, '_blank')
    
    // Clean up the blob URL after a short delay to allow the tab to load
    setTimeout(() => {
      window.URL.revokeObjectURL(blobUrl)
    }, 1000)
    
  } catch (err) {
    console.error('Failed to open attachment:', err)
    error.value = 'Failed to open attachment. Please try again.'
  }
}

// File upload event handlers
const onFilesSelected = (files) => {
  selectedFiles.value = files
  console.log('🔥 DEBUG: Files selected for manager upload:', files.length)
}

const onUploadSuccess = (response) => {
  console.log('✅ DEBUG: Manager upload successful:', response)
  selectedFiles.value = []
  
  // Clear the file upload component
  if (fileUpload.value) {
    fileUpload.value.clearFiles()
  }
  
  // Reload attachments to show the newly uploaded files
  loadAttachments()
  
  // Clear any previous errors
  error.value = ''
}

const onUploadError = (errorMessage) => {
  console.error('❌ DEBUG: Manager upload failed:', errorMessage)
  error.value = errorMessage
}

// File type utilities (same as customer side)
const getFileTypeIcon = (mimeType) => {
  if (mimeType.startsWith('image/')) return '🖼️'
  if (mimeType.includes('pdf')) return '📄'
  if (mimeType.includes('word') || mimeType.includes('document')) return '📝'
  if (mimeType.includes('excel') || mimeType.includes('spreadsheet')) return '📊'
  if (mimeType.includes('powerpoint') || mimeType.includes('presentation')) return '📽️'
  if (mimeType.includes('zip') || mimeType.includes('archive')) return '🗜️'
  if (mimeType.startsWith('video/')) return '🎥'
  if (mimeType.startsWith('audio/')) return '🎵'
  return '📎'
}

const getFileTypeColor = (mimeType) => {
  if (mimeType.startsWith('image/')) return 'bg-green-500'
  if (mimeType.includes('pdf')) return 'bg-red-500'
  if (mimeType.includes('word') || mimeType.includes('document')) return 'bg-blue-500'
  if (mimeType.includes('excel') || mimeType.includes('spreadsheet')) return 'bg-green-600'
  if (mimeType.includes('powerpoint') || mimeType.includes('presentation')) return 'bg-orange-500'
  if (mimeType.includes('zip') || mimeType.includes('archive')) return 'bg-purple-500'
  if (mimeType.startsWith('video/')) return 'bg-pink-500'
  if (mimeType.startsWith('audio/')) return 'bg-indigo-500'
  return 'bg-gray-500'
}

const getFileTypeBadgeClass = (mimeType) => {
  if (mimeType.startsWith('image/')) return 'bg-green-100 text-green-800'
  if (mimeType.includes('pdf')) return 'bg-red-100 text-red-800'
  if (mimeType.includes('word') || mimeType.includes('document')) return 'bg-blue-100 text-blue-800'
  if (mimeType.includes('excel') || mimeType.includes('spreadsheet')) return 'bg-green-100 text-green-800'
  if (mimeType.includes('powerpoint') || mimeType.includes('presentation')) return 'bg-orange-100 text-orange-800'
  if (mimeType.includes('zip') || mimeType.includes('archive')) return 'bg-purple-100 text-purple-800'
  if (mimeType.startsWith('video/')) return 'bg-pink-100 text-pink-800'
  if (mimeType.startsWith('audio/')) return 'bg-indigo-100 text-indigo-800'
  return 'bg-gray-100 text-gray-800'
}

const getFileTypeLabel = (mimeType) => {
  if (mimeType.startsWith('image/')) return 'Image'
  if (mimeType.includes('pdf')) return 'PDF'
  if (mimeType.includes('word') || mimeType.includes('document')) return 'Document'
  if (mimeType.includes('excel') || mimeType.includes('spreadsheet')) return 'Spreadsheet'
  if (mimeType.includes('powerpoint') || mimeType.includes('presentation')) return 'Presentation'
  if (mimeType.includes('zip') || mimeType.includes('archive')) return 'Archive'
  if (mimeType.startsWith('video/')) return 'Video'
  if (mimeType.startsWith('audio/')) return 'Audio'
  return 'File'
}

// Format file size
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Get uploader name
const getUploaderName = (contact) => {
  if (!contact) return 'Unknown'
  return `${contact.first_name || ''} ${contact.last_name || ''}`.trim() || 'Unknown'
}

// Computed
const attachmentCount = computed(() => attachments.value.length)

// Watchers
watch(() => props.ticket, (newTicket) => {
  if (newTicket && newTicket.ticket_id) {
    loadAttachments()
  }
}, { immediate: true })

// Lifecycle
onMounted(() => {
  loadAttachments()
})

// Expose count for parent component
defineExpose({
  attachmentCount
})
</script>
