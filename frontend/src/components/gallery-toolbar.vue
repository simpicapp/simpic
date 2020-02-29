<template>
  <div>
    <aside class="selectionbar" v-if="selectionCount > 0">
      {{ selectionCount }} selected
      <button @click="handleAddToAlbum">Add to album</button>
      <button @click="handleRemoveFromAlbum" v-if="!!album">Remove from album</button>
      <button @click="handleDelete">Delete</button>
      <button @click="$emit('clear-selection')">Clear selection</button>
    </aside>

    <DeleteDialog :what="`${selectionCount} photo${selectionCount === 1 ? '' : 's'}`"
                  @close="showConfirmation = false"
                  @yes="doDelete"
                  v-if="showConfirmation">
    </DeleteDialog>
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

<script lang="ts">
  import {EventBus} from './bus'
  import Axios from 'axios'
  import DeleteDialog from './delete-dialog.vue'
  import Vue from 'vue'

  export default Vue.extend({
    components: {
      DeleteDialog
    },
    props: [
      'selectionCount',
      'selection',
      'album'
    ],
    data() {
      return {
        showConfirmation: false
      }
    },
    methods: {
      doDelete() {
        Axios.post('/photos/delete', {photos: Object.keys(this.selection)}).then(() => {
          EventBus.$emit('toast', this.selectionCount + ' photo' + (this.selectionCount === 1 ? '' : 's') + ' deleted');
          EventBus.$emit('refresh-gallery')
        })
      },
      handleAddToAlbum() {
        new Promise((resolve, reject) => {
          EventBus.$emit('pick-album', resolve, reject)
        }).then(album => Axios.post('/albums/' + album + '/photos', {
          // eslint-disable-next-line @typescript-eslint/camelcase
          add_photos: Object.keys(this.selection)
        }).then(() => {
          EventBus.$emit('toast', this.selectionCount + ' photo' + (this.selectionCount === 1 ? '' : 's') + ' added to album');
          EventBus.$emit('album-updated', album);
          this.$emit('clear-selection')
        }))
      },
      handleDelete() {
        this.showConfirmation = true
      },
      handleRemoveFromAlbum() {
        // eslint-disable-next-line @typescript-eslint/camelcase
        Axios.post('/albums/' + this.album + '/photos', {remove_photos: Object.keys(this.selection)}).then(() => {
          EventBus.$emit('toast', this.selectionCount + ' photo' + (this.selectionCount === 1 ? '' : 's') + ' removed from album');
          EventBus.$emit('album-updated', this.album);
          this.$emit('clear-selection')
        })
      }
    }
  })
</script>
