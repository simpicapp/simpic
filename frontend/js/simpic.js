import Axios from 'axios'
import Vue from 'vue'
import VueRouter from 'vue-router'

import Album from './album'
import Albums from './albums'
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
            path: '/timeline/photo/:id',
            props: true
          }
        ],
        component: Timeline,
        path: '/timeline/'
      },
      {
        component: Albums,
        path: '/albums/'
      },
      {
        children: [
          {
            component: Lightbox,
            path: '/albums/:album/photo/:id',
            props: true
          }
        ],
        component: Album,
        path: '/albums/:id',
        props: true
      },
      {
        path: '/',
        redirect: '/timeline/'
      }
    ]
  }),
  data: {
    gitSHA: '',
    loggedIn: false,
    username: '',
    version: '1.0+dev'
  },
  el: '#main',
  components: { App },
  render: (h) => h(App),
  methods: {
    checkUser () {
      Axios.get('/users/me').then(({ data: { username } }) => {
        this.username = username
        this.loggedIn = true
      }).catch(() => {
        this.loggedIn = false
      })
    }
  },
  created () {
    this.checkUser()
  }
})
