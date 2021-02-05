import { RatifyClient, ACCESS_TOKEN, ID_TOKEN } from "@daystram/ratify-client";
import notification from "../helper/notification";
import urls from "./url";
import axios from "axios";

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
          axios
            .post(
              urls().register(),
              { id_token: authManager.getToken(ID_TOKEN) },
              { headers: { Authorization: `Bearer ${auth().token()}` } }
            )
            .then(() => {
              const lastRoute = sessionStorage.getItem("last_route");
              if (lastRoute) {
                sessionStorage.removeItem("last_route");
                window.location.replace(lastRoute);
              } else {
                window.location.replace("/");
              }
            })
            .catch(() =>
              authManager.logout().then(() => window.location.replace("/"))
            );
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

const refreshAuth = function (destinationPath) {
  sessionStorage.setItem("last_route", destinationPath);
  authManager.reset();
  authManager.authorize(true);
};

export { authManager, refreshAuth };
