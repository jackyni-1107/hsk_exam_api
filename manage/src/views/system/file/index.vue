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
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
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

    <el-dialog v-model="stDlg" :title="stMode === 'create' ? '新增存储' : '编辑存储'" width="520px" destroy-on-close>
      <el-form ref="stFormRef" :model="stForm" :rules="stRules" label-width="120px">
        <el-form-item label="存储类型" prop="storage_type">
          <el-input v-model="stForm.storage_type" :disabled="stMode === 'edit'" placeholder="如 local、s3" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="stForm.name" />
        </el-form-item>
        <el-form-item label="配置 JSON" prop="config_json">
          <el-input v-model="stForm.config_json" type="textarea" :rows="8" />
        </el-form-item>
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
import { reactive, ref, watch, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchFileList,
  deleteFile,
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
const stForm = reactive({
  storage_type: '',
  name: '',
  config_json: '{}',
  cleanup_before_days: 0,
})
const stRules: FormRules = {
  storage_type: [{ required: true, message: '必填', trigger: 'blur' }],
  name: [{ required: true, message: '必填', trigger: 'blur' }],
  config_json: [{ required: true, message: '必填', trigger: 'blur' }],
}

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

function delFile(row: FileItemRow) {
  ElMessageBox.confirm(`删除文件「${row.filename}」？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteFile(row.id)
      ElMessage.success('已删除')
      loadFiles()
    })
    .catch(() => {})
}

function openStCreate() {
  stMode.value = 'create'
  stEditId.value = 0
  Object.assign(stForm, { storage_type: '', name: '', config_json: '{}', cleanup_before_days: 0 })
  stDlg.value = true
}

function openStEdit(row: StorageConfigItem) {
  stMode.value = 'edit'
  stEditId.value = row.id
  Object.assign(stForm, {
    storage_type: row.storage_type,
    name: row.name,
    config_json: row.config_json || '{}',
    cleanup_before_days: row.cleanup_before_days ?? 0,
  })
  stDlg.value = true
}

async function saveSt() {
  await stFormRef.value?.validate(async (ok) => {
    if (!ok) return
    stSaving.value = true
    try {
      if (stMode.value === 'create') {
        await createStorageConfig({ ...stForm })
        ElMessage.success('已创建')
      } else {
        await updateStorageConfig(stEditId.value, {
          name: stForm.name,
          config_json: stForm.config_json,
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
</style>
