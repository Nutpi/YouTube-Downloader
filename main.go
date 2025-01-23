package main

import (
    "bufio"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "os/exec"
    "encoding/json"
    "os"
    "sort"
    "strconv"
)

type VideoInfo struct {
    ID       string   `json:"id"`
    Title    string   `json:"title"`
    Formats  []Format `json:"formats"`
}

type Format struct {
    Quality string `json:"quality"`
    URL     string `json:"url"`
}

type DownloadProgress struct {
    Total     int64   `json:"total"`
    Current   int64   `json:"current"`
    Progress  float64 `json:"progress"`
}

func main() {
    r := gin.Default()
    
    // 提供静态文件服务
    r.Static("/static", "./static")
    r.LoadHTMLGlob("templates/*")

    // 首页
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })

    // 获取视频信息
    r.POST("/video-info", getVideoInfo)
    
    // 下载视频
    r.POST("/download", downloadVideo)
    
    // 修改下载路由为 GET 方法
    r.GET("/download", downloadVideo)

    r.Run(":8080")
}

func getVideoInfo(c *gin.Context) {
    url := c.PostForm("url")
    
    // 使用 yt-dlp 获取视频信息
    cmd := exec.Command("yt-dlp", "-J", url)
    output, err := cmd.Output()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 解析 JSON 输出
    var rawInfo map[string]interface{}
    if err := json.Unmarshal(output, &rawInfo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 提取格式信息
    formats := make([]Format, 0)
    seenQualities := make(map[int]bool)  // 用于去重

    if formatsRaw, ok := rawInfo["formats"].([]interface{}); ok {
        for _, f := range formatsRaw {
            format := f.(map[string]interface{})
            
            height, hasHeight := format["height"].(float64)
            vcodec, hasVcodec := format["vcodec"].(string)
            
            // 只处理视频格式，并且去重
            if hasHeight && hasVcodec && vcodec != "none" {
                h := int(height)
                if !seenQualities[h] {
                    seenQualities[h] = true
                    formats = append(formats, Format{
                        Quality: fmt.Sprintf("%dp", h),
                        URL:    fmt.Sprintf("%d", h),  // 只保存分辨率数字
                    })
                }
            }
        }
    }

    // 按分辨率从高到低排序
    sort.Slice(formats, func(i, j int) bool {
        qi, _ := strconv.Atoi(formats[i].URL)
        qj, _ := strconv.Atoi(formats[j].URL)
        return qi > qj
    })

    videoInfo := VideoInfo{
        ID:      rawInfo["id"].(string),
        Title:   rawInfo["title"].(string),
        Formats: formats,
    }

    c.JSON(http.StatusOK, videoInfo)
}

func downloadVideo(c *gin.Context) {
    url := c.Query("url")
    quality := c.Query("quality")
    
    // 创建下载目录
    if err := os.MkdirAll("downloads", 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建下载目录"})
        return
    }

    // 修改下载命令参数
    args := []string{
        "-f", fmt.Sprintf("bestvideo[height=%s]+bestaudio", quality[:4]),
        "-o", "downloads/%(title)s.%(ext)s",
        "--merge-output-format", "mp4",
        "--prefer-ffmpeg",              // 使用 FFmpeg 进行合并
        "--ffmpeg-location", "/usr/local/bin/ffmpeg",  // 指定 FFmpeg 路径
        "--no-playlist",
        url,
    }
    
    // 创建命令并设置错误输出
    cmd := exec.Command("yt-dlp", args...)
    cmd.Stderr = os.Stderr
    
    // 获取命令的输出管道
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 启动命令
    if err := cmd.Start(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 设置 SSE 头部
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")

    // 创建扫描器读取输出
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        // 发送进度信息到客户端
        c.SSEvent("progress", scanner.Text())
        c.Writer.Flush()
    }

    // 等待命令完成
    if err := cmd.Wait(); err != nil {
        c.SSEvent("error", err.Error())
        return
    }

    c.SSEvent("complete", "下载完成")
}