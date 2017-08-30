<template>
    <div class="container stream-page">
        <div class="columns">
            <div class="column is-3">
                <left-menu></left-menu>
            </div>
            <div class="column is-6">
                <comments></comments>
                <h2 class="title">Home</h2>
                <div class="columns is-mobile" v-show="newPostFiles.length > 0" style="margin-bottom: 10px;">
                    <div class="column is-one-quarter" v-for="id in newPostFiles">
                        <media :mediaId="id"></media>
                    </div>
                </div>
                <div class="field">
                    <div class="control">
                        <textarea class="textarea" id="post" placeholder="What's new" v-model="newPostData"></textarea>
                    </div>
                </div>
                <div class="field is-grouped" v-show="newPostFiles.length < 4">
                    <div class="file control">
                        <label class="file-label">
                            <input class="file-input" type="file" v-on:change="updateFileName">
                            <span class="file-cta">
                                <span class="file-icon">
                                    <i class="fa fa-upload"></i>
                                </span>
                                <span class="file-label">
                                    Media  (Max 5MB)
                                </span>
                            </span>
                        </label>
                    </div>
                    <p class="control" v-if="mediaForm.file">
                        <a class="button" @click="uploadMedia" v-if="mediaForm.progress == null">
                            Upload
                        </a>
                        <a class="button is-disabled is-warning" v-else>
                            Uploading ({{mediaForm.progress}}%)
                        </a>
                    </p>
                </div>
                <a class="button" @click="sendPost" v-if="!newPostIsSending">Send</a>
                <a class="button is-loading" v-else></a>
                <div id="postlist">
                    <post v-for="post in posts" :post="post"></post>
                </div>
                <p class="has-text-centered" v-show="downloadedItems == 10">
                    <a class="button" v-if="!isLoading" @click="loadPosts()">Load more</a>
                    <button class="button is-loading" v-else>Loading</button>
                </p>
            </div>
        </div>
    </div>

</template>

<script>
import config from '../config'
import Vue from 'vue'
import Post from './stream/Post.vue'
import Comments from './stream/Comments.vue'
import LeftMenu from '../partial/LeftMenu.vue'
import Media from '../partial/Media.vue'

export default {
    components: {
        'post': Post,
        'comments': Comments,
        'left-menu': LeftMenu,
        'media': Media,
    },
    data() {
        return {
            posts: [],
            newPostData: '',
            newPostFiles: [],
            newPostIsSending: false,
            isLoading: true,
            downloadedItems: 10,
            lastItemId: null,
            mediaForm: {
                fileName: 'emtpy',
                file: null,
                progress: null,
            }

        }
    },
    created() {
        this.loadPosts()
    },
    methods: {
        async updateFileName(e) {
            if (e.target.files.length == 0) {
                this.mediaForm.fileName = "emtpy"
                return
            }
            this.mediaForm.fileName = e.target.files[0].name
            this.mediaForm.file = e.target.files[0]
        },
        async uploadMedia() {
            if (!this.mediaForm.file)
                return
            var form = new FormData()
            form.append('media', this.mediaForm.file)
            try {
                var configClient = {
                    headers: {'Content-Type': 'multipart/form-data'},
                    progress: (e) => {
                        if (!e.lengthComputable)
                            return
                        this.mediaForm.progress = Math.round((e.loaded / e.total)*100)
                    },
                }
                var resp = await this.$http.post(config.BASE_URL+'/media', form, configClient)
                this.mediaForm.fileName = "empty"
                this.mediaForm.file = null
                this.newPostFiles.push(resp.body.id)
                console.log(this.newPostFiles)
            } catch (e) {
                console.error(e)
            }
            this.mediaForm.progress = null
        },
        async loadPosts() {
            var url = config.BASE_URL+"/stream/post"
            if (this.posts.length > 0)
                url += "?from=" + this.posts[this.posts.length-1].id
            try {
                this.isLoading = true
                var resp = await Vue.http.get(url)
                this.downloadedItems = resp.body.length
                this.posts = this.posts.concat(resp.body)
            } catch(e) {
                console.error(e)
            }
            this.isLoading = false
        },
        sendPost() {
            if (this.newPostData.trim() == "" && this.newPostFiles.length == 0)
                return
            this.newPostIsSending = true
            let data = JSON.stringify({content: this.newPostData, medias: this.newPostFiles})
            this.$http.post(config.BASE_URL+"/stream/post",data).then((data) => {
                this.newPostData = ""
                this.newPostFiles = []
                this.posts.unshift(data.body)
            }).catch((err) => {
                console.error(err)
            })
            this.newPostIsSending =  false
        },
    },
}
</script>
