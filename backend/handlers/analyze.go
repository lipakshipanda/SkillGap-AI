package handlers

import (
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"

	"skillgap-ai/services"
)

type AnalyzeHandler struct {
	Store *services.Store
}

func NewAnalyzeHandler(store *services.Store) *AnalyzeHandler {
	return &AnalyzeHandler{Store: store}
}

// POST /api/analyze
// Accepts multipart form data with either resume_text or resume_file (PDF),
// plus job_description and optional target_role. Mirrors the FastAPI
// /api/analyze route in routers/analysis.py.
func (h *AnalyzeHandler) Analyze(c *fiber.Ctx) error {
	jobDescription := strings.TrimSpace(c.FormValue("job_description"))
	targetRole := strings.TrimSpace(c.FormValue("target_role"))
	resumeText := strings.TrimSpace(c.FormValue("resume_text"))

	var resolvedResume string

	fileHeader, ferr := c.FormFile("resume_file")
	if ferr == nil && fileHeader != nil {
		if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
			return fiber.NewError(fiber.StatusBadRequest, "Only PDF files are supported")
		}
		file, err := fileHeader.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Could not read uploaded file")
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Could not read uploaded file")
		}

		text, err := services.ExtractTextFromPDF(content)
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		resolvedResume = text
	} else if resumeText != "" {
		resolvedResume = resumeText
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "Provide resume_text or upload a resume_file")
	}

	if len(resolvedResume) < 50 {
		return fiber.NewError(fiber.StatusBadRequest, "Resume text too short")
	}
	if len(jobDescription) < 50 {
		return fiber.NewError(fiber.StatusBadRequest, "Job description too short")
	}

	result, err := services.AnalyzeResume(resolvedResume, jobDescription, targetRole)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, "AI analysis failed: "+err.Error())
	}

	roleForStorage := targetRole
	if roleForStorage == "" {
		roleForStorage = "Not specified"
	}

	record, err := h.Store.SaveAnalysis(*result, services.GetResumeSnippet(resolvedResume, 120), roleForStorage)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save analysis")
	}

	return c.JSON(record.ToResponse())
}
