<template>
  <div class="page">
    <el-card shadow="never">
      <template #header
        ><span>{{ pageTitle }}</span></template
      >
      <template v-if="section === 'template'">
        <el-form :inline="true" class="filter">
          <el-form-item label="编码">
            <el-input v-model="tplQuery.code" clearable style="width: 160px" />
          </el-form-item>
          <el-form-item label="渠道">
            <el-select
              v-model="tplQuery.channel"
              clearable
              style="width: 120px"
            >
              <el-option label="邮件" value="email" />
              <el-option label="短信" value="sms" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadTpl">查询</el-button>
            <el-button
              v-permission="'notification:template_create'"
              type="primary"
              plain
              @click="openTplCreate"
              >新增模板</el-button
            >
          </el-form-item>
        </el-form>
        <el-table v-loading="tplLoading" :data="tplRows" border stripe>
          <el-table-column prop="code" label="编码" width="140" />
          <el-table-column prop="name" label="名称" min-width="120" />
          <el-table-column label="模板类型" width="110">
            <template #default="{ row }">{{
              templateTypeLabel(row.template_type)
            }}</template>
          </el-table-column>
          <el-table-column label="渠道" width="100">
            <template #default="{ row }">{{
              channelLabel(row.channel)
            }}</template>
          </el-table-column>
          <el-table-column
            label="绑定渠道配置"
            min-width="160"
            show-overflow-tooltip
          >
            <template #default="{ row }">{{
              channelConfigLabel(row.channel_config_id)
            }}</template>
          </el-table-column>
          <el-table-column
            prop="third_party_template_id"
            label="第三方模板ID"
            min-width="150"
            show-overflow-tooltip
          />
          <el-table-column
            prop="content"
            label="内容"
            min-width="200"
            show-overflow-tooltip
          />
          <el-table-column label="状态" width="80">
            <template #default="{ row }">
              <el-tag
                :type="row.status === 0 ? 'success' : 'info'"
                size="small"
                >{{ row.status === 0 ? "正常" : "停用" }}</el-tag
              >
            </template>
          </el-table-column>
          <el-table-column
            prop="create_time"
            label="创建时间"
            width="170"
            :formatter="formatUtcForDisplay"
          />
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="{ row }">
              <el-button link type="success" @click="openTplSend(row)"
                >测试发送</el-button
              >
              <el-button
                v-permission="'notification:template_update'"
                link
                type="primary"
                @click="openTplEdit(row)"
                >编辑</el-button
              >
              <el-button
                v-permission="'notification:template_delete'"
                link
                type="danger"
                @click="delTpl(row)"
                >删除</el-button
              >
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
      </template>

      <template v-else-if="section === 'channel'">
        <div class="toolbar">
          <el-button
            v-permission="'notification:channel_config_create'"
            type="primary"
            @click="openChCreate"
            >新增渠道</el-button
          >
          <el-button @click="loadCh">刷新</el-button>
        </div>
        <el-table v-loading="chLoading" :data="chRows" border stripe>
          <el-table-column prop="name" label="名称" min-width="120" />
          <el-table-column label="渠道" width="100">
            <template #default="{ row }">{{
              channelLabel(row.channel)
            }}</template>
          </el-table-column>
          <el-table-column label="提供商" width="120">
            <template #default="{ row }">{{
              providerLabel(row.provider)
            }}</template>
          </el-table-column>
          <el-table-column label="当前" width="80">
            <template #default="{ row }">
              <el-tag v-if="row.is_active === 1" type="success" size="small"
                >是</el-tag
              >
              <span v-else>否</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="create_time"
            label="创建时间"
            width="170"
            :formatter="formatUtcForDisplay"
          />
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="{ row }">
              <el-button
                v-permission="'notification:channel_config_set_active'"
                link
                type="warning"
                @click="setChActive(row)"
                >设为当前</el-button
              >
              <el-button
                v-permission="'notification:channel_config_update'"
                link
                type="primary"
                @click="openChEdit(row)"
                >编辑</el-button
              >
              <el-button
                v-permission="'notification:channel_config_delete'"
                link
                type="danger"
                @click="delCh(row)"
                >删除</el-button
              >
            </template>
          </el-table-column>
        </el-table>
      </template>

      <template v-else>
        <el-form :inline="true" class="filter">
          <el-form-item label="渠道">
            <el-select
              v-model="logQuery.channel"
              clearable
              style="width: 120px"
            >
              <el-option label="邮件" value="email" />
              <el-option label="短信" value="sms" />
            </el-select>
          </el-form-item>
          <el-form-item label="收件人">
            <el-input
              v-model="logQuery.recipient"
              clearable
              style="width: 160px"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadLog">查询</el-button>
          </el-form-item>
        </el-form>
        <el-table v-loading="logLoading" :data="logRows" border stripe>
          <el-table-column prop="template_code" label="模板" width="140" />
          <el-table-column label="渠道" width="90">
            <template #default="{ row }">{{
              channelLabel(row.channel)
            }}</template>
          </el-table-column>
          <el-table-column prop="recipient" label="收件人" min-width="160" />
          <el-table-column label="状态" width="80">
            <template #default="{ row }">
              <el-tag
                :type="row.status === 0 ? 'success' : 'danger'"
                size="small"
                >{{ row.status === 0 ? "成功" : "失败" }}</el-tag
              >
            </template>
          </el-table-column>
          <el-table-column
            prop="error_msg"
            label="错误"
            min-width="160"
            show-overflow-tooltip
          />
          <el-table-column
            prop="create_time"
            label="时间"
            width="170"
            :formatter="formatUtcForDisplay"
          />
        </el-table>
        <el-pagination
          v-model:current-page="logQuery.page"
          v-model:page-size="logQuery.size"
          class="pager"
          :total="logTotal"
          layout="total, prev, pager, next"
          @current-change="loadLog"
        />
      </template>
    </el-card>

    <el-dialog
      v-model="tplDlg"
      :title="tplMode === 'create' ? '新增模板' : '编辑模板'"
      width="560px"
      destroy-on-close
    >
      <el-form
        ref="tplFormRef"
        :model="tplForm"
        :rules="tplRules"
        label-width="88px"
      >
        <el-form-item label="编码" prop="code">
          <el-input v-model="tplForm.code" :disabled="tplMode === 'edit'" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="tplForm.name" />
        </el-form-item>
        <el-form-item label="模板类型" prop="template_type">
          <el-radio-group
            v-model="tplForm.template_type"
            @change="onTplTemplateTypeChange"
          >
            <el-radio :label="1">系统模板</el-radio>
            <el-radio :label="2">第三方模板</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="渠道" prop="channel">
          <el-select
            v-model="tplForm.channel"
            style="width: 100%"
            @change="onTplChannelChange"
          >
            <el-option label="邮件" value="email" />
            <el-option label="短信" value="sms" />
          </el-select>
        </el-form-item>
        <el-form-item label="通知渠道" prop="channel_config_id">
          <el-select v-model="tplForm.channel_config_id" style="width: 100%">
            <el-option
              v-for="item in tplAvailableChannelConfigs"
              :key="item.id"
              :label="`${item.name}（${providerLabel(item.provider)}）`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="tplForm.template_type === 1"
          label="内容"
          prop="content"
        >
          <el-input
            v-model="tplForm.content"
            type="textarea"
            :rows="5"
            placeholder="可使用 {{变量名}} 作为占位符"
          />
        </el-form-item>
        <template v-else>
          <el-form-item label="第三方模板ID" prop="third_party_template_id">
            <el-input v-model="tplForm.third_party_template_id" />
          </el-form-item>
          <el-form-item label="第三方参数">
            <div class="variable-config-list">
              <div
                v-for="(item, index) in tplThirdPartyParamRows"
                :key="index"
                class="variable-config-row"
              >
                <el-input v-model="item.key" placeholder="参数名，如 code" />
                <el-input
                  v-model="item.value"
                  placeholder="参数值，如 123456"
                />
                <el-button
                  link
                  type="danger"
                  @click="removeTplThirdPartyParam(index)"
                  >删除</el-button
                >
              </div>
              <el-button link type="primary" @click="addTplThirdPartyParam"
                >添加参数</el-button
              >
            </div>
          </el-form-item>
        </template>
        <el-form-item label="变量配置">
          <div class="variable-config-list">
            <div
              v-for="(item, index) in tplVariableRows"
              :key="index"
              class="variable-config-row"
            >
              <el-input v-model="item.key" placeholder="变量名，如 name" />
              <el-input v-model="item.value" placeholder="示例/说明，如 张三" />
              <el-button
                link
                type="primary"
                :disabled="!item.key.trim()"
                @click="insertTplVariable(item.key)"
                >插入</el-button
              >
              <el-button link type="danger" @click="removeTplVariable(index)"
                >删除</el-button
              >
            </div>
            <el-button link type="primary" @click="addTplVariable"
              >添加变量</el-button
            >
          </div>
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
        <el-button type="primary" :loading="tplSaving" @click="saveTpl"
          >保存</el-button
        >
      </template>
    </el-dialog>

    <el-dialog
      v-model="chDlg"
      :title="chMode === 'create' ? '新增渠道' : '编辑渠道'"
      width="520px"
      destroy-on-close
    >
      <el-form
        ref="chFormRef"
        :model="chForm"
        :rules="chRules"
        label-width="88px"
      >
        <el-form-item label="渠道" prop="channel">
          <el-select
            v-model="chForm.channel"
            :disabled="chMode === 'edit'"
            style="width: 100%"
            @change="onChChannelChange"
          >
            <el-option label="邮件" value="email" />
            <el-option label="短信" value="sms" />
          </el-select>
        </el-form-item>
        <el-form-item label="提供商" prop="provider">
          <el-select
            v-model="chForm.provider"
            :disabled="chMode === 'edit'"
            style="width: 100%"
            @change="resetChConfig"
          >
            <template v-if="chForm.channel === 'email'">
              <el-option label="SMTP" value="smtp" />
              <el-option label="SendGrid" value="sendgrid" />
            </template>
            <template v-else-if="chForm.channel === 'sms'">
              <el-option label="阿里云" value="aliyun" />
              <el-option label="腾讯云" value="tencent" />
            </template>
          </el-select>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="chForm.name" />
        </el-form-item>
        <template
          v-if="chForm.channel === 'email' && chForm.provider === 'smtp'"
        >
          <el-divider content-position="left">SMTP 配置</el-divider>
          <el-form-item label="Host" prop="host">
            <el-input v-model="chConfig.host" placeholder="smtp.example.com" />
          </el-form-item>
          <el-form-item label="Port" prop="port">
            <el-input-number
              v-model="chConfig.port"
              :min="1"
              :max="65535"
              controls-position="right"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="账号">
            <el-input v-model="chConfig.user" autocomplete="off" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input
              v-model="chConfig.pass"
              type="password"
              show-password
              autocomplete="new-password"
            />
          </el-form-item>
          <el-form-item label="发件人">
            <el-input
              v-model="chConfig.from"
              placeholder="noreply@example.com"
            />
          </el-form-item>
        </template>
        <template
          v-else-if="
            chForm.channel === 'email' && chForm.provider === 'sendgrid'
          "
        >
          <el-divider content-position="left">SendGrid 配置</el-divider>
          <el-form-item label="API Key" prop="sendgrid_api_key">
            <el-input
              v-model="chConfig.sendgrid_api_key"
              type="password"
              show-password
              autocomplete="new-password"
            />
          </el-form-item>
          <el-form-item label="发件邮箱" prop="from">
            <el-input
              v-model="chConfig.from"
              placeholder="noreply@example.com"
            />
          </el-form-item>
          <el-form-item label="发件人名称">
            <el-input
              v-model="chConfig.from_name"
              placeholder="例如：HSK考试中心"
            />
          </el-form-item>
        </template>

        <template
          v-else-if="chForm.channel === 'sms' && chForm.provider === 'aliyun'"
        >
          <el-divider content-position="left">阿里云短信配置</el-divider>
          <el-form-item label="AccessKey" prop="access_key">
            <el-input v-model="chConfig.access_key" autocomplete="off" />
          </el-form-item>
          <el-form-item label="SecretKey" prop="secret_key">
            <el-input
              v-model="chConfig.secret_key"
              type="password"
              show-password
              autocomplete="new-password"
            />
          </el-form-item>
          <el-form-item label="短信签名">
            <el-input v-model="chConfig.sign_name" />
          </el-form-item>
          <el-form-item label="模板编码">
            <el-input v-model="chConfig.template_code" />
          </el-form-item>
        </template>

        <template
          v-else-if="chForm.channel === 'sms' && chForm.provider === 'tencent'"
        >
          <el-divider content-position="left">腾讯云短信配置</el-divider>
          <el-form-item label="SecretId" prop="secret_id">
            <el-input v-model="chConfig.secret_id" autocomplete="off" />
          </el-form-item>
          <el-form-item label="SecretKey" prop="secret_key">
            <el-input
              v-model="chConfig.secret_key"
              type="password"
              show-password
              autocomplete="new-password"
            />
          </el-form-item>
          <el-form-item label="SdkAppId">
            <el-input v-model="chConfig.sdk_app_id" />
          </el-form-item>
          <el-form-item label="短信签名">
            <el-input v-model="chConfig.sign_name" />
          </el-form-item>
          <el-form-item label="模板 ID">
            <el-input v-model="chConfig.template_id" />
          </el-form-item>
        </template>

        <el-divider content-position="left">扩展配置</el-divider>
        <div class="extra-config-list">
          <div
            v-for="(item, index) in chExtraConfigRows"
            :key="index"
            class="extra-config-row"
          >
            <el-input v-model="item.key" placeholder="字段名" />
            <el-input v-model="item.value" placeholder="字段值" />
            <el-button link type="danger" @click="removeChExtraConfig(index)"
              >删除</el-button
            >
          </div>
          <el-button link type="primary" @click="addChExtraConfig"
            >添加扩展字段</el-button
          >
        </div>
      </el-form>
      <template #footer>
        <el-button @click="chDlg = false">取消</el-button>
        <el-button type="primary" :loading="chSaving" @click="saveCh"
          >保存</el-button
        >
      </template>
    </el-dialog>

    <el-dialog
      v-model="sendDlg"
      title="测试发送"
      width="560px"
      destroy-on-close
    >
      <el-form :model="sendForm" label-width="100px">
        <el-form-item label="模板编码">
          <el-input v-model="sendForm.template_code" disabled />
        </el-form-item>
        <el-form-item label="模板类型">
          <el-input
            :model-value="templateTypeLabel(sendForm.template_type)"
            disabled
          />
        </el-form-item>
        <el-form-item v-if="sendForm.template_type === 2" label="第三方模板ID">
          <el-input v-model="sendForm.third_party_template_id" disabled />
        </el-form-item>
        <el-form-item label="渠道">
          <el-input :model-value="channelLabel(sendForm.channel)" disabled />
        </el-form-item>
        <el-form-item label="收件人" required>
          <el-input v-model="sendForm.recipient" />
        </el-form-item>
        <el-form-item
          :label="sendForm.template_type === 2 ? '第三方参数/变量' : '变量配置'"
        >
          <div class="variable-config-list">
            <div
              v-for="(item, index) in sendVariableRows"
              :key="index"
              class="send-variable-row"
            >
              <el-input v-model="item.key" placeholder="变量名，如 name" />
              <el-input v-model="item.value" placeholder="变量值，如 张三" />
              <el-button link type="danger" @click="removeSendVariable(index)"
                >删除</el-button
              >
            </div>
            <el-button link type="primary" @click="addSendVariable"
              >添加变量</el-button
            >
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="sendDlg = false">取消</el-button>
        <el-button type="primary" :loading="sendLoading" @click="doSend"
          >发送</el-button
        >
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, onMounted } from "vue";
import type { FormInstance, FormRules } from "element-plus";
import { ElMessage, ElMessageBox } from "element-plus";
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
} from "@/api/notification";
import { formatUtcForDisplay } from "@/utils/datetime";

type NotificationSection = "template" | "channel" | "log";

const props = withDefaults(
  defineProps<{
    section?: NotificationSection;
  }>(),
  {
    section: "template",
  },
);

const pageTitle = computed(() => {
  const titleMap: Record<NotificationSection, string> = {
    template: "通知模板",
    channel: "通知渠道",
    log: "发送记录",
  };
  return titleMap[props.section];
});

const section = computed(() => props.section);
const tplLoading = ref(false);
const chLoading = ref(false);
const logLoading = ref(false);
const tplRows = ref<TemplateItem[]>([]);
const tplTotal = ref(0);
const chRows = ref<ChannelConfigItem[]>([]);
const logRows = ref<unknown[]>([]);
const logTotal = ref(0);

const tplQuery = reactive({ page: 1, size: 10, code: "", channel: "" });
const logQuery = reactive({ page: 1, size: 10, channel: "", recipient: "" });

const tplDlg = ref(false);
const tplMode = ref<"create" | "edit">("create");
const tplSaving = ref(false);
const tplFormRef = ref<FormInstance>();
const tplEditId = ref(0);
const tplForm = reactive({
  code: "",
  name: "",
  channel: "",
  channel_config_id: 0,
  template_type: 1,
  content: "",
  third_party_template_id: "",
  third_party_template_params: "",
  status: 0,
});
const tplVariableRows = ref<{ key: string; value: string }[]>([]);
const tplThirdPartyParamRows = ref<{ key: string; value: string }[]>([]);
const tplChannelConfigOptions = ref<ChannelConfigItem[]>([]);
const tplRules: FormRules = {
  code: [{ required: true, message: "必填", trigger: "blur" }],
  name: [{ required: true, message: "必填", trigger: "blur" }],
  channel: [{ required: true, message: "必填", trigger: "blur" }],
  channel_config_id: [{ required: true, message: "必填", trigger: "change" }],
};

const chDlg = ref(false);
const chMode = ref<"create" | "edit">("create");
const chSaving = ref(false);
const chFormRef = ref<FormInstance>();
const chEditId = ref(0);
const chForm = reactive({ channel: "", provider: "", name: "" });
const chConfig = reactive({
  host: "",
  port: 25,
  user: "",
  pass: "",
  from: "",
  from_name: "",
  sendgrid_api_key: "",
  access_key: "",
  secret_key: "",
  sign_name: "",
  template_code: "",
  secret_id: "",
  sdk_app_id: "",
  template_id: "",
});
const chExtraConfigRows = ref<{ key: string; value: string }[]>([]);
const chKnownConfigKeys = new Set(Object.keys(chConfig));
const chRules: FormRules = {
  channel: [{ required: true, message: "必填", trigger: "blur" }],
  provider: [{ required: true, message: "必填", trigger: "blur" }],
  name: [{ required: true, message: "必填", trigger: "blur" }],
};

const sendDlg = ref(false);
const sendForm = reactive({
  template_code: "",
  template_type: 1,
  third_party_template_id: "",
  channel: "email",
  recipient: "",
});
const sendVariableRows = ref<{ key: string; value: string }[]>([]);
const sendLoading = ref(false);

function channelLabel(channel: string) {
  const map: Record<string, string> = {
    email: "邮件",
    sms: "短信",
  };
  return map[channel] ?? channel;
}

function providerLabel(provider: string) {
  const map: Record<string, string> = {
    smtp: "SMTP",
    sendgrid: "SendGrid",
    aliyun: "阿里云",
    tencent: "腾讯云",
  };
  return map[provider] ?? provider;
}

function templateTypeLabel(templateType: number) {
  return templateType === 2 ? "第三方模板" : "系统模板";
}

const tplAvailableChannelConfigs = computed(() =>
  tplChannelConfigOptions.value.filter(
    (item) => !tplForm.channel || item.channel === tplForm.channel,
  ),
);

function channelConfigLabel(channelConfigId: number) {
  const target = tplChannelConfigOptions.value.find(
    (item) => item.id === channelConfigId,
  );
  if (!target) return String(channelConfigId ?? "");
  return `${target.name}（${providerLabel(target.provider)}）`;
}

async function loadTpl() {
  tplLoading.value = true;
  try {
    const res = (await fetchTemplates({
      page: tplQuery.page,
      size: tplQuery.size,
      code: tplQuery.code || undefined,
      channel: tplQuery.channel || undefined,
    })) as { data?: { list?: TemplateItem[]; total?: number } };
    tplRows.value = res?.data?.list ?? [];
    tplTotal.value = res?.data?.total ?? 0;
  } finally {
    tplLoading.value = false;
  }
}

async function loadTplChannelOptions() {
  const res = (await fetchChannelConfigs()) as {
    data?: { list?: ChannelConfigItem[] };
  };
  tplChannelConfigOptions.value = res?.data?.list ?? [];
}

async function loadCh() {
  chLoading.value = true;
  try {
    const res = (await fetchChannelConfigs()) as {
      data?: { list?: ChannelConfigItem[] };
    };
    chRows.value = res?.data?.list ?? [];
  } finally {
    chLoading.value = false;
  }
}

async function loadLog() {
  logLoading.value = true;
  try {
    const res = (await fetchNotificationLogs({
      page: logQuery.page,
      size: logQuery.size,
      channel: logQuery.channel || undefined,
      recipient: logQuery.recipient || undefined,
    })) as { data?: { list?: unknown[]; total?: number } };
    logRows.value = res?.data?.list ?? [];
    logTotal.value = res?.data?.total ?? 0;
  } finally {
    logLoading.value = false;
  }
}

function stringifyTplVariableValue(value: unknown) {
  if (value === undefined || value === null) return "";
  if (typeof value === "string") return value;
  if (typeof value === "number" || typeof value === "boolean")
    return String(value);
  return JSON.stringify(value);
}

function fillTplVariables(variables?: string) {
  const value = variables?.trim() ?? "";
  if (!value) {
    tplVariableRows.value = [];
    return;
  }

  try {
    const parsed = JSON.parse(value) as unknown;
    if (parsed && typeof parsed === "object" && !Array.isArray(parsed)) {
      tplVariableRows.value = Object.entries(
        parsed as Record<string, unknown>,
      ).map(([key, val]) => ({
        key,
        value: stringifyTplVariableValue(val),
      }));
      return;
    }
    if (Array.isArray(parsed)) {
      tplVariableRows.value = parsed.map((item) => ({
        key: stringifyTplVariableValue(item),
        value: "",
      }));
      return;
    }
  } catch {
    /* Historical data can be comma-separated instead of JSON. */
  }

  tplVariableRows.value = value
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean)
    .map((key) => ({ key, value: "" }));
}

function buildTplVariables() {
  const rows = tplVariableRows.value
    .map((row) => ({ key: row.key.trim(), value: row.value.trim() }))
    .filter((row) => row.key);

  if (!rows.length) return "";
  if (rows.some((row) => row.value)) {
    return JSON.stringify(
      Object.fromEntries(rows.map((row) => [row.key, row.value])),
    );
  }
  return rows.map((row) => row.key).join(",");
}

function addTplVariable() {
  tplVariableRows.value.push({ key: "", value: "" });
}

function removeTplVariable(index: number) {
  tplVariableRows.value.splice(index, 1);
}

function insertTplVariable(key: string) {
  const variable = key.trim();
  if (!variable) return;
  const placeholder = `{{${variable}}}`;
  tplForm.content = tplForm.content
    ? `${tplForm.content}${placeholder}`
    : placeholder;
}

function openTplCreate() {
  tplMode.value = "create";
  tplEditId.value = 0;
  Object.assign(tplForm, {
    code: "",
    name: "",
    channel: "email",
    channel_config_id: 0,
    template_type: 1,
    content: "",
    third_party_template_id: "",
    third_party_template_params: "",
    status: 0,
  });
  tplVariableRows.value = [];
  tplThirdPartyParamRows.value = [];
  onTplChannelChange();
  tplDlg.value = true;
}

function openTplEdit(row: TemplateItem) {
  tplMode.value = "edit";
  tplEditId.value = row.id;
  Object.assign(tplForm, {
    code: row.code,
    name: row.name,
    channel: row.channel,
    channel_config_id: row.channel_config_id,
    template_type: row.template_type || 1,
    content: row.content,
    third_party_template_id: row.third_party_template_id || "",
    third_party_template_params: row.third_party_template_params || "",
    status: row.status,
  });
  fillTplVariables(row.variables);
  fillTplThirdPartyParams(row.third_party_template_params);
  tplDlg.value = true;
}

async function saveTpl() {
  await tplFormRef.value?.validate(async (ok) => {
    if (!ok) return;
    tplSaving.value = true;
    try {
      const variables = buildTplVariables();
      const thirdPartyTemplateParams = buildTplThirdPartyParams();
      if (tplForm.template_type === 1 && !tplForm.content.trim()) {
        ElMessage.warning("系统模板内容不能为空");
        return;
      }
      if (
        tplForm.template_type === 2 &&
        !tplForm.third_party_template_id.trim()
      ) {
        ElMessage.warning("第三方模板ID不能为空");
        return;
      }
      if (tplMode.value === "create") {
        await createTemplate({
          ...tplForm,
          variables,
          third_party_template_params: thirdPartyTemplateParams,
        });
        ElMessage.success("已创建");
      } else {
        await updateTemplate(tplEditId.value, {
          name: tplForm.name,
          channel: tplForm.channel,
          channel_config_id: tplForm.channel_config_id,
          template_type: tplForm.template_type,
          content: tplForm.content,
          third_party_template_id: tplForm.third_party_template_id,
          third_party_template_params: thirdPartyTemplateParams,
          variables,
          status: tplForm.status,
        });
        ElMessage.success("已保存");
      }
      tplDlg.value = false;
      loadTpl();
    } catch {
      /* */
    } finally {
      tplSaving.value = false;
    }
  });
}

function onTplChannelChange() {
  const available = tplAvailableChannelConfigs.value;
  if (!available.some((item) => item.id === tplForm.channel_config_id)) {
    tplForm.channel_config_id = available[0]?.id ?? 0;
  }
}

function onTplTemplateTypeChange() {
  if (tplForm.template_type === 1) {
    tplForm.third_party_template_id = "";
    tplThirdPartyParamRows.value = [];
    return;
  }
  tplForm.content = "";
}

function fillTplThirdPartyParams(params?: string) {
  const value = params?.trim() ?? "";
  if (!value) {
    tplThirdPartyParamRows.value = [];
    return;
  }
  try {
    const parsed = JSON.parse(value) as Record<string, unknown>;
    tplThirdPartyParamRows.value = Object.entries(parsed).map(([key, val]) => ({
      key,
      value: stringifyTplVariableValue(val),
    }));
  } catch {
    tplThirdPartyParamRows.value = [];
  }
}

function buildTplThirdPartyParams() {
  const rows = tplThirdPartyParamRows.value
    .map((row) => ({ key: row.key.trim(), value: row.value }))
    .filter((row) => row.key);
  if (!rows.length) return "";
  return JSON.stringify(
    Object.fromEntries(rows.map((row) => [row.key, row.value])),
  );
}

function addTplThirdPartyParam() {
  tplThirdPartyParamRows.value.push({ key: "", value: "" });
}

function removeTplThirdPartyParam(index: number) {
  tplThirdPartyParamRows.value.splice(index, 1);
}

function delTpl(row: TemplateItem) {
  ElMessageBox.confirm(`删除模板「${row.code}」？`, "确认", { type: "warning" })
    .then(async () => {
      await deleteTemplate(row.id);
      ElMessage.success("已删除");
      loadTpl();
    })
    .catch(() => {});
}

function fillSendVariablesFromTemplate(variables?: string) {
  const value = variables?.trim() ?? "";
  if (!value) {
    sendVariableRows.value = [];
    return;
  }

  try {
    const parsed = JSON.parse(value) as unknown;
    if (parsed && typeof parsed === "object" && !Array.isArray(parsed)) {
      sendVariableRows.value = Object.entries(
        parsed as Record<string, unknown>,
      ).map(([key, val]) => ({
        key,
        value: stringifyTplVariableValue(val),
      }));
      return;
    }
    if (Array.isArray(parsed)) {
      sendVariableRows.value = parsed.map((item) => ({
        key: stringifyTplVariableValue(item),
        value: "",
      }));
      return;
    }
  } catch {
    /* Historical data can be comma-separated instead of JSON. */
  }

  sendVariableRows.value = value
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean)
    .map((key) => ({ key, value: "" }));
}

function parseVariableRowsFromText(value?: string) {
  const content = value?.trim() ?? "";
  if (!content) return [] as { key: string; value: string }[];
  try {
    const parsed = JSON.parse(content) as unknown;
    if (parsed && typeof parsed === "object" && !Array.isArray(parsed)) {
      return Object.entries(parsed as Record<string, unknown>).map(
        ([key, val]) => ({
          key,
          value: stringifyTplVariableValue(val),
        }),
      );
    }
    if (Array.isArray(parsed)) {
      return parsed.map((item) => ({
        key: stringifyTplVariableValue(item),
        value: "",
      }));
    }
  } catch {
    /* Historical data can be comma-separated instead of JSON. */
  }
  return content
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean)
    .map((key) => ({ key, value: "" }));
}

function fillSendParamsByTemplate(row: TemplateItem) {
  const merged = new Map<string, string>();
  // Third-party template params are treated as defaults.
  for (const item of parseVariableRowsFromText(
    row.third_party_template_params,
  )) {
    if (!item.key.trim()) continue;
    merged.set(item.key.trim(), item.value);
  }
  // Variables overwrite defaults when same key exists.
  for (const item of parseVariableRowsFromText(row.variables)) {
    if (!item.key.trim()) continue;
    merged.set(item.key.trim(), item.value);
  }
  sendVariableRows.value = Array.from(merged.entries()).map(([key, value]) => ({
    key,
    value,
  }));
}

function openTplSend(row: TemplateItem) {
  Object.assign(sendForm, {
    template_code: row.code,
    template_type: row.template_type || 1,
    third_party_template_id: row.third_party_template_id || "",
    channel: row.channel,
    recipient: "",
  });
  fillSendParamsByTemplate(row);
  sendDlg.value = true;
}

function resetChConfig() {
  Object.assign(chConfig, {
    host: "",
    port: 25,
    user: "",
    pass: "",
    from: "",
    from_name: "",
    sendgrid_api_key: "",
    access_key: "",
    secret_key: "",
    sign_name: "",
    template_code: "",
    secret_id: "",
    sdk_app_id: "",
    template_id: "",
  });
  chExtraConfigRows.value = [];
}

function parseChConfigJson(configJson?: string) {
  if (!configJson) return {};
  try {
    const parsed = JSON.parse(configJson) as unknown;
    if (parsed && typeof parsed === "object" && !Array.isArray(parsed)) {
      return parsed as Record<string, unknown>;
    }
  } catch {
    /* Keep the form usable even if historical config is malformed. */
  }
  return {};
}

function stringifyChConfigValue(value: unknown) {
  if (value === undefined || value === null) return "";
  if (typeof value === "string") return value;
  if (typeof value === "number" || typeof value === "boolean")
    return String(value);
  return JSON.stringify(value);
}

function parseChExtraConfigValue(value: string) {
  const trimmed = value.trim();
  if (!trimmed) return "";
  if (
    !["{", "[", '"'].includes(trimmed[0]) &&
    !["true", "false", "null"].includes(trimmed) &&
    Number.isNaN(Number(trimmed))
  ) {
    return value;
  }
  try {
    return JSON.parse(trimmed) as unknown;
  } catch {
    return value;
  }
}

function fillChConfig(configJson?: string) {
  resetChConfig();
  const data = parseChConfigJson(configJson);
  for (const key of Object.keys(chConfig) as Array<keyof typeof chConfig>) {
    if (!(key in data)) continue;
    if (key === "port") {
      const port = Number(data[key]);
      chConfig.port = Number.isFinite(port) && port > 0 ? port : 25;
    } else {
      chConfig[key] = stringifyChConfigValue(data[key]);
    }
  }
  chExtraConfigRows.value = Object.entries(data)
    .filter(([key]) => !chKnownConfigKeys.has(key))
    .map(([key, value]) => ({ key, value: stringifyChConfigValue(value) }));
}

function addChExtraConfig() {
  chExtraConfigRows.value.push({ key: "", value: "" });
}

function removeChExtraConfig(index: number) {
  chExtraConfigRows.value.splice(index, 1);
}

function onChChannelChange() {
  chForm.provider = chForm.channel === "sms" ? "aliyun" : "smtp";
  resetChConfig();
}

function buildChConfigJson() {
  const data: Record<string, unknown> = {};
  const setString = (key: keyof typeof chConfig) => {
    const value = chConfig[key];
    if (typeof value === "string" && value.trim()) data[key] = value.trim();
  };

  if (chForm.channel === "email" && chForm.provider === "smtp") {
    setString("host");
    data.port = chConfig.port || 25;
    setString("user");
    setString("pass");
    setString("from");
  } else if (chForm.channel === "email" && chForm.provider === "sendgrid") {
    setString("sendgrid_api_key");
    setString("from");
    setString("from_name");
  } else if (chForm.channel === "sms" && chForm.provider === "aliyun") {
    setString("access_key");
    setString("secret_key");
    setString("sign_name");
    setString("template_code");
  } else if (chForm.channel === "sms" && chForm.provider === "tencent") {
    setString("secret_id");
    setString("secret_key");
    setString("sdk_app_id");
    setString("sign_name");
    setString("template_id");
  }

  for (const row of chExtraConfigRows.value) {
    const key = row.key.trim();
    if (!key || chKnownConfigKeys.has(key)) continue;
    data[key] = parseChExtraConfigValue(row.value);
  }

  return JSON.stringify(data);
}

function openChCreate() {
  chMode.value = "create";
  chEditId.value = 0;
  Object.assign(chForm, { channel: "email", provider: "smtp", name: "" });
  resetChConfig();
  chDlg.value = true;
}

function openChEdit(row: ChannelConfigItem) {
  chMode.value = "edit";
  chEditId.value = row.id;
  Object.assign(chForm, {
    channel: row.channel,
    provider: row.provider,
    name: row.name,
  });
  fillChConfig(row.config_json);
  chDlg.value = true;
}

async function saveCh() {
  await chFormRef.value?.validate(async (ok) => {
    if (!ok) return;
    chSaving.value = true;
    try {
      const configJson = buildChConfigJson();
      if (chMode.value === "create") {
        await createChannelConfig({ ...chForm, config_json: configJson });
        ElMessage.success("已创建");
      } else {
        await updateChannelConfig(chEditId.value, {
          name: chForm.name,
          config_json: configJson,
        });
        ElMessage.success("已保存");
      }
      chDlg.value = false;
      loadCh();
    } catch {
      /* */
    } finally {
      chSaving.value = false;
    }
  });
}

function delCh(row: ChannelConfigItem) {
  ElMessageBox.confirm(`删除渠道「${row.name}」？`, "确认", { type: "warning" })
    .then(async () => {
      await deleteChannelConfig(row.id);
      ElMessage.success("已删除");
      loadCh();
    })
    .catch(() => {});
}

function setChActive(row: ChannelConfigItem) {
  ElMessageBox.confirm(`将「${row.name}」设为当前渠道？`, "确认", {
    type: "warning",
  })
    .then(async () => {
      await setActiveChannelConfig(row.id);
      ElMessage.success("已设置");
      loadCh();
    })
    .catch(() => {});
}

function addSendVariable() {
  sendVariableRows.value.push({ key: "", value: "" });
}

function removeSendVariable(index: number) {
  sendVariableRows.value.splice(index, 1);
}

function buildSendVariables() {
  const rows = sendVariableRows.value
    .map((row) => ({ key: row.key.trim(), value: row.value }))
    .filter((row) => row.key);

  if (!rows.length) return undefined;
  return JSON.stringify(
    Object.fromEntries(rows.map((row) => [row.key, row.value])),
  );
}

async function doSend() {
  if (!sendForm.template_code || !sendForm.channel || !sendForm.recipient) {
    ElMessage.warning("请填写模板、渠道、收件人");
    return;
  }
  sendLoading.value = true;
  try {
    const res = (await sendNotification({
      template_code: sendForm.template_code,
      channel: sendForm.channel,
      recipient: sendForm.recipient,
      variables: buildSendVariables(),
    })) as { data?: { ok?: boolean } };
    if (res?.data?.ok) ElMessage.success("已提交发送");
    else ElMessage.success("请求已完成");
    sendDlg.value = false;
  } catch (error: unknown) {
    const err = error as {
      message?: string;
      response?: { data?: { message?: string } };
    };
    const msg =
      err?.response?.data?.message ||
      err?.message ||
      "发送失败，请查看后端日志";
    ElMessage.error(msg);
  } finally {
    sendLoading.value = false;
  }
}

function loadCurrentSection() {
  if (section.value === "template") {
    loadTpl();
    loadTplChannelOptions();
  }
  if (section.value === "channel") loadCh();
  if (section.value === "log") loadLog();
}

onMounted(loadCurrentSection);
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
.extra-config-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.variable-config-list {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.variable-config-row {
  display: grid;
  grid-template-columns: 1fr 1fr auto auto;
  gap: 8px;
  align-items: center;
}
.send-variable-row {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 8px;
  align-items: center;
}
.extra-config-row {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 8px;
  align-items: center;
}
@media (max-width: 720px) {
  .variable-config-row,
  .send-variable-row,
  .extra-config-row {
    grid-template-columns: 1fr;
  }
}
</style>
