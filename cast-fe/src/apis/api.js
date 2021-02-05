import axios from "axios";
import { ACCESS_TOKEN } from "@daystram/ratify-client";
import { authManager, refreshAuth } from "../helper/auth";

const baseCDN = "https://cdn.daystram.com/cast";
const baseWS =
  process.env.NODE_ENV === "development"
    ? `${process.env.REACT_APP_DEV_BASE_WS}/api/v1`
    : `wss://${window.location.hostname}/api/v1`;
const apiClient = axios.create({
  baseURL: `${
    process.env.NODE_ENV === "development"
      ? process.env.VUE_APP_DEV_BASE_API
      : ""
  }/api/v1`,
});

apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response.status === 401) {
      refreshAuth(window.location.href);
    }
    return Promise.reject(error);
  }
);

const withAuth = () => ({
  headers: {
    Authorization: `Bearer ${authManager.getToken(ACCESS_TOKEN)}`,
  },
});

export default {
  cdn: {
    thumbnail(hash) {
      return `${baseCDN}/thumbnail/${hash}.jpg`;
    },
    vod(hash) {
      return `${baseCDN}/video/${hash}/manifest.mpd`;
    },
    download(hash) {
      return `${baseCDN}/video/${hash}/video.mp4`;
    },
  },
};
