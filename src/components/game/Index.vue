<template>
  <v-app id='app'>
    <sn-toolbar v-model='nav'></sn-toolbar>
    <sn-nav-drawer v-model='nav' app></sn-nav-drawer>
    <sn-snackbar v-model='snackbar.open'>
      <div class='text-center'>
        <span v-html='snackbar.message'></span>
      </div>
    </sn-snackbar>
    <v-main>
      <v-container fluid>
        <v-card>
          <v-row v-if="image == 'box'">
            <v-col cols='1' style='min-width:80px;max-width:80px'>
              <v-img
                class='ma-2'
                min-width='74'
                max-width='74'
                :src="assetPath(`${$route.params.type}-box`)"
                >
              </v-img>
            </v-col>
            <v-col>
              <v-card-title>
                {{ gameName }}
              </v-card-title>
              <v-card-subtitle>
                {{ status }} Games
              </v-card-subtitle>
            </v-col>
          </v-row>
          <v-row v-else>
            <v-col cols='1' v-if="image == 'button'" style='min-width:80px;max-width:80px'>
              <sn-user-btn v-if='cu' :user='cu'></sn-user-btn>
            </v-col>
            <v-col v-if='cu'>
              <v-card-title>
                {{ cu.name }}
              </v-card-title>
              <v-card-subtitle>
                {{ gameName }} {{ status }} Games
              </v-card-subtitle>
            </v-col>
            <v-col v-else>
              <v-card-title>
                {{ gameName }} {{ status }} Games
              </v-card-title>
            </v-col>
          </v-row>
          <v-data-table
            @click:row="showGame"
            :headers="headers"
            :items="items"
            :loading="loading"
            :options.sync="options"
            loading-text="Loading... Please wait"
            :server-items-length="totalItems"
            :items-per-page='10'
            :footer-props="{ itemsPerPageOptions: [ 10, 25, 50 ] }"
            >
            <template v-if="gameName == ''" v-slot:item.type="{ item }">
              {{typeName(item)}}
            </template>
            <template v-slot:item.creator="{ item }">
              <v-row no-gutters>
                <v-col>
                  <sn-user-btn :user="creator(item)" size="x-small">
                    {{creator(item).name}}
                  </sn-user-btn>
                </v-col>
              </v-row>
            </template>
            <template v-slot:item.players="{ item }">
              <v-row no-gutters>
                <v-col class='my-1' cols='12' md='6' v-for="user in users(item)" :key="`${item.id}-${user.id}`" >
                  <sn-user-btn :user="user" size="x-small">
                    <span :class='userClass(item, user)'>{{user.name}}</span>
                  </sn-user-btn>
                </v-col>
              </v-row>
            </template>
          </v-data-table>
        </v-card>
      </v-container>
    </v-main>
    <sn-footer app></sn-footer>
  </v-app>
</template>

<script>
import UserButton from '@/components/lib/user/Button'
import Toolbar from '@/components/lib/Toolbar'
import NavDrawer from '@/components/lib/NavDrawer'
import Snackbar from '@/components/lib/Snackbar'
import Footer from '@/components/lib/Footer'
import CurrentUser from '@/components/lib/mixins/CurrentUser'
import AssetPaths from '@/components/mixins/AssetPaths'

const _ = require('lodash')
const axios = require('axios')

export default {
  name: 'index',
  mixins: [ CurrentUser, AssetPaths ],
  components: {
    'sn-user-btn': UserButton,
    'sn-toolbar': Toolbar,
    'sn-nav-drawer': NavDrawer,
    'sn-snackbar': Snackbar,
    'sn-footer': Footer
  },
  data () {
    return {
      cursors: [ "" ],
      loading: 'false',
      totalItems: 0,
      options: {},
      items: []
    }
  },
  created () {
    this.fetchData()
  },
  watch: {
    '$route': 'fetchData',
    options: {
      handler (val, oldVal) {
        if (val.itemsPerPage != oldVal.itemsPerPage) {
          this.cursors = [ "" ]
        }
        this.fetchData()
      },
      deep: true,
    },
  },
  mounted () {
    this.fetchData()
  },
  methods: {
    showGame: function (item) {
      this.$router.push({name: 'game', params: { type: this.type(item), id: item.id } })
    },
    type: function (item) {
      switch(item.type) {
        case 1:
          return 'confucius'
        case 2:
          return 'tammany'
        case 3:
          return 'atf'
        case 4:
          return 'got'
        case 5:
          return 'indonesia'
        default:
          return ''
      }
    },
    typeName: function (item) {
      switch(item.type) {
        case 1:
          return 'Confucius'
        case 2:
          return 'Tammany Hall'
        case 3:
          return 'After the Flood'
        case 4:
          return 'Guild of Thieves'
        case 5:
          return 'Indonesia'
        default:
          return ''
      }
    },
    fetchData: _.debounce(function () {
      let self = this
      self.loading = true
      let data = {
        options: self.options,
        forward: self.forward,
        status: self.$route.params.status,
        type: self.$route.params.type,
      }
      if (self.$route.name == "ugames") {
        data.userId = _.toInteger(self.$route.params.uid)
      }
      console.log(`data: ${JSON.stringify(data)}`)
      axios.post('/games', data)
        .then(function (response) {

          if (_.has(response, 'data.message')) {
            self.snackbar.message = response.data.message
            self.snackbar.open = true
          }

          if (_.has(response, 'data.totalItems')) {
            self.totalItems = response.data.totalItems
          }

          if (_.has(response, 'data.forward')) {
            self.forward = response.data.forward
          }

          if (_.has(response, 'data.gheaders')) {
            self.items = response.data.gheaders
          }

          if (_.has(response, 'data.cu')) {
            self.cu = response.data.cu
            self.cuLoading = false
          }

          self.loading = false
        })
        .catch(function () {
          self.loading = false
          self.snackbar.message = 'Server Error.  Try refreshing page.'
          self.snackbar.open = true
        })
    }, 1000),
    action: function (action, id) {
      let self = this
      axios.put(`/game/${action}/${id}`)
        .then(function (response) {

          if (_.has(response, 'data.message')) {
            self.snackbar.message = response.data.message
            self.snackbar.open = true
          }

          if (_.has(response, 'data.header')) {
            let index = _.findIndex(self.items, [ 'id', id ])
            if (index >= 0) {
              if (response.data.header.status === 1) { // recruiting is a status of 1
                self.items.splice(index, 1, response.data.header)
              } else {
                self.items.splice(index, 1)
              }
            }
          }
          self.loading = false
        })
        .catch(function () {
          self.loading = false
          self.snackbar.message = 'Server Error.  Try refreshing page.'
          self.snackbar.open = true
        })
    },
    canAccept: function (id) {
      let self = this
      let item = self.getItem(id)
      return !self.joined(item) && item.status === 1 // recruiting is a status 1
    },
    canDrop: function (id) {
      let self = this
      let item = self.getItem(id)
      return self.joined(item) && item.status === 1 // recruiting is a status 1
    },
    joined: function (item) {
      let self = this
      return _.find(item.users, [ 'id', self.cuid ])
    },
    getItem: function (id) {
      let self = this
      return _.find(self.items, [ 'id', id ])
    },
    publicPrivate: function (item) {
      return item.public ? 'Public' : 'Private'
    },
    creator: function (item) {
      return {
        id: item.creatorId,
        name: item.creatorName,
        emailHash: item.creatorEmailHash,
        gravType: item.creatorGravType
      }
    },
    users: function (item) {
      return _.map(item.userIds, function (id, i) {
        return {
          id: id,
          name: _.nth(item.userNames, i),
          emailHash: _.nth(item.userEmailHashes, i),
          gravType: _.nth(item.userGravTypes, i),
        }
      })
    },
    userClass: function (item, user) {
      const completed = 2
      if (item.status == completed) {
        return this.winnerClass(item, user)
      } 
      return this.cpClass(item, user)
    },
    cpClass: function (item, user) {
      switch (this.type(item)) {
        case 'got':
          var pid = _.indexOf(item.userIds, user.id) + 1
          if (_.includes(item.cpids, pid)) {
            if (this.cuid == user.id) {
              return 'font-weight-black red--text text--darken-4'
            }
            return 'font-weight-black'
          }
          return ''
        default:
          var uIndex = _.indexOf(item.userIds, user.id)
          if (_.includes(item.cpUserIndices, uIndex)) {
            if (this.cuid == user.id) {
              return 'font-weight-black red--text text--darken-4'
            }
            return 'font-weight-black'
          }
          return ''
      }
    },
    winnerClass: function (item, user) {
      let uIndex = _.indexOf(item.userIds, user.id)
      if (_.includes(item.winnerIndices, uIndex)) {
        if (this.cuid == user.id) {
          return 'font-weight-black red--text text--darken-4'
        }
        return 'font-weight-black'
      }
      return ''
    },
  },
  computed: {
    image: function () {
      if (this.$route.name == 'ugames') {
        return 'button'
      }
      if (this.gameName) {
        return 'box'
      }
      return ''
    },
    forward: {
      get: function () {
        return this.cursors[this.options.page-1]
      },
      set: function (value) {
        this.cursors.splice(this.options.page, 1, value)
      }
    },
    headers: function () {
      if (this.gameName == '') {
        return [
          { text: 'ID', align: 'left', sortable: false, value: 'id' },
          { text: 'Type', align: 'left', sortable: false, value: 'type' },
          { text: 'Title', sortable: false, value: 'title' },
          { text: 'Creator', sortable: false, width: '180px', value: 'creator' },
          { text: 'Players', sortable: false, width: this.width, value: 'players' },
          { text: 'Last Updated', sortable: false, value: 'lastUpdated' },
        ]
      }
      return [
        { text: 'ID', align: 'left', sortable: false, value: 'id' },
        { text: 'Title', sortable: false, value: 'title' },
        { text: 'Creator', sortable: false, width: '180px', value: 'creator' },
        { text: 'Players', sortable: false, width: this.width, value: 'players' },
        { text: 'Last Updated', sortable: false, value: 'lastUpdated' },
      ]
    },
    width: function () {
      console.log(`mdAndUp: ${this.$vuetify.breakpoint.mdAndUp}`)
      if (this.$vuetify.breakpoint.mdAndUp) {
        return '30%'
      } else {
        return '180px'
      }
    },
    gameName: function () {
      switch(this.$route.params.type) {
        case 'confucius':
          return 'Confucius'
        case 'tammany':
          return 'Tammany Hall'
        case 'atf':
          return 'After The Flood'
        case 'got':
          return 'Guild of Thieves'
        case 'indonesia':
          return 'Indonesia'
        default:
          return ''
      }
    },
    status: function () {
      return _.capitalize(this.$route.params.status)
    },
    snackbar: {
      get: function () {
        return this.$root.snackbar
      },
      set: function (value) {
        this.$root.snackbar = value
      }
    },
    nav: {
      get: function () {
        return this.$root.nav
      },
      set: function (value) {
        this.$root.nav = value
      }
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
::v-deep tbody tr:nth-of-type(odd) {
  background-color: rgba(0, 0, 0, .04);
}
</style>
