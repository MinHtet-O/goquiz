package postgres

import (
	"context"
	"github.com/lib/pq"
	"goquiz/pkg/model"
	"time"
)

func (m QuestionsModel) Insert(categID int, q model.Question) error {
	args := []interface{}{categID, q.WebIndex, q.Text, pq.Array(q.AnsOptions), q.Codeblock, q.Answer.Option, q.Answer.Explanation, q.URL}
	// TODO: replace query with sql prepared statement
	query := `INSERT INTO questions (category_id,web_index,text,ans_options,code_block,correct_ans_opt,correct_ans_explanation,url) values ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}
	return nil
}

func (m QuestionsModel) GetAll(categId int) ([]model.Question, error) {

	query := `select id, text, code_block, ans_options, correct_ans_opt, correct_ans_explanation from questions where category_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, categId)

	if err != nil {
		return nil, err
	}

	var questions []model.Question

	for rows.Next() {
		var question model.Question
		err := rows.Scan(
			&question.ID,
			&question.Text,
			&question.Codeblock,
			pq.Array(&question.AnsOptions),
			&question.Answer.Option,
			&question.Answer.Explanation,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	// TODO: why do we need to close rows??
	defer rows.Close()
	return questions, nil
}
