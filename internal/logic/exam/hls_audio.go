package exam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"

	"exam/internal/config"
	"exam/internal/consts"
	"exam/internal/dao"
	examentity "exam/internal/model/entity/exam"
	"exam/internal/utility/storage"
)

type hlsPlayTicketPayload struct {
	UserID         int64 `json:"user_id"`
	AttemptID      int64 `json:"attempt_id"`
	ExamQuestionID int64 `json:"exam_question_id"`
	ExamPaperID    int64 `json:"exam_paper_id"`
	MaxSegSnap     int   `json:"max_segment_snap"`
}

type paperHLSConfig struct {
	Id                      int64   `json:"id"`
	AudioHlsPrefix          string  `json:"audio_hls_prefix"`
	AudioHlsSegmentCount    int     `json:"audio_hls_segment_count"`
	AudioHlsSegmentPattern  string  `json:"audio_hls_segment_pattern"`
	AudioHlsKeyObject       string  `json:"audio_hls_key_object"`
	AudioHlsIvHex           string  `json:"audio_hls_iv_hex"`
	AudioHlsSegmentDuration float64 `json:"audio_hls_segment_duration"`
}

func (s *sExam) hlsTicketTTL() time.Duration {
	sec := config.Config.Exam.AudioHlsTicketTTLSeconds
	if sec <= 0 {
		sec = 120
	}
	return time.Duration(sec) * time.Second
}

func (s *sExam) hlsPresignTTL() time.Duration {
	sec := config.Config.Exam.AudioHlsPresignTTLSeconds
	if sec <= 0 {
		sec = 300
	}
	tick := config.Config.Exam.AudioHlsTicketTTLSeconds
	if tick <= 0 {
		tick = 120
	}
	if sec < tick {
		sec = tick
	}
	return time.Duration(sec) * time.Second
}

func storageSupportsPresign(stType string) bool {
	return stType == "oss" || stType == "s3" || stType == "minio"
}

func joinOSSKeys(parts ...string) string {
	var out string
	for _, part := range parts {
		p := strings.Trim(part, "/")
		if p == "" {
			continue
		}
		if out == "" {
			out = p
			continue
		}
		out = out + "/" + p
	}
	return out
}

func minInt3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

// loadPaperHLS 按 mock_examination_paper.id 查卷（与客户端 paper.id 一致）。
func (s *sExam) loadPaperHLS(ctx context.Context, mockExaminationPaperID int64) (*paperHLSConfig, error) {
	c := dao.ExamPaper.Columns()
	return s.pickExamPaperHLS(ctx, c.MockExaminationPaperId, mockExaminationPaperID)
}

// loadPaperHLSByExamPaperRowID 按 exam_paper.id 查卷（HLS 票据内存的是行主键）。
func (s *sExam) loadPaperHLSByExamPaperRowID(ctx context.Context, examPaperID int64) (*paperHLSConfig, error) {
	c := dao.ExamPaper.Columns()
	return s.pickExamPaperHLS(ctx, c.Id, examPaperID)
}

func (s *sExam) pickExamPaperHLS(ctx context.Context, column string, value int64) (*paperHLSConfig, error) {
	var cfg paperHLSConfig
	err := dao.ExamPaper.Ctx(ctx).
		Fields("id,audio_hls_prefix,audio_hls_segment_count,audio_hls_segment_pattern,audio_hls_key_object,audio_hls_iv_hex,audio_hls_segment_duration").
		Where(column, value).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	if strings.TrimSpace(cfg.AudioHlsPrefix) == "" || cfg.AudioHlsSegmentCount <= 0 {
		return nil, gerror.NewCode(consts.CodeExamAudioHlsNotConfigured)
	}
	return &cfg, nil
}

func (s *sExam) assertAttemptInProgress(ctx context.Context, userID, attemptID int64) (*examentity.ExamAttempt, error) {
	return assertAttemptInProgressByUser(ctx, attemptID, userID)
}

func (s *sExam) loadQuestionHLS(ctx context.Context, paperID, questionID int64) (*examentity.ExamQuestion, error) {
	var q examentity.ExamQuestion
	err := dao.ExamQuestion.Ctx(ctx).
		Where("id", questionID).
		Where("exam_paper_id", paperID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&q)
	if err != nil {
		return nil, err
	}
	if q.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	if strings.TrimSpace(q.AudioHlsPrefix) == "" || q.AudioHlsSegmentCount <= 0 {
		return nil, gerror.NewCode(consts.CodeExamAudioHlsNotConfigured)
	}
	return &q, nil
}

// IssueAudioHlsPlay 签发短期播放票据，返回相对 play_url（以 / 开头）。
func (s *sExam) IssueAudioHlsPlay(ctx context.Context, userID, attemptID, questionID int64) (playURL string, expiresAt string, err error) {
	att, err := s.assertAttemptInProgress(ctx, userID, attemptID)
	if err != nil {
		return "", "", err
	}
	q, err := s.loadQuestionHLS(ctx, att.ExamPaperId, questionID)
	if err != nil {
		return "", "", err
	}
	stCfg, _ := storage.GetActiveConfig(ctx)
	if !storageSupportsPresign(stCfg.Type) {
		return "", "", gerror.NewCode(consts.CodeExamStoragePresignOnly)
	}
	maxDB := q.AudioHlsSegmentCount
	ticket := uuid.NewString()
	payload := hlsPlayTicketPayload{
		UserID:         userID,
		AttemptID:      attemptID,
		ExamQuestionID: questionID,
		ExamPaperID:    att.ExamPaperId,
		MaxSegSnap:     maxDB,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}
	ttl := s.hlsTicketTTL()
	key := consts.HlsPlayKeyPrefix + ticket
	if err := g.Redis().SetEX(ctx, key, string(b), int64(ttl.Seconds())); err != nil {
		return "", "", err
	}
	playURL = "/api/client/exam/media/hls/" + ticket + ".m3u8"
	expiresAt = time.Now().UTC().Add(ttl).Format(time.RFC3339)
	return playURL, expiresAt, nil
}

// IssuePaperHlsPlay 基于试卷级 HLS 配置签发短期播放票据，返回相对 play_url（以 / 开头）。
func (s *sExam) IssuePaperHlsPlay(ctx context.Context, userID, paperID int64) (playURL string, expiresAt string, err error) {
	if paperID <= 0 {
		return "", "", gerror.NewCode(consts.CodeInvalidParams)
	}
	stCfg, _ := storage.GetActiveConfig(ctx)
	if !storageSupportsPresign(stCfg.Type) {
		return "", "", gerror.NewCode(consts.CodeExamStoragePresignOnly)
	}
	paperCfg, err := s.loadPaperHLS(ctx, paperID)
	if err != nil {
		return "", "", err
	}
	maxDB := paperCfg.AudioHlsSegmentCount - 1
	if maxDB < 0 {
		return "", "", gerror.NewCode(consts.CodeExamAudioHlsNotConfigured)
	}
	ticket := uuid.NewString()
	payload := hlsPlayTicketPayload{
		UserID:         userID,
		AttemptID:      0,
		ExamQuestionID: 0,
		ExamPaperID:    paperCfg.Id,
		MaxSegSnap:     maxDB,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}
	ttl := s.hlsTicketTTL()
	key := consts.HlsPlayKeyPrefix + ticket
	if err := g.Redis().SetEX(ctx, key, string(b), int64(ttl.Seconds())); err != nil {
		return "", "", err
	}
	playURL = "/api/client/exam/media/hls/" + ticket + ".m3u8"
	expiresAt = time.Now().UTC().Add(ttl).Format(time.RFC3339)
	return playURL, expiresAt, nil
}

// BuildHlsM3U8Playlist 校验 Redis 票据并生成 m3u8（内嵌 presigned URL）。
func (s *sExam) BuildHlsM3U8Playlist(ctx context.Context, ticket string) ([]byte, error) {
	ticket = strings.TrimSuffix(strings.TrimSpace(ticket), ".m3u8")
	if ticket == "" {
		return nil, gerror.NewCode(consts.CodeExamHlsTicketInvalid)
	}
	key := consts.HlsPlayKeyPrefix + ticket
	val, err := g.Redis().Get(ctx, key)
	if err != nil || val.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeExamHlsTicketInvalid)
	}
	var payload hlsPlayTicketPayload
	if err := json.Unmarshal([]byte(val.String()), &payload); err != nil {
		return nil, gerror.NewCode(consts.CodeExamHlsTicketInvalid)
	}
	stCfg, _ := storage.GetActiveConfig(ctx)
	if !storageSupportsPresign(stCfg.Type) {
		return nil, gerror.NewCode(consts.CodeExamStoragePresignOnly)
	}
	if payload.ExamPaperID <= 0 {
		return nil, gerror.NewCode(consts.CodeExamHlsTicketInvalid)
	}
	audioPrefix := ""
	hlsPrefix := ""
	keyObject := ""
	ivHex := ""
	pattern := "%05d.ts"
	extDur := 10.0
	effectiveMax := -1
	if payload.ExamQuestionID > 0 {
		var att examentity.ExamAttempt
		err = dao.ExamAttempt.Ctx(ctx).
			Where("id", payload.AttemptID).
			Where("member_id", payload.UserID).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&att)
		if err != nil {
			return nil, err
		}
		if att.Id == 0 || att.ExamPaperId != payload.ExamPaperID {
			return nil, gerror.NewCode(consts.CodeExamHlsTicketInvalid)
		}
		if att.Status != consts.ExamAttemptInProgress {
			return nil, gerror.NewCode(consts.CodeExamHlsTicketInvalid)
		}
		q, err := s.loadQuestionHLS(ctx, att.ExamPaperId, payload.ExamQuestionID)
		if err != nil {
			return nil, err
		}
		var paper examentity.ExamPaper
		if err := dao.ExamPaper.Ctx(ctx).
			Where("id", att.ExamPaperId).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&paper); err != nil {
			return nil, err
		}
		if paper.Id == 0 {
			return nil, gerror.NewCode(consts.CodeExamHlsTicketInvalid)
		}
		audioPrefix = strings.TrimSpace(paper.AudioHlsPrefix)
		hlsPrefix = strings.TrimSpace(q.AudioHlsPrefix)
		keyObject = strings.TrimSpace(q.AudioHlsKeyObject)
		ivHex = strings.TrimSpace(q.AudioHlsIvHex)
		if v := strings.TrimSpace(q.AudioHlsSegmentPattern); v != "" {
			pattern = v
		}
		if q.AudioHlsSegmentDuration > 0 {
			extDur = q.AudioHlsSegmentDuration
		}
		lastIdx := q.AudioHlsSegmentCount - 1
		effectiveMax = minInt3(payload.MaxSegSnap, lastIdx, lastIdx)
	} else {
		paperCfg, err := s.loadPaperHLSByExamPaperRowID(ctx, payload.ExamPaperID)
		if err != nil {
			return nil, err
		}
		hlsPrefix = strings.TrimSpace(paperCfg.AudioHlsPrefix)
		keyObject = strings.TrimSpace(paperCfg.AudioHlsKeyObject)
		ivHex = strings.TrimSpace(paperCfg.AudioHlsIvHex)
		if v := strings.TrimSpace(paperCfg.AudioHlsSegmentPattern); v != "" {
			pattern = v
		}
		if paperCfg.AudioHlsSegmentDuration > 0 {
			extDur = paperCfg.AudioHlsSegmentDuration
		}
		lastIdx := paperCfg.AudioHlsSegmentCount - 1
		effectiveMax = minInt3(payload.MaxSegSnap, lastIdx, lastIdx)
	}
	if effectiveMax < 0 {
		return nil, gerror.NewCode(consts.CodeExamHlsTicketInvalid)
	}
	// RFC 8216：TARGETDURATION 为最大分片时长向上取整；避免浮点 10.000 变成 9.999… 导致偏小
	targetDur := int(math.Ceil(extDur + 1e-6))
	if targetDur < 1 {
		targetDur = 1
	}
	expire := s.hlsPresignTTL()
	adapter := storage.NewAdapter()
	bucket := stCfg.Bucket

	var buf bytes.Buffer
	buf.WriteString("#EXTM3U\n")
	buf.WriteString("#EXT-X-VERSION:3\n")
	fmt.Fprintf(&buf, "#EXT-X-TARGETDURATION:%d\n", targetDur)

	if keyObject != "" {
		keyPath := joinOSSKeys(audioPrefix, hlsPrefix, keyObject)
		keyURL, err := adapter.PresignGet(ctx, bucket, keyPath, expire)
		if err != nil {
			return nil, gerror.Wrapf(err, "presign hls key bucket=%q object=%q", bucket, keyPath)
		}
		if keyURL == "" {
			return nil, gerror.Newf("presign hls key empty url bucket=%q object=%q", bucket, keyPath)
		}
		if ivHex != "" {
			fmt.Fprintf(&buf, `#EXT-X-KEY:METHOD=AES-128,URI="%s",IV=0x%s`+"\n", keyURL, strings.ToUpper(ivHex))
		} else {
			fmt.Fprintf(&buf, `#EXT-X-KEY:METHOD=AES-128,URI="%s"`+"\n", keyURL)
		}
	}

	for i := 0; i <= effectiveMax; i++ {
		segName := fmt.Sprintf(pattern, i)
		segPath := joinOSSKeys(audioPrefix, hlsPrefix, segName)
		segURL, err := adapter.PresignGet(ctx, bucket, segPath, expire)
		if err != nil {
			return nil, gerror.Wrapf(err, "presign hls segment bucket=%q object=%q index=%d", bucket, segPath, i)
		}
		if segURL == "" {
			return nil, gerror.Newf("presign hls segment empty url bucket=%q object=%q index=%d", bucket, segPath, i)
		}
		fmt.Fprintf(&buf, "#EXTINF:%.3f,\n", extDur)
		buf.WriteString(segURL)
		buf.WriteByte('\n')
	}
	buf.WriteString("#EXT-X-ENDLIST\n")
	return buf.Bytes(), nil
}
