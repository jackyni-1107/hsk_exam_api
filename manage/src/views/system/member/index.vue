<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>客户管理</span>
          <div class="head-actions">
            <el-button @click="onDownloadTemplate">下载模板</el-button>
            <el-button @click="importDlg = true">导入客户</el-button>
            <el-button type="primary" @click="openCreate">新增客户</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="用户名">
          <el-input v-model="query.username" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" style="width: 120px">
            <el-option label="全部" :value="-1" />
            <el-option label="正常" :value="0" />
            <el-option label="停用" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="ID" width="72" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column
          prop="email"
          label="邮箱"
          min-width="160"
          show-overflow-tooltip
        />
        <el-table-column prop="mobile" label="手机" width="120" />
        <el-table-column label="状态" width="88">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 0 ? 'success' : 'info'"
              size="small"
              >{{ row.status === 0 ? "正常" : "停用" }}</el-tag
            >
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="170" :formatter="formatUtcForDisplay" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)"
              >编辑</el-button
            >
            <el-button link type="danger" @click="onDelete(row)"
              >删除</el-button
            >
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

    <el-dialog
      v-model="dlg"
      :title="mode === 'create' ? '新增客户' : '编辑客户'"
      width="480px"
      destroy-on-close
      @closed="resetForm"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="88px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="mode === 'edit'" />
        </el-form-item>
        <el-form-item
          :label="mode === 'create' ? '密码' : '新密码'"
          prop="password"
        >
          <el-input
            v-model="form.password"
            type="password"
            show-password
            :placeholder="mode === 'edit' ? '不修改请留空' : ''"
          />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="手机">
          <el-input v-model="form.mobile" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="0">正常</el-radio>
            <el-radio :label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlg = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit"
          >保存</el-button
        >
      </template>
    </el-dialog>

    <el-dialog
      v-model="importDlg"
      title="导入客户"
      width="520px"
      destroy-on-close
      @closed="onImportDlgClosed"
    >
      <p class="import-hint">
        请先下载 CSV 模板，按列填写后上传。文件须为 UTF-8 编码；首行为表头；密码需符合系统安全策略；单次最多
        2000 条有效数据行。
      </p>
      <el-upload
        drag
        :limit="1"
        accept=".csv,text/csv"
        :auto-upload="false"
        :on-change="onImportFileChange"
        :on-exceed="onImportExceed"
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">将文件拖到此处，或<em>点击选择</em></div>
        <template #tip>
          <div class="el-upload__tip">仅支持 .csv</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="importDlg = false">取消</el-button>
        <el-button type="primary" :loading="importing" :disabled="!importFile" @click="submitImport"
          >开始导入</el-button
        >
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from "vue";
import type { FormInstance, FormRules, UploadFile } from "element-plus";
import { ElMessage, ElMessageBox } from "element-plus";
import { UploadFilled } from "@element-plus/icons-vue";
import {
  fetchMemberList,
  createMember,
  updateMember,
  deleteMember,
  importMembersCsv,
  downloadMemberImportTemplate,
  type MemberItem,
} from "@/api/member";
import { formatUtcForDisplay } from "@/utils/datetime";

const loading = ref(false);
const saving = ref(false);
const rows = ref<MemberItem[]>([]);
const total = ref(0);
const dlg = ref(false);
const mode = ref<"create" | "edit">("create");
const formRef = ref<FormInstance>();
const editId = ref(0);
const importDlg = ref(false);
const importFile = ref<File | null>(null);
const importing = ref(false);

const query = reactive({ page: 1, size: 10, username: "", status: -1 });
const form = reactive({
  username: "",
  password: "",
  nickname: "",
  email: "",
  mobile: "",
  status: 0,
});

const rules = computed<FormRules>(() => ({
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password:
    mode.value === "create"
      ? [{ required: true, message: "请输入密码", trigger: "blur" }]
      : [],
}));

async function loadList() {
  loading.value = true;
  try {
    const res = (await fetchMemberList({
      page: query.page,
      size: query.size,
      username: query.username || undefined,
      status: query.status < 0 ? undefined : query.status,
    })) as { data?: { list?: MemberItem[]; total?: number } };
    rows.value = res?.data?.list ?? [];
    total.value = res?.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function resetQuery() {
  query.page = 1;
  query.size = 10;
  query.username = "";
  query.status = -1;
  loadList();
}

function resetForm() {
  editId.value = 0;
  form.username = "";
  form.password = "";
  form.nickname = "";
  form.email = "";
  form.mobile = "";
  form.status = 0;
  formRef.value?.clearValidate();
}

function openCreate() {
  mode.value = "create";
  resetForm();
  dlg.value = true;
}

function openEdit(row: MemberItem) {
  mode.value = "edit";
  resetForm();
  editId.value = row.id;
  form.username = row.username;
  form.nickname = row.nickname;
  form.email = row.email;
  form.mobile = row.mobile;
  form.status = row.status;
  dlg.value = true;
}

async function submit() {
  await formRef.value?.validate(async (ok) => {
    if (!ok) return;
    saving.value = true;
    try {
      if (mode.value === "create") {
        await createMember({
          username: form.username,
          password: form.password,
          nickname: form.nickname,
          email: form.email,
          mobile: form.mobile,
          status: form.status,
        });
        ElMessage.success("已创建");
      } else {
        await updateMember(editId.value, {
          password: form.password || undefined,
          nickname: form.nickname,
          email: form.email,
          mobile: form.mobile,
          status: form.status,
        });
        ElMessage.success("已保存");
      }
      dlg.value = false;
      await loadList();
    } catch {
      /* */
    } finally {
      saving.value = false;
    }
  });
}

function onDelete(row: MemberItem) {
  ElMessageBox.confirm(`删除客户「${row.username}」？`, "确认", {
    type: "warning",
  })
    .then(async () => {
      await deleteMember(row.id);
      ElMessage.success("已删除");
      await loadList();
    })
    .catch(() => {});
}

async function onDownloadTemplate() {
  try {
    await downloadMemberImportTemplate();
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : "下载失败");
  }
}

function onImportFileChange(uploadFile: UploadFile) {
  const raw = uploadFile.raw;
  importFile.value = raw instanceof File ? raw : null;
}

function onImportExceed() {
  ElMessage.warning("请先移除已选文件，再选择新文件");
}

function onImportDlgClosed() {
  importFile.value = null;
}

async function submitImport() {
  const f = importFile.value;
  if (!f) {
    ElMessage.warning("请选择 CSV 文件");
    return;
  }
  importing.value = true;
  try {
    const res = (await importMembersCsv(f)) as {
      data?: { total: number; success: number; failed: number; errors?: string[] };
    };
    const d = res?.data;
    if (!d) {
      ElMessage.warning("未返回导入结果");
      return;
    }
    if (d.total === 0) {
      ElMessage.warning("文件中没有有效数据行（空行已忽略）");
      return;
    }
    ElMessage.success(`导入完成：成功 ${d.success} 条，失败 ${d.failed} 条`);
    if (d.errors?.length) {
      await ElMessageBox.alert(d.errors.join("\n"), "失败明细", {
        confirmButtonText: "确定",
        type: d.failed > 0 ? "warning" : "info",
      });
    }
    importDlg.value = false;
    await loadList();
  } catch {
    /* 全局已提示 */
  } finally {
    importing.value = false;
  }
}

onMounted(loadList);
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
.head-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: flex-end;
}
.import-hint {
  margin: 0 0 12px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
}
.filter {
  margin-bottom: 12px;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
