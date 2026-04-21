<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>试卷管理（Mock 卷）</span>
          <el-button type="primary" @click="openImport">导入试卷</el-button>
        </div>
      </template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="等级">
          <el-select
            v-model="query.levelId"
            clearable
            placeholder="全部等级"
            style="width: 200px"
            @clear="loadList"
          >
            <el-option
              v-for="lv in filterLevels"
              :key="lv.id"
              :label="
                lv.level_name || lv.app_level_name || `level_id=${lv.level_id}`
              "
              :value="lv.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="导入状态">
          <el-select
            v-model="query.importStatus"
            placeholder="全部"
            clearable
            style="width: 160px"
            @clear="loadList"
          >
            <el-option label="全部" value="" />
            <el-option label="已导入考试内容" value="imported" />
            <el-option label="未导入" value="not_imported" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table
        v-loading="loading"
        :data="rows"
        border
        stripe
        :row-class-name="tableRowClassName"
      >
        <el-table-column
          label="导入状态"
          width="132"
          align="center"
          fixed="left"
        >
          <template #default="{ row }">
            <el-tag
              :type="row.imported ? 'success' : 'warning'"
              effect="dark"
              size="large"
              class="import-status-tag"
            >
              {{ row.imported ? "已导入" : "未导入" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="Mock卷ID" width="100" />
        <el-table-column prop="level_id" label="等级ID" width="88" />
        <el-table-column
          prop="name"
          label="试卷名称"
          min-width="180"
          show-overflow-tooltip
        />
        <el-table-column prop="score_full" label="满分" width="72" />
        <el-table-column prop="time_full" label="时长(分)" width="96" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag size="small" :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? "发布" : "未发布" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="paper_type" label="类型" width="88" />
        <el-table-column prop="mock_type" label="Mock类型" width="96" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)"
              >详情</el-button
            >
            <el-button link type="primary" @click="openImportForRow(row)"
              >导入</el-button
            >
            <el-button link type="primary" @click="openSettings(row)"
              >听力 HLS</el-button
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
          @size-change="applyPage"
          @current-change="applyPage"
        />
      </div>
    </el-card>

    <!-- 导入 -->
    <el-dialog
      v-model="importDlg"
      title="从 index.json 导入"
      width="640px"
      destroy-on-close
      @open="onImportOpen"
    >
      <el-form label-width="140px">
        <el-form-item label="Mock 卷" required>
          <div class="mock-row">
            <el-select
              v-model="importForm.levelFilter"
              placeholder="等级筛选"
              clearable
              style="width: 160px; margin-right: 8px"
              @change="onImportLevelChange"
            >
              <el-option
                v-for="lv in mockLevels"
                :key="lv.id"
                :label="
                  lv.level_name || lv.app_level_name || String(lv.level_id)
                "
                :value="lv.id"
              />
            </el-select>
            <el-select
              v-model="importForm.mock_examination_paper_id"
              placeholder="选择模拟卷"
              filterable
              style="width: 260px"
            >
              <el-option
                v-for="p in mockPapers"
                :key="p.id"
                :label="`${p.id} · ${p.name}`"
                :value="p.id"
              />
            </el-select>
          </div>
          <div class="mock-hint">
            或直接填写 ID：
            <el-input-number
              v-model="importForm.mock_examination_paper_id"
              :min="1"
              :controls="false"
            />
          </div>
        </el-form-item>
        <el-form-item label="导入方式">
          <el-radio-group v-model="importMode">
            <el-radio label="url">index.json URL</el-radio>
            <el-radio label="json">粘贴 JSON</el-radio>
          </el-radio-group>
        </el-form-item>
        <template v-if="importMode === 'url'">
          <el-form-item label="index_url" required>
            <el-input
              v-model="importForm.index_url"
              type="textarea"
              :rows="2"
              placeholder="完整 index.json 地址"
            />
          </el-form-item>
        </template>
        <template v-else>
          <el-form-item label="level" required>
            <el-input v-model="importForm.level" placeholder="如 hsk1" />
          </el-form-item>
          <el-form-item label="paper_id" required>
            <el-input v-model="importForm.paper_id" placeholder="远程目录 ID" />
          </el-form-item>
          <el-form-item label="source_base_url" required>
            <el-input
              v-model="importForm.source_base_url"
              placeholder="以 / 结尾"
            />
          </el-form-item>
          <el-form-item label="index_json" required>
            <el-input
              v-model="importForm.index_json"
              type="textarea"
              :rows="6"
              placeholder="index.json 全文"
            />
          </el-form-item>
        </template>
        <el-form-item label="听力 HLS 前缀">
          <el-input
            v-model="importForm.audio_hls_prefix"
            placeholder="可选，无首尾 /"
          />
        </el-form-item>
        <el-form-item label="冲突策略">
          <el-radio-group v-model="importForm.conflict_mode">
            <el-radio label="fail">失败（已存在则拒绝）</el-radio>
            <el-radio label="overwrite">覆盖</el-radio>
            <el-radio label="new_copy">新远程目录</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item
          v-if="importForm.conflict_mode === 'new_copy'"
          label="new_paper_id"
          required
        >
          <el-input
            v-model="importForm.new_paper_id"
            placeholder="新的远程 paper_id"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importDlg = false">取消</el-button>
        <el-button type="primary" :loading="importing" @click="submitImport"
          >开始导入</el-button
        >
      </template>
    </el-dialog>

    <!-- 设置（仅已导入 exam 时） -->
    <el-dialog
      v-model="settingsDlg"
      title="听力 HLS 设置"
      width="560px"
      destroy-on-close
    >
      <el-form label-width="168px">
        <el-form-item label="Mock卷ID">
          <span>{{ settingsPaperId }}</span>
        </el-form-item>
        <div class="hint form-top-hint">
          答题时长请在 Mock 卷（time_full 等）中维护；此处仅配置听力 HLS
          相关字段。
        </div>
        <el-form-item label="HLS 前缀">
          <el-input
            v-model="settingsForm.audio_hls_prefix"
            placeholder="无首尾 /"
          />
        </el-form-item>
        <el-form-item label="分片总数">
          <el-input-number
            v-model="settingsForm.audio_hls_segment_count"
            :min="0"
            :controls="false"
            style="width: 100%"
          />
          <div class="hint">0 表示未配置整卷听力 HLS</div>
        </el-form-item>
        <el-form-item label="分片文件名 fmt">
          <el-input
            v-model="settingsForm.audio_hls_segment_pattern"
            placeholder="空则默认 %05d.ts"
          />
        </el-form-item>
        <el-form-item label="密钥对象路径">
          <el-input
            v-model="settingsForm.audio_hls_key_object"
            placeholder="相对 prefix，空表示不加密"
          />
        </el-form-item>
        <el-form-item label="IV 十六进制">
          <el-input
            v-model="settingsForm.audio_hls_iv_hex"
            placeholder="AES-128 IV"
          />
        </el-form-item>
        <el-form-item label="#EXTINF 时长(秒)">
          <el-input-number
            v-model="settingsForm.audio_hls_segment_duration"
            :min="0"
            :precision="3"
            :step="0.1"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="settingsDlg = false">取消</el-button>
        <el-button
          type="primary"
          :loading="savingSettings"
          @click="submitSettings"
          >保存</el-button
        >
      </template>
    </el-dialog>

    <!-- 详情：Mock + 已导入的考试内容 -->
    <el-drawer
      v-model="detailDrawer"
      title="试卷详情"
      size="70%"
      destroy-on-close
    >
      <div v-loading="detailLoading" class="detail-wrap">
        <template v-if="mockDetailPaper">
          <h4 class="sec-title">Mock 卷信息</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{
              mockDetailPaper.id
            }}</el-descriptions-item>
            <el-descriptions-item label="等级ID">{{
              mockDetailPaper.level_id
            }}</el-descriptions-item>
            <el-descriptions-item label="名称" :span="2">{{
              mockDetailPaper.name
            }}</el-descriptions-item>
            <el-descriptions-item label="满分">{{
              mockDetailPaper.score_full
            }}</el-descriptions-item>
            <el-descriptions-item label="时长(分)">{{
              mockDetailPaper.time_full
            }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{
              mockDetailPaper.status
            }}</el-descriptions-item>
            <el-descriptions-item label="试卷类型">{{
              mockDetailPaper.paper_type
            }}</el-descriptions-item>
            <el-descriptions-item label="Mock类型">{{
              mockDetailPaper.mock_type
            }}</el-descriptions-item>
            <el-descriptions-item label="考试内容导入" :span="2">
              <el-tag
                :type="mockDetailPaper.imported ? 'success' : 'warning'"
                effect="dark"
                size="large"
              >
                {{ mockDetailPaper.imported ? "已导入 exam" : "未导入" }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </template>

        <template v-if="examDetail">
          <h4 class="sec-title">已导入的考试内容</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="Mock卷ID">{{
              examDetail.paper.id
            }}</el-descriptions-item>
            <el-descriptions-item label="级别">{{
              examDetail.paper.level
            }}</el-descriptions-item>
            <el-descriptions-item label="远程 paper_id">{{
              examDetail.paper.paper_id
            }}</el-descriptions-item>
            <el-descriptions-item label="标题">{{
              examDetail.paper.title
            }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{
              formatUtcText(examDetail.paper.create_time)
            }}</el-descriptions-item>
            <el-descriptions-item label="资源基址" :span="2">{{
              examDetail.paper.source_base_url
            }}</el-descriptions-item>
            <el-descriptions-item label="HLS 前缀" :span="2">{{
              examDetail.paper.audio_hls_prefix
            }}</el-descriptions-item>
            <el-descriptions-item label="HLS 分片数">{{
              examDetail.paper.audio_hls_segment_count
            }}</el-descriptions-item>
            <el-descriptions-item label="分片 fmt">{{
              examDetail.paper.audio_hls_segment_pattern
            }}</el-descriptions-item>
            <el-descriptions-item label="密钥对象" :span="2">{{
              examDetail.paper.audio_hls_key_object
            }}</el-descriptions-item>
            <el-descriptions-item label="IV hex" :span="2">{{
              examDetail.paper.audio_hls_iv_hex
            }}</el-descriptions-item>
            <el-descriptions-item label="#EXTINF(秒)">{{
              examDetail.paper.audio_hls_segment_duration
            }}</el-descriptions-item>
            <el-descriptions-item label="考前标题" :span="2">{{
              examDetail.paper.prepare_title
            }}</el-descriptions-item>
            <el-descriptions-item label="考前说明" :span="2">{{
              examDetail.paper.prepare_instruction
            }}</el-descriptions-item>
          </el-descriptions>
          <h4 class="sec-title">index.json</h4>
          <pre class="json-preview">{{
            truncate(examDetail.paper.index_json, 8000)
          }}</pre>
          <h4 class="sec-title">大题（{{ examDetail.sections.length }}）</h4>
          <el-collapse>
            <el-collapse-item
              v-for="s in examDetail.sections"
              :key="s.id"
              :title="`${s.sort_order}. ${s.topic_title || s.topic_type} (#${s.id})`"
            >
              <el-descriptions :column="2" size="small" border>
                <el-descriptions-item label="topic_type">{{
                  s.topic_type
                }}</el-descriptions-item>
                <el-descriptions-item label="文件">{{
                  s.topic_items_file
                }}</el-descriptions-item>
                <el-descriptions-item label="part_code">{{
                  s.part_code
                }}</el-descriptions-item>
                <el-descriptions-item label="segment_code">{{
                  s.segment_code
                }}</el-descriptions-item>
              </el-descriptions>
              <div class="sub">topic_json 预览</div>
              <pre class="json-preview sm">{{
                truncate(s.topic_json, 4000)
              }}</pre>
              <div class="sub">题块数 {{ s.blocks?.length ?? 0 }}</div>
            </el-collapse-item>
          </el-collapse>
        </template>
        <el-alert
          v-else-if="detailLoaded && !examDetail"
          type="info"
          :closable="false"
          show-icon
          class="mt"
          title="该 Mock 卷尚未导入考试内容，可使用「导入」从 index.json 写入。"
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import {
  getExamPaperDetail,
  importExamPaper,
  updateExamPaper,
  type ExamPaperDetail,
} from "@/api/exam";
import {
  getMockLevelsList,
  getMockExaminationPapers,
  getMockExaminationPaperDetail,
  type MockLevelItem,
  type MockExaminationPaperItem,
  type MockPaperImportStatusFilter,
} from "@/api/mockAdmin";
import { formatUtcText } from "@/utils/datetime";

const loading = ref(false);
/** 当前页表格数据（Mock 卷） */
const rows = ref<MockExaminationPaperItem[]>([]);
/** 接口返回的全量列表（当前筛选条件下） */
const fullList = ref<MockExaminationPaperItem[]>([]);
const total = ref(0);
const query = reactive({
  page: 1,
  size: 10,
  /** 与 mock examination-paper 的 level_id 一致；空表示全部 */
  levelId: undefined as number | undefined,
  /** 空：全部；imported / not_imported */
  importStatus: "" as MockPaperImportStatusFilter,
});
const filterLevels = ref<MockLevelItem[]>([]);

const importDlg = ref(false);
/** 从表格「导入」打开时预填 Mock 卷（在 onImportOpen 末尾应用，避免被重置清掉） */
const importPreset = ref<{ id: number; level_id: number } | null>(null);
const importing = ref(false);
const importMode = ref<"url" | "json">("url");
const mockLevels = ref<MockLevelItem[]>([]);
const mockPapers = ref<MockExaminationPaperItem[]>([]);
const importForm = reactive({
  levelFilter: undefined as number | undefined,
  mock_examination_paper_id: undefined as number | undefined,
  index_url: "",
  index_json: "",
  level: "",
  paper_id: "",
  source_base_url: "",
  audio_hls_prefix: "",
  conflict_mode: "fail" as "fail" | "overwrite" | "new_copy",
  new_paper_id: "",
});

/** 听力 HLS 表单默认值（新建或未配置时展示/回填） */
type HlsFormDefaults = {
  audio_hls_prefix: string;
  audio_hls_segment_count: number;
  audio_hls_segment_pattern: string;
  audio_hls_key_object: string;
  audio_hls_iv_hex: string;
  audio_hls_segment_duration: number;
};

const HLS_FORM_DEFAULTS: HlsFormDefaults = {
  audio_hls_prefix: "",
  audio_hls_segment_count: 0,
  audio_hls_segment_pattern: "segment_%03d.ts",
  audio_hls_key_object: "static.key",
  audio_hls_iv_hex: "0x0123456789abcdef0123456789abcdef",
  audio_hls_segment_duration: 10,
};

const settingsDlg = ref(false);
const savingSettings = ref(false);
const settingsPaperId = ref(0);
const settingsForm = reactive({ ...HLS_FORM_DEFAULTS });

const detailDrawer = ref(false);
const detailLoading = ref(false);
const detailLoaded = ref(false);
const mockDetailPaper = ref<MockExaminationPaperItem | null>(null);
const examDetail = ref<ExamPaperDetail | null>(null);

function truncate(s: string, max: number) {
  if (!s) return "";
  return s.length <= max ? s : s.slice(0, max) + "\n…（已截断）";
}

function tableRowClassName({ row }: { row: MockExaminationPaperItem }) {
  return row.imported ? "row-paper-imported" : "row-paper-not-imported";
}

function applyPage() {
  const start = (query.page - 1) * query.size;
  rows.value = fullList.value.slice(start, start + query.size);
}

async function loadList() {
  loading.value = true;
  try {
    const params: {
      level_id?: number;
      import_status?: MockPaperImportStatusFilter;
    } = {};
    if (query.levelId !== undefined && query.levelId !== null) {
      params.level_id = query.levelId;
    }
    if (
      query.importStatus === "imported" ||
      query.importStatus === "not_imported"
    ) {
      params.import_status = query.importStatus;
    }
    const res = (await getMockExaminationPapers(
      Object.keys(params).length > 0 ? params : undefined,
    )) as {
      data?: { list?: MockExaminationPaperItem[] };
    };
    fullList.value = res?.data?.list ?? [];
    total.value = fullList.value.length;
    const maxPage = Math.max(1, Math.ceil(total.value / query.size) || 1);
    if (query.page > maxPage) query.page = maxPage;
    applyPage();
  } finally {
    loading.value = false;
  }
}

async function loadFilterLevels() {
  try {
    const res = (await getMockLevelsList()) as {
      data?: { list?: MockLevelItem[] };
    };
    filterLevels.value = res?.data?.list ?? [];
  } catch {
    filterLevels.value = [];
  }
}

function resetQuery() {
  query.page = 1;
  query.size = 10;
  query.levelId = undefined;
  query.importStatus = "";
  loadList();
}

async function onImportOpen() {
  importForm.levelFilter = undefined;
  importForm.mock_examination_paper_id = undefined;
  importForm.index_url = "";
  importForm.index_json = "";
  importForm.level = "";
  importForm.paper_id = "";
  importForm.source_base_url = "";
  importForm.audio_hls_prefix = "";
  importForm.conflict_mode = "fail";
  importForm.new_paper_id = "";
  importMode.value = "url";
  try {
    const res = (await getMockLevelsList()) as {
      data?: { list?: MockLevelItem[] };
    };
    mockLevels.value = res?.data?.list ?? [];
  } catch {
    mockLevels.value = [];
  }
  mockPapers.value = [];
  if (importPreset.value) {
    const p = importPreset.value;
    importPreset.value = null;
    importForm.levelFilter = p.level_id;
    await onImportLevelChange();
    importForm.mock_examination_paper_id = p.id;
  }
}

async function onImportLevelChange() {
  importForm.mock_examination_paper_id = undefined;
  if (!importForm.levelFilter) {
    mockPapers.value = [];
    return;
  }
  try {
    const res = (await getMockExaminationPapers({
      level_id: importForm.levelFilter,
    })) as {
      data?: { list?: MockExaminationPaperItem[] };
    };
    mockPapers.value = res?.data?.list ?? [];
  } catch {
    mockPapers.value = [];
  }
}

function openImport() {
  importDlg.value = true;
}

function openImportForRow(row: MockExaminationPaperItem) {
  importPreset.value = { id: row.id, level_id: row.level_id };
  importDlg.value = true;
}

async function submitImport() {
  const mid = importForm.mock_examination_paper_id;
  if (!mid || mid < 1) {
    ElMessage.warning("请选择或填写 Mock 卷 ID");
    return;
  }
  if (importMode.value === "url" && !importForm.index_url.trim()) {
    ElMessage.warning("请填写 index_url");
    return;
  }
  if (importMode.value === "json") {
    if (
      !importForm.level ||
      !importForm.paper_id ||
      !importForm.source_base_url ||
      !importForm.index_json.trim()
    ) {
      ElMessage.warning(
        "请完整填写 level、paper_id、source_base_url、index_json",
      );
      return;
    }
  }
  if (
    importForm.conflict_mode === "new_copy" &&
    !importForm.new_paper_id.trim()
  ) {
    ElMessage.warning("new_copy 时请填写 new_paper_id");
    return;
  }
  importing.value = true;
  try {
    const payload: Parameters<typeof importExamPaper>[0] = {
      mock_examination_paper_id: mid,
      conflict_mode: importForm.conflict_mode,
      audio_hls_prefix: importForm.audio_hls_prefix || undefined,
      new_paper_id: importForm.new_paper_id || undefined,
    };
    if (importMode.value === "url") {
      payload.index_url = importForm.index_url.trim();
    } else {
      payload.index_json = importForm.index_json;
      payload.level = importForm.level;
      payload.paper_id = importForm.paper_id;
      payload.source_base_url = importForm.source_base_url;
    }
    const res = (await importExamPaper(payload)) as {
      data?: {
        conflict?: boolean;
        existing_examination_paper_id?: number;
        section_count?: number;
        question_count?: number;
      };
    };
    const d = res?.data;
    if (d?.conflict) {
      ElMessage.warning(
        `该 Mock 卷已存在导入记录（mock id=${d.existing_examination_paper_id ?? mid}），未写入。可改用覆盖或新目录。`,
      );
      return;
    }
    ElMessage.success(
      `导入成功：大题 ${d?.section_count ?? 0}，小题 ${d?.question_count ?? 0}`,
    );
    importDlg.value = false;
    await loadList();
  } catch {
    /* axios 已 toast */
  } finally {
    importing.value = false;
  }
}

async function openSettings(row: MockExaminationPaperItem) {
  try {
    const res = (await getExamPaperDetail(row.id)) as {
      data?: ExamPaperDetail;
    };
    const d = res?.data;
    if (!d?.paper) {
      ElMessage.warning("该 Mock 卷尚未导入考试内容，请先做「导入」");
      return;
    }
    settingsPaperId.value = row.id;
    const p = d.paper;
    settingsForm.audio_hls_prefix =
      p.audio_hls_prefix?.trim() || HLS_FORM_DEFAULTS.audio_hls_prefix;
    settingsForm.audio_hls_segment_count =
      p.audio_hls_segment_count ?? HLS_FORM_DEFAULTS.audio_hls_segment_count;
    settingsForm.audio_hls_segment_pattern =
      p.audio_hls_segment_pattern?.trim() ||
      HLS_FORM_DEFAULTS.audio_hls_segment_pattern;
    settingsForm.audio_hls_key_object =
      p.audio_hls_key_object?.trim() || HLS_FORM_DEFAULTS.audio_hls_key_object;
    settingsForm.audio_hls_iv_hex =
      p.audio_hls_iv_hex?.trim() || HLS_FORM_DEFAULTS.audio_hls_iv_hex;
    {
      const dur = p.audio_hls_segment_duration;
      settingsForm.audio_hls_segment_duration =
        dur != null && dur > 0
          ? dur
          : HLS_FORM_DEFAULTS.audio_hls_segment_duration;
    }
    settingsDlg.value = true;
  } catch {
    ElMessage.warning("该 Mock 卷尚未导入考试内容，请先做「导入」");
  }
}

async function submitSettings() {
  savingSettings.value = true;
  try {
    await updateExamPaper({
      id: settingsPaperId.value,
      audio_hls_prefix: settingsForm.audio_hls_prefix,
      audio_hls_segment_count: settingsForm.audio_hls_segment_count,
      audio_hls_segment_pattern: settingsForm.audio_hls_segment_pattern,
      audio_hls_key_object: settingsForm.audio_hls_key_object,
      audio_hls_iv_hex: settingsForm.audio_hls_iv_hex,
      audio_hls_segment_duration: settingsForm.audio_hls_segment_duration,
    });
    ElMessage.success("已保存");
    settingsDlg.value = false;
  } catch {
    /* */
  } finally {
    savingSettings.value = false;
  }
}

async function openDetail(row: MockExaminationPaperItem) {
  mockDetailPaper.value = null;
  examDetail.value = null;
  detailLoaded.value = false;
  detailDrawer.value = true;
  detailLoading.value = true;
  try {
    const mockRes = (await getMockExaminationPaperDetail(row.id)) as {
      data?: { paper?: MockExaminationPaperItem };
    };
    mockDetailPaper.value = mockRes?.data?.paper ?? null;
    try {
      const examRes = (await getExamPaperDetail(row.id)) as {
        data?: ExamPaperDetail;
      };
      examDetail.value = examRes?.data ?? null;
    } catch {
      examDetail.value = null;
    }
  } finally {
    detailLoading.value = false;
    detailLoaded.value = true;
  }
}

onMounted(async () => {
  await loadFilterLevels();
  await loadList();
});
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
.mock-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}
.mock-hint {
  margin-top: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.hint {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.form-top-hint {
  margin: 0 0 12px;
  padding-left: 168px;
  line-height: 1.4;
}
.detail-wrap {
  padding-right: 8px;
}
.sec-title {
  margin: 16px 0 8px;
  font-size: 14px;
}
.json-preview {
  background: var(--el-fill-color-light);
  padding: 8px;
  border-radius: 4px;
  font-size: 12px;
  overflow: auto;
  max-height: 320px;
  white-space: pre-wrap;
  word-break: break-all;
}
.json-preview.sm {
  max-height: 200px;
}
.sub {
  margin: 8px 0 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.mt {
  margin-top: 16px;
}
.import-status-tag {
  font-weight: 600;
  min-width: 4.5em;
  justify-content: center;
}
:deep(.el-table .row-paper-imported > td) {
  background-color: var(--el-color-success-light-9) !important;
}
:deep(.el-table .row-paper-not-imported > td) {
  background-color: var(--el-color-warning-light-9) !important;
}
:deep(.el-table .el-table__body tr.row-paper-imported:hover > td) {
  background-color: var(--el-color-success-light-8) !important;
}
:deep(.el-table .el-table__body tr.row-paper-not-imported:hover > td) {
  background-color: var(--el-color-warning-light-8) !important;
}
</style>
