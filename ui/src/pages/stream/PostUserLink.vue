<template>
    <a v-if="user.username != ''" href="">@{{user.username}}</a>
</template>

<script>
import config from '../../config'
import Vue from 'vue'
import ApiUser from '../../api/user'

export default {
    props: ['profile'],
    data() {
        return {
            user: {
                username: '',
            },
        }
    },
    created() {
        if (this.profile && this.profile.username)
            this.user.username = this.profile.username
    },
    watch: {
        'profile': function () {
            if (this.profile && this.profile.username)
                this.user.username = this.profile.username
        }
    }
}
</script>
