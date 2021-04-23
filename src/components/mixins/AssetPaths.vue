<script>

  const atfBox = require('@/assets/boxes/atf-box.jpg')
  const confuciusBox = require('@/assets/boxes/confucius-box.jpg')
  const gotBox = require('@/assets/boxes/got-box.jpg')
  const indonesiaBox = require('@/assets/boxes/indonesia-box.jpg')
  const tammanyBox = require('@/assets/boxes/tammany-box.jpg')

  export default {
    methods: {
      assetPath: function (asset) {
        switch (asset) {
          case 'atf-box':
            return atfBox
          case 'confucius-box':
            return confuciusBox
          case 'got-box':
            return gotBox
          case 'indonesia-box':
            return indonesiaBox
          case 'tammany-box':
            return tammanyBox
          default:
            console.log(`no asset for "${asset}"`)
            return ''
        }
      }
    }
  }
</script>
