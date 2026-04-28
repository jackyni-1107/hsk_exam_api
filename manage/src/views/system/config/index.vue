<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>配置中心</span>
          <el-button v-permission="'config:create'" type="primary" @click="openCreate">新增配置</el-button>
        </div>
      </template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="分组">
          <el-select v-model="query.group" clearable filterable style="width: 140px">
            <el-option v-for="item in groupOptions" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="键">
          <el-input v-model="query.key" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="ID" width="72" />
        <el-table-column prop="config_key" label="键" min-width="160" show-overflow-tooltip />
        <el-table-column label="值" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">{{ configValuePreview(row) }}</template>
        </el-table-column>
        <el-table-column label="类型" width="100">
          <template #default="{ row }">{{ configTypeLabel(row.config_type) }}</template>
        </el-table-column>
        <el-table-column prop="group_name" label="分组" width="120" />
        <el-table-column prop="remark" label="备注" min-width="140" show-overflow-tooltip />
        <el-table-column prop="create_time" label="创建时间" width="170" :formatter="formatUtcForDisplay" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button v-permission="'config:update'" link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button v-permission="'config:delete'" link type="danger" @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pager">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.size"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          background
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <el-dialog v-model="dlg" :title="mode === 'create' ? '新增配置' : '编辑配置'" width="520px" destroy-on-close @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="96px">
        <el-form-item label="键" prop="config_key">
          <el-input v-model="form.config_key" :disabled="mode === 'edit'" />
        </el-form-item>
        <el-form-item label="值">
          <el-input v-if="normalizedConfigType === 'string'" v-model="form.config_value" type="textarea" :rows="4" />
          <el-input-number
            v-else-if="normalizedConfigType === 'number'"
            v-model="configNumberValue"
            style="width: 100%"
            controls-position="right"
          />
          <el-switch
            v-else-if="normalizedConfigType === 'boolean'"
            v-model="configBooleanValue"
            active-text="开启"
            inactive-text="关闭"
          />
          <div v-else class="json-config-list">
            <div v-for="(item, index) in jsonConfigRows" :key="index" class="json-config-row">
              <el-input v-model="item.key" placeholder="字段名" />
              <el-input v-model="item.value" placeholder="字段值" />
              <el-button link type="danger" @click="removeJsonConfigRow(index)">删除</el-button>
            </div>
            <el-button link type="primary" @click="addJsonConfigRow">添加字段</el-button>
          </div>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.config_type" :disabled="mode === 'edit'" style="width: 100%" @change="onConfigTypeChange">
            <el-option label="字符串" value="string" />
            <el-option label="数字" value="number" />
            <el-option label="布尔" value="boolean" />
            <el-option label="对象" value="json" />
          </el-select>
        </el-form-item>
        <el-form-item label="分组">
          <el-select v-model="form.group_name" filterable style="width: 100%" placeholder="请先在字典 sys_config_group 中配置">
            <el-option v-for="item in groupOptions" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlg = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchConfigList,
  createConfig,
  updateConfig,
  deleteConfig,
  type ConfigItemRow,
} from '@/api/sysConfig'
import { fetchDictDataList, type DictDataItem } from '@/api/sysDict'
import { formatUtcForDisplay } from '@/utils/datetime'

const loading = ref(false)
const saving = ref(false)
const groupLoading = ref(false)
const rows = ref<ConfigItemRow[]>([])
const groupRows = ref<DictDataItem[]>([])
const total = ref(0)
const dlg = ref(false)
const mode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const editId = ref(0)

const query = reactive({ page: 1, size: 10, group: '', key: '' })
const form = reactive({
  config_key: '',
  config_value: '',
  config_type: 'string',
  group_name: 'default',
  remark: '',
})
const configNumberValue = ref(0)
const configBooleanValue = ref(false)
const jsonConfigRows = ref<{ key: string; value: string }[]>([])
const configTypeOptions = new Set(['string', 'number', 'boolean', 'json'])
const normalizedConfigType = computed(() => {
  const type = (form.config_type || 'string').toLowerCase()
  return configTypeOptions.has(type) ? type : 'string'
})
const groupOptions = computed(() => {
  const groups = new Set<string>()
  for (const row of groupRows.value) {
    if (row.status === 0 && row.dict_value) groups.add(row.dict_value)
  }
  if (query.group) groups.add(query.group)
  if (form.group_name) groups.add(form.group_name)
  return Array.from(groups)
})

const rules: FormRules = {
  config_key: [{ required: true, message: '必填', trigger: 'blur' }],
}

function configTypeLabel(type?: string) {
  const map: Record<string, string> = {
    string: '字符串',
    number: '数字',
    boolean: '布尔',
    json: '对象',
  }
  return map[(type || 'string').toLowerCase()] ?? type ?? '字符串'
}

function stringifyConfigValue(value: unknown) {
  if (value === undefined || value === null) return ''
  if (typeof value === 'string') return value
  if (typeof value === 'number' || typeof value === 'boolean') return String(value)
  return JSON.stringify(value)
}

function parseConfigJsonValue(value: string) {
  const trimmed = value.trim()
  if (!trimmed) return ''
  if (!['{', '[', '"'].includes(trimmed[0]) && !['true', 'false', 'null'].includes(trimmed) && Number.isNaN(Number(trimmed))) {
    return value
  }
  try {
    return JSON.parse(trimmed) as unknown
  } catch {
    return value
  }
}

function parseJsonConfigRows(value?: string) {
  const trimmed = value?.trim() || '{}'
  try {
    const parsed = JSON.parse(trimmed) as unknown
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
      jsonConfigRows.value = Object.entries(parsed as Record<string, unknown>).map(([key, val]) => ({
        key,
        value: stringifyConfigValue(val),
      }))
      return
    }
  } catch {
    /* Historical data may not be a JSON object. */
  }
  jsonConfigRows.value = trimmed && trimmed !== '{}' ? [{ key: 'value', value: trimmed }] : []
}

function fillTypedValue(value: string, type: string) {
  const normalizedType = configTypeOptions.has(type.toLowerCase()) ? type.toLowerCase() : 'string'
  form.config_value = value
  if (normalizedType === 'number') {
    const num = Number(value)
    configNumberValue.value = Number.isFinite(num) ? num : 0
  } else if (normalizedType === 'boolean') {
    configBooleanValue.value = value === 'true' || value === '1'
  } else if (normalizedType === 'json') {
    parseJsonConfigRows(value)
  } else {
    jsonConfigRows.value = []
  }
}

function buildConfigValue() {
  if (normalizedConfigType.value === 'number') {
    return String(configNumberValue.value)
  }
  if (normalizedConfigType.value === 'boolean') {
    return configBooleanValue.value ? 'true' : 'false'
  }
  if (normalizedConfigType.value === 'json') {
    const data: Record<string, unknown> = {}
    for (const row of jsonConfigRows.value) {
      const key = row.key.trim()
      if (!key) continue
      data[key] = parseConfigJsonValue(row.value)
    }
    return JSON.stringify(data)
  }
  return form.config_value
}

function configValuePreview(row: ConfigItemRow) {
  const type = (row.config_type || 'string').toLowerCase()
  if (type === 'boolean') return row.config_value === 'true' || row.config_value === '1' ? '开启' : '关闭'
  if (type === 'json') {
    try {
      const parsed = JSON.parse(row.config_value || '{}') as unknown
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
        return Object.entries(parsed as Record<string, unknown>)
          .map(([key, val]) => `${key}: ${stringifyConfigValue(val)}`)
          .join('；') || '{}'
      }
    } catch {
      /* Fallback to stored value. */
    }
  }
  return row.config_value
}

function onConfigTypeChange() {
  fillTypedValue('', normalizedConfigType.value)
}

function addJsonConfigRow() {
  jsonConfigRows.value.push({ key: '', value: '' })
}

function removeJsonConfigRow(index: number) {
  jsonConfigRows.value.splice(index, 1)
}

async function loadConfigGroups() {
  groupLoading.value = true
  try {
    const res = (await fetchDictDataList({
      page: 1,
      size: 200,
      dict_type: 'sys_config_group',
    })) as { data?: { list?: DictDataItem[] } }
    groupRows.value = res?.data?.list ?? []
  } finally {
    groupLoading.value = false
  }
}

async function loadList() {
  loading.value = true
  try {
    const res = (await fetchConfigList({
      page: query.page,
      size: query.size,
      group: query.group || undefined,
      key: query.key || undefined,
    })) as { data?: { list?: ConfigItemRow[]; total?: number } }
    rows.value = res?.data?.list ?? []
    total.value = res?.data?.total ?? 0
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.page = 1
  query.size = 10
  query.group = ''
  query.key = ''
  loadList()
}

function resetForm() {
  editId.value = 0
  form.config_key = ''
  form.config_value = ''
  form.config_type = 'string'
  form.group_name = 'default'
  form.remark = ''
  configNumberValue.value = 0
  configBooleanValue.value = false
  jsonConfigRows.value = []
  formRef.value?.clearValidate()
}

function openCreate() {
  mode.value = 'create'
  resetForm()
  if (!groupOptions.value.length) {
    ElMessage.warning('请先在字典 sys_config_group 中配置分组')
  }
  dlg.value = true
}

function openEdit(row: ConfigItemRow) {
  mode.value = 'edit'
  resetForm()
  editId.value = row.id
  form.config_key = row.config_key
  form.config_value = row.config_value
  form.config_type = row.config_type || 'string'
  form.group_name = row.group_name || 'default'
  form.remark = row.remark
  fillTypedValue(row.config_value, form.config_type)
  dlg.value = true
}

async function submit() {
  await formRef.value?.validate(async (ok) => {
    if (!ok) return
    const configValue = buildConfigValue()
    if (normalizedConfigType.value === 'string' && !configValue.trim()) {
      ElMessage.warning('请填写配置值')
      return
    }
    saving.value = true
    try {
      if (mode.value === 'create') {
        await createConfig({
          config_key: form.config_key,
          config_value: configValue,
          config_type: form.config_type || undefined,
          group_name: form.group_name || undefined,
          remark: form.remark || undefined,
        })
        ElMessage.success('已创建')
      } else {
        await updateConfig(editId.value, {
          config_value: configValue,
          remark: form.remark,
        })
        ElMessage.success('已保存')
      }
      dlg.value = false
      await loadList()
    } catch {
      /* */
    } finally {
      saving.value = false
    }
  })
}

function onDelete(row: ConfigItemRow) {
  ElMessageBox.confirm(`删除配置「${row.config_key}」？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteConfig(row.id)
      ElMessage.success('已删除')
      await loadList()
    })
    .catch(() => {})
}

onMounted(async () => {
  await Promise.all([loadConfigGroups(), loadList()])
})
</script>

<style scoped>
.page {
  padding: 8px 0;
}
.head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.filter {
  margin-bottom: 12px;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
.json-config-list {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.json-config-row {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 8px;
  align-items: center;
}
@media (max-width: 720px) {
  .json-config-row {
    grid-template-columns: 1fr;
  }
}
</style>
