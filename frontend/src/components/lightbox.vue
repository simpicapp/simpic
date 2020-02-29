<template>
  <modal :darker="true" :should-close="close" @close="$router.push('../')">
    <div @click="close = true" id="lightbox" ref="container">
      <div @click.stop.prevent="$emit('go-to-previous-image', id)" id="prev-overlay">
        <span>←</span>
      </div>
      <div id="close">&times; Close</div>
      <canvas :height="height" :width="width" @click.stop ref="canvas"></canvas>
      <div @click.stop.prevent="$emit('go-to-next-image', id)" id="next-overlay">
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

  #next-overlay,
  #prev-overlay {
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
    user-select: none;

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

<script lang="ts">
  import Modal from "./modal.vue";
  import {cache} from "./cache";
  import Vue from "vue";
  import {Data} from "@vue/composition-api/dist/component";
  import {Photo} from "@/model/photo";

  interface State extends Data {
    close: boolean;
    width: number;
    height: number;
    metadata: Photo;
  }

  export default Vue.extend({
    components: {Modal},
    props: {id: String},

    data(): State {
      return {
        close: false,
        height: 0,
        metadata: {
          id: "",
          file_name: "", // eslint-disable-line @typescript-eslint/camelcase
          width: 0,
          height: 0,
        },
        width: 0,
      };
    },

    methods: {
      handleKey(event: KeyboardEvent) {
        if (event.code === "Escape") {
          this.close = true;
        } else if (event.code === "ArrowLeft") {
          this.$emit("go-to-previous-image", this.id);
        } else if (event.code === "ArrowRight") {
          this.$emit("go-to-next-image", this.id);
        }
      },

      handleResize() {
        this.setSize();
        this.$nextTick(this.startLoading);
      },

      setSize() {
        const widthRatio = this.metadata.width / (window.innerWidth * 0.95);
        const heightRatio = this.metadata.height / (window.innerHeight * 0.9);
        const scale = Math.max(1, widthRatio, heightRatio);

        this.width = Math.round(this.metadata.width / scale);
        this.height = Math.round(this.metadata.height / scale);
      },

      context() {
        return (this.$refs.canvas as HTMLCanvasElement).getContext("2d");
      },

      startLoading() {
        const id = this.id;
        cache
          .getMetadata(this.id)
          .then(metadata => {
            if (this.id !== id) {
              throw Error("wrong-id");
            }
            this.metadata = metadata;
            this.setSize();
            return cache.getThumbnail(this.id);
          })
          .then(img => {
            if (this.id !== id) {
              throw Error("wrong-id");
            }
            const ctx = this.context();
            if (ctx) {
              ctx.filter = "blur(4px)";
              ctx.clearRect(0, 0, this.width, this.height);
              ctx.drawImage(img, 0, 0, this.width, this.height);
            }
            return cache.getImage(this.id);
          })
          .then(img => {
            if (this.id !== id) {
              throw Error("wrong-id");
            }
            const ctx = this.context();
            if (ctx) {
              ctx.filter = "none";
              ctx.drawImage(img, 0, 0, this.width, this.height);
            }
          })
          .catch(err => {
            if (err.message !== "wrong-id") {
              console.log(err);
            }
          });
      },
    },

    watch: {
      id() {
        this.startLoading();
      },
    },

    mounted() {
      window.addEventListener("keyup", this.handleKey);
      window.addEventListener("resize", this.handleResize);
      this.$nextTick(this.startLoading);
    },

    destroyed() {
      window.removeEventListener("keyup", this.handleKey);
      window.removeEventListener("resize", this.handleResize);
    },
  });
</script>
