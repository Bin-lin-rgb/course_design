import request from "./request";

export function GetWordList1() {
  return request({
    url: `http://127.0.0.1:8080/basic-api/word/getList1`,
    method: "get"
  });
}