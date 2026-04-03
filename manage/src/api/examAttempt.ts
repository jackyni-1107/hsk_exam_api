import request from './request'

export interface AttemptListItem {
  id: number
  member_id: number
  username: string
  nickname: string
  examination_paper_id: number
  exam_batch_id: number
  mock_level_id: number
  paper_title: string
  paper_level: string
  remote_paper_id: string
  status: number
  objective_score: number
  subjective_score: number
  total_score: number
  has_subjective: number
  started_at: string
  submitted_at: string
  ended_at: string
  create_time: string
}

export interface AttemptDetailAnswer {
  question_id: number
  question_no: number
  stem_text: string
  is_example: number
  is_subjective: number
  score: number
  answer_json: string
  awarded_score: number | null
  objective_correct: boolean | null
  section_id: number
  section_title: string
  options?: AttemptDetailOption[]
}

export interface AttemptDetailOption {
  id: number
  flag: string
  sort_order: number
  is_correct: number
  option_type: string
  content: string
}

export interface AttemptDetail {
  attempt: {
    id: number
    member_id: number
    examination_paper_id: number
    status: number
    duration_seconds: number
    objective_score: number
    subjective_score: number
    total_score: number
    has_subjective: number
    started_at: string
    deadline_at: string
    submitted_at: string
    ended_at: string
    create_time: string
  }
  user: { id: number; username: string; nickname: string }
  paper: { /** mock 卷 id */ id: number; level: string; paper_id: string; title: string }
  answers: AttemptDetailAnswer[]
}

export function getAttemptList(params: {
  page?: number
  size?: number
  level?: string
  examination_paper_id?: number
  exam_batch_id?: number
  status?: number
  username?: string
}) {
  return request.get<any, { data: { list: AttemptListItem[]; total: number } }>('/admin/exam/attempt/list', {
    params,
  })
}

export function getAttemptDetail(id: number) {
  return request.get<any, { data: AttemptDetail }>(`/admin/exam/attempt/${id}`)
}

export function saveAttemptSubjectiveScores(
  id: number,
  items: { question_id: number; score: number }[],
) {
  return request.put<any, { data: { subjective_score: number; total_score: number } }>(
    `/admin/exam/attempt/${id}/subjective-scores`,
    { items },
  )
}
