<template>
    <div class="bg-white shadow rounded-lg overflow-hidden">
        <!-- Header -->
        <div class="px-4 py-3 border-b border-gray-200 bg-gray-50">
        <div class="flex items-center justify-between">
            <h4 class="text-sm font-medium text-gray-900">Extracted Content</h4>
            <div class="flex items-center space-x-3">
            <span class="text-xs text-gray-500">
                {{ editableBlocks?.length || 0 }} block(s) found
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
                <span v-if="block.confidence" class="text-xs text-gray-500">
                    {{ block.confidence }} confidence
                </span>
                <span v-if="block.isModified" class="text-xs text-orange-600 font-medium">
                    ● Modified
                </span>
                </div>

                <!-- Edit Controls -->
                <div class="flex items-center space-x-2">
                <button
                    @click="confirmDelete(index)"
                    class="px-2 py-1 text-xs bg-red-100 text-red-700 rounded hover:bg-red-200"
                    title="Delete this block"
                >
                    <Icon name="heroicons:trash" class="w-3 h-3" />
                </button>

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

            <!-- Block Content -->
            <div class="mt-2">
                <!-- TABLE -->
                <div v-if="block.type === 'Table'">
                <div v-if="block.isEditing" class="space-y-2">
                    <div class="flex gap-2 mb-2">
                    <button class="btn-add-row" @click="addRow(index)">Add Row</button>
                    <button class="btn-add-col" @click="addColumn(index)">Add Column</button>
                    <button class="btn-remove-row" @click="removeRow(index)">Remove Row</button>
                    <button class="btn-remove-col" @click="removeColumn(index)">Remove Column</button>
                    </div>
                    <div :id="`table-editor-${index}`" class="border rounded"></div>
                    <p class="text-xs text-gray-500">
                    ✏️ Click any cell to edit/select. Row/col actions use the selected cell as a target if available.
                    </p>
                </div>
                <div v-else class="overflow-x-auto">
                    <div v-html="block.content" class="reducto-table"></div>
                </div>
                </div>

                <!-- TEXT / OTHERS -->
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
import { ref, computed, watch, nextTick, onBeforeUnmount } from 'vue'
import { TabulatorFull as Tabulator } from 'tabulator-tables'
import 'tabulator-tables/dist/css/tabulator.min.css'
import { useTableParser } from '../../composables/useTableParser'

const props = defineProps({
    allBlocks: { type: Array, default: () => [] },
    filename: { type: String, default: '' },
    currentUserEmail: { type: String, default: 'qc-user@example.com' },
})
const emit = defineEmits(['update:blocks'])

const { parseHTMLTableToJSON, convertJSONToHTML } = useTableParser()
const editableBlocks = ref([])
const originalBlocks = ref([])
const tabulatorInstances = new Map()

watch(
    () => props.allBlocks,
    (newBlocks) => {
        if (newBlocks && newBlocks.length > 0) {
        editableBlocks.value = newBlocks.map((block) => ({
            ...block,
            isEditing: false,
            isModified: false,
            editContent: block.content,
            originalContent: block.content,
        }))
        originalBlocks.value = JSON.parse(JSON.stringify(newBlocks))
        } else {
        editableBlocks.value = []
        originalBlocks.value = []
        }
    },
    { immediate: true }
)

const hasChanges = computed(() => editableBlocks.value.some((b) => b.isModified))

const startEditing = async (index) => {
    const block = editableBlocks.value[index]
    block.isEditing = true

    if (block.type === 'Table') {
        await nextTick()
        const mountId = `table-editor-${index}`
        const el = document.getElementById(mountId)
        if (!el) return
        const { columns, rows } = parseHTMLTableToJSON(block.content || '<table></table>')

        const table = new Tabulator(el, {
            data: rows,
            columns: columns.map((c) => ({ title: c.title, field: c.field, editor: 'input' })),
            layout: 'fitDataStretch',
            movableColumns: true,
            pagination: false,
            selectable: 1, // Enable single row selection
            // selectableRange: true, // Allow selecting cells
            columnDefaults: { headerSort: false },
            cellClick: (e, cell) => {
                // Store the active cell position
                const row = cell.getRow()
                const col = cell.getColumn()
                const rowPos = row.getPosition() // 1-indexed position
                const colField = col.getField()
                
                lastActiveCell.set(index, {
                rowPos: rowPos,
                colField: colField
                })
                
                // Select the row for visual feedback
                row.select()
            },
            cellEdited: (cell) => {
                // Also track when cell is edited (user might edit before clicking add/remove)
                const row = cell.getRow()
                const col = cell.getColumn()
                const rowPos = row.getPosition()
                const colField = col.getField()
                
                lastActiveCell.set(index, {
                rowPos: rowPos,
                colField: colField
                })
            },
        })
        tabulatorInstances.set(index, table)
    }
}

const saveBlock = async (index) => {
    const block = editableBlocks.value[index]
    if (block.type === 'Table') {
        const table = tabulatorInstances.get(index)
        if (!table) return
        const rows = await table.getData()
        const columns = table.getColumns().map((c) => ({
            title: c.getDefinition().title,
            field: c.getField(),
        }))
        const editedHTML = convertJSONToHTML(columns, rows)
        block.originalContent = block.originalContent || block.content
        block.content = editedHTML
        block.isEditing = false
        block.isModified = true

        try {
            await fetch('/api/v1/extraction/save-qc', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                    block_id: block.id || `block-${index}`,
                    page: block.page,
                    type: block.type,
                    filename: props.filename || 'unknown.pdf',
                    original_html: block.originalContent || '',
                    edited_html: editedHTML,
                    table_json: { columns, rows },
                    qc: { edited_by: props.currentUserEmail, edited_at: new Date().toISOString(), notes: '' },
                }),
            })
        } catch (e) {
            console.error('Save table failed', e)
        } finally {
            const t = tabulatorInstances.get(index)
            if (t) t.destroy()
            tabulatorInstances.delete(index)
        }
    } else {
        block.originalContent = block.originalContent || block.content
        block.content = block.editContent
        block.isEditing = false
        block.isModified = true
        try {
            await fetch('/api/v1/extraction/save-qc', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                block_id: block.id || `block-${index}`,
                page: block.page,
                type: block.type,
                filename: props.filename || 'unknown.pdf',
                original_html: block.originalContent || '',
                edited_html: block.content,
                table_json: null,
                qc: { edited_by: props.currentUserEmail, edited_at: new Date().toISOString(), notes: '' },
                }),
            })
        } catch (e) {
            console.error('Save text failed', e)
        }
    }
}

const cancelEditing = (index) => {
    const block = editableBlocks.value[index]
    const t = tabulatorInstances.get(index)
    if (t) t.destroy()
    tabulatorInstances.delete(index)
    lastActiveCell.delete(index)
    block.isEditing = false
    block.editContent = block.content
}

// Enhanced: Add/remove row/column at selection (or end if nothing selected)
const addRow = (idx) => {
    const table = tabulatorInstances.get(idx)
    if (!table) return
    // Try to use last active cell position first
    const activeCell = lastActiveCell.get(idx)
    if (activeCell && activeCell.rowPos) {
        // addRow(data, addToTop, position) - position is 1-indexed
        // Insert AFTER the active row (position + 1)
        table.addRow({}, false, activeCell.rowPos + 1)
        // Update the stored position to account for the new row
        lastActiveCell.set(idx, { ...activeCell, rowPos: activeCell.rowPos + 1 })
        return
    }
    // Fallback: try selected rows
    const selected = table.getSelectedRows()
    if (selected.length > 0) {
        const rowPos = selected[0].getPosition()
        table.addRow({}, false, rowPos + 1)
        return
    }
    // Last resort: add at end
    table.addRow({})
}


// Remove row at the last active cell's position
const removeRow = (idx) => {
    const table = tabulatorInstances.get(idx)
    if (!table) return

    // Try to use last active cell position first
    const activeCell = lastActiveCell.get(idx)
    if (activeCell && activeCell.rowPos) {
        const rowComponent = table.getRowFromPosition(activeCell.rowPos, true)
        if (rowComponent) {
            rowComponent.delete()
            // Update stored position if there are rows remaining
            const data = table.getData()
            if (data.length > 0 && activeCell.rowPos <= data.length) {
                // Row deleted, but if position was at end, adjust
                if (activeCell.rowPos > data.length) {
                    lastActiveCell.set(idx, { ...activeCell, rowPos: data.length })
                }
            } else {
                lastActiveCell.delete(idx)
            }
            return
        }
    }

    // Fallback: try selected rows
    const selected = table.getSelectedRows()
    if (selected.length > 0) {
        selected[0].delete()
        return
    }

    // Last resort: remove last row
    const data = table.getData()
    if (data.length > 0) {
        const rowComponent = table.getRowFromPosition(data.length, true)
        rowComponent?.delete()
    }
}

// Add column AFTER the last active cell's column
const addColumn = (idx) => {
    const table = tabulatorInstances.get(idx)
    if (!table) return

    const cols = table.getColumns()
    let insertAt = cols.length

    // Try to use last active cell column
    const activeCell = lastActiveCell.get(idx)
    if (activeCell && activeCell.colField) {
        const colIndex = cols.findIndex(c => c.getField() === activeCell.colField)
        if (colIndex >= 0) {
            insertAt = colIndex + 1 // Insert AFTER the active column
        }
    }

    const field = `col${cols.length}`
    table.addColumn(
        { title: `Column ${cols.length + 1}`, field, editor: 'input' },
        false,
        insertAt
    )
}

// Remove column at the last active cell's position
const removeColumn = (idx) => {
    const table = tabulatorInstances.get(idx)
    if (!table) return

    const cols = table.getColumns()
    let fieldToDelete = null

    // Try to use last active cell column
    const activeCell = lastActiveCell.get(idx)
    if (activeCell && activeCell.colField) {
        
    }

    if (!fieldToDelete && cols.length > 0) {
        fieldToDelete = cols[cols.length - 1].getField()
    }

    if (fieldToDelete) {
        table.deleteColumn(fieldToDelete)
        // Update stored column reference if needed
        if (activeCell && activeCell.colField === fieldToDelete) {
            lastActiveCell.delete(idx)
        }
    }
}

const saveAllChanges = () => {
    const updatedBlocks = editableBlocks.value.map((b) => ({
        type: b.type,
        content: b.content,
        page: b.page,
        confidence: b.confidence,
        bbox: b.bbox,
    }))
    emit('update:blocks', updatedBlocks)
    editableBlocks.value.forEach((b) => {
        b.originalContent = b.content
        b.isModified = false
    })
    alert('✅ All changes saved!')
}

const confirmDelete = (index) => {
    const block = editableBlocks.value[index]
    const ok = confirm(
        `Are you sure you want to delete this ${block.type} block?\n\n"${String(block.content).substring(0, 100)}..."`
    )
    if (ok) deleteBlock(index)
    }

const deleteBlock = (index) => {
    editableBlocks.value.splice(index, 1)
    emit(
        'update:blocks',
        editableBlocks.value.map((b) => ({
        type: b.type,
        content: b.content,
        page: b.page,
        confidence: b.confidence,
        bbox: b.bbox,
        }))
    )
    alert('🗑️ Block deleted!')
}

const getTypeColor = (type) => {
    const colors = {
        Title: 'bg-purple-100 text-purple-800',
        'Section Header': 'bg-blue-100 text-blue-800',
        Header: 'bg-indigo-100 text-indigo-800',
        Table: 'bg-green-100 text-green-800',
        Text: 'bg-gray-100 text-gray-800',
        Figure: 'bg-yellow-100 text-yellow-800',
        Footer: 'bg-gray-100 text-gray-600',
        'Key Value': 'bg-blue-100 text-blue-800',
    }
    return colors[type] || 'bg-gray-100 text-gray-600'
}

const getTextClass = (type, isEditing) => {
    if (isEditing) return 'resize-y'
    const classes = {
        Title: 'text-xl font-bold text-gray-900',
        'Section Header': 'text-lg font-semibold text-gray-800',
        Header: 'text-base font-medium text-gray-700',
        Text: 'text-sm text-gray-600 whitespace-pre-wrap',
        Figure: 'bg-gray-50 p-4 rounded text-sm text-gray-500',
        Footer: 'text-xs text-gray-400 italic',
        'Key Value': 'bg-blue-50 p-3 rounded text-sm whitespace-pre-wrap',
    }
    return classes[type] || 'text-sm text-gray-600'
}

onBeforeUnmount(() => {
    tabulatorInstances.forEach((t) => t.destroy())
    tabulatorInstances.clear()
})
</script>

<style>
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

.btn-add-row, .btn-add-col, .btn-remove-row, .btn-remove-col {
    padding: 4px 8px;
    border: 1px solid #ccc;
    margin-right: 4px;
    border-radius: 4px;
    background: #f4f4f4;
    font-size: .88em;
}
.btn-add-row, .btn-add-col { background: #e0fce0; }
.btn-remove-row, .btn-remove-col { background: #fee0e0; }
</style>