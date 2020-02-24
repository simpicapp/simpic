<template>
    <aside v-bind:class="'popup ' + position">
        <h2>{{ title }}</h2>
        <a v-if="closeable" class="close" v-on:click="$emit('close')">❌</a>
        <div class="scroller">
            <slot></slot>
        </div>
    </aside>
</template>

<style scoped>
    aside {
        position: fixed;
        background: white;
        border: 2px solid black;
        border-radius: 5px;
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
        text-align: center;
        background-color: black;
        color: white;
        padding: 10px 5px;
        font-size: medium;
        margin: 0;
    }

    .close {
        cursor: pointer;
        position: absolute;
        top: 0;
        right: 0;
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
      position: {
        type: String,
        default: 'bottom_right',
        validator: value => ['bottom_right', 'center'].indexOf(value) !== -1
      },
      title: String
    }
  }
</script>
