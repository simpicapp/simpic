<template>
    <main class="gallery">
        <aside v-if="selecting" class="selectionbar">
            {{ selection.length }} selected
            <button v-on:click="clearSelection">Clear selection</button>
        </aside>
        <p v-if="loading">Loading...</p>
        <router-view></router-view>
        <thumbnail v-for="photo in photos"
                   v-bind:id="photo.id"
                   v-bind:caption="photo.file_name"
                   v-bind:key="photo.id"
                   v-bind:selecting="selecting"
                   v-on:selected="handleSelected"
                   v-on:deselected="handleDeselected"
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
    props: ['endpoint'],
    data: function () {
      return {
        hasMore: true,
        loading: true,
        offset: 0,
        photos: [],
        selection: []
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
      handleDeselected (id) {
        this.selection.splice(this.selection.indexOf(id), 1)
      },
      handleSelected (id) {
        this.selection.push(id)
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
    },
    mounted () {
      this.update()
      EventBus.$on('bottom', this.infiniteScroll)
      EventBus.$on('refresh-gallery', this.refresh)
    }
  }
</script>
