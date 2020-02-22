<template>
    <div class="album">
        <h2>{{ name }}</h2>
        <gallery v-bind:endpoint="'/albums/' + id + '/photos'"></gallery>
    </div>
</template>

<style scoped>
    .album h2 {
        margin-left: 20px;
    }
</style>

<script>
  import Gallery from './gallery'

  export default {
    components: {
      Gallery
    },
    data () {
      return {
        name: ''
      }
    },
    props: ['id'],
    mounted () {
      fetch('/albums/' + this.id)
        .then(res => res.json())
        .then(res => (this.name = res.name))
    }
  }
</script>
