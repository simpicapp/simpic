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
    token: localStorage.getItem('simpicAuthToken')
  },
  el: '#main',
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
  render: (h) => h(App),
  watch: {
    token (newValue) {
      localStorage.setItem('simpicAuthToken', newValue)
    }
  }
})
