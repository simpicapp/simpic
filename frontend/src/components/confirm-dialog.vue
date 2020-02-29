<template>
    <modal :closeable="false" :should-close="shouldClose" @close="onClose">
        <popup :title="title" position="center" :closeable="false">
            <p>{{ body }}</p>
            <div class="buttons">
                <button @click="onYes" :class="{'danger-button': dangerous}">{{ yesText }}</button>
                <button @click="onNo">{{ noText }}</button>
            </div>
        </popup>
    </modal>
</template>

<style lang="scss" scoped>
    .buttons {
        display: grid;
        grid-template-columns: auto auto;
        grid-column-gap: 20px;
        justify-items: stretch;
        margin-top: 30px;
    }
</style>

<script>
  import Modal from './modal'
  import Popup from './popup'

  export default {
    props: {
      body: String,
      dangerous: {
        type: Boolean,
        default: false
      },
      noText: {
        type: String,
        default: 'No'
      },
      title: String,
      yesText: {
        type: String,
        default: 'Yes'
      }
    },
    data () {
      return {
        shouldClose: false
      }
    },
    components: {
      Modal,
      Popup
    },
    methods: {
      onClose () {
        this.$emit('close')
        this.shouldClose = false
      },
      onNo () {
        this.$emit('no')
        this.shouldClose = true
      },
      onYes () {
        this.$emit('yes')
        this.shouldClose = true
      }
    }
  }
</script>
