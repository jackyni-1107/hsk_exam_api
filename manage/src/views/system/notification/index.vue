<template>
  <div class="page">
    <el-card shadow="never">
      <template #header><span>通知管理</span></template>
      <el-tabs v-model="tab">
        <el-tab-pane label="模板" name="tpl">
          <el-form :inline="true" class="filter">
            <el-form-item label="编码">
              <el-input v-model="tplQuery.code" clearable style="width: 160px" />
            </el-form-item>
            <el-form-item label="渠道">
              <el-input v-model="tplQuery.channel" clearable style="width: 120px" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadTpl">查询</el-button>
              <el-button type="primary" plain @click="openTplCreate">新增模板</el-button>
            </el-form-item>
          </el-form>
          <el-table v-loading="tplLoading" :data="tplRows" border stripe>
            <el-table-column prop="code" label="编码" width="140" />
            <el-table-column prop="name" label="名称" min-width="120" />
            <el-table-column prop="channel" label="渠道" width="100" />
            <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 0 ? 'success' : 'info'" size="small">{{ row.status === 0 ? '正常' : '停用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="create_time" label="创建时间" width="170" :formatter="formatUtcForDisplay" />
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openTplEdit(row)">编辑</el-button>
                <el-button link type="danger" @click="delTpl(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="tplQuery.page"
            v-model:page-size="tplQuery.size"
            class="pager"
            :total="tplTotal"
            layout="total, prev, pager, next"
            @current-change="loadTpl"
          />
        </el-tab-pane>

        <el-tab-pane label="渠道配置" name="ch">
          <div class="toolbar">
            <el-button type="primary" @click="openChCreate">新增渠道</el-button>
            <el-button @click="loadCh">刷新</el-button>
          </div>
          <el-table v-loading="chLoading" :data="chRows" border stripe>
            <el-table-column prop="name" label="名称" min-width="120" />
            <el-table-column prop="channel" label="渠道" width="100" />
            <el-table-column prop="provider" label="提供商" width="120" />
            <el-table-column label="当前" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.is_active === 1" type="success" size="small">是</el-tag>
                <span v-else>否</span>
              </template>
            </el-table-column>
            <el-table-column prop="create_time" label="创建时间" width="170" :formatter="formatUtcForDisplay" />
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <el-button link type="warning" @click="setChActive(row)">设为当前</el-button>
                <el-button link type="primary" @click="openChEdit(row)">编辑</el-button>
                <el-button link type="danger" @click="delCh(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="发送记录" name="log">
          <el-form :inline="true" class="filter">
            <el-form-item label="渠道">
              <el-input v-model="logQuery.channel" clearable style="width: 120px" />
            </el-form-item>
            <el-form-item label="收件人">
              <el-input v-model="logQuery.recipient" clearable style="width: 160px" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadLog">查询</el-button>
            </el-form-item>
          </el-form>
          <el-table v-loading="logLoading" :data="logRows" border stripe>
            <el-table-column prop="template_code" label="模板" width="140" />
            <el-table-column prop="channel" label="渠道" width="90" />
            <el-table-column prop="recipient" label="收件人" min-width="160" />
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 0 ? 'success' : 'danger'" size="small">{{ row.status === 0 ? '成功' : '失败' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="error_msg" label="错误" min-width="160" show-overflow-tooltip />
            <el-table-column prop="create_time" label="时间" width="170" :formatter="formatUtcForDisplay" />
          </el-table>
          <el-pagination
            v-model:current-page="logQuery.page"
            v-model:page-size="logQuery.size"
            class="pager"
            :total="logTotal"
            layout="total, prev, pager, next"
            @current-change="loadLog"
          />
        </el-tab-pane>

        <el-tab-pane label="测试发送" name="send">
          <el-form :model="sendForm" label-width="100px" style="max-width: 520px">
            <el-form-item label="模板编码" required>
              <el-input v-model="sendForm.template_code" placeholder="已存在的 template code" />
            </el-form-item>
            <el-form-item label="渠道" required>
              <el-input v-model="sendForm.channel" placeholder="如 email、sms" />
            </el-form-item>
            <el-form-item label="收件人" required>
              <el-input v-model="sendForm.recipient" />
            </el-form-item>
            <el-form-item label="变量 JSON">
              <el-input v-model="sendForm.variables" type="textarea" :rows="4" placeholder='可选，如 {"name":"张三"}' />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="sendLoading" @click="doSend">发送</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="tplDlg" :title="tplMode === 'create' ? '新增模板' : '编辑模板'" width="560px" destroy-on-close>
      <el-form ref="tplFormRef" :model="tplForm" :rules="tplRules" label-width="88px">
        <el-form-item label="编码" prop="code">
          <el-input v-model="tplForm.code" :disabled="tplMode === 'edit'" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="tplForm.name" />
        </el-form-item>
        <el-form-item label="渠道" prop="channel">
          <el-input v-model="tplForm.channel" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="tplForm.content" type="textarea" :rows="5" />
        </el-form-item>
        <el-form-item label="变量说明">
          <el-input v-model="tplForm.variables" placeholder="可选" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="tplForm.status">
            <el-radio :label="0">正常</el-radio>
            <el-radio :label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tplDlg = false">取消</el-button>
        <el-button type="primary" :loading="tplSaving" @click="saveTpl">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="chDlg" :title="chMode === 'create' ? '新增渠道' : '编辑渠道'" width="520px" destroy-on-close>
      <el-form ref="chFormRef" :model="chForm" :rules="chRules" label-width="88px">
        <el-form-item label="渠道" prop="channel">
          <el-input v-model="chForm.channel" :disabled="chMode === 'edit'" />
        </el-form-item>
        <el-form-item label="提供商" prop="provider">
          <el-input v-model="chForm.provider" :disabled="chMode === 'edit'" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="chForm.name" />
        </el-form-item>
        <el-form-item label="配置 JSON" prop="config_json">
          <el-input v-model="chForm.config_json" type="textarea" :rows="6" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="chDlg = false">取消</el-button>
        <el-button type="primary" :loading="chSaving" @click="saveCh">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchTemplates,
  createTemplate,
  updateTemplate,
  deleteTemplate,
  fetchChannelConfigs,
  createChannelConfig,
  updateChannelConfig,
  deleteChannelConfig,
  setActiveChannelConfig,
  fetchNotificationLogs,
  sendNotification,
  type TemplateItem,
  type ChannelConfigItem,
} from '@/api/notification'
import { formatUtcForDisplay } from '@/utils/datetime'

const tab = ref('tpl')
const tplLoading = ref(false)
const chLoading = ref(false)
const logLoading = ref(false)
const tplRows = ref<TemplateItem[]>([])
const tplTotal = ref(0)
const chRows = ref<ChannelConfigItem[]>([])
const logRows = ref<unknown[]>([])
const logTotal = ref(0)

const tplQuery = reactive({ page: 1, size: 10, code: '', channel: '' })
const logQuery = reactive({ page: 1, size: 10, channel: '', recipient: '' })

const tplDlg = ref(false)
const tplMode = ref<'create' | 'edit'>('create')
const tplSaving = ref(false)
const tplFormRef = ref<FormInstance>()
const tplEditId = ref(0)
const tplForm = reactive({
  code: '',
  name: '',
  channel: '',
  content: '',
  variables: '',
  status: 0,
})
const tplRules: FormRules = {
  code: [{ required: true, message: '必填', trigger: 'blur' }],
  name: [{ required: true, message: '必填', trigger: 'blur' }],
  channel: [{ required: true, message: '必填', trigger: 'blur' }],
  content: [{ required: true, message: '必填', trigger: 'blur' }],
}

const chDlg = ref(false)
const chMode = ref<'create' | 'edit'>('create')
const chSaving = ref(false)
const chFormRef = ref<FormInstance>()
const chEditId = ref(0)
const chForm = reactive({ channel: '', provider: '', name: '', config_json: '{}' })
const chRules: FormRules = {
  channel: [{ required: true, message: '必填', trigger: 'blur' }],
  provider: [{ required: true, message: '必填', trigger: 'blur' }],
  name: [{ required: true, message: '必填', trigger: 'blur' }],
  config_json: [{ required: true, message: '必填', trigger: 'blur' }],
}

const sendForm = reactive({ template_code: '', channel: '', recipient: '', variables: '' })
const sendLoading = ref(false)

async function loadTpl() {
  tplLoading.value = true
  try {
    const res = (await fetchTemplates({
      page: tplQuery.page,
      size: tplQuery.size,
      code: tplQuery.code || undefined,
      channel: tplQuery.channel || undefined,
    })) as { data?: { list?: TemplateItem[]; total?: number } }
    tplRows.value = res?.data?.list ?? []
    tplTotal.value = res?.data?.total ?? 0
  } finally {
    tplLoading.value = false
  }
}

async function loadCh() {
  chLoading.value = true
  try {
    const res = (await fetchChannelConfigs()) as { data?: { list?: ChannelConfigItem[] } }
    chRows.value = res?.data?.list ?? []
  } finally {
    chLoading.value = false
  }
}

async function loadLog() {
  logLoading.value = true
  try {
    const res = (await fetchNotificationLogs({
      page: logQuery.page,
      size: logQuery.size,
      channel: logQuery.channel || undefined,
      recipient: logQuery.recipient || undefined,
    })) as { data?: { list?: unknown[]; total?: number } }
    logRows.value = res?.data?.list ?? []
    logTotal.value = res?.data?.total ?? 0
  } finally {
    logLoading.value = false
  }
}

function openTplCreate() {
  tplMode.value = 'create'
  tplEditId.value = 0
  Object.assign(tplForm, { code: '', name: '', channel: '', content: '', variables: '', status: 0 })
  tplDlg.value = true
}

function openTplEdit(row: TemplateItem) {
  tplMode.value = 'edit'
  tplEditId.value = row.id
  Object.assign(tplForm, {
    code: row.code,
    name: row.name,
    channel: row.channel,
    content: row.content,
    variables: row.variables,
    status: row.status,
  })
  tplDlg.value = true
}

async function saveTpl() {
  await tplFormRef.value?.validate(async (ok) => {
    if (!ok) return
    tplSaving.value = true
    try {
      if (tplMode.value === 'create') {
        await createTemplate({ ...tplForm })
        ElMessage.success('已创建')
      } else {
        await updateTemplate(tplEditId.value, {
          name: tplForm.name,
          content: tplForm.content,
          variables: tplForm.variables,
          status: tplForm.status,
        })
        ElMessage.success('已保存')
      }
      tplDlg.value = false
      loadTpl()
    } catch {
      /* */
    } finally {
      tplSaving.value = false
    }
  })
}

function delTpl(row: TemplateItem) {
  ElMessageBox.confirm(`删除模板「${row.code}」？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteTemplate(row.id)
      ElMessage.success('已删除')
      loadTpl()
    })
    .catch(() => {})
}

function openChCreate() {
  chMode.value = 'create'
  chEditId.value = 0
  Object.assign(chForm, { channel: '', provider: '', name: '', config_json: '{}' })
  chDlg.value = true
}

function openChEdit(row: ChannelConfigItem) {
  chMode.value = 'edit'
  chEditId.value = row.id
  Object.assign(chForm, {
    channel: row.channel,
    provider: row.provider,
    name: row.name,
    config_json: row.config_json || '{}',
  })
  chDlg.value = true
}

async function saveCh() {
  await chFormRef.value?.validate(async (ok) => {
    if (!ok) return
    chSaving.value = true
    try {
      if (chMode.value === 'create') {
        await createChannelConfig({ ...chForm })
        ElMessage.success('已创建')
      } else {
        await updateChannelConfig(chEditId.value, { name: chForm.name, config_json: chForm.config_json })
        ElMessage.success('已保存')
      }
      chDlg.value = false
      loadCh()
    } catch {
      /* */
    } finally {
      chSaving.value = false
    }
  })
}

function delCh(row: ChannelConfigItem) {
  ElMessageBox.confirm(`删除渠道「${row.name}」？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteChannelConfig(row.id)
      ElMessage.success('已删除')
      loadCh()
    })
    .catch(() => {})
}

function setChActive(row: ChannelConfigItem) {
  ElMessageBox.confirm(`将「${row.name}」设为当前渠道？`, '确认', { type: 'warning' })
    .then(async () => {
      await setActiveChannelConfig(row.id)
      ElMessage.success('已设置')
      loadCh()
    })
    .catch(() => {})
}

async function doSend() {
  if (!sendForm.template_code || !sendForm.channel || !sendForm.recipient) {
    ElMessage.warning('请填写模板、渠道、收件人')
    return
  }
  sendLoading.value = true
  try {
    const res = (await sendNotification({
      template_code: sendForm.template_code,
      channel: sendForm.channel,
      recipient: sendForm.recipient,
      variables: sendForm.variables || undefined,
    })) as { data?: { ok?: boolean } }
    if (res?.data?.ok) ElMessage.success('已提交发送')
    else ElMessage.success('请求已完成')
  } catch {
    /* */
  } finally {
    sendLoading.value = false
  }
}

watch(tab, (v) => {
  if (v === 'tpl') loadTpl()
  if (v === 'ch') loadCh()
  if (v === 'log') loadLog()
})

onMounted(() => {
  loadTpl()
})
</script>

<style scoped>
.page {
  padding: 8px 0;
}
.filter {
  margin-bottom: 12px;
}
.toolbar {
  margin-bottom: 12px;
}
.pager {
  margin-top: 12px;
  justify-content: flex-end;
  display: flex;
}
</style>
