import Vue from 'vue'
import config from '../config'

var _cache = {}

export default {
    async GetById(id) {
        if (!_cache[id]) {
            var resp = await Vue.http.get(config.BASE_URL+'/user/'+id)
            _cache[id] = resp.body
        }
        return _cache[id]
    },

    async GetMe() {
        var resp = await Vue.http.get(config.BASE_URL+'/user')
        return resp.body
    }
}
