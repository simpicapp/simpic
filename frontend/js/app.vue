<template>
    <div>
        <toolbar></toolbar>
        <router-view></router-view>
        <uploader></uploader>
        <login></login>
        <bottombar></bottombar>
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
    import bottombar from "./footer";

    export default {
        components: {
            uploader,
            login,
            toolbar,
            bottombar
        }
    }
</script>
