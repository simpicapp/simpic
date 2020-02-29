<template>
    <div>
        <div class="toolbar">
            <h2>{{ name }}</h2>
            <ActionIcon v-if="$root.loggedIn" name="trash-alt" :working="deleting"
                        @click="confirmDelete()"></ActionIcon>
        </div>
        <gallery :endpoint="'/albums/' + id + '/photos'" :album="id"></gallery>

        <DeleteDialog v-if="showConfirmation"
                      what="this album"
                      @yes="doDelete"
                      @close="showConfirmation = false">
        </DeleteDialog>
    </div>
</template>

<style lang="scss" scoped>
    h2 {
        margin: 0 20px;
    }

    .toolbar {
        display: flex;
        align-items: center;
        margin-bottom: 10px;
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
  import 'vue-awesome/icons/trash-alt'
  import Axios from 'axios'
  import Gallery from './gallery'
  import { EventBus } from './bus'
  import ActionIcon from './action-icon'
  import DeleteDialog from './delete-dialog'

  export default {
    components: {
      ActionIcon,
      DeleteDialog,
      Gallery
    },
    props: ['id'],
    data () {
      return {
        deleting: false,
        name: '',
        showConfirmation: false
      }
    },
    methods: {
      confirmDelete () {
        this.showConfirmation = true
      },
      doDelete () {
        this.deleting = true
        Axios.delete('/albums/' + this.id).then(() => {
          EventBus.$emit('albums-updated')
          EventBus.$emit('toast', 'Album deleted')
          this.deleting = false
          this.$router.replace('/albums')
        })
      },
      handleAlbumUpdated (album) {
        if (album === this.id) {
          EventBus.$emit('refresh-gallery')
        }
      }
    },
    mounted () {
      Axios.get('/albums/' + this.id).then(({ data: { name } }) => {
        this.name = name
      })

      EventBus.$on('album-updated', this.handleAlbumUpdated)
    },
    destroyed () {
      EventBus.$off('album-updated', this.handleAlbumUpdated)
    }
  }
</script>
