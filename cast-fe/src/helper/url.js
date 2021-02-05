const api_base = `${
  process.env.NODE_ENV === "development"
    ? process.env.REACT_APP_DEV_BASE_API
    : ""
}/api/v1`;

export default function urls() {
  return {
    edit_window: (open) => `${api_base}/p/live/window?open=${open}`,
    window: () => `${api_base}/p/live/window`,

    register: () => `${api_base}/p/auth/check`,

    user_info: () => `${api_base}/p/user/info`,
    subscribe: () => `${api_base}/p/video/subscribe`,

    current_hash: () => {
      let split = window.location.href.split("/");
      return split.pop() || split.pop();
    },
  };
}
