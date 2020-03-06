<template>
  <div>
    <div class="upload-overlay" v-if="dragging && loggedIn">
      <div id="drop-target">
        Drop files here to add to simpic
      </div>
    </div>
    <div class="upload-overlay" v-if="dragging && !loggedIn">
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

<script lang="ts">
  import {debounce} from "lodash-es";
  import Axios from "axios";
  import Popup from "./popup.vue";
  import {defineComponent, reactive, toRefs} from "@vue/composition-api";
  import {useAuthentication} from "@/features/auth";
  import {useDocumentListener} from "@/features/listeners";
  import {EventBus} from "@/features/eventbus";

  export default defineComponent({
    components: {
      Popup,
    },
    setup() {
      const {loggedIn} = useAuthentication();

      interface Upload {
        started: boolean;
        finished: boolean;
        failed: boolean;
        name: string;
        file: File;
      }

      let nextUpload = 0;

      const state = reactive({
        dragging: false,
        files: new Array<Upload>(),
        visible: false,
      });

      const stopDragging = debounce(() => {
        state.dragging = false;
      }, 100);

      function startUpload() {
        const file = state.files[nextUpload];
        const formData = new FormData();
        formData.append("file", file.file);
        file.started = true;

        Axios.post("/api/photos", formData)
          .then(() => {
            file.finished = true;
            EventBus.$emit("upload-complete");
          })
          .catch(e => {
            console.log("Failed to upload file", file, e);
            file.failed = true;
          })
          .finally(() => {
            nextUpload++;
            if (nextUpload <= state.files.length - 1) {
              startUpload();
            }
          });
      }

      function acceptNewFile(file: File) {
        state.files.push({
          failed: false,
          file,
          finished: false,
          name: file.name,
          started: false,
        });

        state.visible = true;

        if (nextUpload === state.files.length - 1) {
          startUpload();
        }
      }

      useDocumentListener("drop", e => {
        e.stopPropagation();
        e.preventDefault();

        stopDragging.cancel();
        state.dragging = false;

        if (e.dataTransfer && loggedIn.value) {
          Array.from(e.dataTransfer.files).forEach(acceptNewFile);
        }
      });

      useDocumentListener("dragover", e => {
        if (e.dataTransfer && e.dataTransfer.types.includes("Files")) {
          e.stopPropagation();
          e.preventDefault();

          if (loggedIn.value) {
            e.dataTransfer.dropEffect = "copy";
          } else {
            e.dataTransfer.effectAllowed = "none";
            e.dataTransfer.dropEffect = "none";
          }

          state.dragging = true;
          stopDragging.cancel();
        }
      });

      useDocumentListener("dragenter", e => {
        if (e.dataTransfer && e.dataTransfer.types.includes("Files")) {
          e.stopPropagation();
          e.preventDefault();
          state.dragging = true;
          stopDragging.cancel();
        }
      });

      useDocumentListener("dragleave", e => {
        e.stopPropagation();
        e.preventDefault();
        stopDragging();
      });

      return {loggedIn, ...toRefs(state)};
    },
  });
</script>
