import fetch from 'node-fetch';
import { CONFIG } from './config.js';
import cookieManager from './cookie-manager.js';

class VideoManager {
  
  // 隐藏单个视频（仅自己可见）
  async hideVideo(aid) {
    const cookieData = await cookieManager.getCookie();
    if (!cookieData.cookie || !cookieData.csrf) {
      throw new Error('Cookie or CSRF token not configured');
    }

    const url = `${CONFIG.API.BASE_URL}${CONFIG.API.HIDE_VIDEO_ENDPOINT}?csrf=${cookieData.csrf}`;
    
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'accept': 'application/json, text/plain, */*',
        'accept-language': 'zh-CN,zh;q=0.9',
        'content-type': 'application/x-www-form-urlencoded; charset=UTF-8',
        'origin': 'https://member.bilibili.com',
        'referer': 'https://member.bilibili.com/platform/upload-manager/article',
        'Cookie': cookieData.cookie,
        'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36'
      },
      body: `aid=${aid}&is_only_self=1&csrf=${cookieData.csrf}`
    });

    const result = await response.json();
    return result;
  }

  // 显示视频（公开）
  async showVideo(aid) {
    const cookieData = await cookieManager.getCookie();
    if (!cookieData.cookie || !cookieData.csrf) {
      throw new Error('Cookie or CSRF token not configured');
    }

    const url = `${CONFIG.API.BASE_URL}${CONFIG.API.HIDE_VIDEO_ENDPOINT}?csrf=${cookieData.csrf}`;
    
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'accept': 'application/json, text/plain, */*',
        'accept-language': 'zh-CN,zh;q=0.9',
        'content-type': 'application/x-www-form-urlencoded; charset=UTF-8',
        'origin': 'https://member.bilibili.com',
        'referer': 'https://member.bilibili.com/platform/upload-manager/article',
        'Cookie': cookieData.cookie,
        'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36'
      },
      body: `aid=${aid}&is_only_self=0&csrf=${cookieData.csrf}`
    });

    const result = await response.json();
    return result;
  }

  // 批量隐藏包含特定标签的视频
  async hideVideosByTag(targetTag, videos) {
    const results = {
      success: [],
      failed: [],
      skipped: []
    };

    for (const video of videos) {
      // 检查视频是否包含目标标签
      if (video.tags && video.tags.some(tag => 
        tag.toLowerCase().includes(targetTag.toLowerCase()))) {
        
        // 跳过充电视频
        if (video.isCharging) {
          results.skipped.push({
            aid: video.aid,
            title: video.title,
            reason: '充电视频不支持隐藏'
          });
          continue;
        }
        
        // 跳过已经隐藏的视频
        if (video.isOnlySelf) {
          results.skipped.push({
            aid: video.aid,
            title: video.title,
            reason: 'Already hidden'
          });
          continue;
        }

        try {
          console.log(`正在隐藏视频: ${video.title} (aid: ${video.aid})`);
          const result = await this.hideVideo(video.aid);
          
          if (result.code === 0) {
            results.success.push({
              aid: video.aid,
              title: video.title
            });
          } else {
            results.failed.push({
              aid: video.aid,
              title: video.title,
              error: result.message || 'Unknown error'
            });
          }
          
          // 添加延迟，避免请求过快
          await new Promise(resolve => setTimeout(resolve, 500));
        } catch (error) {
          results.failed.push({
            aid: video.aid,
            title: video.title,
            error: error.message
          });
        }
      }
    }

    return results;
  }

  // 批量显示包含特定标签的视频
  async showVideosByTag(targetTag, videos) {
    const results = {
      success: [],
      failed: [],
      skipped: []
    };

    for (const video of videos) {
      // 检查视频是否包含目标标签
      if (video.tags && video.tags.some(tag => 
        tag.toLowerCase().includes(targetTag.toLowerCase()))) {
        
        // 跳过充电视频
        if (video.isCharging) {
          results.skipped.push({
            aid: video.aid,
            title: video.title,
            reason: '充电视频不支持此操作'
          });
          continue;
        }
        
        // 跳过已经公开的视频
        if (!video.isOnlySelf) {
          results.skipped.push({
            aid: video.aid,
            title: video.title,
            reason: 'Already public'
          });
          continue;
        }

        try {
          console.log(`正在公开视频: ${video.title} (aid: ${video.aid})`);
          const result = await this.showVideo(video.aid);
          
          if (result.code === 0) {
            results.success.push({
              aid: video.aid,
              title: video.title
            });
          } else {
            results.failed.push({
              aid: video.aid,
              title: video.title,
              error: result.message || 'Unknown error'
            });
          }
          
          // 添加延迟，避免请求过快
          await new Promise(resolve => setTimeout(resolve, 500));
        } catch (error) {
          results.failed.push({
            aid: video.aid,
            title: video.title,
            error: error.message
          });
        }
      }
    }

    return results;
  }

  // 获取所有唯一的标签
  async getAllTags(videos) {
    const tagSet = new Set();
    
    for (const video of videos) {
      if (video.tags && Array.isArray(video.tags)) {
        video.tags.forEach(tag => {
          if (tag && tag.trim()) {
            tagSet.add(tag.trim());
          }
        });
      }
    }
    
    return Array.from(tagSet).sort();
  }

  // 统计每个标签的视频数量
  async getTagStatistics(videos) {
    const tagStats = {};
    
    for (const video of videos) {
      if (video.tags && Array.isArray(video.tags)) {
        video.tags.forEach(tag => {
          if (tag && tag.trim()) {
            const trimmedTag = tag.trim();
            if (!tagStats[trimmedTag]) {
              tagStats[trimmedTag] = {
                count: 0,
                hiddenCount: 0,
                publicCount: 0
              };
            }
            tagStats[trimmedTag].count++;
            if (video.isHidden) {
              tagStats[trimmedTag].hiddenCount++;
            } else {
              tagStats[trimmedTag].publicCount++;
            }
          }
        });
      }
    }
    
    return tagStats;
  }
}

export default new VideoManager();