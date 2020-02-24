<template>
    <main>
        <aside v-if="$root.loggedIn && selecting" class="selectionbar">
            {{ selectionCount }} selected
            <button @click="handleAddToAlbum">Add to album</button>
            <button @click="handleRemoveFromAlbum" v-if="!!album">Remove from album</button>
            <button @click="clearSelection">Clear selection</button>
        </aside>
        <p v-if="loading">Loading...</p>

        <router-view @go-to-previous-image="handleLightboxPrevious"
                     @go-to-next-image="handleLightboxNext"
        ></router-view>

        <thumbnail v-for="photo in photos"
                   :id="photo.id"
                   :caption="photo.file_name"
                   :key="photo.id"
                   :selected="selection[photo.id]"
                   :selecting="selecting"
                   @selected="handleItemSelected"
                   @deselected="handleItemDeselected"
                   @showing-photo="handleLightboxDisplayed"
                   @select-range="handleSelectRange"
        ></thumbnail>
    </main>
</template>

<style scoped>
    main {
        display: flex;
        flex-wrap: wrap;
        align-content: stretch;
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
  import _ from 'lodash'
  import Axios from 'axios'
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
        lastSelection: null,
        loading: true,
        offset: 0,
        photos: [],
        selection: {},
        showing: null
      }
    },
    computed: {
      selecting () {
        return !_.isEmpty(this.selection)
      },
      selectionCount () {
        return Object.keys(this.selection).length
      }
    },
    methods: {
      clearSelection () {
        this.selection = {}
        this.lastSelection = null
      },
      handleAddToAlbum () {
        new Promise((resolve, reject) => {
          EventBus.$emit('pick-album', resolve, reject)
        }).then(album => Axios.post('/albums/' + album + '/photos', {
          add_photos: this.selection
        }).then(() => {
          EventBus.$emit('toast', this.selectionCount + ' photo' + (this.selectionCount === 1 ? '' : 's') + ' added to album')
          EventBus.$emit('album-updated', album)
          this.clearSelection()
        }))
      },
      handleItemDeselected (id) {
        this.$delete(this.selection, id)
        this.lastSelection = null
      },
      handleItemSelected (id) {
        this.$set(this.selection, id, true)
        this.lastSelection = id
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
        Axios.post(this.endpoint, { remove_photos: this.selection }).then(() => {
          EventBus.$emit('toast', this.selection.length + ' photo' + (this.selection.length === 1 ? '' : 's') + ' removed from album')
          EventBus.$emit('album-updated', this.album)
          this.selection = {}
        })
      },
      handleSelectRange (id) {
        if (this.selection.length === 0 || this.lastSelection === null) {
          this.handleItemSelected(id)
        } else {
          const lastIndex = _.findIndex(this.photos, { id: this.lastSelection })
          const ourIndex = _.findIndex(this.photos, { id: id })

          let slice = []
          if (lastIndex < ourIndex) {
            slice = this.photos.slice(lastIndex + 1, ourIndex + 1)
          } else if (ourIndex < lastIndex) {
            slice = this.photos.slice(ourIndex, lastIndex)
          }

          slice.map((p) => p.id).forEach((id) => {
            this.$set(this.selection, id, true)
          })

          this.lastSelection = id
        }
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

        Axios.get(this.endpoint + '?offset=' + this.offset).then(({ data }) => {
          if (this.offset === 0) {
            this.photos = data
          } else {
            this.photos = this.photos.concat(data)
          }
          this.offset = this.offset + data.length
          this.hasMore = data.length > 0
          this.loading = false
        })
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
