package controllers

import (
	"context"
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Question struct {
	ID           uuid.UUID   `json:"id"`
	Description  *string     `json:"description"`
	Title        *string     `json:"title"`
	InputFormat  *string     `json:"input_format"`
	Points       pgtype.Int4 `json:"points"`
	Round        int32       `json:"round"`
	Constraints  *string     `json:"constraints"`
	OutputFormat *string     `json:"output_format"`
}

func GetAllQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fetchedQuestions, err := database.Queries.GetQuestions(ctx)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httphelpers.WriteJSON(w, 200, fetchedQuestions)
}

func GetQuestionById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "question_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	fetchedQuestion, err := database.Queries.GetQuestion(ctx, id)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httphelpers.WriteJSON(w, 200, fetchedQuestion)
}

func GetQuestionsByRound(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := auth.GetUserID(w, r)

	user, err := database.Queries.GetUserById(ctx, id)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
	}

	questions, err := database.Queries.GetQuestionByRound(ctx, user.RoundQualified)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httphelpers.WriteJSON(w, 200, questions)
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var question Question
	err := httphelpers.ParseJSON(r, &question)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	questions, err := database.Queries.CreateQuestion(ctx, db.CreateQuestionParams{
		ID:           uuid.New(),
		Description:  question.Description,
		Title:        question.Title,
		InputFormat:  question.InputFormat,
		Points:       question.Points,
		Round:        question.Round,
		Constraints:  question.Constraints,
		OutputFormat: question.OutputFormat,
	})
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httphelpers.WriteJSON(w, 200, questions)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	questionIdStr := chi.URLParam(r, "question_id")
	question_id, err := uuid.Parse(questionIdStr)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = database.Queries.DeleteQuestion(ctx, question_id)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httphelpers.WriteJSON(w, 200, map[string]string{"message": "Question successfully deleted"})
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var updateQuestion Question
	if err := httphelpers.ParseJSON(r, &updateQuestion); err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	question, err := database.Queries.GetQuestion(ctx, updateQuestion.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			httphelpers.WriteError(w, http.StatusNotFound, err.Error())
		} else {
			httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		}
	}
	nulVal := pgtype.Int4{
		Int32: 0,
		Valid: false,
	}

	if updateQuestion.Description != nil {
		question.Description = updateQuestion.Description
	}
	if updateQuestion.Title != nil {
		question.Title = updateQuestion.Title
	}
	if updateQuestion.InputFormat != nil {
		question.InputFormat = updateQuestion.InputFormat
	}
	if updateQuestion.Points != nulVal {
		question.Points = updateQuestion.Points
	}
	if updateQuestion.Round != 0 {
		question.Round = updateQuestion.Round
	}
	if updateQuestion.Constraints != nil {
		question.Constraints = updateQuestion.Constraints
	}
	if updateQuestion.OutputFormat != nil {
		question.OutputFormat = updateQuestion.OutputFormat
	}

	err = database.Queries.UpdateQuestion(ctx, db.UpdateQuestionParams{
		Description:  question.Description,
		Title:        question.Title,
		InputFormat:  question.InputFormat,
		Points:       question.Points,
		Round:        question.Round,
		Constraints:  question.Constraints,
		OutputFormat: question.OutputFormat,
		ID:           updateQuestion.ID,
	})
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
	}
	httphelpers.WriteJSON(w, http.StatusOK, question)
}
