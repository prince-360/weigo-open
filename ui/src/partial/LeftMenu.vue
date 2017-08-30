<template>
<aside class="menu is-hidden-mobile">
    <p class="menu-label">
        General
    </p>
    <ul class="menu-list">
        <li><router-link to="/stream" v-bind:class="{'is-active': $route.path==='/stream'}" >Home</router-link></li>
        <li><router-link to="/friends" v-bind:class="{'is-active': $route.path==='/friends'}" >Friends</router-link></li>
    </ul>
    <p class="menu-label">
        Settings
    </p>
    <ul class="menu-list">
        <li><router-link to="/profile" v-bind:class="{'is-active': $route.path==='/profile'}" >Profile</router-link></li>
        <li><a @click="logout">Logout</a></li>
    </ul>
</aside>
</template>

<script>
import auth from '../auth'

export default {
    methods: {
        logout() {
            auth.logout()
        }
    }
}
</script>
