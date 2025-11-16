package handlers

import (
	"context"
)

// DeleteBatch обработчик запросов на удаление ссылок.
func (h *Handler) DeleteBatch(ctx context.Context) {
	for {
		select {
		case r := <-h.deleteCh:
			err := h.storage.DeleteBatch(ctx, r.Urls, r.UserID)
			if err != nil {
				h.log.Info("Не удалось удалить пачку", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
