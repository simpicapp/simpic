<template>
    <main class="albums">
        <Album v-for="album in albums" v-bind:key="album.id" v-bind:id="album.id" v-bind:caption="album.name"
               v-bind:cover="album.cover_photo"></Album>
    </main>
</template>

<style>
    .albums {
        display: flex;
        flex-wrap: wrap;
    }
</style>

<script>
  import Album from './album-icon'

  export default {
    components: { Album },
    data () {
      return {
        albums: []
      }
    },
    mounted () {
      fetch('/albums', {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json'
        }
      }).then(res => res.json())
        .then(albums => (this.albums = albums))
    }
  }
</script>
