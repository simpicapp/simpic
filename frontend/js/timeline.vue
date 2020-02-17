<template>
    <main class="timeline">
        <p v-if="loading">Loading...</p>
        <router-view></router-view>
        <thumbnail v-for="photo in photos" v-bind:photo="photo" v-bind:key="photo.id"></thumbnail>
    </main>
</template>

<style>
    .timeline {
        display: flex;
        flex-wrap: wrap;
    }
</style>

<script>
  import { EventBus } from './bus'
  import thumbnail from './thumbnail'

  export default {
    components: {
      thumbnail
    },
    data: function () {
      return {
        loading: true,
        photos: []
      }
    },
    methods: {
      update () {
        const comp = this
        fetch('/timeline', {
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            ...this.$root.authHeaders()
          }
        })
          .then((response) => response.json())
          .then((json) => (comp.photos = json))
          .then(() => (comp.loading = false))
      }
    },
    mounted: function () {
      this.update()
      EventBus.$on('upload-complete', this.update)
    }
  }
</script>
