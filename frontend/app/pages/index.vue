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
    <div class="space-y-8">
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

            <!-- Data Preview - All Tables -->
            <div class="space-y-6">
              <div class="flex items-center justify-between mb-3">
                <h4 class="text-sm font-medium text-gray-900">Data Preview</h4>
                <span class="text-xs text-gray-500">
                  {{ extractionResult.allTables?.length || 0 }} table(s) found
                </span>
              </div>

              <!-- Loop through all tables -->
              <div
                v-for="(table, tableIndex) in extractionResult.allTables"
                :key="tableIndex"
                class="border border-gray-200 rounded-lg overflow-hidden"
              >
                <!-- Table Header -->
                <div class="bg-gray-50 px-4 py-3 border-b border-gray-200 flex items-center justify-between">
                  <div class="flex items-center space-x-3">
                    <Icon name="heroicons:table-cells" class="h-5 w-5 text-gray-500" />
                    <div>
                      <h5 class="text-sm font-semibold text-gray-900">
                        Table {{ tableIndex + 1 }}
                      </h5>
                      <p class="text-xs text-gray-500">
                        {{ table.rows?.length || 0 }} rows × {{ table.headers?.length || 0 }} columns
                        <span v-if="table.confidence" class="ml-2">
                          • Confidence: {{ (table.confidence * 100).toFixed(1) }}%
                        </span>
                      </p>
                    </div>
                  </div>
                  <button
                    @click="deleteTable(tableIndex)"
                    class="text-red-500 hover:text-red-700 text-xs font-medium"
                  >
                    Delete Table
                  </button>
                </div>

                <!-- Table Data with Scroll Container -->
                <div class="relative">
                  <!-- Scroll indicator -->
                  <div class="absolute top-3 right-3 z-20 bg-indigo-100 text-indigo-700 px-3 py-1 rounded-full text-xs font-medium shadow-sm pointer-events-none opacity-75">
                    ← Scroll to see all columns →
                  </div>
                  
                  <div class="overflow-x-auto overflow-y-auto max-h-[600px] border-t border-gray-200" style="scrollbar-width: thin;">
                    <table class="w-full divide-y divide-gray-200" style="table-layout: auto;">
                      <!-- Column Delete Buttons -->
                      <thead class="bg-gray-50 sticky top-0 z-10 shadow-sm">
                        <tr>
                          <th class="px-3 py-2 h-10 w-12 bg-gray-50 sticky left-0 z-20 border-r border-gray-200">
                            <!-- Empty for row delete column -->
                          </th>
                          <th
                            v-for="(header, columnIndex) in table.headers"
                            :key="`header-btn-${columnIndex}`"
                            class="px-3 py-2 h-10 bg-gray-50 min-w-[200px]"
                          >
                            <button
                              @click="deleteColumn(tableIndex, columnIndex)"
                              class="text-red-500 hover:text-red-700 bg-white border border-red-200 rounded-full p-1.5 hover:bg-red-50 transition-colors shadow-sm w-8 h-8 flex items-center justify-center mx-auto"
                              :title="`Delete column ${columnIndex + 1}`"
                            >
                              <Icon name="heroicons:x-mark" class="h-3 w-3" />
                            </button>
                          </th>
                        </tr>
                      <!-- Header Row - Now Deletable -->
                      <tr class="hover:bg-gray-100">
                        <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider bg-gray-50 w-12 sticky left-0 z-20 border-r border-gray-200">
                          <button
                            @click="deleteHeaderRow(tableIndex)"
                            class="text-red-500 hover:text-red-700 bg-white border border-red-200 rounded-full p-1.5 hover:bg-red-50 transition-colors shadow-sm w-8 h-8 flex items-center justify-center"
                            :title="`Delete header row`"
                          >
                            <Icon name="heroicons:trash" class="h-3 w-3" />
                          </button>
                        </th>
                        <th
                          v-for="(header, columnIndex) in table.headers"
                          :key="`header-${columnIndex}`"
                          class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider bg-gray-50 min-w-[200px]"
                        >
                          <input
                            v-model="table.headers[columnIndex]"
                            class="w-full min-w-[180px] bg-transparent border-none focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded px-1"
                            @blur="updateTableData(tableIndex)"
                          />
                        </th>
                      </tr>
                      </thead>
                      <tbody class="bg-white divide-y divide-gray-200">
                        <tr
                          v-for="(row, rowIndex) in table.rows"
                          :key="`row-${rowIndex}`"
                          class="hover:bg-gray-50 transition-colors"
                        >
                          <!-- Delete Row Button - Sticky -->
                          <td class="px-3 py-2 text-sm text-gray-500 w-12 bg-white sticky left-0 z-10 border-r border-gray-200">
                            <button
                              @click="deleteRow(tableIndex, rowIndex)"
                              class="text-red-500 hover:text-red-700 bg-white border border-red-200 rounded-full p-1.5 hover:bg-red-50 transition-colors shadow-sm w-8 h-8 flex items-center justify-center"
                              :title="`Delete row ${rowIndex + 1}`"
                            >
                              <Icon name="heroicons:trash" class="h-3 w-3" />
                            </button>
                          </td>
                          <!-- Editable Data Cells -->
                          <td
                            v-for="(cell, cellIndex) in row"
                            :key="`cell-${cellIndex}`"
                            class="px-3 py-2 text-sm text-gray-900 min-w-[200px]"
                          >
                            <input
                              v-model="table.rows[rowIndex][cellIndex]"
                              class="w-full min-w-[180px] bg-transparent border border-transparent hover:border-gray-300 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 rounded px-2 py-1 transition-colors"
                              @blur="updateTableData(tableIndex)"
                            />
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>

              <!-- Empty state if no tables -->
              <div v-if="!extractionResult.allTables || extractionResult.allTables.length === 0" class="text-center py-8 border border-gray-200 rounded-lg">
                <Icon name="heroicons:table-cells" class="mx-auto h-12 w-12 text-gray-400" />
                <p class="mt-2 text-sm text-gray-500">No tables extracted</p>
              </div>

              <!-- Save to Database Section -->
              <div v-if="extractionResult.allTables && extractionResult.allTables.length > 0" class="mt-8 p-6 bg-gradient-to-r from-indigo-50 to-purple-50 rounded-lg border border-indigo-200">
                <div class="flex items-center justify-between mb-4">
                  <div>
                    <h4 class="text-lg font-semibold text-gray-900">Save to Database</h4>
                    <p class="text-sm text-gray-600 mt-1">
                      Save {{ extractionResult.allTables.length }} table(s) to the database
                    </p>
                  </div>
                  <div class="text-right">
                    <div class="text-2xl font-bold text-indigo-600">
                      {{ extractionResult.allTables.reduce((sum, t) => sum + (t.rows?.length || 0), 0) }}
                    </div>
                    <div class="text-xs text-gray-500">Total Records</div>
                  </div>
                </div>

                <!-- Progress Steps -->
                <div v-if="isSaving" class="mb-6">
                  <div class="flex items-center justify-between mb-2">
                    <span class="text-sm font-medium text-gray-700">Saving Progress</span>
                    <span class="text-sm text-gray-500">{{ saveProgress.current }}/{{ saveProgress.total }}</span>
                  </div>
                  <div class="w-full bg-gray-200 rounded-full h-2 mb-4">
                    <div 
                      class="bg-indigo-600 h-2 rounded-full transition-all duration-300 ease-out"
                      :style="{ width: `${(saveProgress.current / saveProgress.total) * 100}%` }"
                    ></div>
                  </div>
                  
                  <!-- Step Indicators -->
                  <div class="space-y-2">
                    <div 
                      v-for="(step, index) in saveSteps" 
                      :key="index"
                      class="flex items-center space-x-3"
                    >
                      <div 
                        :class="[
                          'w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium',
                          index < saveProgress.current ? 'bg-green-500 text-white' :
                          index === saveProgress.current ? 'bg-indigo-500 text-white animate-pulse' :
                          'bg-gray-300 text-gray-600'
                        ]"
                      >
                        <Icon 
                          v-if="index < saveProgress.current" 
                          name="heroicons:check" 
                          class="h-4 w-4" 
                        />
                        <Icon 
                          v-else-if="index === saveProgress.current" 
                          name="heroicons:arrow-path" 
                          class="h-4 w-4 animate-spin" 
                        />
                        <span v-else>{{ index + 1 }}</span>
                      </div>
                      <span 
                        :class="[
                          'text-sm',
                          index < saveProgress.current ? 'text-green-700' :
                          index === saveProgress.current ? 'text-indigo-700 font-medium' :
                          'text-gray-500'
                        ]"
                      >
                        {{ step }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Save Button -->
                <div class="flex items-center justify-between">
                  <div class="text-sm text-gray-600">
                    <span v-if="!isSaving">
                      Ready to save {{ extractionResult.allTables.filter(t => t.rows && t.rows.length > 0).length }} non-empty table(s)
                      <br>
                      <span class="text-xs text-indigo-600">💡 You can delete the header row to promote the first data row as new header</span>
                    </span>
                    <span v-else class="text-indigo-600">
                      {{ saveProgress.message }}
                    </span>
                  </div>
                  <UiButton
                    variant="primary"
                    size="lg"
                    :loading="isSaving"
                    :disabled="isSaving || !hasValidTables"
                    icon="heroicons:cloud-arrow-up"
                    @click="saveToDatabase"
                  >
                    {{ isSaving ? 'Saving...' : 'Save to Database' }}
                  </UiButton>
                </div>

                <!-- Success/Error Messages -->
                <div v-if="saveResult" class="mt-4">
                  <div 
                    :class="[
                      'p-4 rounded-lg flex items-center space-x-3',
                      saveResult.success ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'
                    ]"
                  >
                    <Icon 
                      :name="saveResult.success ? 'heroicons:check-circle' : 'heroicons:x-circle'"
                      :class="saveResult.success ? 'h-6 w-6 text-green-500' : 'h-6 w-6 text-red-500'"
                    />
                    <div>
                      <p :class="saveResult.success ? 'text-green-800 font-medium' : 'text-red-800 font-medium'">
                        {{ saveResult.message }}
                      </p>
                      <p v-if="saveResult.details" class="text-sm mt-1 text-gray-600">
                        {{ saveResult.details }}
                      </p>
                    </div>
                  </div>
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

// Save to database state
const isSaving = ref(false)
const saveProgress = ref({
  current: 0,
  total: 0,
  message: ''
})
const saveSteps = ref([
  'Validating tables...',
  'Mapping columns to database fields...',
  'Detecting table types...',
  'Preparing data for insertion...',
  'Inserting data into database...',
  'Finalizing save...'
])
const saveResult = ref(null)

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

const hasValidTables = computed(() => {
  if (!extractionResult.value?.allTables) return false
  return extractionResult.value.allTables.some(table => 
    table.rows && table.rows.length > 0 && table.headers && table.headers.length > 0
  )
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
    console.log("results object:", data.results);
    console.log("json_files:", data.results?.json_files);
    console.log("files_count:", data.results?.files_count);

    // --- Improved JSON files detection logic ---
    let jsonFiles = []
    // Defensive: check if data.results exists and is an object
    if (data && data.results && typeof data.results === 'object') {
      // Try to find json_files as array, or fallback to any array of objects with .data
      if (Array.isArray(data.results.json_files)) {
        jsonFiles = data.results.json_files
        console.log("Found json_files as array:", jsonFiles.length, "files");
      } else if (Array.isArray(data.results.files)) {
        // Sometimes backend may use 'files' instead of 'json_files'
        jsonFiles = data.results.files
        console.log("Found files as array:", jsonFiles.length, "files");
      } else if (typeof data.results.json_files === 'object' && data.results.json_files !== null) {
        // Sometimes json_files is a dict of filename: {data: ...}
        jsonFiles = Object.values(data.results.json_files)
        console.log("Found json_files as object:", jsonFiles.length, "files");
      } else {
        console.log("No json_files found in results");
      }
    } else {
      console.log("No results object found in response");
    }

    // Find the first file that has a .data property (robust to backend changes)
    const firstFile = Array.isArray(jsonFiles)
      ? jsonFiles.find(f => f && typeof f === 'object' && f.data)
      : null

    console.log("firstFile found:", firstFile);
    console.log("firstFile.data:", firstFile?.data);

    if (firstFile && firstFile.data) {
      const jsonData = firstFile.data
      const tables = Array.isArray(jsonData.tables) ? jsonData.tables : []
      const summary = typeof jsonData.summary === 'object' && jsonData.summary !== null ? jsonData.summary : {}

      // Process ALL tables with ALL rows (no slicing!)
      const allTables = tables.map(table => ({
        headers: Array.isArray(table.headers) ? [...table.headers] : [],
        rows: Array.isArray(table.rows) ? table.rows.map(row => [...row]) : [],
        confidence: table.confidence || null,
        dimensions: table.dimensions || null
      }))

      // Calculate total records across all tables
      const totalRecords = allTables.reduce((sum, table) => sum + (table.rows?.length || 0), 0)
      const totalFields = allTables.reduce((sum, table) => sum + (table.headers?.length || 0), 0)

      extractionResult.value = {
        summary: {
          records: summary.total_records || totalRecords,
          fields: summary.total_fields || totalFields,
          tables: summary.total_tables || tables.length
        },
        allTables: allTables, // Store all tables with all data
        warnings: [],
        extractionOutput: data.results.extraction_output,
        filename: data.filename,
        processedAt: data.processed_at,
        // Add the full JSON data for advanced features
        rawData: jsonData
      }

      console.log("✅ Loaded", allTables.length, "tables with", totalRecords, "total rows");
    } else {
      // Fallback if no JSON files found - show a message to the user
      console.log("No firstFile or firstFile.data found, using fallback");
      extractionResult.value = {
        summary: {
          records: 0,
          fields: 0,
          tables: 0
        },
        allTables: [],
        warnings: [
          'No valid JSON extraction data found. The PDF was processed but the backend did not return extractable data. Please check the backend logs and output structure.'
        ],
        extractionOutput: data.results && data.results.extraction_output,
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

const deleteRow = (tableIndex, rowIndex) => {
  if (confirm('Are you sure you want to delete this row?')) {
    const table = extractionResult.value.allTables[tableIndex]
    
    // Simply delete the row - no header promotion during editing
    table.rows.splice(rowIndex, 1)
    
    // Update summary - recalculate total records
    const totalRecords = extractionResult.value.allTables.reduce(
      (sum, t) => sum + (t.rows?.length || 0), 0
    )
    extractionResult.value.summary.records = totalRecords
  }
}

const deleteHeaderRow = (tableIndex) => {
  if (confirm('Are you sure you want to delete the header row? The first data row will become the new header.')) {
    const table = extractionResult.value.allTables[tableIndex]
    
    // If there are data rows, promote the first one to header
    if (table.rows && table.rows.length > 0) {
      // Use the first data row as the new header
      const newHeader = [...table.rows[0]]
      table.headers = newHeader
      
      // Remove the promoted row from data
      table.rows.splice(0, 1)
      
      console.log(`📋 Table ${tableIndex + 1}: Promoted first data row to header`)
    } else {
      // If no data rows, just clear the headers
      table.headers = []
      console.log(`📋 Table ${tableIndex + 1}: Cleared headers (no data rows)`)
    }
    
    // Update summary - recalculate total records
    const totalRecords = extractionResult.value.allTables.reduce(
      (sum, t) => sum + (t.rows?.length || 0), 0
    )
    extractionResult.value.summary.records = totalRecords
  }
}

const deleteColumn = (tableIndex, columnIndex) => {
  if (confirm('Are you sure you want to delete this column?')) {
    const table = extractionResult.value.allTables[tableIndex]
    
    // Remove the header
    table.headers.splice(columnIndex, 1)
    
    // Remove the corresponding data from each row
    table.rows.forEach(row => {
      row.splice(columnIndex, 1)
    })
    
    // Update summary - recalculate total fields
    const totalFields = extractionResult.value.allTables.reduce(
      (sum, t) => sum + (t.headers?.length || 0), 0
    )
    extractionResult.value.summary.fields = totalFields
  }
}

const deleteTable = (tableIndex) => {
  if (confirm('Are you sure you want to delete this entire table?')) {
    extractionResult.value.allTables.splice(tableIndex, 1)
    
    // Update summary
    const totalRecords = extractionResult.value.allTables.reduce(
      (sum, t) => sum + (t.rows?.length || 0), 0
    )
    const totalFields = extractionResult.value.allTables.reduce(
      (sum, t) => sum + (t.headers?.length || 0), 0
    )
    extractionResult.value.summary.records = totalRecords
    extractionResult.value.summary.fields = totalFields
    extractionResult.value.summary.tables = extractionResult.value.allTables.length
  }
}

const updateTableData = (tableIndex) => {
  // This function is called when a cell is edited
  // You can add any post-edit logic here (e.g., validation, auto-save)
  console.log(`Table ${tableIndex} updated`)
}

// Save to database functionality
const saveToDatabase = async () => {
  if (!hasValidTables.value) return

  isSaving.value = true
  saveResult.value = null
  saveProgress.value = {
    current: 0,
    total: saveSteps.value.length,
    message: 'Starting save process...'
  }

  try {
    // Step 1: Validate tables
    updateProgress(0, 'Validating tables...')
    await new Promise(resolve => setTimeout(resolve, 500)) // Simulate processing

    // Step 2: Map columns to database fields
    updateProgress(1, 'Mapping columns to database fields...')
    await new Promise(resolve => setTimeout(resolve, 800))

    // Step 3: Detect table types
    updateProgress(2, 'Detecting table types...')
    await new Promise(resolve => setTimeout(resolve, 600))

    // Step 4: Prepare data for insertion
    updateProgress(3, 'Preparing data for insertion...')
    await new Promise(resolve => setTimeout(resolve, 400))

    // Step 5: Insert data into database
    updateProgress(4, 'Inserting data into database...')
    
    // Filter out empty tables
    const validTables = extractionResult.value.allTables.filter(table => 
      table.rows && table.rows.length > 0 && table.headers && table.headers.length > 0
    )

    const response = await fetch('http://localhost:8081/api/v1/extraction/save-to-db', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        tables: validTables,
        filename: selectedFile.value?.name || 'unknown'
      })
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to save to database')
    }

    const result = await response.json()

    // Step 6: Finalize save
    updateProgress(5, 'Finalizing save...')
    await new Promise(resolve => setTimeout(resolve, 300))

    // Success
    saveResult.value = {
      success: true,
      message: `Successfully saved ${result.saved_tables || validTables.length} table(s) to database!`,
      details: result.details || `Total records saved: ${result.total_records || 0}`
    }

  } catch (error) {
    console.error('Save to database failed:', error)
    saveResult.value = {
      success: false,
      message: 'Failed to save to database',
      details: error.message
    }
  } finally {
    isSaving.value = false
  }
}

const updateProgress = (step, message) => {
  saveProgress.value.current = step
  saveProgress.value.message = message
}

// Promote first row to header for all tables before saving
const promoteFirstRowToHeader = () => {
  if (!extractionResult.value?.allTables) return
  
  extractionResult.value.allTables.forEach((table, tableIndex) => {
    if (table.rows && table.rows.length > 0) {
      // Use the first row as the new header
      const newHeader = [...table.rows[0]]
      table.headers = newHeader
      
      // Remove the first row from data (since it's now the header)
      table.rows.splice(0, 1)
      
      console.log(`📋 Table ${tableIndex + 1}: Promoted first row to header`)
    }
  })
  
  // Update summary after header promotion
  const totalRecords = extractionResult.value.allTables.reduce(
    (sum, t) => sum + (t.rows?.length || 0), 0
  )
  extractionResult.value.summary.records = totalRecords
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

<style scoped>
/* Custom scrollbar styling for better UX */
.overflow-x-auto::-webkit-scrollbar,
.overflow-y-auto::-webkit-scrollbar {
  height: 8px;
  width: 8px;
}

.overflow-x-auto::-webkit-scrollbar-track,
.overflow-y-auto::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 4px;
}

.overflow-x-auto::-webkit-scrollbar-thumb,
.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
  transition: background 0.2s;
}

.overflow-x-auto::-webkit-scrollbar-thumb:hover,
.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Smooth scrolling */
.overflow-x-auto,
.overflow-y-auto {
  scroll-behavior: smooth;
}

/* Better input focus states */
input:focus {
  outline: none;
}

/* Sticky column shadow effect */
td.sticky,
th.sticky {
  box-shadow: 2px 0 4px rgba(0, 0, 0, 0.05);
}
</style>
