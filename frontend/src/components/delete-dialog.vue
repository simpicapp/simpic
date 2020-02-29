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

<script>
  import ConfirmDialog from './confirm-dialog'

  import Vue from 'vue'

  export default Vue.extend({
    components: {
      ConfirmDialog
    },
    props: ['what']
  })
</script>
