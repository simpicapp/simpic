<template>
    <a v-bind:href="'/data/photo/' + id" v-on:click.prevent="handleClick">
        <img v-bind:src="'/data/thumb/' + cover" v-bind:alt="caption" v-if="cover != null">
        <p class="caption">{{ caption }}</p>
    </a>
</template>

<style scoped>
    a {
        position: relative;
        margin: 10px 15px;
        border: 1px solid #333333;
        width: 240px;
        height: 250px;
        display: flex;
        flex-direction: column;
        text-decoration: none;
        border-radius: 10px;
        justify-content: space-evenly;
        align-items: center;
        background-color: #E6D1AD;
    }

    a:hover {
        background-color: #f6e1bb;
    }

    img {
        max-width: 200px;
        max-height: 200px;
        border: 1px solid black;
    }

    p {
        margin: 0;
        color: black;
        font-size: large;
        font-weight: bold;
    }
</style>

<script>
  export default {
    props: ['id', 'caption', 'cover'],
    methods: {
      handleClick () {
        this.$router.push({ path: './' + this.id + '/' })
      }
    }
  }
</script>
