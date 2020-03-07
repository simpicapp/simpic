<template>
  <div>
    <aside class="selectionbar" v-if="selectionCount > 0">
      <span>{{ selectionNoun }} selected</span>
      <Icon @click="handleAddToAlbum" name="folder-plus" title="Add to album"></Icon>
      <Icon @click="handleRemoveFromAlbum" name="folder-minus" title="Remove from this album" v-if="!!album"></Icon>
      &middot;
      <Icon @click="handleVisibility" name="eye-slash" title="Change visibility"></Icon>
      <Icon @click="handleDelete" name="trash-alt" title="Delete"></Icon>
      &middot;
      <Icon @click="$emit('clear-selection')" name="times" title="Clear selection"></Icon>
    </aside>

    <DeleteDialog
      :what="`${selectionCount} photo${selectionCount === 1 ? '' : 's'}`"
      @close="showConfirmation = false"
      @yes="doDelete"
      v-if="showConfirmation"
    >
    </DeleteDialog>

    <modal
      :closeable="true"
      :should-close="shouldVisibilityClose"
      @close="showVisibility = false"
      v-if="showVisibility"
    >
      <popup :closeable="true" position="center" title="Update photos" @close="showVisibility = false">
        <form class="update-photos" @submit.prevent="changeVisibility">
          <p>Change visibility of {{ selectionNoun }} to:</p>
          <VisibilitySwitch :visibility="newVisibility" @change="n => (newVisibility = n)"></VisibilitySwitch>
          <input type="submit" value="Change" :disabled="newVisibility < 0" />
        </form>
      </popup>
    </modal>
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

  .update-photos {
    display: grid;
    grid-template-columns: 100%;
    grid-row-gap: 20px;
  }
</style>

<script lang="ts">
  import Axios from "axios";
  import DeleteDialog from "./delete-dialog.vue";
  import Popup from "./popup.vue";
  import Modal from "./modal.vue";
  import VisibilitySwitch from "./visibility-switch.vue";
  import {computed, defineComponent, reactive, toRefs} from "@vue/composition-api";
  import "vue-awesome/icons/folder-plus";
  import "vue-awesome/icons/folder-minus";
  import "vue-awesome/icons/times";
  import "vue-awesome/icons/trash-alt";
  import "vue-awesome/icons/eye-slash";
  import Icon from "vue-awesome/components/Icon.vue";
  import {EventBus} from "@/features/eventbus";

  export default defineComponent({
    components: {
      DeleteDialog,
      Icon,
      Popup,
      Modal,
      VisibilitySwitch,
    },
    props: {
      album: String,
      selectionCount: Number,
      selection: Object,
    },
    setup(props, ctx) {
      const state = reactive({
        showConfirmation: false,
        showVisibility: false,
        shouldVisibilityClose: false,
        newVisibility: -1,
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

        Axios.post("/api/photos/delete", {photos: selectionKeys()}).then(() => {
          EventBus.$emit("toast", `${selectionNoun.value} deleted`);
          EventBus.$emit("refresh-gallery");
        });
      }

      function handleAddToAlbum() {
        new Promise((resolve, reject) => {
          EventBus.$emit("pick-album", resolve, reject);
        })
          .then(album =>
            Axios.post("/api/albums/" + album + "/photos", {
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

      function handleVisibility() {
        state.newVisibility = -1;
        state.shouldVisibilityClose = false;
        state.showVisibility = true;
      }

      function changeVisibility() {
        Axios.post("/api/photos/update", {photos: selectionKeys(), visibility: state.newVisibility}).then(() => {
          state.shouldVisibilityClose = true;
          EventBus.$emit("toast", `${selectionNoun.value} updated`);
          ctx.emit("clear-selection");
        });
      }

      function handleDelete() {
        state.showConfirmation = true;
      }

      function handleRemoveFromAlbum() {
        // eslint-disable-next-line @typescript-eslint/camelcase
        Axios.post("/api/albums/" + props.album + "/photos", {remove_photos: selectionKeys()}).then(() => {
          EventBus.$emit("toast", `${selectionNoun.value} removed from album`);
          EventBus.$emit("album-updated", props.album);
          ctx.emit("clear-selection");
        });
      }

      return {
        handleVisibility,
        handleDelete,
        handleRemoveFromAlbum,
        handleAddToAlbum,
        doDelete,
        changeVisibility,
        selectionNoun,
        ...toRefs(state),
      };
    },
  });
</script>
