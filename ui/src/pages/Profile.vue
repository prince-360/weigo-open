<template>
    <div class="container">
        <div class="columns">
            <div class="column is-3">
                <left-menu></left-menu>
            </div>
            <div class="column is-6">
                <h2 class="title">Profile</h2>
                <div class="card">
                    <div class="card-content">

                        <avatar :profile="user"></avatar>
                        <br>
                        <div class="field is-grouped" >
                            <div class="file control">
                                <label class="file-label">
                                    <input class="file-input" type="file" v-on:change="updateFileName">
                                    <span class="file-cta">
                                        <span class="file-icon">
                                            <i class="fa fa-upload"></i>
                                        </span>
                                        <span class="file-label">
                                            Media  (Max 2MB)
                                        </span>
                                    </span>
                                </label>
                            </div>
                            <p class="control" v-if="avatarForm.file">
                                <a class="button" @click="uploadAvatar" v-if="avatarForm.progress == null">
                                    Upload
                                </a>
                                <a class="button is-disabled is-warning" v-else>
                                    Uploading ({{avatarForm.progress}}%)
                                </a>
                            </p>
                        </div>

                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-content">
                        <h4 class="subtitle is-4">Informations</h4>
                        <p class="notification is-success" v-show="infoForm.updated">Informations updated.</p>
                        <p class="notification is-danger" v-for="error in infoForm.errors">{{ error }}</p>
                        <div class="field">
                            <label class="label">Username</label>
                            <div class="control">
                                <input class="input" type="text" v-model="infoForm.username" required>
                            </div>
                        </div>
                        <button class="button is-danger" @click="changeInfo()">Update informations</button>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-content">
                        <h4 class="subtitle is-4">Change password</h4>
                        <p class="notification is-success" v-show="passwordForm.updated">Password updated.</p>
                        <p class="notification is-danger" v-for="error in passwordForm.errors">{{ error }}</p>
                        <div class="field">
                            <label class="label">Old password</label>
                            <div class="control">
                                <input class="input" type="password" v-model="passwordForm.oldPassword" required>
                            </div>
                        </div>
                        <div class="field">
                            <label class="label">New password</label>
                            <div class="control">
                                <input class="input" type="password" v-model="passwordForm.newPassword" required>
                            </div>
                        </div>
                        <div class="field">
                            <label class="label">New password Confirmation</label>
                            <div class="control">
                                <input class="input" type="password" v-model="passwordForm.newPasswordConfirm" required>
                            </div>
                        </div>
                        <button class="button is-danger" @click="changePassword()">Update password</button>
                    </div>
                </div>
            </div>
        </div>
    </div>

</template>

<script>
import config from '../config'
import eventHub from '../event'
import {user} from '../auth'
import ApiUser from '../api/user'
import LeftMenu from '../partial/LeftMenu.vue'
import Avatar from '../partial/Avatar.vue'

export default {
    components: {
        'left-menu': LeftMenu,
        'avatar': Avatar,
    },
    data() {
        return {
            user: {},
            seachInput : '',
            avatarForm: {
                file: null,
                fileName: 'empty',
                progress: null,
            },
            passwordForm: {
                oldPassword: "",
                newPassword: "",
                newPasswordConfirm: "",
                errors: [],
                updated: false,
            },
            infoForm: {
                username: '',
                errors: [],
                updated: false,
            },
        }
    },
    async created() {
        try {
            this.user = await ApiUser.GetMe()
            this.infoForm.username = user.username
        } catch (e) {

        }
    },
    methods: {
        async updateFileName(e) {
            if (e.target.files.length == 0) {
                this.avatarForm.fileName = "emtpy"
                return
            }
            this.avatarForm.fileName = e.target.files[0].name
            this.avatarForm.file = e.target.files[0]
        },
        async uploadAvatar() {
            if (!this.avatarForm.file)
                return
            var form = new FormData()
            form.append('avatar', this.avatarForm.file)
            try {
                var configClient = {
                    headers: {'Content-Type': 'multipart/form-data'},
                    progress: (e) => {
                        if (!e.lengthComputable)
                            return
                        this.avatarForm.progress = Math.round((e.loaded / e.total)*100)
                    },
                }
                await this.$http.post(config.BASE_URL+'/user/avatar', form, configClient)
                var userObj = await ApiUser.GetMe()
                this.user = userObj
                this.avatarForm.fileName = "empty"
                this.avatarForm.file = null
            } catch (e) {
                console.error(e)
            }
            this.avatarForm.progress = null
        },
        async changePassword() {
            var data = JSON.stringify(this.passwordForm)
            this.passwordForm.errors = []
            this.passwordForm.updated = false
            try {
                var resp = await this.$http.put(config.BASE_URL+'/user/password', data)
                this.passwordForm.oldPassword = ""
                this.passwordForm.newPassword = ""
                this.passwordForm.newPasswordConfirm = ""
                this.passwordForm.updated = true
                setTimeout(() => {this.passwordForm.updated = false}, 2000)
            } catch (resp) {
                if (resp.body.errors) {
                    this.passwordForm.errors = resp.body.errors
                }
            }
        },
        async changeInfo() {
            var data = JSON.stringify(this.infoForm)
            this.infoForm.errors = []
            this.infoForm.updated = false
            try {
                var resp = await this.$http.put(config.BASE_URL+'/user/info', data)
                this.infoForm.updated = true
                user.username = this.infoForm.username
                eventHub.$emit('update-auth')
                setTimeout(() => {this.infoForm.updated = false}, 2000)
            } catch (resp) {
                if (resp.body.errors) {
                    this.infoForm.errors = resp.body.errors
                }
            }
        },

    },
}
</script>

<style scoped>

</style>
