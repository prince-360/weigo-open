<template>
    <div class="columns">
        <form class="column is-half is-offset-one-quarter">
            <h1 class="title">Register</h1>
            <br>
            <p class="notification is-danger" v-for="error in errors">{{ error }}</p>
            <div class="field">
                <label class="label">Email</label>
                <div class="control">
                    <input class="input" type="email" v-model="email" placeholder="Email" required>
                </div>
            </div>
            <div class="field">
                <label class="label">Username</label>
                <div class="control">
                    <input class="input" type="text" v-model="username" placeholder="Username">
                </div>
            </div>
            <div class="field">
                <label class="label">Password</label>
                <div class="control">
                    <input class="input" type="password" v-model="password" placeholder="Password" required>
                </div>
            </div>
            <div class="field">
                <label class="label">Password Confirmation</label>
                <div class="control">
                    <input class="input" type="password" v-model="passwordConfirm" placeholder="Password confirmation" required>
                </div>
            </div>
            <button v-on:click.stop.prevent="register" class="button is-primary">Register</button>
            <router-link to="/login" class="button is-hidden-tablet is-pulled-right">Login</router-link>
        </form>
    </div>

</template>

<script>

import config from '../config'
import eventHub from '../event'

export default {
    data() {
        return {
            email: '',
            username: '',
            password: '',
            passwordConfirm: '',
            errors: [],
        }
    },
    methods: {
        checkForm() {
            if (this.email == '')
                this.errors.push('Email is required');
            if (this.password == '')
                this.errors.push('Password is required');
            if (this.username == '')
                this.errors.push('Username is required');
            if (this.passwordConfirm == '')
                this.errors.push('Password Confirmation is required');
            if (this.passwordConfirm != this.password)
                this.errors.push('Password Confirmation mismatch');
            return this.errors.length == 0
        },
        register() {
            this.errors = []
            if (!this.checkForm())
                return
            this.$http.post(config.BASE_URL+'/user/register', JSON.stringify(this.$data)).then( (data) => {
                eventHub.$emit('add-success-message', `Account ${this.$data.email} has been created.`);
                this.$router.push('/login')
            }).catch((resp) => {
                if (resp.body.errors)
                    this.errors = resp.body.errors
            })
        }
    }
}
</script>
