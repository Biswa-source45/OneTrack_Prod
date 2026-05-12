<template>
  <div class="border-b border-gray-200 last:border-b-0">
    <button
      @click="toggleExpanded"
      class="w-full px-4 py-3 flex items-center justify-between text-left hover:bg-gray-50 transition-colors focus:outline-none focus:bg-gray-50"
      :class="{ 'bg-gray-50': isExpanded }"
    >
      <div class="flex items-center space-x-2">
        <span class="text-sm font-medium text-gray-900">{{ title }}</span>
        <button
          v-if="showEditIcon"
          @click.stop="$emit('edit')"
          class="p-1 rounded hover:bg-gray-200 transition-colors"
          title="Edit"
        >
          <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
      </div>
      <svg
        :class="[
          'w-4 h-4 text-gray-500 transition-transform duration-200',
          isExpanded ? 'transform rotate-180' : ''
        ]"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>
    
    <div
      :class="[
        'overflow-hidden transition-all duration-300 ease-in-out',
        isExpanded ? 'max-h-96 opacity-100' : 'max-h-0 opacity-0'
      ]"
    >
      <div class="px-4 pb-4">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  defaultExpanded: {
    type: Boolean,
    default: true
  },
  showEditIcon: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['edit'])

const isExpanded = ref(props.defaultExpanded)

const toggleExpanded = () => {
  isExpanded.value = !isExpanded.value
}
</script>

<style scoped>
.max-h-96 {
  max-height: 24rem;
}
</style>
