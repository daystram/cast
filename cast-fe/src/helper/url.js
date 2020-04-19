const cdn_base = '/cdn';
const api_base = '/api/v1';

export default function urls() {
  return {
    thumbnail: (hash) => `${cdn_base}/thumbnail/${hash}.jpg`,
    profile: (username) => `${cdn_base}/profile/${username}.jpg`,
    vod: (hash) => `${cdn_base}/${hash}/manifest.mpd`,
    download: (hash) => `${cdn_base}/${hash}/video.mp4`,
    live: (username) => `${api_base}/live/stream/${username}`,
    cast_details: () => `${api_base}/video/details`,
    like: () => `${api_base}/p/video/like`,
    comment: () => `${api_base}/p/video/comment`,
    // chat: (hash, token) => `wss://${window.location.hostname}${api_base}${token && '/p'}/chat/${hash}${token && '?access_token=$' + token}`,
    chat: (hash, token) => `ws://localhost:8080${api_base}${token && '/p'}/chat/${hash}${token && '?access_token=$' + token}`,
    title_check: () => `${api_base}/p/video/check`,
    upload: () => `${api_base}/p/video/upload`,
    edit_video: () => `${api_base}/p/video/edit`,
    delete: () => `${api_base}/p/video/delete`,
    list: () => `${api_base}/video/list`,
    search: () => `${api_base}/video/search`,
    auth_check: () => `${api_base}/auth/check`,
    login: () => `${api_base}/auth/login`,
    signup: () => `${api_base}/auth/signup`,
    verify: () => `${api_base}/auth/verify`,
    resend_verify: () => `${api_base}/auth/resend`,
    user_info: () => `${api_base}/p/user/info`,
    edit_user: () => `${api_base}/p/user/edit`,
  }
}
