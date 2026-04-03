// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mock

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MockExaminationPaper is the golang structure for table mock_examination_paper.
type MockExaminationPaper struct {
	Id                   int64       `json:"id"                     orm:"id"                     description:"id 主键"`                // id 主键
	LevelId              int64       `json:"level_id"               orm:"level_id"               description:"HSK等级"`                // HSK等级
	Name                 string      `json:"name"                   orm:"name"                   description:"试卷名称"`                 // 试卷名称
	NameTrans            string      `json:"name_trans"             orm:"name_trans"             description:"试卷名称国际化"`              // 试卷名称国际化
	ScoreFull            int         `json:"score_full"             orm:"score_full"             description:"满分"`                   // 满分
	Seq                  int         `json:"seq"                    orm:"seq"                    description:"试卷顺序"`                 // 试卷顺序
	ExplainAudio         string      `json:"explain_audio"          orm:"explain_audio"          description:"说明音频"`                 // 说明音频
	TimeFull             int         `json:"time_full"              orm:"time_full"              description:"考试时长"`                 // 考试时长
	ListenReviewDuration int         `json:"listen_review_duration" orm:"listen_review_duration" description:"听力结束后回顾时间"`            // 听力结束后回顾时间
	TimeSheet            int         `json:"time_sheet"             orm:"time_sheet"             description:"答题卡时间"`                // 答题卡时间
	IconUrl              string      `json:"icon_url"               orm:"icon_url"               description:"图标路径"`                 // 图标路径
	ResourceUrl          string      `json:"resource_url"           orm:"resource_url"           description:"资源包路径"`                // 资源包路径
	Version              int         `json:"version"                orm:"version"                description:"资源版本号"`                // 资源版本号
	DeleteFlag           int         `json:"delete_flag"            orm:"delete_flag"            description:"是否删除 0 未删除 1 已删除"`     // 是否删除 0 未删除 1 已删除
	MockType             int         `json:"mock_type"              orm:"mock_type"              description:"1 hsk 2 hskk 3 yct"`   // 1 hsk 2 hskk 3 yct
	ProductBaseId        int         `json:"product_base_id"        orm:"product_base_id"        description:"学习需要付的元商品"`            // 学习需要付的元商品
	Credit               int         `json:"credit"                 orm:"credit"                 description:"所需元商品数量"`              // 所需元商品数量
	BuyProductId         int         `json:"buy_product_id"         orm:"buy_product_id"         description:"当付费券不足时, 需要购买此商品"`     // 当付费券不足时, 需要购买此商品
	CreateTime           *gtime.Time `json:"create_time"            orm:"create_time"            description:"创建时间"`                 // 创建时间
	UpdateTime           *gtime.Time `json:"update_time"            orm:"update_time"            description:"更新时间"`                 // 更新时间
	Status               int         `json:"status"                 orm:"status"                 description:"试卷状态 0 未发布 1 发布 2 下架"` // 试卷状态 0 未发布 1 发布 2 下架
	MemberResource       int         `json:"member_resource"        orm:"member_resource"        description:"0 非会员资源1 会员资源"`        // 0 非会员资源1 会员资源
	PaperType            int         `json:"paper_type"             orm:"paper_type"             description:"1 模拟考试卷 2 在线练习试卷"`     // 1 模拟考试卷 2 在线练习试卷
}
