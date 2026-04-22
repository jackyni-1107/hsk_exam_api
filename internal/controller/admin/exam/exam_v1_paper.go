package exam

import (
	"context"
	exambo "exam/internal/model/bo/exam"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"

	v1 "exam/api/admin/exam/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	papersvc "exam/internal/service/paper"
	"exam/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) PaperList(ctx context.Context, req *v1.PaperListReq) (res *v1.PaperListRes, err error) {
	list, total, err := papersvc.Paper().PaperList(ctx, req.Page, req.Size, req.Level)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.PaperListItem, 0, len(list))
	for _, p := range list {
		item := &v1.PaperListItem{
			Id:                      p.Id,
			MockExaminationPaperId:  p.MockExaminationPaperId,
			Level:                   p.Level,
			PaperId:                 p.PaperId,
			Title:                   p.Title,
			SourceBaseUrl:           p.SourceBaseUrl,
			AudioHlsPrefix:          p.AudioHlsPrefix,
			AudioHlsSegmentCount:    p.AudioHlsSegmentCount,
			AudioHlsSegmentPattern:  p.AudioHlsSegmentPattern,
			AudioHlsKeyObject:       p.AudioHlsKeyObject,
			AudioHlsIvHex:           p.AudioHlsIvHex,
			AudioHlsSegmentDuration: p.AudioHlsSegmentDuration,
		}
		item.CreateTime = utility.ToRFC3339UTC(p.CreateTime)
		out = append(out, item)
	}
	return &v1.PaperListRes{List: out, Total: total}, nil
}

func (c *ControllerV1) PaperDetail(ctx context.Context, req *v1.PaperDetailReq) (res *v1.PaperDetailRes, err error) {
	d, err := papersvc.Paper().PaperDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	var out v1.PaperDetailRes
	if err := gconv.Scan(d, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ControllerV1) PaperSectionDetail(ctx context.Context, req *v1.PaperSectionDetailReq) (res *v1.PaperSectionDetailRes, err error) {
	d, err := papersvc.Paper().PaperDetailSection(ctx, req.Id, req.SectionId)
	if err != nil {
		return nil, err
	}
	var section v1.PaperDetailSection
	if err := gconv.Scan(d, &section); err != nil {
		return nil, err
	}
	return &v1.PaperSectionDetailRes{Section: section}, nil
}

func (c *ControllerV1) PaperImport(ctx context.Context, req *v1.PaperImportReq) (res *v1.PaperImportRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	r, err := papersvc.Paper().ImportFromIndex(ctx, exambo.ImportParams{
		MockExaminationPaperId: req.MockExaminationPaperId,
		Title:                  req.Title,
		AudioHlsPrefix:         req.AudioHlsPrefix,
		ConflictMode:           req.ConflictMode,
		Creator:                creator,
	})
	if err != nil {
		return nil, err
	}
	return &v1.PaperImportRes{
		ExaminationPaperId:         r.ExaminationPaperID,
		Conflict:                   r.Conflict,
		ExistingExaminationPaperId: r.ExistingExaminationPaperID,
		SectionCount:               r.SectionCount,
		QuestionCount:              r.QuestionCount,
	}, nil
}

func (c *ControllerV1) PaperEdit(ctx context.Context, req *v1.PaperEditReq) (res *v1.PaperEditRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = papersvc.Paper().UpdatePaperMeta(ctx, req.ExamPaperId, exambo.PaperMetaAdminUpdate{
		Title:              strings.TrimSpace(req.Title),
		PrepareTitle:       strings.TrimSpace(req.PrepareTitle),
		PrepareInstruction: strings.TrimSpace(req.PrepareInstruction),
		PrepareAudioFile:   strings.TrimSpace(req.PrepareAudioFile),
		SourceBaseURL:      strings.TrimSpace(req.SourceBaseUrl),
		DurationSeconds:    req.DurationSeconds,
	}, updater)
	if err != nil {
		return nil, err
	}
	return &v1.PaperEditRes{}, nil
}

func (c *ControllerV1) PaperPurge(ctx context.Context, req *v1.PaperPurgeReq) (res *v1.PaperPurgeRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	ok, err := middleware.UserHasActiveRoleCode(ctx, d.UserId, consts.RoleCodeSuperAdmin)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, gerror.NewCode(consts.CodePermissionDenied)
	}
	if err := papersvc.Paper().PaperPurgePhysical(ctx, req.ExamPaperId, req.ConfirmText); err != nil {
		return nil, err
	}
	return &v1.PaperPurgeRes{}, nil
}

func (c *ControllerV1) PaperUpdate(ctx context.Context, req *v1.PaperUpdateReq) (res *v1.PaperUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = papersvc.Paper().UpdatePaperSettings(ctx, req.Id, exambo.PaperHlsExamAdminUpdate{
		AudioHlsPrefix:          req.AudioHlsPrefix,
		AudioHlsSegmentCount:    req.AudioHlsSegmentCount,
		AudioHlsSegmentPattern:  req.AudioHlsSegmentPattern,
		AudioHlsKeyObject:       req.AudioHlsKeyObject,
		AudioHlsIvHex:           req.AudioHlsIvHex,
		AudioHlsSegmentDuration: req.AudioHlsSegmentDuration,
	}, updater)
	if err != nil {
		return nil, err
	}
	return &v1.PaperUpdateRes{}, nil
}
