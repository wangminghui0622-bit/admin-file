// Package cloud maintains the in-memory S3 client group of the configured
// storage providers. The group is rebuilt after every provider mutation so
// that changes take effect immediately (mirrors the reference implementation).
package cloud

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// ProviderData is an immutable snapshot of a storage provider row.
type ProviderData struct {
	ID        uint64
	Name      string
	Folder    string
	Bucket    string
	Endpoint  string
	Region    string
	SecretID  string
	SecretKey string
	UseCdn    bool
	CdnURL    string
	IsDefault bool
}

// Group holds the S3 clients keyed by provider name.
type Group struct {
	db  *gorm.DB
	q   *query.Query
	mu  sync.RWMutex
	s3s map[string]*s3.S3
	// data mirrors s3s so callers can read configuration without a DB hit.
	data            map[string]ProviderData
	defaultProvider string
}

// NewGroup loads all enabled providers and builds their S3 clients.
func NewGroup(db *gorm.DB, q *query.Query) *Group {
	g := &Group{db: db, q: q}
	g.Reload()
	return g
}

// Reload rebuilds the client map from fms_storage_providers where state = true.
func (g *Group) Reload() {
	clients := make(map[string]*s3.S3)
	data := make(map[string]ProviderData)
	defaultProvider := ""

	rows, err := g.q.FmsStorageProvider.WithContext(context.Background()).
		Where(g.q.FmsStorageProvider.State.Eq(true)).Find()
	if err != nil {
		logx.Errorw("failed to load storage providers, make sure the database is initialized",
			logx.Field("detail", err.Error()))
	} else {
		for _, row := range rows {
			sess, err := session.NewSession(&aws.Config{
				Region:      aws.String(row.Region),
				Credentials: credentials.NewStaticCredentials(row.SecretID, row.SecretKey, ""),
				Endpoint:    aws.String(row.Endpoint),
			})
			if err != nil {
				logx.Errorw("failed to create s3 session", logx.Field("provider", row.Name),
					logx.Field("detail", err.Error()))
				continue
			}
			pd := ProviderData{
				ID:        uint64(row.ID),
				Name:      row.Name,
				Bucket:    row.Bucket,
				Endpoint:  row.Endpoint,
				Region:    row.Region,
				SecretID:  row.SecretID,
				SecretKey: row.SecretKey,
				UseCdn:    row.UseCdn,
				IsDefault: row.IsDefault,
			}
			if row.Folder != nil {
				pd.Folder = *row.Folder
			}
			if row.CdnURL != nil {
				pd.CdnURL = *row.CdnURL
			}
			clients[row.Name] = s3.New(sess)
			data[row.Name] = pd
			if row.IsDefault {
				defaultProvider = row.Name
			}
		}
	}

	g.mu.Lock()
	g.s3s = clients
	g.data = data
	g.defaultProvider = defaultProvider
	g.mu.Unlock()
}

// Client returns the S3 client and snapshot for the named provider.
func (g *Group) Client(name string) (*s3.S3, ProviderData, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	c, ok := g.s3s[name]
	if !ok {
		return nil, ProviderData{}, false
	}
	return c, g.data[name], true
}

// Provider returns the configuration snapshot for the named provider.
func (g *Group) Provider(name string) (ProviderData, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	p, ok := g.data[name]
	return p, ok
}

// Default returns the provider flagged as default.
func (g *Group) Default() (ProviderData, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.defaultProvider == "" {
		return ProviderData{}, false
	}
	p, ok := g.data[g.defaultProvider]
	return p, ok
}
