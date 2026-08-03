<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-7xl mx-auto mt-8">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-blue-800">Engineer Dashboard & Tickets</h1>
      <button v-if="activeTicketFilter" @click="setTicketFilter(null)" class="text-sm text-blue-600 hover:text-blue-800 underline font-medium">
        Clear Filter (Show All {{ allTickets.length }})
      </button>
    </div>

    <div v-if="showSuccess" class="text-green-600 mt-2 text-center font-semibold mb-4">
      {{ successMessage }}
    </div>

    <!-- Ticket Statistics Row (Row 1) -->
    <div class="mb-8">
      <h2 class="text-lg font-bold text-gray-800 mb-3 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
        </svg>
        Assigned Ticket Overview
      </h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
        <!-- Total Assigned Tickets -->
        <div 
          @click="setTicketFilter(null)" 
          :class="['rounded-lg shadow-sm p-4 cursor-pointer border transition-all', activeTicketFilter === null ? 'bg-blue-50 border-blue-400 ring-2 ring-blue-300' : 'bg-white border-blue-100 hover:shadow hover:border-blue-300']"
        >
          <div class="flex items-center gap-3">
            <span class="p-3 bg-blue-100 rounded-full text-blue-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-blue-800">{{ stats.total_tickets || 0 }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">Total Tickets</div>
            </div>
          </div>
        </div>

        <!-- Open Tickets -->
        <div 
          @click="setTicketFilter('OPEN')" 
          :class="['rounded-lg shadow-sm p-4 cursor-pointer border transition-all', activeTicketFilter === 'OPEN' ? 'bg-orange-50 border-orange-400 ring-2 ring-orange-300' : 'bg-white border-orange-100 hover:shadow hover:border-orange-300']"
        >
          <div class="flex items-center gap-3">
            <span class="p-3 bg-orange-100 rounded-full text-orange-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-orange-800">{{ stats.open_tickets || 0 }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">Open Tickets</div>
            </div>
          </div>
        </div>

        <!-- In Progress Tickets -->
        <div 
          @click="setTicketFilter('IN_PROGRESS')" 
          :class="['rounded-lg shadow-sm p-4 cursor-pointer border transition-all', activeTicketFilter === 'IN_PROGRESS' ? 'bg-yellow-50 border-yellow-400 ring-2 ring-yellow-300' : 'bg-white border-yellow-100 hover:shadow hover:border-yellow-300']"
        >
          <div class="flex items-center gap-3">
            <span class="p-3 bg-yellow-100 rounded-full text-yellow-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-yellow-800">{{ stats.in_progress_tickets || 0 }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">In Progress</div>
            </div>
          </div>
        </div>

        <!-- Resolved Tickets -->
        <div 
          @click="setTicketFilter('RESOLVED')" 
          :class="['rounded-lg shadow-sm p-4 cursor-pointer border transition-all', activeTicketFilter === 'RESOLVED' ? 'bg-teal-50 border-teal-400 ring-2 ring-teal-300' : 'bg-white border-teal-100 hover:shadow hover:border-teal-300']"
        >
          <div class="flex items-center gap-3">
            <span class="p-3 bg-teal-100 rounded-full text-teal-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-teal-800">{{ stats.resolved_tickets || 0 }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">Resolved</div>
            </div>
          </div>
        </div>

        <!-- Closed Tickets -->
        <div 
          @click="setTicketFilter('CLOSED')" 
          :class="['rounded-lg shadow-sm p-4 cursor-pointer border transition-all', activeTicketFilter === 'CLOSED' ? 'bg-gray-100 border-gray-400 ring-2 ring-gray-300' : 'bg-white border-gray-100 hover:shadow hover:border-gray-300']"
        >
          <div class="flex items-center gap-3">
            <span class="p-3 bg-gray-100 rounded-full text-gray-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-gray-800">{{ stats.closed_tickets || 0 }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">Closed</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Task Statistics Row (Row 2) -->
    <div class="mb-8">
      <h2 class="text-lg font-bold text-gray-800 mb-3 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
        Assigned Task Overview
      </h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Total Tasks Card -->
        <div 
          @click="router.push('/engineer/tasks')"
          class="bg-white rounded-lg shadow-sm border border-blue-100 p-4 cursor-pointer hover:shadow hover:border-blue-300 transition-all"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Total Tasks</p>
              <h3 class="text-2xl font-bold text-blue-900 mt-1">{{ stats.total_tasks || 0 }}</h3>
            </div>
            <div class="p-3 bg-blue-100 rounded-full text-blue-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
          </div>
          <p class="text-xs text-gray-500 mt-2">All tasks assigned to you</p>
        </div>

        <!-- In Progress Tasks Card -->
        <div 
          @click="router.push('/engineer/tasks?status=IN_PROGRESS')"
          class="bg-white rounded-lg shadow-sm border border-yellow-100 p-4 cursor-pointer hover:shadow hover:border-yellow-300 transition-all"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">In Progress / Active</p>
              <h3 class="text-2xl font-bold text-yellow-800 mt-1">{{ stats.in_progress_tasks || 0 }}</h3>
            </div>
            <div class="p-3 bg-yellow-100 rounded-full text-yellow-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
          </div>
          <p class="text-xs text-gray-500 mt-2">Active tasks working on</p>
        </div>

        <!-- Manager Verification Pending Card -->
        <div 
          @click="router.push('/engineer/tasks?status=REVIEW_PENDING')"
          class="bg-white rounded-lg shadow-sm border border-purple-100 p-4 cursor-pointer hover:shadow hover:border-purple-300 transition-all"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Manager Verification</p>
              <h3 class="text-2xl font-bold text-purple-800 mt-1">{{ stats.review_pending_tasks || 0 }}</h3>
            </div>
            <div class="p-3 bg-purple-100 rounded-full text-purple-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
          </div>
          <p class="text-xs text-gray-500 mt-2">Awaiting Manager Verification</p>
        </div>

        <!-- Completed Card -->
        <div 
          @click="router.push('/engineer/tasks?status=COMPLETED')"
          class="bg-white rounded-lg shadow-sm border border-green-100 p-4 cursor-pointer hover:shadow hover:border-green-300 transition-all"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Completed Tasks</p>
              <h3 class="text-2xl font-bold text-green-800 mt-1">{{ stats.completed_tasks || 0 }}</h3>
            </div>
            <div class="p-3 bg-green-100 rounded-full text-green-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
          </div>
          <p class="text-xs text-gray-500 mt-2">Verified & Completed</p>
        </div>
      </div>
    </div>
    
    <!-- Tickets List -->
    <div class="space-y-4">
      <div v-for="ticket in displayedTickets" :key="ticket.id" 
           class="bg-white border border-blue-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer"
           @click="openDetail(ticket)">
        
        <!-- Card Content -->
        <div class="p-4">
          <div class="flex items-start justify-between mb-3">
            <div class="flex-1 min-w-0">
              <h3 class="text-lg font-semibold text-blue-900 truncate">
                {{ ticket.subject || 'No Subject' }}
              </h3>
            </div>
          </div>
          
          <div class="flex flex-wrap items-center justify-between gap-2">
            <!-- Left side: Ticket details in a single line -->
            <div class="flex flex-wrap items-center gap-1 text-sm text-blue-900">
              <!-- Ticket ID -->
              <span class="font-mono">{{ ticket.ticket_id }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Contact Name -->
              <span>{{ ticket.contact ? `${ticket.contact.first_name} ${ticket.contact.last_name}` : 'Unknown Contact' }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Account Name -->
              <span>{{ ticket.contact?.account?.account_name || ticket.account?.account_name || 'Unknown Account' }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Date Raised -->
              <span>{{ formatDate(ticket.created_at) }}</span>
            </div>
            
            <!-- Right side: Status and Priority -->
            <div class="flex items-center gap-3">
              <!-- Status Badge -->
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" 
                    :class="getStatusBadgeClass(ticket.ticket_status)">
                {{ ticket.ticket_status }}
              </span>
              
              <!-- Priority Badge -->
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="getPriorityBadgeClass(ticket.priority)">
                {{ ticket.priority || 'Medium' }}
              </span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-if="!displayedTickets.length" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-blue-900">No tickets assigned</h3>
        <p class="mt-1 text-sm text-blue-500">Assigned tickets will appear here.</p>
      </div>
    </div>

    <!-- Error Modal -->
    <Modal v-if="error" @close="error = ''">
      <div class="text-red-700 font-semibold">{{ error }}</div>
    </Modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { formatDateIST } from '../../utils/date';
import Modal from '../ui/Modal.vue';
import { fetchEngineerTickets, fetchEngineerDashboardStats, fetchEngineerTasks } from '../../api/engineer';

const router = useRouter();
const route = useRoute();
const formatDate = formatDateIST;
const allTickets = ref([]);
const activeTicketFilter = ref(null);
const error = ref('');
const showSuccess = ref(false);
const successMessage = ref('');

const stats = ref({
  total_tickets: 0,
  open_tickets: 0,
  in_progress_tickets: 0,
  resolved_tickets: 0,
  closed_tickets: 0,
  total_tasks: 0,
  in_progress_tasks: 0,
  review_pending_tasks: 0,
  completed_tasks: 0,
});

const displayedTickets = computed(() => {
  if (!activeTicketFilter.value) return allTickets.value;
  return allTickets.value.filter(t => {
    const s = (t.ticket_status || '').toUpperCase();
    const filter = activeTicketFilter.value.toUpperCase();
    if (filter === 'IN_PROGRESS') {
      return s === 'IN PROGRESS' || s === 'IN_PROGRESS';
    }
    return s === filter;
  });
});

function setTicketFilter(filterType) {
  activeTicketFilter.value = filterType;
}

// Badge styling functions (matching Manager implementation exactly)
function getStatusBadgeClass(status) {
  switch (status) {
    case 'OPEN': return 'bg-blue-100 text-blue-800';
    case 'MEETING LOCKED IN WITH OEM': return 'bg-indigo-100 text-indigo-800';
    case 'PARTS ORDERED': return 'bg-purple-100 text-purple-800';
    case 'IN PROGRESS':
    case 'IN_PROGRESS': return 'bg-yellow-100 text-yellow-800';
    case 'RESOLVED': return 'bg-green-100 text-green-800';
    case 'CLOSED': return 'bg-gray-100 text-gray-800';
    case 'ON HOLD': return 'bg-orange-100 text-orange-800';
    case 'ESCALATED': return 'bg-red-100 text-red-800';
    default: return 'bg-gray-100 text-gray-800';
  }
}

function getPriorityBadgeClass(priority) {
  switch (priority) {
    case 'High': return 'bg-red-100 text-red-800';
    case 'Medium': return 'bg-yellow-100 text-yellow-800';
    case 'Low': return 'bg-green-100 text-green-800';
    default: return 'bg-gray-100 text-gray-800';
  }
}

// Load tickets and stats on component mount
onMounted(async () => {
  try {
    const result = await fetchEngineerTickets();
    allTickets.value = Array.isArray(result.tickets) ? result.tickets : [];

    if (route.query.status) {
      activeTicketFilter.value = route.query.status.toUpperCase();
    }
  } catch (err) {
    error.value = err.message || 'Failed to fetch tickets.';
  }

  try {
    const statsData = await fetchEngineerDashboardStats();
    if (statsData) {
      stats.value = statsData;
    }
  } catch (err) {
    console.error('Failed to load dashboard stats:', err);
  }

  // Fallback count calculation if stats total tickets is 0 but local tickets exist
  if (allTickets.value.length > 0 && (!stats.value.total_tickets || Number(stats.value.total_tickets) === 0)) {
    stats.value.total_tickets = allTickets.value.length;
    stats.value.open_tickets = allTickets.value.filter(t => (t.ticket_status || '').toUpperCase() === 'OPEN').length;
    stats.value.in_progress_tickets = allTickets.value.filter(t => {
      const s = (t.ticket_status || '').toUpperCase();
      return s === 'IN PROGRESS' || s === 'IN_PROGRESS';
    }).length;
    stats.value.resolved_tickets = allTickets.value.filter(t => (t.ticket_status || '').toUpperCase() === 'RESOLVED').length;
    stats.value.closed_tickets = allTickets.value.filter(t => (t.ticket_status || '').toUpperCase() === 'CLOSED').length;
  }

  // Fetch tasks and ensure task metrics are populated
  try {
    const tasksRes = await fetchEngineerTasks();
    const tasks = Array.isArray(tasksRes?.tasks) ? tasksRes.tasks : (Array.isArray(tasksRes) ? tasksRes : []);
    if (tasks.length > 0 && (!stats.value.total_tasks || Number(stats.value.total_tasks) === 0)) {
      stats.value.total_tasks = tasks.length;
      stats.value.in_progress_tasks = tasks.filter(tk => {
        const s = (tk.task_status || '').toUpperCase();
        return s === 'IN PROGRESS' || s === 'IN_PROGRESS' || s === 'TODO' || s === 'NOT STARTED';
      }).length;
      stats.value.review_pending_tasks = tasks.filter(tk => {
        const s = (tk.task_status || '').toUpperCase();
        return s.includes('REVIEW') || s.includes('PENDING') || s.includes('VERIFICATION');
      }).length;
      stats.value.completed_tasks = tasks.filter(tk => {
        const s = (tk.task_status || '').toUpperCase();
        return s === 'COMPLETED' || s === 'ACCEPTED';
      }).length;
    }
  } catch (err) {
    console.error('Failed to load tasks for dashboard stats:', err);
  }
});

watch(() => route.query.status, (newStatus) => {
  if (!newStatus) {
    activeTicketFilter.value = null;
  } else {
    activeTicketFilter.value = newStatus.toUpperCase();
  }
});

// Navigate to ticket detail page (matching Manager behavior)
function openDetail(ticket) {
  router.push(`/engineer/tickets/${ticket.id}`);
}
</script>
