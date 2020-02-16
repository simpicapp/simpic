<template>
    <popup title="Login" position="center" v-if="visible">
        <form>
            <label for="username">Username</label>
            <input type="text" id="username">
            <label for="password">Password</label>
            <input type="password" id="password">
            <input type="submit" value="Login">
        </form>
    </popup>
</template>

<style scoped>
    form {
        display: grid;
        grid-template-columns: auto auto;
        grid-gap: 20px 20px;
        padding: 20px;
        align-items: center;
    }

    input[type=submit] {
        padding: 5px 0;
        grid-column: span 2;
    }
</style>

<script>
    import popup from "./popup";

    export default {
        components: {
            popup
        },
        data() {
            return {
                visible: false
            }
        }
    }
</script>
