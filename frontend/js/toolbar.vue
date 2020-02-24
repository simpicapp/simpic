<template>
    <header>
        <h1>Simpic</h1>
        <router-link to="/timeline/">Timeline</router-link>
        <router-link to="/albums/">Albums</router-link>
        <span class="grow"></span>
        <button @click="loginClick" v-if="!$root.loggedIn">Login</button>
        <template v-else>
            <div>Logged in as {{$root.username}}</div>
            <button @click="logoutClick">Logout</button>
        </template>
    </header>
</template>

<style lang="scss" scoped>
    header {
        display: flex;
        margin: 0 0 20px 0;
        padding: 10px;
        border-bottom: 1px solid #ccc;
        justify-content: space-between;
        align-items: center;

        h1 {
            text-transform: lowercase;
            margin-right: 20px;
        }

        span {
            flex-grow: 1;
        }

        a {
            margin: 0 10px;
            padding: 5px;
            align-self: end;
            color: black;
            text-decoration: none;

            &:hover {
                color: blue;
                text-decoration: underline;
            }

            /*noinspection CssUnusedSymbol*/
            &.router-link-active {
                font-weight: bold;
                text-decoration: none;
                color: black;
            }
        }

        * {
            margin: 0 10px;
        }
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
