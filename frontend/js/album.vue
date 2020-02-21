<template>
    <gallery v-bind:endpoint="'/albums/' + id + '/photos'"></gallery>
</template>

<script>
  import Gallery from './gallery'

  export default {
    components: {
      Gallery
    },
    props: ['id']
  }
</script>
