<template>
    <div>
        <toolbar></toolbar>
        <timeline></timeline>
        <uploader></uploader>
        <login></login>
        <lightbox></lightbox>
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
    import lightbox from "./lightbox";
    import timeline from "./timeline";
    import uploader from "./uploader";
    import login from "./login";
    import toolbar from "./toolbar";

    export default {
        components: {
            timeline,
            uploader,
            lightbox,
            login,
            toolbar
        }
    }
</script>
