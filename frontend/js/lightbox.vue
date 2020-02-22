<template>
    <div id="lightbox" v-on:click="close()">
        <img v-bind:src="'/data/image/' + id" v-on:click.stop>
    </div>
</template>

<style>
    #lightbox {
        z-index: 900;
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
  export default {
    props: ['id'],
    methods: {
      close () {
        this.$router.push('../')
      },
      handleKey (event) {
        if (event.code === 'Escape') {
          this.close()
        }
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
