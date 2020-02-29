<template>
  <main>
    <gallery-toolbar
      :album="album"
      :selection="selection"
      :selection-count="selectionCount"
      @clear-selection="clearSelection"
      v-if="loggedIn"
    >
    </gallery-toolbar>

    <router-view @go-to-next-image="handleLightboxNext" @go-to-previous-image="handleLightboxPrevious"></router-view>

    <thumbnail
      :caption="photo.file_name"
      :imageId="photo.id"
      :key="photo.id"
      :selected="selection[photo.id]"
      :selecting="selecting"
      @deselected="handleItemDeselected"
      @select-range="handleSelectRange"
      @selected="handleItemSelected"
      v-for="photo in photos"
    ></thumbnail>

    <spinner v-if="loading"></spinner>

    <div class="nothing-here" v-if="!loading && photos.length === 0">
      <div>
        <p>There's nothing here</p>
        <p v-if="!loggedIn">
          You might need to login to see this content.
        </p>
        <p v-else-if="!!album">
          You can upload pictures to Simpic simply by dragging and dropping them into your browser. Give it a try!
        </p>
        <p v-else>
          You can add pictures to albums by selecting them from the timeline.
        </p>
      </div>
    </div>
  </main>
</template>

<style lang="scss" scoped>
  @import "src/assets/css/nothing-here";

  main {
    display: flex;
    flex-wrap: wrap;
    align-content: stretch;
  }
</style>

<script lang="ts">
  import {findIndex} from "lodash-es";
  import Axios from "axios";
  import Thumbnail from "./thumbnail.vue";
  import GalleryToolbar from "./gallery-toolbar.vue";
  import Spinner from "./spinner.vue";
  import {cache} from "./cache";
  import {computed, defineComponent, onMounted, reactive, toRefs} from "@vue/composition-api";
  import {useAuthentication} from "@/features/auth";
  import {useScrollWatcher} from "@/features/scroll";
  import Vue from "vue";
  import {useRouter} from "@/features/router";
  import {Data} from "@vue/composition-api/dist/component";
  import {useEventListener} from "@/features/eventbus";
  import {Photo} from "@/model/photo";

  export default defineComponent({
    components: {
      GalleryToolbar,
      Spinner,
      Thumbnail,
    },
    props: {
      album: String,
      endpoint: String,
    },
    setup(props) {
      const {router} = useRouter();
      const {loggedIn} = useAuthentication();

      interface Selection {
        [key: string]: boolean;
      }

      interface State extends Data {
        hasMore: boolean;
        lastSelection: null | string;
        loading: boolean;
        offset: number;
        photos: Array<Photo>;
        selection: Selection;
      }

      const state: State = reactive({
        hasMore: true,
        lastSelection: null,
        loading: true,
        offset: 0,
        photos: [],
        selection: {},
      });

      function update() {
        state.loading = true;

        Axios.get(props.endpoint + "?offset=" + state.offset).then(({data}) => {
          if (state.offset === 0) {
            state.photos = data;
          } else {
            state.photos = state.photos.concat(data);
          }
          state.offset = state.offset + data.length;
          state.hasMore = data.length > 0;
          state.loading = false;
          cache.storeMetadata(data);
        });
      }

      onMounted(update);

      useScrollWatcher(() => {
        if (!state.loading && state.hasMore) {
          update();
        }
      });

      useEventListener("refresh-gallery", () => {
        state.selection = {};
        state.offset = 0;
        state.hasMore = true;
        update();
      });

      const selectionCount = computed(() => Object.keys(state.selection).length);
      const selecting = computed(() => selectionCount.value > 0);

      function clearSelection() {
        state.selection = {};
        state.lastSelection = null;
      }

      function handleItemDeselected(id: string) {
        Vue.delete(state.selection, id);
        state.lastSelection = null;
      }

      function handleItemSelected(id: string) {
        Vue.set(state.selection, id, true);
        state.lastSelection = id;
      }

      function handleLightboxNext(id: string) {
        const index = (findIndex(state.photos, {id}) + 1) % state.photos.length;
        router.push({path: state.photos[index].id});
      }

      function handleLightboxPrevious(id: string) {
        const index = (findIndex(state.photos, {id}) - 1 + state.photos.length) % state.photos.length;
        router.push({path: state.photos[index].id});
      }

      function handleSelectRange(id: string) {
        if (Object.keys(state.selection).length === 0 || state.lastSelection === null) {
          handleItemSelected(id);
        } else {
          const lastIndex = findIndex(state.photos, {id: state.lastSelection});
          const ourIndex = findIndex(state.photos, {id: id});

          let slice: Photo[] = [];
          if (lastIndex < ourIndex) {
            slice = state.photos.slice(lastIndex + 1, ourIndex + 1);
          } else if (ourIndex < lastIndex) {
            slice = state.photos.slice(ourIndex, lastIndex);
          }

          slice
            .map(p => p.id)
            .forEach(id => {
              Vue.set(state.selection, id, true);
            });

          state.lastSelection = id;
        }
      }

      return {
        loggedIn,
        selecting,
        selectionCount,
        clearSelection,
        handleLightboxNext,
        handleLightboxPrevious,
        handleItemSelected,
        handleItemDeselected,
        handleSelectRange,
        ...toRefs(state),
      };
    },
  });
</script>
