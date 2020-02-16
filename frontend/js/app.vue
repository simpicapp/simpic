<template>
    <main>
        <div id="drop-target" v-if="dragging">
            Drop files here
        </div>
        <timeline></timeline>
        <uploader></uploader>
        <login></login>
        <lightbox></lightbox>
    </main>
</template>

<style>
    body {
        font-family: sans-serif;
    }

    #drop-target {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: lightsteelblue;
        border: 20px dashed midnightblue;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: xx-large;
    }
</style>

<script>
    import lightbox from "./lightbox";
    import timeline from "./timeline";
    import uploader from "./uploader";
    import login from "./login";
    import { EventBus } from './bus';

    export default {
        components: {
            timeline,
            uploader,
            lightbox,
            login
        },
        data() {
            return {
                dragging: false
            }
        },
        methods: {
            dropHandler(e) {
                this.dragging = false;
                e.stopPropagation();
                e.preventDefault();
                EventBus.$emit('files-dropped', e.dataTransfer.files);
            },
            dragOverHandler(e) {
                this.dragging = true;
                e.stopPropagation();
                e.preventDefault();
                e.dataTransfer.dropEffect = 'copy';
            },
            dragStartHandler(e) {
                this.dragging = true;
                e.stopPropagation();
                e.preventDefault();
            },
            dragEndHandler(e) {
                this.dragging = false;
                e.stopPropagation();
                e.preventDefault();
            }
        },
        mounted() {
            document.addEventListener('drop', this.dropHandler);
            document.addEventListener('dragover', this.dragOverHandler);
            document.addEventListener('dragenter', this.dragStartHandler);
            document.addEventListener('dragleave', this.dragEndHandler);
        },
        beforeDestroy() {
            document.removeEventListener('drop', this.dropHandler);
            document.removeEventListener('dragover', this.dragOverHandler);
            document.removeEventListener('dragenter', this.dragStartHandler);
            document.removeEventListener('dragleave', this.dragEndHandler);
        }
    }
</script>
