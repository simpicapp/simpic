<template>
    <main>
        <Album v-for="album in albums"
               :key="album.id"
               :id="album.id"
               :caption="album.name"
               :cover="album.cover_photo">
        </Album>
    </main>
</template>

<style scoped>
    main {
        display: flex;
        flex-wrap: wrap;
    }
</style>

<script>
  import Album from './album-icon'
  import Axios from 'axios'

  export default {
    components: { Album },
    data () {
      return {
        albums: []
      }
    },
    mounted () {
      Axios.get('albums').then(({ data }) => {
        this.albums = data
      })
    }
  }
</script>
