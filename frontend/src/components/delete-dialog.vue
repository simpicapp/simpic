<template>
  <ConfirmDialog
    :body="`Are you sure you wish to delete ${what}?`"
    :dangerous="true"
    @close="$emit('close')"
    @no="$emit('no')"
    @yes="$emit('yes')"
    no-text="No, cancel"
    title="Confirm deletion"
    yes-text="Yes, delete"
  >
  </ConfirmDialog>
</template>

<script lang="ts">
  import ConfirmDialog from "./confirm-dialog.vue";
  import {defineComponent} from "@vue/composition-api";

  export default defineComponent({
    components: {
      ConfirmDialog,
    },
    props: {what: String},
  });
</script>
