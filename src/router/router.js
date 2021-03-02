import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/home/Home'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/logout',
      name: 'logout',
      beforeEnter() {
        window.location.href = '/logout'
      }
    },
    {
      path: '/login',
      name: 'login',
      beforeEnter() {
        window.location.href = '/login'
      }
    },
    {
      path: '/user/show/:uid',
      name: 'user-show',
      beforeEnter(to) {
        let userv = process.env.VUE_APP_USER_DOMAIN
        window.location.href = `${userv}#/show/${to.params.uid}`
      }
    },
    {
      path: '/login',
      name: 'logins',
      beforeEnter(to) {
        let sngHome = process.env.VUE_APP_SNG_HOME
        window.location.href = `${sngHome}${to.params.type}/games/${to.params.status}`
      }
    },
  ]
})
