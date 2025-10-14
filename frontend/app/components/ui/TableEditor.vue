<template>
  <div class="col-span-7 bg-white shadow rounded-lg overflow-hidden">
    <div class="px-4 py-3 border-b border-gray-200 bg-gray-50">
      <div class="flex items-center justify-between">
        <h4 class="text-sm font-medium text-gray-900">Table Editor</h4>
        <div class="flex items-center space-x-2">
          <span class="text-xs text-gray-500">
            Table {{ selectedTableIndex + 1 }} of {{ extractionResult?.allTables?.length || 0 }}
          </span>
          <UiButton
            variant="primary"
            size="sm"
            :loading="isSaving"
            :disabled="!extractionResult?.allTables?.length"
            @click="handleSaveToDatabase"
          >
            {{ isSaving ? 'Saving...' : 'Save to Database' }}
          </UiButton>
        </div>
      </div>
    </div>

    <!-- Table Navigation Sidebar -->
    <div v-if="extractionResult?.allTables?.length > 1" class="px-4 py-2 border-b border-gray-200 bg-gray-50">
      <div class="flex space-x-2 overflow-x-auto">
        <button
          v-for="(table, index) in extractionResult.allTables"
          :key="index"
          @click="selectTable(index)"
          :class="[
            'px-3 py-1 text-xs rounded-full whitespace-nowrap transition-colors',
            selectedTableIndex === index
              ? 'bg-indigo-100 text-indigo-700 border border-indigo-200'
              : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
          ]"
        >
          Table {{ index + 1 }}
          <span class="ml-1 text-gray-400">({{ table.rows?.length || 0 }} rows)</span>
        </button>
      </div>
    </div>

    <!-- Selected Table Content -->
    <div class="overflow-auto">
      <div v-if="selectedTable" class="p-4">
        <!-- Table Actions -->
        <div class="mb-4 flex justify-between items-center">
          <div class="flex space-x-2">
            <UiButton
              variant="danger"
              size="sm"
              @click="deleteTable(selectedTableIndex)"
            >
              Delete Table
            </UiButton>
            <UiButton
              variant="warning"
              size="sm"
              @click="deleteHeaderRow(selectedTableIndex)"
            >
              Delete Header Row
            </UiButton>
          </div>
          <div class="text-sm text-gray-500">
            {{ selectedTable.rows?.length || 0 }} rows × {{ selectedTable.headers?.length || 0 }} columns
          </div>
        </div>

        <!-- Table Display -->
        <div class="overflow-x-auto border border-gray-200 rounded-lg">
          <table class="min-w-full divide-y divide-gray-200">
            <!-- Header Row -->
            <thead class="bg-gray-50 sticky top-0 z-10">
              <tr>
                <th class="sticky left-0 bg-gray-50 px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider border-r border-gray-200">
                  #
                </th>
                <th
                  v-for="(header, headerIndex) in selectedTable.headers"
                  :key="headerIndex"
                  class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[120px]"
                >
                  <div class="flex items-center justify-between">
                    <input
                      v-model="selectedTable.headers[headerIndex]"
                      class="w-full bg-transparent border-none outline-none text-xs font-medium text-gray-900"
                      @blur="updateTableData(selectedTableIndex)"
                    />
                    <button
                      @click="deleteColumn(selectedTableIndex, headerIndex)"
                      class="ml-2 text-gray-400 hover:text-red-500"
                    >
                      <Icon name="heroicons:x-mark" class="h-3 w-3" />
                    </button>
                  </div>
                </th>
              </tr>
            </thead>

            <!-- Data Rows -->
            <tbody class="bg-white divide-y divide-gray-200">
              <tr
                v-for="(row, rowIndex) in selectedTable.rows"
                :key="rowIndex"
                class="hover:bg-gray-50"
              >
                <td class="sticky left-0 bg-white px-3 py-2 text-sm text-gray-500 border-r border-gray-200">
                  {{ rowIndex + 1 }}
                </td>
                <td
                  v-for="(cell, cellIndex) in row"
                  :key="cellIndex"
                  class="px-3 py-2 text-sm text-gray-900 min-w-[120px]"
                >
                  <input
                    v-model="selectedTable.rows[rowIndex][cellIndex]"
                    class="w-full bg-transparent border-none outline-none text-sm"
                    @blur="updateTableData(selectedTableIndex)"
                  />
                </td>
                <td class="px-3 py-2 text-sm text-gray-500">
                  <button
                    @click="deleteRow(selectedTableIndex, rowIndex)"
                    class="text-gray-400 hover:text-red-500"
                  >
                    <Icon name="heroicons:trash" class="h-4 w-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="flex items-center justify-center h-64 text-gray-500">
        <div class="text-center">
          <Icon name="heroicons:table-cells" class="mx-auto h-12 w-12 text-gray-400" />
          <p class="mt-2 text-sm">No table selected</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useTableEditing, useDatabase } from '~/composables/useExtractionApp'

const props = defineProps({
  extractionResult: Object
})

const { 
  selectedTableIndex, 
  selectedTable, 
  selectTable, 
  updateTableData, 
  deleteRow, 
  deleteColumn, 
  deleteHeaderRow, 
  deleteTable 
} = useTableEditing(toRef(props, 'extractionResult'))

const { isSaving, saveToDatabase } = useDatabase()

const handleSaveToDatabase = () => {
  if (props.extractionResult?.allTables) {
    saveToDatabase(props.extractionResult.allTables)
  }
}
</script>
