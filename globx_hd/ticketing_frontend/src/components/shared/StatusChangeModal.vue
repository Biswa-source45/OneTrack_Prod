<template>
  <Modal :open="open" :title="title" @close="$emit('close')">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
        <p class="text-sm text-red-700">{{ errorMessage }}</p>
      </div>
      
      <!-- Status Selection (hidden for approval actions) -->
      <FormField v-if="!hideStatusDropdown" label="New Status" required>
        <select 
          v-model="form.status" 
          class="w-full px-3 py-2 border border-blue-200 rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
          required
        >
          <option value="" disabled>Select Status</option>
          <option v-for="option in statusOptions" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>
      </FormField>

      <!-- Remarks -->
      <FormField label="Remarks" required>
        <textarea 
          v-model="form.remarks" 
          rows="4" 
          class="w-full px-3 py-2 border border-blue-200 rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
          placeholder="Please provide remarks for this status change..."
          required
        ></textarea>
        <p class="text-xs text-gray-500 mt-1">Remarks are required for all status changes and will be visible in the ticket history.</p>
      </FormField>

      <!-- Current Status Info -->
      <div class="p-3 bg-gray-50 rounded-md">
        <p class="text-sm text-gray-600">
          <span class="font-medium">Current Status:</span> 
          <span class="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" 
                :class="getStatusBadgeClass(currentStatus)">
            {{ currentStatus }}
          </span>
        </p>
      </div>

      <!-- Actions -->
      <div class="flex justify-end space-x-3 pt-4">
        <button
          type="button"
          @click="$emit('close')"
          class="px-4 py-2 border border-gray-300 rounded text-gray-700 hover:bg-gray-50 transition"
        >
          Cancel
        </button>
        <button
          type="submit"
          :disabled="submitting || (!hideStatusDropdown && !form.status) || !form.remarks.trim()"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition"
        >
          {{ submitting ? 'Processing...' : (hideStatusDropdown ? 'Submit' : 'Change Status') }}
        </button>
      </div>
    </form>
  </Modal>
</template>

<script setup>
import { ref, reactive, watch } from 'vue';
import Modal from '../ui/Modal.vue';
import FormField from '../ui/FormField.vue';

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  currentStatus: {
    type: String,
    required: true
  },
  statusOptions: {
    type: Array,
    required: true
  },
  preSelectedStatus: {
    type: String,
    default: ''
  },
  hideStatusDropdown: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Change Status'
  }
});

const emit = defineEmits(['close', 'change']);

const submitting = ref(false);
const errorMessage = ref('');

const form = reactive({
  status: '',
  remarks: ''
});

// Reset form when modal opens/closes
watch(() => props.open, (isOpen) => {
  if (isOpen) {
    form.status = props.preSelectedStatus || '';
    form.remarks = '';
    errorMessage.value = '';
  }
});

// Watch for preSelectedStatus changes
watch(() => props.preSelectedStatus, (newStatus) => {
  if (props.open && newStatus) {
    form.status = newStatus;
  }
});

const handleSubmit = async () => {
  if (!form.remarks.trim()) {
    errorMessage.value = 'Please provide remarks.';
    return;
  }

  if (!props.hideStatusDropdown) {
    if (!form.status) {
      errorMessage.value = 'Please select a status and provide remarks.';
      return;
    }
    
    if (form.status === props.currentStatus) {
      errorMessage.value = 'Please select a different status.';
      return;
    }
  }

  submitting.value = true;
  errorMessage.value = '';

  try {
    await emit('change', {
      status: form.status,
      remarks: form.remarks.trim()
    });
    // Modal will be closed by parent component after successful change
  } catch (error) {
    console.error('Status change failed:', error);
    errorMessage.value = error.response?.data?.error || 'Failed to change status. Please try again.';
  } finally {
    submitting.value = false;
  }
};

// Status badge styling function
const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'OPEN': return 'bg-blue-100 text-blue-800'
    case 'IN PROGRESS': return 'bg-yellow-100 text-yellow-800'
    case 'RESOLVED': return 'bg-green-100 text-green-800'
    case 'CLOSED': return 'bg-gray-100 text-gray-800'
    default: return 'bg-gray-100 text-gray-800'
  }
};
</script>
