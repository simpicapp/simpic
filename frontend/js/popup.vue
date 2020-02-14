<template>
    <aside class="popup">
        <h2>{{ title }}</h2>
        <a v-if="closeable" class="close" v-on:click="$emit('close')">❌</a>
        <div class="scroller">
            <slot></slot>
        </div>
    </aside>
</template>

<style>
    .popup {
        position: fixed;
        background: white;
        border: 2px solid black;
        border-radius: 5px;
        box-shadow: lightslategray 10px 10px;
        bottom: 20px;
        right: 20px;
    }

    .popup h2 {
        text-align: center;
        background-color: black;
        color: white;
        padding: 10px 5px;
        font-size: medium;
        margin: 0;
    }

    .popup .close {
        cursor: pointer;
        position: absolute;
        top: 0;
        right: 0;
        padding: 10px;
        display: inline-block;
    }

    .popup .scroller {
        max-height: 400px;
        overflow-y: auto;
    }
</style>

<script>
    export default {
        props: {
            title: String,
            closeable: {
                type: Boolean,
                default: true
            }
        }
    }
</script>
