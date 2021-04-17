import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/home/Home'
import Games from '@/components/game/Index'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/got-home',
      name: 'got-home',
      beforeEnter() {
        let gotHome = process.env.VUE_APP_GOT_HOME
        window.location.href = gotHome
      }
    },
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/games/:status/:type',
      name: 'games',
      component: Games
    },
    {
      path: '/ugames/:uid/:status/:type',
      name: 'ugames',
      component: Games
    },
    {
      path: '/game/:type/:id',
      name: 'game',
      beforeEnter(to) {
        if (to.params.type != 'got') {
          let sngHome = process.env.VUE_APP_SNG_HOME
          window.location.href = `${sngHome}${to.params.type}/game/show/${to.params.id}`
        } else {
          let gotHome = process.env.VUE_APP_GOT_HOME
          window.location.href = `${gotHome}#/game/${to.params.id}`
        }
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
  ]
})
