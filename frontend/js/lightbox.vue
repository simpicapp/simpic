<template>
    <div id="lightbox" v-if="visible" v-on:click="visible = false">
        <img v-bind:src="'/photo/' + photo.id" v-on:click.stop>
    </div>
</template>

<style>
    #lightbox {
        position: fixed;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
        background-color: #000000ee;
        display: flex;
        align-items: center;
        justify-content: center;
        overscroll-behavior: contain;
    }

    #lightbox img {
        max-width: 95%;
        max-height: 95%;
    }
</style>

<script>
    import {EventBus} from "./bus";

    export default {
        data() {
            return {
                visible: false,
                photo: {}
            }
        },
        methods: {
            showPhoto(photo) {
                this.visible = true;
                this.photo = photo;
            }
        },
        mounted() {
            EventBus.$on('show-photo', this.showPhoto)
        },
        beforeDestroy() {
            EventBus.$off('show-photo', this.showPhoto)
        }
    }
</script>
