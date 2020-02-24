<template>
    <div>
        <aside :class="'popup ' + position">
            <h2>{{ title }}</h2>
            <a v-if="closeable" class="close" @click="$emit('close')">&times;</a>
            <div class="scroller">
                <slot></slot>
            </div>
        </aside>
        <transition name="background" appear>
            <div class="background" v-if="modal" @click="handleBackgroundClick"></div>
        </transition>
    </div>
</template>

<style scoped>
    .background {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: 900;
    }

    @supports (backdrop-filter: blur(5px)) {
        .background {
            backdrop-filter: blur(5px) grayscale(80%);
        }

        .background-enter-active, .background-leave-active {
            transition: backdrop-filter 100ms linear;
        }

        .background-enter, .background-leave-to {
            backdrop-filter: blur(0) grayscale(0);
        }
    }

    @supports not (backdrop-filter: blur(5px)) {
        .background {
            background-color: #66666699;
        }

        .background-enter-active, .background-leave-active {
            transition: background-color 100ms linear;
        }

        .background-enter, .background-leave-to {
            background-color: #66666600;
        }
    }

    aside {
        padding: 30px;
        position: fixed;
        background: white;
        border-radius: 2px;
        box-shadow: 10px 8px 8px #00000066;
        z-index: 1000;
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
        color: var(--smaragdine);
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
    }
</style>

<script>
  export default {
    props: {
      closeable: {
        type: Boolean,
        default: true
      },
      modal: {
        type: Boolean,
        default: false
      },
      position: {
        type: String,
        default: 'bottom_right',
        validator: value => ['bottom_right', 'center'].indexOf(value) !== -1
      },
      title: String
    },
    methods: {
      handleBackgroundClick () {
        if (this.closeable) {
          this.$emit('close')
        }
      }
    }
  }
</script>
