<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-7xl mx-auto mt-8">
    <h1 class="text-2xl font-bold text-blue-800 mb-6">My Assigned Tasks</h1>
    <div v-if="showSuccess" class="text-green-600 mt-2 text-center font-semibold">
      {{ successMessage }}
    </div>
    
    <!-- Tasks Card Layout -->
    <div class="space-y-4">
      <div v-for="task in tasks" :key="task.id" 
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
              <span>{{ getCreatorName(task) }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Assigned To -->
              <span>{{ getAssignedUserName(task) }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Due Date -->
              <span>{{ formatDueDate(task.due_date) }}</span>
              <span class="text-blue-400">•</span>
              
              <!-- Date Created -->
              <span>{{ formatDate(task.created_at) }}</span>
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
      <div v-if="!tasks.length" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-blue-900">No tasks assigned</h3>
        <p class="mt-1 text-sm text-blue-500">Assigned tasks will appear here.</p>
      </div>
    </div>

    <!-- Error Modal -->
    <Modal v-if="error" @close="error = ''">
      <div class="text-red-700 font-semibold">{{ error }}</div>
    </Modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { formatDateIST } from '../../utils/date';
import { formatUserName } from '../../utils/user';
import { fetchEngineerTasks } from '../../api/engineer';
import Modal from '../ui/Modal.vue';

const router = useRouter();
const formatDate = formatDateIST;
const tasks = ref([]);
const error = ref('');
const showSuccess = ref(false);
const successMessage = ref('');

// Helper functions (matching Manager implementation exactly)
function getStatusBadgeClass(status) {
  switch (status) {
    case 'TODO': return 'bg-gray-100 text-gray-800';
    case 'IN_PROGRESS': return 'bg-blue-100 text-blue-800';
    case 'COMPLETED': return 'bg-green-100 text-green-800';
    case 'ON_HOLD': return 'bg-yellow-100 text-yellow-800';
    case 'CANCELLED': return 'bg-red-100 text-red-800';
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

// Load tasks on component mount
onMounted(async () => {
  try {
    const result = await fetchEngineerTasks();
    tasks.value = Array.isArray(result.tasks) ? result.tasks : [];
  } catch (err) {
    error.value = err.message || 'Failed to fetch tasks.';
  }
});

// Navigate to task detail page (matching Manager behavior)
function openDetail(task) {
  router.push(`/engineer/tasks/${task.id}`);
}
</script>
