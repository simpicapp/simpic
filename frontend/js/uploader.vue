<template>
    <popup title="Uploading..." id="uploader" v-if="visible" v-on:close="visible = false">
        <table>
            <tbody>
            <tr v-for="file in files">
                <td>{{ file.name }}</td>
                <td v-if="file.failed">Error</td>
                <td v-else-if="file.finished">Done</td>
                <td v-else-if="file.started">Uploading</td>
                <td v-else>Waiting</td>
            </tr>
            </tbody>
        </table>
    </popup>
</template>

<style scoped>
    td {
        padding: 10px;
    }
</style>

<script>
    import {EventBus} from './bus';
    import popup from "./popup";

    export default {
        components: {
            popup
        },
        data() {
            return {
                visible: false,
                nextUpload: 0,
                files: []
            }
        },
        methods: {
            acceptNewFiles(newFiles) {
                this.visible = true;
                [...newFiles].forEach(this.acceptNewFile);
            },
            acceptNewFile(file) {
                let obj = {
                    file,
                    name: file.name,
                    failed: false,
                    started: false,
                    finished: false
                };
                this.files.push(obj);

                if (this.nextUpload === this.files.length - 1) {
                    this.startUpload();
                }
            },
            startUpload() {
                let file = this.files[this.nextUpload];
                let formData = new FormData();
                formData.append('file', file.file);
                file.started = true;

                fetch('/photo', {
                    method: 'POST',
                    body: formData
                }).then(() => {
                    file.finished = true;
                    EventBus.$emit('upload-complete');
                }).catch(() => {
                    file.failed = true;
                }).finally(() => {
                    this.nextUpload++;
                    if (this.nextUpload <= this.files.length - 1) {
                        this.startUpload();
                    }
                });
            }
        },
        created() {
            EventBus.$on('files-dropped', this.acceptNewFiles);
        }
    }
</script>
