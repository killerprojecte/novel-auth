<script lang="ts">
  let { tabs, activeTab = $bindable() } = $props();
  let activeTabClass = "text-primary border-primary";
  let inactiveTabClass =
    "text-gray-500 hover:text-gray-700 border-transparent cursor-pointer";
</script>

<div class="flex justify-center gap-12 border-b border-gray-200">
  {#each tabs as tab}
    <button
      type="button"
      onclick={() => (activeTab = tab)}
      class={[
        "border-b-2 px-1 py-4 text-base font-medium whitespace-nowrap",
        activeTab === tab ? activeTabClass : inactiveTabClass,
      ]}
    >
      {tab}
    </button>
  {/each}
</div>
