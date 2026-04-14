package consts

import (
	"github.com/gogf/gf/v2/errors/gcode"
)

// 业务错误码（>=1000，避免与框架保留码冲突）
// message 使用 i18n key，由前端或 i18n 中间件解析
// 声明顺序：通用错误 → 系统管理类错误 → 考试/Mock 等业务域错误（与数值段无必然对应，仅便于阅读）
var (
	// ---------- 通用（认证、账号、权限、限流等，10001-10020） ----------
	CodeInvalidCredentials     = gcode.New(10001, "err.invalid_credentials", nil)
	CodeUserDisabled           = gcode.New(10002, "err.user_disabled", nil)
	CodePermissionDenied       = gcode.New(10003, "err.permission_denied", nil)
	CodeTokenInvalid           = gcode.New(10004, "err.token_invalid", nil)
	CodeTokenRequired          = gcode.New(10005, "err.token_required", nil)
	CodeUserNotFound           = gcode.New(10006, "err.user_not_found", nil)
	CodeUserExists             = gcode.New(10007, "err.user_exists", nil)
	CodeRoleNotFound           = gcode.New(10008, "err.role_not_found", nil)
	CodeMenuNotFound           = gcode.New(10009, "err.menu_not_found", nil)
	CodeInvalidParams          = gcode.New(10010, "err.invalid_params", nil)
	CodeLoginFailed            = gcode.New(10011, "err.login_failed", nil) // 登录失败（内部错误，不暴露详情）
	CodeAccountLocked          = gcode.New(10012, "err.account_locked", nil)
	CodeTooManyRequests        = gcode.New(10013, "err.too_many_requests", nil)
	CodeCaptchaRequired        = gcode.New(10014, "err.captcha_required", nil)
	CodeCaptchaInvalid         = gcode.New(10015, "err.captcha_invalid", nil)
	CodePasswordWeak           = gcode.New(10016, "err.password_weak", nil)
	CodePasswordExpired        = gcode.New(10017, "err.password_expired", nil)
	CodePasswordReuse          = gcode.New(10018, "err.password_reuse", nil)
	CodeResourceNotFound       = gcode.New(10019, "err.not_found", nil)
	CodeCannotDeleteSuperAdmin = gcode.New(10020, "err.cannot_delete_super_admin", nil)
	CodeDataNotFound           = gcode.New(10021, "err.data_not_found", nil)

	// ---------- 系统（配置、字典、文件、任务、通知、角色与组织，12001-16003） ----------
	// 配置与字典（12001-12004）
	CodeConfigExists             = gcode.New(12001, "err.config_exists", nil)
	CodeConfigNotFound           = gcode.New(12002, "err.config_not_found", nil)
	CodeDictTypeExists           = gcode.New(12003, "err.dict_type_exists", nil)
	CodeCannotDeleteActiveConfig = gcode.New(12004, "err.cannot_delete_active_config", nil)
	// 文件与上传（13001-13008）
	CodeFileRequired       = gcode.New(13001, "err.file_required", nil)
	CodeFileTypeNotAllowed = gcode.New(13002, "err.file_type_not_allowed", nil)
	CodeUploadNotFound     = gcode.New(13003, "err.upload_not_found", nil)
	CodeUploadCompleted    = gcode.New(13004, "err.upload_completed", nil)
	CodeChunksIncomplete   = gcode.New(13005, "err.chunks_incomplete", nil)
	CodeChunkMergeNotImpl  = gcode.New(13006, "err.chunk_merge_not_impl", nil)
	CodeFileNotFound       = gcode.New(13007, "err.file_not_found", nil)
	CodeFileNotPrivate     = gcode.New(13008, "err.file_not_private", nil)
	// 任务调度（14001-14005）
	CodeTaskCodeExists       = gcode.New(14001, "err.task_code_exists", nil)
	CodeCronExprRequired     = gcode.New(14002, "err.cron_expr_required", nil)
	CodeDelaySecondsRequired = gcode.New(14003, "err.delay_seconds_required", nil)
	CodeTaskNotFound         = gcode.New(14004, "err.task_not_found", nil)
	CodeTaskDisabled         = gcode.New(14005, "err.task_disabled", nil)
	// 通知模板与渠道（15001-15005）
	CodeTemplateExists            = gcode.New(15001, "err.template_exists", nil)
	CodeTemplateNotFound          = gcode.New(15002, "err.template_not_found", nil)
	CodeEmailMustUseSmtp          = gcode.New(15003, "err.email_must_use_smtp", nil)
	CodeSmsMustUseAliyunOrTencent = gcode.New(15004, "err.sms_must_use_aliyun_or_tencent", nil)
	CodeUnsupportedChannel        = gcode.New(15005, "err.unsupported_channel", nil)
	// 角色与权限扩展、学员（16001+）
	CodeRoleCodeExists = gcode.New(16001, "err.role_code_exists", nil)
	CodeRoleExists     = gcode.New(16002, "err.role_exists", nil)
	CodeMemberExists   = gcode.New(16003, "err.member_exists", nil)

	// ---------- 业务：考试与 Mock（11001-11203） ----------
	// 考试作答与播放（11001-11009）
	CodeExamAttemptNotFound       = gcode.New(11001, "err.exam_attempt_not_found", nil)
	CodeExamTimeExpired           = gcode.New(11002, "err.exam_time_expired", nil)
	CodeExamAlreadySubmitted      = gcode.New(11003, "err.exam_already_submitted", nil)
	CodeExamAnswerVersionConflict = gcode.New(11004, "err.exam_answer_version_conflict", nil)
	CodeExamTestHelperDisabled    = gcode.New(11005, "err.exam_test_helper_disabled", nil)
	CodeExamNotStarted            = gcode.New(11006, "err.exam_not_started", nil)
	CodeExamAudioHlsNotConfigured = gcode.New(11007, "err.exam_audio_hls_not_configured", nil)
	CodeExamHlsTicketInvalid      = gcode.New(11008, "err.exam_hls_ticket_invalid", nil)
	CodeExamStoragePresignOnly    = gcode.New(11009, "err.exam_storage_presign_only", nil)
	// 考试内容、导入、批次与模拟卷关联（11101-11123）
	CodeExamPaperNotFound              = gcode.New(11101, "err.exam_paper_not_found", nil)
	CodeExamSectionNotFound            = gcode.New(11102, "err.exam_section_not_found", nil)
	CodeExamSectionTopicEmpty          = gcode.New(11103, "err.exam_section_topic_empty", nil)
	CodeExamNewPaperIdRequired         = gcode.New(11104, "err.exam_new_paper_id_required", nil)
	CodeExamNewPaperIdExists           = gcode.New(11105, "err.exam_new_paper_id_exists", nil)
	CodeExamConflictModeInvalid        = gcode.New(11106, "err.exam_conflict_mode_invalid", nil)
	CodeExamIndexMetaRequired          = gcode.New(11107, "err.exam_index_meta_required", nil)
	CodeExamSourceBaseRequired         = gcode.New(11108, "err.exam_source_base_required", nil)
	CodeExamIndexUrlRequired           = gcode.New(11109, "err.exam_index_url_required", nil)
	CodeExamIndexUrlInvalid            = gcode.New(11110, "err.exam_index_url_invalid", nil)
	CodeExamAttemptNotEnded            = gcode.New(11111, "err.exam_attempt_not_ended", nil)
	CodeExamAttemptNoSubjective        = gcode.New(11112, "err.exam_attempt_no_subjective", nil)
	CodeMockExaminationPaperIdRequired = gcode.New(11113, "err.mock_examination_paper_id_required", nil)
	// Mock 卷存在但尚未向 exam_paper 导入 index/题目树（与 CodeMockExamPaperNotFound 区分）
	CodeExamPaperNotImported      = gcode.New(11114, "err.exam_paper_not_imported", nil)
	CodeExamBatchNotFound         = gcode.New(11115, "err.exam_batch_not_found", nil)
	CodeExamBatchTimeInvalid      = gcode.New(11116, "err.exam_batch_time_invalid", nil)
	CodeExamAttemptUseBatchApi    = gcode.New(11117, "err.exam_attempt_use_batch_api", nil)
	CodeExamBatchWindowNotOpen    = gcode.New(11119, "err.exam_batch_window_not_open", nil)
	CodeExamBatchMemberNotFound   = gcode.New(11120, "err.exam_batch_member_not_found", nil)
	CodeExamAttemptExistsForBatch = gcode.New(11121, "err.exam_attempt_exists_for_batch", nil)
	CodeExamBatchPaperNotInBatch  = gcode.New(11122, "err.exam_batch_paper_not_in_batch", nil)
	CodeExamBatchPaperHasMembers  = gcode.New(11123, "err.exam_batch_paper_has_members", nil)
	// Mock 模拟卷（11201-11299）
	CodeMockExamPaperNotFound   = gcode.New(11201, "err.mock_exam_paper_not_found", nil)
	CodeMockLevelNotFound       = gcode.New(11202, "err.mock_level_not_found", nil)
	CodeMockDataDeleteForbidden = gcode.New(11203, "err.mock_data_delete_forbidden", nil)
)
