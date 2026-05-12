import api from './api.js';

// Fetch all engineers with assigned (not CLOSED) tickets count
export async function fetchEngineersWithTickets() {
  const res = await api.get('/manager/engineers-with-tickets');
  console.log('Raw engineers API response:', res);
  return res.data?.engineers;
}

// Fetch tickets assigned to a specific engineer (not CLOSED)
export async function fetchEngineerAssignedTickets(engineerId) {
  const res = await api.get(`/manager/engineers/${engineerId}/tickets`);
  return res.data.tickets;
}
