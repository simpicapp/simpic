<template>
  <div :class="{selecting}" :style="styles" class="thumbnail">
    <a :href="'/data/image/' + imageId" @click.prevent="handleClick">
      <div class="overlay">
        <p class="caption">{{ caption }}</p>
      </div>
      <span :class="{selected}" @click.prevent.stop="handleToggle" class="tickbox" role="button" v-if="loggedIn">
        {{ selected ? "☑" : "☐" }}
      </span>
    </a>
  </div>
</template>

<style lang="scss" scoped>
  .thumbnail {
    flex-grow: 1;
    flex-shrink: 1;
    position: relative;
    height: 200px;
    margin: 3px 3px;
    background-repeat: no-repeat;
    background-position: 50%;
    background-size: cover;
  }

  a {
    display: block;
    height: 200px;
  }

  img {
    max-width: 800px;
  }

  .overlay {
    display: grid;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    grid-template-rows: 25% auto 25%;
    grid-template-columns: 25% auto 25%;
    overflow: hidden;
    opacity: 0;
    transition: opacity 300ms var(--ease-in-cubic);
  }

  .tickbox {
    position: absolute;
    top: 0;
    left: 0;
    width: 50px;
    height: 50px;
    display: flex;
    background-color: #ffffffdd;
    font-size: xx-large;
    align-items: center;
    justify-content: center;
    border-bottom-right-radius: 5px;
    opacity: 0;
    transition: opacity 300ms var(--ease-in-cubic);

    color: #00000066;

    &.selected {
      color: black;
    }
  }

  .thumbnail:hover .overlay,
  .thumbnail:hover .tickbox,
  .thumbnail.selecting .tickbox {
    opacity: 1;
    transition: opacity 300ms var(--ease-out-cubic);
  }

  .caption {
    grid-area: 3 / 1 / 4 / 4;
    text-align: center;
    align-self: end;
    margin: 0;
    padding: 10px 0;
    color: white;
    overflow: hidden;
    min-width: 0;

    @supports (backdrop-filter: blur()) {
      backdrop-filter: blur(10px);
      background-color: #00000099;
    }

    @supports not (backdrop-filter: blur()) {
      background-color: #000000cc;
    }
  }
</style>

<script lang="ts">
  import ThumbnailBackground from "./thumbnail-background.vue";
  import {defineComponent} from "@vue/composition-api";
  import {useRouter} from "@/features/router";
  import {useAuthentication} from "@/features/auth";

  export default defineComponent({
    mixins: [ThumbnailBackground],
    props: {
      imageId: String,
      caption: String,
      selecting: Boolean,
      selected: Boolean,
    },
    setup(props, ctx) {
      const {router} = useRouter();
      const {loggedIn} = useAuthentication();

      function handleToggle() {
        if (props.selected) {
          ctx.emit("deselected", props.imageId);
        } else {
          ctx.emit("selected", props.imageId);
        }
      }

      function handleClick(e: MouseEvent) {
        if (props.selecting && e.ctrlKey) {
          // Ctrl+click during selection is a shortcut for toggling
          handleToggle();
        } else if (props.selecting && e.shiftKey) {
          // Shift+click is a shortcut for range selection
          ctx.emit("select-range", props.imageId);
        } else {
          // Otherwise just show the lightbox
          router.push({path: "photo/" + props.imageId});
        }
      }

      return {handleClick, handleToggle, loggedIn};
    },
  });
</script>
