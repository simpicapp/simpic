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

        fetch('/timeline?offset=' + this.offset, {
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            ...this.$root.authHeaders()
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
    mounted: function () {
      this.update()
      EventBus.$on('upload-complete', this.refresh)
      EventBus.$on('bottom', this.infiniteScroll)
    }
  }
</script>
