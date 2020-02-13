<template>
    <main>
        <div id="drop-target" v-if="dragging">
            Drop files here
        </div>
        <timeline/>
    </main>
</template>

<style>
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
    import timeline from "./timeline";

    export default {
        components: {
            timeline
        },
        data() {
            return {
                dragging: false
            }
        },
        methods: {
            dropHandler: function (e) {
                this.dragging = false;
                e.stopPropagation();
                e.preventDefault();
            },
            dragOverHandler: function (e) {
                this.dragging = true;
                e.stopPropagation();
                e.preventDefault();
                e.dataTransfer.dropEffect = 'copy';
            },
            dragStartHandler: function (e) {
                this.dragging = true;
                e.stopPropagation();
                e.preventDefault();
            },
            dragEndHandler: function (e) {
                this.dragging = false;
                e.stopPropagation();
                e.preventDefault();
            }
        },
        mounted() {
            document.addEventListener('drop', this.dropHandler);
            document.addEventListener('dragover', this.dragOverHandler);
            document.addEventListener('dragstart', this.dragStartHandler);
            document.addEventListener('dragenter', this.dragStartHandler);
            document.addEventListener('dragend', this.dragEndHandler);
            document.addEventListener('dragleave', this.dragEndHandler);
        },
        beforeDestroy() {
            document.removeEventListener('drop', this.dropHandler);
            document.removeEventListener('dragover', this.dragOverHandler);
            document.removeEventListener('dragstart', this.dragStartHandler);
            document.removeEventListener('dragenter', this.dragStartHandler);
            document.removeEventListener('dragend', this.dragEndHandler);
            document.removeEventListener('dragleave', this.dragEndHandler);
        }
    }
</script>
