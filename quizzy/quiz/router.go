package quiz

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"quizzy.app/backend/quizzy/middlewares"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	secured := rt.Group("/quiz", middlewares.RequireAuth, provideStore)
	secured.GET("", handleGetAllUserQuiz)
	secured.POST("", handlePostQuiz)

	quiz := secured.Group("/:quiz-id", provideQuiz)
	quiz.GET("", handleGetQuiz)
	quiz.PATCH("", handlePatchQuiz)
	quiz.GET("/questions", handleGetQuestions)
<<<<<<< HEAD
	quiz.POST("/questions", handlePostQuestion)
	quiz.PUT("/questions/:question-id", provideQuestion, handlePutQuestion)
	quiz.POST("/start", provideCodeResolver, handleStartQuiz)
=======
	quiz.POST("/questions", handlePostQuiz)
>>>>>>> parent of f80ae81 (feat: issues 9, 10, 11 done)
}

func handleGetQuiz(ctx *gin.Context) {
	quiz := useQuiz(ctx)
	ctx.JSON(http.StatusOK, quiz)
}

type QuizzesResponse struct {
<<<<<<< HEAD
	Data  []Quiz `json:"data"`
	Links Links  `json:"_links"`
=======
	Data []Document `json:"data"`
>>>>>>> parent of f80ae81 (feat: issues 9, 10, 11 done)
}

func handleGetAllUserQuiz(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)
	store := useStore(ctx)

	if quizzes, err := store.GetQuizzes(id.Uid); err == nil {
		ctx.JSON(http.StatusOK, QuizzesResponse{
			Data: quizzes,
			Links: Links{
				Create: "/api/quiz",
			},
		})
		return
	}

	ctx.AbortWithStatus(http.StatusInternalServerError)
}

type CreateQuizRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func handlePostQuiz(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)
	store := useStore(ctx)

	// Getting payload from request.
	var req CreateQuizRequest
	if ctx.ShouldBindJSON(&req) != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if code, err := GenerateCode(); err == nil {
		quiz := Quiz{
			Id:          uuid.New().String(),
			Title:       req.Title,
			Description: req.Description,
			Code:        code,
		}

<<<<<<< HEAD
		if err2 := store.Upsert(id.Uid, quiz); err2 == nil {
			ctx.Header("Location", fmt.Sprintf("/api/quiz/%s", quiz.Id))
			ctx.JSON(http.StatusCreated, quiz)
			return
		}
=======
	quiz := Document{
		Uid:         uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
	}

	if err := store.Upsert(id.Uid, quiz); err == nil {
		ctx.Header("Location", ctx.Request.URL.JoinPath(quiz.Uid).RawPath)
		ctx.JSON(http.StatusCreated, quiz)
		return
>>>>>>> parent of f80ae81 (feat: issues 9, 10, 11 done)
	}

	// WARNING / WARNING / WARNING //
	// If this code-path is reached, it means that the requested user has never
	// been registered in our firestore.
	// WARNING / WARNING / WARNING //

	// This may happen if client does some weird things...
	// We should never let the client decide about this process,
	// user registration must be done in single request (Client->Server), or we must use pub/sub
	// from firebase to firestore to store user automatically to avoid data consistency issues.
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

type PatchQuizRequest []FieldPatchOp

func handlePatchQuiz(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)
	store := useStore(ctx)
	quiz := useQuiz(ctx)

	var req PatchQuizRequest
	if ctx.ShouldBindJSON(&req) != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := store.Patch(id.Uid, quiz.Uid, req); err == nil {
		ctx.Status(http.StatusNoContent)
	} else if errors.Is(err, ErrInvalidPatchOperator) || errors.Is(err, ErrInvalidPatchField) {
		ctx.AbortWithStatus(http.StatusBadRequest)
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

type CreateQuestionRequest struct {
	Title   string   `json:"title"`
	Answers []Answer `json:"answers"`
}

func handlePostQuestion(ctx *gin.Context) {
	store := useStore(ctx)
	id := middlewares.UseIdentity(ctx)
	quiz := useQuiz(ctx)

	var req CreateQuestionRequest
	if ctx.ShouldBindJSON(&req) != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Appending question.
	quiz.Questions = append(quiz.Questions, Question{
		Title:   req.Title,
		Uid:     uuid.New().String(),
		Answers: req.Answers,
	})

	if store.Upsert(id.Uid, quiz) != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}

func handleGetQuestions(ctx *gin.Context) {
	quiz := useQuiz(ctx)
	ctx.JSON(http.StatusOK, quiz.Questions)
}
<<<<<<< HEAD

type UnidentifiedAnswer struct {
	Title     string `json:"title"`
	IsCorrect bool   `json:"isCorrect"`
}

type UpdateQuestionRequest struct {
	Title   string               `json:"title"`
	Answers []UnidentifiedAnswer `json:"answers"`
}

func handlePutQuestion(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)
	store := useStore(ctx)
	quiz := useQuiz(ctx)
	question := useQuestion(ctx)

	var payload UpdateQuestionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	question.Title = payload.Title
	question.Answers = []Answer{}
	for _, a := range payload.Answers {
		question.Answers = append(question.Answers, Answer{
			Id:        uuid.New().String(),
			Title:     a.Title,
			IsCorrect: a.IsCorrect,
		})
	}

	if err := store.UpdateQuestion(id.Uid, quiz.Id, question); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func handleStartQuiz(ctx *gin.Context) {
	resolver := useCodeResolver(ctx)
	quiz := useQuiz(ctx)
	if !canStart(&quiz) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := resolver.BindCode(quiz); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Header("Location", fmt.Sprintf("/api/execution/%s", quiz.Code))
	ctx.Status(http.StatusCreated)
}
=======
>>>>>>> parent of f80ae81 (feat: issues 9, 10, 11 done)
