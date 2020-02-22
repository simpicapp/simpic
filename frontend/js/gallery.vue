<template>
    <main class="gallery">
        <aside v-if="selecting" class="selectionbar">
            {{ selection.length }} selected
            <button v-on:click="handleAddToAlbum">Add to album</button>
            <button v-on:click="handleRemoveFromAlbum" v-if="!!album">Remove from album</button>
            <button v-on:click="clearSelection">Clear selection</button>
        </aside>
        <p v-if="loading">Loading...</p>

        <router-view v-on:go-to-previous-image="handleLightboxPrevious"
                     v-on:go-to-next-image="handleLightboxNext"
        ></router-view>

        <thumbnail v-for="photo in photos"
                   v-bind:id="photo.id"
                   v-bind:caption="photo.file_name"
                   v-bind:key="photo.id"
                   v-bind:selecting="selecting"
                   v-on:selected="handleItemSelected"
                   v-on:deselected="handleItemDeselected"
                   v-on:showing-photo="handleLightboxDisplayed"
        ></thumbnail>
    </main>
</template>

<style>
    .gallery {
        display: flex;
        flex-wrap: wrap;
    }

    .selectionbar {
        position: fixed;
        z-index: 800;
        top: 0;
        left: 25%;
        right: 25%;
        border: 2px solid black;
        border-top: 0;
        padding: 25px;
        border-bottom-right-radius: 10px;
        border-bottom-left-radius: 10px;
        background: #ffffff;
        display: flex;
        justify-content: space-between;
    }
</style>

<script>
  import { EventBus } from './bus'
  import thumbnail from './thumbnail'

  export default {
    components: {
      thumbnail
    },
    props: ['album', 'endpoint'],
    data: function () {
      return {
        hasMore: true,
        loading: true,
        offset: 0,
        photos: [],
        selection: [],
        showing: null
      }
    },
    computed: {
      selecting () {
        return this.selection.length > 0
      }
    },
    methods: {
      clearSelection () {
        this.selection = []
      },
      handleAddToAlbum () {
        EventBus.$emit('pick-album')
      },
      handleAlbumSelected (album) {
        if (!!album && this.selection.length > 0) {
          fetch('/albums/' + album + '/photos', {
            body: JSON.stringify({
              add_photos: this.selection
            }),
            headers: {
              Accept: 'application/json',
              'Content-Type': 'application/json'
            },
            method: 'POST'
          }).then(() => {
            EventBus.$emit('album-updated', album)
            this.selection = []
          })
        }
      },
      handleItemDeselected (id) {
        this.selection.splice(this.selection.indexOf(id), 1)
      },
      handleItemSelected (id) {
        this.selection.push(id)
      },
      handleLightboxDisplayed (id) {
        const comp = this
        this.photos.forEach(function (photo, index) {
          if (photo.id === id) {
            comp.showing = index
          }
        })
      },
      handleLightboxNext () {
        this.showing = (this.showing + 1) % this.photos.length
        this.$router.push({ path: this.photos[this.showing].id })
      },
      handleLightboxPrevious () {
        this.showing = (this.showing - 1 + this.photos.length) % this.photos.length
        this.$router.push({ path: this.photos[this.showing].id })
      },
      handleRemoveFromAlbum () {
        fetch(this.endpoint, {
          body: JSON.stringify({
            remove_photos: this.selection
          }),
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
          },
          method: 'POST'
        }).then(() => {
          EventBus.$emit('album-updated', this.album)
          this.selection = []
        })
      },
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
      EventBus.$off('album-selected', this.handleAlbumSelected)
    },
    mounted () {
      this.update()
      EventBus.$on('bottom', this.infiniteScroll)
      EventBus.$on('refresh-gallery', this.refresh)
      EventBus.$on('album-selected', this.handleAlbumSelected)
    }
  }
</script>
