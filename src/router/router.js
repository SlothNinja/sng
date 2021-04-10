import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/home/Home'
import Games from '@/components/game/Index'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/sng-home',
      name: 'sng-home',
      beforeEnter() {
        let sngHome = process.env.VUE_APP_SNG_HOME
        window.location.href = sngHome
      }
    },
    {
      path: '/sng-ratings/:type',
      name: 'sng-ratings',
      beforeEnter(to) {
        let sngHome = process.env.VUE_APP_SNG_HOME
        window.location.href = `${sngHome}ratings/show/${to.params.type}`
      }
    },
    {
      path: '/sng-games/:type/:status',
      name: 'sng-games',
      beforeEnter(to) {
        let sngHome = process.env.VUE_APP_SNG_HOME
        window.location.href = `${sngHome}${to.params.type}/games/${to.params.status}`
      }
    },
    {
      path: '/sng-new-game/:type',
      name: 'sng-new-game',
      beforeEnter(to) {
        let sngHome = process.env.VUE_APP_SNG_HOME
        window.location.href = `${sngHome}${to.params.type}/game/new`
      }
    },
    {
      path: '/got-home',
      name: 'got-home',
      beforeEnter() {
        let gotHome = process.env.VUE_APP_GOT_HOME
        window.location.href = gotHome
      }
    },
    {
      path: '/got-ratings',
      name: 'got-ratings',
      beforeEnter() {
        let gotHome = process.env.VUE_APP_GOT_HOME
        window.location.href = `${gotHome}#/ratings`
      }
    },
    {
      path: '/got-new-game',
      name: 'got-new-game',
      beforeEnter() {
        let gotHome = process.env.VUE_APP_GOT_HOME
        window.location.href = `${gotHome}#/invitation/new`
      }
    },
    {
      path: '/got-join-game',
      name: 'got-join-game',
      beforeEnter() {
        let gotHome = process.env.VUE_APP_GOT_HOME
        window.location.href = `${gotHome}#/invitations`
      }
    },
    {
      path: '/got-games/:status',
      name: 'got-games',
      beforeEnter(to) {
        let sngHome = process.env.VUE_APP_GOT_HOME
        window.location.href = `${sngHome}#/games/${to.params.status}`
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
      path: '/image/atf/ATF-box.jpg',
      name: 'atfbox'
    }
  ]
})
