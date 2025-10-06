<template>
  <div :class="containerClasses">
    <div :class="spinnerClasses">
      <Icon name="heroicons:arrow-path" class="animate-spin" />
    </div>
    <p v-if="message" :class="messageClasses">
      {{ message }}
    </p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg', 'xl'].includes(value)
  },
  message: {
    type: String,
    default: ''
  },
  center: {
    type: Boolean,
    default: false
  },
  overlay: {
    type: Boolean,
    default: false
  }
})

const sizeConfig = computed(() => {
  const configs = {
    sm: {
      spinner: 'h-4 w-4',
      text: 'text-sm'
    },
    md: {
      spinner: 'h-6 w-6',
      text: 'text-base'
    },
    lg: {
      spinner: 'h-8 w-8',
      text: 'text-lg'
    },
    xl: {
      spinner: 'h-12 w-12',
      text: 'text-xl'
    }
  }
  return configs[props.size]
})

const containerClasses = computed(() => {
  const base = 'flex items-center'
  const direction = props.message ? 'flex-col space-y-2' : 'space-x-2'
  const center = props.center ? 'justify-center' : ''
  const overlay = props.overlay ? 'fixed inset-0 bg-white bg-opacity-75 z-50' : ''
  
  return `${base} ${direction} ${center} ${overlay}`
})

const spinnerClasses = computed(() => 
  `text-indigo-600 ${sizeConfig.value.spinner}`
)

const messageClasses = computed(() => 
  `text-gray-600 ${sizeConfig.value.text}`
)
</script>