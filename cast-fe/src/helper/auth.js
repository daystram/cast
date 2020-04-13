import Cookies from "js-cookie";

export default function auth() {
  return {
    is_authenticated: () => {
      return Cookies.get("user")
    },
    username: () => {
      return Cookies.get("user") ? Cookies.get("user").split("|")[0] : "";
    },
    deauthenticate: () => {
      return Cookies.remove("user")
    }
  }
}
