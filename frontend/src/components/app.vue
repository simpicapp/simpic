<template>
  <div>
    <toolbar></toolbar>
    <div id="content">
      <router-view></router-view>
      <uploader></uploader>
      <login></login>
      <album-picker></album-picker>
      <toaster></toaster>
    </div>
    <bottom-bar></bottom-bar>
  </div>
</template>

<style lang="scss">
  @use '../assets/css/vars';
  @import "../assets/css/global";
  @import "../assets/css/reset";

  body {
    font-family: sans-serif;
    margin: 0;
    padding: 0;
  }

  :root {
    --ease-in-cubic: cubic-bezier(0.55, 0.055, 0.675, 0.19);
    --ease-out-cubic: cubic-bezier(0.645, 0.045, 0.355, 1);
    --ease-out-back: cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }

  input {
    padding: 5px;
  }

  input[type="submit"] {
    padding: 10px;
  }

  input[type="submit"]:hover,
  input[type="submit"]:active {
    background-color: vars.$primary;
  }

  #content {
    min-height: calc(100vh - 90px - 80px);
  }
</style>

<script lang="ts">
  import Vue from "vue";
  import {defineComponent} from "@vue/composition-api";

  import Uploader from "./uploader.vue";
  import Login from "./login.vue";
  import Toolbar from "./toolbar.vue";
  import BottomBar from "./footer.vue";
  import AlbumPicker from "./album-picker.vue";
  import Toaster from "./toaster.vue";
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
      Toaster,
      Toolbar,
      Uploader,
    },
    setup() {
      const {checkUser} = useAuthentication();
      checkUser();
    },
  });
</script>
