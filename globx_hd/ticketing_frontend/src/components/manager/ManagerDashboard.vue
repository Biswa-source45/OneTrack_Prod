<template>
  <div class="p-8">
    <h1 class="text-3xl font-bold text-blue-800 mb-8">Manager Dashboard</h1>
    
    <!-- Month Filter -->
    <div class="mb-8">
      <div class="bg-white rounded-lg shadow-md border border-blue-100 p-6">
        <h2 class="text-lg font-semibold text-blue-800 mb-4">Filter by Month</h2>
        <div class="flex gap-4 items-center">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Month</label>
            <select 
              v-model="selectedMonth" 
              @change="loadStats"
              class="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400"
            >
              <option value="">All Months</option>
              <option v-for="month in months" :key="month.value" :value="month.value">
                {{ month.label }}
              </option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Year</label>
            <select 
              v-model="selectedYear" 
              @change="loadStats"
              class="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400"
            >
              <option v-for="year in years" :key="year" :value="year">
                {{ year }}
              </option>
            </select>
          </div>
          <div class="mt-6">
            <button 
              @click="clearFilter"
              class="px-4 py-2 text-sm text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded-md transition-colors"
            >
              Clear Filter
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Ticket Statistics Row -->
    <div class="mb-8">
      <h2 class="text-xl font-bold text-gray-800 mb-4 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
        </svg>
        Ticket Statistics
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6">
        <!-- Total Tickets -->
        <div @click="navigateToTickets()" class="bg-white rounded-lg shadow-md border border-blue-100 p-6 cursor-pointer hover:shadow-lg hover:border-blue-300 transition-all duration-200">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-blue-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-blue-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-blue-800">{{ stats.total_tickets || 0 }}</div>
              <div class="text-sm text-gray-600">Total Tickets</div>
            </div>
          </div>
        </div>

        <!-- Open Tickets -->
        <div @click="navigateToTickets('OPEN')" class="bg-white rounded-lg shadow-md border border-orange-100 p-6 cursor-pointer hover:shadow-lg hover:border-orange-300 transition-all duration-200">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-orange-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-orange-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-orange-800">{{ stats.open_tickets || 0 }}</div>
              <div class="text-sm text-gray-600">Open Tickets</div>
            </div>
          </div>
        </div>

        <!-- In Progress Tickets -->
        <div @click="navigateToTickets('IN PROGRESS')" class="bg-white rounded-lg shadow-md border border-yellow-100 p-6 cursor-pointer hover:shadow-lg hover:border-yellow-300 transition-all duration-200">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-yellow-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-yellow-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-yellow-800">{{ stats.in_progress_tickets || 0 }}</div>
              <div class="text-sm text-gray-600">In Progress</div>
            </div>
          </div>
        </div>

        <!-- Resolved Tickets -->
        <div @click="navigateToTickets('RESOLVED')" class="bg-white rounded-lg shadow-md border border-teal-100 p-6 cursor-pointer hover:shadow-lg hover:border-teal-300 transition-all duration-200">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-teal-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-teal-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-teal-800">{{ stats.resolved_tickets || 0 }}</div>
              <div class="text-sm text-gray-600">Resolved</div>
            </div>
          </div>
        </div>

        <!-- Closed Tickets -->
        <div @click="navigateToTickets('CLOSED')" class="bg-white rounded-lg shadow-md border border-gray-100 p-6 cursor-pointer hover:shadow-lg hover:border-gray-300 transition-all duration-200">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-gray-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-gray-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-gray-800">{{ stats.closed_tickets || 0 }}</div>
              <div class="text-sm text-gray-600">Closed Tickets</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Task Statistics Row -->
    <div class="mb-8">
      <h2 class="text-xl font-bold text-gray-800 mb-4 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
        </svg>
        Task Statistics
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- Total Tasks -->
        <div @click="navigateToTasks()" class="bg-white rounded-lg shadow-md border border-indigo-100 p-6 cursor-pointer hover:shadow-lg hover:border-indigo-300 transition-all duration-200">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-indigo-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-indigo-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-indigo-800">{{ stats.total_tasks || 0 }}</div>
              <div class="text-sm text-gray-600">Total Tasks</div>
            </div>
          </div>
        </div>

        <!-- In Progress / Not Completed -->
        <div @click="navigateToTasks('NOT_COMPLETED')" class="bg-white rounded-lg shadow-md border border-amber-100 p-6 cursor-pointer hover:shadow-lg hover:border-amber-300 transition-all duration-200">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-amber-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-amber-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-amber-800">{{ stats.not_completed_tasks || 0 }}</div>
              <div class="text-sm text-gray-600">In Progress / Active</div>
            </div>
          </div>
        </div>

        <!-- Manager Verification Pending -->
        <div @click="navigateToTasks('REVIEW_PENDING')" class="bg-white rounded-lg shadow-md border border-purple-100 p-6 cursor-pointer hover:shadow-lg hover:border-purple-300 transition-all duration-200">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-purple-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-purple-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-purple-800">{{ stats.review_pending_tasks || 0 }}</div>
              <div class="text-sm text-gray-600">Manager Verification Pending</div>
            </div>
          </div>
        </div>

        <!-- Completed Tasks -->
        <div @click="navigateToTasks('COMPLETED')" class="bg-white rounded-lg shadow-md border border-emerald-100 p-6 cursor-pointer hover:shadow-lg hover:border-emerald-300 transition-all duration-200">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-emerald-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-emerald-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-emerald-800">{{ stats.completed_tasks || 0 }}</div>
              <div class="text-sm text-gray-600">Completed (Accepted)</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Access Links -->
    <div class="bg-white rounded-lg shadow-md border border-blue-100 p-6">
      <h2 class="text-lg font-semibold text-blue-800 mb-4">Quick Access</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <router-link to="/manager/tickets" class="flex items-center gap-3 p-4 rounded-lg hover:bg-blue-50 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
          </svg>
          <span class="font-medium text-blue-800">Manage Tickets</span>
        </router-link>
        <router-link to="/manager/tasks" class="flex items-center gap-3 p-4 rounded-lg hover:bg-blue-50 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
          </svg>
          <span class="font-medium text-blue-800">Manage Tasks</span>
        </router-link>
        <router-link to="/manager/engineers" class="flex items-center gap-3 p-4 rounded-lg hover:bg-blue-50 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a4 4 0 00-3-3.87M9 20H4v-2a4 4 0 013-3.87m6-3.13a4 4 0 11-8 0 4 4 0 018 0zm6 0a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          <span class="font-medium text-blue-800">Manage Engineers</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { fetchDashboardStats } from '../../api/tickets.js'

const router = useRouter()

// Reactive data
const stats = ref({
  total_tickets: 0,
  open_tickets: 0,
  closed_tickets: 0,
  in_progress_tickets: 0,
  resolved_tickets: 0,
  total_tasks: 0,
  not_completed_tasks: 0,
  review_pending_tasks: 0,
  completed_tasks: 0
})

const selectedMonth = ref('')
const selectedYear = ref(new Date().getFullYear())
const loading = ref(false)

// Month options
const months = [
  { value: 1, label: 'January' },
  { value: 2, label: 'February' },
  { value: 3, label: 'March' },
  { value: 4, label: 'April' },
  { value: 5, label: 'May' },
  { value: 6, label: 'June' },
  { value: 7, label: 'July' },
  { value: 8, label: 'August' },
  { value: 9, label: 'September' },
  { value: 10, label: 'October' },
  { value: 11, label: 'November' },
  { value: 12, label: 'December' }
]

// Generate year options (current year and 2 years back)
const currentYear = new Date().getFullYear()
const years = [currentYear, currentYear - 1, currentYear - 2]

// Load statistics
const loadStats = async () => {
  try {
    loading.value = true
    const month = selectedMonth.value || null
    const year = selectedYear.value || null
    
    const data = await fetchDashboardStats(month, year)
    stats.value = data
  } catch (error) {
    console.error('Failed to load dashboard statistics:', error)
  } finally {
    loading.value = false
  }
}

// Clear filter
const clearFilter = () => {
  selectedMonth.value = ''
  selectedYear.value = currentYear
  loadStats()
}

// Navigate to tickets page with optional status filter
const navigateToTickets = (status = null) => {
  if (status) {
    router.push({ path: '/manager/tickets', query: { status } })
  } else {
    router.push('/manager/tickets')
  }
}

// Navigate to tasks page with optional status filter
const navigateToTasks = (status = null) => {
  if (status) {
    router.push({ path: '/manager/tasks', query: { status } })
  } else {
    router.push('/manager/tasks')
  }
}

// Load stats on component mount
onMounted(() => {
  loadStats()
})
</script>

<style scoped>
/* No custom styles needed, uses Tailwind/theme classes */
</style>
