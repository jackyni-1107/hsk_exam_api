import request from './request'

export interface ExamPaperItem {
  /** mock_examination_paper.id */
  id: number
  level: string
  paper_id: string
  title: string
  source_base_url: string
  audio_hls_prefix: string
  audio_hls_segment_count: number
  audio_hls_segment_pattern: string
  audio_hls_key_object: string
  audio_hls_iv_hex: string
  audio_hls_segment_duration: number
  create_time: string
}

export function getExamPaperList(params: { level?: string; page?: number; size?: number }) {
  return request.get<any, { data: { list: ExamPaperItem[]; total: number } }>('/admin/exam/paper/list', {
    params,
  })
}

export interface ExamPaperDetail {
  paper: {
    /** mock_examination_paper.id */
    id: number
    level: string
    paper_id: string
    title: string
    prepare_title: string
    prepare_instruction: string
    prepare_audio_file: string
    source_base_url: string
    audio_hls_prefix: string
    audio_hls_segment_count: number
    audio_hls_segment_pattern: string
    audio_hls_key_object: string
    audio_hls_iv_hex: string
    audio_hls_segment_duration: number
    index_json: string
    create_time: string
  }
  sections: {
    id: number
    sort_order: number
    topic_title: string
    topic_subtitle: string
    topic_type: string
    part_code: number
    segment_code: string
    topic_items_file: string
    topic_json: string
    blocks: {
      id: number
      block_order: number
      group_index: number
      question_description_json: string
      questions: {
        id: number
        sort_in_block: number
        question_no: number
        score: number
        is_example: number
        content_type: string
        audio_file: string
        stem_text: string
        screen_text_json: string
        analysis_json: string
        question_description_json: string
        raw_json: string
        options: {
          id: number
          flag: string
          sort_order: number
          is_correct: number
          option_type: string
          content: string
        }[]
      }[]
    }[]
  }[]
}

export function getExamPaperDetail(id: number) {
  return request.get<any, { data: ExamPaperDetail }>(`/admin/exam/paper/${id}`)
}

export function importExamPaper(data: {
  mock_examination_paper_id: number
  index_url?: string
  index_json?: string
  level?: string
  paper_id?: string
  source_base_url?: string
  audio_hls_prefix?: string
  conflict_mode?: string
  new_paper_id?: string
}) {
  return request.post<
    any,
    {
      data: {
        examination_paper_id: number
        conflict: boolean
        existing_examination_paper_id: number
        section_count: number
        question_count: number
      }
    }
  >('/admin/exam/paper/import', data, { timeout: 120000 })
}

export function updateExamPaper(data: {
  /** mock_examination_paper.id */
  id: number
  audio_hls_prefix?: string
  audio_hls_segment_count: number
  audio_hls_segment_pattern?: string
  audio_hls_key_object?: string
  audio_hls_iv_hex?: string
  audio_hls_segment_duration?: number
}) {
  return request.post<any, { data: Record<string, never> }>('/admin/exam/paper/update', data)
}

/** --- 考试批次（/admin/exam/batch） --- */

export interface ExamBatchListItem {
  id: number
  mock_examination_paper_id: number
  title: string
  exam_start_at: string
  exam_end_at: string
  mock_level_ids: number[]
  member_count: number
  create_time: string
}

export function getExamBatchList(params: {
  mock_examination_paper_id?: number
  page?: number
  size?: number
}) {
  return request.get<any, { data: { list: ExamBatchListItem[]; total: number } }>('/admin/exam/batch/list', {
    params,
  })
}

export function getExamBatchDetail(id: number) {
  return request.get<any, { data: { batch: ExamBatchListItem } }>(`/admin/exam/batch/${id}`)
}

export function createExamBatch(data: {
  mock_examination_paper_id: number
  title?: string
  exam_start_at: string
  exam_end_at: string
  mock_level_ids: number[]
}) {
  return request.post<any, { data: { id: number } }>('/admin/exam/batch', data)
}

export function updateExamBatch(
  id: number,
  data: { title?: string; exam_start_at: string; exam_end_at: string; mock_level_ids: number[] }
) {
  return request.put<any, { data: Record<string, never> }>(`/admin/exam/batch/${id}`, data)
}

export function deleteExamBatch(id: number) {
  return request.delete<any, { data: Record<string, never> }>(`/admin/exam/batch/${id}`)
}

export interface ExamBatchMemberItem {
  member_id: number
  username: string
  nickname: string
  import_time: string
}

export function importExamBatchMembers(batchId: number, member_ids: number[]) {
  return request.post<any, { data: { inserted: number } }>(
    `/admin/exam/batch/${batchId}/members/import`,
    { member_ids }
  )
}

export function getExamBatchMemberList(batchId: number, params?: { page?: number; size?: number }) {
  return request.get<any, { data: { list: ExamBatchMemberItem[]; total: number } }>(
    `/admin/exam/batch/${batchId}/members/list`,
    { params }
  )
}

export function removeExamBatchMembers(batchId: number, member_ids: number[]) {
  return request.post<any, { data: { removed: number } }>(
    `/admin/exam/batch/${batchId}/members/remove`,
    { member_ids }
  )
}
