import api from './api.js';

export function fetchContacts() {
  return api.get('/contacts').then(r => r.data);
}
export function createContact(payload) {
  return api.post('/contacts', payload).then(r => r.data);
}
export function updateContact(id, payload) {
  return api.put(`/contacts/${id}`, payload).then(r => r.data);
}
export function deleteContact(id) {
  return api.delete(`/contacts/${id}`).then(r => r.data);
}
export function fetchContactDesignations() {
  return api.get('/designations/contacts').then(r => r.data);
}
export function fetchAccounts() {
  return api.get('/accounts').then(r => r.data);
}
