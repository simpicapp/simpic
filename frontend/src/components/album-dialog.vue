<template>
  <popup :title="title" @close="handleClosed" position="center">
    <form @submit="handleSubmit">
      <p class="alert" v-if="hasAlert">{{ alert }}</p>
      <label for="visibility">Visibility</label>
      <ul class="visibility" id="visibility">
        <li :class="{selected: visibility === 0}" @click="visibility = 0">
          <Icon name="globe-europe"></Icon>
          <span>Public</span>
        </li>
        <li :class="{selected: visibility === 1}" @click="visibility = 1">
          <Icon name="link"></Icon>
          <span>Unlisted</span>
        </li>
        <li :class="{selected: visibility === 2}" @click="visibility = 2">
          <Icon name="lock"></Icon>
          <span>Private</span>
        </li>
      </ul>
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

  .visibility {
    display: flex;
    justify-content: space-between;
    border-radius: 5px;
    background-color: #eeeeee;
  }

  .visibility li {
    padding: 5px 10px;
    border-left: 1px solid #dddddd;
    display: flex;
    align-items: center;
    transition: all 200ms linear;
    cursor: pointer;

    &:first-child {
      border-bottom-left-radius: 5px;
      border-top-left-radius: 5px;
      border-left: 0;
    }

    &:last-child {
      border-bottom-right-radius: 5px;
      border-top-right-radius: 5px;
    }

    &.selected,
    &:hover {
      background-color: #333333;
      color: white;
    }

    span {
      margin-left: 5px;
    }
  }
</style>

<script lang="ts">
  import Axios from "axios";
  import Popup from "./popup.vue";
  import {computed, defineComponent, reactive, toRefs, watch} from "@vue/composition-api";
  import {useAlert} from "@/features/alert";
  import "vue-awesome/icons/globe-europe";
  import "vue-awesome/icons/link";
  import "vue-awesome/icons/lock";
  import Icon from "vue-awesome/components/Icon.vue";

  export default defineComponent({
    components: {Icon, Popup},
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
