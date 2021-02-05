const api_base = `${
  process.env.NODE_ENV === "development"
    ? process.env.REACT_APP_DEV_BASE_API
    : ""
}/api/v1`;

export default function urls() {
  return {
    user_info: () => `${api_base}/p/user/info`,
    subscribe: () => `${api_base}/p/video/subscribe`,

    current_hash: () => {
      let split = window.location.href.split("/");
      return split.pop() || split.pop();
    },
  };
}
