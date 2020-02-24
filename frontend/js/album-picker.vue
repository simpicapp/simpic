<template>
    <popup title="Select an Album" position="center" v-bind:modal="true" v-if="visible" v-on:close="handleClosed">
        <div class="album-picker">
            <template v-for="album in albums">
                <img v-bind:key="album.id"
                     v-if="album.cover_photo"
                     v-bind:src="'/data/thumb/' + album.cover_photo"
                     v-bind:alt="album.name"
                     v-on:click="handleAlbumSelected(album.id)">
                <span v-bind:key="album.id" v-else></span>
                <div v-bind:key="album.id + '.name'" v-on:click="handleAlbumSelected(album.id)">
                    <span>{{ album.name }}</span>
                </div>
            </template>
            <div class="icon" v-on:click="handleNewAlbumSelected()"><span>⊕</span></div>
            <div v-on:click="handleNewAlbumSelected()"><span>Create new album...</span></div>
        </div>
    </popup>
</template>

<style scoped>
    .album-picker {
        display: grid;
        grid-template-columns: 50px 200px;
        grid-auto-rows: 2em;
        grid-gap: 20px 20px;
        align-items: center;
    }

    img {
        max-height: 2em;
        max-width: 50px;
        margin-right: 10px;
        border: 1px solid black;
        overflow: hidden;
        cursor: pointer;
    }

    .album-picker div {
        align-self: stretch;
        display: flex;
        align-items: center;
        cursor: pointer;
    }

    div.icon {
        justify-content: center;
    }
</style>

<script>
  import Axios from 'axios'
  import Popup from './popup'
  import { EventBus } from './bus'

  export default {
    components: { Popup },
    data () {
      return {
        albums: [],
        reject () {},
        resolve () {},
        visible: false
      }
    },
    methods: {
      handleAlbumSelected (albumId) {
        this.visible = false
        this.resolve(albumId)
      },
      handleClosed () {
        this.visible = false
        this.reject()
      },
      handleNewAlbumSelected () {
        this.visible = false
        EventBus.$emit('create-album', this.resolve, this.reject)
      },
      show (resolve, reject) {
        this.albums = []
        this.resolve = resolve
        this.reject = reject
        this.visible = true

        Axios.get('/albums').then(({ data }) => {
          this.albums = data
        })
      }
    },
    created () {
      EventBus.$on('pick-album', this.show)
    },
    beforeDestroy () {
      EventBus.$off('pick-album', this.show)
    }
  }
</script>
