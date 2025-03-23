import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// import Vant from 'vant'
// import 'vant/lib/index.css'


// // Vuetify
// import 'vuetify/styles'
// import { createVuetify } from 'vuetify'
// import * as components from 'vuetify/components'
// import * as directives from 'vuetify/directives'
// const vuetify = createVuetify({
//     components,
//     directives,
//   })
  

createApp(App).
use(router).
// use(Vant).
// use(vuetify).
mount('#app')
