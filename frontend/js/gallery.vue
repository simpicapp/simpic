<template>
    <main class="gallery">
        <p v-if="loading">Loading...</p>
        <router-view></router-view>
        <thumbnail v-for="photo in photos"
                   v-bind:id="photo.id"
                   v-bind:caption="photo.file_name"
                   v-bind:key="photo.id"></thumbnail>
    </main>
</template>

<style>
    .gallery {
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
    props: ['endpoint'],
    data: function () {
      return {
        hasMore: true,
        loading: true,
        offset: 0,
        photos: []
      }
    },
    methods: {
      infiniteScroll () {
        if (!this.loading && this.hasMore) {
          this.update()
        }
      },
      refresh () {
        this.offset = 0
        this.hasMore = true
        this.update()
      },
      update () {
        this.loading = true

        fetch(this.endpoint + '?offset=' + this.offset, {
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
          }
        }).then((response) => response.json())
          .then((json) => {
            if (this.offset === 0) {
              this.photos = json
            } else {
              this.photos = this.photos.concat(json)
            }
            this.offset = this.offset + json.length
            this.hasMore = json.length > 0
          })
          .then(() => (this.loading = false))
      }
    },
    beforeDestroy () {
      EventBus.$off('bottom', this.infiniteScroll)
      EventBus.$off('refresh-gallery', this.refresh)
    },
    mounted () {
      this.update()
      EventBus.$on('bottom', this.infiniteScroll)
      EventBus.$on('refresh-gallery', this.refresh)
    }
  }
</script>
