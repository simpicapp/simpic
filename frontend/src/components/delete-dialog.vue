<template>
  <ConfirmDialog :body="`Are you sure you wish to delete ${what}?`"
                 :dangerous="true"
                 @close="$emit('close')"
                 @no="$emit('no')"
                 @yes="$emit('yes')"
                 no-text="No, cancel"
                 title="Confirm deletion"
                 yes-text="Yes, delete">
  </ConfirmDialog>
</template>

<script lang="ts">
  import ConfirmDialog from './confirm-dialog.vue'
  import Vue from 'vue'

  export default Vue.extend({
    components: {
      ConfirmDialog
    },
    props: ['what']
  })
</script>
