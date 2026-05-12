import api from './api.js';

// Upload files to a ticket
export async function uploadTicketAttachments(ticketId, files) {
  console.log('🚨 DEBUG: uploadTicketAttachments called with:', { ticketId, fileCount: files.length });
  
  const formData = new FormData();
  
  // Add files to FormData
  files.forEach((file, index) => {
    console.log(`🔥 DEBUG: Adding file ${index + 1}: ${file.name}, size: ${file.size}`);
    formData.append('files', file);
  });
  
  const url = `/tickets/${ticketId}/attachments`;
  console.log('🚨 DEBUG: Making POST request to:', url);
  
  const response = await api.post(url, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
  
  console.log('✅ DEBUG: Upload response:', response.data);
  return response.data;
}

// Get all attachments for a ticket
export async function getTicketAttachments(ticketId) {
  const response = await api.get(`/tickets/${ticketId}/attachments`);
  return response.data;
}

// Download an attachment
export async function downloadAttachment(attachmentId) {
  const response = await api.get(`/attachments/${attachmentId}/download`, {
    responseType: 'blob'
  });
  return response;
}

// Delete an attachment
export async function deleteAttachment(attachmentId) {
  const response = await api.delete(`/attachments/${attachmentId}`);
  return response.data;
}

// Helper function to create download URL for blob
export function createDownloadUrl(blob, filename) {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}
