<template>
  <Icon :name="icon" :spin="working" @click="$emit('click')" scale="1.5"></Icon>
</template>

<style lang="scss">
  @use "src/assets/css/vars";

  .fa-icon {
    transition: color 100ms linear;

    &:hover {
      color: vars.$primary;
    }
  }
</style>

<script>
  import 'vue-awesome/icons/spinner'
  import Icon from 'vue-awesome/components/Icon'
  import Vue from 'vue'

  export default Vue.extend({
    props: ['name', 'working'],
    components: {
      Icon
    },
    computed: {
      icon () {
        return this.working ? 'spinner' : this.name
      }
    }
  })
</script>
