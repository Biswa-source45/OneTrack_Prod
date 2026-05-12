<template>
  <div class="min-h-screen bg-gray-50 py-6">
    <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
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
            <h1 class="text-2xl font-bold text-gray-900">Edit Ticket Properties</h1>
            <p class="text-sm text-gray-600" v-if="ticketData">{{ ticketData.ticket_id }}</p>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
        <div class="flex items-center">
          <svg class="w-5 h-5 text-red-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-sm text-red-700">{{ error }}</p>
        </div>
      </div>

      <!-- Edit Form -->
      <div v-else-if="ticketData">
        <TicketForm
          :is-edit-mode="true"
          :ticket-data="ticketData"
          @success="handleSuccess"
          @cancel="goBack"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchTicketFullDetails } from '../../api/tickets'
import TicketForm from './TicketForm.vue'

const route = useRoute()
const router = useRouter()

const ticketId = computed(() => route.params.id)
const ticketData = ref(null)
const loading = ref(false)
const error = ref('')

// Load ticket data
const loadTicketData = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetchTicketFullDetails(ticketId.value)
    ticketData.value = response.ticket
  } catch (err) {
    console.error('Failed to load ticket:', err)
    error.value = 'Failed to load ticket data. Please try again.'
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

const handleSuccess = () => {
  // Navigate back to ticket details after successful edit
  const currentPath = router.currentRoute.value.path
  const basePath = currentPath.includes('/manager/') ? '/manager' : '/engineer'
  router.push(`${basePath}/tickets/${ticketId.value}`)
}

// Load data on mount
onMounted(() => {
  loadTicketData()
})
</script>
