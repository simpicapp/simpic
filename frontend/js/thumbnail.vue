<template>
    <div class="thumbnail">
        <a v-bind:href="'/photo/' + photo.id" v-on:click.prevent="handleClick">
            <img v-bind:src="'/thumbnail/' + photo.id" v-bind:alt="photo.file_name">
            <div class="overlay">
                <p class="caption">{{ photo.file_name }}</p>
            </div>
        </a>
    </div>
</template>

<style>
    .thumbnail {
        position: relative;
        margin: 10px 5px;
    }

    .thumbnail img {
        max-width: 800px;
    }

    .thumbnail:hover .overlay {
        display: grid;
    }

    .thumbnail .overlay {
        display: none;
        position: absolute;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        grid-template-rows: 25% 50% 25%;
        grid-template-columns: 25% 50% 25%;
        overflow: hidden;
    }

    .thumbnail .caption {
        background-color: #000000cc;
        grid-area: 3 / 1 / 4 / 4;
        text-align: center;
        align-self: end;
        margin: 0;
        padding: 5px;
        color: white;
        font-weight: bold;
        overflow: hidden;
        min-width: 0;
    }
</style>

<script>
    export default {
        props: ['photo'],
        methods: {
            handleClick(e) {
                this.$router.push({name: 'lightbox', params: {'id': this.photo.id}});
            }
        }
    }
</script>
