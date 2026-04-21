<template>
  <div class="hls-debug">
    <el-alert type="warning" show-icon :closable="false" class="tip">
      <p>
        临时联调页：输入 m3u8 完整地址播放。主列表可走同源代理，例如
        <code>/api/client/exam/media/hls/{ticket}.m3u8</code>
        （需要有效 ticket）。playlist 内的分片或 <code>encryption.key</code>
        如果是对象存储绝对地址，浏览器会直接跨域请求。
      </p>
      <p class="tip-p">
        <strong>本地绕过 CORS：</strong>
        复制 <code>manage/.env.development.example</code> 为
        <code>.env.development.local</code>，设置
        <code>VITE_HLS_STORAGE_ORIGIN=https://你的存储展示域名</code>，
        然后重启 <code>npm run dev</code>。开发模式下，本页会把这个域名下的请求改走 Vite 反代。
      </p>
      <p class="tip-p">
        <strong>正式环境：</strong>
        在 OSS / S3 / MinIO 或前置 CDN 上配置 CORS，允许管理端来源，并暴露
        <code>Content-Type</code>、<code>Content-Length</code> 等头；GET/HEAD 也要包含预检所需头。
        如果仍然出现 <code>403</code>，优先检查预签名是否过期、Host 是否与签名绑定的一致，以及存储策略 / IAM 权限。
      </p>
    </el-alert>

    <el-form label-width="100px" class="form">
      <el-form-item label="m3u8 地址">
        <el-input
          v-model="url"
          type="textarea"
          :rows="2"
          placeholder="https://... 或 /api/client/exam/media/hls/xxx.m3u8"
          clearable
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="loading" @click="play">播放</el-button>
        <el-button @click="stop">停止</el-button>
      </el-form-item>
    </el-form>

    <div class="player-wrap">
      <video ref="videoRef" class="video" controls playsinline />
    </div>

    <el-alert v-if="errorMsg" type="error" :title="errorMsg" show-icon class="err" />
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import Hls, { ErrorTypes, Events, type HlsConfig } from 'hls.js'

const HLS_PROXY_PREFIX = '/__hls_storage_proxy'

const url = ref('')
const videoRef = ref<HTMLVideoElement | null>(null)
const loading = ref(false)
const errorMsg = ref('')

let hls: Hls | null = null

const hlsStorageOrigin = (import.meta.env.VITE_HLS_STORAGE_ORIGIN || '').trim().replace(/\/$/, '')
const devHlsProxyEnabled = Boolean(import.meta.env.DEV && hlsStorageOrigin)

function destroyPlayer() {
  if (hls) {
    hls.destroy()
    hls = null
  }
  const video = videoRef.value
  if (video) {
    video.removeAttribute('src')
    video.load()
  }
  errorMsg.value = ''
}

function play() {
  errorMsg.value = ''
  const src = url.value.trim()
  if (!src) {
    errorMsg.value = '请填写 m3u8 地址'
    return
  }

  const video = videoRef.value
  if (!video) return

  destroyPlayer()
  loading.value = true

  const resolved =
    src.startsWith('http') || src.startsWith('//')
      ? src
      : `${window.location.origin}${src.startsWith('/') ? '' : '/'}${src}`

  const done = () => {
    loading.value = false
  }

  if (Hls.isSupported()) {
    const hlsOpts: Partial<HlsConfig> = {
      enableWorker: true,
      lowLatencyMode: false,
    }

    if (devHlsProxyEnabled && hlsStorageOrigin) {
      hlsOpts.xhrSetup = (xhr, requestUrl) => {
        if (requestUrl.startsWith(hlsStorageOrigin)) {
          try {
            const parsed = new URL(requestUrl)
            xhr.open('GET', `${HLS_PROXY_PREFIX}${parsed.pathname}${parsed.search}`, true)
          } catch {
            xhr.open('GET', requestUrl, true)
          }
        }
      }
    }

    hls = new Hls(hlsOpts)
    hls.on(Events.ERROR, (_event, data) => {
      if (!data.fatal) return

      errorMsg.value = data.type === ErrorTypes.NETWORK_ERROR ? '网络错误' : '播放错误'
      switch (data.type) {
        case ErrorTypes.NETWORK_ERROR:
          hls?.startLoad()
          break
        case ErrorTypes.MEDIA_ERROR:
          hls?.recoverMediaError()
          break
        default:
          destroyPlayer()
          break
      }
    })

    hls.loadSource(resolved)
    hls.attachMedia(video)
    video.play().catch(() => {}).finally(done)
    return
  }

  if (video.canPlayType('application/vnd.apple.mpegurl')) {
    video.src = resolved
    video.addEventListener('loadedmetadata', done, { once: true })
    video.play().catch(() => {}).finally(done)
    return
  }

  errorMsg.value = '当前浏览器不支持 HLS，请使用 Chrome / Edge 或 Safari'
  done()
}

function stop() {
  destroyPlayer()
  loading.value = false
}

onBeforeUnmount(() => {
  destroyPlayer()
})
</script>

<style scoped>
.hls-debug {
  max-width: 960px;
}

.tip {
  margin-bottom: 16px;
}

.tip code {
  font-size: 12px;
}

.tip-p {
  margin: 8px 0 0;
  line-height: 1.5;
}

.tip p:first-child {
  margin: 0;
}

.form {
  margin-bottom: 16px;
}

.player-wrap {
  background: #0f172a;
  border-radius: 8px;
  overflow: hidden;
}

.video {
  display: block;
  width: 100%;
  max-height: 480px;
}

.err {
  margin-top: 12px;
}
</style>
