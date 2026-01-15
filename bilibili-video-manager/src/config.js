export const CONFIG = {
  // B站API配置
  API: {
    BASE_URL: 'https://member.bilibili.com',
    ARCHIVES_ENDPOINT: '/x/web/archives',
    DEFAULT_PAGE_SIZE: 50,
    HIDE_VIDEO_ENDPOINT: '/x/vu/web/edit/visibility'  // 隐藏视频接口
  },
  
  // 请求头配置
  HEADERS: {
    'accept': 'application/json, text/javascript, */*; q=0.01',
    'accept-language': 'zh-CN,zh;q=0.9',
    'priority': 'u=1, i',
    'referer': 'https://member.bilibili.com/platform/upload-manager/article',
    'sec-ch-ua': '"Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': '"Windows"',
    'sec-fetch-dest': 'empty',
    'sec-fetch-mode': 'cors',
    'sec-fetch-site': 'same-origin',
    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36',
    'x-requested-with': 'XMLHttpRequest'
  }
};