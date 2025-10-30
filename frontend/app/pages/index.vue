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

              <!-- Right Side - Table Editor (60%) -->
              <div class="col-span-7">
                <UiContentViewer
                  :all-blocks="extractionResult?.allBlocks || []"
                  @update:blocks="handleBlocksUpdated"
                />
              </div>
            </div>

            <!-- Save to Database Section -->
            <div
              v-if="extractionResult.allTables && extractionResult.allTables.length > 0"
              class="mt-8 p-6 bg-gradient-to-r from-indigo-50 to-purple-50 rounded-lg border border-indigo-200"
            >
              <div class="flex items-center justify-between mb-4">
                <div>
                  <h4 class="text-lg font-semibold text-gray-900">Save to Database</h4>
                  <p class="text-sm text-gray-600 mt-1">
                    Save {{ extractionResult.allTables.length }} table(s) to the database
                  </p>
                </div>
                <div class="text-right">
                  <div class="text-2xl font-bold text-indigo-600">
                    {{
                      extractionResult.allTables.reduce((sum, t) => sum + (t.rows?.length || 0), 0)
                    }}
                  </div>
                  <div class="text-xs text-gray-500">Total Records</div>
                </div>
              </div>

              <!-- Progress Steps -->
              <div v-if="isSaving" class="mb-6">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-medium text-gray-700">Saving Progress</span>
                  <span class="text-sm text-gray-500"
                    >{{ saveProgress.current }}/{{ saveProgress.total }}</span
                  >
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
                        index < saveProgress.current
                          ? 'bg-green-500 text-white'
                          : index === saveProgress.current
                            ? 'bg-indigo-500 text-white animate-pulse'
                            : 'bg-gray-300 text-gray-600',
                      ]"
                    >
                      {{ index + 1 }}
                    </div>
                    <span
                      :class="[
                        'text-sm',
                        index < saveProgress.current
                          ? 'text-green-700'
                          : index === saveProgress.current
                            ? 'text-indigo-700'
                            : 'text-gray-500',
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
                  <p>
                    Ready to save {{ extractionResult.allTables.length }} table(s) with
                    {{
                      extractionResult.allTables.reduce((sum, t) => sum + (t.rows?.length || 0), 0)
                    }}
                    total records
                  </p>
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
  total: 0,
})
const saveSteps = ref(['Validating data', 'Mapping fields', 'Saving to database', 'Finalizing'])

// Methods

// Handle events from FileUpload component
const handleFileSelected = (file) => {
  console.log('File selected:', file.name)
  selectedFile.value = file // Keep your existing selectedFile logic
}

const handleFileCleared = () => {
  console.log('File cleared')
  selectedFile.value = null // Keep your existing clear logic
}

const handleUploadRequested = (file) => {
  console.log('Upload requested for:', file.name)
  uploadAndExtract() // Call your existing upload function
}

// ADD THIS NEW FUNCTION HERE:
const handleBlocksUpdated = (updatedBlocks) => {
  console.log('📝 Blocks updated from ContentViewer:', updatedBlocks)
  
  if (extractionResult.value) {
    extractionResult.value.allBlocks = updatedBlocks
  }
}

const uploadAndExtract = async () => {
  if (!selectedFile.value) return

  isUploading.value = true

  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)

    const response = await fetch('http://localhost:8081/api/v1/extraction/process-pdf', {
      method: 'POST',
      body: formData,
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    console.log('🔍 Full backend response:', result)

    // NEW: Reducto backend returns allBlocks directly
    if (result.allBlocks && result.allBlocks.length > 0) {
      console.log('📊 Found blocks:', result.allBlocks.length)

      extractionResult.value = {
        filename: result.filename,
        allBlocks: result.allBlocks, // All content blocks
      }

      // Set PDF URL for viewer
      // pdfUrl.value = `http://localhost:8081/api/v1/extraction/pdf/${result.filename}`

      // Find first block with page number
      const firstBlock = result.allBlocks.find((b) => b.page)
      if (firstBlock) {
        // currentTablePage.value = firstBlock.page
      }

      alert(`✅ Success! Found ${result.allBlocks.length} content block(s)`)
    } else {
      console.log('⚠️ No tables found')
      console.log('🔍 Response keys:', Object.keys(result))

      extractionResult.value = {
        filename: result.filename,
        allTables: [],
      }

      alert('⚠️ No tables extracted from PDF')
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
    extractionResult.value.allTables[tableIndex] = {
      ...extractionResult.value.allTables[tableIndex],
    }
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
        table.rows.forEach((row) => {
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

// ===== MERGE/UNMERGE HANDLERS =====
const handleTableUnmerged = (tableIndex) => {
  if (!extractionResult.value?.allTables?.[tableIndex]) {
    console.error('Table not found for unmerge')
    return
  }

  const mergedTable = extractionResult.value.allTables[tableIndex]

  // Check if table is actually merged
  if (!mergedTable.merged_from_tables || mergedTable.merged_from_tables.length === 0) {
    alert('This table is not merged')
    return
  }

  // TODO: Implement unmerge logic
  // For now, just alert the user
  alert(
    `Unmerge feature coming soon!\n\nThis table was merged from: ${mergedTable.merged_from_tables.join(', ')}\n\nYou can manually delete rows or columns if needed.`
  )

  console.log('Unmerge requested for table:', tableIndex)
  console.log('Original tables:', mergedTable.merged_from_tables)
}

const handleTablesMerged = (data) => {
  const { table1Index, table2Index } = data

  if (
    !extractionResult.value?.allTables?.[table1Index] ||
    !extractionResult.value?.allTables?.[table2Index]
  ) {
    alert('Invalid tables selected for merge')
    return
  }

  const table1 = extractionResult.value.allTables[table1Index]
  const table2 = extractionResult.value.allTables[table2Index]

  // Simple horizontal merge
  try {
    const mergedTable = {
      id: `${table1.id}_merged_${table2.id}`,
      page: table1.page,
      method: table1.method,
      confidence: Math.min(table1.confidence || 0, table2.confidence || 0),
      dimensions: `${Math.max(table1.rows?.length || 0, table2.rows?.length || 0)}x${(table1.headers?.length || 0) + (table2.headers?.length || 0)}`,

      // Combine headers
      headers: [...(table1.headers || []), ...(table2.headers || [])],
      normalized_headers: [
        ...(table1.normalized_headers || table1.headers || []),
        ...(table2.normalized_headers || table2.headers || []),
      ],
      original_headers: [
        ...(table1.original_headers || table1.headers || []),
        ...(table2.original_headers || table2.headers || []),
      ],

      // Merge rows
      rows: [],

      // Metadata
      merged_from_tables: [table1.id, table2.id],
      merge_confidence: 100, // Manual merge = 100% confidence
      merge_info: {
        table1_id: table1.id,
        table2_id: table2.id,
        table1_rows: table1.rows?.length || 0,
        table2_rows: table2.rows?.length || 0,
        manual_merge: true,
      },
    }

    // Merge rows row-by-row
    const maxRows = Math.max(table1.rows?.length || 0, table2.rows?.length || 0)
    for (let i = 0; i < maxRows; i++) {
      const row1 = table1.rows?.[i] || Array(table1.headers?.length || 0).fill('')
      const row2 = table2.rows?.[i] || Array(table2.headers?.length || 0).fill('')
      mergedTable.rows.push([...row1, ...row2])
    }

    // Replace table1 with merged table and remove table2
    extractionResult.value.allTables[table1Index] = mergedTable
    extractionResult.value.allTables.splice(table2Index, 1)

    // Select the merged table
    selectedTableIndex.value = table1Index

    console.log('✅ Tables merged successfully')
    alert(`✅ Tables ${table1Index + 1} and ${table2Index + 1} merged successfully!`)
  } catch (error) {
    console.error('Merge failed:', error)
    alert('Failed to merge tables. Please try again.')
  }
}

// ===== ADD/INSERT ROW/COLUMN HANDLERS =====
const handleRowAdded = (data) => {
  const { tableIndex } = data

  if (!extractionResult.value?.allTables?.[tableIndex]) {
    console.error('Table not found')
    return
  }

  const table = extractionResult.value.allTables[tableIndex]
  const columnCount = table.headers?.length || 0

  // Create new empty row
  const newRow = Array(columnCount).fill('')

  // Add at the end
  if (!table.rows) {
    table.rows = []
  }
  table.rows.push(newRow)

  handleTableUpdated(tableIndex)
  console.log(`✅ Added new row at end of table ${tableIndex + 1}`)
}

const handleColumnAdded = (data) => {
  const { tableIndex } = data

  if (!extractionResult.value?.allTables?.[tableIndex]) {
    console.error('Table not found')
    return
  }

  const table = extractionResult.value.allTables[tableIndex]

  // Add new header
  const newColumnName = `Column ${(table.headers?.length || 0) + 1}`
  if (!table.headers) {
    table.headers = []
  }
  table.headers.push(newColumnName)

  // Add empty cell to each row
  if (table.rows) {
    table.rows.forEach((row) => {
      row.push('')
    })
  }

  handleTableUpdated(tableIndex)
  console.log(`✅ Added new column at end of table ${tableIndex + 1}`)
}

const handleRowInserted = (data) => {
  const { tableIndex, rowIndex, position } = data

  if (!extractionResult.value?.allTables?.[tableIndex]) {
    console.error('Table not found')
    return
  }

  const table = extractionResult.value.allTables[tableIndex]
  const columnCount = table.headers?.length || 0

  // Create new empty row
  const newRow = Array(columnCount).fill('')

  // Insert at position
  const insertIndex = position === 'above' ? rowIndex : rowIndex + 1
  table.rows.splice(insertIndex, 0, newRow)

  handleTableUpdated(tableIndex)
  console.log(`✅ Inserted new row ${position} row ${rowIndex + 1} in table ${tableIndex + 1}`)
}

const handleColumnInserted = (data) => {
  const { tableIndex, columnIndex, position } = data

  if (!extractionResult.value?.allTables?.[tableIndex]) {
    console.error('Table not found')
    return
  }

  const table = extractionResult.value.allTables[tableIndex]

  // Insert new header
  const insertIndex = position === 'left' ? columnIndex : columnIndex + 1
  const newColumnName = `Column ${insertIndex + 1}`
  table.headers.splice(insertIndex, 0, newColumnName)

  // Insert empty cell in each row
  if (table.rows) {
    table.rows.forEach((row) => {
      row.splice(insertIndex, 0, '')
    })
  }

  handleTableUpdated(tableIndex)
  console.log(
    `✅ Inserted new column ${position} column ${columnIndex + 1} in table ${tableIndex + 1}`
  )
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
    await new Promise((resolve) => setTimeout(resolve, 500))

    // Step 2: Mapping fields
    saveProgress.value.current = 2
    await new Promise((resolve) => setTimeout(resolve, 500))

    const response = await fetch('http://localhost:8081/api/v1/extraction/save-to-db', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        tables: extractionResult.value.allTables,
      }),
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    // Step 3: Saving to database
    saveProgress.value.current = 3
    await new Promise((resolve) => setTimeout(resolve, 1000))

    const result = await response.json()
    console.log('Save completed:', result)

    // Step 4: Finalizing
    saveProgress.value.current = 4
    await new Promise((resolve) => setTimeout(resolve, 500))

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
