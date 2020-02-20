<template>
    <popup title="Login" position="center" v-if="visible" v-on:close="visible = false">
        <form v-on:submit="doLogin">
            <p class="alert" v-if="alert.length > 0">{{ alert }}</p>
            <label for="username">Username</label>
            <input type="text" id="username" v-model="username" v-bind:disabled="loggingIn">
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
        grid-gap: 20px 20px;
        padding: 20px;
        align-items: center;
    }

    input[type=submit] {
        padding: 5px 0;
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
        fetch('/login', {
          body: JSON.stringify({
            password: this.password,
            username: this.username
          }),
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
          },
          method: 'POST'
        }).then((response) => {
          if (response.ok) {
            this.$root.checkUser()
            this.visible = false
          } else {
            return response.json().then((json) => {
              throw new Error(json.error)
            })
          }
        }).catch((error) => {
          this.alert = error.message
        }).finally(() => {
          this.loggingIn = false
        })
      },
      show () {
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
