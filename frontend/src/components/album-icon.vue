<template>
  <a :href="'#/albums/' + id + '/'" :style="styles" @click.prevent="handleClick">
    <div class="caption"><span>{{ caption }}</span> <span>{{ photos }}</span></div>
  </a>
</template>

<style lang="scss" scoped>
  @use 'src/assets/css/vars';

  $border-radius: 10px;

  a {
    position: relative;
    margin: 10px 15px;
    min-width: 300px;
    max-width: 300px;
    height: 200px;
    display: block;
    text-decoration: none;
    border-radius: $border-radius;
    background-repeat: no-repeat;
    background-position: 50%;
    background-size: cover;
    box-shadow: 10px -10px 0px #99999966;

    &:hover {
      background-size: 110%;
      box-shadow: 10px -10px 0px #66666666;

      .caption {
        font-weight: bold;
      }
    }
  }

  .caption {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 10px;
    margin: 0;
    color: white;
    font-size: large;
    background-color: rgba(vars.$primary, 0.85);
    width: 100%;
    border-bottom-left-radius: $border-radius;
    border-bottom-right-radius: $border-radius;
    display: flex;
    justify-content: space-between;
  }
</style>

<script lang="ts">
  import ThumbnailBackground from './thumbnail-background.vue'
  import Vue from 'vue'

  export default Vue.extend({
    mixins: [ThumbnailBackground],
    props: ['id', 'caption', 'photos', 'imageId'],
    methods: {
      handleClick() {
        this.$router.push({path: './' + this.id + '/'})
      }
    }
  })
</script>
