<template>
    <div id="lightbox" v-on:click="close()">
        <div id="prev-overlay" v-on:click.stop.prevent="$emit('go-to-previous-image')">
            <span>←</span>
        </div>
        <img v-bind:src="'/data/image/' + id" v-on:click.stop>
        <div id="next-overlay" v-on:click.stop.prevent="$emit('go-to-next-image')">
            <span>→</span>
        </div>
    </div>
</template>

<style scoped>
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

    img {
        max-width: 95%;
        max-height: 95%;
    }

    #next-overlay, #prev-overlay {
        position: fixed;
        top: 0;
        bottom: 0;
        width: 200px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: xx-large;
        color: #eeeeee;
        cursor: pointer;
    }

    #next-overlay:hover, #prev-overlay:hover {
        background: #ffffff33;
    }

    #next-overlay {
        right: 0;
    }

    #prev-overlay {
        left: 0;
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
        } else if (event.code === 'ArrowLeft') {
          this.$emit('go-to-previous-image')
        } else if (event.code === 'ArrowRight') {
          this.$emit('go-to-next-image')
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
