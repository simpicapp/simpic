<template>
  <gallery endpoint="/timeline"></gallery>
</template>

<script lang="ts">
  import {EventBus} from "../components/bus";
  import Gallery from "../components/gallery.vue";
  import {defineComponent} from "@vue/composition-api";
  import {useEventListener} from "@/features/eventbus";

  export default defineComponent({
    components: {
      Gallery,
    },
    setup() {
      useEventListener("upload-complete", () => EventBus.$emit("refresh-gallery"));
    },
  });
</script>
