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
  import {defineComponent, ref} from "@vue/composition-api";

  export default defineComponent({
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
    components: {
      Modal,
      Popup,
    },
    setup(_, ctx) {
      const shouldClose = ref(false);

      function emitAndClose(eventName: string, close: boolean) {
        return () => {
          ctx.emit(eventName);
          shouldClose.value = close;
        };
      }

      const onClose = emitAndClose("close", false);
      const onYes = emitAndClose("yes", true);
      const onNo = emitAndClose("no", true);

      return {shouldClose, onClose, onYes, onNo};
    },
  });
</script>
