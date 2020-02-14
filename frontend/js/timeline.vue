<template>
    <div class="timeline">
        <p v-if="loading">Loading...</p>
        <img v-for="photo in photos" v-bind:key="photo.id" v-bind:src="'/thumbnail/' + photo.id"
             v-bind:alt="photo.file_name">
    </div>
</template>

<style>

</style>

<script>
    import {EventBus} from './bus';

    export default {
        data: function () {
            return {
                loading: true,
                photos: []
            }
        },
        methods: {
            update() {
                const comp = this;
                fetch('/timeline')
                    .then((response) => response.json())
                    .then((json) => comp.photos = json)
                    .then(() => comp.loading = false)
            }
        },
        mounted: function () {
            this.update();
            EventBus.$on('upload-complete', this.update);
        }
    }
</script>
