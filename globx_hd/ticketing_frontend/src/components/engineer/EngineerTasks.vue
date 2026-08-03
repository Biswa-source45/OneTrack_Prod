<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-7xl mx-auto mt-8">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-blue-800">My Assigned Tasks</h1>
      <button v-if="activeFilter" @click="setFilter(null)" class="text-sm text-blue-600 hover:text-blue-800 underline font-medium">
        Clear Filter (Show All {{ allTasks.length }})
      </button>
    </div>

    <div v-if="showSuccess" class="text-green-600 mt-2 text-center font-semibold mb-4">
      {{ successMessage }}
    </div>

    <!-- Ticket Statistics Row (Row 1) -->
    <div class="mb-6">
      <h2 class="text-lg font-bold text-gray-800 mb-3 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
        </svg>
        Assigned Ticket Overview
      </h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
        <!-- Total Assigned Tickets -->
        <div @click="router.push('/engineer/tickets')" class="bg-white rounded-lg shadow-sm border border-blue-100 p-4 cursor-pointer hover:shadow hover:border-blue-300 transition-all">
          <div class="flex items-center gap-3">
            <span class="p-3 bg-blue-100 rounded-full text-blue-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-blue-800">{{ stats.total_tickets }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">Total Tickets</div>
            </div>
          </div>
        </div>

        <!-- Open Tickets -->
        <div @click="router.push('/engineer/tickets?status=OPEN')" class="bg-white rounded-lg shadow-sm border border-orange-100 p-4 cursor-pointer hover:shadow hover:border-orange-300 transition-all">
          <div class="flex items-center gap-3">
            <span class="p-3 bg-orange-100 rounded-full text-orange-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-orange-800">{{ stats.open_tickets }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">Open Tickets</div>
            </div>
          </div>
        </div>

        <!-- In Progress Tickets -->
        <div @click="router.push('/engineer/tickets?status=IN_PROGRESS')" class="bg-white rounded-lg shadow-sm border border-yellow-100 p-4 cursor-pointer hover:shadow hover:border-yellow-300 transition-all">
          <div class="flex items-center gap-3">
            <span class="p-3 bg-yellow-100 rounded-full text-yellow-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-yellow-800">{{ stats.in_progress_tickets }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">In Progress</div>
            </div>
          </div>
        </div>

        <!-- Resolved Tickets -->
        <div @click="router.push('/engineer/tickets?status=RESOLVED')" class="bg-white rounded-lg shadow-sm border border-teal-100 p-4 cursor-pointer hover:shadow hover:border-teal-300 transition-all">
          <div class="flex items-center gap-3">
            <span class="p-3 bg-teal-100 rounded-full text-teal-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-teal-800">{{ stats.resolved_tickets }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">Resolved</div>
            </div>
          </div>
        </div>

        <!-- Closed Tickets -->
        <div @click="router.push('/engineer/tickets?status=CLOSED')" class="bg-white rounded-lg shadow-sm border border-gray-100 p-4 cursor-pointer hover:shadow hover:border-gray-300 transition-all">
          <div class="flex items-center gap-3">
            <span class="p-3 bg-gray-100 rounded-full text-gray-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </span>
            <div>
              <div class="text-2xl font-bold text-gray-800">{{ stats.closed_tickets }}</div>
              <div class="text-xs font-medium text-gray-500 uppercase">Closed</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Task Statistics Row (Row 2) -->
    <div class="mb-6">
      <h2 class="text-lg font-bold text-gray-800 mb-3 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
        Assigned Task Overview
      </h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Total Tasks Card -->
        <div 
          @click="setFilter(null)"
          :class="['p-4 rounded-lg border transition-all cursor-pointer shadow-sm', activeFilter === null ? 'bg-blue-50 border-blue-400 ring-2 ring-blue-300' : 'bg-white border-gray-200 hover:border-blue-300 hover:shadow']"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Total Tasks</p>
              <h3 class="text-2xl font-bold text-blue-900 mt-1">{{ totalCount }}</h3>
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
          @click="setFilter('IN_PROGRESS')"
          :class="['p-4 rounded-lg border transition-all cursor-pointer shadow-sm', activeFilter === 'IN_PROGRESS' ? 'bg-yellow-50 border-yellow-400 ring-2 ring-yellow-300' : 'bg-white border-gray-200 hover:border-yellow-300 hover:shadow']"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">In Progress / Active</p>
              <h3 class="text-2xl font-bold text-yellow-800 mt-1">{{ inProgressCount }}</h3>
            </div>
            <div class="p-3 bg-yellow-100 rounded-full text-yellow-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
          </div>
          <p class="text-xs text-gray-500 mt-2">Not Started or In Progress</p>
        </div>

        <!-- Manager Verification Pending Card -->
        <div 
          @click="setFilter('REVIEW_PENDING')"
          :class="['p-4 rounded-lg border transition-all cursor-pointer shadow-sm', activeFilter === 'REVIEW_PENDING' ? 'bg-purple-50 border-purple-400 ring-2 ring-purple-300' : 'bg-white border-gray-200 hover:border-purple-300 hover:shadow']"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Manager Verification</p>
              <h3 class="text-2xl font-bold text-purple-800 mt-1">{{ reviewPendingCount }}</h3>
            </div>
            <div class="p-3 bg-purple-100 rounded-full text-purple-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
          </div>
          <p class="text-xs text-gray-500 mt-2">Awaiting Manager Approval</p>
        </div>

        <!-- Completed Card -->
        <div 
          @click="setFilter('COMPLETED')"
          :class="['p-4 rounded-lg border transition-all cursor-pointer shadow-sm', activeFilter === 'COMPLETED' ? 'bg-green-50 border-green-400 ring-2 ring-green-300' : 'bg-white border-gray-200 hover:border-green-300 hover:shadow']"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Completed Tasks</p>
              <h3 class="text-2xl font-bold text-green-800 mt-1">{{ completedCount }}</h3>
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
    
    <!-- Tasks Card Layout -->
    <div class="space-y-4">
      <div v-for="task in displayedTasks" :key="task.id" 
           class="bg-white border border-blue-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer"
           @click="openDetail(task)">
        
        <!-- Card Content -->
        <div class="p-4">
          <div class="flex items-start justify-between mb-3">
            <div class="flex-1 min-w-0">
              <h3 class="text-lg font-semibold text-blue-900 truncate">
                {{ task.subject || 'No Subject' }}
              </h3>
            </div>
          </div>
          
          <div class="flex flex-wrap items-center justify-between gap-2">
            <!-- Left side: Task details in a single line -->
            <div class="flex flex-wrap items-center gap-1 text-sm text-blue-900">
              <!-- Created By -->
              <span>Created: {{ getCreatorName(task) }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Assigned To -->
              <span>Assigned: {{ getAssignedUserName(task) }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Due Date -->
              <span>Due: {{ formatDueDate(task.due_date) }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Date Created -->
              <span>Date: {{ formatDate(task.created_at) }}</span>
            </div>
            
            <!-- Right side: Status and Priority -->
            <div class="flex items-center gap-3">
              <!-- Status Badge -->
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" 
                    :class="getStatusBadgeClass(task.task_status)">
                {{ task.task_status }}
              </span>
              
              <!-- Priority Badge -->
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="getPriorityBadgeClass(task.priority)">
                {{ task.priority || 'Medium' }}
              </span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-if="!displayedTasks.length" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-blue-900">No tasks found</h3>
        <p class="mt-1 text-sm text-blue-500">
          {{ activeFilter ? 'No tasks match the selected filter criteria.' : 'Assigned tasks will appear here.' }}
        </p>
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
import { formatUserName } from '../../utils/user';
import { fetchEngineerTasks, fetchEngineerDashboardStats } from '../../api/engineer';
import Modal from '../ui/Modal.vue';

const router = useRouter();
const route = useRoute();
const formatDate = formatDateIST;
const allTasks = ref([]);
const activeFilter = ref(null);
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

// Computed count metrics
const totalCount = computed(() => allTasks.value.length);

const inProgressCount = computed(() => {
  return allTasks.value.filter(t => {
    const s = (t.task_status || '').toUpperCase();
    return s === 'NOT STARTED' || s === 'TODO' || s === 'IN PROGRESS' || s === 'IN_PROGRESS' || s === 'ON_HOLD' || s === 'DEFERRED';
  }).length;
});

const reviewPendingCount = computed(() => {
  return allTasks.value.filter(t => {
    const s = (t.task_status || '').toUpperCase();
    return s.includes('REVIEW') || s.includes('PENDING') || s.includes('VERIFICATION');
  }).length;
});

const completedCount = computed(() => {
  return allTasks.value.filter(t => {
    const s = (t.task_status || '').toUpperCase();
    return s === 'COMPLETED' || s === 'ACCEPTED';
  }).length;
});

// Sorted tasks (Newest to Oldest) and filtered by active filter
const displayedTasks = computed(() => {
  let filtered = [...allTasks.value];

  if (activeFilter.value === 'IN_PROGRESS') {
    filtered = filtered.filter(t => {
      const s = (t.task_status || '').toUpperCase();
      return s === 'NOT STARTED' || s === 'TODO' || s === 'IN PROGRESS' || s === 'IN_PROGRESS' || s === 'ON_HOLD' || s === 'DEFERRED';
    });
  } else if (activeFilter.value === 'REVIEW_PENDING') {
    filtered = filtered.filter(t => {
      const s = (t.task_status || '').toUpperCase();
      return s.includes('REVIEW') || s.includes('PENDING') || s.includes('VERIFICATION');
    });
  } else if (activeFilter.value === 'COMPLETED') {
    filtered = filtered.filter(t => {
      const s = (t.task_status || '').toUpperCase();
      return s === 'COMPLETED' || s === 'ACCEPTED';
    });
  }

  // Sort from newest to oldest by created_at or id
  return filtered.sort((a, b) => {
    const dateA = new Date(a.created_at || 0).getTime();
    const dateB = new Date(b.created_at || 0).getTime();
    if (dateA !== dateB) return dateB - dateA;
    return (b.id || 0) - (a.id || 0);
  });
});

function setFilter(filterType) {
  activeFilter.value = filterType;
}

// Helper functions (matching Manager implementation exactly)
function getStatusBadgeClass(status) {
  const s = (status || '').toUpperCase();
  if (s === 'TODO' || s === 'NOT STARTED') return 'bg-gray-100 text-gray-800';
  if (s === 'IN PROGRESS' || s === 'IN_PROGRESS') return 'bg-blue-100 text-blue-800';
  if (s.includes('REVIEW') || s.includes('PENDING') || s.includes('VERIFICATION')) return 'bg-purple-100 text-purple-800';
  if (s === 'COMPLETED' || s === 'ACCEPTED') return 'bg-green-100 text-green-800';
  if (s === 'ON_HOLD' || s === 'DEFERRED') return 'bg-yellow-100 text-yellow-800';
  if (s === 'CANCELLED') return 'bg-red-100 text-red-800';
  return 'bg-gray-100 text-gray-800';
}

function getPriorityBadgeClass(priority) {
  switch (priority) {
    case 'High': return 'bg-red-100 text-red-800';
    case 'Medium': return 'bg-yellow-100 text-yellow-800';
    case 'Low': return 'bg-green-100 text-green-800';
    default: return 'bg-gray-100 text-gray-800';
  }
}

function getCreatorName(task) {
  if (!task.creator) return 'Unknown';
  return formatUserName(task.creator);
}

function getAssignedUserName(task) {
  if (!task.assigned_user) return 'Unassigned';
  return formatUserName(task.assigned_user);
}

function formatDueDate(dueDate) {
  if (!dueDate) return 'No due date';
  return formatDate(dueDate);
}

// Load tasks and stats on component mount
onMounted(async () => {
  try {
    const result = await fetchEngineerTasks();
    allTasks.value = Array.isArray(result.tasks) ? result.tasks : [];

    if (route.query.status) {
      const q = route.query.status.toUpperCase();
      if (q.includes('REVIEW') || q.includes('PENDING')) activeFilter.value = 'REVIEW_PENDING';
      else if (q === 'COMPLETED') activeFilter.value = 'COMPLETED';
      else if (q === 'NOT_COMPLETED' || q === 'IN_PROGRESS') activeFilter.value = 'IN_PROGRESS';
    }
  } catch (err) {
    error.value = err.message || 'Failed to fetch tasks.';
  }

  try {
    const statsData = await fetchEngineerDashboardStats();
    if (statsData) {
      stats.value = statsData;
    }
  } catch (err) {
    console.error('Failed to load dashboard stats:', err);
  }
});

watch(() => route.query.status, (newStatus) => {
  if (!newStatus) {
    activeFilter.value = null;
  } else {
    const q = newStatus.toUpperCase();
    if (q.includes('REVIEW') || q.includes('PENDING')) activeFilter.value = 'REVIEW_PENDING';
    else if (q === 'COMPLETED') activeFilter.value = 'COMPLETED';
    else if (q === 'NOT_COMPLETED' || q === 'IN_PROGRESS') activeFilter.value = 'IN_PROGRESS';
  }
});

// Navigate to task detail page (matching Manager behavior)
function openDetail(task) {
  router.push(`/engineer/tasks/${task.id}`);
}
</script>
