<template>
    <popup title="Login" position="center" v-bind:modal="true" v-if="visible" v-on:close="visible = false">
        <form v-on:submit="doLogin">
            <p class="alert" v-if="alert.length > 0">{{ alert }}</p>
            <label for="username">Username</label>
            <input type="text" id="username" v-model="username" v-bind:disabled="loggingIn" v-focus>
            <label for="password">Password</label>
            <input type="password" id="password" v-model="password" v-bind:disabled="loggingIn">
            <input type="submit" value="Login" v-bind:disabled="loggingIn">
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
  import { EventBus } from './bus'
  import popup from './popup'

  export default {
    components: {
      popup
    },
    data () {
      return {
        alert: '',
        loggingIn: false,
        password: '',
        username: '',
        visible: false
      }
    },
    methods: {
      doLogin () {
        this.alert = ''
        this.loggingIn = true

        Axios.post('/login', {
          password: this.password,
          username: this.username
        }).then(() => {
          this.$root.checkUser()
          this.visible = false
          this.username = ''
          this.password = ''
        }).catch((error) => {
          if (error.response) {
            this.alert = error.response.data.error
          } else {
            this.alert = error.message
          }
        }).finally(() => {
          this.loggingIn = false
        })
      },
      show () {
        this.alert = ''
        this.visible = true
      }
    },
    created () {
      EventBus.$on('login', this.show)
    },
    beforeDestroy () {
      EventBus.$off('login', this.show)
    }
  }
</script>
