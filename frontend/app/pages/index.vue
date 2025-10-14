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
              <div class="col-span-7 bg-white shadow rounded-lg overflow-hidden">
                <div class="px-4 py-3 border-b border-gray-200 bg-gray-50">
                  <div class="flex items-center justify-between">
                    <h4 class="text-sm font-medium text-gray-900">Table Editor</h4>
                    <span class="text-xs text-gray-500">
                      {{ extractionResult.allTables?.length || 0 }} table(s) found
                    </span>
                  </div>
                </div>
                
                <!-- Table Navigation Sidebar -->
                <div class="px-4 py-2 border-b border-gray-200 bg-gray-50">
                  <div class="flex space-x-2 overflow-x-auto">
                    <button
                      v-for="(table, tableIndex) in extractionResult.allTables"
                      :key="tableIndex"
                      @click="selectTable(tableIndex)"
                      :class="[
                        'px-3 py-1 text-xs rounded-full border transition-colors',
                        selectedTableIndex === tableIndex
                          ? 'bg-indigo-100 text-indigo-700 border-indigo-300'
                          : 'bg-white text-gray-600 border-gray-300 hover:bg-gray-50'
                      ]"
                    >
                      Table {{ tableIndex + 1 }}
                    </button>
                  </div>
                </div>

                <!-- Selected Table Content -->
                <div class="overflow-auto">
                  <div v-if="selectedTable" class="p-4">

                    <!-- Single Table Display -->
                    <div v-if="selectedTable" class="border border-gray-200 rounded-lg overflow-hidden">
                      <!-- Table Header -->
                      <div class="bg-gray-50 px-4 py-3 border-b border-gray-200 flex items-center justify-between">
                        <div class="flex items-center space-x-3">
                          <Icon name="heroicons:table-cells" class="h-5 w-5 text-gray-500" />
                          <div>
                            <h5 class="text-sm font-semibold text-gray-900">
                              Table {{ selectedTableIndex + 1 }}
                            </h5>
                            <p class="text-xs text-gray-500">
                              {{ selectedTable.rows?.length || 0 }} rows × {{ selectedTable.headers?.length || 0 }} columns
                              <span v-if="selectedTable.confidence" class="ml-2">
                                • Confidence: {{ (selectedTable.confidence).toFixed(1) }}%
                              </span>
                            </p>
                          </div>
                        </div>
                        <button
                          @click="deleteTable(selectedTableIndex)"
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
                                  v-for="(header, columnIndex) in selectedTable.headers"
                                  :key="`header-btn-${columnIndex}`"
                                  class="px-3 py-2 h-10 bg-gray-50 min-w-[200px]"
                                >
                                  <button
                                    @click="deleteColumn(selectedTableIndex, columnIndex)"
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
                                  @click="deleteHeaderRow(selectedTableIndex)"
                                  class="text-red-500 hover:text-red-700 bg-white border border-red-200 rounded-full p-1.5 hover:bg-red-50 transition-colors shadow-sm w-8 h-8 flex items-center justify-center"
                                  :title="`Delete header row`"
                                >
                                  <Icon name="heroicons:trash" class="h-3 w-3" />
                                </button>
                              </th>
                              <th
                                v-for="(header, columnIndex) in selectedTable.headers"
                                :key="`header-${columnIndex}`"
                                class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider bg-gray-50 min-w-[200px]"
                              >
                                <input
                                  v-model="selectedTable.headers[columnIndex]"
                                  class="w-full min-w-[180px] bg-transparent border-none focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded px-1"
                                  @blur="updateTableData(selectedTableIndex)"
                                />
                              </th>
                            </tr>
                            </thead>
                            <tbody class="bg-white divide-y divide-gray-200">
                              <tr
                                v-for="(row, rowIndex) in selectedTable.rows"
                                :key="`row-${rowIndex}`"
                                class="hover:bg-gray-50 transition-colors"
                              >
                                <!-- Delete Row Button - Sticky -->
                                <td class="px-3 py-2 text-sm text-gray-500 w-12 bg-white sticky left-0 z-10 border-r border-gray-200">
                                  <button
                                    @click="deleteRow(selectedTableIndex, rowIndex)"
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
                                    v-model="selectedTable.rows[rowIndex][cellIndex]"
                                    class="w-full min-w-[180px] bg-transparent border border-transparent hover:border-gray-300 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 rounded px-2 py-1 transition-colors"
                                    @blur="updateTableData(selectedTableIndex)"
                                  />
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Empty state if no tables -->
                <div v-if="!extractionResult.allTables || extractionResult.allTables.length === 0" class="text-center py-8 border border-gray-200 rounded-lg">
                  <Icon name="heroicons:table-cells" class="mx-auto h-12 w-12 text-gray-400" />
                  <p class="mt-2 text-sm text-gray-500">No tables extracted</p>
                </div>
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
  
  // If we have a selected table with page information, show that specific page
  if (selectedTable.value?.page) {
    return `http://localhost:8081/api/v1/extraction/pdf/${extractionResult.value.filename}#page=${selectedTable.value.page}`
  }
  
  // Otherwise, show the full PDF
  return `http://localhost:8081/api/v1/extraction/pdf/${extractionResult.value.filename}`
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

// Split-screen functionality methods
const selectTable = (tableIndex) => {
  selectedTableIndex.value = tableIndex
  console.log(`Selected table ${tableIndex + 1}, navigating to page ${selectedTable.value?.page || 1}`)
}

// Table editing methods
const updateTableData = (tableIndex) => {
  if (extractionResult.value?.allTables?.[tableIndex]) {
    // Trigger reactivity update
    extractionResult.value.allTables[tableIndex] = { ...extractionResult.value.allTables[tableIndex] }
  }
}

const deleteRow = (tableIndex, rowIndex) => {
  if (extractionResult.value?.allTables?.[tableIndex]?.rows) {
    extractionResult.value.allTables[tableIndex].rows.splice(rowIndex, 1)
    updateTableData(tableIndex)
  }
}

const deleteColumn = (tableIndex, columnIndex) => {
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
    
    updateTableData(tableIndex)
  }
}

const deleteHeaderRow = (tableIndex) => {
  if (extractionResult.value?.allTables?.[tableIndex]?.rows?.length > 0) {
    // Remove the first row (header row)
    extractionResult.value.allTables[tableIndex].rows.shift()
    
    // Promote the new first row to be the header
    if (extractionResult.value.allTables[tableIndex].rows.length > 0) {
      const newHeader = extractionResult.value.allTables[tableIndex].rows[0]
      extractionResult.value.allTables[tableIndex].headers = [...newHeader]
      // Remove the promoted row from data
      extractionResult.value.allTables[tableIndex].rows.shift()
    }
    
    updateTableData(tableIndex)
  }
}

const deleteTable = (tableIndex) => {
  if (extractionResult.value?.allTables) {
    extractionResult.value.allTables.splice(tableIndex, 1)
    
    // Adjust selected table index if needed
    if (selectedTableIndex.value >= extractionResult.value.allTables.length) {
      selectedTableIndex.value = Math.max(0, extractionResult.value.allTables.length - 1)
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