<template>
  <modal :should-close="close" @close="visible = false" v-if="visible">
    <popup @close="close = true" position="center" title="Login">
      <form @submit="doLogin">
        <p class="alert" v-if="alert.length > 0">{{ alert }}</p>
        <label for="username">Username</label>
        <input :disabled="loggingIn" id="username" type="text" v-focus v-model="username" />
        <label for="password">Password</label>
        <input :disabled="loggingIn" id="password" type="password" v-model="password" />
        <input :disabled="loggingIn" type="submit" value="Login" />
      </form>
    </popup>
  </modal>
</template>

<style lang="scss" scoped>
  form {
    display: grid;
    grid-template-columns: auto auto;
    grid-gap: 30px 20px;
    align-items: center;
  }

  input[type="submit"] {
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

<script lang="ts">
  import Popup from "./popup.vue";
  import Modal from "./modal.vue";
  import {defineComponent, reactive, toRefs} from "@vue/composition-api";
  import {useAuthentication} from "@/features/auth";
  import {useEventListener} from "@/features/eventbus";

  export default defineComponent({
    components: {
      Modal,
      Popup,
    },
    setup() {
      const {login} = useAuthentication();

      const state = reactive({
        alert: "",
        close: false,
        loggingIn: false,
        password: "",
        username: "",
        visible: false,
      });

      function doLogin() {
        state.alert = "";
        state.loggingIn = true;

        login(state.username, state.password)
          .then(() => {
            state.close = true;
            state.username = "";
            state.password = "";
          })
          .catch(error => {
            if (error.response) {
              state.alert = error.response.data.error;
            } else {
              state.alert = error.message;
            }
          })
          .finally(() => {
            state.loggingIn = false;
          });
      }

      useEventListener("login", () => {
        state.close = false;
        state.alert = "";
        state.visible = true;
      });

      return {doLogin, ...toRefs(state)};
    },
  });
</script>
