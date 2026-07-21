import api from './api.js';

export function loginUser(data) {
  return api.post(`/login/user`, data);
}

export function loginContact(data) {
  return api.post(`/login/contact`, data);
}

export function resetUserPassword(data) {
  return api.post(`/reset-password/user`, data);
}

export function resetContactPassword(data) {
  return api.post(`/reset-password/contact`, data);
}

// Master data APIs (client-side pagination can fetch all and paginate locally for now)
export function fetchProducts() {
  return api.get(`/products`).then(r => r.data);
}
export function createProduct(payload) {
  return api.post(`/products`, payload).then(r => r.data);
}
export function updateProduct(id, payload) {
  return api.put(`/products/${id}`, payload).then(r => r.data);
}
export function deleteProduct(id) {
  return api.delete(`/products/${id}`).then(r => r.data);
}

export function fetchUserRoles() {
  return api.get(`/roles`).then(r => r.data);
}
export function createUserRole(payload) {
  return api.post(`/roles`, payload).then(r => r.data);
}

export function fetchUserDesignations() {
  return api.get(`/designations/users`).then(r => r.data);
}
export function createUserDesignation(payload) {
  return api.post(`/designations/users`, payload).then(r => r.data);
}

export function fetchContactDesignations() {
  return api.get(`/designations/contacts`).then(r => r.data);
}
export function createContactDesignation(payload) {
  return api.post(`/designations/contacts`, payload).then(r => r.data);
}

// Issues
export function fetchIssues(params) {
  return api.get(`/issues`, { params }).then(r => r.data);
}
export function createIssue(payload) {
  return api.post(`/issues`, payload).then(r => r.data);
}
export function updateIssue(id, payload) {
  return api.put(`/issues/${id}`, payload).then(r => r.data);
}
export function deleteIssue(id) {
  return api.delete(`/issues/${id}`).then(r => r.data);
}

// Users
export function fetchUsers() {
  return api.get(`/users`).then(r => r.data);
}
export function createUser(payload) {
  return api.post(`/users`, payload).then(r => r.data);
}
export function updateUser(id, payload) {
  return api.put(`/users/${id}`, payload).then(r => r.data);
}
export function deleteUser(id) {
  return api.delete(`/users/${id}`).then(r => r.data);
}

// Accounts
export function fetchAccounts() {
  return api.get(`/accounts`).then(r => r.data);
}
export function createAccount(payload) {
  return api.post(`/accounts`, payload).then(r => r.data);
}
export function updateAccount(id, payload) {
  return api.put(`/accounts/${id}`, payload).then(r => r.data);
}
export function deleteAccount(id) {
  return api.delete(`/accounts/${id}`).then(r => r.data);
}

// Contacts
export function fetchContacts() {
  return api.get(`/contacts`).then(r => r.data);
}
export function createContact(payload) {
  return api.post(`/contacts`, payload).then(r => r.data);
}
export function updateContact(id, payload) {
  return api.put(`/contacts/${id}`, payload).then(r => r.data);
}
export function deleteContact(id) {
  return api.delete(`/contacts/${id}`).then(r => r.data);
}

export function resetManagedUserPassword(id, payload) {
  return api.patch(`/users/${id}/password`, payload).then(r => r.data);
}
