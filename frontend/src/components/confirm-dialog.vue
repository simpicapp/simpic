<template>
  <modal :closeable="false" :should-close="shouldClose" @close="onClose">
    <popup :closeable="false" :title="title" position="center">
      <p>{{ body }}</p>
      <div class="buttons">
        <button :class="{'danger-button': dangerous}" @click="onYes">{{ yesText }}</button>
        <button @click="onNo">{{ noText }}</button>
      </div>
    </popup>
  </modal>
</template>

<style lang="scss" scoped>
  .buttons {
    display: grid;
    grid-template-columns: auto auto;
    grid-column-gap: 20px;
    justify-items: stretch;
    margin-top: 30px;
  }
</style>

<script lang="ts">
  import Modal from "./modal.vue";
  import Popup from "./popup.vue";
  import Vue from "vue";

  export default Vue.extend({
    props: {
      body: String,
      dangerous: {
        type: Boolean,
        default: false,
      },
      noText: {
        type: String,
        default: "No",
      },
      title: String,
      yesText: {
        type: String,
        default: "Yes",
      },
    },
    data() {
      return {
        shouldClose: false,
      };
    },
    components: {
      Modal,
      Popup,
    },
    methods: {
      onClose() {
        this.$emit("close");
        this.shouldClose = false;
      },
      onNo() {
        this.$emit("no");
        this.shouldClose = true;
      },
      onYes() {
        this.$emit("yes");
        this.shouldClose = true;
      },
    },
  });
</script>
