<template>
  <Modal :open="open" title="Add Product" @close="$emit('close')">
    <FormLayout @submit="save">
      <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
        <p class="text-sm text-red-700">{{ errorMessage }}</p>
      </div>
      <FormField label="Product Name" :error="errors.product_name">
        <input 
          v-model="form.product_name" 
          type="text" 
          class="w-full border border-blue-200 rounded px-3 py-2" 
          required 
        />
      </FormField>
      <FormField label="Description">
        <input 
          v-model="form.product_description" 
          type="text" 
          class="w-full border border-blue-200 rounded px-3 py-2" 
        />
      </FormField>
      <template #actions>
        <Button variant="secondary" type="button" @click="$emit('close')">Cancel</Button>
        <Button type="submit" :disabled="submitting">{{ submitting ? 'Saving...' : 'Save' }}</Button>
      </template>
    </FormLayout>
  </Modal>
</template>

<script setup>
import { ref, reactive, watch } from 'vue';
import Modal from '../ui/Modal.vue';
import FormLayout from '../ui/FormLayout.vue';
import FormField from '../ui/FormField.vue';
import Button from '../ui/Button.vue';
import { createProduct } from '../../api/auth';

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  initialName: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['close', 'created']);

const submitting = ref(false);
const errorMessage = ref('');

const form = reactive({
  product_name: '',
  product_description: ''
});

const errors = reactive({
  product_name: ''
});

// Set initial product name when modal opens
watch(() => props.open, (isOpen) => {
  if (isOpen && props.initialName) {
    form.product_name = props.initialName;
  }
});

// Reset form when modal closes
watch(() => props.open, (isOpen) => {
  if (!isOpen) {
    Object.assign(form, {
      product_name: '',
      product_description: ''
    });
    Object.assign(errors, {
      product_name: ''
    });
    errorMessage.value = '';
  }
});

const save = async () => {
  if (submitting.value) return;
  
  // Validate form
  errors.product_name = form.product_name ? '' : 'Product name is required';
  if (errors.product_name) return;
  
  submitting.value = true;
  errorMessage.value = '';
  
  try {
    const payload = {
      product_name: form.product_name,
      product_description: form.product_description
    };
    
    const newProduct = await createProduct(payload);
    emit('created', newProduct);
    emit('close');
  } catch (error) {
    console.error('Failed to create product:', error);
    errorMessage.value = error.response?.data?.error || 'Failed to create product. Please try again.';
  } finally {
    submitting.value = false;
  }
};
</script>
