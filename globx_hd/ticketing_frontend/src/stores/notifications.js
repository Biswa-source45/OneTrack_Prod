import { defineStore } from 'pinia';
import { notificationApi } from '../api/notifications';
import websocketService from '../services/websocket';

export const useNotificationStore = defineStore('notifications', {
  state: () => ({
    notifications: [],
    unreadCount: 0,
    loading: false,
    error: null,
    filters: {
      category: 'all', // 'all', 'ticket', 'task', 'communication', 'system'
      priority: 'all',  // 'all', 'low', 'normal', 'high', 'urgent'
      isRead: 'all'     // 'all', 'read', 'unread'
    },
    pagination: {
      currentPage: 1,
      limit: 20,
      hasMore: true
    },
    wsConnected: false,
    wsUnsubscribers: []
  }),

  getters: {
    // Get notifications filtered by current filters
    filteredNotifications: (state) => {
      return state.notifications.filter(notification => {
        if (state.filters.category !== 'all' && notification.category !== state.filters.category) {
          return false;
        }
        if (state.filters.priority !== 'all' && notification.priority !== state.filters.priority) {
          return false;
        }
        if (state.filters.isRead !== 'all') {
          const isRead = notification.is_read;
          if (state.filters.isRead === 'read' && !isRead) return false;
          if (state.filters.isRead === 'unread' && isRead) return false;
        }
        return true;
      });
    },

    // Get unread notifications
    unreadNotifications: (state) => {
      return state.notifications.filter(n => !n.is_read);
    },

    // Get notifications by category
    getNotificationsByCategory: (state) => {
      return (category) => state.notifications.filter(n => n.category === category);
    },

    // Get notifications by priority
    getNotificationsByPriority: (state) => {
      return (priority) => state.notifications.filter(n => n.priority === priority);
    },

    // Check if there are any high priority unread notifications
    hasUrgentNotifications: (state) => {
      return state.notifications.some(n => !n.is_read && (n.priority === 'high' || n.priority === 'urgent'));
    }
  },

  actions: {
    // Fetch notifications with current filters and pagination
    async fetchNotifications(reset = false) {
      if (reset) {
        this.notifications = [];
        this.pagination.currentPage = 1;
        this.pagination.hasMore = true;
      }

      if (!this.pagination.hasMore && !reset) {
        return;
      }

      this.loading = true;
      this.error = null;

      try {
        const offset = (this.pagination.currentPage - 1) * this.pagination.limit;
        const params = {
          limit: this.pagination.limit,
          offset: offset,
          category: this.filters.category,
          priority: this.filters.priority,
          is_read: this.filters.isRead
        };

        const response = await notificationApi.getNotifications(params);
        const newNotifications = response.data.notifications || [];

        if (reset) {
          this.notifications = newNotifications;
        } else {
          // Avoid duplicates when appending
          const existingIds = new Set(this.notifications.map(n => n.id));
          const uniqueNew = newNotifications.filter(n => !existingIds.has(n.id));
          this.notifications.push(...uniqueNew);
        }

        // Update pagination
        this.pagination.hasMore = newNotifications.length === this.pagination.limit;
        if (this.pagination.hasMore) {
          this.pagination.currentPage++;
        }

      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch notifications';
        console.error('Error fetching notifications:', error);
      } finally {
        this.loading = false;
      }
    },

    // Fetch unread count
    async fetchUnreadCount() {
      try {
        const response = await notificationApi.getUnreadCount();
        this.unreadCount = response.data.unread_count || 0;
      } catch (error) {
        console.error('Error fetching unread count:', error);
      }
    },

    // Mark a specific notification as read
    async markAsRead(notificationId) {
      try {
        await notificationApi.markAsRead(notificationId);
        
        // Update local state
        const notification = this.notifications.find(n => n.id === notificationId);
        if (notification && !notification.is_read) {
          notification.is_read = true;
          notification.read_at = new Date().toISOString();
          this.unreadCount = Math.max(0, this.unreadCount - 1);
        }
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to mark notification as read';
        console.error('Error marking notification as read:', error);
      }
    },

    // Mark all notifications as read
    async markAllAsRead() {
      try {
        await notificationApi.markAllAsRead();
        
        // Update local state
        let markedCount = 0;
        this.notifications.forEach(notification => {
          if (!notification.is_read) {
            notification.is_read = true;
            notification.read_at = new Date().toISOString();
            markedCount++;
          }
        });
        
        this.unreadCount = 0;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to mark all notifications as read';
        console.error('Error marking all notifications as read:', error);
      }
    },

    // Delete a notification
    async deleteNotification(notificationId) {
      try {
        await notificationApi.deleteNotification(notificationId);
        
        // Remove from local state
        const index = this.notifications.findIndex(n => n.id === notificationId);
        if (index !== -1) {
          const notification = this.notifications[index];
          if (!notification.is_read) {
            this.unreadCount = Math.max(0, this.unreadCount - 1);
          }
          this.notifications.splice(index, 1);
        }
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to delete notification';
        console.error('Error deleting notification:', error);
      }
    },

    // Update filters and refresh notifications
    async updateFilters(newFilters) {
      this.filters = { ...this.filters, ...newFilters };
      await this.fetchNotifications(true);
    },

    // Reset filters to default
    async resetFilters() {
      this.filters = {
        category: 'all',
        priority: 'all',
        isRead: 'all'
      };
      await this.fetchNotifications(true);
    },

    // Load more notifications (for infinite scroll)
    async loadMore() {
      if (!this.loading && this.pagination.hasMore) {
        await this.fetchNotifications(false);
      }
    },

    // Initialize WebSocket connection for real-time notifications
    initializeWebSocket() {
      // Connect to WebSocket
      websocketService.connect();

      // Listen for connection state changes
      const connectionUnsubscriber = websocketService.onConnectionStateChange((state) => {
        this.wsConnected = (state === 'connected');
        console.log('🔌 WebSocket connection state:', state);
      });
      this.wsUnsubscribers.push(connectionUnsubscriber);

      // Handle new notifications
      const newNotificationUnsubscriber = websocketService.on('notification.new', (data) => {
        console.log('🔔 New notification received:', data);
        
        // Add notification to the beginning of the list
        this.notifications.unshift({
          id: data.id,
          title: data.title,
          message: data.message,
          notification_type: data.notification_type,
          priority: data.priority,
          category: data.category,
          related_id: data.related_id,
          related_type: data.related_type,
          is_read: data.is_read || false,
          created_at: data.created_at
        });

        // Show toast notification
        this.showNewNotificationToast(1);
      });
      this.wsUnsubscribers.push(newNotificationUnsubscriber);

      // Handle unread count updates
      const countUpdateUnsubscriber = websocketService.on('count.update', (data) => {
        console.log('📊 Unread count updated:', data.unread_count);
        this.unreadCount = data.unread_count;
      });
      this.wsUnsubscribers.push(countUpdateUnsubscriber);

      // Handle notification marked as read
      const readUnsubscriber = websocketService.on('notification.read', (data) => {
        console.log('✅ Notification marked as read:', data.notification_id);
        
        const notification = this.notifications.find(n => n.id === data.notification_id);
        if (notification && !notification.is_read) {
          notification.is_read = true;
          notification.read_at = new Date().toISOString();
        }
      });
      this.wsUnsubscribers.push(readUnsubscriber);

      // Handle all notifications marked as read
      const allReadUnsubscriber = websocketService.on('notification.all_read', () => {
        console.log('✅ All notifications marked as read');
        
        this.notifications.forEach(notification => {
          if (!notification.is_read) {
            notification.is_read = true;
            notification.read_at = new Date().toISOString();
          }
        });
      });
      this.wsUnsubscribers.push(allReadUnsubscriber);

      // Handle notification deletion
      const deleteUnsubscriber = websocketService.on('notification.deleted', (data) => {
        console.log('🗑️ Notification deleted:', data.notification_id);
        
        const index = this.notifications.findIndex(n => n.id === data.notification_id);
        if (index !== -1) {
          this.notifications.splice(index, 1);
        }
      });
      this.wsUnsubscribers.push(deleteUnsubscriber);
    },

    // Cleanup WebSocket connection
    cleanupWebSocket() {
      // Unsubscribe from all WebSocket events
      this.wsUnsubscribers.forEach(unsubscribe => unsubscribe());
      this.wsUnsubscribers = [];
      
      // Disconnect WebSocket
      websocketService.disconnect();
      this.wsConnected = false;
    },

    // Show toast notification for new notifications
    showNewNotificationToast(count) {
      // This can be implemented with your preferred toast library
      // For now, just emit a custom event
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new CustomEvent('new-notifications', {
          detail: { count }
        }));
      }
    },

    // Initialize notification store
    async initialize() {
      // Fetch initial data
      await Promise.all([
        this.fetchNotifications(true),
        this.fetchUnreadCount()
      ]);
      
      // Initialize WebSocket for real-time updates
      this.initializeWebSocket();
    },

    // Cleanup when store is no longer needed
    cleanup() {
      this.cleanupWebSocket();
    },

    // Clear error state
    clearError() {
      this.error = null;
    },

    // Get WebSocket connection status
    isWebSocketConnected() {
      return this.wsConnected;
    },

    // Get notification priority color class
    getPriorityColorClass(priority) {
      const colorMap = {
        low: 'text-gray-500',
        normal: 'text-blue-500',
        high: 'text-orange-500',
        urgent: 'text-red-500'
      };
      return colorMap[priority] || colorMap.normal;
    },

    // Get notification category icon
    getCategoryIcon(category) {
      const iconMap = {
        ticket: '🎫',
        task: '📋',
        communication: '💬',
        system: '⚙️'
      };
      return iconMap[category] || '📢';
    },

    // Format notification time
    formatNotificationTime(timestamp) {
      const date = new Date(timestamp);
      const now = new Date();
      const diffInMinutes = Math.floor((now - date) / (1000 * 60));
      
      if (diffInMinutes < 1) return 'Just now';
      if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
      
      const diffInHours = Math.floor(diffInMinutes / 60);
      if (diffInHours < 24) return `${diffInHours}h ago`;
      
      const diffInDays = Math.floor(diffInHours / 24);
      if (diffInDays < 7) return `${diffInDays}d ago`;
      
      return date.toLocaleDateString();
    }
  }
});
