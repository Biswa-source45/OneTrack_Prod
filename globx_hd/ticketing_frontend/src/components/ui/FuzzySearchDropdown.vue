<template>
  <div class="relative">
    <!-- Search Input -->
    <input
      ref="searchInput"
      v-model="searchQuery"
      type="text"
      :placeholder="placeholder"
      class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
      @focus="onFocus"
      @input="onInput"
      @keydown="onKeydown"
    />
    
    <!-- Dropdown -->
    <div 
      v-if="showDropdown && (filteredItems.length > 0 || searchQuery.length > 0)" 
      class="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-64 overflow-y-auto"
    >
      <div 
        v-for="(item, index) in filteredItems" 
        :key="getItemKey(item)"
        @click="selectItem(item)"
        :class="[
          'p-3 cursor-pointer border-b border-gray-100 last:border-b-0',
          index === highlightedIndex ? 'bg-blue-100' : 'hover:bg-blue-50'
        ]"
      >
        <slot name="item" :item="item" :searchQuery="searchQuery">
          {{ getItemDisplay(item) }}
        </slot>
      </div>
      
      <div v-if="filteredItems.length === 0 && searchQuery.length > 0" class="p-3 text-gray-500 text-center">
        No results found for "{{ searchQuery }}"
        <span v-if="allowCreate" 
              @click="$emit('create-new', searchQuery)"
              class="ml-2 text-blue-600 hover:text-blue-800 cursor-pointer underline">
          Create New?
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';

const props = defineProps({
  items: {
    type: Array,
    required: true
  },
  modelValue: {
    type: [Object, String, Number],
    default: null
  },
  placeholder: {
    type: String,
    default: 'Search...'
  },
  displayKey: {
    type: String,
    default: 'name'
  },
  searchKeys: {
    type: Array,
    default: () => ['name']
  },
  keyProperty: {
    type: String,
    default: 'id'
  },
  allowCreate: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue', 'select', 'create-new']);

const searchQuery = ref('');
const showDropdown = ref(false);
const highlightedIndex = ref(-1);
const searchInput = ref(null);

// Initialize search query with selected item display text
watch(() => props.modelValue, (newValue, oldValue) => {
  console.log('FuzzySearchDropdown modelValue changed:', { newValue, oldValue, displayKey: props.displayKey });
  if (newValue) {
    const displayText = getItemDisplay(newValue);
    console.log('Setting searchQuery to:', displayText);
    searchQuery.value = displayText;
  } else {
    searchQuery.value = '';
  }
}, { immediate: true, deep: true });

// Fuzzy search function
const fuzzyMatch = (text, query) => {
  if (!query) return true;
  
  const textLower = text.toLowerCase();
  const queryLower = query.toLowerCase();
  
  // Exact match gets highest priority
  if (textLower.includes(queryLower)) {
    return { score: 100, match: true };
  }
  
  // Fuzzy matching - check if all characters of query exist in text in order
  let textIndex = 0;
  let queryIndex = 0;
  let matches = 0;
  
  while (textIndex < text.length && queryIndex < query.length) {
    if (textLower[textIndex] === queryLower[queryIndex]) {
      matches++;
      queryIndex++;
    }
    textIndex++;
  }
  
  if (queryIndex === query.length) {
    // All query characters found in order
    const score = (matches / text.length) * 50; // Lower score than exact match
    return { score, match: true };
  }
  
  return { score: 0, match: false };
};

// Filtered and sorted items based on search query
const filteredItems = computed(() => {
  if (!searchQuery.value.trim()) {
    return props.items;
  }
  
  const query = searchQuery.value.trim();
  const results = [];
  
  props.items.forEach(item => {
    let bestScore = 0;
    let hasMatch = false;
    
    // Search across all specified keys
    props.searchKeys.forEach(key => {
      const value = getNestedValue(item, key);
      if (value) {
        const result = fuzzyMatch(String(value), query);
        if (result.match) {
          hasMatch = true;
          bestScore = Math.max(bestScore, result.score);
        }
      }
    });
    
    if (hasMatch) {
      results.push({ item, score: bestScore });
    }
  });
  
  // Sort by score (highest first)
  return results
    .sort((a, b) => b.score - a.score)
    .map(result => result.item);
});

// Helper function to get nested object values
const getNestedValue = (obj, path) => {
  return path.split('.').reduce((current, key) => current?.[key], obj);
};

// Get display text for an item
const getItemDisplay = (item) => {
  if (!item) return '';
  return getNestedValue(item, props.displayKey) || '';
};

// Get unique key for an item
const getItemKey = (item) => {
  return getNestedValue(item, props.keyProperty) || item;
};

// Handle focus - clear search to show all items
const onFocus = () => {
  searchQuery.value = '';
  showDropdown.value = true;
  highlightedIndex.value = -1;
};

// Handle input changes
const onInput = () => {
  highlightedIndex.value = -1;
  showDropdown.value = true;
};

// Handle keyboard navigation
const onKeydown = (event) => {
  if (!showDropdown.value) return;
  
  switch (event.key) {
    case 'ArrowDown':
      event.preventDefault();
      highlightedIndex.value = Math.min(highlightedIndex.value + 1, filteredItems.value.length - 1);
      break;
    case 'ArrowUp':
      event.preventDefault();
      highlightedIndex.value = Math.max(highlightedIndex.value - 1, -1);
      break;
    case 'Enter':
      event.preventDefault();
      if (highlightedIndex.value >= 0 && filteredItems.value[highlightedIndex.value]) {
        selectItem(filteredItems.value[highlightedIndex.value]);
      }
      break;
    case 'Escape':
      showDropdown.value = false;
      highlightedIndex.value = -1;
      restoreSelectedValue();
      searchInput.value?.blur();
      break;
  }
};

// Select an item
const selectItem = (item) => {
  emit('update:modelValue', item);
  emit('select', item);
  searchQuery.value = getItemDisplay(item);
  showDropdown.value = false;
  highlightedIndex.value = -1;
  searchInput.value?.blur();
};

// Restore selected value when closing dropdown without selection
const restoreSelectedValue = () => {
  if (props.modelValue) {
    searchQuery.value = getItemDisplay(props.modelValue);
  } else {
    searchQuery.value = '';
  }
};

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    showDropdown.value = false;
    highlightedIndex.value = -1;
    restoreSelectedValue();
  }
};

// Add global click listener
import { onMounted, onUnmounted } from 'vue';

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>
