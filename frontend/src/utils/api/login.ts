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
    url: `/api/v1/user`,
  });
}

// export function updateProfile(data: any) {
//   return request({
//     url: `/api/v1/user/profile`,
//     method: "put",
//     data: data,
//   });
// }
