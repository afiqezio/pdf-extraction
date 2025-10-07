<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="buttonClasses"
    @click="handleClick"
  >
    <Icon 
      v-if="loading" 
      name="heroicons:arrow-path" 
      class="animate-spin -ml-1 mr-2 h-4 w-4" 
    />
    <Icon 
      v-else-if="icon" 
      :name="icon" 
      class="-ml-1 mr-2 h-4 w-4" 
    />
    <slot />
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'button',
    validator: (value) => ['button', 'submit', 'reset'].includes(value)
  },
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'outline', 'ghost', 'danger'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl'].includes(value)
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  icon: {
    type: String,
    default: ''
  },
  block: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

const buttonClasses = computed(() => {
  const base = 'inline-flex items-center justify-center font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors duration-200'
  
  const sizes = {
    xs: 'px-2.5 py-1.5 text-xs',
    sm: 'px-3 py-2 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-4 py-2 text-base',
    xl: 'px-6 py-3 text-base'
  }
  
  const variants = {
    primary: 'text-white bg-indigo-600 hover:bg-indigo-700 focus:ring-indigo-500 disabled:bg-indigo-300',
    secondary: 'text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:ring-indigo-500 disabled:bg-indigo-50',
    outline: 'text-indigo-700 bg-white border border-indigo-300 hover:bg-indigo-50 focus:ring-indigo-500 disabled:bg-gray-50',
    ghost: 'text-gray-700 bg-transparent hover:bg-gray-100 focus:ring-gray-500 disabled:bg-transparent',
    danger: 'text-white bg-red-600 hover:bg-red-700 focus:ring-red-500 disabled:bg-red-300'
  }
  
  const width = props.block ? 'w-full' : ''
  const disabled = (props.disabled || props.loading) ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'
  
  return `${base} ${sizes[props.size]} ${variants[props.variant]} ${width} ${disabled}`
})

const handleClick = (event) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>