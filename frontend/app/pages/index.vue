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
            <!-- Split Screen Layout -->
            <div class="grid grid-cols-12 gap-6">
              <!-- Left Side - PDF Viewer (40%) -->
              <div class="col-span-5">
                <PdfViewer 
                  :pdf-url="pdfUrl"
                  :current-page="currentTablePage"
                />
              </div>

              <!-- Right Side - Table Editor (60%) -->
              <div class="col-span-7">
                <TableEditor 
                  :all-tables="extractionResult?.allTables || []"
                  :selected-table-index="selectedTableIndex"
                  @table-selected="handleTableSelected"
                  @table-updated="handleTableUpdated"
                  @table-deleted="handleTableDeleted"
                />
              </div>
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
                      {{ index + 1 }}
                    </div>
                    <span 
                      :class="[
                        'text-sm',
                        index < saveProgress.current ? 'text-green-700' :
                        index === saveProgress.current ? 'text-indigo-700' :
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
                  <p>Ready to save {{ extractionResult.allTables.length }} table(s) with {{ extractionResult.allTables.reduce((sum, t) => sum + (t.rows?.length || 0), 0) }} total records</p>
                </div>
                <UiButton
                  variant="primary"
                  size="lg"
                  :loading="isSaving"
                  @click="saveToDatabase"
                  icon="heroicons:cloud-arrow-up"
                >
                  {{ isSaving ? 'Saving to Database...' : 'Save to Database' }}
                </UiButton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import FileUpload from '~/components/ui/FileUpload.vue'
import PdfViewer from '~/components/ui/PdfViewer.vue'
import TableEditor from '~/components/ui/TableEditor.vue'

// Reactive state
const selectedFile = ref(null)
const isUploading = ref(false)
const extractionResult = ref(null)
const recentExtractions = ref([])

// Split-screen functionality
const selectedTableIndex = ref(0)
const selectedTable = computed(() => {
  if (!extractionResult.value?.allTables || extractionResult.value.allTables.length === 0) {
    return null
  }
  return extractionResult.value.allTables[selectedTableIndex.value]
})

const currentTablePage = computed(() => {
  if (!selectedTable.value) return null
  return selectedTable.value.page || 1
})

const pdfUrl = computed(() => {
  if (!extractionResult.value?.filename) return null
  const baseUrl = `http://localhost:8081/api/v1/extraction/pdf/${extractionResult.value.filename}`
  
  // Add page anchor to jump to the specific page in the PDF viewer
  if (currentTablePage.value) {
    return `${baseUrl}#page=${currentTablePage.value}`
  }
  
  return baseUrl
})

// Save to database state
const isSaving = ref(false)
const saveProgress = ref({
  current: 0,
  total: 0
})
const saveSteps = ref([
  'Validating data',
  'Mapping fields',
  'Saving to database',
  'Finalizing'
])

// Methods

// Handle events from FileUpload component
const handleFileSelected = (file) => {
  console.log('File selected:', file.name)
  selectedFile.value = file  // Keep your existing selectedFile logic
}

const handleFileCleared = () => {
  console.log('File cleared')
  selectedFile.value = null  // Keep your existing clear logic
}

const handleUploadRequested = (file) => {
  console.log('Upload requested for:', file.name)
  uploadAndExtract()  // Call your existing upload function
}

const uploadAndExtract = async () => {
  if (!selectedFile.value) return

  isUploading.value = true
  
  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)

    const response = await fetch('http://localhost:8081/api/v1/extraction/process-pdf', {
      method: 'POST',
      body: formData
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    console.log('🔍 Full backend response:', result)
    
    // Extract the data from the response structure
    // The backend returns: { results: { json_files: [...] } }
    const jsonFiles = result.results?.json_files || result.json_files
    
    if (jsonFiles && jsonFiles.length > 0) {
      console.log('📊 Found JSON files:', jsonFiles.length)
      console.log('📋 First file data:', jsonFiles[0].data)
      console.log('📋 Tables in first file:', jsonFiles[0].data.tables)
      
      extractionResult.value = {
        filename: result.filename,
        allTables: jsonFiles[0].data.tables || []
      }
    } else {
      console.log('⚠️ No JSON files found in response')
      console.log('🔍 Available keys in result:', Object.keys(result))
      console.log('🔍 Available keys in result.results:', result.results ? Object.keys(result.results) : 'No results object')
      
      extractionResult.value = {
        filename: result.filename,
        allTables: []
      }
    }
    
    console.log('✅ Final extractionResult:', extractionResult.value)
    
    // Reset table selection
    selectedTableIndex.value = 0
  } catch (error) {
    console.error('Upload failed:', error)
    alert('Upload failed. Please try again.')
  } finally {
    isUploading.value = false
  }
}

// ===== TABLE EDITOR EVENT HANDLERS =====
const handleTableSelected = (tableIndex) => {
  selectedTableIndex.value = tableIndex
  console.log(`Selected table ${tableIndex + 1}`)
}

const handleTableUpdated = (tableIndex) => {
  if (extractionResult.value?.allTables?.[tableIndex]) {
    // Trigger reactivity update
    extractionResult.value.allTables[tableIndex] = { ...extractionResult.value.allTables[tableIndex] }
  }
}

const handleTableDeleted = (data) => {
  const { type, tableIndex, rowIndex, columnIndex } = data
  
  if (type === 'row') {
    if (extractionResult.value?.allTables?.[tableIndex]?.rows) {
      extractionResult.value.allTables[tableIndex].rows.splice(rowIndex, 1)
      handleTableUpdated(tableIndex)
    }
  } else if (type === 'column') {
    if (extractionResult.value?.allTables?.[tableIndex]) {
      const table = extractionResult.value.allTables[tableIndex]
      
      // Remove header
      if (table.headers) {
        table.headers.splice(columnIndex, 1)
      }
      
      // Remove column from all rows
      if (table.rows) {
        table.rows.forEach(row => {
          if (row[columnIndex] !== undefined) {
            row.splice(columnIndex, 1)
          }
        })
      }
      
      handleTableUpdated(tableIndex)
    }
  } else if (type === 'header') {
    if (extractionResult.value?.allTables?.[tableIndex]?.rows?.length > 0) {
      // Promote the first data row to be the new header
      const newHeader = extractionResult.value.allTables[tableIndex].rows[0]
      extractionResult.value.allTables[tableIndex].headers = [...newHeader]
      // Remove the promoted row from data (since it's now the header)
      extractionResult.value.allTables[tableIndex].rows.shift()
      
      handleTableUpdated(tableIndex)
    }
  } else if (type === 'table') {
    if (extractionResult.value?.allTables) {
      extractionResult.value.allTables.splice(tableIndex, 1)
      
      // Adjust selected table index if needed
      if (selectedTableIndex.value >= extractionResult.value.allTables.length) {
        selectedTableIndex.value = Math.max(0, extractionResult.value.allTables.length - 1)
      }
    }
  }
}

const saveToDatabase = async () => {
  if (!extractionResult.value?.allTables || extractionResult.value.allTables.length === 0) {
    alert('No tables to save')
    return
  }

  isSaving.value = true
  saveProgress.value = { current: 0, total: saveSteps.value.length }
  
  try {
    // Step 1: Validating data
    saveProgress.value.current = 1
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // Step 2: Mapping fields
    saveProgress.value.current = 2
    await new Promise(resolve => setTimeout(resolve, 500))
    
    const response = await fetch('http://localhost:8081/api/v1/extraction/save-to-db', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        tables: extractionResult.value.allTables
      })
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    // Step 3: Saving to database
    saveProgress.value.current = 3
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const result = await response.json()
    console.log('Save completed:', result)
    
    // Step 4: Finalizing
    saveProgress.value.current = 4
    await new Promise(resolve => setTimeout(resolve, 500))
    
    alert('Data saved to database successfully!')
  } catch (error) {
    console.error('Save failed:', error)
    alert('Save failed. Please try again.')
  } finally {
    isSaving.value = false
    saveProgress.value = { current: 0, total: 0 }
  }
}

// Lifecycle
onMounted(() => {
  // Any initialization code can go here
})
</script>