<template>
    <div>
        <toolbar></toolbar>
        <router-view></router-view>
        <uploader></uploader>
        <login></login>
        <bottombar></bottombar>
    </div>
</template>

<style>
    body {
        font-family: sans-serif;
        margin: 0;
        padding: 0;
    }
</style>

<script>
  import _ from 'lodash'
  import uploader from './uploader'
  import login from './login'
  import toolbar from './toolbar'
  import bottombar from './footer'
  import { EventBus } from './bus'

  export default {
    components: {
      bottombar,
      login,
      toolbar,
      uploader
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
        const bottomOfPage = visible + scrollY >= pageHeight
        return bottomOfPage || pageHeight < visible
      },
      emitBottom: _.throttle(() => EventBus.$emit('bottom'), 250),
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
  }
</script>
