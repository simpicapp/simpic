<template>
    <main>
        <Album v-for="album in albums"
               :key="album.id"
               :id="album.id"
               :caption="album.name"
               :imageId="album.cover_photo"
               :photos="album.photos">
        </Album>

        <div class="nothing-here" v-if="!loading && albums.length === 0">
            <div>
                <p>There's nothing here</p>
                <p v-if="!$root.loggedIn">
                    You might need to login to see this content.
                </p>
                <p v-else>
                    You can create albums by selecting some photos in the timeline.
                </p>
            </div>
        </div>
    </main>
</template>

<style lang="scss" scoped>
    @import "src/assets/css/nothing-here";

    main {
        display: flex;
        flex-wrap: wrap;
        padding: 20px;
    }
</style>

<script>
  import Album from './album-icon'
  import Axios from 'axios'
  import { EventBus } from './bus'

  export default {
    components: { Album },
    data () {
      return {
        albums: [],
        loading: true
      }
    },
    methods: {
      refresh () {
        Axios.get('albums').then(({ data }) => {
          this.albums = data
          this.loading = false
        })
      }
    },
    mounted () {
      this.refresh()
      EventBus.$on('albums-updated', this.refresh)
    },
    beforeDestroy () {
      EventBus.$off('albums-updated', this.refresh)
    }
  }
</script>
