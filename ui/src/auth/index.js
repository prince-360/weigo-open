import config from '../config'
import Header from '../partial/Header.vue'
import eventHub from '../event'
import Vue from 'vue'
import VueRouter from 'vue-router'

var router = null

export var isAuthenticated = false;
export var user = {id: 0, username: '', avatar: '', email: ''}

export default {

    login(token) {
        localStorage.setItem('token',token);

        Vue.http.get(config.BASE_URL+"/user",{
            headers: {
                'Authorization': 'Bearer '+token,
            },
        }).then((data) => {
            isAuthenticated = true
            user = data.body
            eventHub.$emit('update-auth')
            router.push('/stream')
        }).catch((err) => {
            console.error(err)
        })
    },

    isLoggedIn() {
        var token = this.getToken()
        if (token == "0.0.0")
            return
        Vue.http.get(config.BASE_URL+"/user",{
            headers: {
                'Authorization': 'Bearer '+token,
            },
        }).then((data) => {
            isAuthenticated = true
            user = data.body
            eventHub.$emit('update-auth')
        }).catch((err) => {
            if (err.status == 0) {
                localStorage.removeItem('token')
                router.push('/login')
            }
            if (err.status == 401) {
                localStorage.removeItem('token')
                router.push('/login')
            }
        })
    },

    setRouter(_router) {
        router = _router
    },

    getToken() {
        var token = localStorage.getItem('token')
        if (!token) {
            return "0.0.0"
        }
        return token
    },

    logout() {
        isAuthenticated = false
        user = {}
        eventHub.$emit('update-auth')
        localStorage.removeItem('token')
        router.push('/login')
    }
}
