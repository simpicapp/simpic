<template>
    <header>
        <h1>Simpic</h1>
        <button v-on:click="loginClick">Login</button>
    </header>
</template>

<style>
    header {
        display: flex;
        margin: 0;
        padding: 10px;
        border-bottom: 1px solid #ccc;
        justify-content: space-between;
    }

    header h1 {
        margin: 0;
        text-transform:  lowercase;
        flex-grow: 1;
    }
</style>

<script>
    import {EventBus} from "./bus";

    export default {
        methods: {
            loginClick() {
                EventBus.$emit('login')
            }
        }
    }
</script>
