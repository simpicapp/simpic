<template>
    <div>
        <timeline></timeline>
        <uploader></uploader>
        <login></login>
        <lightbox></lightbox>
    </div>
</template>

<style>
    body {
        font-family: sans-serif;
    }
</style>

<script>
    import lightbox from "./lightbox";
    import timeline from "./timeline";
    import uploader from "./uploader";
    import login from "./login";

    export default {
        components: {
            timeline,
            uploader,
            lightbox,
            login
        }
    }
</script>
