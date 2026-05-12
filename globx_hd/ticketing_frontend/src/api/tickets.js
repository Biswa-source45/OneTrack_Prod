import api from './api.js';

export async function fetchProducts() {
  const res = await api.get('/products');
  return res.data;
}

export async function fetchIssues() {
  const res = await api.get('/issues');
  return res.data;
}

export async function createTicket(payload) {
  const res = await api.post('/tickets', payload);
  return res.data;
}

export async function fetchMyTickets(contactId) {
  const res = await api.get(`/tickets?contact_id=${contactId}`);
  return res.data;
}

export async function fetchTicketDetail(ticketId) {
  const res = await api.get(`/tickets/${ticketId}`);
  return res.data;
}

// Manager: Fetch all tickets
export async function fetchAllTickets() {
  const res = await api.get('/manager/tickets');
  return res.data;
}

// Manager: Fetch dashboard statistics
export async function fetchDashboardStats(month = null, year = null) {
  let url = '/manager/dashboard/stats';
  const params = new URLSearchParams();
  
  if (month && year) {
    params.append('month', month);
    params.append('year', year);
  }
  
  if (params.toString()) {
    url += '?' + params.toString();
  }
  
  const res = await api.get(url);
  return res.data;
}

// Manager: Assign engineer to ticket
export async function assignEngineer(ticketId, engineerId) {
  // Use updateTicket API for assigning/unassigning engineer
  const res = await api.put(`/manager/tickets/${ticketId}`, { assigned_engineer: engineerId });
  return res.data;
}

// Manager: Change ticket status
export async function changeTicketStatus(ticketId, status, remarks = '') {
  const res = await api.patch(`/manager/tickets/${ticketId}/status`, { status, remarks });
  return res.data;
}

// Manager: List engineers
export async function fetchEngineers() {
  const res = await api.get('/manager/engineers');
  return res.data;
}

// Manager: Update ticket (details, product, status, engineer)
export async function updateTicket(ticketId, payload) {
  const res = await api.put(`/manager/tickets/${ticketId}`, payload);
  return res.data;
}

// Manager: Create ticket on behalf of customer
export async function managerCreateTicket(payload) {
  const res = await api.post('/manager/tickets', payload);
  return res.data;
}

// Engineer: Create ticket on behalf of customer (unassigned - only managers can assign)
export async function engineerCreateTicket(payload) {
  const res = await api.post('/engineer/tickets', payload);
  return res.data;
}

// Fetch all users for ticket owner dropdown
export async function fetchUsers() {
  const res = await api.get('/users');
  return res.data;
}

// Enhanced ticket details with all related data
export async function fetchTicketFullDetails(ticketId) {
  const res = await api.get(`/tickets/${ticketId}/full`);
  return res.data;
}

// ===== TICKET COMMENTS API =====
export const ticketComments = {
  // Fetch all comments for a ticket
  fetch: async (ticketId, params = {}) => {
    const res = await api.get(`/tickets/${ticketId}/comments`, { params });
    return res.data;
  },
  
  // Create new comment or resolution
  create: async (ticketId, data) => {
    const res = await api.post(`/tickets/${ticketId}/comments`, data);
    return res.data;
  },
  
  // Update existing comment
  update: async (ticketId, commentId, data) => {
    const res = await api.put(`/tickets/${ticketId}/comments/${commentId}`, data);
    return res.data;
  },
  
  // Delete comment
  delete: async (ticketId, commentId) => {
    const res = await api.delete(`/tickets/${ticketId}/comments/${commentId}`);
    return res.data;
  }
};

// ===== TICKET CALLS API =====
export const ticketCalls = {
  // Fetch all calls for a ticket
  fetch: async (ticketId, params = {}) => {
    const res = await api.get(`/tickets/${ticketId}/calls`, { params });
    return res.data;
  },
  
  // Schedule new call
  schedule: async (ticketId, data) => {
    const res = await api.post(`/tickets/${ticketId}/calls`, data);
    return res.data;
  },
  
  // Update call details
  update: async (ticketId, callId, data) => {
    const res = await api.put(`/tickets/${ticketId}/calls/${callId}`, data);
    return res.data;
  },
  
  // Complete call
  complete: async (ticketId, callId, notes) => {
    const res = await api.patch(`/tickets/${ticketId}/calls/${callId}/complete`, { notes });
    return res.data;
  },
  
  // Cancel call
  cancel: async (ticketId, callId) => {
    const res = await api.patch(`/tickets/${ticketId}/calls/${callId}/cancel`);
    return res.data;
  },
  
  // Close call with remarks (managers only)
  close: async (ticketId, callId, remarks) => {
    const res = await api.patch(`/tickets/${ticketId}/calls/${callId}/close`, { remarks });
    return res.data;
  },
  
  // Upload attachments for a call
  uploadAttachments: async (ticketId, callId, formData) => {
    const res = await api.post(`/tickets/${ticketId}/calls/${callId}/attachments`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    return res.data;
  },
  
  // Get attachments for a call
  getAttachments: async (ticketId, callId) => {
    const res = await api.get(`/tickets/${ticketId}/calls/${callId}/attachments`);
    return res.data;
  },
  
  // Download call attachment
  downloadAttachment: async (attachmentId) => {
    const res = await api.get(`/call-attachments/${attachmentId}/download`, {
      responseType: 'blob'
    });
    return res;
  },
  
  // Delete call attachment
  deleteAttachment: async (attachmentId) => {
    const res = await api.delete(`/call-attachments/${attachmentId}`);
    return res.data;
  }
};

// ===== TICKET APPROVALS API =====
export const ticketApprovals = {
  // Fetch all approvals for a ticket
  fetch: async (ticketId, params = {}) => {
    const res = await api.get(`/tickets/${ticketId}/approvals`, { params });
    return res.data;
  },

  // Create a new approval request
  create: async (ticketId, approvalData) => {
    const res = await api.post(`/tickets/${ticketId}/approvals`, approvalData);
    return res.data;
  },

  // Approve an approval request
  approve: async (ticketId, approvalId, remarks) => {
    const res = await api.patch(`/tickets/${ticketId}/approvals/${approvalId}/approve`, { remarks });
    return res.data;
  },

  // Reject an approval request
  reject: async (ticketId, approvalId, remarks) => {
    const res = await api.patch(`/tickets/${ticketId}/approvals/${approvalId}/reject`, { remarks });
    return res.data;
  }
};

// Helper functions for easier access
export const fetchTicketComments = (ticketId, params = {}) => ticketComments.fetch(ticketId, params)
export const fetchTicketCalls = (ticketId, params = {}) => ticketCalls.fetch(ticketId, params)
export const fetchTicketApprovals = (ticketId, params = {}) => ticketApprovals.fetch(ticketId, params)

// ===== TICKET ACTIVITIES API =====
export const ticketActivities = {
  // Fetch activity history
  fetch: async (ticketId, params = {}) => {
    const res = await api.get(`/tickets/${ticketId}/activities`, { params });
    return res.data;
  },
  
  // Fetch timeline view
  timeline: async (ticketId) => {
    const res = await api.get(`/tickets/${ticketId}/timeline`);
    return res.data;
  }
};

// ===== MANAGER TICKET DELETE API =====
export async function deleteTicket(ticketId) {
  const res = await api.delete(`/manager/tickets/${ticketId}`);
  return res.data;
}
