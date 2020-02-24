<template>
    <header>
        <h1>Simpic</h1>
        <router-link to="/timeline/">Timeline</router-link>
        <router-link to="/albums/">Albums</router-link>
        <span class="grow"></span>
        <button v-on:click="loginClick" v-if="!$root.loggedIn">Login</button>
        <div v-else>Logged in as {{$root.username}}</div>
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

    h1 {
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
        color: gray;
        text-decoration: none;
    }

    header * {
        margin: 0;
    }
</style>

<script>
  import { EventBus } from './bus'

  export default {
    methods: {
      loginClick () {
        EventBus.$emit('login')
      }
    }
  }
</script>
