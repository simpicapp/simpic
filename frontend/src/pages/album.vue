<template>
  <div>
    <div class="toolbar">
      <h2>{{ name }}</h2>
      <div class="buttons">
        <ActionIcon
          :working="false"
          @click="handleEditClicked"
          name="edit"
          title="Change this album's name or visibility"
          v-if="loggedIn"
        ></ActionIcon>
        <ActionIcon
          :working="deleting"
          @click="confirmDelete"
          name="trash-alt"
          title="Delete this album"
          v-if="loggedIn"
        ></ActionIcon>
      </div>
    </div>
    <gallery :album="id" :endpoint="'/albums/' + id + '/photos'"></gallery>

    <DeleteDialog @close="showConfirmation = false" @yes="doDelete" v-if="showConfirmation" what="this album">
    </DeleteDialog>

    <modal :closeable="true" :should-close="editShouldClose" @close="showEdit = false" v-if="showEdit">
      <AlbumDialog
        :id="id"
        :initialName="name"
        :initialVisibility="visibility"
        @close="editShouldClose = true"
        @created="handleEditFinished"
      ></AlbumDialog>
    </modal>
  </div>
</template>

<style lang="scss" scoped>
  h2 {
    margin: 0 20px;
  }

  .toolbar {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
  }

  .buttons {
    display: grid;
    grid-template-columns: auto auto;
    grid-column-gap: 10px;
    justify-items: stretch;
  }
</style>

<script lang="ts">
  import "vue-awesome/icons/trash-alt";
  import "vue-awesome/icons/edit";
  import Axios from "axios";
  import Gallery from "../components/gallery.vue";
  import ActionIcon from "../components/action-icon.vue";
  import AlbumDialog from "../components/album-dialog.vue";
  import DeleteDialog from "../components/delete-dialog.vue";
  import Modal from "../components/modal.vue";
  import {defineComponent, onMounted, reactive, toRefs} from "@vue/composition-api";
  import {useRouter} from "@/features/router";
  import {useAuthentication} from "@/features/auth";
  import {EventBus, useEventListener} from "@/features/eventbus";

  export default defineComponent({
    components: {
      ActionIcon,
      DeleteDialog,
      Gallery,
      Modal,
      AlbumDialog,
    },
    props: {
      id: String,
    },
    setup(props) {
      const {router} = useRouter();
      const {loggedIn} = useAuthentication();

      useEventListener("album-updated", (album: string) => {
        if (album === props.id) {
          EventBus.$emit("refresh-gallery");
        }
      });

      const state = reactive({
        name: "",
        showConfirmation: false,
        deleting: false,
        editShouldClose: false,
        showEdit: false,
        visibility: 0,
      });

      function confirmDelete() {
        state.showConfirmation = true;
      }

      function doDelete() {
        state.deleting = true;
        Axios.delete("/albums/" + props.id).then(() => {
          EventBus.$emit("albums-updated");
          EventBus.$emit("toast", "Album deleted");
          state.deleting = false;
          router.replace("/albums");
        });
      }

      function update() {
        Axios.get("/albums/" + props.id).then(({data: {name, visibility}}) => {
          state.name = name;
          state.visibility = visibility;
        });
      }

      onMounted(() => {
        update();
      });

      function handleEditClicked() {
        state.editShouldClose = false;
        state.showEdit = true;
      }

      function handleEditFinished() {
        EventBus.$emit("toast", "Album has been updated");
        state.editShouldClose = true;
        update();
      }

      return {confirmDelete, doDelete, handleEditClicked, handleEditFinished, loggedIn, ...toRefs(state)};
    },
  });
</script>
