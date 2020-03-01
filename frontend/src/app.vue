<template>
  <div>
    <toolbar></toolbar>
    <div id="content" v-if="ready">
      <router-view></router-view>
      <uploader></uploader>
      <login></login>
      <album-picker></album-picker>
      <toaster></toaster>
    </div>
    <spinner v-else></spinner>
    <bottom-bar></bottom-bar>
  </div>
</template>

<style lang="scss">
  @use 'assets/css/vars';
  @import "assets/css/global";
  @import "assets/css/reset";

  #content {
    min-height: calc(100vh - 90px - 80px);
  }
</style>

<script lang="ts">
  import Vue from "vue";
  import {defineComponent, onMounted, reactive, toRefs} from "@vue/composition-api";

  import Uploader from "./components/uploader.vue";
  import Login from "./components/login.vue";
  import Toolbar from "./components/toolbar.vue";
  import BottomBar from "./components/footer.vue";
  import AlbumPicker from "./components/album-picker.vue";
  import Spinner from "./components/spinner.vue";
  import Toaster from "./components/toaster.vue";
  import {useAuthentication} from "@/features/auth";

  Vue.directive("focus", {
    inserted: function(el) {
      el.focus();
    },
  });

  export default defineComponent({
    components: {
      AlbumPicker,
      BottomBar,
      Login,
      Spinner,
      Toaster,
      Toolbar,
      Uploader,
    },
    setup() {
      const state = reactive({
        ready: false,
      });

      const {checkUser} = useAuthentication();

      onMounted(() => {
        checkUser().finally(() => {
          state.ready = true;
        });
      });

      return toRefs(state);
    },
  });
</script>
