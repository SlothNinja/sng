<template>
  <v-app id='app'>
    <sn-toolbar v-model='nav'></sn-toolbar>
    <sn-nav-drawer v-model='nav' app></sn-nav-drawer>
    <v-main>
      <v-container fluid>
      <v-card>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols='12' md='6' >
                <div class='title font-weight-black text-center'>Welcome to SlothNinja Games</div>
                <div class='body-1 pb-2'>
                  SlothNinja Games is a play-by-web (PBW) site that permits registered members to play board
                  and card games with other members in a turn-based manner. The site provides a graphical
                  interface where players can see the current game state and execute their actions
                  via a web browser. The site keeps track of the game state, enforces legal moves, and optionally
                  sends email notifications to players when it is their turn to move.
                </div>
                <div class='body-1 pb-2'>
                  Registration is required in order to play. Registration is free, but requires a Google Account.
                </div>
                <div class='body-1 pb-2'>
                  Due to the PBW nature of the site, games typically take days, weeks, or even months to complete
                  depending on how frequently players take their turns. 
                </div>
                <div class='body-1'>
                  There is a general expectation among the community that players will make
                  at least one move each day. If you plan to consistently take longer than a day per move,
                  please create a new game and clearly indicate your desired playing speed in the title.
                </div>
              </v-col>

              <v-col cols='12' md='6'>
                <v-card 
                     min-width='320'
                     class='my-4'
                     v-for="(game, index) in games"
                     :key='index'
                     :to='game.homelink' >
                     <v-card-text>
                       <v-row>
                         <v-col align='left' cols='4'>
                           <v-img
                             :src="game.image"
                             height='100'
                             contain
                             >
                           </v-img>
                         </v-col>
                  <v-col cols='8'>
                    <v-card-title>{{game.name}}</v-card-title>
                    <v-card-actions v-if='cu'>
                      <v-row>
                        <v-col>
                          <v-btn small color='blue' dark :to='game.playlink'>Play</v-btn>
                        </v-col>
                        <v-col>
                          <v-btn small color='blue' dark :to='game.createlink'>Create</v-btn>
                        </v-col>
                      </v-row>
                    </v-card-actions>
                  </v-col>
                       </v-row>
                     </v-card-text>
                </v-card>
              </v-col>

            </v-row>
          </v-container>
        </v-card-text>
      </v-card>
      </v-container>
    </v-main>
    <sn-footer></sn-footer>
  </v-app>
</template>

<script>
  import Toolbar from '@/components/Toolbar'
  import NavDrawer from '@/components/NavDrawer'
  import Footer from '@/components/Footer'
  import CurrentUser from '@/components/mixins/CurrentUser'

  const _ = require('lodash')
  const axios = require('axios')

  export default {
    name: 'home',
    mixins: [ CurrentUser ],
    data: function () {
      return {
        nav: false
      }
    },
    components: {
      'sn-toolbar': Toolbar,
      'sn-nav-drawer': NavDrawer,
      'sn-footer': Footer
    },
    created () {
      let self = this
      self.fetchData()
    },
    methods: {
      myUpdate: function (data) {
        let self = this

        if (_.has(data, 'cu')) {
          self.cu = data.cu
          self.cuLoading = false
        }
      },
      fetchData: function () {
        let self = this
        self.loading = true
        axios.get("/home")
          .then(function (response) {
            var data = _.get(response, 'data', false)
            if (data) {
              self.myUpdate(data)
            }
            self.loading = false
          })
          .catch(function () {
            self.loading = false
            self.sbMessage = 'Server Error.  Try refreshing page.'
            self.sbOpen = true
        })
      },
    },
    computed: {
      games: function () {
        return [
          {
            name: 'After the Flood',
            image: 'images/atf/ATF-box.jpg',
            homelink: { name: 'sng-home' },
            playlink: { name: 'games', params: { type: 'atf', status: 'running' } },
            createlink: { name: 'sng-new-game', params: { type: 'atf' } },
          },
          {
            name: 'Confucius',
            image: 'images/confucius/confucius-box.jpg',
            homelink: { name: 'sng-home' },
            playlink: { name: 'games', params: { type: 'confucius', status: 'running' } },
            createlink: { name: 'sng-new-game', params: { type: 'confucius' } },
          },
          {
            name: 'Guild of Thieves',
            image: 'images/got/got-box.jpg',
            homelink: { name: 'got-home' },
            playlink: { name: 'got-games', params: { status: 'running' } },
            createlink: { name: 'got-new-game' },
          },
          {
            name: 'Indonesia',
            image: 'images/indonesia/indonesia-box.jpg',
            homelink: { name: 'sng-home' },
            playlink: { name: 'games', params: { type: 'indonesia', status: 'running' } },
            createlink: { name: 'sng-new-game', params: { type: 'indonesia' } },
          },
          {
            name: 'Tammany Hall',
            image: 'images/tammany/tammany-box.jpg',
            homelink: { name: 'sng-home' },
            playlink: { name: 'games', params: { type: 'tammany', status: 'running' } },
            createlink: { name: 'sng-new-game', params: { type: 'tammany' } },
          },
        ]
      }
    }
  }
</script>
