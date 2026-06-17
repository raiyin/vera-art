package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

// MiscHandler handles miscellaneous HTTP requests (admin, tags, materials, bases, consent, master classes).
type MiscHandler struct {
	adminService       port.AdminService
	tagService         port.TagService
	materialService    port.MaterialService
	baseService        port.BaseService
	consentService     port.ConsentService
	masterClassService port.MasterClassService
}

// NewMiscHandler creates a new MiscHandler.
func NewMiscHandler(
	adminService port.AdminService,
	tagService port.TagService,
	materialService port.MaterialService,
	baseService port.BaseService,
	consentService port.ConsentService,
	masterClassService port.MasterClassService,
) *MiscHandler {
	return &MiscHandler{
		adminService:       adminService,
		tagService:         tagService,
		materialService:    materialService,
		baseService:        baseService,
		consentService:     consentService,
		masterClassService: masterClassService,
	}
}

// ---------------------------------------------------------------------------
// Admin Dashboard
// ---------------------------------------------------------------------------

// GetDashboardStats returns admin dashboard statistics.
func (h *MiscHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.adminService.GetDashboardStats(c.Request.Context())
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ---------------------------------------------------------------------------
// Tags
// ---------------------------------------------------------------------------

// GetTags returns a list of all tags.
func (h *MiscHandler) GetTags(c *gin.Context) {
	tags, err := h.tagService.GetTags(c.Request.Context())
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.TagResponse, len(tags))
	for i, t := range tags {
		responses[i] = dto.TagResponse{
			ID:        t.ID,
			NameRu:    t.NameRu,
			NameEn:    t.NameEn,
			Slug:      t.Slug,
			CreatedAt: t.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"tags": responses})
}

// CreateTag creates a new tag.
func (h *MiscHandler) CreateTag(c *gin.Context) {
	var req dto.CreateTagRequest
	if !BindJSON(c, &req) {
		return
	}

	tag := &domain.Tag{
		NameRu: req.NameRu,
		NameEn: req.NameEn,
		Slug:   req.Slug,
	}

	if err := h.tagService.CreateTag(c.Request.Context(), tag); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusCreated, dto.TagResponse{
		ID:        tag.ID,
		NameRu:    tag.NameRu,
		NameEn:    tag.NameEn,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
	})
}

// UpdateTag updates a tag.
func (h *MiscHandler) UpdateTag(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateTagRequest
	if !BindJSON(c, &req) {
		return
	}

	tag := &domain.Tag{
		ID:     id,
		NameRu: req.NameRu,
		NameEn: req.NameEn,
		Slug:   req.Slug,
	}

	if err := h.tagService.CreateTag(c.Request.Context(), tag); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag updated successfully"})
}

// DeleteTag deletes a tag.
func (h *MiscHandler) DeleteTag(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.tagService.DeleteTag(c.Request.Context(), id); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

// ---------------------------------------------------------------------------
// Materials
// ---------------------------------------------------------------------------

// GetMaterials returns a list of all materials.
func (h *MiscHandler) GetMaterials(c *gin.Context) {
	materials, err := h.materialService.GetMaterials(c.Request.Context())
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.MaterialResponse, len(materials))
	for i, m := range materials {
		responses[i] = dto.MaterialResponse{
			ID:        m.ID,
			NameRu:    m.NameRu,
			NameEn:    m.NameEn,
			CreatedAt: m.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"materials": responses})
}

// CreateMaterial creates a new material.
func (h *MiscHandler) CreateMaterial(c *gin.Context) {
	var req dto.CreateMaterialRequest
	if !BindJSON(c, &req) {
		return
	}

	material := &domain.Material{
		NameRu: req.NameRu,
		NameEn: req.NameEn,
	}

	if err := h.materialService.CreateMaterial(c.Request.Context(), material); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusCreated, dto.MaterialResponse{
		ID:        material.ID,
		NameRu:    material.NameRu,
		NameEn:    material.NameEn,
		CreatedAt: material.CreatedAt,
	})
}

// DeleteMaterial deletes a material.
func (h *MiscHandler) DeleteMaterial(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.materialService.DeleteMaterial(c.Request.Context(), id); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material deleted successfully"})
}

// ---------------------------------------------------------------------------
// Bases
// ---------------------------------------------------------------------------

// GetBases returns a list of all bases.
func (h *MiscHandler) GetBases(c *gin.Context) {
	bases, err := h.baseService.GetBases(c.Request.Context())
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.BaseResponse, len(bases))
	for i, b := range bases {
		responses[i] = dto.BaseResponse{
			ID:        b.ID,
			NameRu:    b.NameRu,
			NameEn:    b.NameEn,
			CreatedAt: b.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"bases": responses})
}

// CreateBase creates a new base.
func (h *MiscHandler) CreateBase(c *gin.Context) {
	var req dto.CreateBaseRequest
	if !BindJSON(c, &req) {
		return
	}

	base := &domain.Base{
		NameRu: req.NameRu,
		NameEn: req.NameEn,
	}

	if err := h.baseService.CreateBase(c.Request.Context(), base); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusCreated, dto.BaseResponse{
		ID:        base.ID,
		NameRu:    base.NameRu,
		NameEn:    base.NameEn,
		CreatedAt: base.CreatedAt,
	})
}

// DeleteBase deletes a base.
func (h *MiscHandler) DeleteBase(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.baseService.DeleteBase(c.Request.Context(), id); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Base deleted successfully"})
}

// ---------------------------------------------------------------------------
// Consent
// ---------------------------------------------------------------------------

// SaveConsent records user consent.
func (h *MiscHandler) SaveConsent(c *gin.Context) {
	var req dto.ConsentRequest
	if !BindJSON(c, &req) {
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = int64(req.UserID)
	}

	ipAddress := req.IPAddress
	if ipAddress == "" {
		ipAddress = c.ClientIP()
	}

	// Record analytics consent
	if err := h.consentService.RecordConsent(c.Request.Context(), userID, "analytics", req.AnalyticsConsent, ipAddress); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	// Record marketing consent
	if err := h.consentService.RecordConsent(c.Request.Context(), userID, "marketing", req.MarketingConsent, ipAddress); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.ConsentResponse{
		Success:          true,
		Message:          "Consent saved successfully",
		ConsentGiven:     true,
		ConsentVersion:   "1.0",
		AnalyticsAllowed: req.AnalyticsConsent,
		MarketingAllowed: req.MarketingConsent,
	})
}

// GetConsentStatus returns the current consent status for the user.
func (h *MiscHandler) GetConsentStatus(c *gin.Context) {
	userID := c.GetInt64("user_id")

	consents, err := h.consentService.GetUserConsents(c.Request.Context(), userID)
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	analyticsAllowed := false
	marketingAllowed := false

	for _, consent := range consents {
		if consent.ConsentType == "analytics" && consent.Granted {
			analyticsAllowed = true
		}
		if consent.ConsentType == "marketing" && consent.Granted {
			marketingAllowed = true
		}
	}

	c.JSON(http.StatusOK, dto.ConsentResponse{
		Success:          true,
		ConsentGiven:     len(consents) > 0,
		ConsentVersion:   "1.0",
		AnalyticsAllowed: analyticsAllowed,
		MarketingAllowed: marketingAllowed,
	})
}

// ---------------------------------------------------------------------------
// Master Classes
// ---------------------------------------------------------------------------

// GetMasterClasses returns a list of published master classes.
func (h *MiscHandler) GetMasterClasses(c *gin.Context) {
	masterClasses, total, err := h.masterClassService.GetMasterClasses(c.Request.Context())
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.MasterClassResponse, len(masterClasses))
	for i, mc := range masterClasses {
		responses[i] = dto.MasterClassResponse{
			ID:          mc.ID,
			Title:       mc.Title,
			Description: mc.Description,
			Price:       mc.Price,
			ImagePath:   mc.ImagePath,
			VideoURL:    mc.VideoURL,
			Status:      mc.Status,
			CreatedAt:   mc.CreatedAt,
			UpdatedAt:   mc.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"master_classes": responses,
		"total":          total,
	})
}

// GetMasterClassByID returns a master class by ID.
func (h *MiscHandler) GetMasterClassByID(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	mc, err := h.masterClassService.GetMasterClassByID(c.Request.Context(), id)
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.MasterClassResponse{
		ID:          mc.ID,
		Title:       mc.Title,
		Description: mc.Description,
		Price:       mc.Price,
		ImagePath:   mc.ImagePath,
		VideoURL:    mc.VideoURL,
		Status:      mc.Status,
		CreatedAt:   mc.CreatedAt,
		UpdatedAt:   mc.UpdatedAt,
	})
}

// CreateMasterClass creates a new master class.
func (h *MiscHandler) CreateMasterClass(c *gin.Context) {
	var req dto.CreateMasterClassRequest
	if !BindJSON(c, &req) {
		return
	}

	mc := &domain.MasterClass{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		ImagePath:   req.ImagePath,
		VideoURL:    req.VideoURL,
		Status:      req.Status,
	}

	if err := h.masterClassService.CreateMasterClass(c.Request.Context(), mc); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusCreated, dto.MasterClassResponse{
		ID:          mc.ID,
		Title:       mc.Title,
		Description: mc.Description,
		Price:       mc.Price,
		ImagePath:   mc.ImagePath,
		VideoURL:    mc.VideoURL,
		Status:      mc.Status,
		CreatedAt:   mc.CreatedAt,
		UpdatedAt:   mc.UpdatedAt,
	})
}

// UpdateMasterClass updates a master class.
func (h *MiscHandler) UpdateMasterClass(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateMasterClassRequest
	if !BindJSON(c, &req) {
		return
	}

	mc := &domain.MasterClass{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		ImagePath:   req.ImagePath,
		VideoURL:    req.VideoURL,
		Status:      req.Status,
	}

	if err := h.masterClassService.UpdateMasterClass(c.Request.Context(), mc); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.MasterClassResponse{
		ID:          mc.ID,
		Title:       mc.Title,
		Description: mc.Description,
		Price:       mc.Price,
		ImagePath:   mc.ImagePath,
		VideoURL:    mc.VideoURL,
		Status:      mc.Status,
		CreatedAt:   mc.CreatedAt,
		UpdatedAt:   mc.UpdatedAt,
	})
}

// DeleteMasterClass deletes a master class.
func (h *MiscHandler) DeleteMasterClass(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.masterClassService.DeleteMasterClass(c.Request.Context(), id); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Master class deleted successfully"})
}

// Ensure unused imports are not flagged.
var _ = strconv.Itoa
