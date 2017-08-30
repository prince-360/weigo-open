<template>
    <div>
        <article class="media post-card-stream">
            <figure class="media-left">
                <p class="image is-64x64">
                    <avatar :profile="post.profile"></avatar>
                </p>
            </figure>
            <div class="media-content">
                <div class="content">
                    <p>
                        <strong><post-user-link :profile="post.profile"></post-user-link></strong> <small class="post-date">{{ displayDate(post.created_at) }}</small>
                        <span class="content">{{post.content}}</span>
                    </p>
                    <p v-if="post.medias && post.medias.length > 0">
                        <div class="columns is-multiline is-mobile" style="margin-bottom: 10px;" v-if="post.medias.length == 1">
                            <div class="column" v-for="id in post.medias">
                                <media :mediaId="id"></media>
                            </div>
                        </div>
                        <div class="columns is-multiline is-mobile" style="margin-bottom: 10px;" v-else>
                            <div class="column is-half" v-for="id in post.medias">
                                <media :mediaId="id"></media>
                            </div>
                        </div>
                    </p>
                    <p>
                        <a v-on:click="like()" v-bind:class="{ 'is-liked': post.is_liked }" class="link-social">
                            <i class="fa fa-heart" aria-hidden="true" v-if="post.is_liked"></i>
                            <i class="fa fa-heart-o" aria-hidden="true" v-else></i>
                            &nbsp;<span v-if="post.like_count" class="post-count">{{post.like_count}}</span>
                        </a>
                        <a @click="showComments()" class="link-social">
                            <i class="fa fa-commenting-o" aria-hidden="true"></i>
                            &nbsp;<span v-if="post.comment_count" class="post-count">{{post.comment_count}}</span>
                        </a>
                    </p>
                </div>
            </div>
        </article>

    </div>
</template>

<script>
import Vue from 'vue'
import config from '../../config'
import PostUserLink from './PostUserLink.vue'
import Avatar from '../../partial/Avatar.vue'
import eventHub from '../../event'
import moment from 'moment'
import Media from '../../partial/Media.vue'

var post_cache = {}

export default {
    props: ['post'],
    components: {
        'post-user-link': PostUserLink,
        'avatar': Avatar,
        'media': Media,
    },
    data() {
        return {}
    },
    methods: {
        like() {
            var url = ""
            if (this.post.is_liked)
                url = config.BASE_URL+'/stream/post/'+this.post.id+'/unlike'
            else
                url = config.BASE_URL+'/stream/post/'+this.post.id+'/like'
            this.$http.get(url).then((data)=> {
                this.post.is_liked = !this.post.is_liked
                if (this.post.is_liked) {
                    this.post.like_count += 1
                } else {
                    this.post.like_count -= 1
                }
            }).catch((e)=>{
                console.error(e)
            }).finally(()=>{
            })
        },
        displayDate(t) {
            return moment(t).fromNow();
        },
        showComments() {
            eventHub.$emit('showComments', this.post.id, this.post)
        },
    }

}
</script>

<style scoped>
.is-liked {
    color: hsl(348, 100%, 61%);
}
.content {
    display: block;
    margin-top: 5px;
}
.post-card-stream {
    border: solid 1px #eee;
    margin: 10px 0;
    padding: 10px 20px;
}
.post-date {
    color: #aaa;
    font-size: 0.8rem;
}
.post-count {
    font-size: 0.8em;
}
.link-social {
    margin-right: 20px;
}
</style>
