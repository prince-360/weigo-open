<template>
    <div class="container">
        <div class="columns">
            <div class="column is-3">
                <left-menu></left-menu>
            </div>
            <div class="column is-6">
                <h2 class="title">Friends</h2>
                <div class="field">
                    <div class="control has-icons-left has-icons-right">
                        <input class="input" placeholder="Looking for someone?" v-model="seachInput" v-on:keyup.enter="searchUser">
                        <span class="icon is-small is-left">
                            <i class="fa fa-user"></i>
                        </span>
                        <span class="icon is-small is-right">
                            <i class="fa fa-search"></i>
                        </span>
                    </div>
                </div>
                <div v-show="results.length > 0">
                    <h4 class="h4">Results:</h4>
                    <p class="search-user-result" v-for="result in results">
                        <post-user-link :profile="result"></post-user-link>
                        <a class="button is-white is-small is-pulled-right" @click="addFriend(result.id)">
                            <span class="icon is-small">
                                <i class="fa fa-user-plus"></i>
                            </span>
                        </a>
                    </p>
                    <hr>
                </div>
                <p class="search-user-result" v-for="friend in friends">
                    <post-user-link :profile="friend"></post-user-link>
                    <a class="button is-white is-small is-pulled-right" @click="removeFriend(friend.id)">
                        <span class="icon is-small">
                            <i class="fa fa-user-times"></i>
                        </span>
                    </a>
                </p>
            </div>
        </div>
    </div>

</template>

<script>
import config from '../config'
import LeftMenu from '../partial/LeftMenu.vue'
import PostUserLink from './stream/PostUserLink.vue'

export default {
    components: {
        'left-menu': LeftMenu,
        'post-user-link': PostUserLink,
    },
    data() {
        return {
            seachInput : '',
            results: [],
            friends: [],
        }
    },
    created() {
        this.listFriends()
        //this.searchUserDebounce =  _.debounce(this.searchUser, 5000).bind(this)
    },
    methods: {
        async searchUser() {
            var content = JSON.stringify({'contains': this.seachInput})
            try {
                var resp = await this.$http.post(config.BASE_URL+'/user/search',content)
                this.results = resp.body
            } catch (e) {
                console.error(e)
            }
        },
        async listFriends() {
            var url = config.BASE_URL+'/user/friendship'
            try {
                var resp = await this.$http.get(url)
                this.friends = resp.body
            } catch (e) {
                console.error(e)
            }
        },
        async addFriend(profile_id) {
            var url = config.BASE_URL+'/user/friendship/'+profile_id
            try {
                await this.$http.post(url, "")
                this.listFriends()
            } catch (e) {
                console.error(e)
            }
        },
        async removeFriend(profile_id) {
            var url = config.BASE_URL+'/user/friendship/'+profile_id
            try {
                await this.$http.delete(url, "")
                this.listFriends()
            } catch (e) {
                console.error(e)
            }
        }
    },
    watch: {

    },
}
</script>

<style scoped>
.search-user-result {
    background-color: white;
    padding: 10px;
    border-radius: 2px;
    margin-top: 5px;
}
</style>
