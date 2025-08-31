import fetch from 'node-fetch';
import fs from 'fs/promises';
import path from 'path';
import { fileURLToPath } from 'url';
import { CONFIG } from './config.js';
import cookieManager from './cookie-manager.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

class BilibiliVideoFetcher {
  constructor() {
    this.allVideos = [];
    this.statistics = {
      totalVideos: 0,
      publishedVideos: 0,
      unpublishedVideos: 0,
      publishingVideos: 0,
      totalViews: 0,
      totalLikes: 0,
      totalCoins: 0,
      totalFavorites: 0,
      totalDanmaku: 0,
      totalComments: 0,
      totalShares: 0,
      videosByCategory: {},
      fetchTime: new Date().toISOString()
    };
  }

  async fetchVideos(pageNum = 1, pageSize = CONFIG.API.DEFAULT_PAGE_SIZE, keyword = '') {
    let url = `${CONFIG.API.BASE_URL}${CONFIG.API.ARCHIVES_ENDPOINT}?status=is_pubing%2Cpubed%2Cnot_pubed&pn=${pageNum}&ps=${pageSize}&coop=1&interactive=1`;
    if (keyword) {
      url += `&keyword=${encodeURIComponent(keyword)}`;
    }
    
    const cookieData = await cookieManager.getCookie();
    if (!cookieData.cookie) {
      throw new Error('Cookie not configured. Please set cookie first.');
    }
    
    try {
      const response = await fetch(url, {
        headers: {
          ...CONFIG.HEADERS,
          'Cookie': cookieData.cookie
        }
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      
      if (data.code !== 0) {
        throw new Error(`API error: ${data.message}`);
      }

      return data.data;
    } catch (error) {
      console.error(`获取第${pageNum}页数据失败:`, error);
      throw error;
    }
  }

  async fetchAllVideos() {
    console.log('开始获取所有视频数据...');
    
    // 先获取第一页，了解总数
    const firstPageData = await this.fetchVideos(1, CONFIG.API.DEFAULT_PAGE_SIZE);
    
    // 更新统计信息
    this.statistics.publishedVideos = firstPageData.class.pubed || 0;
    this.statistics.unpublishedVideos = firstPageData.class.not_pubed || 0;
    this.statistics.publishingVideos = firstPageData.class.is_pubing || 0;
    this.statistics.totalVideos = this.statistics.publishedVideos + 
                                 this.statistics.unpublishedVideos + 
                                 this.statistics.publishingVideos;
    
    // 处理第一页的视频
    if (firstPageData.arc_audits && firstPageData.arc_audits.length > 0) {
      this.processVideos(firstPageData.arc_audits);
    }
    
    // 计算总页数
    const totalPages = Math.ceil(this.statistics.totalVideos / CONFIG.API.DEFAULT_PAGE_SIZE);
    
    console.log(`总视频数: ${this.statistics.totalVideos}, 总页数: ${totalPages}`);
    
    // 获取剩余页面
    for (let page = 2; page <= totalPages; page++) {
      console.log(`正在获取第${page}/${totalPages}页...`);
      try {
        const pageData = await this.fetchVideos(page, CONFIG.API.DEFAULT_PAGE_SIZE);
        if (pageData.arc_audits && pageData.arc_audits.length > 0) {
          this.processVideos(pageData.arc_audits);
        }
        // 添加延迟，避免请求过快
        await new Promise(resolve => setTimeout(resolve, 500));
      } catch (error) {
        console.error(`获取第${page}页失败，继续下一页`);
      }
    }
    
    console.log(`成功获取${this.allVideos.length}个视频数据`);
    return this.allVideos;
  }

  processVideos(videos) {
    for (const video of videos) {
      const processedVideo = {
        // 基本信息
        aid: video.Archive.aid,
        bvid: video.Archive.bvid,
        title: video.Archive.title,
        cover: video.Archive.cover,
        desc: video.Archive.desc,
        duration: video.Archive.duration,
        
        // 发布信息
        state: video.Archive.state,
        stateDesc: video.Archive.state_desc,
        pubTime: video.Archive.ptime ? new Date(video.Archive.ptime * 1000).toISOString() : null,
        createTime: video.Archive.ctime ? new Date(video.Archive.ctime * 1000).toISOString() : null,
        
        // 分类信息
        tid: video.Archive.tid,
        typename: video.typename,
        parentTname: video.parent_tname,
        
        // 统计数据
        stats: {
          view: video.stat.view || 0,
          danmaku: video.stat.danmaku || 0,
          reply: video.stat.reply || 0,
          favorite: video.stat.favorite || 0,
          coin: video.stat.coin || 0,
          share: video.stat.share || 0,
          like: video.stat.like || 0
        },
        
        // 其他信息
        tags: video.Archive.tag ? video.Archive.tag.split(',') : [],
        copyright: video.Archive.copyright,
        isCharging: video.Archive.charging_pay === 1,  // 充电视频标识
        isOnlySelf: video.Archive.is_only_self === 1,   // 仅自己可见标识
        url: `https://www.bilibili.com/video/${video.Archive.bvid}`
      };
      
      this.allVideos.push(processedVideo);
      
      // 更新统计数据
      if (video.Archive.state === 0) { // 只统计已发布的视频
        this.statistics.totalViews += processedVideo.stats.view;
        this.statistics.totalLikes += processedVideo.stats.like;
        this.statistics.totalCoins += processedVideo.stats.coin;
        this.statistics.totalFavorites += processedVideo.stats.favorite;
        this.statistics.totalDanmaku += processedVideo.stats.danmaku;
        this.statistics.totalComments += processedVideo.stats.reply;
        this.statistics.totalShares += processedVideo.stats.share;
        
        // 按分类统计
        const category = processedVideo.parentTname || '未分类';
        if (!this.statistics.videosByCategory[category]) {
          this.statistics.videosByCategory[category] = 0;
        }
        this.statistics.videosByCategory[category]++;
      }
    }
  }

  async saveData() {
    const dataDir = path.join(__dirname, '..', 'data');
    
    // 确保数据目录存在
    await fs.mkdir(dataDir, { recursive: true });
    
    // 保存视频数据
    const videosPath = path.join(dataDir, 'videos.json');
    await fs.writeFile(videosPath, JSON.stringify(this.allVideos, null, 2));
    console.log(`视频数据已保存到: ${videosPath}`);
    
    // 保存统计数据
    const statsPath = path.join(dataDir, 'statistics.json');
    await fs.writeFile(statsPath, JSON.stringify(this.statistics, null, 2));
    console.log(`统计数据已保存到: ${statsPath}`);
  }
}

// 主函数
async function main() {
  const fetcher = new BilibiliVideoFetcher();
  
  try {
    await fetcher.fetchAllVideos();
    await fetcher.saveData();
    
    console.log('\n=== 统计信息 ===');
    console.log(`总视频数: ${fetcher.statistics.totalVideos}`);
    console.log(`已发布: ${fetcher.statistics.publishedVideos}`);
    console.log(`未发布: ${fetcher.statistics.unpublishedVideos}`);
    console.log(`发布中: ${fetcher.statistics.publishingVideos}`);
    console.log(`\n总播放量: ${fetcher.statistics.totalViews.toLocaleString()}`);
    console.log(`总点赞数: ${fetcher.statistics.totalLikes.toLocaleString()}`);
    console.log(`总硬币数: ${fetcher.statistics.totalCoins.toLocaleString()}`);
    console.log(`总收藏数: ${fetcher.statistics.totalFavorites.toLocaleString()}`);
    console.log(`\n分类统计:`);
    for (const [category, count] of Object.entries(fetcher.statistics.videosByCategory)) {
      console.log(`  ${category}: ${count}个视频`);
    }
  } catch (error) {
    console.error('获取视频数据失败:', error);
    process.exit(1);
  }
}

// 如果直接运行此文件
if (import.meta.url === `file://${process.argv[1]}`) {
  main();
}

export { BilibiliVideoFetcher };