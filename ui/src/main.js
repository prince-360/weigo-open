import Vue from 'vue'

import App from './App.vue'
import Home from './Home.vue'
import Login from './pages/Login.vue'
import LoginOauth from './pages/LoginOauth.vue'
import LoginSetup from './pages/LoginSetup.vue'
import Register from './pages/Register.vue'
import Friends from './pages/Friends.vue'
import Stream from './pages/Stream.vue'
import Profile from './pages/Profile.vue'

import VueRouter from 'vue-router'
import VueResource from 'vue-resource'
import auth from './auth'

if (process.env.NODE_ENV == 'production') {
    console.error = () => {}
}

Vue.use(VueRouter)
Vue.use(VueResource)

const routes = [
  { path: '/', redirect: '/stream' },
  { path: '/login', component: Login },
  { path: '/login/auth/:oauthkey', component: LoginOauth, props: true },
  { path: '/login/setup/:oauthkey', component: LoginSetup, props: true },
  { path: '/register', component: Register },
  { path: '/stream', component: Stream },
  { path: '/friends', component: Friends },
  { path: '/profile', component: Profile },
  { path: '/home', component: Home },
]

var router = new VueRouter({
    routes
})
auth.setRouter(router)

Vue.http.interceptors.push((request, next)  => {
    var token = auth.getToken()
    if (token)
        request.headers.set('Authorization', 'Bearer '+token)
    next((response) => {
        if(response.status == 0 ) {
            auth.logout()
            return
        }
        if(response.status == 401 ) {
            auth.logout()
            return
        }
    });
});

auth.isLoggedIn()

new Vue({
  el: '#app',
  render: h => h(App),
  router,
})
