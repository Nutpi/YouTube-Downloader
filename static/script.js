async function getVideoInfo() {
    const videoUrl = document.getElementById('videoUrl').value;
    if (!videoUrl) {
        alert('请输入视频链接');
        return;
    }

    try {
        const response = await fetch('/video-info', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `url=${encodeURIComponent(videoUrl)}`
        });

        const data = await response.json();
        if (data.error) {
            alert(data.error);
            return;
        }

        document.getElementById('videoTitle').textContent = data.title;
        const qualitySelect = document.getElementById('qualitySelect');
        qualitySelect.innerHTML = '<option value="">选择视频质量</option>';
        
        data.formats.forEach(format => {
            const option = document.createElement('option');
            option.value = format.quality;
            option.textContent = format.quality;
            qualitySelect.appendChild(option);
        });

        document.getElementById('videoInfo').style.display = 'block';
        document.getElementById('progressContainer').style.display = 'none';
    } catch (error) {
        alert('获取视频信息失败：' + error.message);
    }
}

async function downloadVideo() {
    const videoUrl = document.getElementById('videoUrl').value;
    const quality = document.getElementById('qualitySelect').value;
    
    if (!quality) {
        alert('请选择视频质量');
        return;
    }

    document.getElementById('progressContainer').style.display = 'block';
    document.getElementById('progressText').textContent = '准备下载...';
    document.getElementById('progress').style.width = '0%';

    const eventSource = new EventSource(`/download?url=${encodeURIComponent(videoUrl)}&quality=${encodeURIComponent(quality)}`);

    eventSource.addEventListener('progress', function(e) {
        const progressText = e.data;
        document.getElementById('progressText').textContent = progressText;
        
        // 尝试从进度文本中提取百分比
        const match = progressText.match(/(\d+\.?\d*)%/);
        if (match) {
            const percent = parseFloat(match[1]);
            document.getElementById('progress').style.width = `${percent}%`;
        }
    });

    eventSource.addEventListener('complete', function(e) {
        eventSource.close();
        document.getElementById('progressText').textContent = '下载完成';
        document.getElementById('progress').style.width = '100%';
        addToHistory(document.getElementById('videoTitle').textContent, quality);
    });

    eventSource.addEventListener('error', function(e) {
        eventSource.close();
        document.getElementById('progressText').textContent = '下载失败：' + e.data;
        document.getElementById('progress').style.width = '0%';
    });
}

function addToHistory(title, quality) {
    const historyList = document.getElementById('historyList');
    const li = document.createElement('li');
    const now = new Date().toLocaleString();
    li.innerHTML = `
        <span>${title} (${quality})</span>
        <span>${now}</span>
    `;
    historyList.insertBefore(li, historyList.firstChild);
}