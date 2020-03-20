const cdn_base = 'https://storage.googleapis.com/cast-uploaded-videos';
const api_base = '${api_base}';

export default function urls() {
  return {
    thumbnail: (hash) => `${cdn_base}/thumbnail/${hash}.jpg`,
    profile: (username) => `${cdn_base}/profile/${username}.png`,
    vod: (hash) => `${cdn_base}/${hash}/manifest.mpd`,
    live: (username) => `${api_base}/live/stream/${username}`,
    upload: () => `${api_base}/p/video/upload`,
    list: () => `${api_base}/video/list`,
    auth_check: () => `${api_base}/auth/check`,
    login: () => `${api_base}/auth/login`,
    signup: () => `${api_base}/auth/signup`,
  }
}
