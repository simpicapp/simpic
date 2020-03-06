<template>
  <gallery endpoint="/api/timeline"></gallery>
</template>

<script lang="ts">
  import Gallery from "../components/gallery.vue";
  import {defineComponent, onMounted} from "@vue/composition-api";
  import {EventBus, useEventListener} from "@/features/eventbus";
  import {useTitle} from "@/features/title";

  export default defineComponent({
    components: {
      Gallery,
    },
    setup() {
      const {setTitle} = useTitle();
      useEventListener("upload-complete", () => EventBus.$emit("refresh-gallery"));
      onMounted(() => setTitle(1, "Timeline"));
    },
  });
</script>
