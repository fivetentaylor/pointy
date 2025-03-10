package prompts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/fivetentaylor/pointy/pkg/client"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/teamreviso/freeplay"
	"gorm.io/gorm"
)

func Refresh(ctx context.Context) error {
	f, err := client.NewFreeplayClientFromEnv()
	if err != nil {
		return err
	}

	q := env.Query(ctx)

	ctx, cancel := context.WithCancel(context.Background())
	refresher := &refresher{
		ctx:    ctx,
		client: f,
		cancel: cancel,
		logger: slog.Default(),
		query:  q,
	}

	return refresher.Refresh()
}

type refresher struct {
	ctx      context.Context
	duration time.Duration
	client   *client.Freeplay
	cancel   context.CancelFunc
	logger   *slog.Logger
	query    *query.Query
}

func (p *refresher) Poll() {
	for {
		select {
		case <-time.After(p.duration):
			err := p.Refresh()
			if err != nil {
				p.logger.Error("[prompts] error refreshing prompts", slog.Any("error", err))
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *refresher) Refresh() error {
	defer func() {
		if r := recover(); r != nil {
			p.logger.Error("[prompts] PANIC refreshing prompts", slog.Any("panic", r), slog.Any("stack", string(debug.Stack())))
		}
	}()

	prompts, err := p.client.All()
	if err != nil {
		p.logger.Error("[prompts] error getting all prompts", slog.Any("error", err))
		return err
	}

	for _, prompt := range prompts {
		p.logger.Info(
			"refreshing prompt",
			slog.Any("prompt", prompt.PromptTemplateName),
			slog.Any("content", prompt.Content),
			slog.Any("system", prompt.SystemContent),
		)
		err = p.UpsertPrompt(prompt)
		if err != nil {
			p.logger.Error("[prompts] error upserting prompt", slog.Any("error", err), slog.Any("prompt", prompt))
			return err
		}
	}

	return nil
}

func (p *refresher) UpsertPrompt(prompt freeplay.Prompt) error {
	tbl := p.query.Prompt

	record, err := tbl.Where(tbl.PromptName.Eq(prompt.PromptTemplateName)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if record == nil {
		record = &models.Prompt{
			PromptName: prompt.PromptTemplateName,
		}

		p.logger.Info("creating prompt", slog.Any("record", record))
	} else {
		p.logger.Info("updating prompt", slog.Any("record", record))
	}

	record.Version = prompt.PromptTemplateVersionID
	if prompt.SystemContent != nil {
		record.SystemContent = prompt.SystemContent
	}
	if prompt.SystemContent == nil {
		for _, content := range prompt.Content {
			if content.Role == "system" {
				record.SystemContent = &content.Content
				break
			}
		}
	}

	bts, err := json.Marshal(prompt.Content)
	if err != nil {
		return fmt.Errorf("error marshaling prompt content: %w", err)
	}

	record.ContentJSON = string(bts)

	record.Provider = prompt.Metadata.Provider
	record.ModelName = prompt.Metadata.Model
	record.Temperature = prompt.Metadata.Params.Temperature
	record.MaxTokens = int32(prompt.Metadata.Params.MaxTokens)
	record.TopP = prompt.Metadata.Params.TopP

	err = tbl.Save(record)
	if err != nil {
		return err
	}

	p.logger.Info("saved prompt", slog.Any("record", record))
	return nil
}
