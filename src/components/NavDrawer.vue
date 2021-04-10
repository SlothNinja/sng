<template>
  <v-navigation-drawer
    clipped
    v-model='nav'
    light
    app
  >
  <v-list-item v-if='cu'>
      <v-list-item-icon>
        <sn-user-btn size='x-small' :user='cu' ></sn-user-btn>
      </v-list-item-icon>
      <v-list-item-content>
        <v-list-item-title>
          {{cu.name}}
        </v-list-item-title>
      </v-list-item-content>
    </v-list-item>

    <v-list-item v-else :to="{ name: 'login'}" >
      <v-list-item-icon>
        <v-icon>mdi-login</v-icon>
      </v-list-item-icon>
      <v-list-item-content>
        <v-list-item-title>Login</v-list-item-title>
      </v-list-item-content>
    </v-list-item>

    <v-divider></v-divider>

    <v-list
      dense
      nav
    >
      <v-list-item :to="{ name: 'sng-home' }" exact>
        <v-list-item-icon>
          <v-icon>mdi-home</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>Home</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item :to="{ name: 'ugames', params: { status: 'running', type: 'all', uid: cuid} }" exact>
        <v-list-item-icon>
          <v-icon>mdi-account-details</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>Your Games</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <template v-if='cu'>
        <v-list-group
          no-action
          prepend-icon='mdi-pencil'
          >
          <template v-slot:activator>
            <v-list-item-title>Create</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.createlink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-plus'
          >
          <template v-slot:activator>
            <v-list-item-title>Join</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.joinlink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-play'
          >
          <template v-slot:activator>
            <v-list-item-title>Play</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.playlink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-check'
          >
          <template v-slot:activator>
            <v-list-item-title>Completed</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.completedlink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-star'
          >
          <template v-slot:activator>
            <v-list-item-title>Top Players</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.ratingslink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-item :to="{ name: 'logout'}" >
          <v-list-item-icon>
            <v-icon>mdi-logout</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Logout</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
  import UserButton from '@/components/user/Button'
  import CurrentUser from '@/components/mixins/CurrentUser'

  export default {
    mixins: [ CurrentUser ],
    name: 'nav-drawer',
    props: [ 'value' ],
    components: {
      'sn-user-btn': UserButton
    },
    computed: {
      items: function () {
        return [
          { 
            createlink: { name: 'sng-new-game', params: { type: 'atf' } },
            joinlink: { name: 'sng-games', params: { type: 'atf', status: 'recruiting' } },
            playlink: { name: 'games', params: { type: 'atf', status: 'running' } },
            completedlink: { name: 'games', params: { type: 'atf', status: 'completed' } },
            ratingslink: { name: 'sng-ratings', params: { type: 'atf' } },
            title: "After the Flood"
          },
          { 
            createlink: { name: 'sng-new-game', params: { type: 'confucius' } },
            joinlink: { name: 'sng-games', params: { type: 'confucius', status: 'recruiting' } },
            playlink: { name: 'games', params: { type: 'confucius', status: 'running' } },
            completedlink: { name: 'games', params: { type: 'confucius', status: 'completed' } },
            ratingslink: { name: 'sng-ratings', params: { type: 'confucius' } },
            title: "Confucius"
          },
          { 
            createlink: { name: 'got-new-game' },
            joinlink: { name: 'got-join-game' },
            playlink: { name: 'got-games', params: { status: 'running' } },
            completedlink: { name: 'got-games', params: { status: 'completedlink' } },
            ratingslink: { name: 'got-ratings' },
            title: "Guild of Thieves"
          },
          { 
            createlink: { name: 'sng-new-game', params: { type: 'indonesia' } },
            joinlink: { name: 'sng-games', params: { type: 'indonesia', status: 'recruiting' } },
            playlink: { name: 'games', params: { type: 'indonesia', status: 'running' } },
            completedlink: { name: 'games', params: { type: 'indonesia', status: 'completed' } },
            ratingslink: { name: 'sng-ratings', params: { type: 'indonesia' } },
            title: "Indonesia"
          },
          { 
            createlink: { name: 'sng-new-game', params: { type: 'tammany' } },
            joinlink: { name: 'sng-games', params: { type: 'tammany', status: 'recruiting' } },
            playlink: { name: 'games', params: { type: 'tammany', status: 'running' } },
            completedlink: { name: 'games', params: { type: 'tammany', status: 'completed' } },
            ratingslink: { name: 'sng-ratings', params: { type: 'tammany' } },
            title: "Tammany Hall"
          },
        ]
      },
      nav: {
        get: function () {
          return this.value
        },
        set: function (value) {
          this.$emit('input', value)
        }
      }
    }
  }
</script>
