<template>
  <main>
    <Album
      :caption="album.name"
      :id="album.id"
      :imageId="album.cover_photo"
      :key="album.id"
      :photos="album.photos"
      v-for="album in albums"
    >
    </Album>

    <div class="nothing-here" v-if="!loading && albums.length === 0">
      <div>
        <p>There's nothing here</p>
        <p v-if="!loggedIn">
          You might need to login to see this content.
        </p>
        <p v-else>
          You can create albums by selecting some photos in the timeline.
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
    padding: 20px;
  }
</style>

<script lang="ts">
  import Album from "./album-icon.vue";
  import Axios from "axios";
  import {defineComponent, onMounted, reactive, toRefs} from "@vue/composition-api";
  import {useEventListener} from "@/features/eventbus";
  import {useAuthentication} from "@/features/auth";

  export default defineComponent({
    components: {Album},
    setup() {
      const {loggedIn} = useAuthentication();

      const state = reactive({
        albums: [],
        loading: true,
      });

      function refresh() {
        Axios.get("albums").then(({data}) => {
          state.albums = data;
          state.loading = false;
        });
      }

      useEventListener("albums-updated", refresh);
      useEventListener("user-changed", refresh);
      onMounted(refresh);

      return {loggedIn, ...toRefs(state)};
    },
  });
</script>
