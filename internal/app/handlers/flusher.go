package handlers

import (
	"context"
)

func (h *Handler) DeleteBatch(ctx context.Context) {
	for r := range h.deleteCh {
		err := h.storage.DeleteBatch(ctx, r.Urls, r.UserID)
		if err != nil {
			h.log.Info("Не удалось удалить пачку", err)
		}
	}
}
