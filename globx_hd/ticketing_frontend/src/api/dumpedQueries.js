import axios from 'axios';

const API_URL = 'http://localhost:8080';

// Helper to get auth header
function authHeader() {
  const token = localStorage.getItem('token');
  return { headers: { Authorization: `Bearer ${token}` } };
}

export async function fetchDumpedQueries(page = 1, limit = 10, status = '') {
  const params = { page, limit };
  if (status) params.status = status;
  
  const response = await axios.get(`${API_URL}/manager/dumped-queries`, {
    params,
    ...authHeader()
  });
  return response.data;
}

export async function fetchDumpedQuery(id) {
  const response = await axios.get(`${API_URL}/manager/dumped-queries/${id}`, authHeader());
  return response.data;
}

export async function updateDumpedQueryStatus(id, status) {
  const response = await axios.patch(`${API_URL}/manager/dumped-queries/${id}/status`, { status }, authHeader());
  return response.data;
}

export async function deleteDumpedQuery(id) {
  const response = await axios.delete(`${API_URL}/manager/dumped-queries/${id}`, authHeader());
  return response.data;
}
