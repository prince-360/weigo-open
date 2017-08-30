<template>
    <div class="columns">
        <form class="column is-half is-offset-one-quarter">
            <h1 class="title">Setup account</h1>
            <br>
            <p class="notification is-danger" v-for="error in errors">{{ error }}</p>
            <div class="field">
                <label class="label">Username</label>
                <div class="control">
                    <input class="input" type="text" v-model="username" placeholder="Username">
                </div>
            </div>
            <button v-on:click.stop.prevent="register" class="button is-primary">Register</button>
        </form>
    </div>
</template>

<script>
import config from '../config'
import auth from '../auth'
import eventHub from '../event'

export default {
    props: ['oauthkey'],
    data() {
        return {
            errors: [],
            username: '',
        }
    },
    async created() {
        try {
            let resp = await this.$http.get(config.BASE_URL+"/user/login/setup/"+this.oauthkey)
            this.username = resp.body.username
        } catch (e) {
            console.error(e)
        }
    },
    methods: {
        async register() {
            try {
                let data = {username: this.username}
                let resp = await this.$http.post(config.BASE_URL+"/user/login/setup/"+this.oauthkey, data)
                this.$router.push('/login/auth/'+this.oauthkey)
            } catch (e) {
                if (e.status == 400) {
                    this.errors = e.body.errors
                }
            }
        }
    }
}
</script>
