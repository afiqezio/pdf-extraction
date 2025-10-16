<template>
  <!-- Right Side - Table Editor (60%) -->
  <div class="bg-white shadow rounded-lg overflow-hidden">
    <div class="px-4 py-3 border-b border-gray-200 bg-gray-50">
      <div class="flex items-center justify-between">
        <h4 class="text-sm font-medium text-gray-900">Table Editor</h4>
        <span class="text-xs text-gray-500">
          {{ allTables?.length || 0 }} table(s) found
        </span>
      </div>
    </div>
    
    <!-- Table Navigation Sidebar -->
    <div class="px-4 py-2 border-b border-gray-200 bg-gray-50">
      <div class="flex space-x-2 overflow-x-auto">
        <button
          v-for="(table, tableIndex) in allTables"
          :key="tableIndex"
          @click="selectTable(tableIndex)"
          :class="[
            'px-3 py-1 text-xs rounded-full border transition-colors',
            props.selectedTableIndex === tableIndex
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
                  Table {{ props.selectedTableIndex + 1 }}
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
    <div v-if="!allTables || allTables.length === 0" class="text-center py-8 border border-gray-200 rounded-lg">
      <Icon name="heroicons:table-cells" class="mx-auto h-12 w-12 text-gray-400" />
      <p class="mt-2 text-sm text-gray-500">No tables extracted</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// ===== PROPS =====
const props = defineProps({
  allTables: {
    type: Array,
    default: () => []
  },
  selectedTableIndex: {
    type: Number,
    default: 0
  }
})

// ===== EMITS =====
const emit = defineEmits([
  'table-selected',
  'table-updated', 
  'table-deleted'
])

// ===== COMPUTED =====
const selectedTable = computed(() => {
  if (!props.allTables || props.allTables.length === 0) {
    return null
  }
  return props.allTables[props.selectedTableIndex]
})

const allTables = computed(() => props.allTables)
const selectedTableIndex = computed(() => props.selectedTableIndex)

// ===== METHODS =====
const selectTable = (tableIndex) => {
  emit('table-selected', tableIndex)
}

const updateTableData = (tableIndex) => {
  emit('table-updated', tableIndex)
}

const deleteRow = (tableIndex, rowIndex) => {
  emit('table-deleted', { type: 'row', tableIndex, rowIndex })
}

const deleteColumn = (tableIndex, columnIndex) => {
  emit('table-deleted', { type: 'column', tableIndex, columnIndex })
}

const deleteHeaderRow = (tableIndex) => {
  emit('table-deleted', { type: 'header', tableIndex })
}

const deleteTable = (tableIndex) => {
  emit('table-deleted', { type: 'table', tableIndex })
}
</script>
