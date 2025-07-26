<script lang="ts">
  let { rules = undefined, children } = $props();

  let root: HTMLElement;
  let validateError = $state("");
  function validate() {
    const input = root.getElementsByTagName("input").item(0);
    if (!input) {
      validateError = "";
      return;
    }
    if (rules) {
      const result = rules(input.value);
      validateError = result === true ? "" : `* ${result}`;
    } else {
      validateError = "";
    }
    input.setCustomValidity(validateError);
  }
</script>

<div class="relative w-auto">
  <div class="flex" bind:this={root} oninput={validate} onblur={validate}>
    {@render children?.()}
  </div>
  <div class="text-error mt-1 text-left text-xs text-red-600">
    {validateError}
  </div>
</div>
