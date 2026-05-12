<template>
  <div class="h-full flex flex-col">
    <!-- Header -->
    <div class="border-b border-gray-200 p-6">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-medium text-gray-900">Comments</h3>
        <button
          @click="toggleCommentForm"
          class="p-2 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100 transition-colors"
          :title="showCommentForm ? 'Hide comment form' : 'Add a comment'"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.955 8.955 0 01-4.126-.98L3 20l1.98-5.874A8.955 8.955 0 013 12a8 8 0 018-8c4.418 0 8 3.582 8 8z" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Comments List -->
    <div class="flex-1 overflow-y-auto p-6">
      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <div v-else-if="error" class="text-center py-8 text-red-600">
        <p>{{ error }}</p>
        <button @click="loadComments" class="mt-2 text-blue-600 hover:text-blue-800">
          Try again
        </button>
      </div>

      <div v-else>
        <!-- Original Task Details (First Message) -->
        <div v-if="task && task.description" class="bg-blue-50 border border-blue-200 rounded-lg p-4 shadow-sm mb-6">
          <div class="flex items-start space-x-3">
            <!-- User Avatar -->
            <div class="flex-shrink-0">
              <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
                {{ getUserInitials(task.creator) }}
              </div>
            </div>
            
            <!-- Original Task Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center space-x-2">
                  <p class="text-sm font-medium text-gray-900">{{ getCreatorName(task.creator) }}</p>
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                    Original Task
                  </span>
                </div>
                <span class="text-xs text-gray-500">{{ formatDateTime(task.created_at) }}</span>
              </div>
              
              <div class="prose prose-sm max-w-none">
                <p class="text-gray-700 whitespace-pre-wrap">{{ task.description }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- No Additional Comments Message -->
        <div v-if="comments.length === 0" class="text-center py-8 text-gray-500">
          <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.955 8.955 0 01-4.126-.98L3 20l1.98-5.874A8.955 8.955 0 013 12c0-4.418 3.582-8 8-8s8 3.582 8 8z" />
          </svg>
          <p class="text-lg font-medium">No additional comments yet</p>
          <p class="text-sm">Use the comment icon above to add comments.</p>
        </div>

        <!-- Comment Items -->
        <div v-else class="space-y-4">
          <div v-for="comment in comments" :key="comment.id" class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
            <div class="flex items-start space-x-3">
              <!-- User Avatar -->
              <div class="flex-shrink-0">
                <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
                  {{ getUserInitials(comment.user) }}
                </div>
              </div>
              
              <!-- Comment Content -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center space-x-2">
                    <p class="text-sm font-medium text-gray-900">
                      {{ getUserName(comment.user) }}
                    </p>
                    <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                      Comment
                    </span>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span class="text-xs text-gray-500">{{ formatDateTime(comment.created_at) }}</span>
                    <div class="flex space-x-1">
                      <!-- Edit Button -->
                      <button
                        v-if="canEditComment(comment) && editingComment?.id !== comment.id"
                        @click="editComment(comment)"
                        class="text-gray-400 hover:text-gray-600 p-1 rounded"
                        title="Edit comment"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                      </button>
                      <!-- Delete Button -->
                      <button
                        v-if="canEditComment(comment)"
                        @click="deleteComment(comment)"
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
                
                <!-- Comment Content -->
                <div class="mt-3">
            <div v-if="editingComment?.id === comment.id">
              <textarea
                v-model="editContent"
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
              ></textarea>
              <div class="flex justify-end space-x-2 mt-2">
                <button
                  @click="cancelEdit"
                  class="px-3 py-1 text-sm text-gray-600 hover:text-gray-800"
                >
                  Cancel
                </button>
                <button
                  @click="saveEdit"
                  :disabled="saving"
                  class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                >
                  {{ saving ? 'Saving...' : 'Save' }}
                </button>
              </div>
            </div>
                  <div v-else class="prose prose-sm max-w-none">
                    <p class="text-gray-700 whitespace-pre-wrap">{{ comment.content }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Comment Form -->
    <div v-if="showCommentForm" class="border-t border-gray-200 p-6 bg-gray-50">
      <form @submit.prevent="addComment" class="space-y-4">
        <div>
          <textarea
            v-model="newComment.content"
            rows="3"
            placeholder="Add a comment..."
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          ></textarea>
        </div>
        
        <div class="flex items-center justify-end">
          <button
            type="submit"
            :disabled="submitting || !newComment.content.trim()"
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ submitting ? 'Adding...' : 'Add Comment' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Mark as Completed Section -->
    <div class="border-t border-gray-200 p-6 bg-gray-50">
      <div class="flex items-center justify-end">
        <button
          @click="markCompleted"
          class="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
        >
          Mark as Completed
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { taskComments } from '../../../api/tasks'
import { formatDateTime } from '../../../utils/date'
import { formatUserName } from '../../../utils/user'
import { useAuthStore } from '../../../stores/auth'
import { decodeJWT } from '../../../utils/jwt'

const props = defineProps({
  taskId: {
    type: [String, Number],
    required: true
  },
  task: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['comment-added', 'mark-completed'])

// Auth store
const auth = useAuthStore()

// Check if current user can add comments
const canAddComments = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return false
  
  const claims = decodeJWT(token)
  return claims?.type === 'user'
})

// State
const comments = ref([])
const loading = ref(false)
const submitting = ref(false)
const saving = ref(false)
const error = ref('')
const showCommentForm = ref(false)

// New comment form
const newComment = ref({
  content: ''
})

// Edit comment
const editingComment = ref(null)
const editContent = ref('')

// Toggle comment form
const toggleCommentForm = () => {
  showCommentForm.value = !showCommentForm.value
}

// Load comments
const loadComments = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await taskComments.fetch(props.taskId, { limit: 100 })
    comments.value = response.comments || []
  } catch (err) {
    console.error('Failed to load comments:', err)
    error.value = 'Failed to load comments'
  } finally {
    loading.value = false
  }
}

// Add comment
const addComment = async () => {
  if (!newComment.value.content.trim()) return
  
  submitting.value = true
  error.value = ''
  
  try {
    const response = await taskComments.create(props.taskId, {
      content: newComment.value.content.trim(),
      is_internal: false
    })
    
    // Add new comment to the list
    comments.value.unshift(response.comment)
    
    // Reset form and hide it
    newComment.value = {
      content: ''
    }
    showCommentForm.value = false
    
    // Emit event to parent to refresh task data
    emit('comment-added')
    
  } catch (err) {
    console.error('Failed to add comment:', err)
    if (err.response?.status === 401) {
      error.value = 'You are not authorized to add comments'
    } else {
      error.value = `Failed to add comment: ${err.response?.data?.error || err.message}`
    }
  } finally {
    submitting.value = false
  }
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
    const response = await taskComments.update(props.taskId, editingComment.value.id, {
      content: editContent.value.trim()
    })
    
    // Update comment in the list
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
    await taskComments.delete(props.taskId, comment.id)
    
    // Remove comment from the list
    const index = comments.value.findIndex(c => c.id === comment.id)
    if (index !== -1) {
      comments.value.splice(index, 1)
    }
  } catch (err) {
    console.error('Failed to delete comment:', err)
    error.value = 'Failed to delete comment'
  }
}

// Mark as completed
const markCompleted = () => {
  emit('mark-completed')
}

// Check if user can edit comment
const canEditComment = (comment) => {
  const token = localStorage.getItem('token')
  if (!token) return false
  
  const claims = decodeJWT(token)
  // User can edit their own comments or manager can edit any comment
  // JWT uses 'sub' field for user ID, comment uses 'user_id'
  return claims?.sub === comment.user_id || claims?.type === 'manager'
}

// Utility functions
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

const getCreatorName = (creator) => {
  if (!creator) return 'Unknown'
  return formatUserName(creator)
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
