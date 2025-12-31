<template>
  <div class="bg-white shadow rounded-lg overflow-hidden">
    <div class="border-b border-gray-200">
      <nav class="-mb-px flex" aria-label="Tabs">
        <button
          v-for="tab in tabs"
          :key="tab.name"
          @click="currentTab = tab.name"
          :class="[
            currentTab === tab.name
              ? 'border-indigo-500 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
            'w-1/3 py-4 px-1 text-center border-b-2 font-medium text-sm'
          ]"
        >
          {{ tab.name }}
        </button>
      </nav>
    </div>

    <div class="p-4">
      <div v-if="currentTab === 'Interactive'">
        <UiContentViewer :all-blocks="allBlocks" />
      </div>
      <div v-else-if="currentTab === 'JSON'">
        <div class="relative">
             <button @click="copyJson" class="absolute top-2 right-2 text-xs bg-gray-100 hover:bg-gray-200 px-2 py-1 rounded text-gray-600 border">Copy JSON</button>
             <pre class="bg-gray-50 p-4 rounded text-xs overflow-auto max-h-[600px] font-mono border">{{ jsonString }}</pre>
        </div>
      </div>
       <div v-else-if="currentTab === 'Markdown'">
        <div class="relative">
            <button @click="copyMarkdown" class="absolute top-2 right-2 text-xs bg-gray-100 hover:bg-gray-200 px-2 py-1 rounded text-gray-600 border">Copy Markdown</button>
            <div class="prose max-w-none p-4 bg-white rounded border max-h-[600px] overflow-auto">
                <pre class="whitespace-pre-wrap font-mono text-sm">{{ markdownString }}</pre>
            </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  data: {
    type: Array, // The job result data (array of pages)
    required: true
  }
})

const currentTab = ref('Interactive')
const tabs = [
  { name: 'Interactive' },
  { name: 'JSON' },
  { name: 'Markdown' }
]

const allBlocks = computed(() => {
    // Flatten pages to blocks for ContentViewer
    if (Array.isArray(props.data)) {
         return props.data.flatMap(page => extractBlocksFromPage(page))
    }
    return []
})

const jsonString = computed(() => JSON.stringify(props.data, null, 2))

const markdownString = computed(() => {
    return allBlocks.value.map(b => {
        if (b.type === 'Table') return `\n\n### Table (Page ${b.page})\n${b.content}\n`
        if (b.type === 'Title') return `\n# ${b.content}\n`
        if (b.type === 'Section Header') return `\n## ${b.content}\n`
        if (b.type === 'Header') return `\n### ${b.content}\n`
        return `\n${b.content}\n`
    }).join('')
})

function extractBlocksFromPage(page) {
    let blocks = []
    // Support both direct block list (if already processed) or Reducto structure
    if (page.blocks && Array.isArray(page.blocks)) {
        return page.blocks;
    }

    // Reducto result structure
    if (page.result) {
        let chunks = []
        if (page.result.parse?.result?.chunks) chunks = page.result.parse.result.chunks
        else if (page.result.chunks) chunks = page.result.chunks
        
        chunks.forEach(chunk => {
            if (chunk.blocks) {
                blocks.push(...chunk.blocks.map(b => ({
                    ...b, 
                    page: page.page_number || page.page || 1,
                    // ContentViewer expects specific fields
                    isEditing: false,
                    isModified: false
                })))
            }
        })
    } else if (page.chunks) {
        // Flat chunks structure
         page.chunks.forEach(chunk => {
            if (chunk.blocks) {
                blocks.push(...chunk.blocks.map(b => ({
                    ...b, 
                    page: page.page_number || 1
                })))
            }
        })
    }
    
    return blocks
}

function copyJson() {
    navigator.clipboard.writeText(jsonString.value)
    alert('JSON copied to clipboard')
}

function copyMarkdown() {
    navigator.clipboard.writeText(markdownString.value)
    alert('Markdown copied to clipboard')
}
</script>

