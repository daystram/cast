import Cookies from "js-cookie";

const COOKIE_KEY = "cast";

export default function auth() {
  return {
    is_authenticated: () => {
      return Cookies.get(COOKIE_KEY)
    },
    username: () => {
      return Cookies.get(COOKIE_KEY) ? Cookies.get(COOKIE_KEY).split("|")[0] : "";
    },
    token: () => {
      return Cookies.get(COOKIE_KEY) ? Cookies.get(COOKIE_KEY).split("|")[1] : "";
    },
    deauthenticate: () => {
      return Cookies.remove(COOKIE_KEY)
    }
  }
}
