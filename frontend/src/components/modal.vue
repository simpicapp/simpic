<template>
  <transition @after-enter="transitionFinished = true" appear name="fade">
    <div :class="{darker}" @click="handleBackgroundClick" class="background">
      <transition @after-leave="$emit('close')" name="fade">
        <div @click.stop v-show="showContent">
          <slot></slot>
        </div>
      </transition>
    </div>
  </transition>
</template>

<style lang="scss" scoped>
  .fade-enter,
  .fade-leave-to {
    opacity: 0;
  }

  .fade-enter-active,
  .fade-leave-active {
    transition: 200ms opacity ease-out;
  }

  .background {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 900;
    transform: translateZ(0);
    background-color: #cccccc99;

    &.darker {
      background-color: #000000ee;
    }
  }
</style>

<script lang="ts">
  import {computed, defineComponent, reactive, toRefs} from "@vue/composition-api";
  import {useWindowListener} from "@/features/listeners";

  // Modal has a double stage animation: background, then content.
  //
  // The state here is a bit complex:
  //
  //  shouldClose is a prop used by the parent element to tell us we should begin closing
  //
  //  closing is an internal data value used when we decide we should begin closing ourselves (i.e., escape/click)
  //
  //  showContent is a computed value over those that enables the actual modal content only if we're not closing
  //  and the outer transition has finished.

  export default defineComponent({
    props: {
      closeable: Boolean,
      darker: Boolean,
      shouldClose: Boolean,
    },
    setup(props) {
      const state = reactive({
        closing: false,
        transitionFinished: false,
      });

      function handleBackgroundClick() {
        if (props.closeable) {
          state.closing = true;
        }
      }

      useWindowListener("keyup", event => {
        if (props.closeable && event.code === "Escape") {
          state.closing = false;
        }
      });

      const showContent = computed(() => !props.shouldClose && !state.closing && state.transitionFinished);

      return {handleBackgroundClick, showContent, ...toRefs(state)};
    },
  });
</script>
