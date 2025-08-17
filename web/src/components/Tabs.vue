<script setup lang="ts">
interface Props {
  tabs: string[];
}

defineProps<Props>();
const model = defineModel<string>();

const activeTabClass = 'text-primary border-primary';
const inactiveTabClass =
  'text-gray-500 hover:text-gray-700 border-transparent cursor-pointer';

const getTabClass = (tab: string) => {
  return [
    'border-b-2 pl-1 pr-1 pt-4 pb-4 text-base font-medium whitespace-nowrap',
    model.value === tab ? activeTabClass : inactiveTabClass,
  ];
};
</script>

<template>
  <div class="flex justify-center gap-12 border-b border-gray-200">
    <button
      v-for="tab in tabs"
      :key="tab"
      type="button"
      @click="model = tab"
      :class="getTabClass(tab)"
    >
      {{ tab }}
    </button>
  </div>
</template>
