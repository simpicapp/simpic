<template>
    <modal @close="$router.push('../')" :should-close="close" :darker="true">
        <div id="lightbox" @click="close = true">
            <div id="prev-overlay" @click.stop.prevent="$emit('go-to-previous-image')">
                <span>←</span>
            </div>
            <div id="close">&times; Close</div>
            <img :src="'/data/image/' + id" @click.stop>
            <div id="next-overlay" @click.stop.prevent="$emit('go-to-next-image')">
                <span>→</span>
            </div>
        </div>
    </modal>
</template>

<style lang="scss" scoped>
    #lightbox {
        z-index: 950;
        position: fixed;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        overscroll-behavior: contain;
        flex-direction: column;
    }

    img {
        max-width: 95%;
        max-height: 90%;
    }

    #close {
        position: fixed;
        top: 10px;
        right: 220px;
        color: #999;
        padding: 10px;
        cursor: pointer;

        &:hover {
            color: white;
        }
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
        color: #999;
        cursor: pointer;

        &:hover {
            background: #ffffff33;
            color: white;
        }
    }

    #next-overlay {
        right: 0;
    }

    #prev-overlay {
        left: 0;
    }
</style>

<script>
  import Modal from './modal'

  export default {
    components: { Modal },
    props: ['id'],
    data () {
      return {
        close: false
      }
    },
    methods: {
      handleKey (event) {
        if (event.code === 'Escape') {
          this.close = true
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
