<template>
    <popup title="Create new album" :modal="true" position="center" v-if="visible" @close="handleClosed">
        <form @submit="doCreate">
            <p class="alert" v-if="alert.length > 0">{{ alert }}</p>
            <label for="name">Name</label>
            <input type="text" id="name" v-model="name" placeholder="My Holiday" v-focus>
            <input type="submit" value="Create">
        </form>
    </popup>
</template>

<style scoped>
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
  import Popup from './popup'
  import { EventBus } from './bus'

  export default {
    components: { Popup },
    data () {
      return {
        alert: '',
        name: '',
        reject () {},
        resolve () {},
        visible: false
      }
    },
    methods: {
      doCreate () {
        Axios.post('/albums', { name: this.name }).then(({ data: { id } }) => {
          this.visible = false
          this.resolve(id)
        }).catch((error) => {
          if (error.response) {
            this.alert = error.response.data.error
          } else {
            this.alert = error.message
          }
        })
      },
      handleClosed () {
        this.visible = false
        this.reject()
      },
      show (resolve, reject) {
        this.name = ''
        this.resolve = resolve
        this.reject = reject
        this.visible = true
      }
    },
    created () {
      EventBus.$on('create-album', this.show)
    },
    beforeDestroy () {
      EventBus.$off('create-album', this.show)
    }
  }
</script>
