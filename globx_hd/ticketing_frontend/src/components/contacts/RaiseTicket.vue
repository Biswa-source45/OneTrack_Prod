<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-xl mx-auto mt-8">
    <h2 class="text-2xl font-bold text-blue-800 mb-4">Raise a Ticket</h2>
    <form @submit.prevent="onSubmit" class="space-y-4">
      <FormField label="Product" required>
        <select v-model="form.productId" class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400">
          <option value="" disabled>Select Product</option>
          <option v-for="product in products" :key="product.id" :value="product.id">{{ product.product_name }}</option>
        </select>
      </FormField>
      <FormField label="Subject" required>
        <input v-model="form.subject" type="text" class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400" placeholder="Enter ticket subject..." />
      </FormField>
      <FormField label="Ticket Details" required>
        <textarea v-model="form.ticketDetails" rows="4" class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400" placeholder="Describe your issue..."></textarea>
      </FormField>
      
      <!-- File Upload Section -->
      <FormField label="Attach Files">
        <FileUpload 
          ref="fileUpload"
          :ticket-id="createdTicketId"
          :show-upload-button="false"
          @files-selected="onFilesSelected"
          @upload-success="onUploadSuccess"
          @upload-error="onUploadError"
        />
      </FormField>
      
      <Button 
        type="submit" 
        :disabled="isSubmitting"
        class="w-full bg-blue-700 text-white py-2 rounded hover:bg-blue-800 transition disabled:opacity-50"
      >
        {{ isSubmitting ? 'Creating Ticket...' : 'Submit Ticket' }}
      </Button>
    </form>
  <div v-if="showSuccess" class="text-green-600 mt-4 text-center font-semibold">Ticket raised successfully!</div>
    <Modal v-if="error" @close="error = ''">
      <div class="text-red-700 font-semibold">{{ error }}</div>
    </Modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../../stores/auth';
import { fetchProducts, createTicket } from '../../api/tickets';
import FormField from '../ui/FormField.vue';
import Button from '../ui/Button.vue';
import Modal from '../ui/Modal.vue';
import FileUpload from '../shared/FileUpload.vue';

const auth = useAuthStore();
const form = ref({
  productId: '',
  subject: '',
  ticketDetails: '',
});
const products = ref([]);
const showSuccess = ref(false);
const error = ref('');
const isSubmitting = ref(false);
const createdTicketId = ref(null);
const selectedFiles = ref([]);
const fileUpload = ref(null);

onMounted(async () => {
  try {
    products.value = await fetchProducts();
  } catch (err) {
    error.value = err.message || 'Failed to load master data.';
  }
});

async function onSubmit() {
  error.value = '';
  if (!form.value.productId || !form.value.subject || !form.value.ticketDetails) {
    error.value = 'Please fill all required fields.';
    return;
  }
  
  isSubmitting.value = true;
  
  try {
    const payload = {
      product_id: Number(form.value.productId),
      subject: form.value.subject,
      ticket_details: form.value.ticketDetails,
    };
    
    // Create ticket first
    const response = await createTicket(payload);
    const ticketId = response.ticket.ticket_id;
    createdTicketId.value = ticketId;
    
    // Upload files if any are selected
    console.log('🔥 DEBUG: selectedFiles.value.length:', selectedFiles.value.length);
    console.log('🔥 DEBUG: createdTicketId.value:', createdTicketId.value);
    if (selectedFiles.value.length > 0) {
      console.log('🔥 DEBUG: Uploading files...');
      // Pass the ticket ID directly to uploadFiles method
      await fileUpload.value.uploadFiles(ticketId);
    } else {
      console.log('🔥 DEBUG: No files to upload');
    }
    
    showSuccess.value = true;
    form.value = { productId: '', subject: '', ticketDetails: '' };
    createdTicketId.value = null;
    selectedFiles.value = [];
    
    if (fileUpload.value) {
      fileUpload.value.clearFiles();
    }
    
    setTimeout(() => { showSuccess.value = false; }, 3000);
    
  } catch (err) {
    error.value = err.message || 'Failed to raise ticket.';
  } finally {
    isSubmitting.value = false;
  }
}

function onFilesSelected(files) {
  console.log('🔥 DEBUG: onFilesSelected called with', files.length, 'files');
  selectedFiles.value = files;
}

function onUploadSuccess(response) {
  console.log('Files uploaded successfully:', response);
  if (response.errors && response.errors.length > 0) {
    error.value = 'Some files failed to upload: ' + response.errors.join(', ');
  }
}

function onUploadError(errorMessage) {
  error.value = 'File upload failed: ' + errorMessage;
}
</script>
