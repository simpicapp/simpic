import {reactive, toRefs} from "@vue/composition-api";
import Axios from "axios";
import {EventBus} from "@/features/eventbus";

const state = reactive({
  username: "",
  loggedIn: false,
  wasLoggedIn: undefined as boolean | undefined,
});

export function useAuthentication() {
  function checkUserChanged() {
    if (state.loggedIn !== state.wasLoggedIn) {
      EventBus.$emit("user-changed");
      state.wasLoggedIn = state.loggedIn;
    }
  }

  function checkUser() {
    return Axios.get("/users/me")
      .then(({data}) => {
        state.username = data.username;
        state.loggedIn = true;
        checkUserChanged();
      })
      .catch(() => {
        state.loggedIn = false;
        checkUserChanged();
      });
  }

  function logout() {
    return Axios.get("/logout").then(() => {
      state.loggedIn = false;
      EventBus.$emit("toast", "You have been logged out");
      checkUserChanged();
    });
  }

  function login(username: string, password: string) {
    return Axios.post("/login", {
      username: username,
      password: password,
    }).then(() => {
      EventBus.$emit("toast", "You are now logged in");
      return checkUser();
    });
  }

  return {
    ...toRefs(state),
    login,
    logout,
    checkUser,
  };
}
