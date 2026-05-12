import api from './api';

export const notificationApi = {
  // Get notifications with filters and pagination
  getNotifications(params = {}) {
    const queryParams = new URLSearchParams();
    
    if (params.limit) queryParams.append('limit', params.limit);
    if (params.offset) queryParams.append('offset', params.offset);
    if (params.category && params.category !== 'all') queryParams.append('category', params.category);
    if (params.priority && params.priority !== 'all') queryParams.append('priority', params.priority);
    if (params.is_read && params.is_read !== 'all') queryParams.append('is_read', params.is_read);
    
    const queryString = queryParams.toString();
    const url = queryString ? `/notifications?${queryString}` : '/notifications';
    
    return api.get(url);
  },

  // Get unread notification count
  getUnreadCount() {
    return api.get('/notifications/unread-count');
  },

  // Mark a specific notification as read
  markAsRead(notificationId) {
    return api.patch(`/notifications/${notificationId}/read`);
  },

  // Mark all notifications as read
  markAllAsRead() {
    return api.patch('/notifications/mark-all-read');
  },

  // Delete a notification
  deleteNotification(notificationId) {
    return api.delete(`/notifications/${notificationId}`);
  }
};

export default notificationApi;
