<template>
  <div class="file-upload-container">
    <!-- File Upload Area -->
    <div 
      class="file-upload-area"
      :class="{ 'drag-over': isDragOver, 'has-files': selectedFiles.length > 0 }"
      @drop="handleDrop"
      @dragover.prevent="isDragOver = true"
      @dragleave.prevent="isDragOver = false"
      @click="triggerFileInput"
    >
      <input
        ref="fileInput"
        type="file"
        multiple
        accept=".jpg,.jpeg,.png,.gif,.webp,.pdf,.doc,.docx,.txt,.rtf,.xls,.xlsx,.csv"
        @change="handleFileSelect"
        class="hidden"
      />
      
      <div class="upload-content">
        <svg class="upload-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
        </svg>
        <p class="upload-text">
          <span class="font-semibold">Click to upload</span> or drag and drop
        </p>
        <p class="upload-subtext">
          Images, Documents up to 3MB each
        </p>
        <p class="allowed-types">
          Allowed: JPG, PNG, PDF, DOC, DOCX, TXT, XLS, XLSX, CSV
        </p>
      </div>
    </div>

    <!-- Selected Files List -->
    <div v-if="selectedFiles.length > 0" class="selected-files">
      <h4 class="files-title">Selected Files ({{ selectedFiles.length }})</h4>
      <div class="files-list">
        <div 
          v-for="(file, index) in selectedFiles" 
          :key="index"
          class="file-item"
          :class="{ 'file-error': file.error }"
        >
          <div class="file-info">
            <svg class="file-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            <div class="file-details">
              <p class="file-name">{{ file.name }}</p>
              <p class="file-size">{{ formatFileSize(file.size) }}</p>
              <p v-if="file.error" class="file-error-text">{{ file.error }}</p>
            </div>
          </div>
          <button 
            @click="removeFile(index)"
            class="remove-file-btn"
            type="button"
          >
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Upload Button -->
    <div v-if="selectedFiles.length > 0 && !isUploading && showUploadButton" class="upload-actions">
      <button
        @click="uploadFiles()"
        class="upload-button"
        :disabled="getValidFiles().length === 0"
      >
        <svg class="upload-button-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 12l2 2 4-4"></path>
        </svg>
        Upload {{ getValidFiles().length }} File{{ getValidFiles().length !== 1 ? 's' : '' }}
      </button>
    </div>

    <!-- Upload Progress -->
    <div v-if="isUploading" class="upload-progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
      </div>
      <p class="progress-text">Uploading files... {{ uploadProgress }}%</p>
    </div>

    <!-- Error Messages -->
    <div v-if="uploadErrors.length > 0" class="error-messages">
      <h4 class="error-title">Upload Errors:</h4>
      <ul class="error-list">
        <li v-for="(error, index) in uploadErrors" :key="index" class="error-item">
          {{ error }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FileUpload',
  props: {
    ticketId: {
      type: String,
      default: null
    },
    showUploadButton: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      selectedFiles: [],
      isDragOver: false,
      isUploading: false,
      uploadProgress: 0,
      uploadErrors: [],
      maxFileSize: 3 * 1024 * 1024, // 3MB
      allowedTypes: ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.pdf', '.doc', '.docx', '.txt', '.rtf', '.xls', '.xlsx', '.csv']
    }
  },
  methods: {
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    
    handleFileSelect(event) {
      const files = Array.from(event.target.files);
      this.processFiles(files);
    },
    
    handleDrop(event) {
      event.preventDefault();
      this.isDragOver = false;
      const files = Array.from(event.dataTransfer.files);
      this.processFiles(files);
    },
    
    processFiles(files) {
      const validFiles = [];
      
      files.forEach(file => {
        const fileObj = {
          file: file,
          name: file.name,
          size: file.size,
          error: null
        };
        
        // Validate file
        const validation = this.validateFile(file);
        if (validation.isValid) {
          validFiles.push(fileObj);
        } else {
          fileObj.error = validation.error;
          validFiles.push(fileObj);
        }
      });
      
      this.selectedFiles = [...this.selectedFiles, ...validFiles];
      this.$emit('files-selected', this.getValidFiles());
    },
    
    validateFile(file) {
      // Check file size
      if (file.size > this.maxFileSize) {
        return { isValid: false, error: 'File size exceeds 3MB limit' };
      }
      
      // Check file type
      const fileName = file.name.toLowerCase();
      const isValidType = this.allowedTypes.some(type => fileName.endsWith(type));
      if (!isValidType) {
        return { isValid: false, error: 'File type not allowed' };
      }
      
      return { isValid: true };
    },
    
    removeFile(index) {
      this.selectedFiles.splice(index, 1);
      this.$emit('files-selected', this.getValidFiles());
    },
    
    getValidFiles() {
      return this.selectedFiles.filter(f => !f.error).map(f => f.file);
    },
    
    formatFileSize(bytes) {
      if (bytes === 0) return '0 Bytes';
      const k = 1024;
      const sizes = ['Bytes', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    },
    
    async uploadFiles(ticketId = null) {
      const finalTicketId = ticketId || this.ticketId;
      console.log('🔥 DEBUG: uploadFiles called with ticketId:', finalTicketId);
      if (!finalTicketId) {
        this.$emit('upload-error', 'Ticket ID is required for file upload');
        return;
      }
      
      const validFiles = this.getValidFiles();
      console.log('🔥 DEBUG: validFiles count:', validFiles.length);
      if (validFiles.length === 0) {
        this.$emit('upload-error', 'No valid files to upload');
        return;
      }
      
      this.isUploading = true;
      this.uploadProgress = 0;
      this.uploadErrors = [];
      
      try {
        const { uploadTicketAttachments } = await import('../../api/attachments');
        
        // Simulate progress for now (real progress tracking would need axios interceptor)
        const progressInterval = setInterval(() => {
          if (this.uploadProgress < 90) {
            this.uploadProgress += 10;
          }
        }, 100);
        
        const response = await uploadTicketAttachments(finalTicketId, validFiles);
        
        clearInterval(progressInterval);
        this.uploadProgress = 100;
        
        this.isUploading = false;
        this.selectedFiles = [];
        
        if (response.errors && response.errors.length > 0) {
          this.uploadErrors = response.errors;
        }
        
        this.$emit('upload-success', response);
        
      } catch (error) {
        this.isUploading = false;
        const errorMessage = error.response?.data?.error || 'Failed to upload files';
        this.uploadErrors = [errorMessage];
        this.$emit('upload-error', errorMessage);
      }
    },
    
    clearFiles() {
      this.selectedFiles = [];
      this.uploadErrors = [];
      this.$emit('files-selected', []);
    }
  }
}
</script>

<style scoped>
.file-upload-container {
  @apply space-y-4;
}

.file-upload-area {
  @apply border-2 border-dashed border-gray-300 rounded-lg p-8 text-center cursor-pointer transition-colors duration-200;
}

.file-upload-area:hover {
  @apply border-blue-400 bg-blue-50;
}

.file-upload-area.drag-over {
  @apply border-blue-500 bg-blue-100;
}

.file-upload-area.has-files {
  @apply border-green-400 bg-green-50;
}

.upload-content {
  @apply space-y-2;
}

.upload-icon {
  @apply w-12 h-12 mx-auto text-gray-400;
}

.upload-text {
  @apply text-sm text-gray-600;
}

.upload-subtext {
  @apply text-xs text-gray-500;
}

.allowed-types {
  @apply text-xs text-gray-400 mt-1;
}

.selected-files {
  @apply space-y-3;
}

.files-title {
  @apply text-sm font-medium text-gray-700;
}

.files-list {
  @apply space-y-2;
}

.file-item {
  @apply flex items-center justify-between p-3 bg-gray-50 rounded-lg border;
}

.file-item.file-error {
  @apply bg-red-50 border-red-200;
}

.file-info {
  @apply flex items-center space-x-3;
}

.file-icon {
  @apply w-5 h-5 text-gray-500;
}

.file-details {
  @apply space-y-1;
}

.file-name {
  @apply text-sm font-medium text-gray-700;
}

.file-size {
  @apply text-xs text-gray-500;
}

.file-error-text {
  @apply text-xs text-red-600;
}

.remove-file-btn {
  @apply p-1 text-gray-400 hover:text-red-500 transition-colors;
}

.remove-file-btn svg {
  @apply w-4 h-4;
}

.upload-progress {
  @apply space-y-2;
}

.progress-bar {
  @apply w-full bg-gray-200 rounded-full h-2;
}

.progress-fill {
  @apply bg-blue-600 h-2 rounded-full transition-all duration-300;
}

.progress-text {
  @apply text-sm text-gray-600 text-center;
}

.error-messages {
  @apply space-y-2;
}

.error-title {
  @apply text-sm font-medium text-red-700;
}

.error-list {
  @apply space-y-1;
}

.error-item {
  @apply text-sm text-red-600;
}

.upload-actions {
  @apply flex justify-center mt-4;
}

.upload-button {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors;
}

.upload-button-icon {
  @apply w-4 h-4 mr-2;
}

.hidden {
  display: none;
}
</style>
