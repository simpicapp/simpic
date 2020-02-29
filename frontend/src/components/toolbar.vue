<template>
  <header>
    <h1><img alt="Simpic" src="../assets/logo.png"></h1>
    <router-link to="/timeline/"><span>Timeline</span></router-link>
    <router-link to="/albums/"><span>Albums</span></router-link>
    <span class="grow"></span>
    <button @click="loginClick" v-if="!loggedIn">Login</button>
    <template v-else>
      <div class="username">Logged in as {{username}}</div>
      <button @click="logout">Logout</button>
    </template>
  </header>
</template>

<style lang="scss" scoped>
  @use 'src/assets/css/vars';

  header {
    display: flex;
    margin: 0 0 20px 0;
    background-color: vars.$primary;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 1px 2px black;
    height: 70px;
    overflow: visible;

    h1 {
      text-transform: lowercase;
      margin: 10px 20px 0 10px;
      padding: 0;
    }

    img {
      margin-top: 28px;
      min-height: 66px;
      height: 66px;
    }

    .grow {
      flex-grow: 1;
    }

    a {
      align-self: stretch;
      margin: 0 10px;
      padding: 5px;
      width: 125px;
      color: white;
      text-decoration: none;
      font-size: large;
      display: flex;
      align-items: flex-end;
      justify-content: center;
      border-left: 1px solid darken(vars.$primary, 10%);
      border-right: 1px solid darken(vars.$primary, 10%);
      background-color: vars.$primary;
      transition: background-color 200ms ease-in;

      &.router-link-active, &:hover {
        text-decoration: none;
        color: white;
        background: lighten(vars.$primary, 5%);
      }
    }

    .username {
      color: white;
    }

    button {
      border: 1px solid white;
      border-radius: 2px;
      color: white;
      padding: 5px 20px;
      cursor: pointer;
      margin-right: 20px;
      background-color: vars.$primary;
      transition: background-color 200ms ease-in;

      &:hover, &:active {
        background: lighten(vars.$primary, 5%);
      }
    }

    * {
      margin: 0 10px;
    }
  }
</style>

<script>
  import { EventBus } from './bus'
  import { defineComponent } from '@vue/composition-api'
  import { useAuthentication } from '@/features/auth'

  export default defineComponent({
    methods: {
      loginClick () {
        EventBus.$emit('login')
      },
    },
    setup () {
      const { loggedIn, username, logout } = useAuthentication()

      return {
        loggedIn,
        username,
        logout
      }
    }
  })
</script>
