package paper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/sync/singleflight"

	"exam/internal/consts"
)

func paperForExamInitRedisKey(examPaperId int64) string {
	return fmt.Sprintf(consts.ExamPaperInitCacheKeyFmt, examPaperId)
}

func paperForExamLegacyRedisKey(examPaperId int64) string {
	return fmt.Sprintf(consts.ExamPaperLegacyCacheKeyFmt, examPaperId)
}

func paperForExamSectionRedisKey(examPaperId, sectionId int64) string {
	return fmt.Sprintf(consts.ExamPaperSectionCacheKeyFmt, examPaperId, sectionId)
}

// sectionTopicCacheVersion 用于在不修改 consts 的前提下对 section topic 缓存做版本隔离。
// 升级读/写路径语义（例如去掉读时的 YCT 剥离与 EOID 兜底扫描）后，bump 本常量可令历史脏缓存整体失效。
const sectionTopicCacheVersion = "v2"

func paperSectionTopicRedisKey(examPaperId, sectionId int64) string {
	return fmt.Sprintf(consts.ExamPaperSectionTopicCacheKeyFmt, examPaperId, sectionId) + ":" + sectionTopicCacheVersion
}

func paperPrepareRedisKey(mockPaperID int64) string {
	return fmt.Sprintf(consts.ExamPaperPrepareCacheKeyFmt, mockPaperID)
}

var paperForExamInitSF singleflight.Group
var paperForExamSectionSF singleflight.Group
var paperSectionTopicSF singleflight.Group
var paperPrepareSF singleflight.Group

// InvalidatePaperForExamCache 试卷树变更后删除考前相关缓存（初始化 + 各 section 详情 + section topic + 历史整卷 key）。
func InvalidatePaperForExamCache(ctx context.Context, examPaperId int64) {
	if examPaperId <= 0 {
		return
	}
	initKey := paperForExamInitRedisKey(examPaperId)
	legacyKey := paperForExamLegacyRedisKey(examPaperId)
	if _, err := g.Redis().Del(ctx, initKey, legacyKey); err != nil {
		g.Log().Warningf(ctx, "paper for-exam redis del init/legacy %s %s: %v", initKey, legacyKey, err)
	}
	invalidatePaperSectionExamCachesByPaper(ctx, examPaperId)
	invalidateByPattern(ctx, fmt.Sprintf(consts.ExamPaperSectionTopicCachePattern, examPaperId))
	invalidateSectionTopicMemCacheByPaper(examPaperId)
}

// InvalidatePaperPrepareCache 删除试卷准备阶段缓存。
func InvalidatePaperPrepareCache(ctx context.Context, mockPaperID int64) {
	if mockPaperID <= 0 {
		return
	}
	key := paperPrepareRedisKey(mockPaperID)
	if _, err := g.Redis().Del(ctx, key); err != nil {
		g.Log().Warningf(ctx, "paper prepare redis del %s: %v", key, err)
	}
}

// InvalidatePaperSectionForExamCache 删除单个 section 的考前详情缓存（精确 key）。
func InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId, sectionId int64) {
	if examPaperId <= 0 || sectionId <= 0 {
		return
	}
	key := paperForExamSectionRedisKey(examPaperId, sectionId)
	if _, err := g.Redis().Del(ctx, key); err != nil {
		g.Log().Warningf(ctx, "paper for-exam redis del section %s: %v", key, err)
	}
	topicKey := paperSectionTopicRedisKey(examPaperId, sectionId)
	if _, err := g.Redis().Del(ctx, topicKey); err != nil {
		g.Log().Warningf(ctx, "paper for-exam redis del section topic %s: %v", topicKey, err)
	}
	invalidateSectionTopicMemCache(examPaperId, sectionId)
}

func invalidatePaperSectionExamCachesByPaper(ctx context.Context, examPaperId int64) {
	invalidateByPattern(ctx, fmt.Sprintf(consts.ExamPaperSectionCachePattern, examPaperId))
}

func invalidateByPattern(ctx context.Context, pattern string) {
	var cursor int64
	for {
		v, err := g.Redis().Do(ctx, "SCAN", cursor, "MATCH", pattern, "COUNT", 100)
		if err != nil {
			g.Log().Warningf(ctx, "redis scan %s cursor=%d: %v", pattern, cursor, err)
			return
		}
		arr := v.Interfaces()
		if len(arr) < 2 {
			return
		}
		nextCursor := gvar.New(arr[0]).Int64()
		keys := gvar.New(arr[1]).Strings()
		if len(keys) > 0 {
			if _, err := g.Redis().Del(ctx, keys...); err != nil {
				g.Log().Warningf(ctx, "redis del keys pattern %s: %v", pattern, err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}

func redisGetPaperForExamInit(ctx context.Context, rkey string) *PaperDetailForExamInitTree {
	val, err := g.Redis().Get(ctx, rkey)
	if err != nil {
		g.Log().Warningf(ctx, "paper for-exam init redis get %s: %v", rkey, err)
		return nil
	}
	if val.IsEmpty() {
		return nil
	}
	var out PaperDetailForExamInitTree
	if err := json.Unmarshal([]byte(val.String()), &out); err != nil {
		g.Log().Warningf(ctx, "paper for-exam init redis unmarshal %s: %v", rkey, err)
		return nil
	}
	return &out
}

func redisSetPaperForExamInitJSON(ctx context.Context, rkey string, jsonStr string) {
	if err := g.Redis().SetEX(ctx, rkey, jsonStr, consts.PaperForExamCacheTTLSeconds); err != nil {
		g.Log().Warningf(ctx, "paper for-exam init redis setex %s: %v", rkey, err)
	}
}

func redisGetPaperSectionForExam(ctx context.Context, rkey string) *SectionDetailForExamView {
	val, err := g.Redis().Get(ctx, rkey)
	if err != nil {
		g.Log().Warningf(ctx, "paper for-exam section redis get %s: %v", rkey, err)
		return nil
	}
	if val.IsEmpty() {
		return nil
	}
	var out SectionDetailForExamView
	if err := json.Unmarshal([]byte(val.String()), &out); err != nil {
		g.Log().Warningf(ctx, "paper for-exam section redis unmarshal %s: %v", rkey, err)
		return nil
	}
	return &out
}

func redisSetPaperSectionForExamJSON(ctx context.Context, rkey string, jsonStr string) {
	if err := g.Redis().SetEX(ctx, rkey, jsonStr, consts.PaperForExamCacheTTLSeconds); err != nil {
		g.Log().Warningf(ctx, "paper for-exam section redis setex %s: %v", rkey, err)
	}
}

func redisGetSectionTopicJSON(ctx context.Context, rkey string) string {
	val, err := g.Redis().Get(ctx, rkey)
	if err != nil {
		g.Log().Warningf(ctx, "section topic redis get %s: %v", rkey, err)
		return ""
	}
	return val.String()
}

func redisSetSectionTopicJSON(ctx context.Context, rkey string, jsonStr string) {
	if err := g.Redis().SetEX(ctx, rkey, jsonStr, consts.PaperForExamCacheTTLSeconds); err != nil {
		g.Log().Warningf(ctx, "section topic redis setex %s: %v", rkey, err)
	}
}

func redisGetPrepareJSON(ctx context.Context, rkey string) string {
	val, err := g.Redis().Get(ctx, rkey)
	if err != nil {
		g.Log().Warningf(ctx, "paper prepare redis get %s: %v", rkey, err)
		return ""
	}
	return val.String()
}

func redisSetPrepareJSON(ctx context.Context, rkey string, jsonStr string) {
	if err := g.Redis().SetEX(ctx, rkey, jsonStr, consts.PaperForExamCacheTTLSeconds); err != nil {
		g.Log().Warningf(ctx, "paper prepare redis setex %s: %v", rkey, err)
	}
}

func hlsPlayRedisKey(ticket string) string {
	return consts.HlsPlayKeyPrefix + ticket
}

func redisSetHlsPlayTicket(ctx context.Context, redisKey string, payloadJSON string, ttlSeconds int64) error {
	return g.Redis().SetEX(ctx, redisKey, payloadJSON, ttlSeconds)
}

func redisGetHlsPlayTicket(ctx context.Context, redisKey string) (*gvar.Var, error) {
	return g.Redis().Get(ctx, redisKey)
}
