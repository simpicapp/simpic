<template>
    <ConfirmDialog :body="`Are you sure you wish to delete ${what}?`"
                   title="Confirm deletion"
                   yes-text="Yes, delete"
                   no-text="No, cancel"
                   :dangerous="true"
                   @yes="$emit('yes')"
                   @no="$emit('no')"
                   @close="$emit('close')">
    </ConfirmDialog>
</template>

<script>
  import ConfirmDialog from './confirm-dialog'

  export default {
    components: {
      ConfirmDialog
    },
    props: ['what']
  }
</script>
