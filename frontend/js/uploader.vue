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
        <popup title="Uploading..." id="uploader" v-if="visible" v-on:close="visible = false">
            <table>
                <tbody>
                <tr v-for="file in files" v-bind:key="file.name">
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

<style scoped>
    td {
        padding: 10px;
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
  import { EventBus } from './bus'
  import popup from './popup'

  export default {
    components: {
      popup
    },
    data () {
      return {
        dragging: false,
        files: [],
        nextUpload: 0,
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

        if (this.nextUpload === this.files.length - 1) {
          this.startUpload()
        }
      },
      dragEndHandler (e) {
        e.stopPropagation()
        e.preventDefault()
        this.dragging = false
      },
      dragOverHandler (e) {
        e.stopPropagation()
        e.preventDefault()

        if (this.$root.loggedIn) {
          e.dataTransfer.dropEffect = 'copy'
        } else {
          e.dataTransfer.effectAllowed = 'none'
          e.dataTransfer.dropEffect = 'none'
        }

        this.dragging = true
      },
      dragStartHandler (e) {
        e.stopPropagation()
        e.preventDefault()
        this.dragging = true
      },
      dropHandler (e) {
        e.stopPropagation()
        e.preventDefault()

        this.dragging = false

        if (this.$root.loggedIn) {
          this.visible = true;
          [...e.dataTransfer.files].forEach(this.acceptNewFile)
        }
      },
      startUpload () {
        const file = this.files[this.nextUpload]
        const formData = new FormData()
        formData.append('file', file.file)
        file.started = true

        fetch('/photo', {
          body: formData,
          method: 'POST'
        }).then(() => {
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
      document.addEventListener('dragenter', this.dragStartHandler)
      document.addEventListener('dragleave', this.dragEndHandler)
    },
    beforeDestroy () {
      document.removeEventListener('drop', this.dropHandler)
      document.removeEventListener('dragover', this.dragOverHandler)
      document.removeEventListener('dragenter', this.dragStartHandler)
      document.removeEventListener('dragleave', this.dragEndHandler)
    }
  }
</script>
