<template>
    <header>
        <h1>Simpic</h1>
        <button v-on:click="loginClick" v-if="!$root.loggedIn">Login</button>
        <div v-else>Logged in as {{$root.username}}</div>
    </header>
</template>

<style>
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
        flex-grow: 1;
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
