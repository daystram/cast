import axios from "axios";
import { ACCESS_TOKEN } from "@daystram/ratify-client";
import { authManager, refreshAuth } from "../helper/auth";

const baseAPI = `${
  process.env.NODE_ENV === "development"
    ? process.env.REACT_APP_DEV_BASE_API
    : ""
}/api/v1`;
const baseCDN = "https://cdn.daystram.com/cast";
const baseWS =
  process.env.NODE_ENV === "development"
    ? `${process.env.REACT_APP_DEV_BASE_WS}/api/v1`
    : `wss://${window.location.hostname}/api/v1`;

const apiClient = axios.create({
  baseURL: baseAPI,
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
  live: {
    stream(username) {
      return `${baseAPI}/live/stream/${username}`;
    },
    window: {
      status() {
        return apiClient.get(`/p/live/window`, withAuth());
      },
      set(open) {
        return apiClient.put(`/p/live/window?open=${open}`, {}, withAuth());
      },
    },
  },
  ws: {
    chat(hash) {
      const token = authManager.getToken(ACCESS_TOKEN);
      return `${baseWS}${token && "/p"}/ws/chat/${hash}${
        token && "?access_token=" + token
      }`;
    },
    notification() {
      return `${baseWS}/p/ws/notification?access_token=${authManager.getToken(
        ACCESS_TOKEN
      )}`;
    },
  },
  cast: {
    list(params) {
      return apiClient.get(`/video/list`, { params });
    },
    listCurated(params) {
      return apiClient.get(`/p/video/list`, { params, ...withAuth() });
    },
    search(params) {
      return apiClient.get(`/video/search`, { params });
    },
    detail(params) {
      return apiClient.get(`/video/details`, { params });
    },
    upload(form, onUploadProgress) {
      return apiClient.post(`/p/video/upload`, form, {
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "multipart/form-data",
          ...withAuth().headers,
        },
        onUploadProgress,
      });
    },
    edit(form) {
      return apiClient.put(`/p/video/edit`, form, {
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "multipart/form-data",
          ...withAuth().headers,
        },
      });
    },
    remove(hash) {
      return apiClient.delete(`/p/video/delete`, {
        params: { hash },
        ...withAuth(),
      });
    },
    like(data) {
      return apiClient.post(`/p/video/like`, data, withAuth());
    },
    comment(data) {
      return apiClient.post(`/p/video/comment`, data, withAuth());
    },
    titleCheck(title) {
      return apiClient.get(`/p/video/check`, {
        params: { title },
        ...withAuth(),
      });
    },
  },
  user: {
    detail() {
      return apiClient.get(`/p/user/info`, withAuth());
    },
    subscribe(data) {
      return apiClient.post(`/p/video/subscribe`, data, withAuth());
    },
  },
  auth: {
    register(idToken) {
      return apiClient.post(`/p/auth/check`, { id_token: idToken }, withAuth());
    },
  },
};
