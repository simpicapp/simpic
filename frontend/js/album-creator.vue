<template>
    <popup title="Create new album" position="center" @close="handleClosed">
        <form @submit="doCreate">
            <p class="alert" v-if="alert.length > 0">{{ alert }}</p>
            <label for="name">Name</label>
            <input type="text" id="name" v-model="name" placeholder="My Holiday" v-focus>
            <input type="submit" value="Create">
        </form>
    </popup>
</template>

<style lang="scss" scoped>
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

  export default {
    components: { Popup },
    data () {
      return {
        alert: '',
        name: ''
      }
    },
    methods: {
      doCreate () {
        Axios.post('/albums', { name: this.name }).then(({ data: { id } }) => {
          this.$emit('created', id)
          this.name = ''
        }).catch((error) => {
          if (error.response) {
            this.alert = error.response.data.error
          } else {
            this.alert = error.message
          }
        })
      },
      handleClosed () {
        this.$emit('close')
        this.name = ''
      }
    }
  }
</script>
