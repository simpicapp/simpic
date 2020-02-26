<template>
    <main>
        <gallery-toolbar v-if="$root.loggedIn"
                         :album="album"
                         :selection="selection"
                         :selection-count="selectionCount"
                         @clear-selection="clearSelection">
        </gallery-toolbar>

        <router-view @go-to-previous-image="handleLightboxPrevious"
                     @go-to-next-image="handleLightboxNext"
        ></router-view>

        <thumbnail v-for="photo in photos"
                   :imageId="photo.id"
                   :caption="photo.file_name"
                   :key="photo.id"
                   :selected="selection[photo.id]"
                   :selecting="selecting"
                   @selected="handleItemSelected"
                   @deselected="handleItemDeselected"
                   @select-range="handleSelectRange"
        ></thumbnail>

        <spinner v-if="loading"></spinner>

        <div class="nothing-here" v-if="!loading && photos.length === 0">
            <div>
                <p>There's nothing here</p>
                <p v-if="!$root.loggedIn">
                    You might need to login to see this content.
                </p>
                <p v-else-if="!!album">
                    You can upload pictures to Simpic simply by dragging and dropping them into your browser.
                    Give it a try!
                </p>
                <p v-else>
                    You can add pictures to albums by selecting them from the timeline.
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
        align-content: stretch;
    }
</style>

<script>
  import _ from 'lodash'
  import Axios from 'axios'
  import { EventBus } from './bus'
  import Thumbnail from './thumbnail'
  import GalleryToolbar from './gallery-toolbar'
  import Spinner from './spinner'
  import { cache } from './cache'

  export default {
    components: {
      GalleryToolbar,
      Spinner,
      Thumbnail
    },
    props: ['album', 'endpoint'],
    data: function () {
      return {
        hasMore: true,
        lastSelection: null,
        loading: true,
        offset: 0,
        photos: [],
        selection: {}
      }
    },
    computed: {
      selecting () {
        return this.selectionCount > 0
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
      handleItemDeselected (id) {
        this.$delete(this.selection, id)
        this.lastSelection = null
      },
      handleItemSelected (id) {
        this.$set(this.selection, id, true)
        this.lastSelection = id
      },
      handleLightboxNext (id) {
        const index = (_.findIndex(this.photos, { id }) + 1) % this.photos.length
        this.$router.push({ path: this.photos[index].id })
      },
      handleLightboxPrevious (id) {
        const index = (_.findIndex(this.photos, { id }) - 1 + this.photos.length) % this.photos.length
        this.$router.push({ path: this.photos[index].id })
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
        this.selection = {}
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
          cache.storeMetadata(data)
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
