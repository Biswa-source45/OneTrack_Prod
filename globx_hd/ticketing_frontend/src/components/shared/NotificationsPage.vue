<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="py-6">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-900">Notifications</h1>
              <p class="mt-1 text-sm text-gray-500">
                Stay updated with your tickets, tasks, and system updates
              </p>
            </div>
            <div class="flex items-center gap-4">
              <div class="text-sm text-gray-500">
                <span class="font-medium">{{ unreadCount }}</span> unread
              </div>
              <button
                v-if="unreadCount > 0"
                @click="markAllAsRead"
                class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
                :disabled="loading"
              >
                Mark All Read
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="flex flex-wrap items-center gap-4">
          <!-- Category Filter -->
          <div class="flex items-center gap-2">
            <label class="text-sm font-medium text-gray-700">Category:</label>
            <select
              v-model="filters.category"
              @change="updateFilters"
              class="border border-gray-300 rounded-md px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="all">All Categories</option>
              <option value="ticket">Tickets</option>
              <option value="task">Tasks</option>
              <option value="communication">Communication</option>
              <option value="system">System</option>
            </select>
          </div>

          <!-- Priority Filter -->
          <div class="flex items-center gap-2">
            <label class="text-sm font-medium text-gray-700">Priority:</label>
            <select
              v-model="filters.priority"
              @change="updateFilters"
              class="border border-gray-300 rounded-md px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="all">All Priorities</option>
              <option value="urgent">Urgent</option>
              <option value="high">High</option>
              <option value="normal">Normal</option>
              <option value="low">Low</option>
            </select>
          </div>

          <!-- Read Status Filter -->
          <div class="flex items-center gap-2">
            <label class="text-sm font-medium text-gray-700">Status:</label>
            <select
              v-model="filters.isRead"
              @change="updateFilters"
              class="border border-gray-300 rounded-md px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="all">All</option>
              <option value="unread">Unread</option>
              <option value="read">Read</option>
            </select>
          </div>

          <!-- Reset Filters -->
          <button
            @click="resetFilters"
            class="text-sm text-gray-600 hover:text-gray-800 underline"
          >
            Reset Filters
          </button>
        </div>
      </div>
    </div>

    <!-- Notifications List -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-8">
      <!-- Loading State -->
      <div v-if="loading && notifications.length === 0" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p class="text-gray-500 mt-4">Loading notifications...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-12">
        <div class="text-red-600 mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <p class="text-gray-600 mb-4">{{ error }}</p>
        <button
          @click="refreshNotifications"
          class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
        >
          Try Again
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredNotifications.length === 0" class="text-center py-12">
        <div class="text-gray-400 mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No notifications found</h3>
        <p class="text-gray-500">
          {{ hasActiveFilters ? 'Try adjusting your filters to see more notifications.' : 'You\'re all caught up! New notifications will appear here.' }}
        </p>
      </div>

      <!-- Notifications -->
      <div v-else class="space-y-4">
        <div
          v-for="notification in filteredNotifications"
          :key="notification.id"
          class="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow"
          :class="{ 'border-l-4 border-l-blue-500': !notification.is_read }"
        >
          <div class="p-6">
            <div class="flex items-start justify-between">
              <!-- Content -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-3 mb-2">
                  <!-- Category Icon -->
                  <span class="text-xl">{{ getCategoryIcon(notification.category) }}</span>
                  
                  <!-- Title -->
                  <h3 class="text-lg font-medium text-gray-900" :class="{ 'font-semibold': !notification.is_read }">
                    {{ notification.title }}
                  </h3>
                  
                  <!-- Priority Badge -->
                  <span
                    v-if="notification.priority !== 'normal'"
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="getPriorityBadgeClass(notification.priority)"
                  >
                    {{ notification.priority.toUpperCase() }}
                  </span>
                  
                  <!-- Unread Indicator -->
                  <div
                    v-if="!notification.is_read"
                    class="w-3 h-3 bg-blue-500 rounded-full"
                  ></div>
                </div>
                
                <!-- Message -->
                <p class="text-gray-700 mb-4">{{ notification.message }}</p>
                
                <!-- Metadata -->
                <div class="flex items-center gap-4 text-sm text-gray-500">
                  <span>{{ formatTime(notification.created_at) }}</span>
                  <span class="capitalize">{{ notification.category }}</span>
                  <span v-if="notification.read_at && notification.is_read">
                    Read {{ formatTime(notification.read_at) }}
                  </span>
                </div>
              </div>
              
              <!-- Actions -->
              <div class="flex items-center gap-2 ml-4">
                <!-- Mark as Read/Unread -->
                <button
                  @click="toggleReadStatus(notification)"
                  class="p-2 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100 transition-colors"
                  :title="notification.is_read ? 'Mark as unread' : 'Mark as read'"
                >
                  <svg v-if="notification.is_read" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 19v-8.93a2 2 0 01.89-1.664l7-4.666a2 2 0 012.22 0l7 4.666A2 2 0 0121 10.07V19M3 19a2 2 0 002 2h14a2 2 0 002-2M3 19l6.75-4.5M21 19l-6.75-4.5M12 12v7" />
                  </svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </button>
                
                <!-- Navigate to Related Item -->
                <button
                  v-if="notification.related_type && notification.related_id"
                  @click="navigateToRelatedItem(notification)"
                  class="p-2 text-gray-400 hover:text-blue-600 rounded-full hover:bg-blue-50 transition-colors"
                  title="View related item"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </button>
                
                <!-- Delete -->
                <button
                  @click="deleteNotification(notification.id)"
                  class="p-2 text-gray-400 hover:text-red-600 rounded-full hover:bg-red-50 transition-colors"
                  title="Delete notification"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Load More -->
      <div v-if="hasMore && !loading" class="text-center mt-8">
        <button
          @click="loadMore"
          class="bg-white border border-gray-300 text-gray-700 px-6 py-2 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
        >
          Load More
        </button>
      </div>

      <!-- Loading More -->
      <div v-if="loading && notifications.length > 0" class="text-center mt-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useNotificationStore } from '../../stores/notifications';

const router = useRouter();
const notificationStore = useNotificationStore();

// Local state
const filters = ref({
  category: 'all',
  priority: 'all',
  isRead: 'all'
});

// Computed properties
const notifications = computed(() => notificationStore.notifications);
const filteredNotifications = computed(() => notificationStore.filteredNotifications);
const unreadCount = computed(() => notificationStore.unreadCount);
const loading = computed(() => notificationStore.loading);
const error = computed(() => notificationStore.error);
const hasMore = computed(() => notificationStore.pagination.hasMore);

const hasActiveFilters = computed(() => {
  return filters.value.category !== 'all' || 
         filters.value.priority !== 'all' || 
         filters.value.isRead !== 'all';
});

// Methods
const updateFilters = async () => {
  await notificationStore.updateFilters(filters.value);
};

const resetFilters = async () => {
  filters.value = {
    category: 'all',
    priority: 'all',
    isRead: 'all'
  };
  await notificationStore.resetFilters();
};

const refreshNotifications = async () => {
  await notificationStore.fetchNotifications(true);
};

const markAllAsRead = async () => {
  await notificationStore.markAllAsRead();
};

const toggleReadStatus = async (notification) => {
  if (notification.is_read) {
    // Note: Backend doesn't support marking as unread, so we'll just mark as read
    // This could be enhanced in the future
    return;
  } else {
    await notificationStore.markAsRead(notification.id);
  }
};

const deleteNotification = async (notificationId) => {
  if (confirm('Are you sure you want to delete this notification?')) {
    await notificationStore.deleteNotification(notificationId);
  }
};

const loadMore = async () => {
  await notificationStore.loadMore();
};

const navigateToRelatedItem = (notification) => {
  try {
    const metadata = typeof notification.metadata === 'string' 
      ? JSON.parse(notification.metadata) 
      : notification.metadata || {};
    
    // Navigate based on related type
    if (notification.related_type === 'ticket' && notification.related_id) {
      router.push(`/tickets/${notification.related_id}`);
    } else if (notification.related_type === 'task' && notification.related_id) {
      router.push(`/tasks/${notification.related_id}`);
    }
  } catch (error) {
    console.error('Error navigating to related item:', error);
  }
};

const getCategoryIcon = (category) => {
  return notificationStore.getCategoryIcon(category);
};

const formatTime = (timestamp) => {
  return notificationStore.formatNotificationTime(timestamp);
};

const getPriorityBadgeClass = (priority) => {
  const classes = {
    urgent: 'bg-red-100 text-red-800',
    high: 'bg-orange-100 text-orange-800',
    normal: 'bg-blue-100 text-blue-800',
    low: 'bg-gray-100 text-gray-800'
  };
  return classes[priority] || classes.normal;
};

// Lifecycle
onMounted(async () => {
  // Initialize notifications
  await notificationStore.initialize();
  
  // Sync local filters with store
  filters.value = { ...notificationStore.filters };
});

onUnmounted(() => {
  // Cleanup polling when leaving the page
  notificationStore.cleanup();
});
</script>
