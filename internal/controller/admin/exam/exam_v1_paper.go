package exam

import (
	"context"
	exambo "exam/internal/model/bo/exam"

	"github.com/gogf/gf/v2/util/gconv"

	v1 "exam/api/admin/exam/v1"
	"exam/internal/middleware"
	"exam/internal/service/exam"
	"exam/internal/utility"
	"exam/internal/utility/exampaper"
)

func (c *ControllerV1) PaperList(ctx context.Context, req *v1.PaperListReq) (res *v1.PaperListRes, err error) {
	list, total, err := exam.Exam().PaperList(ctx, req.Page, req.Size, req.Level)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.PaperListItem, 0, len(list))
	for _, p := range list {
		item := &v1.PaperListItem{
			Id:                      p.MockExaminationPaperId,
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
	paper, err := exampaper.ByMockID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	d, err := exam.Exam().PaperDetail(ctx, paper.Id)
	if err != nil {
		return nil, err
	}
	var out v1.PaperDetailRes
	if err := gconv.Scan(d, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ControllerV1) PaperImport(ctx context.Context, req *v1.PaperImportReq) (res *v1.PaperImportRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	r, err := exam.Exam().ImportFromIndex(ctx, exambo.ImportParams{
		MockExaminationPaperId: req.MockExaminationPaperId,
		IndexURL:               req.IndexUrl,
		IndexJSON:              req.IndexJson,
		Level:                  req.Level,
		PaperID:                req.PaperId,
		SourceBaseURL:          req.SourceBaseUrl,
		AudioHlsPrefix:         req.AudioHlsPrefix,
		ConflictMode:           req.ConflictMode,
		NewPaperID:             req.NewPaperId,
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

func (c *ControllerV1) PaperUpdate(ctx context.Context, req *v1.PaperUpdateReq) (res *v1.PaperUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	paper, err := exampaper.ByMockID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	err = exam.Exam().UpdatePaperSettings(ctx, paper.Id, exambo.PaperHlsExamAdminUpdate{
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
