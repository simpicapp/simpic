<template>
    <gallery endpoint="/timeline"></gallery>
</template>

<script>
  import { EventBus } from './bus'
  import Gallery from './gallery'

  export default {
    components: {
      Gallery
    },
    methods: {
      refresh () {
        EventBus.$emit('refresh-gallery')
      }
    },
    beforeDestroy () {
      EventBus.$off('upload-complete', this.refresh)
    },
    mounted () {
      EventBus.$on('upload-complete', this.refresh)
    }
  }
</script>
