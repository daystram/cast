const base = 'https://storage.googleapis.com/cast-uploaded-videos';
export default function urls() {
  return {
    thumbnail: (hash) => `${base}/thumbnail/${hash}.jpg`,
    profile: (username) => `${base}/profile/${username}.png`,
    vod: (hash) => `${base}/${hash}/manifest.mpd`,
    live: (username) => `/api/v1/live/stream/${username}`,
    upload: () => `/api/v1/p/video/upload`,
    list: () => `/api/v1/video/list`,
  }
}
