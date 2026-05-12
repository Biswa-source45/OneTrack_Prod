// Engineer API integration for frontend
import api from './api.js';

export async function fetchEngineerTickets() {
  const res = await api.get('/engineer/tickets');
  return res.data;
}

export async function resolveEngineerTicket(ticketId) {
  const res = await api.patch(`/engineer/tickets/${ticketId}/status`, { status: 'RESOLVED' });
  return res.data;
}

// Change ticket status (engineer)
export async function changeEngineerTicketStatus(ticketId, status, remarks = '') {
  const res = await api.patch(`/engineer/tickets/${ticketId}/status`, { status, remarks });
  return res.data;
}

// ===== ENGINEER TASK API =====

// Fetch tasks assigned to engineer
export async function fetchEngineerTasks() {
  const res = await api.get('/engineer/tasks');
  return res.data;
}

// Get task by ID (engineer)
export async function fetchEngineerTask(taskId) {
  const res = await api.get(`/engineer/tasks/${taskId}`);
  return res.data;
}

// Change task status (engineer)
export async function changeEngineerTaskStatus(taskId, status) {
  const res = await api.patch(`/engineer/tasks/${taskId}/status`, { status });
  return res.data;
}
