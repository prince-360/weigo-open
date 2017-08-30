<template>
    <section class="section" v-if="messages.length">
        <div class="columns is-gapless">
            <div class="column is-half is-offset-one-quarter">
                <div class="message is-success" v-for="message in messages">
                    <div class="message-body">{{message}}</div>
                </div>
            </div>
        </div>
    </section>
</template>

<script>

import eventHub from './event'

export default {
    name: 'messages',
    data() {
        return {
            messages: []
        }
    },
    created() {
        eventHub.$on('add-success-message', this.addSuccessMessage)
    },
    beforeDestroy() {
        eventHub.$off('add-success-message', this.addSuccessMessage)
    },
    methods: {
        addSuccessMessage(txt) {
            this.messages.push(txt)
            setTimeout(()=>{this.messages.pop()}, 5000)
        }
    }
}
</script>
