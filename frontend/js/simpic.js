import Vue from './vendor/vue'
import VueRouter from './vendor/vue-router'

import App from './app'
import Lightbox from './lightbox'
import Timeline from './timeline'

Vue.use(VueRouter);

new Vue({
    router: new VueRouter({
        routes: [
            {
                path: '/',
                component: Timeline,
                children: [
                    {
                        path: 'photo/:id',
                        name: 'lightbox',
                        component: Lightbox,
                        props: true
                    }
                ]
            }
        ]
    }),
    data: {
        token: localStorage.getItem("simpicAuthToken"),
    },
    el: '#main',
    template: `<App></App>`,
    components: {App},
    methods: {
        authHeaders() {
            if (this.loggedIn) {
                return {'Authorization': 'Bearer ' + this.token}
            } else {
                return {}
            }
        }
    },
    computed: {
        loggedIn() {
            return !!this.token;
        },
        username() {
            if (this.token) {
                return JSON.parse(atob(this.token.split(".")[1]))['unm']
            } else {
                return null
            }
        }
    },
    watch: {
        token(newValue) {
            localStorage.setItem('simpicAuthToken', newValue)
        }
    }
});
