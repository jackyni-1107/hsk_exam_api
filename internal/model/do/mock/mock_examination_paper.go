// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mock

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MockExaminationPaper is the golang structure of table mock_examination_paper for DAO operations like Where/Data.
type MockExaminationPaper struct {
	g.Meta               `orm:"table:mock_examination_paper, do:true"`
	Id                   any         // id 主键
	LevelId              any         // HSK等级
	Name                 any         // 试卷名称
	NameTrans            any         // 试卷名称国际化
	ScoreFull            any         // 满分
	Seq                  any         // 试卷顺序
	ExplainAudio         any         // 说明音频
	TimeFull             any         // 考试时长
	ListenReviewDuration any         // 听力结束后回顾时间
	TimeSheet            any         // 答题卡时间
	IconUrl              any         // 图标路径
	ResourceUrl          any         // 资源包路径
	Version              any         // 资源版本号
	DeleteFlag           any         // 是否删除 0 未删除 1 已删除
	MockType             any         // 1 hsk 2 hskk 3 yct
	ProductBaseId        any         // 学习需要付的元商品
	Credit               any         // 所需元商品数量
	BuyProductId         any         // 当付费券不足时, 需要购买此商品
	CreateTime           *gtime.Time // 创建时间
	UpdateTime           *gtime.Time // 更新时间
	Status               any         // 试卷状态 0 未发布 1 发布 2 下架
	MemberResource       any         // 0 非会员资源1 会员资源
	PaperType            any         // 1 模拟考试卷 2 在线练习试卷
}
