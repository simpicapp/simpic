<template>
    <div id="uploader" v-if="visible">
        <h2>Uploading...</h2>
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
    </div>
</template>

<style>
    #uploader {
        position: absolute;
        bottom: 20px;
        right: 20px;
        background: white;
        border: 2px solid black;
        border-radius: 2px;
        box-shadow: cornflowerblue 5px 5px;
    }

    #uploader h2 {
        text-align: center;
        border-bottom: 1px solid black;
        font-size: medium;
        margin: 0;
    }
</style>

<script>
    import {EventBus} from './bus';

    export default {
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
