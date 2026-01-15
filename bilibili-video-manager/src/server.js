import express from 'express';
import cors from 'cors';
import fs from 'fs/promises';
import path from 'path';
import { fileURLToPath } from 'url';
import { BilibiliVideoFetcher } from './fetch-videos.js';
import cookieManager from './cookie-manager.js';
import videoManager from './video-manager.js';
import videoBatchManager from './video-batch-manager.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const app = express();
const PORT = 3000;

// 中间件
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(express.static(path.join(__dirname, '..', 'public')));

// API路由 - 获取视频列表
app.get('/api/videos', async (req, res) => {
    try {
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        res.json(JSON.parse(data));
    } catch (error) {
        console.error('读取视频数据失败:', error);
        res.status(500).json({ error: '无法读取视频数据' });
    }
});

// API路由 - 获取统计数据
app.get('/api/statistics', async (req, res) => {
    try {
        const statsPath = path.join(__dirname, '..', 'data', 'statistics.json');
        const data = await fs.readFile(statsPath, 'utf-8');
        res.json(JSON.parse(data));
    } catch (error) {
        console.error('读取统计数据失败:', error);
        res.status(500).json({ error: '无法读取统计数据' });
    }
});

// API路由 - 刷新数据
app.post('/api/refresh', async (req, res) => {
    try {
        console.log('开始刷新数据...');
        const fetcher = new BilibiliVideoFetcher();
        await fetcher.fetchAllVideos();
        await fetcher.saveData();
        
        res.json({ 
            success: true, 
            message: '数据刷新成功',
            totalVideos: fetcher.statistics.totalVideos
        });
    } catch (error) {
        console.error('刷新数据失败:', error);
        res.status(500).json({ 
            success: false, 
            error: '数据刷新失败', 
            details: error.message 
        });
    }
});

// API路由 - 获取视频详情
app.get('/api/video/:bvid', async (req, res) => {
    try {
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        const videos = JSON.parse(data);
        const video = videos.find(v => v.bvid === req.params.bvid);
        
        if (video) {
            res.json(video);
        } else {
            res.status(404).json({ error: '视频不存在' });
        }
    } catch (error) {
        console.error('获取视频详情失败:', error);
        res.status(500).json({ error: '无法获取视频详情' });
    }
});

// API路由 - 获取分类统计
app.get('/api/statistics/categories', async (req, res) => {
    try {
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        const videos = JSON.parse(data);
        
        // 按分类统计
        const categoryStats = {};
        videos.forEach(video => {
            if (video.state === 0) { // 只统计已发布的视频
                const category = video.parentTname || '未分类';
                if (!categoryStats[category]) {
                    categoryStats[category] = {
                        count: 0,
                        totalViews: 0,
                        totalLikes: 0,
                        totalCoins: 0,
                        totalFavorites: 0
                    };
                }
                categoryStats[category].count++;
                categoryStats[category].totalViews += video.stats.view || 0;
                categoryStats[category].totalLikes += video.stats.like || 0;
                categoryStats[category].totalCoins += video.stats.coin || 0;
                categoryStats[category].totalFavorites += video.stats.favorite || 0;
            }
        });
        
        res.json(categoryStats);
    } catch (error) {
        console.error('获取分类统计失败:', error);
        res.status(500).json({ error: '无法获取分类统计' });
    }
});

// API路由 - 设置Cookie
app.post('/api/cookie', async (req, res) => {
    try {
        const { cookie } = req.body;
        if (!cookie) {
            return res.status(400).json({ error: 'Cookie is required' });
        }
        
        // 验证cookie
        const validation = await cookieManager.validateCookie(cookie);
        if (!validation.valid) {
            return res.status(400).json({ 
                error: `Invalid cookie: missing ${validation.missing}` 
            });
        }
        
        // 保存cookie
        const result = await cookieManager.saveCookie(cookie);
        res.json({ 
            success: true, 
            message: 'Cookie saved successfully',
            csrf: result.csrf 
        });
    } catch (error) {
        console.error('保存Cookie失败:', error);
        res.status(500).json({ error: '保存Cookie失败' });
    }
});

// API路由 - 获取Cookie状态
app.get('/api/cookie/status', async (req, res) => {
    try {
        const cookieData = await cookieManager.getCookie();
        res.json({
            hasCookie: !!cookieData.cookie,
            updatedAt: cookieData.updatedAt
        });
    } catch (error) {
        res.status(500).json({ error: '获取Cookie状态失败' });
    }
});

// API路由 - 获取所有标签
app.get('/api/tags', async (req, res) => {
    try {
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        const videos = JSON.parse(data);
        
        const tags = await videoManager.getAllTags(videos);
        const stats = await videoManager.getTagStatistics(videos);
        
        res.json({ tags, stats });
    } catch (error) {
        console.error('获取标签失败:', error);
        res.status(500).json({ error: '无法获取标签' });
    }
});

// API路由 - 批量隐藏视频（按标签）
app.post('/api/videos/hide-by-tag', async (req, res) => {
    try {
        const { tag } = req.body;
        if (!tag) {
            return res.status(400).json({ error: 'Tag is required' });
        }
        
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        const videos = JSON.parse(data);
        
        const results = await videoManager.hideVideosByTag(tag, videos);
        
        res.json({
            success: true,
            message: `处理完成`,
            results
        });
    } catch (error) {
        console.error('批量隐藏视频失败:', error);
        res.status(500).json({ error: '批量隐藏视频失败: ' + error.message });
    }
});

// API路由 - 批量显示视频（按标签）
app.post('/api/videos/show-by-tag', async (req, res) => {
    try {
        const { tag } = req.body;
        if (!tag) {
            return res.status(400).json({ error: 'Tag is required' });
        }
        
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        const videos = JSON.parse(data);
        
        const results = await videoManager.showVideosByTag(tag, videos);
        
        res.json({
            success: true,
            message: `处理完成`,
            results
        });
    } catch (error) {
        console.error('批量显示视频失败:', error);
        res.status(500).json({ error: '批量显示视频失败: ' + error.message });
    }
});

// API路由 - 隐藏单个视频
app.post('/api/video/:aid/hide', async (req, res) => {
    try {
        const { aid } = req.params;
        const result = await videoManager.hideVideo(aid);
        
        if (result.code === 0) {
            res.json({ success: true, message: '视频已隐藏' });
        } else {
            res.status(400).json({ error: result.message || '隐藏失败' });
        }
    } catch (error) {
        console.error('隐藏视频失败:', error);
        res.status(500).json({ error: '隐藏视频失败: ' + error.message });
    }
});

// API路由 - 显示单个视频
app.post('/api/video/:aid/show', async (req, res) => {
    try {
        const { aid } = req.params;
        const result = await videoManager.showVideo(aid);
        
        if (result.code === 0) {
            res.json({ success: true, message: '视频已公开' });
        } else {
            res.status(400).json({ error: result.message || '公开失败' });
        }
    } catch (error) {
        console.error('公开视频失败:', error);
        res.status(500).json({ error: '公开视频失败: ' + error.message });
    }
});

// API路由 - 批量隐藏所有普通视频
app.post('/api/videos/batch-hide-normal', async (req, res) => {
    try {
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        const videos = JSON.parse(data);
        
        const results = await videoBatchManager.hideAllNormalVideos(videos);
        
        res.json({
            success: true,
            message: '批量隐藏完成',
            results
        });
    } catch (error) {
        console.error('批量隐藏失败:', error);
        res.status(500).json({ error: '批量隐藏失败: ' + error.message });
    }
});

// API路由 - 批量删除所有充电视频
app.post('/api/videos/batch-delete-charging', async (req, res) => {
    try {
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        const videos = JSON.parse(data);
        
        const results = await videoBatchManager.deleteAllChargingVideos(videos);
        
        res.json({
            success: true,
            message: '批量删除完成',
            results
        });
    } catch (error) {
        console.error('批量删除失败:', error);
        res.status(500).json({ error: '批量删除失败: ' + error.message });
    }
});

// API路由 - 获取视频统计
app.get('/api/videos/stats', async (req, res) => {
    try {
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        const videos = JSON.parse(data);
        
        const stats = await videoBatchManager.getVideoStats(videos);
        
        res.json(stats);
    } catch (error) {
        console.error('获取统计失败:', error);
        res.status(500).json({ error: '获取统计失败: ' + error.message });
    }
});

// API路由 - 删除单个视频
app.delete('/api/video/:aid', async (req, res) => {
    try {
        const { aid } = req.params;
        const result = await videoBatchManager.deleteVideo(aid);
        
        if (result.code === 0) {
            res.json({ success: true, message: '视频已删除' });
        } else {
            res.status(400).json({ error: result.message || '删除失败' });
        }
    } catch (error) {
        console.error('删除视频失败:', error);
        res.status(500).json({ error: '删除视频失败: ' + error.message });
    }
});

// API路由 - 获取时间趋势数据
app.get('/api/statistics/timeline', async (req, res) => {
    try {
        const videosPath = path.join(__dirname, '..', 'data', 'videos.json');
        const data = await fs.readFile(videosPath, 'utf-8');
        const videos = JSON.parse(data);
        
        // 按月统计
        const monthlyStats = {};
        videos.forEach(video => {
            if (video.pubTime && video.state === 0) {
                const date = new Date(video.pubTime);
                const monthKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
                
                if (!monthlyStats[monthKey]) {
                    monthlyStats[monthKey] = {
                        count: 0,
                        totalViews: 0,
                        totalLikes: 0
                    };
                }
                monthlyStats[monthKey].count++;
                monthlyStats[monthKey].totalViews += video.stats.view || 0;
                monthlyStats[monthKey].totalLikes += video.stats.like || 0;
            }
        });
        
        // 转换为数组并排序
        const timeline = Object.entries(monthlyStats)
            .map(([month, stats]) => ({ month, ...stats }))
            .sort((a, b) => a.month.localeCompare(b.month));
        
        res.json(timeline);
    } catch (error) {
        console.error('获取时间趋势失败:', error);
        res.status(500).json({ error: '无法获取时间趋势' });
    }
});

// 启动服务器
app.listen(PORT, async () => {
    console.log(`服务器运行在 http://localhost:${PORT}`);
    
    // 检查是否存在数据文件
    const dataDir = path.join(__dirname, '..', 'data');
    const videosPath = path.join(dataDir, 'videos.json');
    const statsPath = path.join(dataDir, 'statistics.json');
    
    try {
        await fs.access(videosPath);
        await fs.access(statsPath);
        console.log('数据文件已存在');
    } catch {
        console.log('数据文件不存在，正在初始化...');
        console.log('请先运行 npm run fetch 获取数据');
        
        // 创建空数据文件
        await fs.mkdir(dataDir, { recursive: true });
        await fs.writeFile(videosPath, '[]');
        await fs.writeFile(statsPath, JSON.stringify({
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
        }, null, 2));
    }
});