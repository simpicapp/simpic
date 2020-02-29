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

<script lang="ts">
  import "vue-awesome/icons/spinner";
  import Icon from "vue-awesome/components/Icon.vue";
  import {computed, defineComponent} from "@vue/composition-api";

  export default defineComponent({
    props: {
      name: String,
      working: Boolean,
    },
    components: {
      Icon,
    },
    setup(props) {
      const icon = computed(() => (props.working ? "spinner" : props.name));
      return {icon};
    },
  });
</script>
