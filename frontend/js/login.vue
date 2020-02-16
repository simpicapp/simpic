<template>
    <popup title="Login" position="center" v-if="visible" v-on:close="visible = false">
        <form v-on:submit="doLogin">
            <p class="alert" v-if="alert.length > 0">{{ alert }}</p>
            <label for="username">Username</label>
            <input type="text" id="username" v-model="username" v-bind:disabled="loggingIn">
            <label for="password">Password</label>
            <input type="password" id="password" v-model="password" v-bind:disabled="loggingIn">
            <input type="submit" value="Login" v-bind:disabled="loggingIn">
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

    .alert {
        margin: 0;
        padding: 5px 10px;
        grid-column: span 2;
        background-color: darkred;
        color: white;
        font-weight: bold;
        text-align: center;
        border-radius: 15px;
        white-space: pre-line;
    }
</style>

<script>
    import {EventBus} from "./bus";
    import popup from "./popup";

    export default {
        components: {
            popup
        },
        data() {
            return {
                visible: false,
                loggingIn: false,
                username: "",
                password: "",
                alert: ""
            }
        },
        methods: {
            show() {
                this.visible = true;
            },
            doLogin() {
                this.alert = "";
                this.loggingIn = true;
                fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        username: this.username,
                        password: this.password
                    })
                }).then((response) => {
                    if (response.ok) {
                        return response.json();
                    } else {
                        return response.json().then((json) => {
                            throw new Error(json.error)
                        });
                    }
                }).then((json) => {
                    this.$root.loggedIn = true;
                    this.$root.token = json.token;
                    this.$root.username = this.username;
                    this.visible = false;
                }).catch((error) => {
                    this.alert = error.message;
                }).finally(() => {
                    this.loggingIn = false;
                })
            }
        },
        created() {
            EventBus.$on('login', this.show);
        }
        ,
        beforeDestroy() {
            EventBus.$off('login', this.show);
        }
    }
</script>
