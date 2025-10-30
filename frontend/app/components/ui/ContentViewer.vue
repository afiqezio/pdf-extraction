<template>
    <div class="bg-white shadow rounded-lg overflow-hidden">
        <!-- Header -->
        <div class="px-4 py-3 border-b border-gray-200 bg-gray-50">
        <div class="flex items-center justify-between">
            <h4 class="text-sm font-medium text-gray-900">Extracted Content</h4>
            <div class="flex items-center space-x-3">
            <span class="text-xs text-gray-500">
                {{ allBlocks?.length || 0 }} block(s) found
            </span>
            <button
                v-if="hasChanges"
                @click="saveAllChanges"
                class="px-3 py-1 bg-green-600 text-white text-xs rounded hover:bg-green-700"
            >
                Save All Changes
            </button>
            </div>
        </div>
        </div>

        <!-- Content Area -->
        <div class="overflow-auto max-h-screen p-4">
        <div v-if="!editableBlocks || editableBlocks.length === 0" class="text-center py-12">
            <Icon name="heroicons:document-text" class="mx-auto h-12 w-12 text-gray-400" />
            <h3 class="mt-2 text-sm font-medium text-gray-900">No content extracted</h3>
            <p class="mt-1 text-sm text-gray-500">Upload a file to see extracted content here</p>
        </div>

        <!-- Render each block by type -->
        <div v-else class="space-y-4">
            <div 
            v-for="(block, index) in editableBlocks" 
            :key="index"
            class="border rounded-lg p-4 transition-shadow"
            :class="block.isEditing ? 'ring-2 ring-blue-500 shadow-lg' : 'hover:shadow-md'"
            >
            <!-- Block Header -->
            <div class="flex items-center justify-between mb-2">
                <div class="flex items-center space-x-2">
                <span 
                    class="px-2 py-1 text-xs font-medium rounded"
                    :class="getTypeColor(block.type)"
                >
                    {{ block.type }}
                </span>
                <span class="text-xs text-gray-500">Page {{ block.page }}</span>
                <span 
                    v-if="block.confidence"
                    class="text-xs text-gray-500"
                >
                    {{ block.confidence }} confidence
                </span>
                <span 
                    v-if="block.isModified"
                    class="text-xs text-orange-600 font-medium"
                >
                    ● Modified
                </span>
                </div>
                
                <!-- Edit Controls -->
            <div class="flex items-center space-x-2">
                <!-- Delete Button - Always visible -->
                <button
                    @click="confirmDelete(index)"
                    class="px-2 py-1 text-xs bg-red-100 text-red-700 rounded hover:bg-red-200"
                    title="Delete this block"
                >
                    <Icon name="heroicons:trash" class="w-3 h-3" />
                </button>
                
                <!-- Edit/Save/Cancel Buttons -->
                <button
                    v-if="!block.isEditing"
                    @click="startEditing(index)"
                    class="px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded hover:bg-blue-200"
                >
                    Edit
                </button>
                <template v-else>
                    <button
                    @click="saveBlock(index)"
                    class="px-2 py-1 text-xs bg-green-100 text-green-700 rounded hover:bg-green-200"
                    >
                    Save
                    </button>
                    <button
                    @click="cancelEditing(index)"
                    class="px-2 py-1 text-xs bg-gray-100 text-gray-700 rounded hover:bg-gray-200"
                    >
                    Cancel
                    </button>
                </template>
                </div>
            </div>

            <!-- Block Content - Editable -->
            <div class="mt-2">
                <!-- Table - Special handling with Tabulator -->
                <div v-if="block.type === 'Table'">
                    <div v-if="block.isEditing" class="space-y-2">
                        <!-- Tabulator container -->
                        <div 
                        :id="`table-editor-${index}`"
                        class="border rounded"
                        ></div>
                        <p class="text-xs text-gray-500">
                        ✏️ Click any cell to edit. Changes are saved when you click "Save"
                        </p>
                    </div>
                    <div 
                        v-else
                        class="overflow-x-auto"
                    >
                        <div 
                        v-html="block.content" 
                        class="reducto-table"
                        ></div>
                    </div>
                </div>

                <!-- Text Content - Title, Header, Text, etc. -->
                <div v-else>
                <textarea
                    v-if="block.isEditing"
                    v-model="block.editContent"
                    class="w-full min-h-[100px] p-2 border rounded"
                    :class="getTextClass(block.type, true)"
                ></textarea>
                
                <div v-else :class="getTextClass(block.type, false)">
                    {{ block.content }}
                </div>
                </div>
            </div>
            </div>
        </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { TabulatorFull as Tabulator } from 'tabulator-tables'
import 'tabulator-tables/dist/css/tabulator.min.css'
import { useTableParser } from '~/composables/useTableParser'

const { parseHTMLTableToJSON, convertJSONToHTML } = useTableParser()

const props = defineProps({
allBlocks: {
    type: Array,
    default: () => []
}
})

const emit = defineEmits(['update:blocks'])

// Create editable copies of blocks
const editableBlocks = ref([])
const originalBlocks = ref([])
const tabulatorInstances = ref({})

// Initialize editable blocks when props change
watch(() => props.allBlocks, (newBlocks) => {
if (newBlocks && newBlocks.length > 0) {
    editableBlocks.value = newBlocks.map(block => ({
    ...block,
    isEditing: false,
    isModified: false,
    editContent: block.content,
    originalContent: block.content
    }))
    originalBlocks.value = JSON.parse(JSON.stringify(newBlocks))
}
}, { immediate: true })

// Check if there are any unsaved changes
const hasChanges = computed(() => {
return editableBlocks.value.some(block => block.isModified)
})

// Start editing a block
const startEditing = async (index) => {
  const block = editableBlocks.value[index]
  block.isEditing = true
  block.editContent = block.content
  
  // Special handling for tables
  if (block.type === 'Table') {
    await nextTick()
    
    console.log('Raw table HTML:', block.content)

    const tableElement = document.getElementById(`table-editor-${index}`)
    if (tableElement) {
      const tableData = parseHTMLTableToJSON(block.content)

      console.log("table data", tableData)
      
      tabulatorInstances.value[index] = new Tabulator(tableElement, {
        // data: tableData.rows,
        // columns: tableData.columns,
        data: tableData.columns,
        columns: tableData.columns,
        layout: 'fitData',
        height: 'auto',
        reactiveData: true,
        editable: true,
      })
    }
  }
}

// Save a single block
const saveBlock = (index) => {
  const block = editableBlocks.value[index]
  
  // Special handling for tables
  if (block.type === 'Table' && tabulatorInstances.value[index]) {
    const data = tabulatorInstances.value[index].getData()
    const columns = tabulatorInstances.value[index].getColumnDefinitions()
    
    // Convert back to HTML
    block.content = convertJSONToHTML(columns, data)
    
    // Destroy the tabulator instance
    tabulatorInstances.value[index].destroy()
    delete tabulatorInstances.value[index]
  } else {
    // Regular text content
    block.content = block.editContent
  }
  
  block.isEditing = false
  block.isModified = block.content !== block.originalContent
}

// Cancel editing
const cancelEditing = (index) => {
  const block = editableBlocks.value[index]
  
  // Destroy tabulator instance if exists
  if (tabulatorInstances.value[index]) {
    tabulatorInstances.value[index].destroy()
    delete tabulatorInstances.value[index]
  }
  
  block.editContent = block.content
  block.isEditing = false
}

// Save all changes
const saveAllChanges = () => {
// Emit the updated blocks to parent
const updatedBlocks = editableBlocks.value.map(block => ({
    type: block.type,
    content: block.content,
    page: block.page,
    confidence: block.confidence,
    bbox: block.bbox
}))

emit('update:blocks', updatedBlocks)

// Reset modified flags
editableBlocks.value.forEach(block => {
    block.originalContent = block.content
    block.isModified = false
})

alert('✅ All changes saved!')
}

// Confirm and delete a block
const confirmDelete = (index) => {
    const block = editableBlocks.value[index]
    const confirmed = confirm(`Are you sure you want to delete this ${block.type} block?\n\n"${block.content.substring(0, 100)}..."`)

    if (confirmed) {
    deleteBlock(index)
    }
}

// Delete a block
const deleteBlock = (index) => {
    editableBlocks.value.splice(index, 1)

    // Mark as having changes so user can save
    // Emit immediately or wait for save all
    emit('update:blocks', editableBlocks.value.map(block => ({
    type: block.type,
    content: block.content,
    page: block.page,
    confidence: block.confidence,
    bbox: block.bbox
    })))

    alert('🗑️ Block deleted!')
}

// Get color classes for block type badge
const getTypeColor = (type) => {
const colors = {
    'Title': 'bg-purple-100 text-purple-800',
    'Section Header': 'bg-blue-100 text-blue-800',
    'Header': 'bg-indigo-100 text-indigo-800',
    'Table': 'bg-green-100 text-green-800',
    'Text': 'bg-gray-100 text-gray-800',
    'Figure': 'bg-yellow-100 text-yellow-800',
    'Footer': 'bg-gray-100 text-gray-600',
    'Key Value': 'bg-blue-100 text-blue-800',
}
return colors[type] || 'bg-gray-100 text-gray-600'
}

// Get appropriate text styling for content type
const getTextClass = (type, isEditing) => {
if (isEditing) {
    return 'resize-y'
}

const classes = {
    'Title': 'text-xl font-bold text-gray-900',
    'Section Header': 'text-lg font-semibold text-gray-800',
    'Header': 'text-base font-medium text-gray-700',
    'Text': 'text-sm text-gray-600 whitespace-pre-wrap',
    'Figure': 'bg-gray-50 p-4 rounded text-sm text-gray-500',
    'Footer': 'text-xs text-gray-400 italic',
    'Key Value': 'bg-blue-50 p-3 rounded text-sm whitespace-pre-wrap',
}
return classes[type] || 'text-sm text-gray-600'
}
</script>

<style>
/* Style for Reducto HTML tables */
.reducto-table table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.875rem;
}

.reducto-table th {
    background-color: #f3f4f6;
    border: 1px solid #d1d5db;
    padding: 0.5rem 1rem;
    text-align: left;
    font-weight: 600;
}

.reducto-table td {
    border: 1px solid #d1d5db;
    padding: 0.5rem 1rem;
}

.reducto-table tr:hover {
    background-color: #f9fafb;
}

/* Tabulator table editor styles */
.tabulator {
  font-size: 0.875rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
}

.tabulator-cell {
  padding: 0.5rem 1rem;
}

.tabulator-editing {
  background-color: #fef3c7 !important;
}

.tabulator-header {
  background-color: #f3f4f6;
  font-weight: 600;
}
</style>