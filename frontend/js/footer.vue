<template>
    <footer>
        <a href="https://simpic.app/">Simpic - simple self-hosted picture manager</a> &middot;
        <span class="version">pre-mvp development build</span>
    </footer>
</template>

<style lang="scss" scoped>
    footer {
        border-top: 1px solid #ccc;
        background-color: #eee;
        margin-top: 10px;
        height: 70px;
        display: flex;
        align-items: center;
        column-gap: 10px;
        padding-left: 20px;
    }

    a {
        color: black;
        transition: color 300ms linear;
        &:hover {
            color: blue;
        }
    }

    .version {
        color: #333;
    }
</style>
