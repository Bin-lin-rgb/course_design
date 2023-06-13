import request from "./request";

export function GetClassification() {
  return request({
    url: `/api/v1/community`,
  });
}

export function GetPostList() {
  return request({
    url: `/api/v1/postlist`,
  });
}

export function GetPostDetails(data:any) {
  return request({
    url: `/api/v1/post/${data.id}`,
  });
}

export function GetPostListWithTime(data: any) {
  return request({
    url: `/api/v1/postlist/${data.nextTime}`,
  });
}

export function PostArticle(data: any) {
  return request({
    url: `/api/v1/post`,
    method: "post",
    data: data,
  });
}

export function GetPostListByUserId() {
  return request({
    url: `/api/v1/user/bloglist`,
  });
}

export function DeleteBlog(data:any) {
  return request({
    url: `/api/v1/user/delete/${data.id}`,
    method: "delete",
  });
}

export function UpdateBlog(data:any) {
  return request({
    url: `/api/v1/user/updateList/${data.id}`,
    method: "delete",
  });
}
