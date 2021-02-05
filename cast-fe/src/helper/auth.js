import { RatifyClient, ID_TOKEN } from "@daystram/ratify-client";
import notification from "../helper/notification";
import api from "../apis/api";

const CLIENT_ID = process.env.REACT_APP_CLIENT_ID;
const ISSUER = process.env.REACT_APP_OAUTH_ISSUER;
const REDIRECT_URI = `${window.location.origin}/callback`;

export const authManager = new RatifyClient({
  clientId: CLIENT_ID,
  redirectUri: REDIRECT_URI,
  issuer: ISSUER,
  storage: localStorage,
});

export const login = function () {
  authManager.authorize();
  return null;
};

export const logout = function () {
  authManager.logout().then(() => {
    notification().destroy();
    window.location.replace("/");
  });
  return null;
};

export const callback = function () {
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
      api.auth
        .register(authManager.getToken(ID_TOKEN))
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
};

export const refreshAuth = function (destinationPath) {
  sessionStorage.setItem("last_route", destinationPath);
  authManager.reset();
  authManager.authorize(true);
};
