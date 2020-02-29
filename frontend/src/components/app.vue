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
  @import '../assets/css/global';
  @import '../assets/css/reset';

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

  input[type=submit] {
    padding: 10px;
  }

  input[type=submit]:hover, input[type=submit]:active {
    background-color: vars.$primary;
  }

  #content {
    min-height: calc(100vh - 90px - 80px)
  }
</style>

<script>
  import { throttle } from 'lodash-es'
  import Vue from 'vue'

  import Uploader from './uploader'
  import Login from './login'
  import Toolbar from './toolbar'
  import BottomBar from './footer'
  import { EventBus } from './bus'
  import AlbumPicker from './album-picker'
  import Toaster from './toaster'

  Vue.directive('focus', {
    inserted: function (el) {
      el.focus()
    }
  })

  export default Vue.extend({
    components: {
      AlbumPicker,
      BottomBar,
      Login,
      Toaster,
      Toolbar,
      Uploader
    },
    data () {
      return {
        bottom: false
      }
    },
    methods: {
      bottomVisible () {
        const scrollY = window.scrollY
        const visible = document.documentElement.clientHeight
        const pageHeight = document.documentElement.scrollHeight
        return visible + scrollY >= pageHeight - 400
      },
      emitBottom: throttle(() => EventBus.$emit('bottom'), 250),
      handleScroll () {
        this.bottom = this.bottomVisible()
      }
    },
    beforeDestroy () {
      window.removeEventListener('scroll', this.handleScroll)
    },
    created () {
      window.addEventListener('scroll', this.handleScroll)
    },
    watch: {
      bottom (newVal) {
        if (newVal) {
          this.emitBottom()
        }
      }
    }
  })
</script>
