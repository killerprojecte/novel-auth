<script lang="ts">
  import { roundClass } from "./util";

  let {
    type = "text",
    placeholder = "",
    round = "full",
    value = $bindable(),
    class: className = "",
  } = $props();
</script>

<input
  {type}
  {placeholder}
  bind:value
  autocomplete="off"
  class={[
    "box-border h-12 w-full border border-[#8d8d8d] px-4 text-base",
    "focus:border-primary invalid:border-red-600 focus:outline-none invalid:focus:border-red-600",
    roundClass(round),
    className,
  ]}
/>
