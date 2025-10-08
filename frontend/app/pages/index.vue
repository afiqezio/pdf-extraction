<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
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
    <div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
      <!-- Upload Section -->
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
                  CSV, XLSX, JSON, XML, TXT, PDF files up to 500MB
                </p>
              </div>
            </div>
          </div>

          <!-- Selected File Info -->
          <div v-if="selectedFile" class="mt-4 p-4 bg-gray-50 rounded-lg">
            <div class="flex items-center justify-between">
              <div class="flex items-center">
                <Icon 
                  :name="getFileIcon(selectedFile.type)" 
                  class="h-8 w-8 text-gray-400 mr-3" 
                />
                <div>
                  <p class="text-sm font-medium text-gray-900">{{ selectedFile.name }}</p>
                  <p class="text-xs text-gray-500">
                    {{ formatFileSize(selectedFile.size) }} • {{ selectedFile.type }}
                  </p>
                </div>
              </div>
              <button
                type="button"
                class="text-gray-400 hover:text-gray-600"
                @click="removeFile"
              >
                <Icon name="heroicons:x-mark" class="h-5 w-5" />
              </button>
            </div>
          </div>

          <!-- Extraction Options -->
          <!-- <div v-if="selectedFile" class="mt-6">
            <h4 class="text-sm font-medium text-gray-900 mb-3">Extraction Options</h4>
            <div class="space-y-3">
              <div class="flex items-center">
                <input
                  id="auto-detect"
                  v-model="extractionOptions.autoDetect"
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label for="auto-detect" class="ml-2 text-sm text-gray-700">
                  Auto-detect data structure
                </label>
              </div>
              <div class="flex items-center">
                <input
                  id="include-metadata"
                  v-model="extractionOptions.includeMetadata"
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label for="include-metadata" class="ml-2 text-sm text-gray-700">
                  Include file metadata
                </label>
              </div>
              <div class="flex items-center">
                <input
                  id="validate-data"
                  v-model="extractionOptions.validateData"
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label for="validate-data" class="ml-2 text-sm text-gray-700">
                  Validate extracted data
                </label>
              </div>
            </div>
          </div> -->

          <!-- Upload Button -->
          <div v-if="selectedFile" class="mt-6">
            <UiButton
              variant="primary"
              size="lg"
              :loading="isUploading"
              :disabled="!selectedFile"
              icon="heroicons:arrow-up-tray"
              block
              @click="uploadAndExtract"
            >
              {{ isUploading ? 'Extracting...' : 'Extract Data' }}
            </UiButton>
          </div>
        </div>
      </div>

      <!-- Results Section -->
      <div class="bg-white shadow rounded-lg">
        <div class="px-4 py-5 sm:p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-medium text-gray-900">Extraction Results</h3>
          </div>

          <!-- Loading State -->
          <div v-if="isUploading" class="text-center py-12">
            <UiLoading size="lg" message="Processing your file..." />
          </div>

          <!-- Empty State -->
          <div v-else-if="!extractionResult" class="text-center py-12">
            <Icon name="heroicons:document-magnifying-glass" class="mx-auto h-12 w-12 text-gray-400" />
            <h3 class="mt-2 text-sm font-medium text-gray-900">No extraction results</h3>
            <p class="mt-1 text-sm text-gray-500">
              Upload a file to see extraction results here
            </p>
          </div>

          <!-- Results Content -->
          <div v-else class="space-y-4">
            <!-- Summary Stats -->
            <!-- <div class="grid grid-cols-2 gap-4">
              <div class="bg-blue-50 p-4 rounded-lg">
                <div class="flex items-center">
                  <Icon name="heroicons:table-cells" class="h-6 w-6 text-blue-600" />
                  <div class="ml-3">
                    <p class="text-sm font-medium text-blue-900">Records Found</p>
                    <p class="text-2xl font-bold text-blue-600">{{ extractionResult.summary.records }}</p>
                  </div>
                </div>
              </div>
              <div class="bg-green-50 p-4 rounded-lg">
                <div class="flex items-center">
                  <Icon name="heroicons:list-bullet" class="h-6 w-6 text-green-600" />
                  <div class="ml-3">
                    <p class="text-sm font-medium text-green-900">Fields Detected</p>
                    <p class="text-2xl font-bold text-green-600">{{ extractionResult.summary.fields }}</p>
                  </div>
                </div>
              </div>
            </div> -->

            <!-- Data Preview -->
            <div>
              <h4 class="text-sm font-medium text-gray-900 mb-3">Data Preview</h4>
              <div class="border border-gray-200 rounded-lg overflow-hidden">
                <div class="overflow-x-auto">
                  <table class="min-w-full divide-y divide-gray-200">
                    <!-- Hidden header row for column delete buttons -->
                    <thead class="bg-transparent">
                      <tr>
                        <th class="px-3 py-2 h-10 w-12">
                          <!-- Empty header for delete row column -->
                        </th>
                        <th
                          v-for="(field, columnIndex) in extractionResult.preview.fields"
                          :key="columnIndex"
                          class="px-3 py-2 h-10"
                        >
                          <button
                            @click="deleteColumn(columnIndex)"
                            class="text-red-500 hover:text-red-700 bg-white border border-red-200 rounded-full p-1.5 hover:bg-red-50 transition-colors shadow-sm w-8 h-8 flex items-center justify-center mx-auto"
                            :title="`Delete column ${columnIndex + 1}`"
                          >
                            <Icon name="heroicons:x-mark" class="h-3 w-3" />
                          </button>
                        </th>
                      </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                      <tr
                        v-for="(row, index) in extractionResult.preview.rows"
                        :key="index"
                        class="hover:bg-gray-50"
                      >
                        <!-- Delete row button column -->
                        <td class="px-3 py-2 whitespace-nowrap text-sm text-gray-900 w-12">
                          <button
                            @click="deleteRow(index)"
                            class="text-red-500 hover:text-red-700 bg-white border border-red-200 rounded-full p-1.5 hover:bg-red-50 transition-colors shadow-sm w-8 h-8 flex items-center justify-center"
                            :title="`Delete row ${index + 1}`"
                          >
                            <Icon name="heroicons:trash" class="h-3 w-3" />
                          </button>
                        </td>
                        <!-- Data columns -->
                        <td
                          v-for="(value, fieldIndex) in row"
                          :key="fieldIndex"
                          class="px-3 py-2 whitespace-nowrap text-sm text-gray-900"
                        >
                          {{ value }}
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            <!-- Field Analysis -->
            <!-- <div>
              <h4 class="text-sm font-medium text-gray-900 mb-3">Field Analysis</h4>
              <div class="space-y-2">
                <div
                  v-for="field in extractionResult.analysis"
                  :key="field.name"
                  class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
                >
                  <div>
                    <p class="text-sm font-medium text-gray-900">{{ field.name }}</p>
                    <p class="text-xs text-gray-500">{{ field.type }} • {{ field.nullCount }} nulls</p>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span
                      :class="getFieldTypeBadgeClass(field.type)"
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ field.type }}
                    </span>
                  </div>
                </div>
              </div>
            </div> -->

            <!-- Warnings/Errors -->
            <!-- <div v-if="extractionResult.warnings?.length > 0" class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <div class="flex">
                <Icon name="heroicons:exclamation-triangle" class="h-5 w-5 text-yellow-400" />
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-yellow-800">Warnings</h3>
                  <div class="mt-2 text-sm text-yellow-700">
                    <ul class="list-disc list-inside space-y-1">
                      <li v-for="warning in extractionResult.warnings" :key="warning">
                        {{ warning }}
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
            </div> -->
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Extractions -->
    <!-- <div v-if="recentExtractions.length > 0" class="mt-8">
      <div class="bg-white shadow rounded-lg">
        <div class="px-4 py-5 sm:p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Recent Extractions</h3>
          <div class="space-y-3">
            <div
              v-for="extraction in recentExtractions"
              :key="extraction.id"
              class="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50"
            >
              <div class="flex items-center">
                <Icon 
                  :name="getFileIcon(extraction.fileType)" 
                  class="h-8 w-8 text-gray-400 mr-3" 
                />
                <div>
                  <p class="text-sm font-medium text-gray-900">{{ extraction.fileName }}</p>
                  <p class="text-xs text-gray-500">
                    {{ extraction.records }} records • {{ formatDate(extraction.createdAt) }}
                  </p>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <span
                  :class="getStatusBadgeClass(extraction.status)"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                >
                  {{ extraction.status }}
                </span>
                <UiButton
                  variant="ghost"
                  size="sm"
                  icon="heroicons:eye"
                  @click="viewExtraction(extraction.id)"
                >
                  View
                </UiButton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div> -->
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'


// Reactive state
const selectedFile = ref(null)
const isUploading = ref(false)
const isDragOver = ref(false)
const extractionResult = ref(null)
const recentExtractions = ref([])

// File input refs
const fileInput = ref(null)
const dropZone = ref(null)

// Extraction options
const extractionOptions = ref({
  autoDetect: true,
  includeMetadata: true,
  validateData: true
})

// Computed
const dropZoneClasses = computed(() => {
  const base = 'border-2 border-dashed rounded-lg p-6 text-center cursor-pointer transition-colors duration-200'
  const state = isDragOver.value 
    ? 'border-indigo-400 bg-indigo-50' 
    : 'border-gray-300 hover:border-gray-400'
  return `${base} ${state}`
})

// Methods
const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    validateAndSetFile(file)
  }
}

const handleDragOver = (event) => {
  event.preventDefault()
  isDragOver.value = true
}

const handleDragLeave = (event) => {
  event.preventDefault()
  isDragOver.value = false
}

const handleDrop = (event) => {
  event.preventDefault()
  isDragOver.value = false
  
  const files = event.dataTransfer.files
  if (files.length > 0) {
    validateAndSetFile(files[0])
  }
}

const validateAndSetFile = (file) => {
  // File size validation (10MB limit)
  const maxSize = 500 * 1024 * 1024
  if (file.size > maxSize) {
    alert('File size must be less than 500MB')
    return
  }

  // File type validation
  const allowedTypes = [
    'text/csv',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/json',
    'application/xml',
    'text/xml',
    'text/plain',
    'application/pdf'
  ]

  if (!allowedTypes.includes(file.type) && !file.name.match(/\.(csv|xlsx|xls|json|xml|txt|pdf)$/i)) {
    alert('Please select a valid file type (CSV, XLSX, JSON, XML, TXT, PDF)')
    return
  }

  selectedFile.value = file
}

const removeFile = () => {
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const uploadAndExtract = async () => {
  if (!selectedFile.value) return

  isUploading.value = true
  extractionResult.value = null

  try {
    // Create FormData for file upload
    const formData = new FormData()
    formData.append('file', selectedFile.value)

    console.log("calling process pdf");
    // Call the real API endpoint
    const response = await fetch('http://localhost:8081/api/v1/extraction/process-pdf', {
      method: 'POST',
      body: formData
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Extraction failed')
    }

    const data = await response.json()

    console.log("response data", data);
    
    // Process the real extraction results
    const markdownFiles = data.results.markdown_files || []
    const firstFile = markdownFiles[0]
    
    if (firstFile && firstFile.content) {
      // Parse the markdown content to extract table information
      const content = firstFile.content
      const tableMatches = content.match(/### Table \d+/g) || []
      const tableCount = tableMatches.length
      
      // Extract sample data from the first table
      const tableStart = content.indexOf('|')
      const tableEnd = content.indexOf('---', tableStart)
      const tableContent = content.substring(tableStart, tableEnd)
      const tableLines = tableContent.split('\n').filter(line => line.includes('|'))
      
      // Parse table headers and first few rows
      const headers = tableLines[0]?.split('|').map(h => h.trim()).filter(h => h) || []
      const rows = tableLines.slice(1, 4).map(line => 
        line.split('|').map(cell => cell.trim()).filter(cell => cell)
      ).filter(row => row.length > 0)

      extractionResult.value = {
        summary: {
          records: tableCount,
          fields: headers.length,
          tables: tableCount
        },
        preview: {
          fields: headers,
          rows: rows
        },
        analysis: headers.map(header => ({
          name: header,
          type: 'String',
          nullCount: 0
        })),
        warnings: [],
        extractionOutput: data.results.extraction_output,
        filename: data.filename,
        processedAt: data.processed_at
      }
    } else {
      // Fallback if no tables found
      extractionResult.value = {
        summary: {
          records: 0,
          fields: 0,
          tables: 0
        },
        preview: {
          fields: [],
          rows: []
        },
        analysis: [],
        warnings: ['No tables found in the PDF'],
        extractionOutput: data.results.extraction_output,
        filename: data.filename,
        processedAt: data.processed_at
      }
    }

    // Add to recent extractions
    recentExtractions.value.unshift({
      id: Date.now(),
      fileName: selectedFile.value.name,
      fileType: selectedFile.value.type,
      records: extractionResult.value.summary.records,
      status: 'Completed',
      createdAt: new Date()
    })

  } catch (error) {
    console.error('Extraction failed:', error)
    alert(`Extraction failed: ${error.message}`)
  } finally {
    isUploading.value = false
  }
}

// Utility functions
const getFileIcon = (fileType) => {
  if (fileType.includes('csv') || fileType.includes('excel')) return 'heroicons:table-cells'
  if (fileType.includes('json')) return 'heroicons:code-bracket'
  if (fileType.includes('xml')) return 'heroicons:document-text'
  if (fileType.includes('pdf')) return 'heroicons:document'
  return 'heroicons:document'
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const deleteRow = (rowIndex) => {
  if (confirm('Are you sure you want to delete this row?')) {
    extractionResult.value.preview.rows.splice(rowIndex, 1)
    // Update summary
    extractionResult.value.summary.records = extractionResult.value.preview.rows.length
  }
}

const deleteColumn = (columnIndex) => {
  if (confirm('Are you sure you want to delete this column?')) {
    // Remove the field from fields array
    extractionResult.value.preview.fields.splice(columnIndex, 1)
    
    // Remove the corresponding data from each row
    extractionResult.value.preview.rows.forEach(row => {
      row.splice(columnIndex, 1)
    })
    
    // Update summary
    extractionResult.value.summary.fields = extractionResult.value.preview.fields.length
  }
}

// Lifecycle
onMounted(() => {
  // Load recent extractions
})

// SEO
useHead({
  title: 'Data Extraction - Workbench Platform',
  meta: [
    { name: 'description', content: 'Upload and extract data from various file formats with intelligent analysis and validation.' }
  ]
})
</script>
