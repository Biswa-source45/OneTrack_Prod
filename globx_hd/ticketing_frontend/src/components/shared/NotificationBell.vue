<template>
  <div class="relative">
    <!-- Notification Bell Button -->
    <button
      @click="toggleDropdown"
      class="relative flex items-center justify-center w-10 h-10 rounded-full bg-blue-700 hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-300 transition-colors"
      :class="{ 'ring-2 ring-blue-300': isOpen }"
    >
      <!-- Bell Icon -->
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
      </svg>
      
      <!-- Notification Badge -->
      <span
        v-if="unreadCount > 0"
        class="absolute -top-1 -right-1 inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-white transform translate-x-1/2 -translate-y-1/2 rounded-full"
        :class="hasUrgent ? 'bg-red-500 animate-pulse' : 'bg-red-500'"
      >
        {{ displayCount }}
      </span>
    </button>

    <!-- Dropdown Menu -->
    <div
      v-if="isOpen"
      class="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-lg border border-gray-200 z-50"
      @click.stop
    >
      <!-- Header -->
      <div class="px-4 py-3 border-b border-gray-200 flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">Notifications</h3>
        <div class="flex items-center gap-2">
          <button
            v-if="unreadCount > 0"
            @click="markAllAsRead"
            class="text-sm text-blue-600 hover:text-blue-800 font-medium"
            :disabled="loading"
          >
            Mark all read
          </button>
          <button
            @click="openNotificationsPage"
            class="text-sm text-gray-600 hover:text-gray-800"
          >
            View all
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="px-4 py-6 text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        <p class="text-sm text-gray-500 mt-2">Loading notifications...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="px-4 py-6 text-center">
        <p class="text-sm text-red-600">{{ error }}</p>
        <button
          @click="refreshNotifications"
          class="mt-2 text-sm text-blue-600 hover:text-blue-800"
        >
          Try again
        </button>
      </div>

      <!-- Notifications List -->
      <div v-else-if="recentNotifications.length > 0" class="max-h-96 overflow-y-auto">
        <div
          v-for="notification in recentNotifications"
          :key="notification.id"
          class="px-4 py-3 border-b border-gray-100 hover:bg-gray-50 cursor-pointer transition-colors"
          :class="{ 'bg-blue-50': !notification.is_read }"
          @click="handleNotificationClick(notification)"
        >
          <div class="flex items-start gap-3">
            <!-- Category Icon -->
            <div class="flex-shrink-0 mt-1">
              <span class="text-lg">{{ getCategoryIcon(notification.category) }}</span>
            </div>
            
            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between">
                <h4 class="text-sm font-medium text-gray-900 truncate" :class="{ 'font-semibold': !notification.is_read }">
                  {{ notification.title }}
                </h4>
                <div class="flex items-center gap-1 ml-2">
                  <!-- Priority Indicator -->
                  <div
                    v-if="notification.priority === 'high' || notification.priority === 'urgent'"
                    class="w-2 h-2 rounded-full"
                    :class="notification.priority === 'urgent' ? 'bg-red-500' : 'bg-orange-500'"
                  ></div>
                  <!-- Unread Indicator -->
                  <div
                    v-if="!notification.is_read"
                    class="w-2 h-2 bg-blue-500 rounded-full"
                  ></div>
                </div>
              </div>
              
              <p class="text-sm text-gray-600 mt-1 line-clamp-2">
                {{ notification.message }}
              </p>
              
              <p class="text-xs text-gray-400 mt-2">
                {{ formatTime(notification.created_at) }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="px-4 py-8 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-gray-300 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
        <p class="text-sm text-gray-500">No notifications yet</p>
      </div>

      <!-- Footer -->
      <div class="px-4 py-3 border-t border-gray-200 text-center">
        <button
          @click="openNotificationsPage"
          class="text-sm text-blue-600 hover:text-blue-800 font-medium"
        >
          View all notifications
        </button>
      </div>
    </div>

    <!-- Backdrop -->
    <div
      v-if="isOpen"
      class="fixed inset-0 z-40"
      @click="closeDropdown"
    ></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useNotificationStore } from '../../stores/notifications';

const router = useRouter();
const notificationStore = useNotificationStore();

// Component state
const isOpen = ref(false);

// Computed properties
const unreadCount = computed(() => notificationStore.unreadCount);
const loading = computed(() => notificationStore.loading);
const error = computed(() => notificationStore.error);
const hasUrgent = computed(() => notificationStore.hasUrgentNotifications);

// Display count (max 99+)
const displayCount = computed(() => {
  return unreadCount.value > 99 ? '99+' : unreadCount.value.toString();
});

// Get recent notifications (max 5 for dropdown)
const recentNotifications = computed(() => {
  return notificationStore.notifications.slice(0, 5);
});

// Methods
const toggleDropdown = () => {
  isOpen.value = !isOpen.value;
  if (isOpen.value && recentNotifications.value.length === 0) {
    refreshNotifications();
  }
};

const closeDropdown = () => {
  isOpen.value = false;
};

const refreshNotifications = async () => {
  await notificationStore.fetchNotifications(true);
};

const markAllAsRead = async () => {
  await notificationStore.markAllAsRead();
};

const handleNotificationClick = async (notification) => {
  // Mark as read if not already read
  if (!notification.is_read) {
    await notificationStore.markAsRead(notification.id);
  }
  
  // Navigate based on notification type and related data
  navigateToRelatedItem(notification);
  
  closeDropdown();
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
    } else {
      // Default to notifications page
      openNotificationsPage();
    }
  } catch (error) {
    console.error('Error navigating to related item:', error);
    openNotificationsPage();
  }
};

const openNotificationsPage = () => {
  router.push('/notifications');
  closeDropdown();
};

const getCategoryIcon = (category) => {
  return notificationStore.getCategoryIcon(category);
};

const formatTime = (timestamp) => {
  return notificationStore.formatNotificationTime(timestamp);
};

// Handle click outside to close dropdown
const handleClickOutside = (event) => {
  if (isOpen.value && !event.target.closest('.relative')) {
    closeDropdown();
  }
};

// Lifecycle
onMounted(() => {
  // Initialize notifications (will connect WebSocket)
  notificationStore.fetchUnreadCount();
  
  // Add click outside listener
  document.addEventListener('click', handleClickOutside);
  
  // Listen for new notifications (from WebSocket via store)
  window.addEventListener('new-notifications', handleNewNotifications);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
  window.removeEventListener('new-notifications', handleNewNotifications);
});

const handleNewNotifications = (event) => {
  // Show visual feedback for new notifications
  if (event.detail.count > 0) {
    // Could add animation or sound here
    console.log(`${event.detail.count} new notification(s) received`);
    
    // Optionally play notification sound
    // const audio = new Audio('/notification.mp3');
    // audio.play().catch(e => console.log('Audio play failed:', e));
  }
};
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
