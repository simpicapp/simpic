<template>
  <transition appear name="boing">
    <aside @click="hide" v-if="visible">
      {{message}}
    </aside>
  </transition>
</template>

<style lang="scss" scoped>
  @use 'src/assets/css/vars';

  aside {
    position: fixed;
    bottom: 100px;
    left: 30%;
    right: 30%;
    background-color: vars.$primary;
    color: white;
    font-size: large;
    text-align: center;
    padding: 30px;
    border-radius: 2px;
    box-shadow: 10px 8px 8px #00000066;
    z-index: 2000;
    transform: translateY(0);
  }

  .boing-enter-active {
    transition: all 200ms var(--ease-out-back);
  }

  .boing-leave-active {
    transition: all 100ms linear;
  }

  .boing-enter, .boing-leave-to {
    transform: translateY(150px)
  }
</style>

<script lang="ts">
  import {EventBus} from './bus'
  import Vue from 'vue'

  export default Vue.extend({
    data() {
      return {
        message: '',
        visible: false
      }
    },
    methods: {
      hide() {
        this.visible = false
      },
      showToast(toast: string) {
        this.message = toast;
        this.visible = true;
        setTimeout(this.hide, 3500)
      }
    },
    mounted() {
      EventBus.$on('toast', this.showToast)
    },
    destroyed() {
      EventBus.$off('toast', this.showToast)
    }
  })
</script>
