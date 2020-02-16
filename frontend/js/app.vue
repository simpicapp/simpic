<template>
    <div>
        <toolbar></toolbar>
        <router-view></router-view>
        <uploader></uploader>
        <login></login>
    </div>
</template>

<style>
    body {
        font-family: sans-serif;
        margin: 0;
        padding: 0;
    }
</style>

<script>
    import uploader from "./uploader";
    import login from "./login";
    import toolbar from "./toolbar";

    export default {
        components: {
            uploader,
            login,
            toolbar
        }
    }
</script>
