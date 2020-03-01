<template>
  <div>
    <aside class="selectionbar" v-if="selectionCount > 0">
      <span>{{ selectionNoun }} selected</span>
      <Icon name="folder-plus" @click="handleAddToAlbum" title="Add to album"></Icon>
      <Icon name="folder-minus" @click="handleRemoveFromAlbum" v-if="!!album" title="Remove from this album"></Icon>
      &middot;
      <Icon name="trash-alt" @click="handleDelete" title="Delete"></Icon>
      &middot;
      <Icon name="times" @click="$emit('clear-selection')" title="Clear selection"></Icon>
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
  @use '../assets/css/vars';

  .selectionbar {
    position: fixed;
    z-index: 800;
    top: 20px;
    right: 20px;
    border: 1px solid vars.$primary;
    box-shadow: 0 0 10px black;
    padding: 10px 25px;
    border-radius: 10px;
    background: #ffffff;
    display: flex;
    justify-content: space-between;
    align-items: center;

    span {
      margin-right: 20px;
    }

    svg {
      box-sizing: content-box;
      width: 24px;
      height: 24px;
      padding: 10px;
      cursor: pointer;

      &:hover {
        color: vars.$primary;
      }
    }
  }
</style>

<script lang="ts">
  import Axios from "axios";
  import DeleteDialog from "./delete-dialog.vue";
  import {computed, defineComponent, reactive, toRefs} from "@vue/composition-api";
  import "vue-awesome/icons/folder-plus";
  import "vue-awesome/icons/folder-minus";
  import "vue-awesome/icons/times";
  import "vue-awesome/icons/trash-alt";
  import Icon from "vue-awesome/components/Icon.vue";
  import {EventBus} from "@/features/eventbus";

  export default defineComponent({
    components: {
      DeleteDialog,
      Icon,
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

      return {handleDelete, handleRemoveFromAlbum, handleAddToAlbum, doDelete, selectionNoun, ...toRefs(state)};
    },
  });
</script>
