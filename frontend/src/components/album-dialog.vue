<template>
  <popup :title="title" @close="handleClosed" position="center">
    <form @submit.prevent="handleSubmit">
      <p class="alert" v-if="hasAlert">{{ alert }}</p>
      <label for="visibility">Visibility</label>
      <VisibilitySwitch :visibility="visibility" @change="n => (visibility = n)"></VisibilitySwitch>
      <label for="name">Name</label>
      <input id="name" placeholder="My Holiday" type="text" v-focus v-model="name" />
      <input :value="action" type="submit" />
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
  import VisibilitySwitch from "./visibility-switch.vue";
  import {computed, defineComponent, reactive, toRefs, watch} from "@vue/composition-api";
  import {useAlert} from "@/features/alert";

  export default defineComponent({
    components: {Popup, VisibilitySwitch},
    props: {
      id: String,
      initialName: String,
      initialVisibility: Number,
    },
    setup(props, ctx) {
      const {alert, hasAlert, setAlert} = useAlert();
      const state = reactive({
        name: props.initialName || "",
        visibility: props.initialVisibility || 0,
      });

      const isNew = computed(() => !props.id);
      const title = computed(() => (isNew.value ? "Create new album" : "Edit album"));
      const action = computed(() => (isNew.value ? "Create" : "Update"));

      function handleSubmit() {
        const url = isNew.value ? "/api/albums" : `/api/albums/${props.id}`;
        Axios.post(url, {name: state.name, visibility: state.visibility})
          .then(({data: {id}}) => {
            ctx.emit("created", id);
            state.name = "";
          })
          .catch(setAlert);
      }

      function handleClosed() {
        ctx.emit("close");
        state.name = "";
      }

      watch(
        () => props.initialName,
        () => {
          console.log("Initial name updated ", props.initialName, state.name);
          state.name = props.initialName || "";
        }
      );
      watch(
        () => props.initialVisibility,
        () => {
          state.visibility = props.initialVisibility || 0;
        }
      );

      return {handleSubmit, handleClosed, alert, hasAlert, title, action, ...toRefs(state)};
    },
  });
</script>
