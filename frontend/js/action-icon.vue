<template>
    <Icon :spin="working" :name="icon" scale="1.5" @click="$emit('click')"></Icon>
</template>

<style lang="scss">
    @use "../css/vars";

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

  export default {
    props: ['name', 'working'],
    components: {
      Icon
    },
    computed: {
      icon () {
        return this.working ? 'spinner' : this.name
      }
    }
  }
</script>
