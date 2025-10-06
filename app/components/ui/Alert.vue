<template>
  <div v-if="show" :class="alertClasses">
    <div class="flex">
      <div class="flex-shrink-0">
        <Icon :name="iconName" :class="iconClasses" />
      </div>
      <div class="ml-3 flex-1">
        <h3 v-if="title" :class="titleClasses">
          {{ title }}
        </h3>
        <div :class="messageClasses">
          <slot>
            {{ message }}
          </slot>
        </div>
      </div>
      <div v-if="dismissible" class="ml-auto pl-3">
        <div class="-mx-1.5 -my-1.5">
          <button
            type="button"
            :class="dismissClasses"
            @click="dismiss"
          >
            <span class="sr-only">Dismiss</span>
            <Icon name="heroicons:x-mark" class="h-5 w-5" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['success', 'error', 'warning', 'info'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    default: ''
  },
  dismissible: {
    type: Boolean,
    default: false
  },
  modelValue: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'dismiss'])

const show = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const typeConfig = computed(() => {
  const configs = {
    success: {
      bg: 'bg-green-50',
      border: 'border-green-200',
      icon: 'heroicons:check-circle',
      iconColor: 'text-green-400',
      titleColor: 'text-green-800',
      messageColor: 'text-green-700',
      dismissColor: 'text-green-500 hover:bg-green-100 focus:ring-green-600'
    },
    error: {
      bg: 'bg-red-50',
      border: 'border-red-200',
      icon: 'heroicons:x-circle',
      iconColor: 'text-red-400',
      titleColor: 'text-red-800',
      messageColor: 'text-red-700',
      dismissColor: 'text-red-500 hover:bg-red-100 focus:ring-red-600'
    },
    warning: {
      bg: 'bg-yellow-50',
      border: 'border-yellow-200',
      icon: 'heroicons:exclamation-triangle',
      iconColor: 'text-yellow-400',
      titleColor: 'text-yellow-800',
      messageColor: 'text-yellow-700',
      dismissColor: 'text-yellow-500 hover:bg-yellow-100 focus:ring-yellow-600'
    },
    info: {
      bg: 'bg-blue-50',
      border: 'border-blue-200',
      icon: 'heroicons:information-circle',
      iconColor: 'text-blue-400',
      titleColor: 'text-blue-800',
      messageColor: 'text-blue-700',
      dismissColor: 'text-blue-500 hover:bg-blue-100 focus:ring-blue-600'
    }
  }
  return configs[props.type]
})

const alertClasses = computed(() => 
  `rounded-md p-4 border ${typeConfig.value.bg} ${typeConfig.value.border}`
)

const iconName = computed(() => typeConfig.value.icon)

const iconClasses = computed(() => 
  `h-5 w-5 ${typeConfig.value.iconColor}`
)

const titleClasses = computed(() => 
  `text-sm font-medium ${typeConfig.value.titleColor}`
)

const messageClasses = computed(() => {
  const base = `text-sm ${typeConfig.value.messageColor}`
  return props.title ? `${base} mt-1` : base
})

const dismissClasses = computed(() => 
  `inline-flex rounded-md p-1.5 focus:outline-none focus:ring-2 focus:ring-offset-2 ${typeConfig.value.dismissColor}`
)

const dismiss = () => {
  show.value = false
  emit('dismiss')
}
</script>