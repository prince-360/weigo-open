<template>
    <div class="avatar-holder" :style="getAvatarPic"></div>
</template>

<script>
import ApiUser from '../api/user'

export default {
    props: ['profile', 'profileId'],
    data() {
        return {
            user: {
                id: 0,
                avatar: '',
            },
        }
    },
    async created() {
        if (this.profile) {
            this.user = this.profile
            return
        }
        if (this.profileId) {
            this.loadProfile()
            return
        }
    },
    methods: {
        async loadProfile() {
            try {
                this.user = await ApiUser.GetById(this.profileId)
            } catch (e) {
                console.error(e)
            }
        },
    },
    computed: {
        getAvatarPic() {
            if (!this.user)
                return `background-image:url(/uploaded/avatar/default.png)`
            if (!this.user.avatar)
                return `background-image:url(/uploaded/avatar/default.png)`
            return `background-image:url(/uploaded/avatar/${this.user.id}-${this.user.avatar}.jpg)`
        }
    },
    watch: {
        'profile': function () {
            if (this.profile)
                this.user = this.profile
        },
        'profileId': function () {
            if (this.profileId)
                this.loadProfile()
        }
    }
}
</script>
