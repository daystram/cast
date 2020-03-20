const base = 'https://storage.googleapis.com/cast-uploaded-videos';
export default function urls() {
  return {
    thumbnail: (hash) => `${base}/thumbnail/${hash}.jpg`,
    profile: (username) => `${base}/profile/${username}.png`,
    vod: (hash) => `${base}/${hash}/manifest.mpd`,
    live: (username) => `http://localhost:8080/api/v1/live/stream/${username}`,
    upload: () => `/p/video/upload`,
    list: () => `/video/list`,
  }
}
