import fetch from 'node-fetch';
import { CONFIG } from './config.js';
import cookieManager from './cookie-manager.js';
import videoManager from './video-manager.js';

class VideoBatchManager {
  
  // 删除单个视频
  async deleteVideo(aid) {
    const cookieData = await cookieManager.getCookie();
    if (!cookieData.cookie || !cookieData.csrf) {
      throw new Error('Cookie or CSRF token not configured');
    }

    // B站删除视频的API
    const url = `${CONFIG.API.BASE_URL}/x/web/archive/delete?csrf=${cookieData.csrf}`;
    
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
      body: `aid=${aid}&csrf=${cookieData.csrf}`
    });

    const result = await response.json();
    return result;
  }

  // 批量隐藏所有普通视频（非充电视频）
  async hideAllNormalVideos(videos) {
    const results = {
      success: [],
      failed: [],
      skipped: [],
      totalProcessed: 0
    };

    for (const video of videos) {
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
          reason: '已经是仅自己可见'
        });
        continue;
      }

      // 跳过未发布的视频
      if (video.state !== 0) {
        results.skipped.push({
          aid: video.aid,
          title: video.title,
          reason: '视频未发布'
        });
        continue;
      }

      try {
        console.log(`正在隐藏视频: ${video.title} (aid: ${video.aid})`);
        const result = await videoManager.hideVideo(video.aid);
        
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
        
        results.totalProcessed++;
        
        // 添加延迟，避免请求过快
        await new Promise(resolve => setTimeout(resolve, 1000));
      } catch (error) {
        results.failed.push({
          aid: video.aid,
          title: video.title,
          error: error.message
        });
      }
    }

    return results;
  }

  // 批量删除所有充电视频
  async deleteAllChargingVideos(videos) {
    const results = {
      success: [],
      failed: [],
      skipped: [],
      totalProcessed: 0
    };

    // 筛选出充电视频
    const chargingVideos = videos.filter(v => v.isCharging);
    
    if (chargingVideos.length === 0) {
      return {
        ...results,
        message: '没有找到充电视频'
      };
    }

    // 再次确认
    console.log(`找到 ${chargingVideos.length} 个充电视频准备删除`);

    for (const video of chargingVideos) {
      try {
        console.log(`正在删除充电视频: ${video.title} (aid: ${video.aid})`);
        const result = await this.deleteVideo(video.aid);
        
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
        
        results.totalProcessed++;
        
        // 删除操作需要更长的延迟
        await new Promise(resolve => setTimeout(resolve, 2000));
      } catch (error) {
        results.failed.push({
          aid: video.aid,
          title: video.title,
          error: error.message
        });
      }
    }

    return results;
  }

  // 获取视频统计信息
  async getVideoStats(videos) {
    const stats = {
      total: videos.length,
      normal: 0,
      charging: 0,
      hidden: 0,
      published: 0,
      unpublished: 0
    };

    for (const video of videos) {
      if (video.isCharging) {
        stats.charging++;
      } else {
        stats.normal++;
      }

      if (video.isOnlySelf) {
        stats.hidden++;
      }

      if (video.state === 0) {
        stats.published++;
      } else {
        stats.unpublished++;
      }
    }

    return stats;
  }
}

export default new VideoBatchManager();