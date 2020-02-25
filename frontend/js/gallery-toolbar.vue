<template>
    <div>
        <aside class="selectionbar">
            {{ selectionCount }} selected
            <button @click="handleAddToAlbum">Add to album</button>
            <button @click="handleRemoveFromAlbum" v-if="!!album">Remove from album</button>
            <button @click="handleDelete">Delete</button>
            <button @click="$emit('clear-selection')">Clear selection</button>
        </aside>

        <modal v-if="showConfirmation" :closeable="false" :should-close="closeConfirmation"
               @close="showConfirmation = false">
            <popup title="Confirm deletion" position="center" :closeable="false">
                <p>Are you sure you wish to delete {{ selectionCount }} photo{{ selectionCount === 1 ? '' : 's'}}?</p>
                <div class="buttons">
                    <button @click="doDelete" class="danger-button">Yes, delete</button>
                    <button @click="cancelDelete">No, cancel</button>
                </div>
            </popup>
        </modal>
    </div>
</template>

<style lang="scss" scoped>
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

    .buttons {
        display: grid;
        grid-template-columns: auto auto;
        grid-column-gap: 20px;
        justify-items: stretch;
        margin-top: 30px;
    }

</style>

<script>
  import { EventBus } from './bus'
  import Axios from 'axios'
  import Modal from './modal'
  import Popup from './popup'

  export default {
    components: {
      Modal,
      Popup
    },
    props: [
      'selectionCount',
      'selection',
      'album'
    ],
    data () {
      return {
        closeConfirmation: false,
        showConfirmation: false
      }
    },
    methods: {
      cancelDelete () {
        this.closeConfirmation = true
      },
      doDelete () {
        this.closeConfirmation = true
        Axios.post('/photos/delete', { photos: Object.keys(this.selection) }).then(() => {
          EventBus.$emit('toast', this.selectionCount + ' photo' + (this.selectionCount === 1 ? '' : 's') + ' deleted')
          this.$emit('clear-selection')
          EventBus.$emit('refresh-gallery')
        })
      },
      handleAddToAlbum () {
        new Promise((resolve, reject) => {
          EventBus.$emit('pick-album', resolve, reject)
        }).then(album => Axios.post('/albums/' + album + '/photos', {
          add_photos: Object.keys(this.selection)
        }).then(() => {
          EventBus.$emit('toast', this.selectionCount + ' photo' + (this.selectionCount === 1 ? '' : 's') + ' added to album')
          EventBus.$emit('album-updated', album)
          this.$emit('clear-selection')
        }))
      },
      handleDelete () {
        this.closeConfirmation = false
        this.showConfirmation = true
      },
      handleRemoveFromAlbum () {
        Axios.post('/albums/' + this.album + '/photos', { remove_photos: Object.keys(this.selection) }).then(() => {
          EventBus.$emit('toast', this.selectionCount + ' photo' + (this.selectionCount === 1 ? '' : 's') + ' removed from album')
          EventBus.$emit('album-updated', this.album)
          this.$emit('clear-selection')
        })
      }
    }
  }
</script>
