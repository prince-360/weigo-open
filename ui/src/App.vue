<template>
    <main>
        <header-view></header-view>
        <div id="container">
            <div class="grey-content">
                <messages></messages>
                <section class="section">
                    <router-view></router-view>
                </section>
            </div>
        </div>
    </main>
</template>

<script>
import Header from './partial/Header.vue'
import Messages from './Messages.vue'

export default {
    data() {
        return {}
    },
    components: {
        'header-view': Header,
        'messages': Messages,
    }
}
</script>
