import Vue from 'vue'
import VueRouter from 'vue-router'
import './composition-plugin'

import App from './components/app.vue'
import {useRouter} from "@/features/router";

const {router} = useRouter();
Vue.use(VueRouter);

new Vue({
  router,
  data: {
    gitSHA: '',
    version: '1.0+dev'
  },
  render: (h) => h(App)
}).$mount('#main');
