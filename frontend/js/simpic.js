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
        loggedIn: false,
        token: "",
        username: ""
    },
    el: '#main',
    template: `<App></App>`,
    components: {App},
});
