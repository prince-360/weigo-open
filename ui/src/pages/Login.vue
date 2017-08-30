<template>
    <div>
        <div class="columns">
            <form class="column is-half is-offset-one-quarter">
                <h1 class="title">Login</h1>
                <br>
                <p class="notification is-danger" v-for="error in errors">{{ error }}</p>
                <div class="field">
                    <label class="label">Username</label>
                    <div class="control">
                        <input class="input" type="text" v-model="username" placeholder="Username">
                    </div>
                </div>
                <div class="field">
                    <label class="label">Password</label>
                    <div class="control">
                        <input class="input" type="password" v-model="password" placeholder="Password">
                    </div>
                </div>
                <button v-on:click.stop.prevent="login" class="button is-primary">Login</button>
                <router-link to="/register" class="button is-hidden-tablet is-pulled-right">Register</router-link>
            </form>
        </div>
        <div class="columns">
            <p class="column has-text-centered">
                <!-- <a class="button" href="https://accounts.google.com/o/oauth2/v2/auth?client_id=830278333120-i9nsr2tqegmi41fqcjo3vmik11tcuca6.apps.googleusercontent.com&redirect_uri=https://weigo.tuxlinuxien.com/user/oauth/google/callback&response_type=code&scope=email%20profile"> -->
                <a class="button" :href="googleKey">
                    Login with Google&nbsp;<i class="fa fa-google" aria-hidden="true"></i>
                </a>
            </p>
        </div>
        <div class="columns">
            <p class="column has-text-centered">
                <a class="button" :href="githubKey">
                    Login with Github&nbsp;<i class="fa fa-github" aria-hidden="true"></i>
                </a>
            </p>
        </div>
    </div>
</template>

<script>
import config from '../config'
import auth from '../auth'
import eventHub from '../event'

export default {
    data() {
        return {
            username: '',
            password: '',
            errors: [],
            githubKey: `http://github.com/login/oauth/authorize?client_id=${config.GITHUB_CLIENT_KEY}&redirect_uri=${config.GITHUB_REDIRECT}&scope=user:email`,
            googleKey: `https://accounts.google.com/o/oauth2/v2/auth?client_id=${config.GOOGLE_CLIENT_KEY}&redirect_uri=${config.GOOGLE_REDIRECT}&response_type=code&scope=email%20profile`,
        }
    },
    methods: {
        login() {
            this.$http.post(config.BASE_URL+'/user/login', JSON.stringify(this.$data)).then( (data) => {
                auth.login(data.body.token, this.$http, null)
            }).catch((resp) => {
                if (resp.status == 404) {
                    this.errors = ["This user doesn't exist."]
                }
                if (resp.status == 401) {
                    this.errors = ["Invalid password."]
                }
            })
        }
    }
}
</script>
