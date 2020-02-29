<template>
    <transition name="fade" appear @after-enter="transitionFinished = true">
        <div class="background" :class="{darker}" @click="handleBackgroundClick">
            <transition name="fade" @after-leave="$emit('close')">
                <div @click.stop v-show="showContent">
                    <slot></slot>
                </div>
            </transition>
        </div>
    </transition>
</template>

<style lang="scss" scoped>
    .fade-enter, .fade-leave-to {
        opacity: 0;
    }

    .fade-enter-active, .fade-leave-active {
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

<script>
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
  export default {
    props: {
      closeable: {
        type: Boolean,
        default: true
      },
      darker: {
        type: Boolean,
        default: false
      },
      shouldClose: {
        type: Boolean,
        default: false
      }
    },
    data () {
      return {
        closing: false,
        transitionFinished: false
      }
    },
    methods: {
      handleBackgroundClick () {
        if (this.closeable) {
          this.closing = true
        }
      },
      handleKey (event) {
        if (event.code === 'Escape') {
          this.closing = false
        }
      }
    },
    computed: {
      showContent () {
        return !this.shouldClose && !this.closing && this.transitionFinished
      }
    },
    mounted () {
      window.addEventListener('keyup', this.handleKey)
    },
    destroyed () {
      window.removeEventListener('keyup', this.handleKey)
    }
  }
</script>
