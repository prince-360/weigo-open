<template>
    <div class="columns">
        <div class="column">
            <p>login you in...</p>
        </div>
    </div>
</template>

<script>
import config from '../config'
import auth from '../auth'
import eventHub from '../event'

export default {
    props: ['oauthkey'],
    async created() {
        try {
            let resp = await this.$http.get(config.BASE_URL+"/user/login/"+this.oauthkey)
            auth.login(resp.body.token)
        } catch (e) {
            console.error(e)
        }
    }
}
</script>
