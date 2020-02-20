import Vue from 'vue'
import VueRouter from 'vue-router'

import App from './app'
import Lightbox from './lightbox'
import Timeline from './timeline'

Vue.use(VueRouter)

// eslint-disable-next-line no-new
new Vue({
  router: new VueRouter({
    routes: [
      {
        children: [
          {
            component: Lightbox,
            name: 'lightbox',
            path: '/timeline/photo/:id',
            props: true
          }
        ],
        component: Timeline,
        path: '/timeline/'
      },
      {
        path: '/',
        redirect: '/timeline/'
      }
    ]
  }),
  data: {
    loggedIn: false,
    username: ''
  },
  el: '#main',
  components: { App },
  render: (h) => h(App),
  methods: {
    checkUser () {
      fetch('/users/me')
        .then(res => {
          if (res.ok) {
            res.json().then(json => {
              this.username = json.username
              this.loggedIn = true
            })
          } else {
            this.loggedIn = false
          }
        })
    }
  },
  created () {
    this.checkUser()
  }
})
