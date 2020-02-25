<template>
    <main>
        <Album v-for="album in albums"
               :key="album.id"
               :id="album.id"
               :caption="album.name"
               :cover="album.cover_photo">
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
    @import "../css/nothing-here";

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
        albums: [],
        loading: true
      }
    },
    mounted () {
      Axios.get('albums').then(({ data }) => {
        this.albums = data
        this.loading = false
      })
    }
  }
</script>
