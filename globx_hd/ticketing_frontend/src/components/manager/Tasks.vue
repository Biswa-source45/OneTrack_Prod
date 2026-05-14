<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-7xl mx-auto mt-8">
    <h1 class="text-2xl font-bold text-neutral-dark mb-6">All Tasks</h1>
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
          <div class="mb-3">
            <div class="flex-1 min-w-0">
              <h3 class="text-lg font-semibold text-neutral-dark truncate">
                {{ task.subject || 'No Subject' }}
              </h3>
            </div>
          </div>
          
          <div class="flex flex-wrap items-center justify-between gap-2">
            <!-- Left side: Task details in a single line -->
            <div class="flex flex-wrap items-center gap-1 text-sm text-neutral-dark">
              <!-- Created By -->
              <span>{{ getCreatorName(task) }}</span>
              <span class="text-brand-teal">•</span>
              
              <!-- Assigned To -->
              <span>{{ getAssignedUserName(task) }}</span>
              <span class="text-brand-teal">•</span>
              
              <!-- Due Date -->
              <span>{{ formatDueDate(task.due_date) }}</span>
              <span class="text-brand-teal">•</span>
              
              <!-- Date Created -->
              <span>{{ formatDate(task.created_at) }}</span>
            </div>
            
            <!-- Right side: Status, Priority, and Delete -->
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
              
              <!-- Delete Button -->
              <button 
                @click.stop="openDelete(task)"
                class="p-1.5 rounded-full hover:bg-red-100 transition-colors group"
                title="Delete Task"
              >
                <svg class="w-4 h-4 text-red-600 group-hover:text-red-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1-1H9a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-if="!tasks.length" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-brand-teal" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-neutral-dark">No tasks found</h3>
        <p class="mt-1 text-sm text-neutral-medium">Get started by creating a new task.</p>
      </div>
    </div>

    <!-- Edit Details Modal -->
    <Modal :open="showEdit" title="Edit Task Details" @close="closeEdit">
      <div class="mb-3">
        <label class="block mb-1 text-neutral-dark font-medium">Subject</label>
        <input v-model="editForm.subject" type="text" class="w-full border rounded px-3 py-2" />
      </div>
      <div class="mb-3">
        <label class="block mb-1 text-neutral-dark font-medium">Priority</label>
        <select v-model="editForm.priority" class="w-full border rounded px-3 py-2">
          <option value="High">High</option>
          <option value="Medium">Medium</option>
          <option value="Low">Low</option>
        </select>
      </div>
      <div class="mb-3">
        <label class="block mb-1 text-neutral-dark font-medium">Due Date</label>
        <input v-model="editForm.due_date" type="date" class="w-full border rounded px-3 py-2" />
      </div>
      <div class="mb-3">
        <label class="block mb-1 text-neutral-dark font-medium">Description</label>
        <textarea v-model="editForm.description" rows="3" class="w-full border rounded px-3 py-2"></textarea>
      </div>
      <div class="flex gap-2 justify-end">
        <Button size="sm" @click="updateTaskDetails">Save</Button>
        <Button size="sm" variant="secondary" @click="closeEdit">Cancel</Button>
      </div>
      <div v-if="editError" class="text-red-600 mt-2">{{ editError }}</div>
    </Modal>

    <!-- Assign User Modal -->
    <Modal :open="showAssign" @close="closeAssign">
      <div class="p-4 w-full max-w-md">
        <h2 class="text-lg font-semibold text-neutral-dark mb-2">Assign Task</h2>
        <select v-model="selectedUser" class="w-full border rounded px-3 py-2 mb-4">
          <option value="">Unassigned</option>
          <option v-for="user in users" :key="user.id" :value="String(user.id)">
            {{ formatUserName(user) }}
          </option>
        </select>
        <div class="flex gap-2 justify-end">
          <Button size="sm" @click="assignUserToTask">Assign</Button>
          <Button size="sm" variant="secondary" @click="closeAssign">Cancel</Button>
        </div>
        <div v-if="assignError" class="text-red-600 mt-2">{{ assignError }}</div>
      </div>
    </Modal>

    <!-- Change Status Modal -->
    <Modal :open="showStatus" @close="closeStatus">
      <div class="p-4 w-full max-w-md">
        <h2 class="text-lg font-semibold text-neutral-dark mb-2">Change Task Status</h2>
        <select v-model="selectedStatus" class="w-full border rounded px-3 py-2 mb-4">
          <option disabled value="">Please select status</option>
          <option v-for="status in statusOptions" :key="status" :value="status">{{ status }}</option>
        </select>
        <div class="flex gap-2 justify-end">
          <Button size="sm" @click="changeStatusOfTask">Change</Button>
          <Button size="sm" variant="secondary" @click="closeStatus">Cancel</Button>
        </div>
        <div v-if="statusError" class="text-red-600 mt-2">{{ statusError }}</div>
      </div>
    </Modal>

    <!-- Delete Confirmation Modal -->
    <Modal :open="showDelete" title="Delete Task" @close="closeDelete">
      <div class="p-4 w-full max-w-md">
        <div class="flex items-center mb-4">
          <svg class="w-12 h-12 text-red-600 mr-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.314 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <div>
            <h3 class="text-lg font-semibold text-red-800 mb-1">Delete Task</h3>
            <p class="text-sm text-gray-600">Are you sure you want to delete task <span class="font-semibold text-red-700">"{{ selectedTask?.subject }}"</span>?</p>
          </div>
        </div>
        <div class="flex gap-2 justify-end">
          <Button size="sm" variant="danger" @click="deleteTaskFromList">Delete</Button>
          <Button size="sm" variant="secondary" @click="closeDelete">Cancel</Button>
        </div>
        <div v-if="deleteError" class="text-red-600 mt-2 text-sm">{{ deleteError }}</div>
      </div>
    </Modal>

  </div>
</template>

<script setup>
import { formatDateIST } from '@/utils/date';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { fetchTasks, deleteTask, updateTask } from '@/api/tasks';
import { fetchUsers } from '@/api/tickets';
import { formatUserName } from '@/utils/user';
import Button from '@/components/ui/Button.vue';
import Modal from '@/components/ui/Modal.vue';

const router = useRouter();
const formatDate = formatDateIST;
const tasks = ref([]);
const users = ref([]);
const showEdit = ref(false);
const showAssign = ref(false);
const showStatus = ref(false);
const showDelete = ref(false);
const selectedTask = ref(null);
const selectedUser = ref('');
const selectedStatus = ref('');
const assignError = ref('');
const statusError = ref('');
const editError = ref('');
const deleteError = ref('');
const editForm = ref({ subject: '', priority: '', due_date: '', description: '' });
const statusOptions = [
  'Not Started',
  'In Progress',
  'Completed',
  'Deferred',
  'Waiting on someone else',
  'Canceled'
];

const showSuccess = ref(false);
const successMessage = ref('');

// Load tasks
async function loadTasks() {
  try {
    const response = await fetchTasks();
    tasks.value = response.tasks || [];
  } catch (e) {
    console.error('Failed to load tasks:', e);
  }
}

// Load users for assignment
async function loadUsers() {
  try {
    const response = await fetchUsers();
    users.value = Array.isArray(response) ? response : (response?.users || []);
  } catch (e) {
    console.error('Failed to load users:', e);
  }
}

// Edit task functions
async function openEdit(task) {
  selectedTask.value = task;
  editError.value = '';
  editForm.value = {
    subject: task.subject || '',
    priority: task.priority || '',
    due_date: task.due_date ? formatDateForInput(task.due_date) : '',
    description: task.description || ''
  };
  showEdit.value = true;
}

function closeEdit() {
  showEdit.value = false;
  selectedTask.value = null;
  editForm.value = { subject: '', priority: '', due_date: '', description: '' };
}

async function updateTaskDetails() {
  try {
    await updateTask(selectedTask.value.id, {
      subject: editForm.value.subject,
      priority: editForm.value.priority,
      due_date: editForm.value.due_date || null,
      description: editForm.value.description
    });
    showEdit.value = false;
    await loadTasks();
    successMessage.value = 'Task updated successfully!';
    showSuccess.value = true;
    setTimeout(() => { showSuccess.value = false; }, 2000);
  } catch (e) {
    editError.value = e.response?.data?.error || 'Failed to update task.';
  }
}

// Assign user functions
async function openAssign(task) {
  selectedTask.value = task;
  assignError.value = '';
  await loadUsers();
  selectedUser.value = task.assigned_to ? String(task.assigned_to) : '';
  showAssign.value = true;
}

function closeAssign() {
  showAssign.value = false;
  selectedTask.value = null;
  selectedUser.value = '';
}

async function assignUserToTask() {
  let userId = selectedUser.value ? Number(selectedUser.value) : null;
  try {
    await updateTask(selectedTask.value.id, {
      assigned_to: userId
    });
    showAssign.value = false;
    await loadTasks();
    successMessage.value = 'User assigned successfully!';
    showSuccess.value = true;
    setTimeout(() => { showSuccess.value = false; }, 2000);
  } catch (e) {
    assignError.value = e.response?.data?.error || 'Failed to assign user.';
  }
}

// Change status functions
function openStatus(task) {
  selectedTask.value = task;
  selectedStatus.value = task.task_status || '';
  statusError.value = '';
  showStatus.value = true;
}

function closeStatus() {
  showStatus.value = false;
  selectedTask.value = null;
  selectedStatus.value = '';
}

async function changeStatusOfTask() {
  if (!selectedStatus.value) {
    statusError.value = 'Please select a status.';
    return;
  }
  try {
    await updateTask(selectedTask.value.id, {
      task_status: selectedStatus.value
    });
    showStatus.value = false;
    await loadTasks();
    successMessage.value = 'Status changed successfully!';
    showSuccess.value = true;
    setTimeout(() => { showSuccess.value = false; }, 2000);
  } catch (e) {
    statusError.value = e.response?.data?.error || 'Failed to change status.';
  }
}

// Navigation
async function openDetail(task) {
  router.push(`/manager/tasks/${task.id}`);
}

// Delete task functions
function openDelete(task) {
  selectedTask.value = task;
  deleteError.value = '';
  showDelete.value = true;
}

function closeDelete() {
  showDelete.value = false;
  selectedTask.value = null;
  deleteError.value = '';
}

async function deleteTaskFromList() {
  if (!selectedTask.value) {
    deleteError.value = 'No task selected for deletion.';
    return;
  }
  
  try {
    // Call backend API to delete task
    await deleteTask(selectedTask.value.id);
    
    console.log(`🗑️ TASK DELETED: "${selectedTask.value.subject}" (ID: ${selectedTask.value.id})`);
    
    // Close modal
    showDelete.value = false;
    
    // Refresh task list from backend
    await loadTasks();
    
    // Show success message
    successMessage.value = `Task "${selectedTask.value.subject}" deleted successfully!`;
    showSuccess.value = true;
    setTimeout(() => { showSuccess.value = false; }, 3000);
    
    selectedTask.value = null;
    
  } catch (e) {
    console.error('Delete task error:', e);
    deleteError.value = e.response?.data?.error || 'Failed to delete task.';
  }
}

// Utility functions
function getStatusBadgeClass(status) {
  switch (status) {
    case 'Not Started': return 'bg-gray-400 text-white shadow-sm';
    case 'In Progress': return 'bg-brand-cyan text-white shadow-sm';
    case 'Completed': return 'bg-emerald-500 text-white shadow-sm';
    case 'Deferred': return 'bg-amber-500 text-white shadow-sm';
    case 'Waiting on someone else': return 'bg-purple-500 text-white shadow-sm';
    case 'Canceled': return 'bg-red-500 text-white shadow-sm';
    default: return 'bg-gray-400 text-white shadow-sm';
  }
}

function getPriorityBadgeClass(priority) {
  switch (priority) {
    case 'High': return 'bg-orange-500 text-white shadow-sm';
    case 'Medium': return 'bg-amber-500 text-white shadow-sm';
    case 'Low': return 'bg-gray-400 text-white shadow-sm';
    default: return 'bg-gray-400 text-white shadow-sm';
  }
}

function getAssignedUserName(task) {
  if (!task.assigned_user) return 'Unassigned';
  return formatUserName(task.assigned_user);
}

function getCreatorName(task) {
  if (!task.creator) return 'Unknown';
  return formatUserName(task.creator);
}

function formatDueDate(dateStr) {
  if (!dateStr) return 'No due date';
  return formatDate(dateStr);
}

function formatDateForInput(dateStr) {
  if (!dateStr) return '';
  return new Date(dateStr).toISOString().split('T')[0];
}

onMounted(() => {
  loadTasks();
});
</script>

<style scoped>
table {
  border-collapse: collapse;
}
th, td {
  border-bottom: 1px solid #e0e7ff;
}
</style>
