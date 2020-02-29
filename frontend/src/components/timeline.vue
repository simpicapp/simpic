<template>
  <gallery endpoint="/timeline"></gallery>
</template>

<script lang="ts">
  import {EventBus} from "./bus";
  import Gallery from "./gallery.vue";
  import Vue from "vue";

  export default Vue.extend({
    components: {
      Gallery,
    },
    methods: {
      refresh() {
        EventBus.$emit("refresh-gallery");
      },
    },
    beforeDestroy() {
      EventBus.$off("upload-complete", this.refresh);
    },
    mounted() {
      EventBus.$on("upload-complete", this.refresh);
    },
  });
</script>
