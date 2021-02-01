const cdn_base = "/cdn";
const api_base = `${
  process.env.NODE_ENV === "development"
    ? process.env.REACT_APP_DEV_BASE_API
    : ""
}/api/v1`;
const ws_base =
  process.env.NODE_ENV === "development"
    ? `${process.env.REACT_APP_DEV_BASE_WS}/api/v1`
    : `wss://${window.location.hostname}/api/v1`;

export default function urls() {
  return {
    thumbnail: (hash) => `${cdn_base}/thumbnail/${hash}.jpg`,
    profile: (username) => `${cdn_base}/profile/${username}.jpg`,
    vod: (hash) => `${cdn_base}/${hash}/manifest.mpd`,
    download: (hash) => `${cdn_base}/${hash}/video.mp4`,
    live: (username) => `${api_base}/live/stream/${username}`,
    cast_details: () => `${api_base}/video/details`,
    like: () => `${api_base}/p/video/like`,
    subscribe: () => `${api_base}/p/video/subscribe`,
    comment: () => `${api_base}/p/video/comment`,
    chat: (hash, token) =>
      `${ws_base}${token && "/p"}/ws/chat/${hash}${
        token && "?access_token=" + token
      }`,
    notification: (token) =>
      `${ws_base}/p/ws/notification?access_token=${token}`,
    title_check: () => `${api_base}/p/video/check`,
    upload: () => `${api_base}/p/video/upload`,
    edit_video: () => `${api_base}/p/video/edit`,
    edit_window: (open) => `${api_base}/p/live/window?open=${open}`,
    window: () => `${api_base}/p/live/window`,
    delete: () => `${api_base}/p/video/delete`,
    list: () => `${api_base}/video/list`,
    list_authed: () => `${api_base}/p/video/list`,
    search: () => `${api_base}/video/search`,
    register: () => `${api_base}/p/auth/check`,
    user_info: () => `${api_base}/p/user/info`,
    current_hash: () => {
      let split = window.location.href.split("/");
      return split.pop() || split.pop();
    },
  };
}
