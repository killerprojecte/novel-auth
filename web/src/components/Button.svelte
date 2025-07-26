<script lang="ts">
  import { roundClass } from "./util";

  let {
    text,
    round = "full",
    class: className = "",
    loading = false,
    ...rest
  } = $props();
</script>

<button
  {...rest}
  class={[
    "bg-primary h-12 w-full cursor-pointer border-0",
    "text-base font-medium text-white transition-all duration-200",
    "ease-in-out outline-none active:bg-[#0c7a43]",
    { loading: loading },
    roundClass(round),
    className,
  ]}>{text}</button
>

<style>
  .loading {
    position: relative;
    color: transparent !important;

    &:after {
      content: "";
      position: absolute;
      width: 10px;
      height: 10px;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      border-radius: 50%;
      animation: loading-dot-flashing 1.0s infinite ease alternate;
    }
  }

  @keyframes loading-dot-flashing {
    0% {
      box-shadow:
        16px 0 rgba(255, 255, 255, 0.3),
        -16px 0 rgba(255, 255, 255, 0.3);
      background: rgba(255, 255, 255, 0.3);
    }

    100% {
      box-shadow:
        16px 0 rgba(255, 255, 255, 0.7),
        -16px 0 rgba(255, 255, 255, 0.7);
      background: rgba(255, 255, 255, 0.7);
    }
  }
</style>
