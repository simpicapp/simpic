<template>
  <ul class="visibility" id="visibility">
    <li
      :class="{selected: visibility === 0}"
      @click="$emit('change', 0)"
      title="Can be seen by anyone, whether they're logged in or not"
    >
      <Icon name="globe-europe"></Icon>
      <span>Public</span>
    </li>

    <li
      :class="{selected: visibility === 1}"
      @click="$emit('change', 1)"
      title="Can only be seen by users who are logged in to Simpic, unless accessed directly by its ID"
    >
      <Icon name="link"></Icon>
      <span>Unlisted</span>
    </li>

    <li
      :class="{selected: visibility === 2}"
      @click="$emit('change', 2)"
      title="Can only be seen by users who are logged in to Simpic, with no exceptions"
    >
      <Icon name="lock"></Icon>
      <span>Private</span>
    </li>
  </ul>
</template>

<style lang="scss" scoped>
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
    user-select: none;

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
  import "vue-awesome/icons/globe-europe";
  import "vue-awesome/icons/link";
  import "vue-awesome/icons/lock";
  import Icon from "vue-awesome/components/Icon.vue";
  import {defineComponent} from "@vue/composition-api";

  export default defineComponent({
    props: {
      visibility: Number,
    },
    components: {
      Icon,
    },
  });
</script>
