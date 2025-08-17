<script setup lang="ts">
import { roundClass } from './util';

interface Props {
  type?: string;
  placeholder?: string;
  round?: string;
  class?: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  placeholder: '',
  round: 'full',
  class: '',
});
const model = defineModel<string>();

const inputClass = computed(() => [
  'box-border h-12 w-full border border-[#8d8d8d] pl-4 pr-4 text-base',
  'focus:border-primary invalid:border-red-600 focus:outline-none invalid:focus:border-red-600',
  roundClass(props.round),
  props.class,
]);
</script>

<template>
  <input
    v-model="model"
    :type="props.type"
    :placeholder="props.placeholder"
    autocomplete="off"
    :class="inputClass"
  />
</template>
