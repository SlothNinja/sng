<template>
  <v-navigation-drawer
    clipped
    v-model='drawer'
    light
    app
  >
    <v-list-item>
      <v-list-item-content>
        <v-list-item-title class='title font-weight-black text-center'>
          SlothNinja Games
        </v-list-item-title>
      </v-list-item-content>
    </v-list-item>

    <v-divider></v-divider>

    <v-list
      dense
      nav
    >
      <v-list-item :to="{ name: 'home' }" exact>
        <v-list-item-icon>
          <v-icon>mdi-home</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>Home</v-list-item-title>
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
          <v-list-item href='/atf/game/new'>
            <v-list-item-title>After the Flood</v-list-item-title>
          </v-list-item>
          <v-list-item href='/confucius/game/new'>
            <v-list-item-title>Confucius</v-list-item-title>
          </v-list-item>
          <v-list-item href='https://got.slothninja.com/#/new'>
            <v-list-item-title>Guild of Thieves</v-list-item-title>
          </v-list-item>
          <v-list-item href='/indonesia/game/new'>
            <v-list-item-title>Indonesia</v-list-item-title>
          </v-list-item>
          <v-list-item href='/tammany/game/new'>
            <v-list-item-title>Tammany</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-plus'
          >
          <template v-slot:activator>
            <v-list-item-title>Join</v-list-item-title>
          </template>
          <v-list-item href='/atf/games/recruiting'>
            <v-list-item-title>After the Flood</v-list-item-title>
          </v-list-item>
          <v-list-item href='/confucius/games/recruiting'>
            <v-list-item-title>Confucius</v-list-item-title>
          </v-list-item>
          <v-list-item href='https://got.slothninja.com/#/invitations'>
            <v-list-item-title>Guild of Thieves</v-list-item-title>
          </v-list-item>
          <v-list-item href='/indonesia/games/recruiting'>
            <v-list-item-title>Indonesia</v-list-item-title>
          </v-list-item>
          <v-list-item href='/tammany/games/recruiting'>
            <v-list-item-title>Tammany</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-play'
          >
          <template v-slot:activator>
            <v-list-item-title>Play</v-list-item-title>
          </template>
          <v-list-item href='/atf/games/running'>
            <v-list-item-title>After the Flood</v-list-item-title>
          </v-list-item>
          <v-list-item href='/confucius/games/running'>
            <v-list-item-title>Confucius</v-list-item-title>
          </v-list-item>
          <v-list-item href='https://got.slothninja.com/#/games/running'>
            <v-list-item-title>Guild of Thieves</v-list-item-title>
          </v-list-item>
          <v-list-item href='/indonesia/games/running'>
            <v-list-item-title>Indonesia</v-list-item-title>
          </v-list-item>
          <v-list-item href='/tammany/games/running'>
            <v-list-item-title>Tammany</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-check'
          >
          <template v-slot:activator>
            <v-list-item-title>Completed</v-list-item-title>
          </template>
          <v-list-item href='/atf/games/completed'>
            <v-list-item-title>After the Flood</v-list-item-title>
          </v-list-item>
          <v-list-item href='/confucius/games/completed'>
            <v-list-item-title>Confucius</v-list-item-title>
          </v-list-item>
          <v-list-item href='https://got.slothninja.com/#/invitations'>
            <v-list-item-title>Guild of Thieves</v-list-item-title>
          </v-list-item>
          <v-list-item href='/indonesia/games/completed'>
            <v-list-item-title>Indonesia</v-list-item-title>
          </v-list-item>
          <v-list-item href='/tammany/games/completed'>
            <v-list-item-title>Tammany</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-star'
          >
          <template v-slot:activator>
            <v-list-item-title>Top Players</v-list-item-title>
          </template>
          <v-list-item href='/ratings/show/atf'>
            <v-list-item-title>After the Flood</v-list-item-title>
          </v-list-item>
          <v-list-item href='/ratings/show/confucius'>
            <v-list-item-title>Confucius</v-list-item-title>
          </v-list-item>
          <v-list-item href='https://got.slothninja.com/#/ratings'>
            <v-list-item-title>Guild of Thieves</v-list-item-title>
          </v-list-item>
          <v-list-item href='/ratings/show/indonesia'>
            <v-list-item-title>Indonesia</v-list-item-title>
          </v-list-item>
          <v-list-item href='/ratings/show/tammany'>
            <v-list-item-title>Tammany</v-list-item-title>
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
  import CurrentUser from '@/components/mixins/CurrentUser'

  export default {
    mixins: [ CurrentUser ],
    name: 'nav-drawer',
    props: [ 'value' ],
    methods: {
      logout: function () {
        var self = this
        self.$router.push({ name: 'logout'})
      }
    },
    computed: {
      drawer: {
        get: function () {
          var self = this
          return self.value
        },
        set: function (value) {
          var self = this
          self.$emit('input', value)
        }
      }
    }
  }
</script>
