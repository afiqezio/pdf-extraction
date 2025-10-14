<template>
  <!-- Left Side - PDF Viewer (40%) -->
  <div class="bg-white shadow rounded-lg overflow-hidden">
    <div class="px-4 py-3 border-b border-gray-200 bg-gray-50">
      <div class="flex items-center justify-between">
        <h4 class="text-sm font-medium text-gray-900">PDF Viewer</h4>
        <div class="flex items-center space-x-2">
          <span class="text-xs text-gray-500">
            Page {{ currentPage || 1 }}
          </span>
          <div class="flex space-x-1">
            <button 
              @click="zoomOut"
              class="p-1 text-gray-400 hover:text-gray-600"
              title="Zoom Out"
            >
              <Icon name="heroicons:minus" class="h-4 w-4" />
            </button>
            <span class="text-xs text-gray-500 px-1">{{ Math.round(zoomLevel * 100) }}%</span>
            <button 
              @click="zoomIn"
              class="p-1 text-gray-400 hover:text-gray-600"
              title="Zoom In"
            >
              <Icon name="heroicons:plus" class="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
    <div class="h-full overflow-auto bg-gray-100">
      <div 
        v-if="pdfUrl" 
        class="pdf-container"
        :style="{ transform: `scale(${zoomLevel})`, transformOrigin: 'top left' }"
      >
        <iframe
          :key="pdfUrl"
          :src="pdfUrl"
          class="w-full h-full border-0"
          style="min-height: 800px;"
        />
      </div>
      <div v-else class="flex items-center justify-center h-full text-gray-500">
        <div class="text-center">
          <Icon name="heroicons:document-text" class="mx-auto h-12 w-12 text-gray-400" />
          <p class="mt-2 text-sm">No PDF available</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// ===== PROPS =====
const props = defineProps({
  pdfUrl: {
    type: String,
    default: null
  },
  currentPage: {
    type: Number,
    default: 1
  }
})

// ===== EMITS =====
const emit = defineEmits([
  'zoom-in',
  'zoom-out'
])

// ===== REACTIVE STATE =====
const zoomLevel = ref(1.0)

// ===== METHODS =====
const zoomIn = () => {
  if (zoomLevel.value < 2.0) {
    zoomLevel.value = Math.min(zoomLevel.value + 0.25, 2.0)
  }
  emit('zoom-in', zoomLevel.value)
}

const zoomOut = () => {
  if (zoomLevel.value > 0.5) {
    zoomLevel.value = Math.max(zoomLevel.value - 0.25, 0.5)
  }
  emit('zoom-out', zoomLevel.value)
}
</script>