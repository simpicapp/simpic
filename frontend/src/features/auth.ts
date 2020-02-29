import {reactive, toRefs} from "@vue/composition-api";
import Axios from "axios";
import {EventBus} from "@/components/bus";

const state = reactive({
  username: '',
  loggedIn: false
});

export function useAuthentication() {
  function checkUser() {
    return Axios.get('/users/me').then(({data}) => {
      state.username = data.username;
      state.loggedIn = true;
    }).catch(() => {
      state.loggedIn = false;
    })
  }

  function logout() {
    return Axios.get('/logout').then(() => {
      state.loggedIn = false;
      EventBus.$emit('toast', 'You have been logged out')
    })
  }

  function login(username: string, password: string) {
    return Axios.post('/login', {
      username: username,
      password: password
    }).then(() => {
      EventBus.$emit('toast', 'You are now logged in');
      return checkUser();
    })
  }

  return {
    ...toRefs(state),
    login,
    logout,
    checkUser
  };
}
