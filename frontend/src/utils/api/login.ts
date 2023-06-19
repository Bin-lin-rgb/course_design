import request from "./request";

export function login(data: any) {
  return request({
    url: `/api/login-api/login`,
    method: "post",
    data: data,
  });
}

export function register(data: any) {
  return request({
    url: `/api/login-api/register`,
    method: "post",
    data: data,
  });
}

export function getUserInfo() {
  return request({
    url: `/api/basic-api/getUserInfo`,
  });
}
export function updateScore(data: any) {
  return request({
    url: `/api/basic-api/setGrade`,
    method: "post",
    data: data,
  });
}
