import { defineStore } from 'pinia';
import api from '../services/api';

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [],
    unreadCount: 0,
  }),
  actions: {
    async fetchNotifications() {
      try {
        const response = await api.get('/notifications');
        this.notifications = response.data;
        this.unreadCount = this.notifications.filter(n => !n.is_read).length;
      } catch (error) {
        console.error('Error fetching notifications:', error);
      }
    },
    async markAsRead(notificationId) {
      try {
        await api.post(`/notifications/${notificationId}/read`);
        const notification = this.notifications.find(n => n.id === notificationId);
        if (notification) {
          notification.is_read = true;
          this.unreadCount--;
        }
      } catch (error) {
        console.error('Error marking notification as read:', error);
      }
    },
    async markAllAsRead() {
      try {
        await api.post('/notifications/read-all');
        this.notifications.forEach(n => n.is_read = true);
        this.unreadCount = 0;
      } catch (error) {
        console.error('Error marking all notifications as read:', error);
      }
    },
  },
});
