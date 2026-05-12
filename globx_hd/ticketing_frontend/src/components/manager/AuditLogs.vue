<template>
  <div class="p-6 bg-white rounded-lg shadow-md w-full overflow-hidden">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900">Audit Logs</h1>
        <p class="text-gray-600 mt-1">View and filter system audit trail</p>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6 w-full">
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div class="flex items-center">
            <div class="p-2 bg-blue-100 rounded-lg">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Total Logs</p>
              <p class="text-xl font-semibold text-gray-900">{{ stats.total_logs || 0 }}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div class="flex items-center">
            <div class="p-2 bg-red-100 rounded-lg">
              <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Critical</p>
              <p class="text-xl font-semibold text-gray-900">{{ stats.critical_logs || 0 }}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div class="flex items-center">
            <div class="p-2 bg-yellow-100 rounded-lg">
              <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Failed</p>
              <p class="text-xl font-semibold text-gray-900">{{ stats.failed_logs || 0 }}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div class="flex items-center">
            <div class="p-2 bg-green-100 rounded-lg">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Auth Events</p>
              <p class="text-xl font-semibold text-gray-900">{{ stats.auth_events || 0 }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-3 mb-4 w-full">
        <div class="grid grid-cols-2 lg:grid-cols-5 gap-2 w-full">
          <!-- Search -->
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">Search</label>
            <input
              v-model="filters.search"
              type="text"
              placeholder="Search..."
              class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
              @input="debouncedFetch"
            />
          </div>

          <!-- Action Filter -->
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">Action</label>
            <select
              v-model="filters.action"
              class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
              @change="fetchLogs"
            >
              <option value="">All Actions</option>
              <optgroup label="Authentication">
                <option value="USER_LOGIN_SUCCESS">User Login Success</option>
                <option value="USER_LOGIN_FAILURE">User Login Failure</option>
                <option value="CONTACT_LOGIN_SUCCESS">Contact Login Success</option>
                <option value="CONTACT_LOGIN_FAILURE">Contact Login Failure</option>
                <option value="LOGOUT">Logout</option>
                <option value="PASSWORD_RESET">Password Reset</option>
              </optgroup>
              <optgroup label="User Management">
                <option value="USER_CREATED">User Created</option>
                <option value="USER_UPDATED">User Updated</option>
                <option value="USER_DELETED">User Deleted</option>
              </optgroup>
              <optgroup label="Contact Management">
                <option value="CONTACT_CREATED">Contact Created</option>
                <option value="CONTACT_UPDATED">Contact Updated</option>
                <option value="CONTACT_DELETED">Contact Deleted</option>
              </optgroup>
              <optgroup label="Account Management">
                <option value="ACCOUNT_CREATED">Account Created</option>
                <option value="ACCOUNT_UPDATED">Account Updated</option>
                <option value="ACCOUNT_DELETED">Account Deleted</option>
              </optgroup>
              <optgroup label="Ticket Management">
                <option value="TICKET_CREATED">Ticket Created</option>
                <option value="TICKET_UPDATED">Ticket Updated</option>
                <option value="TICKET_DELETED">Ticket Deleted</option>
                <option value="TICKET_STATUS_CHANGED">Ticket Status Changed</option>
                <option value="TICKET_ASSIGNED">Ticket Assigned</option>
              </optgroup>
              <optgroup label="Task Management">
                <option value="TASK_CREATED">Task Created</option>
                <option value="TASK_UPDATED">Task Updated</option>
                <option value="TASK_DELETED">Task Deleted</option>
                <option value="TASK_STATUS_CHANGED">Task Status Changed</option>
              </optgroup>
              <optgroup label="Master Data">
                <option value="PRODUCT_CREATED">Product Created</option>
                <option value="PRODUCT_UPDATED">Product Updated</option>
                <option value="PRODUCT_DELETED">Product Deleted</option>
                <option value="PRODUCT_ISSUE_CREATED">Product Issue Created</option>
                <option value="PRODUCT_ISSUE_UPDATED">Product Issue Updated</option>
                <option value="PRODUCT_ISSUE_DELETED">Product Issue Deleted</option>
              </optgroup>
              <optgroup label="System">
                <option value="EMAIL_TO_TICKET">Email to Ticket</option>
                <option value="DUMPED_QUERY_CREATED">Dumped Query Created</option>
                <option value="DUMPED_QUERY_RESOLVED">Dumped Query Resolved</option>
                <option value="DUMPED_QUERY_DELETED">Dumped Query Deleted</option>
              </optgroup>
            </select>
          </div>

          <!-- Entity Type Filter -->
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">Entity</label>
            <select
              v-model="filters.entity_type"
              class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
              @change="fetchLogs"
            >
              <option value="">All Entities</option>
              <option value="user">User</option>
              <option value="contact">Contact</option>
              <option value="account">Account</option>
              <option value="ticket">Ticket</option>
              <option value="task">Task</option>
              <option value="product">Product</option>
              <option value="product_issue">Product Issue</option>
              <option value="role">Role</option>
              <option value="designation">Designation</option>
              <option value="attachment">Attachment</option>
              <option value="comment">Comment</option>
              <option value="dumped_query">Dumped Query</option>
              <option value="system">System</option>
            </select>
          </div>

          <!-- Status Filter -->
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">Status</label>
            <select
              v-model="filters.status"
              class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
              @change="fetchLogs"
            >
              <option value="">All Status</option>
              <option value="success">Success</option>
              <option value="failure">Failure</option>
              <option value="pending">Pending</option>
            </select>
          </div>

          <!-- Clear Filters -->
          <div class="flex items-end">
            <button
              @click="clearFilters"
              class="w-full px-2 py-1.5 text-sm text-gray-700 bg-gray-100 hover:bg-gray-200 rounded transition-colors"
            >
              Clear
            </button>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="bg-white rounded-lg shadow-sm border border-gray-200 p-8 w-full">
        <div class="flex justify-center items-center">
          <svg class="animate-spin h-8 w-8 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span class="ml-3 text-gray-600">Loading audit logs...</span>
        </div>
      </div>

      <!-- Logs Table -->
      <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden w-full">
        <div class="overflow-x-auto w-full">
          <table class="w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-28">Timestamp</th>
                <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-24">Actor</th>
                <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-32">Action</th>
                <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-24">Entity</th>
                <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-20">Status</th>
                <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-20">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-if="logs.length === 0">
                <td colspan="7" class="px-3 py-8 text-center text-gray-500">
                  No audit logs found matching your criteria
                </td>
              </tr>
              <tr v-for="log in logs" :key="log.id" class="hover:bg-gray-50">
                <td class="px-3 py-3 text-xs text-gray-500">
                  {{ formatDate(log.created_at) }}
                </td>
                <td class="px-3 py-3">
                  <div class="text-xs font-medium text-gray-900 truncate">{{ log.actor_name || 'System' }}</div>
                  <div class="text-xs text-gray-500 truncate">{{ log.actor_type }}</div>
                </td>
                <td class="px-3 py-3">
                  <span :class="getActionBadgeClass(log.action)" class="px-2 py-1 text-xs font-medium rounded-full inline-block truncate max-w-full">
                    {{ formatAction(log.action) }}
                  </span>
                </td>
                <td class="px-3 py-3">
                  <div class="text-xs font-medium text-gray-900 truncate">{{ formatEntityType(log.entity_type) }}</div>
                  <div class="text-xs text-gray-500 truncate">{{ log.entity_name || '-' }}</div>
                </td>
                <td class="px-3 py-3 text-xs text-gray-500">
                  <div class="truncate">{{ log.description || '-' }}</div>
                </td>
                <td class="px-3 py-3">
                  <span :class="getStatusBadgeClass(log.status)" class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ log.status }}
                  </span>
                </td>
                <td class="px-3 py-3 text-xs">
                  <button
                    @click="viewDetails(log)"
                    class="text-blue-600 hover:text-blue-800 font-medium whitespace-nowrap"
                  >
                    View
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div class="bg-gray-50 px-3 py-3 border-t border-gray-200 flex flex-wrap items-center justify-between gap-2">
          <div class="text-xs text-gray-700">
            <span class="font-medium">{{ pagination.total_count }}</span> results
          </div>
          <div class="flex items-center space-x-1">
            <button
              @click="goToPage(pagination.current_page - 1)"
              :disabled="pagination.current_page <= 1"
              :class="pagination.current_page <= 1 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-200'"
              class="px-2 py-1 text-xs font-medium text-gray-700 bg-white border border-gray-300 rounded"
            >
              Prev
            </button>
            <span class="px-2 py-1 text-xs text-gray-700">
              {{ pagination.current_page }}/{{ pagination.total_pages }}
            </span>
            <button
              @click="goToPage(pagination.current_page + 1)"
              :disabled="pagination.current_page >= pagination.total_pages"
              :class="pagination.current_page >= pagination.total_pages ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-200'"
              class="px-2 py-1 text-xs font-medium text-gray-700 bg-white border border-gray-300 rounded"
            >
              Next
            </button>
          </div>
        </div>
      </div>

      <!-- Details Modal -->
      <div v-if="showModal" class="fixed inset-0 z-50 overflow-y-auto">
        <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:p-0">
          <div class="fixed inset-0 transition-opacity" @click="closeModal">
            <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
          </div>

          <div class="inline-block w-full max-w-2xl my-8 overflow-hidden text-left align-middle transition-all transform bg-white rounded-lg shadow-xl">
            <!-- Modal Header -->
            <div class="px-6 py-4 bg-blue-600 text-white">
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold">Audit Log Details</h3>
                <button @click="closeModal" class="text-white hover:text-gray-200">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </div>
            </div>

            <!-- Modal Body -->
            <div class="px-6 py-4 max-h-96 overflow-y-auto" v-if="selectedLog">
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <p class="text-sm font-medium text-gray-500">Timestamp</p>
                  <p class="text-sm text-gray-900">{{ formatDate(selectedLog.created_at) }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Status</p>
                  <span :class="getStatusBadgeClass(selectedLog.status)" class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ selectedLog.status }}
                  </span>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Actor</p>
                  <p class="text-sm text-gray-900">{{ selectedLog.actor_name || 'System' }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Actor Type</p>
                  <p class="text-sm text-gray-900">{{ selectedLog.actor_type }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Action</p>
                  <span :class="getActionBadgeClass(selectedLog.action)" class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ formatAction(selectedLog.action) }}
                  </span>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Entity Type</p>
                  <p class="text-sm text-gray-900">{{ formatEntityType(selectedLog.entity_type) }}</p>
                </div>
                <div class="col-span-2">
                  <p class="text-sm font-medium text-gray-500">Entity Name</p>
                  <p class="text-sm text-gray-900">{{ selectedLog.entity_name || '-' }}</p>
                </div>
                <div class="col-span-2">
                  <p class="text-sm font-medium text-gray-500">Description</p>
                  <p class="text-sm text-gray-900">{{ selectedLog.description || '-' }}</p>
                </div>
                <div class="col-span-2">
                  <p class="text-sm font-medium text-gray-500">IP Address</p>
                  <p class="text-sm text-gray-900">{{ selectedLog.ip_address || '-' }}</p>
                </div>
                <div class="col-span-2">
                  <p class="text-sm font-medium text-gray-500">User Agent</p>
                  <p class="text-sm text-gray-900 break-all">{{ selectedLog.user_agent || '-' }}</p>
                </div>
              </div>

              <!-- Old Values -->
              <div v-if="selectedLog.old_values" class="mt-4">
                <p class="text-sm font-medium text-gray-500 mb-2">Old Values</p>
                <pre class="bg-gray-100 p-3 rounded-lg text-xs overflow-x-auto">{{ formatJSON(selectedLog.old_values) }}</pre>
              </div>

              <!-- New Values -->
              <div v-if="selectedLog.new_values" class="mt-4">
                <p class="text-sm font-medium text-gray-500 mb-2">New Values</p>
                <pre class="bg-gray-100 p-3 rounded-lg text-xs overflow-x-auto">{{ formatJSON(selectedLog.new_values) }}</pre>
              </div>

              <!-- Changes Summary -->
              <div v-if="selectedLog.changes_summary" class="mt-4">
                <p class="text-sm font-medium text-gray-500 mb-2">Changes Summary</p>
                <p class="text-sm text-gray-900 bg-yellow-50 p-3 rounded-lg">{{ selectedLog.changes_summary }}</p>
              </div>

              <!-- Error Message -->
              <div v-if="selectedLog.error_message" class="mt-4">
                <p class="text-sm font-medium text-red-500 mb-2">Error Message</p>
                <p class="text-sm text-red-700 bg-red-50 p-3 rounded-lg">{{ selectedLog.error_message }}</p>
              </div>
            </div>

            <!-- Modal Footer -->
            <div class="px-6 py-4 bg-gray-50 flex justify-end">
              <button
                @click="closeModal"
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import auditLogsAPI from '@/api/auditLogs'

export default {
  name: 'AuditLogs',
  setup() {
    const logs = ref([])
    const loading = ref(true)
    const showModal = ref(false)
    const selectedLog = ref(null)
    
    const stats = reactive({
      total_logs: 0,
      critical_logs: 0,
      failed_logs: 0,
      auth_events: 0
    })

    const filters = reactive({
      search: '',
      action: '',
      entity_type: '',
      status: ''
    })

    const pagination = reactive({
      current_page: 1,
      total_pages: 1,
      total_count: 0,
      limit: 20
    })

    let debounceTimer = null

    const fetchLogs = async () => {
      loading.value = true
      try {
        const params = {
          page: pagination.current_page,
          limit: pagination.limit
        }
        
        if (filters.search) params.search = filters.search
        if (filters.action) params.action = filters.action
        if (filters.entity_type) params.entity_type = filters.entity_type
        if (filters.status) params.status = filters.status

        const response = await auditLogsAPI.getAuditLogs(params)
        
        if (response) {
          logs.value = response.logs || []
          if (response.pagination) {
            pagination.current_page = response.pagination.current_page || 1
            pagination.total_pages = response.pagination.total_pages || 1
            pagination.total_count = response.pagination.total_count || 0
            pagination.limit = response.pagination.limit || 20
          }
        }
      } catch (error) {
        console.error('Error fetching audit logs:', error)
        logs.value = []
      } finally {
        loading.value = false
      }
    }

    const fetchStats = async () => {
      try {
        const response = await auditLogsAPI.getStats()
        if (response && response.stats) {
          stats.total_logs = response.stats.total_logs || 0
          stats.critical_logs = response.stats.critical_logs || 0
          stats.failed_logs = response.stats.failed_logs || 0
          stats.auth_events = response.stats.auth_events || 0
        }
      } catch (error) {
        console.error('Error fetching stats:', error)
      }
    }

    const debouncedFetch = () => {
      clearTimeout(debounceTimer)
      debounceTimer = setTimeout(() => {
        pagination.current_page = 1
        fetchLogs()
      }, 300)
    }

    const clearFilters = () => {
      filters.search = ''
      filters.action = ''
      filters.entity_type = ''
      filters.status = ''
      pagination.current_page = 1
      fetchLogs()
    }

    const goToPage = (page) => {
      if (page >= 1 && page <= pagination.total_pages) {
        pagination.current_page = page
        fetchLogs()
      }
    }

    const viewDetails = (log) => {
      selectedLog.value = log
      showModal.value = true
    }

    const closeModal = () => {
      showModal.value = false
      selectedLog.value = null
    }

    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString('en-IN', {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    }

    const formatAction = (action) => {
      if (!action) return '-'
      return action.replace(/_/g, ' ')
    }

    const formatEntityType = (entityType) => {
      if (!entityType) return '-'
      return entityType.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
    }

    const formatJSON = (jsonString) => {
      if (!jsonString) return '-'
      try {
        const parsed = typeof jsonString === 'string' ? JSON.parse(jsonString) : jsonString
        return JSON.stringify(parsed, null, 2)
      } catch {
        return jsonString
      }
    }

    const getActionBadgeClass = (action) => {
      if (!action) return 'bg-gray-100 text-gray-800'
      
      if (action.includes('DELETE') || action.includes('FAILURE')) {
        return 'bg-red-100 text-red-800'
      }
      if (action.includes('CREATE')) {
        return 'bg-green-100 text-green-800'
      }
      if (action.includes('UPDATE') || action.includes('CHANGED')) {
        return 'bg-blue-100 text-blue-800'
      }
      if (action.includes('LOGIN') || action.includes('LOGOUT')) {
        return 'bg-purple-100 text-purple-800'
      }
      return 'bg-gray-100 text-gray-800'
    }

    const getStatusBadgeClass = (status) => {
      switch (status?.toLowerCase()) {
        case 'success':
          return 'bg-green-100 text-green-800'
        case 'failure':
          return 'bg-red-100 text-red-800'
        case 'pending':
          return 'bg-yellow-100 text-yellow-800'
        default:
          return 'bg-gray-100 text-gray-800'
      }
    }

    onMounted(() => {
      fetchLogs()
      fetchStats()
    })

    return {
      logs,
      loading,
      showModal,
      selectedLog,
      stats,
      filters,
      pagination,
      fetchLogs,
      debouncedFetch,
      clearFilters,
      goToPage,
      viewDetails,
      closeModal,
      formatDate,
      formatAction,
      formatEntityType,
      formatJSON,
      getActionBadgeClass,
      getStatusBadgeClass
    }
  }
}
</script>
