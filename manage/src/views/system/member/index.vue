<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>会员管理（Member）</span>
          <el-button type="primary" @click="openCreate">新增会员</el-button>
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
        <el-table-column prop="create_time" label="创建时间" width="170" />
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
      :title="mode === 'create' ? '新增会员' : '编辑会员'"
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
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from "vue";
import type { FormInstance, FormRules } from "element-plus";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  fetchMemberList,
  createMember,
  updateMember,
  deleteMember,
  type MemberItem,
} from "@/api/member";

const loading = ref(false);
const saving = ref(false);
const rows = ref<MemberItem[]>([]);
const total = ref(0);
const dlg = ref(false);
const mode = ref<"create" | "edit">("create");
const formRef = ref<FormInstance>();
const editId = ref(0);

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
  ElMessageBox.confirm(`删除会员「${row.username}」？`, "确认", {
    type: "warning",
  })
    .then(async () => {
      await deleteMember(row.id);
      ElMessage.success("已删除");
      await loadList();
    })
    .catch(() => {});
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
.filter {
  margin-bottom: 12px;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
