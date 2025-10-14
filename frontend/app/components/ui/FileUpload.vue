<template>
  <div class="bg-white shadow rounded-lg">
    <div class="px-4 py-5 sm:p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Upload File</h3>
      
      <!-- File Upload Area -->
      <div
        ref="dropZone"
        :class="dropZoneClasses"
        @dragover.prevent="handleDragOver"
        @dragleave.prevent="handleDragLeave"
        @drop.prevent="handleDrop"
        @click="triggerFileInput"
      >
        <input
          ref="fileInput"
          type="file"
          class="hidden"
          accept=".csv,.xlsx,.xls,.json,.xml,.txt,.pdf"
          @change="handleFileSelect"
        />
        
        <div class="text-center">
          <Icon 
            name="heroicons:cloud-arrow-up" 
            class="mx-auto h-12 w-12 text-gray-400" 
          />
          <div class="mt-4">
            <p class="text-sm text-gray-600">
              <span class="font-medium text-indigo-600 hover:text-indigo-500 cursor-pointer">
                Click to upload
              </span>
              or drag and drop
            </p>
            <p class="text-xs text-gray-500 mt-1">
              PDF, CSV, Excel, JSON, XML, TXT files
            </p>
          </div>
        </div>
      </div>

      <!-- Selected File Display -->
      <div v-if="selectedFile" class="mt-4 p-4 bg-gray-50 rounded-lg">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <Icon name="heroicons:document" class="h-8 w-8 text-indigo-600" />
            <div>
              <p class="text-sm font-medium text-gray-900">{{ selectedFile.name }}</p>
              <p class="text-xs text-gray-500">{{ formatFileSize(selectedFile.size) }}</p>
            </div>
          </div>
          <button
            @click="clearFile"
            class="text-gray-400 hover:text-gray-600"
          >
            <Icon name="heroicons:x-mark" class="h-5 w-5" />
          </button>
        </div>
      </div>

      <!-- Upload Button -->
      <div v-if="selectedFile" class="mt-6">
        <UiButton
          variant="primary"
          size="lg"
          :loading="isUploading"
          :disabled="!selectedFile"
          icon="heroicons:arrow-up-tray"
          block
          @click="handleUpload"
        >
          {{ isUploading ? 'Extracting...' : 'Extract Data' }}
        </UiButton>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

// ===== PROPS =====
// Define what this component receives from parent
const props = defineProps({
  isUploading: {
    type: Boolean,
    default: false
  }
})

// ===== EMITS =====
// Define what this component sends to parent
const emit = defineEmits([
  'file-selected',    // When user selects a file
  'file-cleared',     // When user clears the file
  'upload-requested'  // When user clicks upload button
])

// ===== REACTIVE STATE =====
const selectedFile = ref(null)
const isDragOver = ref(false)

// ===== REFS =====
const fileInput = ref(null)
const dropZone = ref(null)

// ===== COMPUTED =====
const dropZoneClasses = computed(() => [
  'border-2 border-dashed rounded-lg p-12 text-center cursor-pointer transition-colors',
  isDragOver.value 
    ? 'border-indigo-500 bg-indigo-50' 
    : 'border-gray-300 hover:border-gray-400'
])

// ===== METHODS =====

// Drag and Drop Methods
const handleDragOver = (e) => {
  e.preventDefault()
  isDragOver.value = true
}

const handleDragLeave = (e) => {
  e.preventDefault()
  isDragOver.value = false
}

const handleDrop = (e) => {
  e.preventDefault()
  isDragOver.value = false
  
  const files = e.dataTransfer.files
  if (files.length > 0) {
    handleFileSelect({ target: { files } })
  }
}

// File Selection Methods
const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    selectedFile.value = file
    // Tell parent component that a file was selected
    emit('file-selected', file)
  }
}

const clearFile = () => {
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
  // Tell parent component that file was cleared
  emit('file-cleared')
}

// Utility Methods
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Upload Method
const handleUpload = () => {
  if (selectedFile.value) {
    // Tell parent component that upload was requested
    emit('upload-requested', selectedFile.value)
  }
}
</script>