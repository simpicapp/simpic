import Vue from './vendor/vue'
import VueRouter from './vendor/vue-router'

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
            path: 'photo/:id',
            props: true
          }
        ],
        component: Timeline,
        path: '/'
      }
    ]
  }),
  data: {
    token: localStorage.getItem('simpicAuthToken')
  },
  el: '#main',
  template: '<App></App>',
  components: { App },
  methods: {
    authHeaders () {
      if (this.loggedIn) {
        return { Authorization: 'Bearer ' + this.token }
      } else {
        return {}
      }
    }
  },
  computed: {
    loggedIn () {
      return !!this.token
    },
    username () {
      if (this.token) {
        return JSON.parse(atob(this.token.split('.')[1])).sub
      } else {
        return null
      }
    }
  },
  watch: {
    token (newValue) {
      localStorage.setItem('simpicAuthToken', newValue)
    }
  }
})
