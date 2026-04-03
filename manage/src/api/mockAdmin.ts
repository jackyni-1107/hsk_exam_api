import request from './request'

export interface MockLevelItem {
  id: number
  level_id: number
  level_type: number
  type_name: string
  level_name: string
  app_level_name: string
  exam_show_status: number
  homework_show_status: number
}

export function getMockLevelsList() {
  return request.get<any, { data: { list: MockLevelItem[] } }>('/admin/mock/levels/list')
}

export interface MockExaminationPaperItem {
  id: number
  level_id: number
  name: string
  score_full: number
  time_full: number
  status: number
  paper_type: number
  mock_type: number
  /** 是否已在 exam 域导入（exam_paper） */
  imported: boolean
}

export type MockPaperImportStatusFilter = '' | 'imported' | 'not_imported'

export function getMockExaminationPapers(params?: {
  level_id?: number
  import_status?: MockPaperImportStatusFilter
}) {
  return request.get<any, { data: { list: MockExaminationPaperItem[] } }>(
    '/admin/mock/examination-paper/list',
    { params },
  )
}

export function getMockExaminationPaperDetail(id: number) {
  return request.get<any, { data: { paper: MockExaminationPaperItem } }>(
    `/admin/mock/examination-paper/${id}`,
  )
}
