<template>
  <div class="page">
    <el-card shadow="never" class="card">
      <template #header>
        <div class="head">
          <span>字典类型</span>
          <el-button type="primary" @click="openTypeCreate">新增类型</el-button>
        </div>
      </template>
      <el-form :inline="true" class="filter">
        <el-form-item label="类型编码">
          <el-input v-model="typeQuery.dict_type" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadTypes">查询</el-button>
        </el-form-item>
      </el-form>
      <el-table
        v-loading="typeLoading"
        :data="typeRows"
        border
        stripe
        highlight-current-row
        @current-change="onTypeRow"
      >
        <el-table-column prop="dict_name" label="名称" min-width="140" />
        <el-table-column prop="dict_type" label="类型编码" min-width="160" />
        <el-table-column label="状态" width="88">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'" size="small">{{ row.status === 0 ? '正常' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="170" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openTypeEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="delType(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="typeQuery.page"
        v-model:page-size="typeQuery.size"
        class="pager"
        :total="typeTotal"
        layout="total, prev, pager, next"
        @current-change="loadTypes"
      />
    </el-card>

    <el-card v-if="currentDictType" shadow="never" class="card data-card">
      <template #header>
        <div class="head">
          <span>字典数据 — {{ currentDictType }}</span>
          <el-button type="primary" @click="openDataCreate">新增数据</el-button>
        </div>
      </template>
      <el-table v-loading="dataLoading" :data="dataRows" border stripe>
        <el-table-column prop="dict_label" label="标签" min-width="120" />
        <el-table-column prop="dict_value" label="值" min-width="120" />
        <el-table-column prop="sort" label="排序" width="72" />
        <el-table-column label="状态" width="88">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'" size="small">{{ row.status === 0 ? '正常' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="170" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDataEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="delData(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="dataQuery.page"
        v-model:page-size="dataQuery.size"
        class="pager"
        :total="dataTotal"
        layout="total, prev, pager, next"
        @current-change="loadData"
      />
    </el-card>
    <el-empty v-else description="请在上表选择一条字典类型" class="hint" />

    <el-dialog v-model="typeDlg" :title="typeMode === 'create' ? '新增字典类型' : '编辑类型'" width="480px" destroy-on-close>
      <el-form ref="typeFormRef" :model="typeForm" :rules="typeRules" label-width="96px">
        <el-form-item label="名称" prop="dict_name">
          <el-input v-model="typeForm.dict_name" />
        </el-form-item>
        <el-form-item label="类型编码" prop="dict_type">
          <el-input v-model="typeForm.dict_type" :disabled="typeMode === 'edit'" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="typeForm.status">
            <el-radio :label="0">正常</el-radio>
            <el-radio :label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="typeForm.remark" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="typeDlg = false">取消</el-button>
        <el-button type="primary" :loading="typeSaving" @click="saveType">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="dataDlg" :title="dataMode === 'create' ? '新增字典数据' : '编辑数据'" width="480px" destroy-on-close>
      <el-form ref="dataFormRef" :model="dataForm" :rules="dataRules" label-width="88px">
        <el-form-item label="标签" prop="dict_label">
          <el-input v-model="dataForm.dict_label" />
        </el-form-item>
        <el-form-item label="值" prop="dict_value">
          <el-input v-model="dataForm.dict_value" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="dataForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="dataForm.status">
            <el-radio :label="0">正常</el-radio>
            <el-radio :label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="dataForm.remark" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dataDlg = false">取消</el-button>
        <el-button type="primary" :loading="dataSaving" @click="saveData">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchDictTypeList,
  createDictType,
  updateDictType,
  deleteDictType,
  fetchDictDataList,
  createDictData,
  updateDictData,
  deleteDictData,
  type DictTypeItem,
  type DictDataItem,
} from '@/api/sysConfig'

const typeLoading = ref(false)
const dataLoading = ref(false)
const typeRows = ref<DictTypeItem[]>([])
const typeTotal = ref(0)
const dataRows = ref<DictDataItem[]>([])
const dataTotal = ref(0)
const currentDictType = ref('')

const typeQuery = reactive({ page: 1, size: 10, dict_type: '' })
const dataQuery = reactive({ page: 1, size: 10 })

const typeDlg = ref(false)
const typeMode = ref<'create' | 'edit'>('create')
const typeSaving = ref(false)
const typeFormRef = ref<FormInstance>()
const typeEditId = ref(0)
const typeForm = reactive({ dict_name: '', dict_type: '', status: 0, remark: '' })
const typeRules: FormRules = {
  dict_name: [{ required: true, message: '必填', trigger: 'blur' }],
  dict_type: [{ required: true, message: '必填', trigger: 'blur' }],
}

const dataDlg = ref(false)
const dataMode = ref<'create' | 'edit'>('create')
const dataSaving = ref(false)
const dataFormRef = ref<FormInstance>()
const dataEditId = ref(0)
const dataForm = reactive({ dict_label: '', dict_value: '', sort: 0, status: 0, remark: '' })
const dataRules: FormRules = {
  dict_label: [{ required: true, message: '必填', trigger: 'blur' }],
  dict_value: [{ required: true, message: '必填', trigger: 'blur' }],
}

async function loadTypes() {
  typeLoading.value = true
  try {
    const res = (await fetchDictTypeList({
      page: typeQuery.page,
      size: typeQuery.size,
      dict_type: typeQuery.dict_type || undefined,
    })) as { data?: { list?: DictTypeItem[]; total?: number } }
    typeRows.value = res?.data?.list ?? []
    typeTotal.value = res?.data?.total ?? 0
  } finally {
    typeLoading.value = false
  }
}

async function loadData() {
  if (!currentDictType.value) return
  dataLoading.value = true
  try {
    const res = (await fetchDictDataList({
      page: dataQuery.page,
      size: dataQuery.size,
      dict_type: currentDictType.value,
    })) as { data?: { list?: DictDataItem[]; total?: number } }
    dataRows.value = res?.data?.list ?? []
    dataTotal.value = res?.data?.total ?? 0
  } finally {
    dataLoading.value = false
  }
}

function onTypeRow(row: DictTypeItem | null) {
  if (!row) return
  currentDictType.value = row.dict_type
  dataQuery.page = 1
  loadData()
}

function openTypeCreate() {
  typeMode.value = 'create'
  typeEditId.value = 0
  Object.assign(typeForm, { dict_name: '', dict_type: '', status: 0, remark: '' })
  typeDlg.value = true
}

function openTypeEdit(row: DictTypeItem) {
  typeMode.value = 'edit'
  typeEditId.value = row.id
  Object.assign(typeForm, { dict_name: row.dict_name, dict_type: row.dict_type, status: row.status, remark: '' })
  typeDlg.value = true
}

async function saveType() {
  await typeFormRef.value?.validate(async (ok) => {
    if (!ok) return
    typeSaving.value = true
    try {
      if (typeMode.value === 'create') {
        await createDictType({ ...typeForm })
        ElMessage.success('已创建')
      } else {
        await updateDictType(typeEditId.value, {
          dict_name: typeForm.dict_name,
          status: typeForm.status,
          remark: typeForm.remark,
        })
        ElMessage.success('已保存')
      }
      typeDlg.value = false
      loadTypes()
    } catch {
      /* */
    } finally {
      typeSaving.value = false
    }
  })
}

function delType(row: DictTypeItem) {
  ElMessageBox.confirm(`删除字典类型「${row.dict_type}」？其下数据需先清理。`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteDictType(row.id)
      ElMessage.success('已删除')
      if (currentDictType.value === row.dict_type) currentDictType.value = ''
      loadTypes()
    })
    .catch(() => {})
}

function openDataCreate() {
  if (!currentDictType.value) {
    ElMessage.warning('请先选择字典类型')
    return
  }
  dataMode.value = 'create'
  dataEditId.value = 0
  Object.assign(dataForm, { dict_label: '', dict_value: '', sort: 0, status: 0, remark: '' })
  dataDlg.value = true
}

function openDataEdit(row: DictDataItem) {
  dataMode.value = 'edit'
  dataEditId.value = row.id
  Object.assign(dataForm, {
    dict_label: row.dict_label,
    dict_value: row.dict_value,
    sort: row.sort,
    status: row.status,
    remark: '',
  })
  dataDlg.value = true
}

async function saveData() {
  await dataFormRef.value?.validate(async (ok) => {
    if (!ok) return
    dataSaving.value = true
    try {
      if (dataMode.value === 'create') {
        await createDictData({
          dict_type: currentDictType.value,
          dict_label: dataForm.dict_label,
          dict_value: dataForm.dict_value,
          sort: dataForm.sort,
          status: dataForm.status,
          remark: dataForm.remark,
        })
        ElMessage.success('已创建')
      } else {
        await updateDictData(dataEditId.value, {
          dict_label: dataForm.dict_label,
          dict_value: dataForm.dict_value,
          sort: dataForm.sort,
          status: dataForm.status,
          remark: dataForm.remark,
        })
        ElMessage.success('已保存')
      }
      dataDlg.value = false
      loadData()
    } catch {
      /* */
    } finally {
      dataSaving.value = false
    }
  })
}

function delData(row: DictDataItem) {
  ElMessageBox.confirm(`删除「${row.dict_label}」？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteDictData(row.id)
      ElMessage.success('已删除')
      loadData()
    })
    .catch(() => {})
}

onMounted(loadTypes)
</script>

<style scoped>
.page {
  padding: 8px 0;
}
.card {
  margin-bottom: 16px;
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
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}
.hint {
  margin: 24px 0;
}
</style>
