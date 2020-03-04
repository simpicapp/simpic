<template>
  <modal :should-close="close" @close="visible = false" v-if="visible">
    <popup @close="close = true" position="center" title="Select an Album" v-if="selecting">
      <div class="album-picker">
        <template v-for="album in albums">
          <img
            :alt="album.name"
            :key="album.id"
            :src="'/data/thumb/' + album.cover_photo"
            @click="handleAlbumSelected(album.id)"
            v-if="album.cover_photo"
          />
          <span :key="album.id" v-else></span>
          <div :key="album.id + '.name'" @click="handleAlbumSelected(album.id)">
            <span>{{ album.name }}</span>
          </div>
        </template>
        <div @click="handleNewAlbumSelected" class="icon"><span>⊕</span></div>
        <div @click="handleNewAlbumSelected"><span>Create new album...</span></div>
      </div>
    </popup>
    <album-dialog @close="handleClosed" @created="handleAlbumSelected" v-else></album-dialog>
  </modal>
</template>

<style lang="scss" scoped>
  .album-picker {
    display: grid;
    grid-template-columns: 50px 200px;
    grid-auto-rows: 2em;
    grid-gap: 20px 20px;
    align-items: center;

    div {
      align-self: stretch;
      display: flex;
      align-items: center;
      cursor: pointer;

      &.icon {
        justify-content: center;
      }
    }
  }

  img {
    max-height: 2em;
    max-width: 50px;
    margin-right: 10px;
    border: 1px solid black;
    overflow: hidden;
    cursor: pointer;
  }

  form {
    display: grid;
    grid-template-columns: auto auto;
    grid-gap: 30px 20px;
    align-items: center;
  }

  input[type="submit"] {
    grid-column: span 2;
  }

  .alert {
    margin: 0;
    padding: 5px 10px;
    grid-column: span 2;
    background-color: darkred;
    color: white;
    font-weight: bold;
    text-align: center;
    border-radius: 15px;
    white-space: pre-line;
  }
</style>

<script lang="ts">
  import Axios from "axios";
  import Modal from "./modal.vue";
  import Popup from "./popup.vue";
  import AlbumDialog from "./album-dialog.vue";
  import {defineComponent, reactive, toRefs} from "@vue/composition-api";
  import {useEventListener} from "@/features/eventbus";

  interface Album {
    id: string;
    name: string;
    cover_photo: string;
    owner_id: number;
    created: string;
    photos: number;
  }

  export default defineComponent({
    components: {
      AlbumDialog,
      Modal,
      Popup,
    },
    setup() {
      const state = reactive({
        albums: new Array<Album>(),
        close: false,
        selecting: false,
        visible: false,
      });

      const promise = {
        reject() {
          // noop
        },

        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        resolve(_: string) {
          // noop
        },
      };

      useEventListener("pick-album", (resolve: (_: string) => void, reject: () => void) => {
        state.albums = [];
        state.close = false;
        state.selecting = true;
        state.visible = true;

        promise.resolve = resolve;
        promise.reject = reject;

        Axios.get("/albums").then(({data}) => {
          state.albums = data;
        });
      });

      function handleAlbumSelected(albumId: string) {
        state.visible = false;
        promise.resolve(albumId);
      }

      function handleClosed() {
        state.visible = false;
        promise.reject();
      }

      function handleNewAlbumSelected() {
        state.selecting = false;
      }

      return {handleAlbumSelected, handleNewAlbumSelected, handleClosed, ...toRefs(state)};
    },
  });
</script>
