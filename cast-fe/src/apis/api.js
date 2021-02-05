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
    upload(form) {
      return apiClient.post(`/p/video/upload`, form, {
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "multipart/form-data",
          ...withAuth().headers,
        },
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
};
