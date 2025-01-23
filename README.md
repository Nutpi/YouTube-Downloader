# YouTube 视频下载器

一个基于 Go 语言开发的 YouTube 视频下载工具，支持多种分辨率下载和音视频合并。

## 功能特点

- 支持多种视频分辨率选择
- 自动合并音视频
- 实时显示下载进度
- 简洁的 Web 界面
- 下载历史记录

## 系统要求

- Go 1.21 或更高版本
- yt-dlp
- FFmpeg

## 安装依赖

### 1. 安装 Go
访问 https://go.dev/dl/ 下载并安装 Go

### 2. 安装 yt-dlp
- MacOS: `brew install yt-dlp`
- Linux: `sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && sudo chmod a+rx /usr/local/bin/yt-dlp`
- Windows: `winget install yt-dlp`

### 3. 安装 FFmpeg
- MacOS: `brew install ffmpeg`
- Linux: `sudo apt install ffmpeg`
- Windows: 访问 https://ffmpeg.org/download.html 下载并安装

## 运行项目

1. 克隆项目：
```bash
git clone [你的仓库地址]
cd youtube-downloader
```

## 使用方法

1. 在浏览器中访问 `http://localhost:8080`
2. 输入 YouTube 视频链接
3. 选择所需的视频质量
4. 点击下载按钮

## 开发相关

如果你想贡献代码，请确保：

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 致谢

- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)