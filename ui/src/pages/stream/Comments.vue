<template>
    <div class="modal" v-bind:class="{ 'is-active': isActive }">
        <div class="modal-background" @click="close()"></div>
        <div class="modal-card">
            <section class="modal-card-body">
                <p>
                    <strong><post-user-link :profile="post.profile"></post-user-link></strong> <small class="post-date">{{ displayDate(post.created_at) }}</small>
                    <span class="content">{{post.content}}</span>
                </p>
                <br>
                <div class="field">
                    <div class="control">
                        <textarea class="textarea" v-model="commentContent" placeholder="Comment this post">{{commentContent}}</textarea>
                    </div>
                </div>
                <button v-on:click="sendComment" class="button is-primary">Send</button>
                <br>
                <br>
                <div class="h4">Comments:</div>
                <p v-for="comment in comments" class="comment-from">
                    <strong><post-user-link :profile="comment.profile"></post-user-link></strong> <small class="post-date">{{ displayDate(comment.created_at) }}</small>
                    <span class="content">{{comment.content}}</span>
                </p>
            </section>
        </div>
        <button class="modal-close is-large" @click="close()"></button>
    </div>
</template>

<script>
import eventHub from '../../event'
import config from '../../config'
import PostUserLink from './PostUserLink.vue'
import moment from 'moment'

export default {
    components: {
        'post-user-link': PostUserLink,
    },
    data() {
        return {
            isActive: false,
            post: {},
            postId: 0,
            comments: [],
            commentContent: "",
        }
    },
    created() {
        eventHub.$on('showComments', this.showComments)
    },
    destroyed() {
        eventHub.$off('showComments', this.showComments)
    },
    methods: {
        close() {
            this.isActive = false
        },
        async sendComment() {
            try {
                var postdata = JSON.stringify({content: this.commentContent})
                await this.$http.post(config.BASE_URL+'/stream/post/'+this.postId+'/comment',postdata)
                this.getComments()
                this.commentContent = ""
                this.post.comment_count += 1
            } catch (e) {
                console.error(e)
            }
        },
        async getComments() {
            try {
                var resp = await this.$http.get(config.BASE_URL+'/stream/post/'+this.postId+'/comment')
                this.comments = resp.body
            } catch (e) {
                console.error(e)
            }
        },
        displayDate(t) {
            return moment(t).fromNow();
        },
        showComments(postId, postData) {
            this.isActive = true
            this.post = postData
            this.postId = postId
            this.commentContent = ""
            this.getComments()
        }
    },

}
</script>

<style scoped>
.comment-from {
    padding: 5px 0;
    border-bottom: solid 1px #eee;
}
.content {
    display: block;
    margin-top: 5px;
}
.post-date {
    color: #aaa;
    font-size: 0.8rem;
}
</style>
