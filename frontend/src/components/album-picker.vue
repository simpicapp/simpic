<template>
    <modal v-if="visible" @close="visible = false" :should-close="close">
        <popup title="Select an Album" position="center" @close="close = true" v-if="selecting">
            <div class="album-picker">
                <template v-for="album in albums">
                    <img :key="album.id"
                         v-if="album.cover_photo"
                         :src="'/data/thumb/' + album.cover_photo"
                         :alt="album.name"
                         @click="handleAlbumSelected(album.id)">
                    <span :key="album.id" v-else></span>
                    <div :key="album.id + '.name'" @click="handleAlbumSelected(album.id)">
                        <span>{{ album.name }}</span>
                    </div>
                </template>
                <div class="icon" @click="handleNewAlbumSelected"><span>⊕</span></div>
                <div @click="handleNewAlbumSelected"><span>Create new album...</span></div>
            </div>
        </popup>
        <album-creator v-else
                @created="handleAlbumSelected"
                @close="handleClosed">
        </album-creator>
    </modal>
</template>

<style lang="scss" scoped>
    .album-picker {
        display: grid;
        grid-template-columns: 50px 200px;
        grid-auto-rows: 2em;
        grid-gap: 20px 20px;
        align-items: center;

        div {
            align-self: stretch;
            display: flex;
            align-items: center;
            cursor: pointer;

            &.icon {
                justify-content: center;
            }
        }
    }

    img {
        max-height: 2em;
        max-width: 50px;
        margin-right: 10px;
        border: 1px solid black;
        overflow: hidden;
        cursor: pointer;
    }

    form {
        display: grid;
        grid-template-columns: auto auto;
        grid-gap: 30px 20px;
        align-items: center;
    }

    input[type=submit] {
        grid-column: span 2;
    }

    .alert {
        margin: 0;
        padding: 5px 10px;
        grid-column: span 2;
        background-color: darkred;
        color: white;
        font-weight: bold;
        text-align: center;
        border-radius: 15px;
        white-space: pre-line;
    }
</style>

<script>
  import Axios from 'axios'
  import Modal from './modal'
  import Popup from './popup'
  import { EventBus } from './bus'
  import AlbumCreator from './album-creator'

  export default {
    components: {
      AlbumCreator,
      Modal,
      Popup
    },
    data () {
      return {
        albums: [],
        alert: '',
        close: false,
        name: '',
        reject () { /* noop */ },
        resolve () { /* noop */ },
        selecting: true,
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
        this.selecting = false
      },
      show (resolve, reject) {
        this.albums = []
        this.resolve = resolve
        this.reject = reject
        this.close = false
        this.selecting = true
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
