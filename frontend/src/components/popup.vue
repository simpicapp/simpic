<template>
  <aside :class="'popup ' + position">
    <h2>{{ title }}</h2>
    <a @click="$emit('close')" class="close" v-if="closeable">&times;</a>
    <div class="scroller">
      <slot></slot>
    </div>
  </aside>
</template>

<style lang="scss" scoped>
  @use 'src/assets/css/vars';

  aside {
    padding: 30px;
    position: fixed;
    background: white;
    border-radius: 2px;
    box-shadow: 10px 8px 8px #00000066;
    z-index: 1000;
    transform: translateZ(0);
  }

  /*noinspection CssUnusedSymbol*/
  .bottom_right {
    bottom: 20px;
    right: 20px;
  }

  /*noinspection CssUnusedSymbol*/
  .center {
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
  }

  h2 {
    color: vars.$primary;
    padding: 0;
    font-size: large;
    margin: 0 0 30px 0;
  }

  .close {
    cursor: pointer;
    position: absolute;
    top: 20px;
    right: 20px;
    padding: 10px;
    display: inline-block;
  }

  .scroller {
    max-height: 400px;
    overflow-y: auto;

    /* Edges of buttons and form fields sometimes get cut off without this... */
    margin-right: 1px;
  }
</style>

<script lang="ts">
  import {defineComponent} from "@vue/composition-api";

  export default defineComponent({
    props: {
      closeable: {
        type: Boolean,
        default: true,
      },
      position: {
        type: String,
        default: "bottom_right",
        validator: value => ["bottom_right", "center"].indexOf(value) !== -1,
      },
      title: String,
    },
  });
</script>
