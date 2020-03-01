<template>
  <div>
    <aside class="selectionbar" v-if="selectionCount > 0">
      {{ selectionCount }} selected
      <button @click="handleAddToAlbum">Add to album</button>
      <button @click="handleRemoveFromAlbum" v-if="!!album">Remove from album</button>
      <button @click="handleDelete">Delete</button>
      <button @click="$emit('clear-selection')">Clear selection</button>
    </aside>

    <DeleteDialog
      :what="`${selectionCount} photo${selectionCount === 1 ? '' : 's'}`"
      @close="showConfirmation = false"
      @yes="doDelete"
      v-if="showConfirmation"
    >
    </DeleteDialog>
  </div>
</template>

<style lang="scss" scoped>
  .selectionbar {
    position: fixed;
    z-index: 800;
    top: 0;
    left: 25%;
    right: 25%;
    border: 2px solid black;
    border-top: 0;
    padding: 25px;
    border-bottom-right-radius: 10px;
    border-bottom-left-radius: 10px;
    background: #ffffff;
    display: flex;
    justify-content: space-between;
  }

  .buttons {
    display: grid;
    grid-template-columns: auto auto;
    grid-column-gap: 20px;
    justify-items: stretch;
    margin-top: 30px;
  }
</style>

<script lang="ts">
  import {EventBus} from "./bus";
  import Axios from "axios";
  import DeleteDialog from "./delete-dialog.vue";
  import {computed, defineComponent, reactive, toRefs} from "@vue/composition-api";

  export default defineComponent({
    components: {
      DeleteDialog,
    },
    props: {
      album: String,
      selectionCount: Number,
      selection: Object,
    },
    setup(props, ctx) {
      const state = reactive({
        showConfirmation: false,
      });

      const selectionNoun = computed(() => props.selectionCount + " photo" + (props.selectionCount === 1 ? "" : "s"));

      function selectionKeys() {
        if (props.selection) {
          return Object.keys(props.selection);
        } else {
          return [];
        }
      }

      function doDelete() {
        if (!props.selection) {
          return;
        }

        Axios.post("/photos/delete", {photos: selectionKeys()}).then(() => {
          EventBus.$emit("toast", `${selectionNoun.value} deleted`);
          EventBus.$emit("refresh-gallery");
        });
      }

      function handleAddToAlbum() {
        new Promise((resolve, reject) => {
          EventBus.$emit("pick-album", resolve, reject);
        })
          .then(album =>
            Axios.post("/albums/" + album + "/photos", {
              // eslint-disable-next-line @typescript-eslint/camelcase
              add_photos: selectionKeys(),
            }).then(() => {
              EventBus.$emit("toast", `${selectionNoun.value} added to album`);
              EventBus.$emit("album-updated", album);
              ctx.emit("clear-selection");
            })
          )
          .catch(() => {
            /* ignore */
          });
      }

      function handleDelete() {
        state.showConfirmation = true;
      }

      function handleRemoveFromAlbum() {
        // eslint-disable-next-line @typescript-eslint/camelcase
        Axios.post("/albums/" + props.album + "/photos", {remove_photos: selectionKeys()}).then(() => {
          EventBus.$emit("toast", `${selectionNoun.value} removed from album`);
          EventBus.$emit("album-updated", props.album);
          ctx.emit("clear-selection");
        });
      }

      return {handleDelete, handleRemoveFromAlbum, handleAddToAlbum, doDelete, ...toRefs(state)};
    },
  });
</script>
