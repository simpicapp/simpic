<template>
  <div>
    <div class="upload-overlay" v-if="dragging && $root.loggedIn">
      <div id="drop-target">
        Drop files here to add to simpic
      </div>
    </div>
    <div class="upload-overlay" v-if="dragging && !$root.loggedIn">
      <div id="upload-login-prompt">
        You need to login before uploading files
      </div>
    </div>
    <popup @close="visible = false" id="uploader" title="Uploading..." v-if="visible">
      <table>
        <tbody>
        <tr :key="file.name" v-for="file in files">
          <td>{{ file.name }}</td>
          <td v-if="file.failed">Error</td>
          <td v-else-if="file.finished">Done</td>
          <td v-else-if="file.started">Uploading</td>
          <td v-else>Waiting</td>
        </tr>
        </tbody>
      </table>
    </popup>
  </div>
</template>

<style lang="scss" scoped>
  td:first-child {
    padding-right: 10px;
  }

  .upload-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #cccccccc;
  }

  .upload-overlay > div {
    width: 50%;
    height: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: xx-large;
  }

  #drop-target {
    background-color: lightsteelblue;
    border: 10px dashed midnightblue;
    border-radius: 10px;
  }

  #upload-login-prompt {
    background-color: #ffcccc;
    border: 1px solid red;
  }
</style>

<script>
  import { debounce } from 'lodash-es'
  import Axios from 'axios'
  import { EventBus } from './bus'
  import popup from './popup'
  import Vue from 'vue'

  export default Vue.extend({
    components: {
      popup
    },
    data () {
      return {
        dragging: false,
        files: [],
        nextUpload: 0,
        stopDragging: debounce(function () {
          this.dragging = false
        }, 100),
        visible: false
      }
    },
    methods: {
      acceptNewFile (file) {
        this.files.push({
          failed: false,
          file,
          finished: false,
          name: file.name,
          started: false
        })

        this.visible = true

        if (this.nextUpload === this.files.length - 1) {
          this.startUpload()
        }
      },
      dragEnterHandler (e) {
        if (e.dataTransfer.types.includes('Files')) {
          e.stopPropagation()
          e.preventDefault()
          this.dragging = true
          this.stopDragging.cancel()
        }
      },
      dragLeaveHandler (e) {
        e.stopPropagation()
        e.preventDefault()
        this.stopDragging()
      },
      dragOverHandler (e) {
        if (e.dataTransfer.types.includes('Files')) {
          e.stopPropagation()
          e.preventDefault()

          if (this.$root.loggedIn) {
            e.dataTransfer.dropEffect = 'copy'
          } else {
            e.dataTransfer.effectAllowed = 'none'
            e.dataTransfer.dropEffect = 'none'
          }

          this.dragging = true
          this.stopDragging.cancel()
        }
      },
      dropHandler (e) {
        e.stopPropagation()
        e.preventDefault()

        this.stopDragging.cancel()
        this.dragging = false

        if (this.$root.loggedIn) {
          [...e.dataTransfer.files].forEach(this.acceptNewFile)
        }
      },
      startUpload () {
        const file = this.files[this.nextUpload]
        const formData = new FormData()
        formData.append('file', file.file)
        file.started = true

        Axios.post('/photos', formData).then(() => {
          file.finished = true
          EventBus.$emit('upload-complete')
        }).catch((e) => {
          console.log('Failed to upload file', file, e)
          file.failed = true
        }).finally(() => {
          this.nextUpload++
          if (this.nextUpload <= this.files.length - 1) {
            this.startUpload()
          }
        })
      }
    },
    mounted () {
      document.addEventListener('drop', this.dropHandler)
      document.addEventListener('dragover', this.dragOverHandler)
      document.addEventListener('dragenter', this.dragEnterHandler)
      document.addEventListener('dragleave', this.dragLeaveHandler)
    },
    beforeDestroy () {
      document.removeEventListener('drop', this.dropHandler)
      document.removeEventListener('dragover', this.dragOverHandler)
      document.removeEventListener('dragenter', this.dragEnterHandler)
      document.removeEventListener('dragleave', this.dragLeaveHandler)
    }
  })
</script>
