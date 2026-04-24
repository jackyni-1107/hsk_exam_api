package paper

import (
	"fmt"
	"sync"
	"time"

	exambo "exam/internal/model/bo/exam"
	"golang.org/x/sync/singleflight"
)

const adminPaperDetailCacheTTL = 2 * time.Minute

type cachedAdminPaperDetailEntry struct {
	data     *exambo.PaperDetailTree
	cachedAt time.Time
}

type cachedAdminPaperSectionEntry struct {
	data     *exambo.SectionDetailView
	cachedAt time.Time
}

var (
	adminPaperDetailCache  sync.Map // paperID -> *cachedAdminPaperDetailEntry
	adminPaperSectionCache sync.Map // "paperID:sectionID" -> *cachedAdminPaperSectionEntry
	adminPaperDetailSF     singleflight.Group
	adminPaperSectionSF    singleflight.Group
)

func loadAdminPaperDetailFromCache(paperID int64) *exambo.PaperDetailTree {
	if v, ok := adminPaperDetailCache.Load(paperID); ok {
		e := v.(*cachedAdminPaperDetailEntry)
		if time.Since(e.cachedAt) < adminPaperDetailCacheTTL {
			return e.data
		}
		adminPaperDetailCache.Delete(paperID)
	}
	return nil
}

func storeAdminPaperDetailToCache(paperID int64, data *exambo.PaperDetailTree) {
	if data == nil {
		return
	}
	adminPaperDetailCache.Store(paperID, &cachedAdminPaperDetailEntry{
		data:     data,
		cachedAt: time.Now(),
	})
}

func loadAdminPaperSectionFromCache(paperID, sectionID int64) *exambo.SectionDetailView {
	key := adminPaperSectionCacheKey(paperID, sectionID)
	if v, ok := adminPaperSectionCache.Load(key); ok {
		e := v.(*cachedAdminPaperSectionEntry)
		if time.Since(e.cachedAt) < adminPaperDetailCacheTTL {
			return e.data
		}
		adminPaperSectionCache.Delete(key)
	}
	return nil
}

func storeAdminPaperSectionToCache(paperID, sectionID int64, data *exambo.SectionDetailView) {
	if data == nil {
		return
	}
	adminPaperSectionCache.Store(adminPaperSectionCacheKey(paperID, sectionID), &cachedAdminPaperSectionEntry{
		data:     data,
		cachedAt: time.Now(),
	})
}

func invalidateAdminPaperDetailCacheByPaper(paperID int64) {
	adminPaperDetailCache.Delete(paperID)
	prefix := fmt.Sprintf("%d:", paperID)
	adminPaperSectionCache.Range(func(k, _ any) bool {
		if ks, ok := k.(string); ok && len(ks) >= len(prefix) && ks[:len(prefix)] == prefix {
			adminPaperSectionCache.Delete(ks)
		}
		return true
	})
}

func adminPaperSectionCacheKey(paperID, sectionID int64) string {
	return fmt.Sprintf("%d:%d", paperID, sectionID)
}
