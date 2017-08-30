<template>
    <div class="top-header">
        <div class="container">
            <header>
                <nav class="navbar">
                    <div class="navbar-brand">
                        <a class="navbar-item" :href="homeUrl">
                            <img src="/dist/image/logo.png" alt="Weigo logo" width="" height="28">
                        </a>
                        <div class="navbar-burger burger is-hidden-tablet" :class="{'is-active': burgerIsActive}" @click="burgerIsActive = !burgerIsActive" data-target="navMenu">
                            <span></span>
                            <span></span>
                            <span></span>
                        </div>
                    </div>
                    <div class="is-hidden-tablet">
                        <div class="navbar-menu" :class="{'is-active': burgerIsActive}" id="navMenu" v-if="auth.isAuthenticated">
                            <div class="navbar-start">
                                <div class="navbar-item has-dropdown is-hoverable">
                                    <a class="navbar-link  is-active" href="">General</a>
                                    <div class="navbar-dropdown ">
                                        <router-link class="navbar-item" to="/stream">Home</router-link>
                                        <router-link class="navbar-item" to="/friends">Friends</router-link>
                                    </div>
                                </div>
                                <div class="navbar-item has-dropdown is-hoverable">
                                    <a class="navbar-link is-active">Settings</a>
                                    <div class="navbar-dropdown ">
                                        <router-link to="/profile" class="navbar-item">Profile</router-link>
                                        <a class="navbar-item" @click="logout">Logout</a>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="navbar-menu" :class="{'is-active': burgerIsActive}" id="navMenu" v-else>
                            <div class="navbar-start">
                                <router-link class="navbar-item" to="/login">Login</router-link>
                                <router-link class="navbar-item" to="/register">Register</router-link>
                            </div>
                        </div>
                    </div>
                    <div class="navbar-menu" v-if="auth.isAuthenticated">
                        <div class="navbar-end">
                            <span class="navbar-item">
                                {{auth.user.username}}&nbsp;
                            </span>
                        </div>
                    </div>
                    <div class="navbar-menu" v-else>
                        <div class="navbar-end">
                            <router-link class="navbar-item" v-bind:class="{'is-active': $route.path==='/register'}" to="/register">Register</router-link>
                            <router-link class="navbar-item"  v-bind:class="{'is-active': $route.path==='/login'}" to="/login">Login</router-link>
                        </div>
                    </div>
                </nav>
            </header>
        </div>
    </div>
</template>

<script>

import auth, {isAuthenticated, user} from '../auth'
import eventHub from '../event'
import config from '../config'

export default {
    name: 'header-view',
    data() {
        return {
            burgerIsActive: false,
            auth:{
                isAuthenticated: false,
                user: user,
            },
            homeUrl: config.HOME_URL,
        }
    },

    created() {
        eventHub.$on('update-auth', this.updateAuth)
    },
    beforeDestroy() {
        eventHub.$off('update-auth', this.updateAuth)
    },
    methods: {
        updateAuth() {
            this.auth.isAuthenticated = isAuthenticated
            this.auth.user = user
        },
        logout() {
            auth.logout()
        }
    }
}
</script>

<style scoped>
.top-header {
    background-color: white;
}
</style>
