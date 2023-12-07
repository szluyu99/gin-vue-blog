<script setup>
import { reactive } from 'vue'

const props = defineProps({
  show: { type: Boolean, default: false },
  position: {
    type: String,
    default: 'top',
    validator: value => ['top', 'bottom'].includes(value),
  },
  align: {
    type: String,
    default: 'center',
    validator: value => ['left', 'center', 'right'].includes(value),
  },
  timeout: { type: Number, default: 2500 },
  queue: { type: Boolean, default: true },
  zIndex: { type: Number, default: 100 },
  closeable: { type: Boolean, default: false },
})

const flux = reactive({
  /** @type { Array<{ show, id, content, type }> } */
  events: [],
  success: content => flux.add('success', content),
  info: content => flux.add('info', content),
  warning: content => flux.add('warning', content),
  error: content => flux.add('error', content),
  /**
   * @param {'success' | 'info' | 'warning' | 'error'} type
   * @param {string} content
   */
  show: (type, content) => flux.add(type, content),
  add: (type, content) => {
    if (!props.queue)
      flux.events = []

    setTimeout(() => {
      const event = { id: Date.now(), content, type }
      flux.events.push(event)
      setTimeout(() => flux.remove(event), props.timeout)
    }, 100)
  },
  remove: (event) => {
    flux.events = flux.events.filter(e => e.id !== event.id)
  },
})

defineExpose({
  show: flux.show,
  success: flux.success,
  info: flux.info,
  warning: flux.warning,
  error: flux.error,
})
</script>

<template>
  <Teleport to="body">
    <div
      v-show="flux.events.length"
      class="pointer-events-none fixed w-4/5 sm:w-[400px]"
      :class="{
        'left-1/2 -translate-x-1/2': align === 'center',
        'left-16': align === 'left',
        'right-16': align === 'right',
        'top-4': position === 'top',
        'bottom-6 ': position === 'bottom',
      }"
      :style="{ zIndex }"
    >
      <TransitionGroup
        tag="ul"
        enter-active-class="transition ease-out duration-200"
        leave-active-class="transition ease-in duration-200 absolute w-full"
        :enter-class="position === 'bottom'
          ? 'transform translate-y-3 opacity-0'
          : 'transform -translate-y-3 opacity-0'"
        enter-to-class="transform translate-y-0 opacity-100"
        leave-class="transform translate-y-0 opacity-100"
        :leave-to-class="position === 'bottom'
          ? 'transform translate-y-1/4 opacity-0'
          : 'transform -translate-y-1/4 opacity-0'"
        move-class="ease-in-out duration-200"
        class="inline-block w-full"
      >
        <li
          v-for="event in flux.events" :key="event.id"
          :class="{
            'pb-2': position === 'bottom',
            'pt-2': position === 'top',
          }"
        >
          <slot :type="event.type" :content="event.content">
            <div class="pointer-events-auto w-full overflow-hidden rounded-lg bg-white ring-1 ring-black ring-opacity-5">
              <div class="flex justify-between px-4 py-3">
                <div class="flex items-center">
                  <div
                    class="mr-6 h-6 w-6"
                    :class="{
                      'i-mdi:check-circle text-green': event.type === 'success',
                      'i-mdi:information-outline text-blue': event.type === 'info',
                      'i-mdi:alert-outline text-yellow': event.type === 'warning',
                      'i-mdi:alert-circle-outline text-red': event.type === 'error',
                    }"
                  />
                  <div class="ml-1">
                    <slot>
                      <div> {{ event.content }} </div>
                    </slot>
                  </div>
                </div>
                <button
                  v-if="closeable"
                  class="i-mdi:close h-5 w-5 flex items-center justify-center rounded-full rounded-full p-1 font-bold text-gray-400 hover:text-gray-600"
                  @click="flux.remove(event)"
                />
              </div>
            </div>
          </slot>
        </li>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
