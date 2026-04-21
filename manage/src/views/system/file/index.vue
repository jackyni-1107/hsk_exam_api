<template>
  <div class="page">
    <el-card shadow="never">
      <template #header><span>文件中心</span></template>
      <el-tabs v-model="tab">
        <el-tab-pane label="文件列表" name="files">
          <el-form :inline="true" class="filter">
            <el-form-item label="文件名">
              <el-input v-model="fileQuery.filename" clearable style="width: 200px" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadFiles">查询</el-button>
            </el-form-item>
            <el-form-item>
              <el-upload :show-file-list="false" accept="*/*" :http-request="handleFileUploadRequest">
                <el-button type="success">上传</el-button>
              </el-upload>
            </el-form-item>
          </el-form>
          <el-table v-loading="fileLoading" :data="fileRows" border stripe>
            <el-table-column prop="id" label="ID" width="72" />
            <el-table-column prop="filename" label="文件名" min-width="160" show-overflow-tooltip />
            <el-table-column prop="path" label="路径" min-width="200" show-overflow-tooltip />
            <el-table-column label="大小" width="100">
              <template #default="{ row }">{{ formatSize(row.size) }}</template>
            </el-table-column>
            <el-table-column prop="mime_type" label="类型" width="120" show-overflow-tooltip />
            <el-table-column label="私有" width="72">
              <template #default="{ row }">{{ row.is_private === 1 ? '是' : '否' }}</template>
            </el-table-column>
            <el-table-column prop="create_time" label="上传时间" width="170" :formatter="formatUtcForDisplay" />
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="downloadRow(row)">下载</el-button>
                <el-button link type="danger" @click="delFile(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="fileQuery.page"
            v-model:page-size="fileQuery.size"
            class="pager"
            :total="fileTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @size-change="loadFiles"
            @current-change="loadFiles"
          />
        </el-tab-pane>

        <el-tab-pane label="存储配置" name="storage">
          <div class="toolbar">
            <el-button type="primary" @click="openStCreate">新增配置</el-button>
            <el-button @click="loadStorage">刷新</el-button>
          </div>
          <el-table v-loading="stLoading" :data="stRows" border stripe>
            <el-table-column prop="name" label="名称" min-width="140" />
            <el-table-column prop="storage_type" label="类型" width="120" />
            <el-table-column label="当前" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.is_active === 1" type="success" size="small">是</el-tag>
                <span v-else>否</span>
              </template>
            </el-table-column>
            <el-table-column prop="cleanup_before_days" label="清理(天)" width="100" />
            <el-table-column prop="create_time" label="创建时间" width="170" :formatter="formatUtcForDisplay" />
            <el-table-column label="操作" width="240" fixed="right">
              <template #default="{ row }">
                <el-button link type="warning" @click="setStActive(row)">设为当前</el-button>
                <el-button link type="primary" @click="openStEdit(row)">编辑</el-button>
                <el-button link type="danger" @click="delSt(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="stDlg" :title="stMode === 'create' ? '新增存储' : '编辑存储'" width="680px" destroy-on-close>
      <el-form ref="stFormRef" :model="stForm" :rules="stRules" label-width="132px">
        <el-form-item label="存储类型" prop="storage_type">
          <el-select
            v-model="stForm.storage_type"
            :disabled="stMode === 'edit'"
            style="width: 100%"
            placeholder="请选择"
            filterable
            allow-create
            default-first-option
          >
            <el-option label="local（本地磁盘）" value="local" />
            <el-option label="s3" value="s3" />
            <el-option label="oss（S3 兼容）" value="oss" />
            <el-option label="minio" value="minio" />
          </el-select>
          <div class="field-hint">编辑时不可改类型；可输入列表外的自定义类型（与后端约定一致即可）。</div>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="stForm.name" />
        </el-form-item>

        <template v-if="isLocalStorage">
          <el-divider content-position="left">本地配置</el-divider>
          <el-form-item label="存储根路径" prop="base_path">
            <el-input v-model="stForm.base_path" placeholder="默认 ./storage" />
          </el-form-item>
          <el-form-item label="访问域名" prop="public_base_url">
            <el-input
              v-model="stForm.public_base_url"
              clearable
              placeholder="可选，填写后对外链接使用该前缀（如 https://files.example.com/static）"
            />
            <div class="field-hint">本地存储时用于拼接公开访问 URL，与对象存储的「访问域名」语义类似。</div>
          </el-form-item>
        </template>

        <template v-else-if="showObjectStorageForm">
          <el-divider content-position="left">对象存储</el-divider>
          <el-form-item label="Endpoint" prop="endpoint">
            <el-input v-model="stForm.endpoint" placeholder="如 https://minio.example.com:9000" />
          </el-form-item>
          <el-form-item label="Bucket" prop="bucket">
            <el-input v-model="stForm.bucket" placeholder="桶名称" />
          </el-form-item>
          <el-form-item label="Access Key" prop="access_key">
            <el-input v-model="stForm.access_key" autocomplete="off" />
          </el-form-item>
          <el-form-item label="Secret Key" prop="secret_key">
            <el-input v-model="stForm.secret_key" type="password" show-password autocomplete="new-password" />
          </el-form-item>
          <el-form-item label="Region" prop="region">
            <el-input v-model="stForm.region" placeholder="如 us-east-1，留空则后端默认" />
          </el-form-item>
          <el-form-item label="Path 风格 URL" prop="s3_force_path_style">
            <el-switch v-model="stForm.s3_force_path_style" active-text="开启（路径含 /bucket/）" inactive-text="关闭（推荐，virtual-hosted）" />
            <div class="field-hint">
              关闭时预签名地址为 <code>https://&lt;bucket&gt;.&lt;endpoint主机&gt;/对象路径</code>，路径段不再重复桶名；旧版 MinIO 或 IP Endpoint 若异常可改为开启。
            </div>
          </el-form-item>
          <el-form-item label="预签名版本" prop="presign_signature_version">
            <el-select v-model="stForm.presign_signature_version" style="width: 100%" placeholder="S3 兼容预签名">
              <el-option label="v4（SigV4，默认）" value="v4" />
              <el-option label="v3（与 v4 相同）" value="v3" />
              <el-option label="v2（SigV2，仅 GET）" value="v2" />
            </el-select>
            <div class="field-hint">
              v3 与 v4 同为 SigV4；v2 适用于部分旧版兼容端点。
            </div>
          </el-form-item>
          <el-form-item label="访问域名" prop="public_base_url">
            <el-input
              v-model="stForm.public_base_url"
              clearable
              placeholder="https://cdn.example.com（仅协议+主机+端口，可选）"
            />
            <div class="field-hint">
              预签名生成后替换链接中的协议与主机，不改变路径与 Query。SigV4 与 Host 绑定，随意换域可能导致 403。
            </div>
          </el-form-item>
        </template>

        <template v-else>
          <el-alert type="info" show-icon :closable="false" class="type-hint">
            请先选择或填写存储类型；保存时将生成服务端所需配置（界面不展示原始 JSON）。
          </el-alert>
        </template>

        <el-divider />
        <el-form-item label="清理早于(天)">
          <el-input-number v-model="stForm.cleanup_before_days" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="stDlg = false">取消</el-button>
        <el-button type="primary" :loading="stSaving" @click="saveSt">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch, onMounted, computed } from 'vue'
import type { FormInstance, FormRules, UploadRequestOptions } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchFileList,
  deleteFile,
  uploadSysFile,
  downloadSysFile,
  fetchStorageConfigs,
  createStorageConfig,
  updateStorageConfig,
  deleteStorageConfig,
  setActiveStorageConfig,
  type FileItemRow,
  type StorageConfigItem,
} from '@/api/file'
import { formatUtcForDisplay } from '@/utils/datetime'

const tab = ref('files')
const fileLoading = ref(false)
const stLoading = ref(false)
const fileRows = ref<FileItemRow[]>([])
const fileTotal = ref(0)
const stRows = ref<StorageConfigItem[]>([])
const fileQuery = reactive({ page: 1, size: 10, filename: '' })

const stDlg = ref(false)
const stMode = ref<'create' | 'edit'>('create')
const stSaving = ref(false)
const stFormRef = ref<FormInstance>()
const stEditId = ref(0)
type PresignSig = 'v2' | 'v3' | 'v4'

const stForm = reactive({
  storage_type: '',
  name: '',
  base_path: './storage',
  endpoint: '',
  bucket: '',
  access_key: '',
  secret_key: '',
  region: '',
  presign_signature_version: 'v4' as PresignSig,
  /** true = path-style（/bucket/key）；false = virtual-hosted（路径仅对象 key） */
  s3_force_path_style: false,
  public_base_url: '',
  cleanup_before_days: 0,
})

const isLocalStorage = computed(() => stForm.storage_type.trim().toLowerCase() === 'local')
/** 非 local 且已选类型时展示对象存储表单项（含 s3/oss/minio 及自定义 S3 兼容类型） */
const showObjectStorageForm = computed(() => {
  const t = stForm.storage_type.trim()
  return t.length > 0 && !isLocalStorage.value
})

const stRules = computed<FormRules>(() => {
  const base: FormRules = {
    storage_type: [{ required: true, message: '必填', trigger: 'change' }],
    name: [{ required: true, message: '必填', trigger: 'blur' }],
  }
  if (showObjectStorageForm.value) {
    base.endpoint = [{ required: true, message: '必填', trigger: 'blur' }]
    base.bucket = [{ required: true, message: '必填', trigger: 'blur' }]
    base.access_key = [{ required: true, message: '必填', trigger: 'blur' }]
    base.secret_key = [{ required: true, message: '必填', trigger: 'blur' }]
  }
  return base
})

function formatSize(n: number) {
  if (n < 1024) return `${n} B`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
  return `${(n / 1024 / 1024).toFixed(2)} MB`
}

async function loadFiles() {
  fileLoading.value = true
  try {
    const res = (await fetchFileList({
      page: fileQuery.page,
      size: fileQuery.size,
      filename: fileQuery.filename || undefined,
    })) as { data?: { list?: FileItemRow[]; total?: number } }
    fileRows.value = res?.data?.list ?? []
    fileTotal.value = res?.data?.total ?? 0
  } finally {
    fileLoading.value = false
  }
}

async function loadStorage() {
  stLoading.value = true
  try {
    const res = (await fetchStorageConfigs()) as { data?: { list?: StorageConfigItem[] } }
    stRows.value = res?.data?.list ?? []
  } finally {
    stLoading.value = false
  }
}

async function handleFileUploadRequest(opt: UploadRequestOptions & { onError?: (error: Error) => void }) {
  try {
    await uploadSysFile(opt.file as File, 0)
    ElMessage.success('上传成功')
    loadFiles()
    opt.onSuccess?.({} as never)
  } catch (e) {
    opt.onError?.(e instanceof Error ? e : new Error('上传失败'))
  }
}

async function downloadRow(row: FileItemRow) {
  try {
    await downloadSysFile(row.id, row.filename)
  } catch (e: unknown) {
    ElMessage.error(e instanceof Error ? e.message : '下载失败')
  }
}

function delFile(row: FileItemRow) {
  ElMessageBox.confirm(`删除文件「${row.filename}」？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteFile(row.id)
      ElMessage.success('已删除')
      loadFiles()
    })
    .catch(() => {})
}

function strField(obj: Record<string, unknown>, key: string): string {
  const v = obj[key]
  return typeof v === 'string' ? v : ''
}

/** 从后端返回的 config_json 填充表单项（不在界面展示原始 JSON） */
function parseStorageConfigFromBackendJson(jsonStr: string) {
  let obj: Record<string, unknown> = {}
  try {
    const raw = jsonStr?.trim() || '{}'
    const o = JSON.parse(raw) as unknown
    if (typeof o === 'object' && o !== null && !Array.isArray(o)) obj = o as Record<string, unknown>
  } catch {
    obj = {}
  }
  stForm.base_path = strField(obj, 'base_path') || './storage'
  stForm.endpoint = strField(obj, 'endpoint')
  stForm.bucket = strField(obj, 'bucket')
  stForm.access_key = strField(obj, 'access_key')
  stForm.secret_key = strField(obj, 'secret_key')
  stForm.region = strField(obj, 'region')
  stForm.public_base_url = strField(obj, 'public_base_url')
  const v = strField(obj, 'presign_signature_version')
    .trim()
    .toLowerCase()
  if (v === 'v2' || v === '2' || v === 'sigv2' || v === 's3v2') stForm.presign_signature_version = 'v2'
  else if (v === 'v3' || v === '3' || v === 'sigv3' || v === 's3v3') stForm.presign_signature_version = 'v3'
  else stForm.presign_signature_version = 'v4'
  if ('s3_force_path_style' in obj && typeof obj.s3_force_path_style === 'boolean') {
    stForm.s3_force_path_style = obj.s3_force_path_style as boolean
  } else {
    stForm.s3_force_path_style = true
  }
}

/** 根据表单项生成提交给后端的 config_json */
function buildStorageConfigJsonForSave(): string {
  const t = stForm.storage_type.trim().toLowerCase()
  if (t === 'local') {
    const o: Record<string, string> = {
      base_path: stForm.base_path.trim() || './storage',
    }
    const pub = stForm.public_base_url.trim()
    if (pub) o.public_base_url = pub
    return JSON.stringify(o)
  }
  if (showObjectStorageForm.value) {
    const o: Record<string, unknown> = {
      endpoint: stForm.endpoint.trim(),
      bucket: stForm.bucket.trim(),
      access_key: stForm.access_key.trim(),
      secret_key: stForm.secret_key.trim(),
      presign_signature_version: stForm.presign_signature_version,
      s3_force_path_style: stForm.s3_force_path_style,
    }
    const region = stForm.region.trim()
    if (region) o.region = region
    const pub = stForm.public_base_url.trim()
    if (pub) o.public_base_url = pub
    return JSON.stringify(o)
  }
  throw new Error('请选择存储类型')
}

function openStCreate() {
  stMode.value = 'create'
  stEditId.value = 0
  Object.assign(stForm, {
    storage_type: '',
    name: '',
    base_path: './storage',
    endpoint: '',
    bucket: '',
    access_key: '',
    secret_key: '',
    region: '',
    presign_signature_version: 'v4' as PresignSig,
    s3_force_path_style: false,
    public_base_url: '',
    cleanup_before_days: 0,
  })
  stDlg.value = true
}

function openStEdit(row: StorageConfigItem) {
  stMode.value = 'edit'
  stEditId.value = row.id
  Object.assign(stForm, {
    storage_type: row.storage_type,
    name: row.name,
    cleanup_before_days: row.cleanup_before_days ?? 0,
  })
  parseStorageConfigFromBackendJson(row.config_json || '{}')
  stDlg.value = true
}

async function saveSt() {
  await stFormRef.value?.validate(async (ok) => {
    if (!ok) return
    const pub = stForm.public_base_url.trim()
    if (pub) {
      try {
        const u = new URL(pub)
        if (!u.protocol || !u.host) throw new Error('bad')
      } catch {
        ElMessage.error('访问域名须为完整 URL（含协议与主机，例如 https://cdn.example.com）')
        return
      }
    }
    let mergedJson: string
    try {
      mergedJson = buildStorageConfigJsonForSave()
    } catch {
      ElMessage.error('请完善存储类型与对应配置项')
      return
    }
    stSaving.value = true
    try {
      if (stMode.value === 'create') {
        await createStorageConfig({
          storage_type: stForm.storage_type,
          name: stForm.name,
          config_json: mergedJson,
          cleanup_before_days: stForm.cleanup_before_days,
        })
        ElMessage.success('已创建')
      } else {
        await updateStorageConfig(stEditId.value, {
          name: stForm.name,
          config_json: mergedJson,
          cleanup_before_days: stForm.cleanup_before_days,
        })
        ElMessage.success('已保存')
      }
      stDlg.value = false
      loadStorage()
    } catch {
      /* */
    } finally {
      stSaving.value = false
    }
  })
}

function delSt(row: StorageConfigItem) {
  ElMessageBox.confirm(`删除存储「${row.name}」？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteStorageConfig(row.id)
      ElMessage.success('已删除')
      loadStorage()
    })
    .catch(() => {})
}

function setStActive(row: StorageConfigItem) {
  ElMessageBox.confirm(`将「${row.name}」设为当前存储？`, '确认', { type: 'warning' })
    .then(async () => {
      await setActiveStorageConfig(row.id)
      ElMessage.success('已设置')
      loadStorage()
    })
    .catch(() => {})
}

watch(tab, (v) => {
  if (v === 'files') loadFiles()
  if (v === 'storage') loadStorage()
})

onMounted(() => loadFiles())
</script>

<style scoped>
.page {
  padding: 8px 0;
}
.filter,
.toolbar {
  margin-bottom: 12px;
}
.pager {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}
.field-hint {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.45;
}
.type-hint {
  margin-bottom: 8px;
}
</style>
