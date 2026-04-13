<template>
  <div class="hls-debug">
    <el-alert type="warning" show-icon :closable="false" class="tip">
      <p>
        临时联调页：输入 m3u8 完整地址播放。主列表可走同源代理，例如
        <code>/api/client/exam/media/hls/&#123;ticket&#125;.m3u8</code>
        （需有效 ticket）。playlist 内的分片、<code>encryption.key</code> 若为对象存储<strong>绝对地址</strong>，浏览器会跨域拉取。
      </p>
      <p class="tip-p">
        <strong>本地绕过 CORS：</strong>复制 <code>manage/.env.development.example</code> 为
        <code>.env.development.local</code>，设置
        <code>VITE_HLS_STORAGE_ORIGIN=https://你的存储展示域名</code>（与预签名 URL 的协议+主机一致），重启
        <code>npm run dev</code>。本页在开发模式下会把该源下的请求改走 Vite 反代。
      </p>
      <p class="tip-p">
        <strong>正式环境：</strong>在 OSS / S3 / MinIO 桶（或前置 CDN）配置 CORS，允许管理端来源，并暴露
        <code>Content-Type</code>、<code>Content-Length</code> 等；GET/HEAD 需包含预检所需头。
        文件中心存储里可配置 <code>presign_signature_version</code>（<code>v2</code> / <code>v3</code> / <code>v4</code>，其中 v3 与 v4 同为
        SigV4）与 <code>public_base_url</code>：服务端在生成预签名后会把链接的协议与主机替换为该域名（路径与 Query 不变）。
        若仍为 <code>403</code>，检查预签名是否过期、公网域名是否与签名所用 Host 一致（SigV4 与 Host 绑定；AWS 勿随意换域）、桶策略/IAM；可在网络面板查看响应体 XML 中
        <code>Code</code>：<code>SignatureDoesNotMatch</code> 多为 Host/签名不一致，<code>AccessDenied</code> 多为策略或权限。
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
import { ref, onBeforeUnmount } from 'vue'
import Hls from 'hls.js'

const HLS_PROXY_PREFIX = '/__hls_storage_proxy'

const url = ref('')
const videoRef = ref<HTMLVideoElement | null>(null)
const loading = ref(false)
const errorMsg = ref('')

let hls: Hls | null = null

/** 开发环境且配置了 VITE_HLS_STORAGE_ORIGIN 时，将对象存储直链改为走 Vite 反代（同域，无 CORS） */
const hlsStorageOrigin = (import.meta.env.VITE_HLS_STORAGE_ORIGIN || '').trim().replace(/\/$/, '')
const devHlsProxyEnabled = Boolean(import.meta.env.DEV && hlsStorageOrigin)

function destroyPlayer() {
  if (hls) {
    hls.destroy()
    hls = null
  }
  const v = videoRef.value
  if (v) {
    v.removeAttribute('src')
    v.load()
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

  const resolved = src.startsWith('http') || src.startsWith('//') ? src : `${window.location.origin}${src.startsWith('/') ? '' : '/'}${src}`

  const done = () => {
    loading.value = false
  }

  if (Hls.isSupported()) {
    const hlsOpts: Partial<Hls.Options> = {
      enableWorker: true,
      lowLatencyMode: false,
    }
    if (devHlsProxyEnabled && hlsStorageOrigin) {
      hlsOpts.xhrSetup = (xhr, requestUrl) => {
        if (requestUrl.startsWith(hlsStorageOrigin)) {
          try {
            const u = new URL(requestUrl)
            xhr.open('GET', `${HLS_PROXY_PREFIX}${u.pathname}${u.search}`, true)
          } catch {
            xhr.open('GET', requestUrl, true)
          }
        }
      }
    }
    hls = new Hls(hlsOpts)
    hls.on(Hls.Events.ERROR, (_e, data) => {
      if (data.fatal) {
        errorMsg.value = data.type === Hls.ErrorTypes.NETWORK_ERROR ? '网络错误' : '播放错误'
        switch (data.type) {
          case Hls.ErrorTypes.NETWORK_ERROR:
            hls?.startLoad()
            break
          case Hls.ErrorTypes.MEDIA_ERROR:
            hls?.recoverMediaError()
            break
          default:
            destroyPlayer()
            break
        }
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

  errorMsg.value = '当前浏览器不支持 HLS（请用 Chrome/Edge + hls.js 或 Safari）'
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
