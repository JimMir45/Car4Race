import fs from 'fs/promises';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

class CookieManager {
  constructor() {
    this.cookiePath = path.join(__dirname, '..', 'data', 'cookie.json');
  }

  async saveCookie(cookie) {
    const dataDir = path.join(__dirname, '..', 'data');
    await fs.mkdir(dataDir, { recursive: true });
    
    // 从cookie中提取csrf token (bili_jct)
    const csrfMatch = cookie.match(/bili_jct=([^;]+)/);
    const csrf = csrfMatch ? csrfMatch[1] : '';
    
    const cookieData = {
      cookie: cookie,
      csrf: csrf,
      updatedAt: new Date().toISOString()
    };
    
    await fs.writeFile(this.cookiePath, JSON.stringify(cookieData, null, 2));
    return cookieData;
  }

  async getCookie() {
    try {
      const data = await fs.readFile(this.cookiePath, 'utf-8');
      return JSON.parse(data);
    } catch (error) {
      return { cookie: '', csrf: '', updatedAt: null };
    }
  }

  async validateCookie(cookie) {
    // 检查cookie是否包含必要的字段
    const requiredFields = ['SESSDATA', 'bili_jct', 'DedeUserID'];
    for (const field of requiredFields) {
      if (!cookie.includes(field)) {
        return { valid: false, missing: field };
      }
    }
    return { valid: true };
  }
}

export default new CookieManager();