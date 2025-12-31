<template>
  <div class="max-w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Page Header -->
    <div class="mb-8">
      <div class="md:flex md:items-center md:justify-between">
        <div class="flex-1 min-w-0">
          <h1 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Data Extraction
          </h1>
          <p class="mt-1 text-sm text-gray-500">
            Upload files to extract and analyze data from various formats
          </p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="space-y-8">
      <!-- Upload Section -->
      <FileUpload
        :is-uploading="isUploading"
        @file-selected="handleFileSelected"
        @file-cleared="handleFileCleared"
        @upload-requested="handleUploadRequested"
      />

      <!-- Results Section -->
      <div class="bg-white shadow rounded-lg">
        <div class="px-4 py-5 sm:p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-medium text-gray-900">Extraction Results</h3>
             <div v-if="jobStatus" class="flex items-center space-x-2">
                <span class="text-sm text-gray-500">Status:</span>
                <span :class="statusClass" class="px-2 py-1 rounded-full text-xs font-medium uppercase">{{ jobStatus }}</span>
            </div>
          </div>

          <!-- Loading State -->
          <div v-if="isUploading || (jobStatus && jobStatus !== 'completed' && jobStatus !== 'failed')" class="text-center py-12">
            <UiLoading size="lg" message="Processing your file..." />
            <p class="mt-2 text-sm text-gray-500">This may take a few minutes for large documents.</p>
          </div>

          <!-- Error State -->
          <div v-else-if="jobStatus === 'failed'" class="text-center py-12">
             <Icon name="heroicons:exclamation-circle" class="mx-auto h-12 w-12 text-red-400" />
             <h3 class="mt-2 text-sm font-medium text-gray-900">Extraction Failed</h3>
             <p class="mt-1 text-sm text-red-500">{{ extractionError }}</p>
          </div>

          <!-- Empty State -->
          <div v-else-if="!extractionResult" class="text-center py-12">
            <Icon
              name="heroicons:document-magnifying-glass"
              class="mx-auto h-12 w-12 text-gray-400"
            />
            <h3 class="mt-2 text-sm font-medium text-gray-900">No extraction results</h3>
            <p class="mt-1 text-sm text-gray-500">Upload a file to see extraction results here</p>
          </div>

          <!-- Results Content -->
          <div v-else class="space-y-4">
            <!-- Split Screen Layout -->
            <div class="grid grid-cols-12 gap-6">
              <!-- Left Side - PDF Viewer (40%) -->
              <div class="col-span-5">
                <PdfViewer :pdf-url="pdfUrl" :current-page="currentTablePage" />
              </div>

              <!-- Right Side - Extraction Results (60%) -->
              <div class="col-span-7">
                <ExtractionResults :data="extractionResult.data || []" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import FileUpload from '~/components/ui/FileUpload.vue'
import PdfViewer from '~/components/ui/PdfViewer.vue'
import ExtractionResults from '~/components/ui/ExtractionResults.vue'

// Reactive state
const selectedFile = ref(null)
const isUploading = ref(false)
const extractionResult = ref(null)
const jobStatus = ref(null)
const extractionError = ref(null)

// Computed
const pdfUrl = computed(() => {
  if (!extractionResult.value?.file_path) return null
  // Extract filename from path
  const filename = extractionResult.value.file_path.split('/').pop()
  // Adjust port if needed, assuming backend is on 8081 based on previous code
  return `http://localhost:8081/api/v1/extraction/pdf/${filename}`
})

const currentTablePage = ref(1) // Simple page tracking for now

const statusClass = computed(() => {
    switch (jobStatus.value) {
        case 'completed': return 'bg-green-100 text-green-800'
        case 'processing': return 'bg-blue-100 text-blue-800'
        case 'failed': return 'bg-red-100 text-red-800'
        default: return 'bg-gray-100 text-gray-800'
    }
})

// Methods
const handleFileSelected = (file) => {
  selectedFile.value = file
}

const handleFileCleared = () => {
  selectedFile.value = null
}

const handleUploadRequested = () => {
  uploadAndExtract()
}

const uploadAndExtract = async () => {
  if (!selectedFile.value) return

  isUploading.value = true
  jobStatus.value = 'processing'
  extractionError.value = null
  extractionResult.value = null

  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)

    // 1. Upload and start job
    const response = await fetch('http://localhost:8081/api/v1/extraction/upload', {
      method: 'POST',
      body: formData,
    })

    if (!response.ok) {
        const err = await response.json()
        throw new Error(err.error || `HTTP error! status: ${response.status}`)
    }

    const { job_id } = await response.json()
    console.log('Job started:', job_id)

    // 2. Poll for status
    await pollJob(job_id)

  } catch (error) {
    console.error('Upload failed:', error)
    jobStatus.value = 'failed'
    extractionError.value = error.message
    alert('Upload failed: ' + error.message)
  } finally {
    isUploading.value = false
  }
}

const pollJob = async (jobId) => {
    const pollInterval = 2000 // 2 seconds
    
    const checkStatus = async () => {
        try {
            const response = await fetch(`http://localhost:8081/api/v1/extraction/jobs/${jobId}`)
            if (!response.ok) throw new Error('Failed to check job status')
            
            const job = await response.json()
            console.log('Job status:', job.status)
            jobStatus.value = job.status
            
            if (job.status === 'completed') {
                extractionResult.value = job
                return
            } else if (job.status === 'failed') {
                extractionError.value = job.error || 'Extraction failed'
                return
            }
            
            // Continue polling
            setTimeout(checkStatus, pollInterval)
        } catch (e) {
            console.error('Polling error:', e)
            jobStatus.value = 'failed'
            extractionError.value = e.message
        }
    }
    
    await checkStatus()
}
</script>
