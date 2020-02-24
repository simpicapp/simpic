<template>
    <header>
        <h1>Simpic</h1>
        <router-link to="/timeline/">Timeline</router-link>
        <router-link to="/albums/">Albums</router-link>
        <span class="grow"></span>
        <button v-on:click="loginClick" v-if="!$root.loggedIn">Login</button>
        <template v-else>
            <div>Logged in as {{$root.username}}</div>
            <button v-on:click="logoutClick">Logout</button>
        </template>
    </header>
</template>

<style scoped>
    header {
        display: flex;
        margin: 0 0 20px 0;
        padding: 10px;
        border-bottom: 1px solid #ccc;
        justify-content: space-between;
        align-items: center;
    }

    header h1 {
        text-transform: lowercase;
        margin-right: 20px;
    }

    span {
        flex-grow: 1;
    }

    header a {
        margin: 0 10px;
        padding: 5px;
        align-self: end;
        color: black;
        text-decoration: none;
    }

    a:hover {
        color: blue;
        text-decoration: underline;
    }

    /*noinspection CssUnusedSymbol*/
    a.router-link-active {
        font-weight: bold;
        text-decoration: none;
        color: black;
    }

    header * {
        margin: 0 10px;
    }
</style>

<script>
  import Axios from 'axios'
  import { EventBus } from './bus'

  export default {
    methods: {
      loginClick () {
        EventBus.$emit('login')
      },
      logoutClick () {
        Axios.get('/logout').then(() => {
          this.$root.loggedIn = false
          EventBus.$emit('toast', 'You have been logged out')
        })
      }
    }
  }
</script>
