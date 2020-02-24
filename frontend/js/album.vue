<template>
    <div>
        <h2>{{ name }}</h2>
        <gallery :endpoint="'/albums/' + id + '/photos'" :album="id"></gallery>
    </div>
</template>

<style scoped>
    h2 {
        margin-left: 20px;
    }
</style>

<script>
  import Axios from 'axios'
  import Gallery from './gallery'
  import { EventBus } from './bus'

  export default {
    components: {
      Gallery
    },
    data () {
      return {
        name: ''
      }
    },
    methods: {
      handleAlbumUpdated (album) {
        if (album === this.id) {
          EventBus.$emit('refresh-gallery')
        }
      }
    },
    props: ['id'],
    mounted () {
      Axios.get('/albums/' + this.id).then(({ data: { name } }) => {
        this.name = name
      })

      EventBus.$on('album-updated', this.handleAlbumUpdated)
    },
    destroyed () {
      EventBus.$off('album-updated', this.handleAlbumUpdated)
    }
  }
</script>
