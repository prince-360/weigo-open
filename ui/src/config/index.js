var BASE_URL_VAR = '//localhost:1314'
var GOOGLE_CLIENT_KEY_VAR='830278333120-uaisbretd5hjb5hra28i4b48la8t4tp7.apps.googleusercontent.com'
var GITHUB_CLIENT_KEY_VAR='fcf6a048247e9555ba8b'

var GOOGLE_REDIRECT_VAR='http://localhost:1314/user/oauth/google/callback'
var GITHUB_REDIRECT_VAR='http://localhost:1314/user/oauth/github/callback'

var HOME_URL_VAR = '/index.dev.html#/'

if (process.env.NODE_ENV == 'production') {
    BASE_URL_VAR = ''
    GOOGLE_CLIENT_KEY_VAR='830278333120-i9nsr2tqegmi41fqcjo3vmik11tcuca6.apps.googleusercontent.com'
    GITHUB_CLIENT_KEY_VAR='95c7d20c53f82bb8a9e7'
    GOOGLE_REDIRECT_VAR='https://weigo.tuxlinuxien.com/user/oauth/google/callback'
    GITHUB_REDIRECT_VAR='https://weigo.tuxlinuxien.com/user/oauth/github/callback'
    HOME_URL_VAR='/#/'
}

export default {
    BASE_URL: BASE_URL_VAR,
    GOOGLE_CLIENT_KEY:GOOGLE_CLIENT_KEY_VAR,
    GITHUB_CLIENT_KEY:GITHUB_CLIENT_KEY_VAR,
    GOOGLE_REDIRECT:GOOGLE_REDIRECT_VAR,
    GITHUB_REDIRECT:GITHUB_REDIRECT_VAR,
    HOME_URL:HOME_URL_VAR,
}
