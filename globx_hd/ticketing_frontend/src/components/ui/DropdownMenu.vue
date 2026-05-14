<template>
  <div class="relative inline-block text-left" @keydown.escape="open=false">
    <slot name="button" :toggle="toggle" :open="open" />
    <transition enter-active-class="transition ease-out duration-100" enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75" leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
      <div v-if="open" class="absolute right-0 mt-2 w-40 origin-top-right rounded-lg bg-white shadow-xl border border-gray-200 focus:outline-none z-50 backdrop-blur-sm">
        <slot />
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref } from 'vue';
const open = ref(false);
function toggle() { open.value = !open.value; }
</script>
