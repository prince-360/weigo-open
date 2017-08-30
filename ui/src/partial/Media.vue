<template>
    <div class="media-holder" :style="getMediaPic"></div>
</template>

<script>
export default {
    props: {
        mediaUrl: {
            default: '',
        },
        mediaId: {
            default: '',
        },
    },
    computed: {
        getMediaPic() {
            if (this.mediaUrl)
                return `background-image:url(${this.mediaUrl})`
            if (this.mediaId)
                return `background-image:url(/uploaded/media/${this.mediaId})`
            return `background-image:url(/uploaded/media/default.png)`
        }
    },
}
</script>
