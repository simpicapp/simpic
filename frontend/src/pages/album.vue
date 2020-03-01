<template>
  <div>
    <div class="toolbar">
      <h2>{{ name }}</h2>
      <ActionIcon :working="deleting" @click="confirmDelete" name="trash-alt" v-if="loggedIn"></ActionIcon>
    </div>
    <gallery :album="id" :endpoint="'/albums/' + id + '/photos'"></gallery>

    <DeleteDialog @close="showConfirmation = false" @yes="doDelete" v-if="showConfirmation" what="this album">
    </DeleteDialog>
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
    grid-column-gap: 20px;
    justify-items: stretch;
    margin-top: 30px;
  }
</style>

<script lang="ts">
  import "vue-awesome/icons/trash-alt";
  import Axios from "axios";
  import Gallery from "../components/gallery.vue";
  import ActionIcon from "../components/action-icon.vue";
  import DeleteDialog from "../components/delete-dialog.vue";
  import {defineComponent, onMounted, reactive, toRefs} from "@vue/composition-api";
  import {useRouter} from "@/features/router";
  import {useAuthentication} from "@/features/auth";
  import {EventBus, useEventListener} from "@/features/eventbus";

  export default defineComponent({
    components: {
      ActionIcon,
      DeleteDialog,
      Gallery,
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

      onMounted(() => {
        Axios.get("/albums/" + props.id).then(({data: {name}}) => {
          state.name = name;
        });
      });

      return {confirmDelete, doDelete, loggedIn, ...toRefs(state)};
    },
  });
</script>
