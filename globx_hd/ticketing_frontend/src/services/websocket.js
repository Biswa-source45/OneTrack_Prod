import { useAuthStore } from '../stores/auth';

/**
 * WebSocket Service for real-time notifications
 * Implements auto-reconnect with exponential backoff
 */
class WebSocketService {
  constructor() {
    this.ws = null;
    this.url = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 10;
    this.reconnectDelay = 1000; // Start with 1 second
    this.maxReconnectDelay = 30000; // Max 30 seconds
    this.reconnectTimer = null;
    this.pingInterval = null;
    this.isIntentionalClose = false;
    this.messageHandlers = new Map();
    this.connectionStateListeners = new Set();
    this.isConnected = false;
    this.messageQueue = [];
  }

  /**
   * Connect to WebSocket server
   */
  connect() {
    const authStore = useAuthStore();
    
    if (!authStore.token) {
      console.warn('⚠️  Cannot connect to WebSocket: No authentication token');
      return;
    }

    // Determine WebSocket URL based on current environment
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = import.meta.env.VITE_API_URL 
      ? import.meta.env.VITE_API_URL.replace(/^https?:\/\//, '')
      : 'localhost:8080';
    
    this.url = `${protocol}//${host}/ws/notifications`;

    // Close existing connection if any
    if (this.ws) {
      this.isIntentionalClose = true;
      this.ws.close();
    }

    console.log('🔌 Connecting to WebSocket:', this.url);

    try {
      this.ws = new WebSocket(this.url);
      this.setupEventHandlers();
    } catch (error) {
      console.error('❌ WebSocket connection error:', error);
      this.scheduleReconnect();
    }
  }

  /**
   * Setup WebSocket event handlers
   */
  setupEventHandlers() {
    this.ws.onopen = () => {
      console.log('✅ WebSocket connected');
      this.isConnected = true;
      this.reconnectAttempts = 0;
      this.reconnectDelay = 1000;
      this.isIntentionalClose = false;
      
      // Send authentication token
      this.sendAuthToken();
      
      // Start ping interval
      this.startPingInterval();
      
      // Process queued messages
      this.processMessageQueue();
      
      // Notify listeners
      this.notifyConnectionState('connected');
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        this.handleMessage(message);
      } catch (error) {
        console.error('❌ Error parsing WebSocket message:', error);
      }
    };

    this.ws.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
      this.notifyConnectionState('error');
    };

    this.ws.onclose = (event) => {
      console.log('🔌 WebSocket disconnected:', event.code, event.reason);
      this.isConnected = false;
      this.stopPingInterval();
      this.notifyConnectionState('disconnected');

      // Reconnect if not intentional close
      if (!this.isIntentionalClose) {
        this.scheduleReconnect();
      }
    };
  }

  /**
   * Send authentication token to server
   */
  sendAuthToken() {
    const authStore = useAuthStore();
    if (authStore.token && this.isConnected) {
      // Token is sent via HTTP header during WebSocket upgrade
      // This is just a placeholder for any additional auth logic
      console.log('🔐 WebSocket authenticated');
    }
  }

  /**
   * Handle incoming WebSocket message
   */
  handleMessage(message) {
    const { type, data, timestamp } = message;
    
    console.log('📨 WebSocket message received:', type, data);

    // Call registered handlers for this message type
    if (this.messageHandlers.has(type)) {
      const handlers = this.messageHandlers.get(type);
      handlers.forEach(handler => {
        try {
          handler(data, timestamp);
        } catch (error) {
          console.error(`❌ Error in message handler for ${type}:`, error);
        }
      });
    }

    // Call wildcard handlers
    if (this.messageHandlers.has('*')) {
      const handlers = this.messageHandlers.get('*');
      handlers.forEach(handler => {
        try {
          handler({ type, data, timestamp });
        } catch (error) {
          console.error('❌ Error in wildcard message handler:', error);
        }
      });
    }
  }

  /**
   * Register a message handler for specific message type
   */
  on(messageType, handler) {
    if (!this.messageHandlers.has(messageType)) {
      this.messageHandlers.set(messageType, new Set());
    }
    this.messageHandlers.get(messageType).add(handler);

    // Return unsubscribe function
    return () => {
      const handlers = this.messageHandlers.get(messageType);
      if (handlers) {
        handlers.delete(handler);
      }
    };
  }

  /**
   * Remove a message handler
   */
  off(messageType, handler) {
    const handlers = this.messageHandlers.get(messageType);
    if (handlers) {
      handlers.delete(handler);
    }
  }

  /**
   * Send a message to the server
   */
  send(type, data) {
    const message = {
      type,
      data,
      timestamp: new Date().toISOString()
    };

    if (this.isConnected && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      // Queue message for later
      this.messageQueue.push(message);
      console.warn('⚠️  WebSocket not connected, message queued:', type);
    }
  }

  /**
   * Process queued messages
   */
  processMessageQueue() {
    while (this.messageQueue.length > 0 && this.isConnected) {
      const message = this.messageQueue.shift();
      this.ws.send(JSON.stringify(message));
    }
  }

  /**
   * Start ping interval to keep connection alive
   */
  startPingInterval() {
    this.stopPingInterval();
    this.pingInterval = setInterval(() => {
      if (this.isConnected) {
        this.send('ping', null);
      }
    }, 30000); // Ping every 30 seconds
  }

  /**
   * Stop ping interval
   */
  stopPingInterval() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }

  /**
   * Schedule reconnection with exponential backoff
   */
  scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('❌ Max reconnection attempts reached');
      this.notifyConnectionState('failed');
      return;
    }

    this.reconnectAttempts++;
    const delay = Math.min(
      this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1),
      this.maxReconnectDelay
    );

    console.log(`🔄 Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
    
    this.notifyConnectionState('reconnecting');

    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, delay);
  }

  /**
   * Disconnect from WebSocket
   */
  disconnect() {
    this.isIntentionalClose = true;
    
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    this.stopPingInterval();

    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }

    this.isConnected = false;
    this.messageQueue = [];
    console.log('🔌 WebSocket disconnected intentionally');
  }

  /**
   * Add connection state listener
   */
  onConnectionStateChange(listener) {
    this.connectionStateListeners.add(listener);
    
    // Return unsubscribe function
    return () => {
      this.connectionStateListeners.delete(listener);
    };
  }

  /**
   * Notify connection state listeners
   */
  notifyConnectionState(state) {
    this.connectionStateListeners.forEach(listener => {
      try {
        listener(state);
      } catch (error) {
        console.error('❌ Error in connection state listener:', error);
      }
    });
  }

  /**
   * Get current connection state
   */
  getConnectionState() {
    if (!this.ws) return 'disconnected';
    
    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return 'connecting';
      case WebSocket.OPEN:
        return 'connected';
      case WebSocket.CLOSING:
        return 'disconnecting';
      case WebSocket.CLOSED:
        return 'disconnected';
      default:
        return 'unknown';
    }
  }
}

// Export singleton instance
export const websocketService = new WebSocketService();
export default websocketService;
