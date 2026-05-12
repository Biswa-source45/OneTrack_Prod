<template>
  <Modal :open="open" title="Add Account" @close="$emit('close')">
    <FormLayout @submit="save">
      <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
        <p class="text-sm text-red-700">{{ errorMessage }}</p>
      </div>
      
      <FormField label="Account Name" required>
        <input 
          v-model="form.account_name" 
          class="w-full border border-blue-200 rounded px-3 py-2"
          placeholder="Enter account name"
        />
      </FormField>
      
      <FormField label="Account Owner">
        <input 
          v-model="form.account_owner" 
          class="w-full border border-blue-200 rounded px-3 py-2"
          placeholder="Enter account owner name"
        />
      </FormField>
      
      <FormField label="Address">
        <textarea 
          v-model="form.address" 
          class="w-full border border-blue-200 rounded px-3 py-2"
          rows="3"
          placeholder="Enter account address"
        ></textarea>
      </FormField>
      
      <template #actions>
        <Button variant="secondary" type="button" @click="$emit('close')">Cancel</Button>
        <Button type="submit" :disabled="submitting">{{ submitting ? 'Saving...' : 'Save' }}</Button>
      </template>
    </FormLayout>
  </Modal>
</template>

<script setup>
import { ref, watch } from 'vue'
import { createAccount } from '../../api/auth.js'
import Modal from '../ui/Modal.vue'
import FormLayout from '../ui/FormLayout.vue'
import FormField from '../ui/FormField.vue'
import Button from '../ui/Button.vue'

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  initialName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'created'])

const form = ref({
  account_name: '',
  account_owner: '',
  address: ''
})

const submitting = ref(false)
const errorMessage = ref('')

// Watch for initial name changes
watch(() => props.initialName, (newName) => {
  if (newName) {
    form.value.account_name = newName
  }
}, { immediate: true })

// Reset form when modal opens/closes
watch(() => props.open, (isOpen) => {
  if (isOpen) {
    form.value.account_name = props.initialName || ''
    form.value.account_owner = ''
    form.value.address = ''
    errorMessage.value = ''
    submitting.value = false // Reset submitting state
  } else {
    // Reset form when modal closes
    form.value.account_name = ''
    form.value.account_owner = ''
    form.value.address = ''
    errorMessage.value = ''
    submitting.value = false
  }
})

const save = async () => {
  if (!form.value.account_name.trim()) {
    errorMessage.value = 'Account name is required'
    return
  }

  if (submitting.value) {
    return // Prevent double submission
  }

  submitting.value = true
  errorMessage.value = ''

  try {
    const newAccount = await createAccount({
      account_name: form.value.account_name.trim(),
      account_owner: form.value.account_owner.trim(),
      address: form.value.address.trim()
    })
    
    emit('created', newAccount)
    emit('close')
  } catch (error) {
    console.error('Failed to create account:', error)
    errorMessage.value = error.response?.data?.error || 'Failed to create account'
  } finally {
    submitting.value = false
  }
}
</script>
