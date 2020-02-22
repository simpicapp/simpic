<template>
    <popup title="Select an Album" position="center" v-if="visible" v-on:close="handleAlbumSelected(null)">
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

<style>
    .album-picker {
        display: grid;
        grid-template-columns: 50px 200px;
        grid-auto-rows: 2em;
        padding: 20px;
        grid-gap: 10px 10px;
        align-items: center;
    }

    .album-picker img {
        max-height: 2em;
        max-width: 5em;
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

    .album-picker div.icon {
        justify-content: center;
    }
</style>

<script>
  import Popup from './popup'
  import { EventBus } from './bus'

  export default {
    components: { Popup },
    data () {
      return {
        albums: [],
        visible: false
      }
    },
    methods: {
      handleAlbumSelected (albumId) {
        this.visible = false
        EventBus.$emit('album-selected', albumId)
      },
      handleNewAlbumSelected () {
        this.visible = false
        EventBus.$emit('create-album')
      },
      show () {
        this.albums = []
        this.visible = true
        fetch('/albums', {
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
          }
        }).then(res => res.json())
          .then(albums => (this.albums = albums))
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
