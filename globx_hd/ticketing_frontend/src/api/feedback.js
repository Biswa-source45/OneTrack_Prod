import api from './api.js';

// payload: { title, description, pageUrl, image (File|Blob|null) }
export function createFeedback({ title, description, pageUrl, image }) {
  const fd = new FormData();
  fd.append('title', title);
  fd.append('description', description);
  if (pageUrl) fd.append('page_url', pageUrl);
  if (image) fd.append('image', image, image.name || 'screenshot.jpg');
  // Explicit multipart header: api.js defaults to JSON, which makes axios serialize FormData as JSON.
  return api.post('/feedback', fd, { headers: { 'Content-Type': 'multipart/form-data' }, timeout: 60000 }).then(r => r.data);
}

export function fetchMyFeedback() {
  return api.get('/feedback/mine').then(r => r.data);
}

// Super Admin only
export function fetchAllFeedback(status) {
  return api.get('/feedback', { params: status ? { status } : {} }).then(r => r.data);
}

export function updateFeedbackStatus(id, status) {
  return api.patch(`/feedback/${id}/status`, { status }).then(r => r.data);
}

// Image needs the auth header, so fetch as a blob and hand back an object URL.
export async function fetchFeedbackImageUrl(id) {
  const r = await api.get(`/feedback/${id}/image`, { responseType: 'blob' });
  return URL.createObjectURL(r.data);
}

export const FEEDBACK_STATUSES = ['OPEN', 'IN PROGRESS', 'COMPLETED'];

export const FEEDBACK_STATUS_CLASSES = {
  'OPEN': 'bg-amber-100 text-amber-800 border-amber-200',
  'IN PROGRESS': 'bg-blue-100 text-blue-800 border-blue-200',
  'COMPLETED': 'bg-green-100 text-green-800 border-green-200',
};
