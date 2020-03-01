<template>
  <popup @close="handleClosed" position="center" title="Create new album">
    <form @submit="doCreate">
      <p class="alert" v-if="hasAlert">{{ alert }}</p>
      <label for="name">Name</label>
      <input id="name" placeholder="My Holiday" type="text" v-focus v-model="name" />
      <input type="submit" value="Create" />
    </form>
  </popup>
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
  import Axios from "axios";
  import Popup from "./popup.vue";
  import {defineComponent, ref} from "@vue/composition-api";
  import {useAlert} from "@/features/alert";

  export default defineComponent({
    components: {Popup},
    setup(_, ctx) {
      const {alert, hasAlert, setAlert} = useAlert();
      const name = ref("");

      function doCreate() {
        Axios.post("/albums", {name: name.value})
          .then(({data: {id}}) => {
            ctx.emit("created", id);
            name.value = "";
          })
          .catch(error => setAlert(error));
      }

      function handleClosed() {
        ctx.emit("close");
        name.value = "";
      }

      return {doCreate, handleClosed, alert, hasAlert, name};
    },
  });
</script>
