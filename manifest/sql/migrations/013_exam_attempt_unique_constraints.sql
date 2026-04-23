-- 013: attempt reliability constraints.
-- Only active rows (delete_flag = 0) participate in uniqueness. Deleted rows
-- are ignored through a generated NULL flag, so soft-delete history is not
-- blocked by the unique indexes.

-- 1) Keep one active answer row per attempt/question before adding the index.
-- Preference: manually graded row, higher version, newer update_time, larger id.
UPDATE `exam_attempt_answer` loser
INNER JOIN `exam_attempt_answer` keeper
    ON keeper.`attempt_id` = loser.`attempt_id`
    AND keeper.`exam_question_id` = loser.`exam_question_id`
    AND keeper.`delete_flag` = 0
    AND loser.`delete_flag` = 0
    AND keeper.`id` <> loser.`id`
    AND (
        (keeper.`awarded_score` IS NOT NULL AND loser.`awarded_score` IS NULL)
        OR (
            (keeper.`awarded_score` IS NOT NULL) = (loser.`awarded_score` IS NOT NULL)
            AND keeper.`version` > loser.`version`
        )
        OR (
            (keeper.`awarded_score` IS NOT NULL) = (loser.`awarded_score` IS NOT NULL)
            AND keeper.`version` = loser.`version`
            AND COALESCE(keeper.`update_time`, keeper.`create_time`, '1970-01-01') > COALESCE(loser.`update_time`, loser.`create_time`, '1970-01-01')
        )
        OR (
            (keeper.`awarded_score` IS NOT NULL) = (loser.`awarded_score` IS NOT NULL)
            AND keeper.`version` = loser.`version`
            AND COALESCE(keeper.`update_time`, keeper.`create_time`, '1970-01-01') = COALESCE(loser.`update_time`, loser.`create_time`, '1970-01-01')
            AND keeper.`id` > loser.`id`
        )
    )
SET loser.`delete_flag` = 1,
    loser.`updater` = 'migration_013',
    loser.`update_time` = NOW(3);

ALTER TABLE `exam_attempt_answer`
    ADD COLUMN `active_unique_flag` tinyint GENERATED ALWAYS AS (
        CASE WHEN `delete_flag` = 0 THEN 1 ELSE NULL END
    ) STORED;

CREATE UNIQUE INDEX `uk_exam_attempt_answer_active_question`
    ON `exam_attempt_answer` (`attempt_id`, `exam_question_id`, `active_unique_flag`);

-- 2) Keep one active batch attempt per member/batch/paper before adding the index.
-- Historical non-batch attempts (exam_batch_id = 0) are not constrained here.
-- Preference: most advanced status, newer update_time, larger id.
UPDATE `exam_attempt` loser
INNER JOIN `exam_attempt` keeper
    ON keeper.`member_id` = loser.`member_id`
    AND keeper.`exam_batch_id` = loser.`exam_batch_id`
    AND keeper.`exam_paper_id` = loser.`exam_paper_id`
    AND keeper.`delete_flag` = 0
    AND loser.`delete_flag` = 0
    AND keeper.`exam_batch_id` > 0
    AND loser.`exam_batch_id` > 0
    AND keeper.`id` <> loser.`id`
    AND (
        keeper.`status` > loser.`status`
        OR (
            keeper.`status` = loser.`status`
            AND COALESCE(keeper.`update_time`, keeper.`create_time`, '1970-01-01') > COALESCE(loser.`update_time`, loser.`create_time`, '1970-01-01')
        )
        OR (
            keeper.`status` = loser.`status`
            AND COALESCE(keeper.`update_time`, keeper.`create_time`, '1970-01-01') = COALESCE(loser.`update_time`, loser.`create_time`, '1970-01-01')
            AND keeper.`id` > loser.`id`
        )
    )
SET loser.`delete_flag` = 1,
    loser.`updater` = 'migration_013',
    loser.`update_time` = NOW(3);

ALTER TABLE `exam_attempt`
    ADD COLUMN `active_unique_flag` tinyint GENERATED ALWAYS AS (
        CASE WHEN `delete_flag` = 0 AND `exam_batch_id` > 0 THEN 1 ELSE NULL END
    ) STORED;

CREATE UNIQUE INDEX `uk_exam_attempt_active_member_batch_paper`
    ON `exam_attempt` (`member_id`, `exam_batch_id`, `exam_paper_id`, `active_unique_flag`);
