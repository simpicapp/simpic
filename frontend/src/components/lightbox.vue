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
  import {Photo} from "@/model/photo";
  import {useWindowListener} from "@/features/listeners";
  import {defineComponent, reactive, ref, toRefs, watch} from "@vue/composition-api";

  export default defineComponent({
    components: {Modal},
    props: {id: String},

    setup(props, ctx) {
      const canvas = ref(null as HTMLCanvasElement | null);
      const state = reactive({
        close: false,
        width: 0,
        height: 0,
        metadata: null as Photo | null,
      });

      function setSize() {
        if (!state.metadata) {
          return;
        }

        const widthRatio = state.metadata.width / (window.innerWidth * 0.95);
        const heightRatio = state.metadata.height / (window.innerHeight * 0.9);
        const scale = Math.max(1, widthRatio, heightRatio);

        state.width = Math.round(state.metadata.width / scale);
        state.height = Math.round(state.metadata.height / scale);
      }

      function context() {
        return canvas.value && canvas.value.getContext("2d");
      }

      function startLoading() {
        const id = props.id;
        if (!id) {
          return;
        }

        cache
          .getMetadata(id)
          .then(metadata => {
            if (props.id !== id) {
              throw Error("wrong-id");
            }
            state.metadata = metadata;
            setSize();
            return cache.getThumbnail(id);
          })
          .then(img => {
            if (props.id !== id) {
              throw Error("wrong-id");
            }
            const ctx = context();
            if (ctx) {
              ctx.filter = "blur(4px)";
              ctx.clearRect(0, 0, state.width, state.height);
              ctx.drawImage(img, 0, 0, state.width, state.height);
            }
            return cache.getImage(id);
          })
          .then(img => {
            if (props.id !== id) {
              throw Error("wrong-id");
            }
            const ctx = context();
            if (ctx) {
              ctx.filter = "none";
              ctx.drawImage(img, 0, 0, state.width, state.height);
            }
          })
          .catch(err => {
            if (err.message !== "wrong-id") {
              console.log(err);
            }
          });
      }

      useWindowListener("keyup", (event: KeyboardEvent) => {
        if (event.code === "Escape") {
          state.close = true;
        } else if (event.code === "ArrowLeft") {
          ctx.emit("go-to-previous-image", props.id);
        } else if (event.code === "ArrowRight") {
          ctx.emit("go-to-next-image", props.id);
        }
      });

      useWindowListener("resize", () => {
        setSize();
        startLoading();
      });

      watch(() => props.id, startLoading);

      return {canvas, ...toRefs(state)};
    },
  });
</script>
