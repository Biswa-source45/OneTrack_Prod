<template>
  <div class="relative">
    <!-- Display Mode -->
    <div 
      v-if="!isEditing"
      @click="startEdit"
      class="cursor-pointer hover:bg-gray-50 rounded px-1 py-0.5 transition-colors"
      :class="displayClass"
    >
      {{ displayValue }}
    </div>
    
    <!-- Edit Mode -->
    <div v-else class="relative">
      <!-- Date Picker Mode -->
      <input
        v-if="isDatePicker"
        ref="inputRef"
        v-model="selectedValue"
        type="date"
        @change="saveChange"
        @blur="cancelEdit"
        class="w-full px-2 py-1 border border-blue-400 rounded focus:outline-none focus:ring-2 focus:ring-blue-400 text-sm"
        :class="selectClass"
      />
      <!-- Regular Select Mode -->
      <select
        v-else
        ref="selectRef"
        v-model="selectedValue"
        @change="saveChange"
        @blur="cancelEdit"
        class="w-full px-2 py-1 border border-blue-400 rounded focus:outline-none focus:ring-2 focus:ring-blue-400 text-sm"
        :class="selectClass"
      >
        <option v-if="allowEmpty" value="">{{ emptyLabel }}</option>
        <option 
          v-for="option in options" 
          :key="getOptionValue(option)" 
          :value="getOptionValue(option)"
        >
          {{ getOptionLabel(option) }}
        </option>
      </select>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  options: {
    type: Array,
    required: true
  },
  optionValue: {
    type: String,
    default: 'value'
  },
  optionLabel: {
    type: String,
    default: 'label'
  },
  displayClass: {
    type: String,
    default: 'text-sm font-medium text-gray-900'
  },
  selectClass: {
    type: String,
    default: ''
  },
  emptyLabel: {
    type: String,
    default: 'Select...'
  },
  allowEmpty: {
    type: Boolean,
    default: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  isDatePicker: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const isEditing = ref(false)
const selectedValue = ref(props.modelValue)
const selectRef = ref(null)
const inputRef = ref(null)

// Computed display value
const displayValue = computed(() => {
  if (props.loading) return 'Loading...'
  
  if (!props.modelValue) {
    return props.emptyLabel
  }
  
  // For date picker, display the formatted date directly
  if (props.isDatePicker) {
    try {
      const date = new Date(props.modelValue)
      if (isNaN(date.getTime())) return props.emptyLabel
      return date.toLocaleDateString()
    } catch (error) {
      return props.emptyLabel
    }
  }
  
  // For regular dropdowns, find in options
  const option = props.options.find(opt => getOptionValue(opt) == props.modelValue)
  return option ? getOptionLabel(option) : props.emptyLabel
})

// Helper functions for option handling
const getOptionValue = (option) => {
  return typeof option === 'object' ? option[props.optionValue] : option
}

const getOptionLabel = (option) => {
  return typeof option === 'object' ? option[props.optionLabel] : option
}

// Watch for external changes to modelValue
watch(() => props.modelValue, (newValue) => {
  selectedValue.value = newValue
})

// Start editing
const startEdit = async () => {
  if (props.loading) return
  
  isEditing.value = true
  selectedValue.value = props.modelValue
  
  await nextTick()
  if (props.isDatePicker && inputRef.value) {
    inputRef.value.focus()
  } else if (selectRef.value) {
    selectRef.value.focus()
  }
}

// Save changes
const saveChange = () => {
  if (selectedValue.value !== props.modelValue) {
    emit('update:modelValue', selectedValue.value)
    emit('change', selectedValue.value)
  }
  isEditing.value = false
}

// Cancel editing
const cancelEdit = () => {
  selectedValue.value = props.modelValue
  isEditing.value = false
}
</script>
