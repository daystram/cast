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
      `wss://${window.location.hostname}${api_base}${
        token && "/p"
      }/ws/chat/${hash}${token && "?access_token=" + token}`,
    // chat: (hash, token) => `ws://localhost:8080${api_base}${token && '/p'}/ws/chat/${hash}${token && '?access_token=' + token}`,
    notification: (token) =>
      `wss://${window.location.hostname}${api_base}/p/ws/notification?access_token=${token}`,
    // notification: (token) => `ws://localhost:8080${api_base}/p/ws/notification?access_token=${token}`,
    title_check: () => `${api_base}/p/video/check`,
    upload: () => `${api_base}/p/video/upload`,
    edit_video: () => `${api_base}/p/video/edit`,
    edit_window: (open) => `${api_base}/p/live/window?open=${open}`,
    window: () => `${api_base}/p/live/window`,
    delete: () => `${api_base}/p/video/delete`,
    list: () => `${api_base}/video/list`,
    list_authed: () => `${api_base}/p/video/list`,
    search: () => `${api_base}/video/search`,
    auth_check: () => `${api_base}/auth/check`,
    login: () => `${api_base}/auth/login`,
    signup: () => `${api_base}/auth/signup`,
    verify: () => `${api_base}/auth/verify`,
    resend_verify: () => `${api_base}/auth/resend`,
    check_reset: () => `${api_base}/auth/validate_reset`,
    reset_password: () => `${api_base}/auth/forget`,
    update_password: () => `${api_base}/auth/update`,
    user_info: () => `${api_base}/p/user/info`,
    edit_user: () => `${api_base}/p/user/edit`,
    current_hash: () => {
      let split = window.location.href.split("/");
      return split.pop() || split.pop();
    },
  };
}
