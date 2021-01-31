import { RatifyClient, ACCESS_TOKEN } from "@daystram/ratify-client";
import notification from "../helper/notification";

const CLIENT_ID = process.env.REACT_APP_CLIENT_ID;
const ISSUER = process.env.REACT_APP_OAUTH_ISSUER;
const REDIRECT_URI = `${window.location.origin}/callback`;

const authManager = new RatifyClient({
  clientId: CLIENT_ID,
  redirectUri: REDIRECT_URI,
  issuer: ISSUER,
  storage: localStorage,
});

export default function auth() {
  return {
    login() {
      authManager.authorize();
      return null;
    },
    logout() {
      authManager.logout().then(() => {
        notification().destroy();
        window.location.replace("/");
      });
      return null;
    },
    callback() {
      const params = new URLSearchParams(document.location.search);
      const code = params.get("code");
      const state = params.get("state");
      if (!code || !state || !authManager.checkState(state)) {
        window.location.replace("/");
        return null;
      }
      authManager
        .redeemToken(code)
        .then(() => {
          const lastRoute = sessionStorage.getItem("last_route");
          if (lastRoute) {
            sessionStorage.removeItem("last_route");
            window.location.replace({
              path: lastRoute?.toString(),
            });
          } else {
            window.location.replace("/");
          }
        })
        .catch(() => {
          window.location.replace("/");
        });
      return null;
    },
    is_authenticated() {
      return authManager.isAuthenticated();
    },
    user() {
      return authManager.getUser();
    },
    username() {
      return authManager.getUser()?.preferred_username;
    },
    token() {
      return authManager.getToken(ACCESS_TOKEN);
    },
  };
}

export const refreshAuth = function (destinationPath) {
  sessionStorage.setItem("last_route", destinationPath);
  authManager.reset();
  authManager.authorize(true);
};
